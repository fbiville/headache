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

package core_test

import (
	"strings"

	"github.com/fbiville/headache/internal/pkg/core"
	"github.com/fbiville/headache/internal/pkg/core_mocks"
	"github.com/fbiville/headache/internal/pkg/fs"
	"github.com/fbiville/headache/internal/pkg/fs_mocks"
	"github.com/fbiville/headache/internal/pkg/helper_mocks"
	. "github.com/fbiville/headache/internal/pkg/vcs"
	"github.com/fbiville/headache/internal/pkg/vcs_mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration parser", func() {
	var (
		t                     GinkgoTInterface
		fileReader            *fs_mocks.FileReader
		fileWriter            *fs_mocks.FileWriter
		fileSystem            *fs.FileSystem
		versioningClient      *vcs_mocks.VersioningClient
		tracker               *core_mocks.ExecutionTracker
		pathMatcher           *fs_mocks.PathMatcher
		clock                 *helper_mocks.Clock
		initialChanges        []FileChange
		includes              []string
		excludes              []string
		resultingChanges      []FileChange
		systemConfiguration   *core.SystemConfiguration
		data                  map[string]string
		revision              string
		configurationResolver *core.ConfigurationResolver
	)

	BeforeEach(func() {
		t = GinkgoT()
		fileReader = new(fs_mocks.FileReader)
		fileWriter = new(fs_mocks.FileWriter)
		fileSystem = &fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}
		versioningClient = new(vcs_mocks.VersioningClient)
		tracker = new(core_mocks.ExecutionTracker)
		pathMatcher = new(fs_mocks.PathMatcher)
		clock = new(helper_mocks.Clock)
		initialChanges = []FileChange{{Path: "hello-world.go"}, {Path: "license.txt"}}
		includes = []string{"../fixtures/hello_*.go"}
		excludes = []string{}
		resultingChanges = []FileChange{initialChanges[0]}
		systemConfiguration = &core.SystemConfiguration{
			FileSystem:       fileSystem,
			Clock:            clock,
			VersioningClient: versioningClient,
		}
		data = map[string]string{
			"Owner": "ACME Labs",
		}
		revision = "some-sha"
		configurationResolver = &core.ConfigurationResolver{
			SystemConfiguration: systemConfiguration,
			ExecutionTracker:    tracker,
			PathMatcher:         pathMatcher,
		}
	})

	AfterEach(func() {
		fileReader.AssertExpectations(t)
		fileWriter.AssertExpectations(t)
		versioningClient.AssertExpectations(t)
		tracker.AssertExpectations(t)
		pathMatcher.AssertExpectations(t)
		clock.AssertExpectations(t)
	})

	It("pre-computes the final configuration", func() {
		configuration := &core.Configuration{
			HeaderFile:   "some-header",
			CommentStyle: "SlashSlash",
			Includes:     includes,
			Excludes:     excludes,
			TemplateData: data,
		}
		tracker.On("RetrieveVersionedTemplate", configuration).
			Return(unchangedHeaderContents("Copyright {{.Year}} {{.Owner}}\n\nSome fictional license", data, revision), nil)
		versioningClient.On("GetChanges", revision).Return(initialChanges, nil)
		pathMatcher.On("MatchFiles", initialChanges, includes, excludes, fileSystem).Return(resultingChanges)
		versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

		changeSet, err := configurationResolver.ResolveEagerly(configuration)

		Expect(err).To(BeNil())
		Expect(changeSet.HeaderContents).To(Equal("// Copyright {{.YearRange}} ACME Labs\n//\n// Some fictional license"))
		Expect(onlyPaths(changeSet.Files)).To(Equal([]FileChange{{Path: "hello-world.go"}}))
	})

	It("pre-computes the final configuration with dashdash comment style", func() {
		configuration := &core.Configuration{
			HeaderFile:   "some-header",
			CommentStyle: "DashDash",
			Includes:     includes,
			Excludes:     excludes,
			TemplateData: data,
		}
		tracker.On("RetrieveVersionedTemplate", configuration).
			Return(unchangedHeaderContents("Copyright {{.Year}} {{.Owner}}\n\nSome fictional license", data, revision), nil)
		versioningClient.On("GetChanges", revision).Return(initialChanges, nil)
		pathMatcher.On("MatchFiles", initialChanges, includes, excludes, fileSystem).Return(resultingChanges)
		versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

		changeSet, err := configurationResolver.ResolveEagerly(configuration)

		Expect(err).To(BeNil())
		Expect(changeSet.HeaderContents).To(Equal("-- Copyright {{.YearRange}} ACME Labs\n--\n-- Some fictional license"))
		Expect(onlyPaths(changeSet.Files)).To(Equal([]FileChange{{Path: "hello-world.go"}}))
	})

	It("pre-computes the final configuration with semicolon comment style", func() {
		configuration := &core.Configuration{
			HeaderFile:   "some-header",
			CommentStyle: "SemiColon",
			Includes:     includes,
			Excludes:     excludes,
			TemplateData: data,
		}
		tracker.On("RetrieveVersionedTemplate", configuration).
			Return(unchangedHeaderContents("Copyright {{.Year}} {{.Owner}}\n\nSome fictional license", data, revision), nil)
		versioningClient.On("GetChanges", revision).Return(initialChanges, nil)
		pathMatcher.On("MatchFiles", initialChanges, includes, excludes, fileSystem).Return(resultingChanges)
		versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

		changeSet, err := configurationResolver.ResolveEagerly(configuration)

		Expect(err).To(BeNil())
		Expect(changeSet.HeaderContents).To(Equal("; Copyright {{.YearRange}} ACME Labs\n;\n; Some fictional license"))
		Expect(onlyPaths(changeSet.Files)).To(Equal([]FileChange{{Path: "hello-world.go"}}))
	})

	It("pre-computes the final configuration with hash comment style", func() {
		configuration := &core.Configuration{
			HeaderFile:   "some-header",
			CommentStyle: "Hash",
			Includes:     includes,
			Excludes:     excludes,
			TemplateData: data,
		}
		tracker.On("RetrieveVersionedTemplate", configuration).
			Return(unchangedHeaderContents("Copyright {{.Year}} {{.Owner}}\n\nSome fictional license", data, revision), nil)
		versioningClient.On("GetChanges", revision).Return(initialChanges, nil)
		pathMatcher.On("MatchFiles", initialChanges, includes, excludes, fileSystem).Return(resultingChanges)
		versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

		changeSet, err := configurationResolver.ResolveEagerly(configuration)

		Expect(err).To(BeNil())
		Expect(changeSet.HeaderContents).To(Equal("# Copyright {{.YearRange}} ACME Labs\n#\n# Some fictional license"))
		Expect(onlyPaths(changeSet.Files)).To(Equal([]FileChange{{Path: "hello-world.go"}}))
	})

	It("pre-computes the final configuration with REM comment style", func() {
		configuration := &core.Configuration{
			HeaderFile:   "some-header",
			CommentStyle: "REM",
			Includes:     includes,
			Excludes:     excludes,
			TemplateData: data,
		}
		tracker.On("RetrieveVersionedTemplate", configuration).
			Return(unchangedHeaderContents("Copyright {{.Year}} {{.Owner}}\n\nSome fictional license", data, revision), nil)
		versioningClient.On("GetChanges", revision).Return(initialChanges, nil)
		pathMatcher.On("MatchFiles", initialChanges, includes, excludes, fileSystem).Return(resultingChanges)
		versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

		changeSet, err := configurationResolver.ResolveEagerly(configuration)

		Expect(err).To(BeNil())
		Expect(changeSet.HeaderContents).To(Equal("REM Copyright {{.YearRange}} ACME Labs\nREM\nREM Some fictional license"))
		Expect(onlyPaths(changeSet.Files)).To(Equal([]FileChange{{Path: "hello-world.go"}}))
	})

	It("pre-computes the final configuration with SlashStarStar comment style", func() {
		configuration := &core.Configuration{
			HeaderFile:   "some-header",
			CommentStyle: "SlashStarStar",
			Includes:     includes,
			Excludes:     excludes,
			TemplateData: data,
		}
		tracker.On("RetrieveVersionedTemplate", configuration).
			Return(unchangedHeaderContents("Copyright {{.Year}} {{.Owner}}\n\nSome fictional license", data, revision), nil)
		versioningClient.On("GetChanges", revision).Return(initialChanges, nil)
		pathMatcher.On("MatchFiles", initialChanges, includes, excludes, fileSystem).Return(resultingChanges)
		versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

		changeSet, err := configurationResolver.ResolveEagerly(configuration)

		Expect(err).To(BeNil())
		Expect(changeSet.HeaderContents).To(Equal("/**\n * Copyright {{.YearRange}} ACME Labs\n *\n * Some fictional license\n */"))
		Expect(onlyPaths(changeSet.Files)).To(Equal([]FileChange{{Path: "hello-world.go"}}))
	})

	It("pre-computes the header contents with SlashStar comment style", func() {
		configuration := &core.Configuration{
			HeaderFile:   "some-header",
			CommentStyle: "SlashStar",
			Includes:     includes,
			Excludes:     excludes,
			TemplateData: data,
		}
		tracker.On("RetrieveVersionedTemplate", configuration).
			Return(unchangedHeaderContents("Copyright {{.Year}} {{.Owner}}\n\nSome fictional license", data, revision), nil)
		versioningClient.On("GetChanges", revision).Return(initialChanges, nil)
		pathMatcher.On("MatchFiles", initialChanges, includes, excludes, fileSystem).Return(resultingChanges)
		versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

		changeSet, err := configurationResolver.ResolveEagerly(configuration)

		Expect(err).To(BeNil())
		Expect(changeSet.HeaderContents).To(Equal(`/*
 * Copyright {{.YearRange}} ACME Labs
 *
 * Some fictional license
 */`))
		Expect(onlyPaths(changeSet.Files)).To(Equal([]FileChange{{Path: "hello-world.go"}}))
	})

	It("pre-computes a regex that allows to detect headers", func() {
		configuration := &core.Configuration{
			HeaderFile:   "some-header",
			CommentStyle: "SlashStar",
			Includes:     includes,
			Excludes:     excludes,
			TemplateData: data,
		}
		tracker.On("RetrieveVersionedTemplate", configuration).
			Return(unchangedHeaderContents("Copyright {{.Year}} {{.Owner}}", data, revision), nil)
		versioningClient.On("GetChanges", revision).Return(initialChanges, nil)
		pathMatcher.On("MatchFiles", initialChanges, includes, excludes, fileSystem).Return(resultingChanges)
		versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

		changeSet, err := configurationResolver.ResolveEagerly(configuration)

		regex := changeSet.HeaderRegex
		Expect(err).To(BeNil())
		Expect(changeSet.HeaderContents).To(Equal(`/*
 * Copyright {{.YearRange}} ACME Labs
 */`))
		Expect(regex.MatchString(changeSet.HeaderContents)).To(BeTrue(), "Regex should match contents")
		Expect(regex.MatchString("// Copyright 2018 ACME Labs")).To(BeTrue(),
			"Regex should match contents with different comment style")
		Expect(regex.MatchString("# Copyright 2018 ACME Labs")).To(BeTrue(),
			"Regex should match contents with different comment style")
		Expect(regex.MatchString(`/*
 * Copyright 2018-2042 ACME World corporation
 */`)).To(BeTrue(), "Regex should match contents with different data")
		Expect(regex.MatchString("// Copyright 2009-2012 ACME!")).To(BeTrue(),
			"Regex should match contents with different data and comment style")
		Expect(regex.MatchString("# Copyright 2009-2012 ACME!")).To(BeTrue(),
			"Regex should match contents with different data and comment style")
	})

	It("computes the header regex based on previous configuration", func() {
		configuration := &core.Configuration{
			HeaderFile:   "some-header",
			CommentStyle: "SlashSlash",
			Includes:     includes,
			Excludes:     excludes,
			TemplateData: data,
		}
		tracker.On("RetrieveVersionedTemplate", configuration).
			Return(&core.VersionedHeaderTemplate{
				Current:  template("new\nheader {{.Owner}}", map[string]string{"Owner": "Someone"}),
				Revision: revision,
				Previous: template("{{.Notice}} - old\nheader", map[string]string{"Notice": "Redding"}),
			}, nil)
		pathMatcher.On("ScanAllFiles", includes, excludes, fileSystem).Return(resultingChanges, nil)
		versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

		changeSet, err := configurationResolver.ResolveEagerly(configuration)

		regex := changeSet.HeaderRegex
		Expect(err).To(BeNil())
		Expect(regex.MatchString("// Redding - old\n// header")).To(BeTrue(),
			"Regex should match headers generated from previous run")
	})
})

func unchangedHeaderContents(lines string, data map[string]string, revision string) *core.VersionedHeaderTemplate {
	unchangedTemplate := template(lines, data)
	return &core.VersionedHeaderTemplate{
		Current:  unchangedTemplate,
		Revision: revision,
		Previous: unchangedTemplate,
	}
}

func template(lines string, data map[string]string) *core.HeaderTemplate {
	unchangedTemplate := &core.HeaderTemplate{
		Lines: strings.Split(lines, "\n"),
		Data:  data,
	}
	return unchangedTemplate
}

func onlyPaths(changes []FileChange) []FileChange {
	result := make([]FileChange, len(changes))
	for i := range changes {
		result[i] = FileChange{
			Path: changes[i].Path,
		}
	}
	return result
}
