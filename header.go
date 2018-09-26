package main

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-zglob"
	"io/ioutil"
	"os"
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
			fileContents = fileContents[:matchLocation[0]] + fileContents[matchLocation[1]:]
		}

		newContents := append([]byte(fmt.Sprintf("%s%s", config.HeaderContents, "\n\n")), []byte(fileContents)...)
		writeToFile(config, file, newContents)
	}
}

func writeToFile(config *configuration, file string, newContents []byte) {
	var writer = config.writer
	if writer == nil {
		openFile, err := os.OpenFile(file, os.O_WRONLY, os.ModeAppend)
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
