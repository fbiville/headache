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
	"github.com/fbiville/headache/fs"
	jsonsch "github.com/xeipuuv/gojsonschema"
	"log"
)

func main() {
	configFile := flag.String("configuration", "headache.json", "Path to configuration file")

	flag.Parse()

	systemConfig := DefaultSystemConfiguration()
	fileSystem := systemConfig.FileSystem
	rawConfiguration := readConfiguration(configFile, fileSystem.FileReader)
	executionTracker := &fs.ExecutionVcsTracker{
		Versioning: systemConfig.VersioningClient.GetClient(),
		FileSystem: fileSystem,
		Clock:      systemConfig.Clock,
	}
	matcher := &fs.ZglobPathMatcher{}
	configuration, err := ParseConfiguration(rawConfiguration, systemConfig, executionTracker, matcher)
	if err != nil {
		log.Fatalf("headache configuration error, cannot parse\n\t%v\n", err)
	}
	Run(configuration, fileSystem)
}

func readConfiguration(configFile *string, reader fs.FileReader) Configuration {
	validateConfiguration(reader, configFile)

	file, err := reader.Read(*configFile)
	if err != nil {
		log.Fatalf("headache configuration error, cannot read file:\n\t%v", err)
	}
	result := Configuration{}
	err = json.Unmarshal(file, &result)
	if err != nil {
		log.Fatalf("headache configuration error, cannot unmarshall JSON:\n\t%v", err)
	}
	return result
}

func validateConfiguration(reader fs.FileReader, configFile *string) {
	jsonSchemaValidator := JsonSchemaValidator{
		Schema:     loadSchema(),
		FileReader: reader,
	}
	validationError := jsonSchemaValidator.Validate("file://" + *configFile)
	if validationError != nil {
		log.Fatalf("headache configuration error, validation failed\n\t%s\n", validationError)
	}
}

func loadSchema() *jsonsch.Schema {
	schema, err := jsonsch.NewSchema(jsonsch.NewReferenceLoader("file://docs/schema.json"))
	if err != nil {
		log.Fatalf("headache configuration error, cannot load JSON schema\n\t%v\n", err)
	}
	return schema
}
