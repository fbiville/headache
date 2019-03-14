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
	"log"
	"regexp"
)

func DefaultSystemConfiguration() *SystemConfiguration {
	return &SystemConfiguration{
		VersioningClient: &vcs.Client{
			Vcs: &vcs.Git{},
		},
		FileSystem: fs.DefaultFileSystem(),
		Clock:      helper.SystemClock{},
	}
}

type SystemConfiguration struct {
	VersioningClient vcs.VersioningClient
	FileSystem       *fs.FileSystem
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
	currentConfig *Configuration,
	system *SystemConfiguration,
	tracker ExecutionTracker,
	pathMatcher fs.PathMatcher) (*ChangeSet, error) {

	versionedTemplate, err := tracker.RetrieveVersionedTemplate(currentConfig)
	if err != nil {
		return nil, err
	}

	contents, err := ParseTemplate(versionedTemplate, ParseCommentStyle(currentConfig.CommentStyle))
	if err != nil {
		return nil, err
	}

	changes, err := getAffectedFiles(currentConfig, system, versionedTemplate, pathMatcher)
	if err != nil {
		return nil, err
	}

	return &ChangeSet{
		HeaderContents: contents.actualContent,
		HeaderRegex:    contents.detectionRegex,
		Files:          changes,
	}, nil
}

func getAffectedFiles(config *Configuration,
	sysConfig *SystemConfiguration,
	versionedTemplate *VersionedHeaderTemplate,
	pathMatcher fs.PathMatcher) ([]vcs.FileChange, error) {

	versioningClient := sysConfig.VersioningClient
	fileSystem := sysConfig.FileSystem
	var (
		changes []vcs.FileChange
		err     error
	)

	if versionedTemplate.RequiresFullScan() {
		log.Print("Unable to get last execution revision, triggering a full scan")
		changes, err = pathMatcher.ScanAllFiles(config.Includes, config.Excludes, fileSystem)
		if err != nil {
			return nil, err
		}
	} else {
		revision := versionedTemplate.Revision
		log.Printf("Scanning changes since revision %s", revision)
		fileChanges, err := versioningClient.GetChanges(revision)
		if err != nil {
			return nil, err
		}
		changes = pathMatcher.MatchFiles(fileChanges, config.Includes, config.Excludes, fileSystem)
	}
	return versioningClient.AddMetadata(changes, sysConfig.Clock)
}
