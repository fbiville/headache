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

package core_test

import (
	"github.com/fbiville/headache/core"
	"github.com/fbiville/headache/fs"
	"github.com/fbiville/headache/fs_mocks"
	"github.com/fbiville/headache/helper_mocks"
	. "github.com/fbiville/headache/vcs"
	"github.com/fbiville/headache/vcs_mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"testing"
)

func TestConfigurationInitWithSlashSlashStyle(t *testing.T) {
	I := NewGomegaWithT(t)
	controller := gomock.NewController(t)
	defer controller.Finish()

	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}
	versioningClient := new(vcs_mocks.VersioningClient)
	defer versioningClient.AssertExpectations(t)
	tracker := new(fs_mocks.ExecutionTracker)
	defer tracker.AssertExpectations(t)
	pathMatcher := new(fs_mocks.PathMatcher)
	defer pathMatcher.AssertExpectations(t)
	clock := new(helper_mocks.Clock)
	defer clock.AssertExpectations(t)

	initialChanges := []FileChange{{Path: "hello-world.go"}, {Path: "license.txt"}}
	includes := []string{"../fixtures/hello_*.go"}
	excludes := []string{}
	resultingChanges := []FileChange{initialChanges[0]}
	fileReader.On("Read", "some-header").
		Return([]byte("Copyright {{.Year}} {{.Owner}}\n\nSome fictional license"), nil)
	tracker.On("GetLastExecutionRevision").Return("some-sha", nil)
	versioningClient.On("GetChanges", "some-sha").Return(initialChanges, nil)
	pathMatcher.On("MatchFiles", initialChanges, includes, excludes, fileSystem).Return(resultingChanges)
	versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

	systemConfiguration := core.SystemConfiguration{
		FileSystem:       fileSystem,
		Clock:            clock,
		VersioningClient: versioningClient,
	}
	configuration := core.Configuration{
		HeaderFile:   "some-header",
		CommentStyle: "SlashSlash",
		Includes:     includes,
		Excludes:     excludes,
		TemplateData: map[string]string{
			"Owner": "ACME Labs",
		}}

	changeSet, err := core.ParseConfiguration(configuration, systemConfiguration, tracker, pathMatcher)

	I.Expect(err).To(BeNil())
	I.Expect(changeSet.HeaderContents).To(Equal("// Copyright {{.Year}} ACME Labs\n//\n// Some fictional license"))
	I.Expect(onlyPaths(changeSet.Files)).To(Equal([]FileChange{{Path: "hello-world.go"}}))
}

func TestConfigurationInitWithSlashStarStyle(t *testing.T) {
	I := NewGomegaWithT(t)
	controller := gomock.NewController(t)
	defer controller.Finish()

	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}
	versioningClient := new(vcs_mocks.VersioningClient)
	defer versioningClient.AssertExpectations(t)
	tracker := new(fs_mocks.ExecutionTracker)
	defer tracker.AssertExpectations(t)
	pathMatcher := new(fs_mocks.PathMatcher)
	defer pathMatcher.AssertExpectations(t)
	clock := new(helper_mocks.Clock)
	defer clock.AssertExpectations(t)

	initialChanges := []FileChange{{Path: "hello-world.go"}, {Path: "license.txt"}}
	includes := []string{"../fixtures/hello_*.go"}
	excludes := []string{}
	resultingChanges := []FileChange{initialChanges[0]}
	fileReader.On("Read", "some-header").
		Return([]byte("Copyright {{.Year}} {{.Owner}}\n\nSome fictional license"), nil)
	tracker.On("GetLastExecutionRevision").Return("some-sha", nil)
	versioningClient.On("GetChanges", "some-sha").Return(initialChanges, nil)
	pathMatcher.On("MatchFiles", initialChanges, includes, excludes, fileSystem).Return(resultingChanges)
	versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

	systemConfiguration := core.SystemConfiguration{
		FileSystem:       fileSystem,
		Clock:            clock,
		VersioningClient: versioningClient,
	}
	configuration := core.Configuration{
		HeaderFile:   "some-header",
		CommentStyle: "SlashStar",
		Includes:     includes,
		Excludes:     excludes,
		TemplateData: map[string]string{
			"Owner": "ACME Labs",
		}}

	changeSet, err := core.ParseConfiguration(configuration, systemConfiguration, tracker, pathMatcher)

	I.Expect(err).To(BeNil())
	I.Expect(changeSet.HeaderContents).To(Equal(`/*
 * Copyright {{.Year}} ACME Labs
 *
 * Some fictional license
 */`))
	I.Expect(onlyPaths(changeSet.Files)).To(Equal([]FileChange{{Path: "hello-world.go"}}))
}

func TestHeaderDetectionRegexComputation(t *testing.T) {
	I := NewGomegaWithT(t)
	controller := gomock.NewController(t)
	defer controller.Finish()

	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}
	versioningClient := new(vcs_mocks.VersioningClient)
	defer versioningClient.AssertExpectations(t)
	tracker := new(fs_mocks.ExecutionTracker)
	defer tracker.AssertExpectations(t)
	pathMatcher := new(fs_mocks.PathMatcher)
	defer pathMatcher.AssertExpectations(t)
	clock := new(helper_mocks.Clock)
	defer clock.AssertExpectations(t)

	initialChanges := []FileChange{{Path: "hello-world.go"}, {Path: "license.txt"}}
	includes := []string{"../fixtures/hello_*.go"}
	excludes := []string{}
	resultingChanges := []FileChange{initialChanges[0]}
	fileReader.On("Read", "some-header").
		Return([]byte("Copyright {{.Year}} {{.Owner}}"), nil)
	tracker.On("GetLastExecutionRevision").Return("some-sha", nil)
	versioningClient.On("GetChanges", "some-sha").Return(initialChanges, nil)
	pathMatcher.On("MatchFiles", initialChanges, includes, excludes, fileSystem).Return(resultingChanges)
	versioningClient.On("AddMetadata", resultingChanges, clock).Return(resultingChanges, nil)

	systemConfiguration := core.SystemConfiguration{
		FileSystem:       fileSystem,
		Clock:            clock,
		VersioningClient: versioningClient,
	}
	configuration := core.Configuration{
		HeaderFile:   "some-header",
		CommentStyle: "SlashStar",
		Includes:     includes,
		Excludes:     excludes,
		TemplateData: map[string]string{
			"Owner": "ACME Labs",
		}}

	changeSet, err := core.ParseConfiguration(configuration, systemConfiguration, tracker, pathMatcher)

	regex := changeSet.HeaderRegex
	I.Expect(err).To(BeNil())
	I.Expect(changeSet.HeaderContents).To(Equal(`/*
 * Copyright {{.Year}} ACME Labs
 */`))
	I.Expect(regex.MatchString(changeSet.HeaderContents)).To(BeTrue(), "Regex should match contents")
	I.Expect(regex.MatchString("// Copyright 2018 ACME Labs")).To(BeTrue(),
		"Regex should match contents with different comment style")
	I.Expect(regex.MatchString(`/*
 * Copyright 2018-2042 ACME World corporation
 */`)).To(BeTrue(), "Regex should match contents with different data")
	I.Expect(regex.MatchString("// Copyright 2009-2012 ACME!")).To(BeTrue(),
		"Regex should match contents with different data and comment style")
}

func TestFailOnReservedYearParameter(t *testing.T) {
	I := NewGomegaWithT(t)
	controller := gomock.NewController(t)
	defer controller.Finish()

	fileReader := new(fs_mocks.FileReader)
	defer fileReader.AssertExpectations(t)
	fileWriter := new(fs_mocks.FileWriter)
	defer fileWriter.AssertExpectations(t)
	fileSystem := fs.FileSystem{FileWriter: fileWriter, FileReader: fileReader}
	versioningClient := new(vcs_mocks.VersioningClient)
	defer versioningClient.AssertExpectations(t)
	tracker := new(fs_mocks.ExecutionTracker)
	defer tracker.AssertExpectations(t)
	pathMatcher := new(fs_mocks.PathMatcher)
	defer pathMatcher.AssertExpectations(t)
	clock := new(helper_mocks.Clock)
	defer clock.AssertExpectations(t)

	includes := []string{"../fixtures/hello_*.go"}
	excludes := []string{}
	fileReader.On("Read", "some-header").
		Return([]byte("Copyright {{.Year}}"), nil)

	systemConfiguration := core.SystemConfiguration{
		FileSystem:       fileSystem,
		Clock:            clock,
		VersioningClient: versioningClient,
	}
	configuration := core.Configuration{
		HeaderFile:   "some-header",
		CommentStyle: "SlashStar",
		Includes:     includes,
		Excludes:     excludes,
		TemplateData: map[string]string{
			"Year": "oopsie - reserved parameter!",
		}}

	changeSet, err := core.ParseConfiguration(configuration, systemConfiguration, tracker, pathMatcher)

	I.Expect(changeSet).To(BeNil())
	I.Expect(err).To(MatchError("Year is a reserved parameter and is automatically computed.\n" +
		"Please remove it from your configuration"))
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
