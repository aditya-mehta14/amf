// SPDX-FileCopyrightText: 2022 Infosys Limited
// Copyright 2019 free5GC.org
//
// SPDX-License-Identifier: Apache-2.0
//

package nrf_cache

import (
	"encoding/json"
	"github.com/omec-project/openapi/Nnrf_NFDiscovery"
	"github.com/omec-project/openapi/models"
	"regexp"
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

	matchFound := true

	if opts.ServiceNames.IsSet() {
		reqServiceNames := opts.ServiceNames.Value().([]models.ServiceName)
		matchCount := 0
		for _, sn := range reqServiceNames {
			for i := 0; i < len(*profile.NfServices); i++ {
				if (*profile.NfServices)[i].ServiceName == sn {
					matchCount++
					break
				}
			}
		}

		if matchCount == 0 {
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

			// Snssai in the smfInfo has priority
			if profile.SmfInfo != nil && profile.SmfInfo.SNssaiSmfInfoList != nil {
				for _, s := range *profile.SmfInfo.SNssaiSmfInfoList {
					if s.SNssai != nil && (*s.SNssai) == snssai {
						matchCount++
					}
				}
			} else if profile.AllowedNssais != nil {
				for _, s := range *profile.AllowedNssais {
					if s == snssai {
						matchCount++
					}
				}
			}

		}

		// if at least one matching snssai has been found
		if matchCount == 0 {
			matchFound = false
		}

	}

	// validate dnn
	if matchFound && opts.Dnn.IsSet() {
		// if a dnn is provided by the upper layer, check for the exact match
		// or wild card match
		dnnMatched := false

		if profile.SmfInfo != nil && profile.SmfInfo.SNssaiSmfInfoList != nil {
		matchDnnLoop:
			for _, s := range *profile.SmfInfo.SNssaiSmfInfoList {
				if s.DnnSmfInfoList != nil {
					for _, d := range *s.DnnSmfInfoList {
						if d.Dnn == opts.Dnn.Value() || d.Dnn == "*" {
							dnnMatched = true
							break matchDnnLoop
						}
					}
				}
			}
		}
		matchFound = dnnMatched
	}

	return matchFound
}

func MatchAusfProfile(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool {
	matchFound := true
	if opts.Supi.IsSet() {
		if profile.AusfInfo != nil && len(profile.AusfInfo.SupiRanges) > 0 {
			matchCount := 0
			for _, s := range profile.AusfInfo.SupiRanges {
				if len(s.Pattern) > 0 {
					r, _ := regexp.Compile(s.Pattern)
					if r.MatchString(opts.Supi.Value()) {
						matchCount++
					}

				} else if s.Start <= opts.Supi.Value() && opts.Supi.Value() < s.End {
					matchCount++
				}
			}

			if matchCount == 0 {
				matchFound = false
			}
		}
	}
	return matchFound
}

func MatchNssfProfile(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool {
	return true
}

func MatchAmfProfile(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool {
	return true
}

func MatchPcfProfile(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool {
	matchFound := true
	if opts.Supi.IsSet() {
		if profile.AusfInfo != nil && len(profile.AusfInfo.SupiRanges) > 0 {
			matchCount := 0
			for _, s := range profile.AusfInfo.SupiRanges {
				if len(s.Pattern) > 0 {
					r, _ := regexp.Compile(s.Pattern)
					if r.MatchString(opts.Supi.Value()) {
						matchCount++
					}

				} else if s.Start <= opts.Supi.Value() && opts.Supi.Value() < s.End {
					matchCount++
				}
			}

			if matchCount == 0 {
				matchFound = false
			}
		}
	}
	return matchFound
}

func MatchUdmProfile(profile *models.NfProfile, opts *Nnrf_NFDiscovery.SearchNFInstancesParamOpts) bool {
	matchFound := true
	if opts.Supi.IsSet() {
		if profile.AusfInfo != nil && len(profile.AusfInfo.SupiRanges) > 0 {
			matchCount := 0
			for _, s := range profile.AusfInfo.SupiRanges {
				if len(s.Pattern) > 0 {
					r, _ := regexp.Compile(s.Pattern)
					if r.MatchString(opts.Supi.Value()) {
						matchCount++
					}

				} else if s.Start <= opts.Supi.Value() && opts.Supi.Value() < s.End {
					matchCount++
				}
			}

			if matchCount == 0 {
				matchFound = false
			}
		}
	}
	return matchFound
}
