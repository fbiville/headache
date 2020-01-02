/*
 * Copyright 2018-2019 Florent Biville (@fbiville)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package core

import (
	"github.com/fbiville/headache/internal/pkg/fs"
	"github.com/fbiville/headache/internal/pkg/vcs"
	"log"
	"regexp"
)

type Configuration struct {
	HeaderFile   string            `json:"headerFile"`
	CommentStyle string            `json:"style"`
	Includes     []string          `json:"includes"`
	Excludes     []string          `json:"excludes"`
	TemplateData map[string]string `json:"data"`
	Path         *string
}

type ChangeSet struct {
	HeaderContents string
	HeaderRegex    *regexp.Regexp
	Files          []vcs.FileChange
}

type ConfigurationResolver struct {
	Environment      *Environment
	ExecutionTracker ExecutionTracker
	PathMatcher      fs.PathMatcher
}

func (resolver *ConfigurationResolver) ResolveEagerly(currentConfig *Configuration) (*ChangeSet, error) {
	versionedTemplate, err := resolver.ExecutionTracker.RetrieveVersionedTemplate(currentConfig)
	if err != nil {
		return nil, err
	}

	contents, err := ParseTemplate(versionedTemplate, ParseCommentStyle(currentConfig.CommentStyle))
	if err != nil {
		return nil, err
	}

	changes, err := resolver.getAffectedFiles(currentConfig, versionedTemplate)
	if err != nil {
		return nil, err
	}

	return &ChangeSet{
		HeaderContents: contents.ActualContent,
		HeaderRegex:    contents.DetectionRegex,
		Files:          changes,
	}, nil
}

func (resolver *ConfigurationResolver) getAffectedFiles(config *Configuration, versionedTemplate *VersionedHeaderTemplate) ([]vcs.FileChange, error) {
	versioningClient := resolver.Environment.VersioningClient
	fileSystem := resolver.Environment.FileSystem
	var (
		changes []vcs.FileChange
		err     error
	)

	if versionedTemplate.RequiresFullScan() {
		if versionedTemplate.Revision == "" {
			log.Print("Unable to get last execution revision, triggering a full scan")
		} else {
			log.Printf("Configuration and/or license header template changed since last execution (%s), triggering a full scan", versionedTemplate.Revision)
		}
		changes, err = resolver.PathMatcher.ScanAllFiles(config.Includes, config.Excludes, fileSystem)
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
		changes = resolver.PathMatcher.MatchFiles(fileChanges, config.Includes, config.Excludes, fileSystem)
	}
	return versioningClient.AddMetadata(changes, resolver.Environment.Clock)
}
