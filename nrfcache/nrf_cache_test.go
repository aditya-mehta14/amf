// SPDX-FileCopyrightText: 2022 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//
package nrf_cache

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	amf_context "github.com/omec-project/amf/context"
	"github.com/omec-project/openapi/Nnrf_NFDiscovery"
	"github.com/omec-project/openapi/models"
	"github.com/stretchr/testify/assert"
)

func NoTestNfPriorityQ(t *testing.T) {
	t.Log("Method Entry")

	pq := newNfProfilePriorityQ()

	smfProfileStr := `{ "ipv4Addresses" : [ "smf" ], "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "smfInfo" : { "sNssaiSmfInfoList" : [ { "sNssai" : { "sst" : 1, "sd" : "010203" }, "dnnSmfInfoList" : [ { "dnn" : "internet" } ] } ] }, "nfServices" : [ { "apiPrefix" : "http://smf:29502", "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "serviceInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bfnsmf-pdusession", "serviceName" : "nsmf-pdusession", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "https://smf:29502/nsmf-pdusession/v1", "expiry" : "2022-08-17T05:31:40.997097141Z" } ], "scheme" : "https", "nfServiceStatus" : "REGISTERED" }, { "scheme" : "https", "nfServiceStatus" : "REGISTERED", "apiPrefix" : "http://smf:29502", "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "serviceInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bfnsmf-event-exposure", "serviceName" : "nsmf-event-exposure", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "https://smf:29502/nsmf-pdusession/v1", "expiry" : "2022-08-17T05:31:40.997097141Z" } ] } ], "nfInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bf", "plmnList" : [ { "mnc" : "93", "mcc" : "208" } ], "sNssais" : [ { "sd" : "010203", "sst" : 1 } ], "nfType" : "SMF", "nfStatus" : "REGISTERED" }`
	ausfProfileStr := `{ "nfServices" : [ { "serviceName" : "nausf-auth", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "ausf", "port" : 29509 } ], "serviceInstanceId" : "57d0a167-5283-4170-bdd8-881076049a81" } ], "nfInstanceId" : "57d0a167-5283-4170-bdd8-881076049a81", "nfType" : "AUSF", "nfStatus" : "REGISTERED", "plmnList" : [ { "mcc" : "208", "mnc" : "93" } ], "ipv4Addresses" : [ "ausf" ], "ausfInfo" : { "groupId" : "ausfGroup001" } }`
	amfProfileStr := `{ "nfServices" : [ { "serviceInstanceId" : "0", "serviceName" : "namf-comm", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "amf", "transport" : "TCP", "port" : 29518 } ], "apiPrefix" : "http://amf:29518" }, { "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "amf", "transport" : "TCP", "port" : 29518 } ], "apiPrefix" : "http://amf:29518", "serviceInstanceId" : "1", "serviceName" : "namf-evts" }, { "apiPrefix" : "http://amf:29518", "serviceInstanceId" : "2", "serviceName" : "namf-mt", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "amf", "transport" : "TCP", "port" : 29518 } ] }, { "serviceInstanceId" : "3", "serviceName" : "namf-loc", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "amf", "transport" : "TCP", "port" : 29518 } ], "apiPrefix" : "http://amf:29518" }, { "serviceInstanceId" : "4", "serviceName" : "namf-oam", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "amf", "transport" : "TCP", "port" : 29518 } ], "apiPrefix" : "http://amf:29518" } ], "nfInstanceId" : "9f7d5a3f-88ab-4525-b31e-334da7faedab", "nfType" : "AMF", "nfStatus" : "REGISTERED", "plmnList" : [ { "mcc" : "208", "mnc" : "93" } ], "sNssais" : [ { "sst" : 1, "sd" : "010203" } ], "ipv4Addresses" : [ "amf" ], "amfInfo" : { "amfSetId" : "3f8", "amfRegionId" : "ca", "guamiList" : [ { "plmnId" : { "mcc" : "208", "mnc" : "93" }, "amfId" : "cafe00" } ], "taiList" : [ { "plmnId" : { "mcc" : "208", "mnc" : "93" }, "tac" : "1" } ] } }`

	var smfProfile models.NfProfile
	err := json.Unmarshal([]byte(smfProfileStr), &smfProfile)

	if err != nil {
		t.Log("Failed to convert smf profile")
		return
	}

	smfProfileItem := newNfProfileItem(&smfProfile, 180)

	pq.push(smfProfileItem)

	var ausfProfile models.NfProfile
	err = json.Unmarshal([]byte(ausfProfileStr), &ausfProfile)

	if err != nil {
		t.Log("Failed to convert ausf profile")
		return
	}

	ausfProfileItem := newNfProfileItem(&ausfProfile, 60)

	pq.push(ausfProfileItem)

	var amfProfile models.NfProfile
	err = json.Unmarshal([]byte(amfProfileStr), &amfProfile)

	if err != nil {
		t.Log("Failed to convert amf profile")
		return
	}

	amfProfileItem := newNfProfileItem(&amfProfile, 120)

	pq.push(amfProfileItem)

	item1 := pq.pop().(*NfProfileItem)
	assert.Equal(t, ausfProfileItem, item1)

	var item *NfProfileItem
	for i := 0; i < 2; i++ {
		item = pq.pop().(*NfProfileItem)
	}
	assert.Empty(t, item)
}

var count int32 = 1

func SEARCH(nrfUri string, targetNfType, requestNfType models.NfType,
	param *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (models.SearchResult, error) {
	fmt.Println("SEARCH")

	var searchResult models.SearchResult

	smfProfileStr := `{ "ipv4Addresses" : [ "smf" ], "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "smfInfo" : { "sNssaiSmfInfoList" : [ { "sNssai" : { "sst" : 1, "sd" : "010203" }, "dnnSmfInfoList" : [ { "dnn" : "internet" } ] } ] }, "nfServices" : [ { "apiPrefix" : "http://smf:29502", "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "serviceInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bfnsmf-pdusession", "serviceName" : "nsmf-pdusession", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "https://smf:29502/nsmf-pdusession/v1", "expiry" : "2022-08-17T05:31:40.997097141Z" } ], "scheme" : "https", "nfServiceStatus" : "REGISTERED" }, { "scheme" : "https", "nfServiceStatus" : "REGISTERED", "apiPrefix" : "http://smf:29502", "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "serviceInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bfnsmf-event-exposure", "serviceName" : "nsmf-event-exposure", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "https://smf:29502/nsmf-pdusession/v1", "expiry" : "2022-08-17T05:31:40.997097141Z" } ] } ], "nfInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bf", "plmnList" : [ { "mnc" : "93", "mcc" : "208" } ], "sNssais" : [ { "sd" : "010203", "sst" : 1 } ], "nfType" : "SMF", "nfStatus" : "REGISTERED" }`
	ausfProfileStr := `{ "nfServices" : [ { "serviceName" : "nausf-auth", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "ausf", "port" : 29509 } ], "serviceInstanceId" : "57d0a167-5283-4170-bdd8-881076049a81" } ], "nfInstanceId" : "57d0a167-5283-4170-bdd8-881076049a81", "nfType" : "AUSF", "nfStatus" : "REGISTERED", "plmnList" : [ { "mcc" : "208", "mnc" : "93" } ], "ipv4Addresses" : [ "ausf" ], "ausfInfo" : { "groupId" : "ausfGroup001" } }`
	amfProfileStr := `{ "nfServices" : [ { "serviceInstanceId" : "0", "serviceName" : "namf-comm", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "amf", "transport" : "TCP", "port" : 29518 } ], "apiPrefix" : "http://amf:29518" }, { "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "amf", "transport" : "TCP", "port" : 29518 } ], "apiPrefix" : "http://amf:29518", "serviceInstanceId" : "1", "serviceName" : "namf-evts" }, { "apiPrefix" : "http://amf:29518", "serviceInstanceId" : "2", "serviceName" : "namf-mt", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "amf", "transport" : "TCP", "port" : 29518 } ] }, { "serviceInstanceId" : "3", "serviceName" : "namf-loc", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "amf", "transport" : "TCP", "port" : 29518 } ], "apiPrefix" : "http://amf:29518" }, { "serviceInstanceId" : "4", "serviceName" : "namf-oam", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "amf", "transport" : "TCP", "port" : 29518 } ], "apiPrefix" : "http://amf:29518" } ], "nfInstanceId" : "9f7d5a3f-88ab-4525-b31e-334da7faedab", "nfType" : "AMF", "nfStatus" : "REGISTERED", "plmnList" : [ { "mcc" : "208", "mnc" : "93" } ], "sNssais" : [ { "sst" : 1, "sd" : "010203" } ], "ipv4Addresses" : [ "amf" ], "amfInfo" : { "amfSetId" : "3f8", "amfRegionId" : "ca", "guamiList" : [ { "plmnId" : { "mcc" : "208", "mnc" : "93" }, "amfId" : "cafe00" } ], "taiList" : [ { "plmnId" : { "mcc" : "208", "mnc" : "93" }, "tac" : "1" } ] } }`

	var smfProfile models.NfProfile
	err := json.Unmarshal([]byte(smfProfileStr), &smfProfile)

	if err != nil {
		return searchResult, nil
	}

	searchResult.NfInstances = append(searchResult.NfInstances, smfProfile)

	var ausfProfile models.NfProfile
	err = json.Unmarshal([]byte(ausfProfileStr), &ausfProfile)

	if err != nil {
		return searchResult, nil
	}
	searchResult.NfInstances = append(searchResult.NfInstances, ausfProfile)

	var amfProfile models.NfProfile
	err = json.Unmarshal([]byte(amfProfileStr), &amfProfile)

	if err != nil {
		return searchResult, nil
	}
	searchResult.NfInstances = append(searchResult.NfInstances, amfProfile)
	searchResult.ValidityPeriod = count * 120
	count++

	return searchResult, nil
}

func TestNfCaching(t *testing.T) {
	t.Log("Method Entry")
	amfSelf := amf_context.AMF_Self()

	smfProfileActual := `{ "ipv4Addresses" : [ "smf" ], "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "smfInfo" : { "sNssaiSmfInfoList" : [ { "sNssai" : { "sst" : 1, "sd" : "010203" }, "dnnSmfInfoList" : [ { "dnn" : "internet" } ] } ] }, "nfServices" : [ { "apiPrefix" : "http://smf:29502", "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "serviceInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bfnsmf-pdusession", "serviceName" : "nsmf-pdusession", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "https://smf:29502/nsmf-pdusession/v1", "expiry" : "2022-08-17T05:31:40.997097141Z" } ], "scheme" : "https", "nfServiceStatus" : "REGISTERED" }, { "scheme" : "https", "nfServiceStatus" : "REGISTERED", "apiPrefix" : "http://smf:29502", "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "serviceInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bfnsmf-event-exposure", "serviceName" : "nsmf-event-exposure", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "https://smf:29502/nsmf-pdusession/v1", "expiry" : "2022-08-17T05:31:40.997097141Z" } ] } ], "nfInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bf", "plmnList" : [ { "mnc" : "93", "mcc" : "208" } ], "sNssais" : [ { "sd" : "010203", "sst" : 1 } ], "nfType" : "SMF", "nfStatus" : "REGISTERED" }`
	ausfProfileActual := `{ "nfServices" : [ { "serviceName" : "nausf-auth", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "1.0.0" } ], "scheme" : "http", "nfServiceStatus" : "REGISTERED", "ipEndPoints" : [ { "ipv4Address" : "ausf", "port" : 29509 } ], "serviceInstanceId" : "57d0a167-5283-4170-bdd8-881076049a81" } ], "nfInstanceId" : "57d0a167-5283-4170-bdd8-881076049a81", "nfType" : "AUSF", "nfStatus" : "REGISTERED", "plmnList" : [ { "mcc" : "208", "mnc" : "93" } ], "ipv4Addresses" : [ "ausf" ], "ausfInfo" : { "groupId" : "ausfGroup001" } }`

	InitNrfCaching(time.Duration(30), SEARCH)
	searchResult, err := SearchNFInstances(amfSelf.NrfUri, models.NfType_SMF, models.NfType_AMF, nil)
	if err == nil {
		t.Log(len(searchResult.NfInstances))
		assert.Equal(t, searchResult.NfInstances[0], smfProfileActual)
	}

	sr, e := SearchNFInstances(amfSelf.NrfUri, models.NfType_AUSF, models.NfType_AMF, nil)
	if e == nil {
		t.Log(len(sr.NfInstances))
		assert.Equal(t, sr.NfInstances[0], ausfProfileActual)
	}

	time.Sleep(180000 * time.Millisecond)
}

func TestCleanupExpiredItems(t *testing.T) {
	t.Log("Method Entry")

	var c *NrfCache
	pq := newNfProfilePriorityQ()

	smfProfileStr := `{ "ipv4Addresses" : [ "smf" ], "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "smfInfo" : { "sNssaiSmfInfoList" : [ { "sNssai" : { "sst" : 1, "sd" : "010203" }, "dnnSmfInfoList" : [ { "dnn" : "internet" } ] } ] }, "nfServices" : [ { "apiPrefix" : "http://smf:29502", "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "serviceInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bfnsmf-pdusession", "serviceName" : "nsmf-pdusession", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "https://smf:29502/nsmf-pdusession/v1", "expiry" : "2022-08-17T05:31:40.997097141Z" } ], "scheme" : "https", "nfServiceStatus" : "REGISTERED" }, { "scheme" : "https", "nfServiceStatus" : "REGISTERED", "apiPrefix" : "http://smf:29502", "allowedPlmns" : [ { "mcc" : "208", "mnc" : "93" } ], "serviceInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bfnsmf-event-exposure", "serviceName" : "nsmf-event-exposure", "versions" : [ { "apiVersionInUri" : "v1", "apiFullVersion" : "https://smf:29502/nsmf-pdusession/v1", "expiry" : "2022-08-17T05:31:40.997097141Z" } ] } ], "nfInstanceId" : "b926f193-1083-49a8-adb3-5fcf57a1f0bf", "plmnList" : [ { "mnc" : "93", "mcc" : "208" } ], "sNssais" : [ { "sd" : "010203", "sst" : 1 } ], "nfType" : "SMF", "nfStatus" : "REGISTERED" }`

	var smfProfile models.NfProfile
	err := json.Unmarshal([]byte(smfProfileStr), &smfProfile)

	if err != nil {
		t.Log("Failed to convert smf profile")
		return
	}

	smfProfileItem := newNfProfileItem(&smfProfile, 60)
	cache := make(map[string]*NfProfileItem)
	cache[smfProfile.NfInstanceId] = smfProfileItem
	pq.push(smfProfileItem)

	c = &NrfCache{
		cache:            cache,
		priorityQ:        pq,
		evictionInterval: defaultCacheTTl,
		done:             make(chan struct{}),
	}

	time.Sleep(70 * time.Second)
	c.cleanupExpiredItems()
	assert.Empty(t, pq)
}
