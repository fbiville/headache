/*
 * Copyright 2018-2019 Florent Biville (@fbiville)
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

package core_test

import (
	"fmt"
	"github.com/fbiville/headache/internal/pkg/core"
	"github.com/fbiville/headache/internal/pkg/fs_mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	json "github.com/xeipuuv/gojsonschema"
)

var _ = Describe("Configuration loader", func() {
	var (
		t          GinkgoTInterface
		fileReader *fs_mocks.FileReader
		loader     core.ConfigurationFileLoader

		configurationUri string
	)

	BeforeEach(func() {
		t = GinkgoT()
		fileReader = new(fs_mocks.FileReader)
		loader = core.ConfigurationFileLoader{
			Reader:         fileReader,
			SchemaLocation: "file://../../../docs/schema.json",
			SchemaLoader:   &LocalSchemaLoader{},
		}
		configurationUri = fmt.Sprintf("file://%s", "headache.json")
	})

	AfterEach(func() {
		fileReader.AssertExpectations(t)
	})

	It("accepts and loads minimal valid configuration", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-file.txt", "style": "slashstar", "includes": ["**/*.go"]}`), nil)

		configuration, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError).To(BeNil())
		Expect(configuration).To(Equal(&core.Configuration{
			HeaderFile:   "some-file.txt",
			CommentStyle: "slashstar",
			Includes:     []string{"**/*.go"},
			Path:         &configurationUri,
		}))
	})

	It("accepts and loads valid configuration with // comment style", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-file.txt", "style": "slashslash", "includes": ["**/*.go"], "data": {"FooBar": "true"}}`), nil)

		configuration, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError).To(BeNil())
		Expect(configuration).To(Equal(&core.Configuration{
			HeaderFile:   "some-file.txt",
			CommentStyle: "slashslash",
			Includes:     []string{"**/*.go"},
			Path:         &configurationUri,
			TemplateData: map[string]string{"FooBar": "true"},
		}))
	})

	It("accepts and loads valid configuration with -- comment style", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-file.txt", "style": "dashdash", "includes": ["**/*.go"]}`), nil)

		configuration, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError).To(BeNil())
		Expect(configuration).To(Equal(&core.Configuration{
			HeaderFile:   "some-file.txt",
			CommentStyle: "dashdash",
			Includes:     []string{"**/*.go"},
			Path:         &configurationUri,
		}))
	})

	It("accepts and loads valid configuration with ; comment style", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-file.txt", "style": "semicolon", "includes": ["**/*.go"]}`), nil)

		configuration, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError).To(BeNil())
		Expect(configuration).To(Equal(&core.Configuration{
			HeaderFile:   "some-file.txt",
			CommentStyle: "semicolon",
			Includes:     []string{"**/*.go"},
			Path:         &configurationUri,
		}))
	})

	It("accepts and loads valid configuration with # comment style", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-file.txt", "style": "hash", "includes": ["**/*.go"]}`), nil)

		configuration, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError).To(BeNil())
		Expect(configuration).To(Equal(&core.Configuration{
			HeaderFile:   "some-file.txt",
			CommentStyle: "hash",
			Includes:     []string{"**/*.go"},
			Path:         &configurationUri,
		}))
	})

	It("accepts and loads valid configuration with REM comment style", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-file.txt", "style": "rem", "includes": ["**/*.go"]}`), nil)

		configuration, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError).To(BeNil())
		Expect(configuration).To(Equal(&core.Configuration{
			HeaderFile:   "some-file.txt",
			CommentStyle: "rem",
			Includes:     []string{"**/*.go"},
			Path:         &configurationUri,
		}))
	})

	It("accepts and loads valid configuration with /** comment style", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-file.txt", "style": "slashstarstar", "includes": ["**/*.go"]}`), nil)

		configuration, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError).To(BeNil())
		Expect(configuration).To(Equal(&core.Configuration{
			HeaderFile:   "some-file.txt",
			CommentStyle: "slashstarstar",
			Includes:     []string{"**/*.go"},
			Path:         &configurationUri,
		}))
	})

	It("accepts and loads comment style names in a case-insensitive way", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-file.txt", "style": "slAshStarStaR", "includes": ["**/*.go"]}`), nil)

		configuration, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError).To(BeNil())
		Expect(configuration).To(Equal(&core.Configuration{
			HeaderFile:   "some-file.txt",
			CommentStyle: "slashstarstar",
			Includes:     []string{"**/*.go"},
			Path:         &configurationUri,
		}))
	})

	It("rejects configuration with missing header file", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"style": "SlashStar", "includes": ["**/*.go"]}`), nil)

		_, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError.Error()).
			To(HaveSuffix("Error with field 'headerFile': String length must be greater than or equal to 1"))
	})

	It("rejects configuration with missing comment style", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-header.txt", "includes": ["**/*.go"]}`), nil)

		_, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError.Error()).
			To(HavePrefix("Error with field 'style': style must be one of the following:"))
	})

	It("rejects configuration with invalid comment style", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-header.txt", "style": "invalid", includes": ["**/*.go"]}`), nil)

		_, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError).To(MatchError("invalid character 'i' looking for beginning of object key string"))
	})

	It("rejects configuration with empty includes", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-header.txt", "style": "SlashStar", "includes": []}`), nil)

		_, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError.Error()).To(Equal("Error with field 'includes': Array must have at least 1 items"))
	})

	It("rejects configuration with empty includes", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-header.txt", "style": "SlashSlash", "includes": []}`), nil)

		_, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError.Error()).To(HaveSuffix("Array must have at least 1 items"))
	})

	It("rejects configuration with reserved year parameter", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-header.txt", "style": "SlashSlash", "includes": ["**/*.*"], "data": {"Year": "2019"}}`), nil)

		_, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError.Error()).To(HaveSuffix("Year is a reserved data parameter and cannot be used"))
	})

	It("rejects configuration with reserved year range parameter", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-header.txt", "style": "SlashSlash", "includes": ["**/*.*"], "data": {"YearRange": "2019"}}`), nil)

		_, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError.Error()).To(HaveSuffix("YearRange is a reserved data parameter and cannot be used"))
	})

	It("rejects configuration with reserved start year parameter", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-header.txt", "style": "SlashSlash", "includes": ["**/*.*"], "data": {"StartYear": "2019"}}`), nil)

		_, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError.Error()).To(HaveSuffix("StartYear is a reserved data parameter and cannot be used"))
	})

	It("rejects configuration with reserved end year parameter", func() {
		fileReader.On("Read", configurationUri).
			Return([]byte(`{"headerFile": "some-header.txt", "style": "SlashSlash", "includes": ["**/*.*"], "data": {"EndYear": "2019"}}`), nil)

		_, validationError := loader.ValidateAndLoad(configurationUri)

		Expect(validationError.Error()).To(HaveSuffix("EndYear is a reserved data parameter and cannot be used"))
	})

})

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

type LocalSchemaLoader struct{}

func (*LocalSchemaLoader) Load(schemaLocation string) *json.Schema {
	schema, err := json.NewSchema(json.NewReferenceLoader(schemaLocation))
	if err != nil {
		panic(err)
	}
	return schema
}
