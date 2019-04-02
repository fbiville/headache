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

package core

import (
	"github.com/fbiville/headache/fs"
	"github.com/fbiville/headache/fs_mocks"
	"github.com/fbiville/headache/vcs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"regexp"
)

var _ = Describe("Headache", func() {
	var (
		t          GinkgoTInterface
		fileReader *fs_mocks.FileReader
		fileWriter *fs_mocks.FileWriter
		fileSystem *fs.FileSystem
		delimiter  string
	)

	BeforeEach(func() {
		t = GinkgoT()
		fileReader = new(fs_mocks.FileReader)
		fileWriter = new(fs_mocks.FileWriter)
		fileSystem = &fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}
		delimiter = "\n\n"
	})

	AfterEach(func() {
		fileReader.AssertExpectations(t)
		fileWriter.AssertExpectations(t)
	})

	It("writes the header of matched files", func() {
		header := "// some multi-line header \n// with some text"
		fakeFile := new(fs_mocks.File)
		fileContents := "hello\nworld"
		fileName := "some-file-1"
		fileReader.On("Read", fileName).
			Return([]byte(fileContents), nil).
			Once()
		fileWriter.On("Open", fileName, os.O_WRONLY|os.O_TRUNC, os.ModeAppend).
			Return(fakeFile, nil).
			Once()
		fakeFile.On(
			"Write",
			[]byte(header+delimiter+fileContents)).Return(nil).Once()
		fakeFile.On("Close").Return(nil).Once()

		configuration := ChangeSet{
			HeaderRegex:    getRegex("some multi-line header", "with some text"),
			HeaderContents: header,
			Files:          []vcs.FileChange{{Path: fileName}},
		}

		Run(&configuration, fileSystem)
	})

	It("updates the header according to the comment style", func() {
		oldHeader := "// some multi-line header \n// with some text"
		newHeader := `/*
* some multi-line header
* with some text
*/`
		fakeFile := new(fs_mocks.File)
		commentlessContents := "hello\nworld"
		fileName := "some-file-1"
		fileReader.On("Read", fileName).
			Return([]byte(oldHeader+delimiter+commentlessContents), nil).
			Once()
		fileWriter.On("Open", fileName, os.O_WRONLY|os.O_TRUNC, os.ModeAppend).
			Return(fakeFile, nil).
			Once()
		fakeFile.On(
			"Write",
			[]byte(newHeader+delimiter+commentlessContents)).Return(nil).Once()
		fakeFile.On("Close").Return(nil).Once()

		configuration := ChangeSet{
			HeaderRegex:    getRegex("some multi-line header", "with some text"),
			HeaderContents: newHeader,
			Files:          []vcs.FileChange{{Path: fileName}},
		}

		Run(&configuration, fileSystem)
	})

	It("updates the header according to the header parameters", func() {
		oldHeader := "// some multi-line header \n// with some text from Soloing Inc."
		newHeader := "// some multi-line header \n// with some text from Pairing Corp."
		fakeFile := new(fs_mocks.File)
		commentlessContents := "hello\nworld"
		fileName := "some-file-1"
		fileReader.On("Read", fileName).
			Return([]byte(oldHeader+delimiter+commentlessContents), nil).
			Once()
		fileWriter.On("Open", fileName, os.O_WRONLY|os.O_TRUNC, os.ModeAppend).
			Return(fakeFile, nil).
			Once()
		fakeFile.On(
			"Write",
			[]byte(newHeader+delimiter+commentlessContents)).Return(nil).Once()
		fakeFile.On("Close").Return(nil).Once()

		configuration := ChangeSet{
			HeaderRegex:    getRegexWithParams(map[string]string{"Company": "Soloing Inc."}, "some multi-line header", "with some text from {{.Company}}"),
			HeaderContents: newHeader,
			Files:          []vcs.FileChange{{Path: fileName}},
		}

		Run(&configuration, fileSystem)
	})

	It("automatically inserts the year", func() {
		header := "// some multi-line header from 2022\n// with some text"
		fakeFile := new(fs_mocks.File)
		fileContents := "hello\nworld"
		fileName := "some-file-1"
		fileReader.On("Read", fileName).
			Return([]byte(fileContents), nil).
			Once()
		fileWriter.On("Open", fileName, os.O_WRONLY|os.O_TRUNC, os.ModeAppend).
			Return(fakeFile, nil).
			Once()
		fakeFile.On(
			"Write",
			[]byte(header+delimiter+fileContents)).Return(nil).Once()
		fakeFile.On("Close").Return(nil).Once()

		configuration := ChangeSet{
			HeaderRegex:    getRegex("some multi-line header from {{.Year}}", "with some text"),
			HeaderContents: header,
			Files:          []vcs.FileChange{{Path: fileName, CreationYear: 2022}},
		}

		Run(&configuration, fileSystem)
	})

	It("automatically inserts the year interval", func() {
		header := "// some multi-line header from 2022-2034\n// with some text"
		fakeFile := new(fs_mocks.File)
		fileContents := "hello\nworld"
		fileName := "some-file-1"
		fileReader.On("Read", fileName).
			Return([]byte(fileContents), nil).
			Once()
		fileWriter.On("Open", fileName, os.O_WRONLY|os.O_TRUNC, os.ModeAppend).
			Return(fakeFile, nil).
			Once()
		fakeFile.On(
			"Write",
			[]byte(header+delimiter+fileContents)).Return(nil).Once()
		fakeFile.On("Close").Return(nil).Once()

		configuration := ChangeSet{
			HeaderRegex:    getRegex("some multi-line header from {{.Year}}", "with some text"),
			HeaderContents: header,
			Files:          []vcs.FileChange{{Path: fileName, CreationYear: 2022, LastEditionYear: 2034}},
		}

		Run(&configuration, fileSystem)
	})

	It("automatically prevents the end year insertion if it's the same as the start year", func() {
		header := "// some multi-line header from 2022\n// with some text"
		fakeFile := new(fs_mocks.File)
		fileContents := "hello\nworld"
		fileName := "some-file-1"
		fileReader.On("Read", fileName).
			Return([]byte(fileContents), nil).
			Once()
		fileWriter.On("Open", fileName, os.O_WRONLY|os.O_TRUNC, os.ModeAppend).
			Return(fakeFile, nil).
			Once()
		fakeFile.On(
			"Write",
			[]byte(header+delimiter+fileContents)).Return(nil).Once()
		fakeFile.On("Close").Return(nil).Once()

		configuration := ChangeSet{
			HeaderRegex:    getRegex("some multi-line header from {{.Year}}", "with some text"),
			HeaderContents: header,
			Files:          []vcs.FileChange{{Path: fileName, CreationYear: 2022, LastEditionYear: 2022}},
		}

		Run(&configuration, fileSystem)
	})

	It("matches similar header and replaces it", func() {
		oldHeader := `/*
 *   Some Header 2022 -   2023 and stuff .
 *
 */`
		newHeader := "// some header 2022-2024 and stuff"
		fakeFile := new(fs_mocks.File)
		fileContents := "hello\nworld"
		fileName := "some-file-1"
		fileReader.On("Read", fileName).
			Return([]byte(oldHeader+"\n\n"+fileContents), nil).
			Once()
		fileWriter.On("Open", fileName, os.O_WRONLY|os.O_TRUNC, os.ModeAppend).
			Return(fakeFile, nil).
			Once()
		fakeFile.On(
			"Write",
			[]byte(newHeader+"\n\n"+fileContents)).Return(nil).Once()
		fakeFile.On("Close").Return(nil).Once()

		configuration := ChangeSet{
			HeaderRegex:    getRegexWithParams(map[string]string{"Year": "2022"}, "some header {{.Year}} and stuff"),
			HeaderContents: newHeader,
			Files:          []vcs.FileChange{{Path: fileName, CreationYear: 2022, LastEditionYear: 2024}},
		}

		Run(&configuration, fileSystem)
	})

	It("preserves existing start year when it is lower than the configured one", func() {
		oldHeader := "// Copyright 2014 ACME"
		newHeader := "// Copyright 2014-2022 ACME"
		fakeFile := new(fs_mocks.File)
		fileContents := "hello\nworld"
		fileName := "some-file-1"
		fileReader.On("Read", fileName).
			Return([]byte(oldHeader+delimiter+fileContents), nil).
			Once()
		fileWriter.On("Open", fileName, os.O_WRONLY|os.O_TRUNC, os.ModeAppend).
			Return(fakeFile, nil).
			Once()
		fakeFile.On(
			"Write",
			[]byte(newHeader+delimiter+fileContents)).Return(nil).Once()
		fakeFile.On("Close").Return(nil).Once()

		configuration := ChangeSet{
			HeaderRegex: getRegexWithParams(map[string]string{
				"Year":    "{{.Year}}",
				"Company": "ACME",
			}, "Copyright {{.Year}} {{.Company}}"),
			HeaderContents: newHeader,
			Files:          []vcs.FileChange{{Path: fileName, CreationYear: 2016, LastEditionYear: 2022}},
		}

		Run(&configuration, fileSystem)
	})

	It("replaces single future copyright header date with single commit year", func() {
		change := vcs.FileChange{
			Path:            "pkg/fileutils/abs_test.go",
			CreationYear:    2018,
			LastEditionYear: 2018,
		}
		header := `/*
 * Copyright 2019 The original author or authors
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
 */`

		startYear, endYear, err := computeCopyrightYears(&change, header)

		Expect(err).NotTo(HaveOccurred())
		Expect(startYear).To(Equal(2018))
		Expect(endYear).To(Equal(2018))
	})
})

func getRegex(headerLines ...string) *regexp.Regexp {
	return getRegexWithParams(map[string]string{}, headerLines...)
}

func getRegexWithParams(params map[string]string, headerLines ...string) *regexp.Regexp {
	regex, err := ComputeDetectionRegex(headerLines, params)
	if err != nil {
		panic(err)
	}
	return regexp.MustCompile(regex)
}
