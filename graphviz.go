package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/awalterschulze/gographviz"
)

func buildGraph(m Metadata) string {
	// Graph
	g := gographviz.NewGraph()

	nodeAttrs := map[string]string{
		"shape": "Mrecord",
	}

	gName := fmt.Sprintf("%q", m.Name)
	if err := g.SetName(gName); err != nil {
		log.Fatal(err)
	}

	labelData := fmt.Sprintf("{%s | id=%s}", gName, m.ID)
	nodeAttrs["label"] = fmt.Sprintf("%q", labelData)

	if err := g.AddNode(gName, gName, nodeAttrs); err != nil {
		log.Fatal(err)
	}
	if err := g.SetDir(true); err != nil {
		log.Fatal(err)
	}

	for _, config := range m.Usage {
		name := fmt.Sprintf("%q", config.ConfigName)
		nameString := config.ConfigName
		if config.ConfigName == "" {
			name = fmt.Sprintf("%q", "WAF Security File")
			nameString = "WAF Security File"
		}

		labelData := fmt.Sprintf("{%s | id=%d}", name, config.ConfigID)
		nodeAttrs["label"] = fmt.Sprintf("%q", labelData)

		if err := g.AddNode(gName, name, nodeAttrs); err != nil {
			log.Fatal(err)
		}
		if err := g.AddEdge(gName, name, true, nil); err != nil {
			log.Fatal(err)
		}

		if len(config.Policies) > 0 {

			if err := g.AddNode(name, fmt.Sprintf("%q", fmt.Sprintf("Security Policies %s", nameString)), nil); err != nil {
				log.Fatal(err)
			}
			if err := g.AddEdge(name, fmt.Sprintf("%q", fmt.Sprintf("Security Policies %s", nameString)), true, nil); err != nil {
				log.Fatal(err)
			}
		}

		if len(config.RatePolicies) > 0 {
			if err := g.AddNode(name, fmt.Sprintf("%q", fmt.Sprintf("Rate Limit Policies %s", nameString)), nil); err != nil {
				log.Fatal(err)
			}
			if err := g.AddEdge(name, fmt.Sprintf("%q", fmt.Sprintf("Rate Limit Policies %s", nameString)), true, nil); err != nil {
				log.Fatal(err)
			}
		}

		if len(config.MatchTargets) > 0 {
			if err := g.AddNode(name, fmt.Sprintf("%q", fmt.Sprintf("Match Targets %s", nameString)), nil); err != nil {
				log.Fatal(err)
			}
			if err := g.AddEdge(name, fmt.Sprintf("%q", fmt.Sprintf("Match Targets %s", nameString)), true, nil); err != nil {
				log.Fatal(err)
			}
		}

		// Policies graph
		for _, policy := range config.Policies {
			pName := fmt.Sprintf("%q", policy.Name)
			labelData := fmt.Sprintf("{%s | id=%s | type=%s | action=%s}", policy.Name, policy.ID, policy.Type, policy.Action)
			nodeAttrs["label"] = fmt.Sprintf("%q", labelData)
			if err := g.AddNode(gName, pName, nodeAttrs); err != nil {
				log.Fatal(err)
			}

			if err := g.AddEdge(fmt.Sprintf("%q", fmt.Sprintf("Security Policies %s", nameString)), pName, true, nil); err != nil {
				log.Fatal(err)
			}
		}

		// Match Targets graph
		for _, mTarget := range config.MatchTargets {
			tName := fmt.Sprintf("%d", mTarget.ID)

			labelData := fmt.Sprintf("{id=%d | securityPolicyID=%s | type=%s}", mTarget.ID, mTarget.SecurityPolicyID, mTarget.Type)
			nodeAttrs["label"] = fmt.Sprintf("%q", labelData)

			if err := g.AddNode(gName, tName, nodeAttrs); err != nil {
				log.Fatal(err)
			}

			if err := g.AddEdge(fmt.Sprintf("%q", fmt.Sprintf("Match Targets %s", nameString)), tName, true, nil); err != nil {
				log.Fatal(err)
			}

			hosts := []string{}
			for _, hostname := range mTarget.Hostnames.([]interface{}) {
				hName := fmt.Sprintf("%s", hostname)

				hosts = append(hosts, hName)
			}

			labelDataHosts := fmt.Sprintf("{Match Target Hostnames | %s}", strings.Join(hosts, " | "))
			nodeAttrs["label"] = fmt.Sprintf("%q", labelDataHosts)

			hostsNodeName := fmt.Sprintf("%s hostnames", tName)
			if err := g.AddNode(gName, fmt.Sprintf("%q", hostsNodeName), nodeAttrs); err != nil {
				log.Fatal(err)
			}

			if err := g.AddEdge(tName, fmt.Sprintf("%q", hostsNodeName), true, nil); err != nil {
				log.Fatal(err)
			}

			paths := []string{}
			for _, path := range mTarget.Paths.([]interface{}) {
				pName := fmt.Sprintf("%s", path)

				paths = append(paths, pName)
			}

			labelDataPaths := fmt.Sprintf("{Match Target Paths | %s}", strings.Join(paths, " | "))
			nodeAttrs["label"] = fmt.Sprintf("%q", labelDataPaths)

			pathsNodeName := fmt.Sprintf("%s paths", tName)
			if err := g.AddNode(gName, fmt.Sprintf("%q", pathsNodeName), nodeAttrs); err != nil {
				log.Fatal(err)
			}

			if err := g.AddEdge(fmt.Sprintf("%q", hostsNodeName), fmt.Sprintf("%q", pathsNodeName), true, nil); err != nil {
				log.Fatal(err)
			}
		}

		// Rate Policies Graph
		for _, rateP := range config.RatePolicies {
			rpName := fmt.Sprintf("%q", rateP.Name)

			labelData := fmt.Sprintf("{%s | id=%d}", rateP.Name, rateP.ID)
			nodeAttrs["label"] = fmt.Sprintf("%q", labelData)

			if err := g.AddNode(gName, rpName, nodeAttrs); err != nil {
				log.Fatal(err)
			}

			if err := g.AddEdge(fmt.Sprintf("%q", fmt.Sprintf("Rate Limit Policies %s", nameString)), rpName, true, nil); err != nil {
				log.Fatal(err)
			}
		}
	}

	return g.String()
}
