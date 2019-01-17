/*
 * Copyright 2018-2019 Florent Biville (@fbiville)
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

package core_test

import (
	"github.com/fbiville/headache/core"
	"github.com/fbiville/headache/fs_mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	json "github.com/xeipuuv/gojsonschema"
	"io"
	"net/http"
	"os"
)

var _ = Describe("Configuration validator", func() {
	var (
		t          GinkgoTInterface
		fileReader *fs_mocks.FileReader
		validator  core.JsonSchemaValidator
	)

	BeforeEach(func() {
		t = GinkgoT()
		fileReader = new(fs_mocks.FileReader)
		validator = core.JsonSchemaValidator{
			FileReader: fileReader,
			Schema:     schemaFrom(json.NewReferenceLoader("file://../docs/schema.json")),
		}
	})

	AfterEach(func() {
		fileReader.AssertExpectations(t)
	})

	It("accepts minimal valid configuration", func() {
		fileReader.On("Open", "docs.json").
			Return(inMemoryFile(`{"headerFile": "some-file.txt", "style": "SlashStar", "includes": ["**/*.go"]}`), nil)

		validationError := validator.Validate("file://docs.json")

		Expect(validationError).To(BeNil())
	})

	It("accepts valid configuration", func() {
		fileReader.On("Open", "docs.json").
			Return(inMemoryFile(`{"headerFile": "some-file.txt", "style": "SlashStar", "includes": ["**/*.go"], "data": {"FooBar": true}}`), nil)

		validationError := validator.Validate("file://docs.json")

		Expect(validationError).To(BeNil())
	})

	It("rejects configuration with missing header file", func() {
		fileReader.On("Open", "docs.json").
			Return(inMemoryFile(`{"style": "SlashStar", "includes": ["**/*.go"]}`), nil)

		validationError := validator.Validate("file://docs.json")

		Expect(validationError.Error()).To(HaveSuffix("headerFile is required"))
	})

	It("rejects configuration with missing comment style", func() {
		fileReader.On("Open", "docs.json").
			Return(inMemoryFile(`{"headerFile": "some-header.txt", "includes": ["**/*.go"]}`), nil)

		validationError := validator.Validate("file://docs.json")

		Expect(validationError.Error()).To(HaveSuffix("style is required"))
	})

	It("rejects configuration with invalid comment style", func() {
		fileReader.On("Open", "docs.json").
			Return(inMemoryFile(`{"headerFile": "some-header.txt", "style": "invalid", includes": ["**/*.go"]}`), nil)

		validationError := validator.Validate("file://docs.json")

		Expect(validationError).To(MatchError("invalid character 'i' looking for beginning of object key string"))
	})

	It("rejects configuration with missing includes", func() {
		fileReader.On("Open", "docs.json").
			Return(inMemoryFile(`{"headerFile": "some-header.txt", "style": "SlashStar"}`), nil)

		validationError := validator.Validate("file://docs.json")

		Expect(validationError.Error()).To(HaveSuffix("includes is required"))
	})

	It("rejects configuration with empty includes", func() {
		fileReader.On("Open", "docs.json").
			Return(inMemoryFile(`{"headerFile": "some-header.txt", "style": "SlashSlash", "includes": []}`), nil)

		validationError := validator.Validate("file://docs.json")

		Expect(validationError.Error()).To(HaveSuffix("Array must have at least 1 items"))
	})

	It("rejects configuration with reserved year parameter", func() {
		fileReader.On("Open", "docs.json").
			Return(inMemoryFile(`{"headerFile": "some-header.txt", "style": "SlashSlash", "includes": ["**/*.*"], "data": {"Year": 2019}}`), nil)

		validationError := validator.Validate("file://docs.json")

		Expect(validationError.Error()).To(HaveSuffix("Year is a reserved data parameter and cannot be used"))
	})

})

func schemaFrom(loader json.JSONLoader) *json.Schema {
	schema, err := json.NewSchema(loader)
	if err != nil {
		panic(err)
	}
	return schema
}

func inMemoryFile(contents string) http.File {
	return InMemoryFile{Contents: []byte(contents)}
}

type InMemoryFile struct {
	Contents []byte
}

func (InMemoryFile) Close() error { return nil }
func (file InMemoryFile) Read(p []byte) (int, error) {
	contents := file.Contents
	end := min(len(p), len(contents))
	copiedBytes := copy(p, contents[0:end])
	if len(p) > len(contents) {
		return copiedBytes, io.EOF
	}
	return copiedBytes, nil
}
func (InMemoryFile) Seek(offset int64, whence int) (int64, error) { panic("not implemented") }
func (InMemoryFile) Readdir(count int) ([]os.FileInfo, error)     { panic("not implemented") }
func (InMemoryFile) Stat() (os.FileInfo, error)                   { panic("not implemented") }

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
