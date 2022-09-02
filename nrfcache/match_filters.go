// SPDX-FileCopyrightText: 2022 Open Networking Foundation <info@opennetworking.org>
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//

package nrf_cache

import (
	"encoding/json"
	"github.com/omec-project/openapi/Nnrf_NFDiscovery"
	"github.com/omec-project/openapi/models"
)

type MatchFilter func(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool

type MatchFilters map[models.NfType]MatchFilter

var matchFilters = MatchFilters{
	models.NfType_SMF:  MatchSmfProfile,
	models.NfType_AUSF: MatchAusfProfile,
	models.NfType_PCF:  MatchPcfProfile,
	models.NfType_NSSF: MatchNssfProfile,
	models.NfType_UDM:  MatchUdmProfile,
	models.NfType_AMF:  MatchAmfProfile,
}

func MatchSmfProfile(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool {

	matchFound := false

	// validate dnn
	if opts.Dnn.IsSet() {
		if profile.SmfInfo != nil && profile.SmfInfo.SNssaiSmfInfoList != nil {
			for _, s := range *profile.SmfInfo.SNssaiSmfInfoList {
				for _, d := range *s.DnnSmfInfoList {
					if d.Dnn == opts.Dnn.Value() {
						matchFound = true
						break
					}
				}
			}
		}
	}

	if matchFound && opts.ServiceNames.IsSet() {
		reqServiceNames := opts.ServiceNames.Value().([]models.ServiceName)
		matchCount := 0
		for _, sn := range reqServiceNames {
			for _, psn := range *profile.NfServices {
				if psn.ServiceName == sn {
					matchCount++
				}
			}
		}

		if matchCount != len(reqServiceNames) {
			matchFound = false
		}
	}

	if matchFound && opts.Snssais.IsSet() {
		reqSnssais := opts.Snssais.Value().([]string)
		matchCount := 0

		for _, reqSnssai := range reqSnssais {
			var snssai models.Snssai
			err := json.Unmarshal([]byte(reqSnssai), &snssai)
			if err != nil {
				return false
			}

			for _, snssaiInfoItem := range *profile.SmfInfo.SNssaiSmfInfoList {
				if *snssaiInfoItem.SNssai == snssai {
					matchCount++
				}
			}

		}

		if matchCount != len(reqSnssais) {
			matchFound = false
		}

	}

	return matchFound
}

func MatchAusfProfile(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool {
	return true
}

func MatchNssfProfile(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool {
	return true
}

func MatchAmfProfile(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool {
	return true
}

func MatchPcfProfile(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool {
	return true
}

func MatchUdmProfile(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool {
	return true
}
