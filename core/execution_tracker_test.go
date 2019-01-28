/*
 * Copyright 2018-2019 Florent Biville (@fbiville)
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
	"errors"
	"fmt"
	"github.com/fbiville/headache/core"
	. "github.com/fbiville/headache/fs"
	"github.com/fbiville/headache/fs_mocks"
	"github.com/fbiville/headache/vcs_mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"reflect"
	"strings"
	"time"
)

type FixedClock struct{}

func (*FixedClock) Now() time.Time {
	return time.Unix(42, 42)
}

var _ = Describe("The execution tracker", func() {

	var (
		t          GinkgoTInterface
		vcs        *vcs_mocks.Vcs
		fileReader *fs_mocks.FileReader
		fileWriter *fs_mocks.FileWriter
		tracker    core.ExecutionVcsTracker
	)

	BeforeEach(func() {
		t = GinkgoT()
		vcs = new(vcs_mocks.Vcs)
		fileReader = new(fs_mocks.FileReader)
		fileWriter = new(fs_mocks.FileWriter)
		tracker = core.ExecutionVcsTracker{
			Versioning: vcs,
			FileSystem: &FileSystem{FileReader: fileReader, FileWriter: fileWriter},
			Clock:      &FixedClock{},
		}
	})

	AfterEach(func() {
		vcs.AssertExpectations(t)
		fileReader.AssertExpectations(t)
		fileWriter.AssertExpectations(t)
	})

	Describe("when retrieving the versioned header template", func() {

		var (
			currentHeaderFile    string
			currentData          map[string]string
			currentConfiguration *core.Configuration
			fakeRepositoryRoot   string
			trackerFilePath      string
		)

		BeforeEach(func() {
			currentHeaderFile = "current-header-file"
			currentData = map[string]string{"foo": "bar"}
			currentConfiguration = &core.Configuration{
				HeaderFile:   currentHeaderFile,
				TemplateData: currentData,
			}
			fakeRepositoryRoot = "/path/to"
			trackerFilePath = fakeRepositoryRoot + "/.headache-run"
		})

		It("returns only the current contents if there were no previous execution", func() {
			currentContents := "some\nheader"
			fileReader.On("Read", currentHeaderFile).Return([]byte(currentContents), nil)
			vcs.On("Root").Return(fakeRepositoryRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			vcs.On("LatestRevision", trackerFilePath).Return("", nil)

			versionedTemplate, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(BeNil())
			Expect(versionedTemplate.Revision).To(BeEmpty())
			Expect(versionedTemplate.Current.Data).To(Equal(currentData))
			Expect(strings.Join(versionedTemplate.Current.Lines, "\n")).To(Equal(currentContents))
			Expect(versionedTemplate.Previous.Data).To(Equal(currentData))
			Expect(strings.Join(versionedTemplate.Previous.Lines, "\n")).To(Equal(currentContents))
		})

		It("returns only the current contents if there were no tracked configuration for backwards compatibility", func() {
			currentContents := "some\nheader"
			fileReader.On("Read", currentHeaderFile).Return([]byte(currentContents), nil)
			vcs.On("Root").Return(fakeRepositoryRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			vcs.On("LatestRevision", trackerFilePath).Return("some-sha", nil)
			fileReader.On("Read", trackerFilePath).Return([]byte("no tracked configuration in here"), nil)

			versionedTemplate, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(BeNil())
			Expect(versionedTemplate.Revision).To(BeEmpty())
			Expect(versionedTemplate.Current.Data).To(Equal(currentData))
			Expect(strings.Join(versionedTemplate.Current.Lines, "\n")).To(Equal(currentContents))
			Expect(versionedTemplate.Previous.Data).To(Equal(currentData))
			Expect(strings.Join(versionedTemplate.Previous.Lines, "\n")).To(Equal(currentContents))
		})

		It("returns the current and previous contents", func() {
			previousConfigFile := "previous-config"
			revision := "some-revision"
			previousHeaderFile := "previous-header"
			currentContents := "some\nheader"
			previousContents := "previous\nheader"
			fileReader.On("Read", currentHeaderFile).Return([]byte(currentContents), nil)
			vcs.On("Root").Return(fakeRepositoryRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			vcs.On("LatestRevision", trackerFilePath).Return(revision, nil)
			fileReader.On("Read", trackerFilePath).Return([]byte("configuration:"+previousConfigFile), nil)
			vcs.On("ShowContentAtRevision", previousConfigFile, revision).Return(fmt.Sprintf(`{
  "headerFile": "%s",
  "data": {"some": "thing"}
}`, previousHeaderFile), nil)
			vcs.On("ShowContentAtRevision", previousHeaderFile, revision).Return(previousContents, nil)

			result, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(BeNil())
			Expect(strings.Join(result.Current.Lines, "\n")).To(Equal(currentContents))
			Expect(result.Current.Data).To(Equal(currentData))
			Expect(result.Revision).To(Equal(revision))
			Expect(strings.Join(result.Previous.Lines, "\n")).To(Equal(previousContents))
			Expect(result.Previous.Data).To(Equal(map[string]string{"some": "thing"}))
		})

		It("fails of the current template cannot be read", func() {
			expectedError := errors.New("read error")
			fileReader.On("Read", currentHeaderFile).Return(nil, expectedError)

			_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(MatchError(expectedError))
		})

		It("fails if the current repository root cannot be determined", func() {
			expectedError := errors.New("root error")
			fileReader.On("Read", currentHeaderFile).Return([]byte("some\nheader"), nil)
			vcs.On("Root").Return("", expectedError)

			_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(MatchError(expectedError))
		})

		It("fails if it cannot get stats on the tracker file", func() {
			expectedError := errors.New("stat error")
			fileReader.On("Read", currentHeaderFile).Return([]byte("some\nheader"), nil)
			vcs.On("Root").Return(fakeRepositoryRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(nil, expectedError)

			_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(MatchError(expectedError))
		})

		It("fails if the tracker file is not a regular file", func() {
			fileReader.On("Read", currentHeaderFile).Return([]byte("some\nheader"), nil)
			vcs.On("Root").Return(fakeRepositoryRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: os.ModeDir}, nil)

			_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(MatchError(fmt.Sprintf("'%s' should be a regular file", trackerFilePath)))
		})

		It("fails if the latest execution's revision cannot be retrieved", func() {
			expectedError := errors.New("revision error")
			fileReader.On("Read", currentHeaderFile).Return([]byte("some\nheader"), nil)
			vcs.On("Root").Return(fakeRepositoryRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			vcs.On("LatestRevision", trackerFilePath).Return("", expectedError)

			_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(MatchError(expectedError))
		})

		It("fails if the tracking file cannot be read", func() {
			expectedError := errors.New("read error")
			fileReader.On("Read", currentHeaderFile).Return([]byte("some\nheader"), nil)
			vcs.On("Root").Return(fakeRepositoryRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			vcs.On("LatestRevision", trackerFilePath).Return("some-revision", nil)
			fileReader.On("Read", trackerFilePath).Return(nil, expectedError)

			_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(MatchError(expectedError))
		})

		It("fails if the tracking file contents cannot be shown at last execution's revision", func() {
			expectedError := errors.New("show at revision error")
			previousConfigFile := "previous-config"
			revision := "some-revision"
			fileReader.On("Read", currentHeaderFile).Return([]byte("some\nheader"), nil)
			vcs.On("Root").Return(fakeRepositoryRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			vcs.On("LatestRevision", trackerFilePath).Return(revision, nil)
			fileReader.On("Read", trackerFilePath).Return([]byte("configuration:"+previousConfigFile), nil)
			vcs.On("ShowContentAtRevision", previousConfigFile, revision).Return("", expectedError)

			_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(MatchError(expectedError))
		})

		It("fails if the previous configuration cannot be unmarshalled", func() {
			previousConfigFile := "previous-config"
			revision := "some-revision"
			fileReader.On("Read", currentHeaderFile).Return([]byte("some\nheader"), nil)
			vcs.On("Root").Return(fakeRepositoryRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			vcs.On("LatestRevision", trackerFilePath).Return(revision, nil)
			fileReader.On("Read", trackerFilePath).Return([]byte("configuration:"+previousConfigFile), nil)
			vcs.On("ShowContentAtRevision", previousConfigFile, revision).Return("not-json", nil)

			_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(reflect.TypeOf(err).String()).To(Equal("*json.SyntaxError"))
		})

		It("fails if the previous header file cannot be read", func() {
			expectedError := errors.New("show at revision error")
			previousConfigFile := "previous-config"
			revision := "some-revision"
			previousHeaderFile := "previous-header"
			fileReader.On("Read", currentHeaderFile).Return([]byte("some\nheader"), nil)
			vcs.On("Root").Return(fakeRepositoryRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			vcs.On("LatestRevision", trackerFilePath).Return(revision, nil)
			fileReader.On("Read", trackerFilePath).Return([]byte("configuration:"+previousConfigFile), nil)
			vcs.On("ShowContentAtRevision", previousConfigFile, revision).Return(fmt.Sprintf(`{
  "headerFile": "%s",
  "data": {}
}`, previousHeaderFile), nil)
			vcs.On("ShowContentAtRevision", previousHeaderFile, revision).Return("", expectedError)

			_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(MatchError(expectedError))
		})
	})

	Describe("when tracking headache execution", func() {

		var (
			repositoryFakeRoot   string
			trackerFilePath      string
			configurationPath    *string
			fileMode             os.FileMode
			expectedFileContents string
		)

		BeforeEach(func() {
			repositoryFakeRoot = "/path/to"
			trackerFilePath = repositoryFakeRoot + "/.headache-run"
			configPath := "/path/to/headache.json"
			configurationPath = &configPath
			fileMode = 0640
			expectedFileContents = fmt.Sprintf(`# Generated by headache | 42 -- commit me!
configuration:%s`, configPath)
		})

		It("saves the timestamp and the path to the runtime configuration file without prior tracking file", func() {
			vcs.On("Root").Return(repositoryFakeRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(nil, os.ErrNotExist)
			fileWriter.On("Write", trackerFilePath, expectedFileContents, fileMode).Return(nil)

			err := tracker.TrackExecution(configurationPath)

			Expect(err).To(BeNil())
		})

		It("saves the timestamp and the path to the runtime configuration file with prior tracking file", func() {
			vcs.On("Root").Return(repositoryFakeRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			fileWriter.On("Write", trackerFilePath, expectedFileContents, fileMode).Return(nil)

			err := tracker.TrackExecution(configurationPath)

			Expect(err).To(BeNil())
		})

		It("fails if the repository root cannot be retrieved", func() {
			rootErr := errors.New("root error")
			vcs.On("Root").Return("", rootErr)

			err := tracker.TrackExecution(configurationPath)

			Expect(err).To(MatchError(rootErr))
		})

		It("fails if the tracker path is not a regular file", func() {
			vcs.On("Root").Return(repositoryFakeRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: os.ModeDir}, nil)

			err := tracker.TrackExecution(configurationPath)

			Expect(err).To(MatchError(fmt.Sprintf("'%s' should be a regular file", trackerFilePath)))
		})

		It("fails if the stat call failed", func() {
			statError := errors.New("stat fail")
			vcs.On("Root").Return(repositoryFakeRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(nil, statError)

			err := tracker.TrackExecution(configurationPath)

			Expect(err).To(MatchError(statError))
		})

		It("fails if the write call failed", func() {
			writeError := errors.New("write fail")
			vcs.On("Root").Return(repositoryFakeRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			fileWriter.On("Write", trackerFilePath, expectedFileContents, fileMode).Return(writeError)

			err := tracker.TrackExecution(configurationPath)

			Expect(err).To(MatchError(writeError))
		})
	})
})
