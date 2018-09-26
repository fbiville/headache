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

package main

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-zglob"
	"io/ioutil"
	"os"
	"strings"
)

func InsertHeader(config *configuration) {
	for _, includePattern := range config.Includes {
		matches, err := zglob.Glob(includePattern)
		if err != nil {
			panic(err)
		}
		insertInMatchedFiles(config, exclude(matches, config.Excludes))
	}
}

func exclude(strings []string, exclusionPatterns []string) []string {
	result := strings[:0]
	for _, str := range strings {
		if !matches(str, exclusionPatterns) {
			result = append(result, str)
		}
	}
	return result
}

func matches(str string, exclusionPatterns []string) bool {
	for _, exclusionPattern := range exclusionPatterns {
		matched, _ := zglob.Match(exclusionPattern, str)
		if matched {
			return true
		}
	}
	return false
}

func insertInMatchedFiles(config *configuration, files []string) {
	for _, file := range files {
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		fileContents := string(bytes)
		matchLocation := config.HeaderRegex.FindStringIndex(fileContents)
		if matchLocation != nil {
			fileContents = strings.TrimLeft(fileContents[:matchLocation[0]] + fileContents[matchLocation[1]:], "\n")
		}

		newContents := append([]byte(fmt.Sprintf("%s%s", config.HeaderContents, "\n\n")), []byte(fileContents)...)
		writeToFile(config, file, newContents)
	}
}

func writeToFile(config *configuration, file string, newContents []byte) {
	var writer = config.writer
	if writer == nil {
		openFile, err := os.OpenFile(file, os.O_WRONLY | os.O_TRUNC, os.ModeAppend)
		if err != nil {
			panic(err)
		}
		_, err = openFile.Write(newContents)
		openFile.Close()
		if err != nil {
			panic(err)
		}
	} else {
		bufferedWriter := bufio.NewWriter(writer)
		bufferedWriter.Write(newContents)
		bufferedWriter.Flush()
	}
}
