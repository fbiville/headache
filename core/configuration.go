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
	"github.com/fbiville/headache/fs"
	"github.com/fbiville/headache/helper"
	"github.com/fbiville/headache/vcs"
	"regexp"
)

func DefaultSystemConfiguration() SystemConfiguration {
	return SystemConfiguration{
		VersioningClient: &vcs.Client{
			Vcs: vcs.Git{},
		},
		FileSystem: fs.DefaultFileSystem(),
		Clock:      helper.SystemClock{},
	}
}

type SystemConfiguration struct {
	VersioningClient vcs.VersioningClient
	FileSystem       fs.FileSystem
	Clock            helper.Clock
}

type Configuration struct {
	HeaderFile   string            `json:"headerFile"`
	CommentStyle string            `json:"style"`
	Includes     []string          `json:"includes"`
	Excludes     []string          `json:"excludes"`
	TemplateData map[string]string `json:"data"`
}

type ChangeSet struct {
	HeaderContents string
	HeaderRegex    *regexp.Regexp
	Files          []vcs.FileChange
}

func ParseConfiguration(
	config Configuration,
	sysConfig SystemConfiguration,
	tracker fs.ExecutionTracker,
	pathMatcher fs.PathMatcher) (*ChangeSet, error) {

	// TODO: config.HeaderFile may have changed since last run - this may cause e.g. to not find the previous header file!
	headerContents, err := tracker.ReadLinesAtLastExecutionRevision(config.HeaderFile)
	if err != nil {
		return nil, err
	}

	// TODO: config.TemplateData may have changed since last run -- this may change the detection regex!
	// note: comment style does not matter as the detection regex is designed to be insensitive to the style in use
	contents, err := ParseTemplate(headerContents, config.TemplateData, ParseCommentStyle(config.CommentStyle))
	if err != nil {
		return nil, err
	}

	changes, err := getFileChanges(config, sysConfig, tracker, pathMatcher)
	if err != nil {
		return nil, err
	}

	return &ChangeSet{
		HeaderContents: contents.actualContent,
		HeaderRegex:    contents.detectionRegex,
		Files:          changes,
	}, nil
}

func getFileChanges(config Configuration,
	sysConfig SystemConfiguration,
	tracker fs.ExecutionTracker,
	pathMatcher fs.PathMatcher) ([]vcs.FileChange, error) {

	versioningClient := sysConfig.VersioningClient
	fileSystem := sysConfig.FileSystem
	revision, err := tracker.GetLastExecutionRevision()
	if err != nil {
		return nil, err
	}
	var changes []vcs.FileChange
	// TODO: centralize this check between here and ExecutionTracker#ReadLinesAtLastExecutionRevision
	if revision == "" {
		changes, err = pathMatcher.ScanAllFiles(config.Includes, config.Excludes, fileSystem)
		if err != nil {
			return nil, err
		}
	} else {
		fileChanges, err := versioningClient.GetChanges(revision)
		if err != nil {
			return nil, err
		}
		changes = pathMatcher.MatchFiles(fileChanges, config.Includes, config.Excludes, fileSystem)
	}
	return versioningClient.AddMetadata(changes, sysConfig.Clock)
}
