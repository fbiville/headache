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

package fs_test

import (
	. "github.com/fbiville/headache/internal/pkg/fs"
	"github.com/fbiville/headache/internal/pkg/fs_mocks"
	"github.com/fbiville/headache/internal/pkg/vcs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"os"
)

var _ = Describe("Path matcher", func() {
	var (
		t          GinkgoTInterface
		fileReader *fs_mocks.FileReader
		fileSystem *FileSystem
		matcher    *ZglobPathMatcher
	)

	BeforeEach(func() {
		t = GinkgoT()
		fileReader = new(fs_mocks.FileReader)
		fileSystem = &FileSystem{FileReader: fileReader}
		matcher = &ZglobPathMatcher{}
	})

	AfterEach(func() {
		fileReader.AssertExpectations(t)
	})

	It("matches paths matching include patterns and not matching exclude patterns", func() {
		includes := []string{"../fixtures/*.txt"}
		excludes := []string{"../fixtures/*_with_header.txt"}
		fileReader.On("Stat", mock.Anything).Return(&FakeFileInfo{FileMode: 0777}, nil)

		Expect(matcher.MatchFiles([]vcs.FileChange{{Path: "../fixtures/bonjour_world.txt"}}, includes, excludes, fileSystem)).
			To(Equal([]vcs.FileChange{{Path: "../fixtures/bonjour_world.txt"}}))
		Expect(matcher.MatchFiles([]vcs.FileChange{{Path: "../fixtures/bonjour_world.go"}}, includes, excludes, fileSystem)).
			To(BeEmpty())
		Expect(matcher.MatchFiles([]vcs.FileChange{{Path: "../fixtures/hello_world_with_header.txt"}}, includes, excludes, fileSystem)).
			To(BeEmpty())
	})

	It("matches only files", func() {
		includes := []string{"../fixture*"}
		excludes := []string{}
		fileReader.On("Stat", mock.Anything).Return(&FakeFileInfo{FileMode: os.ModeDir}, nil)

		Expect(matcher.MatchFiles([]vcs.FileChange{{Path: "../fixtures"}}, includes, excludes, fileSystem)).To(BeEmpty())
	})
})
