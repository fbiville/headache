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
	. "github.com/fbiville/headache/helper"
	tpl "html/template"
	"os"
	"regexp"
	"strings"
)

type templateResult struct {
	actualContent  string
	detectionRegex *regexp.Regexp
}

func parseTemplate(file string, data map[string]string, style CommentStyle) (*templateResult, error) {
	if err := validateData(data); err != nil {
		return nil, err
	}
	data["Year"] = "{{.Year}}" // template will be parsed a second time, file by file
	rawLines, err := readLines(file)
	if err != nil {
		return nil, err
	}
	commentedLines, err := applyComments(rawLines, style)
	if err != nil {
		return nil, err
	}
	template, err := tpl.New("header").Parse(strings.Join(commentedLines, "\n"))
	if err != nil {
		return nil, err
	}
	builder := &strings.Builder{}
	err = template.Execute(builder, data)
	if err != nil {
		return nil, err
	}
	regex, err := ComputeDetectionRegex(rawLines, data)
	if err != nil {
		return nil, err
	}
	return &templateResult{
		actualContent:  builder.String(),
		detectionRegex: regexp.MustCompile(regex),
	}, nil
}

func validateData(data map[string]string) error {
	if _, ok := data["Year"]; ok {
		return fmt.Errorf("Year is a reserved parameter and is automatically computed.\n" +
			"Please remove it from your configuration")
	}
	return nil
}

func readLines(file string) ([]string, error) {
	openFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer UnsafeClose(openFile)

	lines := make([]string, 0)
	scanner := bufio.NewScanner(openFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func applyComments(lines []string, style CommentStyle) ([]string, error) {
	result := make([]string, 0)
	if openingLine := style.GetOpeningString(); openingLine != "" {
		result = append(result, openingLine)
	}
	for _, line := range lines {
		result = append(result, prependLine(style, line))
	}
	if closingLine := style.GetClosingString(); closingLine != "" {
		result = append(result, closingLine)
	}
	return result, nil
}

func prependLine(style CommentStyle, line string) string {
	comment := style.GetString()
	if line == "" {
		return strings.TrimRight(comment, " ")
	}
	return comment + line
}