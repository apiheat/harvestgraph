package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/urfave/cli"
)

func cmdSearch(c *cli.Context) error {
	var configs []ConfigurationMap

	nlID := c.String("id")

	metadata := Metadata{
		ID:   nlID,
		Name: c.String("name"),
	}

	// Read config file
	byteValue := readFileBytes(c.String("map-file"))
	err := json.Unmarshal(byteValue, &configs)

	if err != nil {
		log.Fatal(err)
	}

	for _, config := range configs {
		log.Printf("Executing search for network list ID %s, Name %s in %s configuration", nlID, c.String("name"), config.ConfigName)

		var rpFound, mtFound, spFound int

		//cMap used to append all information for Network list usage in one config
		// cMap := ConfigurationMap{}

		// Read Security Configuration file
		configPath := fmt.Sprintf("%s/%s.json", c.String("source"), strconv.Itoa(config.ConfigID))
		cfile := readFileBytes(configPath)

		// Rate Policies search
		rpFound = ratePolicySearch(nlID, cfile)

		// Match Targets Search
		mtFound = matchTargetSearch(nlID, cfile)

		// Security Policies Search
		spFound = securityPolicySearch(nlID, cfile)

		if rpFound > 0 || mtFound > 0 || spFound > 0 {
			cMap.ConfigID = config.ConfigID
			cMap.ConfigName = config.ConfigName
			metadata.Usage = append(metadata.Usage, cMap)
		}
	}

	// Output metadata json
	if c.String("output") == "json" {
		if c.String("destination") == "" {
			outputJSON(metadata)
			return nil
		}

		jsonF, err := json.MarshalIndent(metadata, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		resultFilePath := fmt.Sprintf("%s/%s.metadata", c.String("destination"), nlID)
		err = ioutil.WriteFile(resultFilePath, jsonF, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Output dot graphs
	if c.String("output") == "dot" {
		s := buildGraph(metadata)

		if c.String("destination") == "" {
			fmt.Println(s)
			return nil
		}

		graphFilePath := fmt.Sprintf("%s/%s.dot", c.String("destination"), nlID)
		err = ioutil.WriteFile(graphFilePath, []byte(s), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
