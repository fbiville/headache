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
	tpl "html/template"
	"regexp"
	"strings"
)

type templateResult struct {
	actualContent  string
	detectionRegex *regexp.Regexp
}

func ParseTemplate(versionedHeader *VersionedHeaderTemplate, style CommentStyle) (*templateResult, error) {
	currentData := injectReservedYearParameter(versionedHeader.Current.Data)
	commentedLines, err := applyComments(versionedHeader.Current.Lines, style)
	if err != nil {
		return nil, err
	}
	template, err := tpl.New("header").Parse(strings.Join(commentedLines, "\n"))
	if err != nil {
		return nil, err
	}
	builder := &strings.Builder{}
	err = template.Execute(builder, currentData)
	if err != nil {
		return nil, err
	}

	previousData := injectReservedYearParameter(versionedHeader.Previous.Data)
	regex, err := ComputeDetectionRegex(versionedHeader.Previous.Lines, previousData)
	if err != nil {
		return nil, err
	}
	return &templateResult{
		actualContent:  builder.String(),
		detectionRegex: regexp.MustCompile(regex),
	}, nil
}

// injects reserved parameter into template data map
// the template will be parsed a second time, file by file, with the computed .Year value
func injectReservedYearParameter(currentData map[string]string) map[string]string {
	currentData["Year"] = "{{.Year}}"
	return currentData
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
