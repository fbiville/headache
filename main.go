/*
 * Copyright 2018 Florent Biville (@fbiville)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"flag"
	. "github.com/fbiville/headache/core"
	"io/ioutil"
)

func main() {
	configFile := flag.String("configuration", "headache.json", "Path to configuration file")

	flag.Parse()

	rawConfiguration := readConfiguration(configFile)
	configuration, err := ParseConfiguration(rawConfiguration)
	if err != nil {
		panic(err)
	}
	Run(configuration)
}

func readConfiguration(configFile *string) Configuration {
	flag.Parse()
	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}
	result := Configuration{}
	err = json.Unmarshal(file, &result)
	if err != nil {
		panic(err)
	}
	return result
}

