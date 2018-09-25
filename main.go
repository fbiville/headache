package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)

func main() {
	rawConfiguration := readConfiguration()
	configuration, err := ParseConfiguration(rawConfiguration)
	if err != nil {
		panic(err)
	}
	InsertHeader(configuration)
}

func readConfiguration() Configuration {
	configFile := flag.String("configuration", "license.json", "Path to configuration file")
	flag.Parse()
	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}
	rawConfiguration := Configuration{}
	json.Unmarshal(file, &rawConfiguration)
	return rawConfiguration
}
