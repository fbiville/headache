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

package fs_test

import (
	. "github.com/fbiville/headache/fs"
	"github.com/fbiville/headache/vcs"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMatch(t *testing.T) {
	I := NewGomegaWithT(t)
	includes := []string{"../fixtures/*.txt"}
	excludes := []string{"../fixtures/*_with_header.txt"}
	fileSystem := FileSystem{FileReader: &OsFileReader{}}
	matcher := &ZglobPathMatcher{}

	matchedChanges := []vcs.FileChange{{Path: "../fixtures/bonjour_world.txt"}}
	I.Expect(matcher.MatchFiles(matchedChanges, includes, excludes, fileSystem)).To(Equal(matchedChanges))
	I.Expect(matcher.MatchFiles([]vcs.FileChange{{Path: "../fixtures/bonjour_world.go"}}, includes, excludes, fileSystem)).To(BeEmpty())
	I.Expect(matcher.MatchFiles([]vcs.FileChange{{Path: "../fixtures/hello_world_with_header.txt"}}, includes, excludes, fileSystem)).To(BeEmpty())
}

func TestMatchOnlyFiles(t *testing.T) {
	I := NewGomegaWithT(t)
	fileSystem := FileSystem{FileReader: &OsFileReader{}}
	matcher := &ZglobPathMatcher{}

	I.Expect(matcher.MatchFiles([]vcs.FileChange{{Path: "../fixtures"}}, []string{"../fixture*"}, []string{}, fileSystem)).To(BeEmpty())
}
