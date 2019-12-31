/*
 * Copyright 2019 Florent Biville (@fbiville)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"encoding/json"
	"github.com/fbiville/headache/fs"
	jsonsch "github.com/xeipuuv/gojsonschema"
	"log"
)

type ConfigurationLoader struct {
	Reader fs.FileReader
}

func (cl *ConfigurationLoader) ReadConfiguration(configFile *string) (*Configuration, error) {
	err := cl.validateConfiguration(configFile)
	if err != nil {
		return nil, err
	}

	payload, err := cl.Reader.Read(*configFile)
	if err != nil {
		return nil, err
	}
	configuration, err := cl.UnmarshalConfiguration(payload)
	if err != nil {
		return nil, err
	}
	configuration.Path = configFile
	return configuration, err
}

func (cl *ConfigurationLoader) UnmarshalConfiguration(configurationPayload []byte) (*Configuration, error) {
	result := Configuration{}
	err := json.Unmarshal(configurationPayload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (cl *ConfigurationLoader) validateConfiguration(configFile *string) error {
	schema := loadSchema()
	if schema == nil {
		return nil
	}
	jsonSchemaValidator := JsonSchemaValidator{
		Schema:     schema,
		FileReader: cl.Reader,
	}
	return jsonSchemaValidator.Validate("file://" + *configFile)
}

func loadSchema() *jsonsch.Schema {
	schema, err := jsonsch.NewSchema(jsonsch.NewReferenceLoader("https://fbiville.github.io/headache/schema.json"))
	if err != nil {
		log.Printf("headache configuration warning: cannot load schema, skipping configuration validation. See reason below:\n\t%v\n", err)
		return nil
	}
	return schema
}
