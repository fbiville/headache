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
	"bufio"
	"fmt"
	"github.com/fbiville/headache/versioning"
	"github.com/sergi/go-diff/diffmatchpatch"
	tpl "html/template"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type VcsChangeGetter func(versioning.Vcs, string, string) (error, []versioning.FileChange)

func DryRun(config *configuration) (string, error) {
	file, err := ioutil.TempFile("", "headache-dry-run")
	if err != nil {
		return "", err
	}
	config.writer = bufio.NewWriter(file)
	defer file.Close()
	insertInMatchedFiles(config)
	return file.Name(), nil
}

func Run(config *configuration) {
	insertInMatchedFiles(config)
}

func insertInMatchedFiles(config *configuration) {
	for _, change := range config.vcsChanges {
		bytes, err := ioutil.ReadFile(change.Path)
		if err != nil {
			panic(err)
		}

		fileContents := string(bytes)
		matchLocation := config.HeaderRegex.FindStringIndex(fileContents)
		existingHeader := ""
		if matchLocation != nil {
			existingHeader = fileContents[matchLocation[0]:matchLocation[1]+1]
			fileContents = strings.TrimLeft(fileContents[:matchLocation[0]]+fileContents[matchLocation[1]:], "\n")
		}

		finalHeaderContent, err := insertYears(config.HeaderContents, change, existingHeader)
		if err != nil {
			panic(err)
		}
		newContents := append([]byte(fmt.Sprintf("%s%s", finalHeaderContent, "\n\n")), []byte(fileContents)...)
		writeToFile(config, change, newContents)
	}
}

func insertYears(template string, change versioning.FileChange, existingHeader string) (string, error) {
	t, err := tpl.New("header-second-pass").Parse(template)
	if err != nil {
		return "", err
	}
	data := make(map[string]string)
	data["Year"] = computeCopyrightYears(change, existingHeader)
	builder := &strings.Builder{}
	err = t.Execute(builder, data)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func computeCopyrightYears(change versioning.FileChange, existingHeader string) string {
	regex := regexp.MustCompile(`(\d{4})(?:\s*-\s*(\d{4}))?`)
	matches := regex.FindStringSubmatch(existingHeader)
	creationYear := change.CreationYear
	if len(matches) > 2 {
		start, _ := strconv.Atoi(matches[1])
		creationYear = start
	}
	year := strconv.Itoa(creationYear)
	lastEditionYear := change.LastEditionYear
	if lastEditionYear != 0 && lastEditionYear != creationYear {
		year += fmt.Sprintf("-%d", lastEditionYear)
	}
	return year
}

func writeToFile(config *configuration, change versioning.FileChange, newContents []byte) {
	file := change.Path
	var writer = config.writer
	if writer != nil {
		appendToDryRunFile(writer, change, newContents)
	} else {
		alterSourceFile(file, newContents)
	}
}

func appendToDryRunFile(writer io.Writer, change versioning.FileChange, newContents []byte) {
	diff := computeDiff(change, newContents)
	if diff != "" {
		bufferedWriter := bufio.NewWriter(writer)
		bufferedWriter.Write([]byte(fmt.Sprintf("file:%s", change.Path)))
		bufferedWriter.Write([]byte("\n---\n"))
		bufferedWriter.Write([]byte(prefixLines(diff, "\t")))
		bufferedWriter.Write([]byte("---\n"))
		bufferedWriter.Flush()
	}

}
func prefixLines(content string, prefix string) string {
	builder := strings.Builder{}
	newline := "\n"
	for _, line := range strings.Split(content, newline) {
		builder.WriteString(prefix)
		builder.WriteString(line)
		builder.WriteString(newline)
	}
	return builder.String()
}

func computeDiff(change versioning.FileChange, newContents []byte) string {
	diffTool := diffmatchpatch.New()
	differences := diffTool.DiffMain(change.ReferenceContent, string(newContents), false)
	if onlyEqualDiffs(differences) {
		return ""
	}
	return diffTool.DiffPrettyText(differences)
}

func onlyEqualDiffs(diffs []diffmatchpatch.Diff) bool {
	for _, diff := range diffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			return false
		}
	}
	return true
}

func alterSourceFile(file string, newContents []byte) {
	openFile, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	_, err = openFile.Write(newContents)
	openFile.Close()
	if err != nil {
		panic(err)
	}
}
