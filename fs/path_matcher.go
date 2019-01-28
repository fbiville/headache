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

package fs

import (
	"github.com/fbiville/headache/vcs"
	"github.com/mattn/go-zglob"
)

type PathMatcher interface {
	ScanAllFiles(includes []string, excludes []string, filesystem *FileSystem) ([]vcs.FileChange, error)
	MatchFiles(changes []vcs.FileChange, includes []string, excludes []string, filesystem *FileSystem) []vcs.FileChange
}

type ZglobPathMatcher struct {}

// Scans all local files based on the provided inclusion and exclusion patterns
func (*ZglobPathMatcher) ScanAllFiles(includes []string, excludes []string, filesystem *FileSystem) ([]vcs.FileChange, error) {
	result := make([]vcs.FileChange, 0)
	for _, includePattern := range includes {
		matches, err := zglob.Glob(includePattern)
		if err != nil {
			return nil, err
		}
		for _, matchedPath := range matches {
			if !isExcluded(matchedPath, excludes, filesystem) {
				result = append(result, vcs.FileChange{
					Path: matchedPath,
				})
			}
		}

	}
	return result, nil
}

// Matches files based on the provided inclusion and exclusion patterns
func (*ZglobPathMatcher) MatchFiles(changes []vcs.FileChange, includes []string, excludes []string, filesystem *FileSystem) []vcs.FileChange {
	result := make([]vcs.FileChange, 0)
	for _, change := range changes {
		if matches(change.Path, includes, excludes, filesystem) {
			result = append(result, change)
		}
	}
	return result
}

func matches(path string, includes []string, excludes []string, filesystem *FileSystem) bool {
	return matchesPattern(path, includes) && !isExcluded(path, excludes, filesystem)
}

func isExcluded(path string, excludes []string, filesystem *FileSystem) bool {
	return !filesystem.IsFile(path) || matchesPattern(path, excludes)
}

func matchesPattern(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, err := zglob.Match(pattern, path); err == nil && matched {
			return true
		}
	}
	return false
}
