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
	"github.com/fbiville/headache/helper"
	"github.com/fbiville/headache/versioning"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHeaderWrite(t *testing.T) {
	I := NewGomegaWithT(t)
	testWriter := temporaryFile("headache-test-run", os.O_RDWR|os.O_CREATE)
	defer helper.UnsafeClose(testWriter.file)
	defer helper.UnsafeDelete(testWriter.file)

	newHeader := `// some multi-line header
// with some text`
	regex, _ := ComputeDetectionRegex([]string{"some multi-line header", "with some text"}, map[string]string{})
	configuration := configuration{
		HeaderRegex:    regexp.MustCompile(regex),
		HeaderContents: newHeader,
		vcsChanges:     []versioning.FileChange{{Path: "../fixtures/hello_world.txt", ReferenceContent: ""}},
	}

	insertInMatchedFiles(&configuration, testWriter)

	I.Expect(readFile(testWriter.file.Name())).To(Equal(`// some multi-line header
// with some text

hello
world
---
`))
}

func TestHeaderDoesNotWriteTwice(t *testing.T) {
	I := NewGomegaWithT(t)
	testWriter := temporaryFile("headache-test-run", os.O_RDWR|os.O_CREATE)
	defer helper.UnsafeClose(testWriter.file)
	defer helper.UnsafeDelete(testWriter.file)

	regex, _ := ComputeDetectionRegex([]string{"some multi-line header", "with some text"}, map[string]string{})
	sourceFile := "../fixtures/hello_world_with_header.txt"
	configuration := configuration{
		HeaderRegex: regexp.MustCompile(regex),
		HeaderContents: `// some multi-line header
// with some text`,
		vcsChanges: []versioning.FileChange{{Path: sourceFile, ReferenceContent: ""}},
	}

	insertInMatchedFiles(&configuration, testWriter)

	I.Expect(readFile(testWriter.file.Name())).To(Equal(readFile(sourceFile) + "\n---\n"),
		"it should rewrite the file as is")
}

func TestHeaderCommentUpdate(t *testing.T) {
	I := NewGomegaWithT(t)
	testWriter := temporaryFile("headache-test-run", os.O_RDWR|os.O_CREATE)
	defer helper.UnsafeClose(testWriter.file)
	defer helper.UnsafeDelete(testWriter.file)

	regex, _ := ComputeDetectionRegex([]string{"some multi-line header", "with some text"}, map[string]string{})
	configuration := configuration{
		HeaderRegex: regexp.MustCompile(regex),
		HeaderContents: `/*
 * some multi-line header
 * with some text
 */`,
		vcsChanges: []versioning.FileChange{{Path: "../fixtures/hello_world_with_header.txt", ReferenceContent: ""}},
	}

	insertInMatchedFiles(&configuration, testWriter)

	I.Expect(readFile(testWriter.file.Name())).To(Equal(`/*
 * some multi-line header
 * with some text
 */

hello
world
---
`), "it should rewrite the file with slashstar style")
}

func TestHeaderDataUpdate(t *testing.T) {
	I := NewGomegaWithT(t)
	testWriter := temporaryFile("headache-test-run", os.O_RDWR|os.O_CREATE)
	defer helper.UnsafeClose(testWriter.file)
	defer helper.UnsafeDelete(testWriter.file)

	regex, _ := ComputeDetectionRegex([]string{"some multi-line header 2017", "with some text from {{.Company}}"},
		map[string]string{
			"Company": "Pairing Corp",
		})
	configuration := configuration{
		HeaderRegex: regexp.MustCompile(regex),
		HeaderContents: `// some multi-line header 2017
// with some text from Pairing Corp`,
		vcsChanges: []versioning.FileChange{{Path: "../fixtures/hello_world_with_parameterized_header.txt", ReferenceContent: ""}},
	}

	insertInMatchedFiles(&configuration, testWriter)

	file := readFile(testWriter.file.Name())
	I.Expect(file).To(Equal(`// some multi-line header 2017
// with some text from Pairing Corp

hello
world
---
`))
}

func TestInsertCreationYearAutomatically(t *testing.T) {
	I := NewGomegaWithT(t)
	testWriter := temporaryFile("headache-test-run", os.O_RDWR|os.O_CREATE)
	defer helper.UnsafeClose(testWriter.file)
	defer helper.UnsafeDelete(testWriter.file)

	regex, _ := ComputeDetectionRegex([]string{"some multi-line header {{.Year}}", "with some text from {{.Company}}"},
		map[string]string{
			"Year":    "{{.Year}}",
			"Company": "Pairing Corp",
		})
	configuration := configuration{
		HeaderRegex: regexp.MustCompile(regex),
		HeaderContents: `// some multi-line header {{.Year}}
// with some text from Pairing Corp`,
		vcsChanges: []versioning.FileChange{{
			Path:             "../fixtures/hello_world.txt",
			CreationYear:     2022,
			ReferenceContent: "",
		}},
	}

	insertInMatchedFiles(&configuration, testWriter)

	I.Expect(readFile(testWriter.file.Name())).To(Equal(`// some multi-line header 2022
// with some text from Pairing Corp

hello
world
---
`))
}

func TestInsertCreationAndLastEditionYearsAutomatically(t *testing.T) {
	I := NewGomegaWithT(t)
	testWriter := temporaryFile("headache-test-run", os.O_RDWR|os.O_CREATE)
	defer helper.UnsafeClose(testWriter.file)
	defer helper.UnsafeDelete(testWriter.file)

	regex, _ := ComputeDetectionRegex([]string{"some multi-line header {{.Year}}", "with some text from {{.Company}}"},
		map[string]string{
			"Year":    "{{.Year}}",
			"Company": "Pairing Corp",
		})
	configuration := configuration{
		HeaderRegex: regexp.MustCompile(regex),
		HeaderContents: `// some multi-line header {{.Year}}
// with some text from Pairing Corp`,
		vcsChanges: []versioning.FileChange{{
			Path:             "../fixtures/hello_world.txt",
			CreationYear:     2022,
			LastEditionYear:  2034,
			ReferenceContent: "",
		}},
	}

	insertInMatchedFiles(&configuration, testWriter)

	I.Expect(readFile(testWriter.file.Name())).To(Equal(`// some multi-line header 2022-2034
// with some text from Pairing Corp

hello
world
---
`))
}

func TestDoesNotInsertLastEditionYearWhenEqualToCreationYear(t *testing.T) {
	I := NewGomegaWithT(t)
	testWriter := temporaryFile("headache-test-run", os.O_RDWR|os.O_CREATE)
	defer helper.UnsafeClose(testWriter.file)
	defer helper.UnsafeDelete(testWriter.file)

	regex, _ := ComputeDetectionRegex([]string{"some multi-line header {{.Year}}", "with some text from {{.Company}}"},
		map[string]string{
			"Year":    "{{.Year}}",
			"Company": "Pairing Corp",
		})
	configuration := configuration{
		HeaderRegex: regexp.MustCompile(regex),
		HeaderContents: `// some multi-line header {{.Year}}
// with some text from Pairing Corp`,
		vcsChanges: []versioning.FileChange{{
			Path:             "../fixtures/hello_world.txt",
			CreationYear:     2022,
			LastEditionYear:  2022,
			ReferenceContent: "",
		}},
	}

	insertInMatchedFiles(&configuration, testWriter)

	I.Expect(readFile(testWriter.file.Name())).To(Equal(`// some multi-line header 2022
// with some text from Pairing Corp

hello
world
---
`))
}

func TestHeaderDryRunOnSeveralFiles(t *testing.T) {
	I := NewGomegaWithT(t)
	testWriter := temporaryFile("headache-test-run", os.O_RDWR|os.O_CREATE)
	defer helper.UnsafeClose(testWriter.file)
	defer helper.UnsafeDelete(testWriter.file)

	regex, _ := ComputeDetectionRegex([]string{"some header {{.Year}}"},
		map[string]string{
			"Year": "{{.Year}}",
		})
	insertInMatchedFiles(&configuration{
		HeaderRegex:    regexp.MustCompile(regex),
		HeaderContents: "// some header {{.Year}}",
		vcsChanges: []versioning.FileChange{{
			Path:             "../fixtures/hello_world.txt",
			CreationYear:     2022,
			LastEditionYear:  2022,
			ReferenceContent: "",
		}, {
			Path:             "../fixtures/bonjour_world.txt",
			CreationYear:     2019,
			LastEditionYear:  2021,
			ReferenceContent: "",
		}},
	}, testWriter)

	s := readFile(testWriter.file.Name())
	I.Expect(s).To(Equal(`// some header 2022

hello
world
---
// some header 2019-2021

bonjour
le world
---
`))

}

func TestSimilarHeaderReplacement(t *testing.T) {
	I := NewGomegaWithT(t)
	testWriter := temporaryFile("headache-test-run", os.O_RDWR|os.O_CREATE)
	defer helper.UnsafeClose(testWriter.file)
	defer helper.UnsafeDelete(testWriter.file)

	regex, _ := ComputeDetectionRegex([]string{"some header {{.Year}} and stuff"},
		map[string]string{
			"Year": "{{.Year}}",
		})
	insertInMatchedFiles(&configuration{
		HeaderRegex:    regexp.MustCompile(regex),
		HeaderContents: "// some header {{.Year}} and stuff",
		vcsChanges: []versioning.FileChange{{
			Path:             "../fixtures/hello_world_similar.txt",
			CreationYear:     2022,
			LastEditionYear:  2022,
			ReferenceContent: "",
		}},
	}, testWriter)

	I.Expect(readFile(testWriter.file.Name())).To(Equal(`// some header 2022 and stuff

hello
world
---
`))
}

func TestPreserveYear(t *testing.T) {
	I := NewGomegaWithT(t)
	testWriter := temporaryFile("headache-test-run", os.O_RDWR|os.O_CREATE)
	defer helper.UnsafeClose(testWriter.file)
	defer helper.UnsafeDelete(testWriter.file)

	regex, _ := ComputeDetectionRegex([]string{"Copyright {{.Year}} {{.Company}}"},
		map[string]string{
			"Year":    "{{.Year}}",
			"Company": "ACME",
		})
	insertInMatchedFiles(&configuration{
		HeaderRegex:    regexp.MustCompile(regex),
		HeaderContents: "// some header {{.Year}} {{.Company}}",
		vcsChanges: []versioning.FileChange{{
			Path:             "../fixtures/hello_world_2014.txt",
			CreationYear:     2016,
			LastEditionYear:  2022,
			ReferenceContent: "",
		}},
	}, testWriter)

	I.Expect(readFile(testWriter.file.Name())).To(Equal(`// some header 2014-2022 

Hello world!!
---
`))
}

func readFile(file string) string {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

type AppendingFile struct {
	file *os.File
}

// open-close are noops because they are managed within the tests
func (osw *AppendingFile) Open(name string, mask int, permissions os.FileMode) (*os.File, error) {
	return osw.file, nil
}
func (*AppendingFile) Close(file *os.File) {
	_, err := file.WriteString("\n---\n")
	if err != nil {
		panic(err)
	}
}

func temporaryFile(name string, umask int) *AppendingFile {
	rand.Seed(time.Now().UTC().UnixNano())
	tempDirectory := os.TempDir()
	if !strings.HasSuffix(tempDirectory, "/") {
		tempDirectory += "/"
	}
	file, err := os.OpenFile(tempDirectory+name+strconv.Itoa(rand.Int()), umask, 0644)
	if err != nil {
		panic(err)
	}
	return &AppendingFile{
		file: file,
	}
}
