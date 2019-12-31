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
	"fmt"
	"github.com/fbiville/headache/fs"
	"github.com/fbiville/headache/vcs"
	tpl "html/template"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Headache struct {
	Fs *fs.FileSystem
}

func (headache *Headache) Run(config *ChangeSet) {
	currentHeaderDetectionRegex := config.HeaderRegex
	newHeaderTemplate := config.HeaderContents
	for _, file := range config.Files {
		headache.UpdateFile(file, currentHeaderDetectionRegex, newHeaderTemplate)
	}
}

func (headache *Headache) UpdateFile(change vcs.FileChange, currentHeaderDetectionRegex *regexp.Regexp, newHeaderTemplate string) {
	path := change.Path
	bytes, err := headache.Fs.FileReader.Read(path)
	if err != nil {
		log.Fatalf("headache execution error, cannot read file %s\n\t%v", path, err)
	}

	fileContents := string(bytes)
	matchLocation := currentHeaderDetectionRegex.FindStringIndex(fileContents)
	existingHeader := ""
	if matchLocation != nil {
		existingHeader = fileContents[matchLocation[0]:matchLocation[1]]
		fileContents = strings.TrimLeft(fileContents[:matchLocation[0]]+fileContents[matchLocation[1]:], "\n")
	}

	finalHeaderContent, err := insertYears(newHeaderTemplate, &change, existingHeader)
	if err != nil {
		log.Fatalf("headache execution error, cannot parse header for file %s\n\t%v", path, err)
	}
	newContents := append([]byte(fmt.Sprintf("%s%s", finalHeaderContent, "\n\n")), []byte(fileContents)...)
	headache.writeToFile(path, newContents)
}

func insertYears(template string, change *vcs.FileChange, existingHeader string) (string, error) {
	t, err := tpl.New("header-second-pass").Parse(template)
	if err != nil {
		return "", err
	}
	data := make(map[string]string)
	startYear, endYear, err := ComputeCopyrightYears(change, existingHeader)
	if err != nil {
		return "", err
	}
	data["YearRange"] = strconv.Itoa(startYear)
	data["StartYear"] = strconv.Itoa(startYear)
	data["EndYear"] = strconv.Itoa(endYear)
	if startYear != endYear {
		data["YearRange"] = fmt.Sprintf("%d-%d", startYear, endYear)
	}
	builder := &strings.Builder{}
	err = t.Execute(builder, data)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

// visible for testing
func ComputeCopyrightYears(change *vcs.FileChange, existingHeader string) (int, int, error) {
	regex := regexp.MustCompile(`(\d{4})(?:\s*-\s*(\d{4}))?`)
	matches := regex.FindStringSubmatch(existingHeader)
	creationYear := change.CreationYear
	if len(matches) > 2 {
		startYearInHeader, err := strconv.Atoi(matches[1])
		if err != nil {
			return 0, 0, err
		}
		if startYearInHeader < creationYear {
			creationYear = startYearInHeader
		}
	}
	lastEditionYear := change.LastEditionYear
	if lastEditionYear != 0 && lastEditionYear != creationYear {
		return creationYear, lastEditionYear, nil
	}
	return creationYear, creationYear, nil
}

func (headache *Headache) writeToFile(path string, newContents []byte) {
	file, err := headache.Fs.FileWriter.Open(path, os.O_WRONLY|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		log.Fatalf("headache execution error, cannot open file %s\n\t%v", path, err)
	}
	defer fs.UnsafeClose(file)
	err = file.Write(newContents)
	if err != nil {
		log.Fatalf("headache execution error, cannot write to file %s\n\t%v", path, err)
	}
}
