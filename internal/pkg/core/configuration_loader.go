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
	"fmt"
	"github.com/fbiville/headache/internal/pkg/fs"
	json_schema "github.com/xeipuuv/gojsonschema"
	"strings"
)

type ConfigurationLoader interface {
	ValidateAndLoad(path string) (*Configuration, error)
	LoadFile(path string) (*Configuration, error)
	LoadBytes(bytes []byte) (*Configuration, error)
}

type ConfigurationFileLoader struct {
	Reader         fs.FileReader
	SchemaLocation string
	SchemaLoader   JsonSchemaLoader
}

func (loader *ConfigurationFileLoader) ValidateAndLoad(path string) (*Configuration, error) {
	schema := loader.SchemaLoader.Load(loader.SchemaLocation)
	configuration, err := loader.LoadFile(path)
	if err != nil {
		return nil, err
	}
	if schema == nil {
		return configuration, nil
	}
	// normalize comment style names
	configuration.CommentStyle = strings.ToLower(configuration.CommentStyle)
	configurationPayload, err := json.Marshal(configuration)
	if err != nil {
		return nil, err
	}
	result, err := schema.Validate(json_schema.NewBytesLoader(configurationPayload))
	if err != nil {
		return nil, err
	}
	errors := result.Errors()
	if len(errors) > 0 {
		return nil, fmt.Errorf(report(errors))
	}
	configuration.Path = &path
	return configuration, nil
}

func (loader *ConfigurationFileLoader) LoadFile(path string) (*Configuration, error) {
	configurationPayload, err := loader.Reader.Read(path)
	if err != nil {
		return nil, err
	}
	return loader.LoadBytes(configurationPayload)
}

func (loader *ConfigurationFileLoader) LoadBytes(configurationPayload []byte) (*Configuration, error) {
	result := Configuration{}
	err := json.Unmarshal(configurationPayload, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func report(errors []json_schema.ResultError) string {
	builder := strings.Builder{}
	for _, validationError := range errors {
		details := validationError.Details()
		field := details["field"]
		builder.WriteString(fmt.Sprintf("Error with field '%s': %s", field, description(field, validationError)))
		builder.WriteString("\n")
	}
	result := builder.String()
	return result[:len(result)-1]
}

func description(field interface{}, validationError json_schema.ResultError) string {
	for _, name := range []string{"Year", "YearRange", "StartYear", "EndYear"} {
		if field == fmt.Sprintf("data.%s", name) {
			return fmt.Sprintf("%s is a reserved data parameter and cannot be used", name)
		}
	}
	return validationError.Description()
}
