// SPDX-FileCopyrightText: 2022 Infosys Limited
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//

package nrf_cache

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"github.com/omec-project/openapi/Nnrf_NFDiscovery"
	"github.com/omec-project/openapi/models"
	"sync"
	"time"
)

func InitCache() int {
	return 1
}

const defaultCacheTTl = time.Hour

type NfProfileItem struct {
	nfProfile  *models.NfProfile
	ttl        time.Duration
	expiryTime time.Time
	index      int // index of the entry in the priority queue
}

func (item *NfProfileItem) isExpired() bool {
	return item.expiryTime.Before(time.Now())
}

func (item *NfProfileItem) updateExpiryTime() {
	item.expiryTime = time.Now().Add(time.Second * item.ttl)
}

func newNfProfileItem(profile *models.NfProfile, ttl time.Duration) *NfProfileItem {
	item := &NfProfileItem{
		nfProfile: profile,
		ttl:       ttl,
	}
	// since nobody is aware yet of this item, it's safe to touch without lock here
	item.updateExpiryTime()
	return item
}

// NfProfilePriorityQ : Priority Queue to store the profile by expiry time
type NfProfilePriorityQ []*NfProfileItem

func (npq NfProfilePriorityQ) Len() int {
	return len(npq)
}

func (npq NfProfilePriorityQ) Less(i, j int) bool {
	return npq[i].expiryTime.Before(npq[j].expiryTime)
}

func (npq NfProfilePriorityQ) Swap(i, j int) {
	npq[i], npq[j] = npq[j], npq[i]
	npq[i].index = i
	npq[j].index = j
}

func (npq NfProfilePriorityQ) root() *NfProfileItem {
	return npq[0]
}

func (npq NfProfilePriorityQ) at(index int) *NfProfileItem {
	return npq[index]
}

func (npq *NfProfilePriorityQ) push(item interface{}) {
	heap.Push(npq, item)
}

func (npq *NfProfilePriorityQ) pop() interface{} {
	if npq.Len() == 0 {
		return nil
	}
	return heap.Pop(npq).(*NfProfileItem)
}

// update modifies the priority and value of an Item in the queue.
func (npq *NfProfilePriorityQ) update(item *NfProfileItem, value *models.NfProfile, ttl time.Duration) {
	item.nfProfile = value
	item.ttl = ttl
	item.updateExpiryTime()
	heap.Fix(npq, item.index)
}

func (npq *NfProfilePriorityQ) remove(item *NfProfileItem) {
	heap.Remove(npq, item.index)
}

func (npq *NfProfilePriorityQ) Push(item interface{}) {
	n := len(*npq)
	entry := item.(*NfProfileItem)
	entry.index = n
	*npq = append(*npq, entry)
}

func (npq *NfProfilePriorityQ) Pop() interface{} {
	old := *npq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*npq = old[0 : n-1]
	return item
}

func newNfProfilePriorityQ() *NfProfilePriorityQ {
	q := &NfProfilePriorityQ{}
	heap.Init(q)
	return q
}

// NrfCache : cache of nf profiles
type NrfCache struct {
	cache map[string]*NfProfileItem // map[nf-instance-id] =*NfProfile

	priorityQ *NfProfilePriorityQ // sorted by expiry time

	evictionInterval time.Duration

	evictionTicker *time.Ticker

	done chan struct{}

	mutex sync.RWMutex
}

func (c *NrfCache) set(nfProfile *models.NfProfile, ttl time.Duration) {

	nfProfileCopy, err := copyNrfProfile(nfProfile)
	if err != nil {
		fmt.Printf("failed to copy nfprofile %s", err)
		return
	}

	c.mutex.Lock()

	item, exists := c.cache[nfProfileCopy.NfInstanceId]
	if exists {
		// if item.isExpired()
		c.priorityQ.update(item, nfProfileCopy, ttl)
	} else {
		newItem := newNfProfileItem(nfProfileCopy, ttl)
		c.cache[nfProfileCopy.NfInstanceId] = newItem
		c.priorityQ.push(newItem)
	}

	c.mutex.Unlock()
}

func (c *NrfCache) get(opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) []models.NfProfile {
	var nfProfiles []models.NfProfile

	c.mutex.RLock()
	for _, element := range c.cache {
		if !element.isExpired() {
			if opts != nil {
				if matchFilters[element.nfProfile.NfType](element.nfProfile, opts) {
					nrfProfile, err := copyNrfProfile(element.nfProfile)
					if err == nil {
						nfProfiles = append(nfProfiles, *nrfProfile)
					}
				}
			} else {
				nrfProfile, err := copyNrfProfile(element.nfProfile)
				if err == nil {
					nfProfiles = append(nfProfiles, *nrfProfile)
				}
			}
		}
	}
	c.mutex.RUnlock()
	return nfProfiles
}

func (c *NrfCache) remove(item *NfProfileItem) {
	c.priorityQ.remove(item)
	delete(c.cache, item.nfProfile.NfInstanceId)
}

func (c *NrfCache) cleanupExpiredItems() {
	fmt.Println("evit items")

	for item := c.priorityQ.at(0); item.isExpired(); {
		fmt.Printf("Item being removed %s", item.nfProfile.NfType)
		c.remove(item)
		if c.priorityQ.Len() == 0 {
			break
		} else {
			item = c.priorityQ.at(0)
		}
	}
}

func (c *NrfCache) purge() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	close(c.done)
	c.priorityQ = newNfProfilePriorityQ()
	c.cache = make(map[string]*NfProfileItem)
	c.evictionTicker.Stop()
}

func (c *NrfCache) startExpiryProcessing() {
	for {
		select {
		case <-c.evictionTicker.C:
			c.mutex.Lock()
			if c.priorityQ.Len() == 0 {
				c.mutex.Unlock()
				continue
			}

			c.cleanupExpiredItems()
			c.mutex.Unlock()

		case <-c.done:
			return
		}
	}
}

func NewNrfCache(duration time.Duration) *NrfCache {
	cache := &NrfCache{
		cache:            make(map[string]*NfProfileItem),
		priorityQ:        newNfProfilePriorityQ(),
		evictionInterval: defaultCacheTTl,
		done:             make(chan struct{}),
	}

	cache.evictionTicker = time.NewTicker(duration * time.Second)

	go cache.startExpiryProcessing()

	return cache
}

func copyNrfProfile(src *models.NfProfile) (*models.NfProfile, error) {
	nrfProfileJSON, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	nrfProfile := models.NfProfile{}
	if err = json.Unmarshal(nrfProfileJSON, &nrfProfile); err != nil {
		return nil, err
	}

	return &nrfProfile, nil
}

type NrfRequest struct {
	targetNfType models.NfType
	searchParams *Nnrf_NFDiscovery.SearchNFInstancesParamOpts
}

type NrfMasterCache struct {
	nfTypeToCacheMap    map[models.NfType]*NrfCache
	evictionInterval    time.Duration
	nrfDiscoveryQueryCb NrfDiscoveryQueryCb

	NrfCommChan chan interface{}

	mutex sync.RWMutex
}

func (c *NrfMasterCache) Get(request NrfRequest) []models.NfProfile {
	c.mutex.RLock()

	var result []models.NfProfile
	cache, exists := c.nfTypeToCacheMap[request.targetNfType]
	if exists {
		result = cache.get(request.searchParams)
	}

	c.mutex.RUnlock()

	return result
}

func (c *NrfMasterCache) Set(profile *models.NfProfile, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	nfCache, exists := c.nfTypeToCacheMap[profile.NfType]
	if !exists {
		nfCache = NewNrfCache(c.evictionInterval)
		if nfCache != nil {
			c.nfTypeToCacheMap[profile.NfType] = nfCache
		}
	}

	if nfCache != nil {
		nfCache.set(profile, duration)
	}

}

var masterCache *NrfMasterCache

type NrfDiscoveryQueryCb func(nrfUri string, targetNfType, requestNfType models.NfType, param *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (models.SearchResult, error)

func InitNrfCaching(interval time.Duration, cb NrfDiscoveryQueryCb) {
	m := &NrfMasterCache{
		nfTypeToCacheMap:    make(map[models.NfType]*NrfCache),
		evictionInterval:    interval,
		nrfDiscoveryQueryCb: cb,

		NrfCommChan: make(chan interface{}), // process request

	}
	masterCache = m
}

func SearchNFInstances(nrfUri string, targetNfType, requestNfType models.NfType,
	param *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (models.SearchResult, error) {
	req := NrfRequest{
		targetNfType: targetNfType,
		searchParams: param,
	}

	var searchResult models.SearchResult
	searchResult.NfInstances = masterCache.Get(req)

	if len(searchResult.NfInstances) > 0 {
		return searchResult, nil
	} else if masterCache.nrfDiscoveryQueryCb != nil {
		searchResult, err := masterCache.nrfDiscoveryQueryCb(nrfUri, targetNfType, requestNfType, param)
		for i := 0; i < len(searchResult.NfInstances); i++ {
			masterCache.Set(&searchResult.NfInstances[i], time.Duration(searchResult.ValidityPeriod))
		}

		return searchResult, err
	}
	return searchResult, nil
}
