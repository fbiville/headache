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
	. "github.com/fbiville/headache/helper"
	"github.com/fbiville/headache/versioning"
	"github.com/mattn/go-zglob"
	"regexp"
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
}

func ParseConfiguration(config Configuration) (*configuration, error) {
	return parseConfiguration(config, versioning.GetVcsChanges)
}

func parseConfiguration(config Configuration,
	getRevisionChanges func(versioning.Vcs, string) ([]versioning.FileChange, error)) (*configuration, error) {

	contents, err := ParseTemplate(config.HeaderFile, config.TemplateData, ParseCommentStyle(config.CommentStyle))
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

func matchChangedFiles(sha string, config Configuration, vcs versioning.Vcs,
	getVersioningChanges func(versioning.Vcs, string) ([]versioning.FileChange, error)) ([]versioning.FileChange, error) {
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
