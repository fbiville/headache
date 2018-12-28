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

package core

import (
	"fmt"
	"github.com/fbiville/headache/fs"
	json "github.com/xeipuuv/gojsonschema"
	"strings"
)

type ConfigurationValidator interface {
	Validate(document string) error
}

type JsonSchemaValidator struct {
	FileReader fs.FileReader
	Schema     *json.Schema
}

func (validator *JsonSchemaValidator) Validate(path string) error {
	documentLoader := json.NewReferenceLoaderFileSystem(path, validator.FileReader)
	result, err := validator.Schema.Validate(documentLoader)
	if err != nil {
		return err
	}
	errors := result.Errors()
	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf(report(errors))
}

func report(errors []json.ResultError) string {
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

func description(field interface{}, validationError json.ResultError) string {
	if field == "data.Year" {
		return "Year is a reserved data parameter and cannot be used"
	}
	return validationError.Description()
}
