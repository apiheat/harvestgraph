package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// PrintJSON pretty print JSON string
func printJSON(str string) {
	var prettyJSON bytes.Buffer

	error := json.Indent(&prettyJSON, []byte(str), "", "    ")
	if error != nil {
		log.Fatal("JSON parse error: ", error)
	}

	fmt.Println(string(prettyJSON.Bytes()))
	return
}

// OutputJSON displays output of query for alerts in JSON format
func outputJSON(input interface{}) {
	b, err := json.Marshal(input)
	if err != nil {
		log.Fatal(err)
	}

	printJSON(string(b))
}

// readFileBytes
func readFileBytes(fpath string) []byte {
	log.Printf("Openning file %s", fpath)

	jsonFile, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	return byteValue
}
