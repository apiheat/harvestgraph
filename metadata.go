package main

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func ratePolicySearch(listID string, cfile []byte, cMap *ConfigurationMap) int {
	var i int

	// Rate Policies search
	for _, name := range gjson.GetBytes(cfile, "ratePolicies").Array() {
		nlc := name.Get("additionalMatchOptions.#(type==NetworkListCondition)")

		cond := fmt.Sprintf("values.#(==%s)", listID)
		present := nlc.Get(cond).Raw

		if present != "" {
			rpObj := RatePolicy{
				ID:   name.Get("id").Int(),
				Name: name.Get("name").String(),
			}

			cMap.RatePolicies = append(cMap.RatePolicies, rpObj)
			i++
		}
	}

	return i
}

func matchTargetSearch(listID string, cfile []byte, cMap *ConfigurationMap) int {
	var i int

	// Match Targets Search
	for _, name := range gjson.GetBytes(cfile, "matchTargets.websiteTargets").Array() {
		netListsPresent := name.Get("bypassNetworkLists")
		if netListsPresent.Exists() {
			cond := fmt.Sprintf("bypassNetworkLists.#(id==%s)", listID)
			present := name.Get(cond).Raw

			if present != "" {
				mtObj := MatchTarget{
					ID:                name.Get("id").Int(),
					Hostnames:         name.Get("hostnames").Value(),
					Paths:             name.Get("filePaths").Value(),
					NegativePathMatch: name.Get("isNegativePathMatch").Bool(),
					SecurityPolicyID:  name.Get("securityPolicy.policyId").String(),
					Type:              "bypass",
				}

				cMap.MatchTargets = append(cMap.MatchTargets, mtObj)
				i++
			}
		}
	}
	return i
}

func networkListSearch(listID, listType, listAction string, ipGeoFirewallNode, spNode gjson.Result, cMap *ConfigurationMap) bool {
	var found bool

	searchString := fmt.Sprintf("%sControls.%sIPNetworkLists", listType, listAction)

	networkList := ipGeoFirewallNode.Get(searchString)
	if networkList.Exists() {
		cond := fmt.Sprintf("networkList.#(==%s)", listID)
		present := networkList.Get(cond).Raw

		if present != "" {
			spObj := SecurityPolicy{
				ID:     spNode.Get("id").String(),
				Name:   spNode.Get("name").String(),
				Type:   fmt.Sprintf("%sControls", listType),
				Action: listAction,
			}

			cMap.Policies = append(cMap.Policies, spObj)
			found = true
		}
	}

	return found
}

func securityPolicySearch(listID string, cfile []byte, cMap *ConfigurationMap) int {
	var i int

	for _, name := range gjson.GetBytes(cfile, "securityPolicies").Array() {
		var ipAllowedFound, ipBlockedFound, geoAllowedFound, geoBlockedFound bool
		ipGeoFirewall := name.Get("ipGeoFirewall")

		// IP
		ipAllowedFound = networkListSearch(listID, "ip", "allowed", ipGeoFirewall, name, cMap)
		ipBlockedFound = networkListSearch(listID, "ip", "blocked", ipGeoFirewall, name, cMap)

		// Geo
		geoAllowedFound = networkListSearch(listID, "geo", "allowed", ipGeoFirewall, name, cMap)
		geoBlockedFound = networkListSearch(listID, "geo", "blocked", ipGeoFirewall, name, cMap)

		if ipAllowedFound || ipBlockedFound || geoAllowedFound || geoBlockedFound {
			i++
		}
	}

	return i
}
