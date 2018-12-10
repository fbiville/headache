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
	"github.com/fbiville/headache/versioning"
	"github.com/mattn/go-zglob"
	tpl "html/template"
	"io"
	"os"
	"regexp"
	"strings"
)

type Configuration struct {
	HeaderFile   string            `json:"headerFile"`
	CommentStyle string            `json:"style"`
	Includes     []string          `json:"includes"`
	Excludes     []string          `json:"excludes"`
	TemplateData map[string]string `json:"data"`
}

type configuration struct {
	HeaderContents string
	HeaderRegex    *regexp.Regexp
	vcsChanges     []versioning.FileChange
	writer         io.Writer
}

func ParseConfiguration(config Configuration) (*configuration, error) {
	return parseConfiguration(config, versioning.GetVcsChanges)
}

func parseConfiguration(config Configuration,
	getRevisionChanges func(versioning.Vcs, string) ([]versioning.FileChange, error)) (*configuration, error) {

	contents, err := parseTemplate(config.HeaderFile, config.TemplateData, newCommentStyle(config.CommentStyle))
	if err != nil {
		return nil, err
	}

	changes, err := getFileChanges(config, getRevisionChanges)
	if err != nil {
		return nil, err
	}

	return &configuration{
		HeaderContents: contents.actualContent,
		HeaderRegex:    contents.detectionRegex,
		vcsChanges:     changes,
	}, nil
}

func getFileChanges(config Configuration,
	getRevisionChanges func(versioning.Vcs, string) ([]versioning.FileChange, error)) ([]versioning.FileChange, error) {
	vcs := versioning.Git{}
	revision, err := versioning.GetLatestExecutionRevision(vcs)
	if err != nil {
		return nil, err
	}
	var changes []versioning.FileChange
	if revision == "" {
		changes, err = matchAllFiles(config.Includes, config.Excludes)
	} else {
		changes, err = matchChangedFiles(revision, config, vcs, getRevisionChanges)
	}
	if err != nil {
		return nil, err
	}
	return versioning.AddMetadata(vcs, changes, revision)
}

func matchAllFiles(includes []string, excludes []string) ([]versioning.FileChange, error) {
	result := make([]versioning.FileChange, 0)
	for _, includePattern := range includes {
		matches, err := zglob.Glob(includePattern)
		if err != nil {
			return nil, err
		}
		for _, matchedPath := range matches {
			if !isExcluded(matchedPath, excludes) {
				result = append(result, versioning.FileChange{
					Path: matchedPath,
				})
			}
		}

	}
	return result, nil
}

func matchChangedFiles(sha string, config Configuration, vcs versioning.Vcs, getVersioningChanges func(versioning.Vcs, string) ([]versioning.FileChange, error)) ([]versioning.FileChange, error) {
	fileChanges, err := getVersioningChanges(vcs, sha)
	if err != nil {
		return nil, err
	}
	return filterFiles(fileChanges, config.Includes, config.Excludes), nil
}

func filterFiles(changes []versioning.FileChange, includes []string, excludes []string) []versioning.FileChange {
	result := make([]versioning.FileChange, 0)
	for _, change := range changes {
		if match(change.Path, includes, excludes) {
			result = append(result, change)
		}
	}
	return result
}

func match(path string, includes []string, excludes []string) bool {
	return Match(path, includes) && !isExcluded(path, excludes)
}

func isExcluded(path string, excludes []string) bool {
	return !IsFile(path) || Match(path, excludes)
}

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
	regex, err := computeDetectionRegex(rawLines, data)
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
	if style.opening() {
		result = append(result, style.open())
	}
	for _, line := range lines {
		result = append(result, style.apply(line))
	}
	if style.closing() {
		result = append(result, style.close())
	}
	return result, nil
}

func computeDetectionRegex(lines []string, data map[string]string) (string, error) {
	regex := regexLines(lines)
	return injectDataRegex(strings.Join(regex, ""), data)
}

func injectDataRegex(result string, data map[string]string) (string, error) {
	template, err := tpl.New("header-regex").Parse(result)
	if err != nil {
		return "", err
	}
	builder := &strings.Builder{}
	err = template.Execute(builder, regexValues(&data))
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func regexLines(lines []string) []string {
	result := make([]string, 0)
	result = append(result, `(?im)(?:\/\*\n)?`)
	for _, line := range lines {
		result = append(result, fmt.Sprintf(`%s\Q%s\E[ \t\.]*\n?`, `(?:\/{2}| \*)[ \t]*`, line))
	}
	result = append(result, `(?:(?:\/{2}| \*) ?\n)*`)
	result = append(result, `(?: \*\/)?`)
	return result
}

func regexValues(data *map[string]string) *map[string]string {
	for k := range *data {
		(*data)[k] = "\\E.*\\Q"
	}
	return data
}

