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

package core_test

import (
	. "github.com/fbiville/headache/core"
	"github.com/fbiville/headache/fs"
	"github.com/fbiville/headache/fs_mocks"
	"github.com/fbiville/headache/vcs"
	"os"
	"regexp"
	"testing"
)

func TestHeaderWrite(t *testing.T) {
	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}

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
		[]byte(header+"\n\n"+fileContents)).Return(nil).Once()
	fakeFile.On("Close").Return(nil).Once()

	configuration := ChangeSet{
		HeaderRegex:    getRegex("some multi-line header", "with some text"),
		HeaderContents: header,
		Files:          []vcs.FileChange{{Path: fileName}},
	}

	Run(&configuration, fileSystem)
}

func TestHeaderCommentStyleUpdate(t *testing.T) {
	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}

	oldHeader := "// some multi-line header \n// with some text"
	newHeader := `/*
* some multi-line header
* with some text
*/`
	fakeFile := new(fs_mocks.File)
	commentlessContents := "hello\nworld"
	fileName := "some-file-1"
	fileReader.On("Read", fileName).
		Return([]byte(oldHeader+"\n\n"+commentlessContents), nil).
		Once()
	fileWriter.On("Open", fileName, os.O_WRONLY|os.O_TRUNC, os.ModeAppend).
		Return(fakeFile, nil).
		Once()
	fakeFile.On(
		"Write",
		[]byte(newHeader+"\n\n"+commentlessContents)).Return(nil).Once()
	fakeFile.On("Close").Return(nil).Once()

	configuration := ChangeSet{
		HeaderRegex:    getRegex("some multi-line header", "with some text"),
		HeaderContents: newHeader,
		Files:          []vcs.FileChange{{Path: fileName}},
	}

	Run(&configuration, fileSystem)
}

func TestHeaderDataUpdate(t *testing.T) {
	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}

	oldHeader := "// some multi-line header \n// with some text from Soloing Inc."
	newHeader := "// some multi-line header \n// with some text from Pairing Corp."
	fakeFile := new(fs_mocks.File)
	commentlessContents := "hello\nworld"
	fileName := "some-file-1"
	fileReader.On("Read", fileName).
		Return([]byte(oldHeader+"\n\n"+commentlessContents), nil).
		Once()
	fileWriter.On("Open", fileName, os.O_WRONLY|os.O_TRUNC, os.ModeAppend).
		Return(fakeFile, nil).
		Once()
	fakeFile.On(
		"Write",
		[]byte(newHeader+"\n\n"+commentlessContents)).Return(nil).Once()
	fakeFile.On("Close").Return(nil).Once()

	configuration := ChangeSet{
		HeaderRegex:    getRegexWithParams(map[string]string{"Company": "Soloing Inc."}, "some multi-line header", "with some text from {{.Company}}"),
		HeaderContents: newHeader,
		Files:          []vcs.FileChange{{Path: fileName}},
	}

	Run(&configuration, fileSystem)
}

func TestAutomaticYearInsertion(t *testing.T) {
	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}

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
		[]byte(header+"\n\n"+fileContents)).Return(nil).Once()
	fakeFile.On("Close").Return(nil).Once()

	configuration := ChangeSet{
		HeaderRegex:    getRegex("some multi-line header from {{.Year}}", "with some text"),
		HeaderContents: header,
		Files:          []vcs.FileChange{{Path: fileName, CreationYear: 2022}},
	}

	Run(&configuration, fileSystem)
}

func TestYearIntervalAutomaticInsertion(t *testing.T) {
	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}

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
		[]byte(header+"\n\n"+fileContents)).Return(nil).Once()
	fakeFile.On("Close").Return(nil).Once()

	configuration := ChangeSet{
		HeaderRegex:    getRegex("some multi-line header from {{.Year}}", "with some text"),
		HeaderContents: header,
		Files:          []vcs.FileChange{{Path: fileName, CreationYear: 2022, LastEditionYear: 2034}},
	}

	Run(&configuration, fileSystem)
}

func TestEndYearSkipWhenEqualToStartYear(t *testing.T) {
	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}

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
		[]byte(header+"\n\n"+fileContents)).Return(nil).Once()
	fakeFile.On("Close").Return(nil).Once()

	configuration := ChangeSet{
		HeaderRegex:    getRegex("some multi-line header from {{.Year}}", "with some text"),
		HeaderContents: header,
		Files:          []vcs.FileChange{{Path: fileName, CreationYear: 2022, LastEditionYear: 2022}},
	}

	Run(&configuration, fileSystem)
}

func TestSimilarHeaderReplacement(t *testing.T) {
	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}

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
}

func TestPreserveEarlierThanVersionedYear(t *testing.T) {
	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}

	oldHeader := "// Copyright 2014 ACME"
	newHeader := "// Copyright 2014-2022 ACME"
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
		HeaderRegex: getRegexWithParams(map[string]string{
			"Year":    "{{.Year}}",
			"Company": "ACME",
		}, "Copyright {{.Year}} {{.Company}}"),
		HeaderContents: newHeader,
		Files:          []vcs.FileChange{{Path: fileName, CreationYear: 2016, LastEditionYear: 2022}},
	}

	Run(&configuration, fileSystem)
}

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
