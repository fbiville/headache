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
	"encoding/base64"
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
			configurationPath := "/path/to/headache.json"
			currentHeaderFile = "current-header-file"
			currentData = map[string]string{"foo": "bar"}
			currentConfiguration = &core.Configuration{
				HeaderFile:   currentHeaderFile,
				TemplateData: currentData,
				Path:         &configurationPath,
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

		It("gets the current config at the previous revision if there were no tracked configuration, for backwards compatibility", func() {
			revision := "some-revision"
			currentContents := "some\nheader"
			fileReader.On("Read", currentHeaderFile).Return([]byte(currentContents), nil)
			vcs.On("Root").Return(fakeRepositoryRoot, nil)
			fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			vcs.On("LatestRevision", trackerFilePath).Return(revision, nil)
			fileReader.On("Read", trackerFilePath).Return([]byte("no tracked configuration in here"), nil)
			currentConfigPreviousContents := fmt.Sprintf(`{
  "headerFile": "%s",
  "data": {"some": "thing"}
}`, *currentConfiguration.Path)
			vcs.On("ShowContentAtRevision", *currentConfiguration.Path, revision).Return(currentConfigPreviousContents, nil)

			versionedTemplate, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

			Expect(err).To(BeNil())
			Expect(versionedTemplate.Revision).To(Equal(revision))
			Expect(versionedTemplate.Current.Data).To(Equal(currentData))
			Expect(strings.Join(versionedTemplate.Current.Lines, "\n")).To(Equal(currentContents))
			Expect(versionedTemplate.Previous.Data).To(Equal(map[string]string{"some": "thing"}))
			Expect(strings.Join(versionedTemplate.Previous.Lines, "\n")).To(Equal(currentConfigPreviousContents))
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
			repositoryFakeRoot string
			trackerFilePath    string

			configurationPath string
			configuration     string

			headerPath string
			header     string

			fileMode os.FileMode

			expectedContents string
		)

		BeforeEach(func() {
			repositoryFakeRoot = "/path/to"
			trackerFilePath = repositoryFakeRoot + "/.headache-run"
			configurationPath = "/path/to/headache.json"
			fileMode = 0640
			headerPath = "./license-header.txt"
			configuration = fmt.Sprintf(`{
  "headerFile": "%s",
  "style": "SlashStar",
  "includes": ["**/*.go"],
  "excludes": ["vendor/**/*", "*_mocks/**/*"],
  "data": {
    "Owner": "Florent Biville (@fbiville)"
  }
}
`, headerPath)
			header = `Copyright {{.YearRange}} {{.Owner}}

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`
			expectedContents = fmt.Sprintf(`# Generated by headache | 42 -- commit me!
encoded_configuration:%s
encoded_header:%s
`, base64Encode(configuration), base64Encode(header))
		})

		Context("without pre-existing tracker file", func() {
			BeforeEach(func() {
				vcs.On("Root").Return(repositoryFakeRoot, nil)
				fileReader.On("Stat", trackerFilePath).Return(nil, os.ErrNotExist)
			})

			It("serializes both the configuration and the configured header's contents", func() {
				fileReader.On("Read", configurationPath).Return([]byte(configuration), nil)
				fileReader.On("Read", headerPath).Return([]byte(header), nil)
				fileWriter.On("Write", trackerFilePath, expectedContents, fileMode).Return(nil)

				err := tracker.TrackExecution(&configurationPath)

				Expect(err).To(BeNil())
			})
		})

		Context("with pre-existing tracker file", func() {
			BeforeEach(func() {
				vcs.On("Root").Return(repositoryFakeRoot, nil)
				fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: os.ModePerm}, nil)
			})

			It("serializes both the configuration and the configured header's contents", func() {
				fileReader.On("Read", configurationPath).Return([]byte(configuration), nil)
				fileReader.On("Read", headerPath).Return([]byte(header), nil)
				fileWriter.On("Write", trackerFilePath, expectedContents, fileMode).Return(nil)

				err := tracker.TrackExecution(&configurationPath)

				Expect(err).To(BeNil())
			})
		})

		Context("with an error while proceeding to repository root", func() {
			var err error

			BeforeEach(func() {
				err = fmt.Errorf("vcs root error")

				vcs.On("Root").Return("", err)
			})

			It("forwards the error with context", func() {
				actual := tracker.TrackExecution(&configurationPath)

				Expect(actual).To(MatchError(err))
			})
		})

		Context("with an error while reading the possibly existing tracker file", func() {
			var err error

			BeforeEach(func() {
				err = fmt.Errorf("tracker file read error")

				vcs.On("Root").Return(repositoryFakeRoot, nil)
				fileReader.On("Stat", trackerFilePath).Return(nil, err)
			})

			It("forwards the error with context", func() {
				actual := tracker.TrackExecution(&configurationPath)

				Expect(actual).To(MatchError(err))
			})
		})

		Context("with an error with a tracker path that is not a file", func() {
			var errorMessage string

			BeforeEach(func() {
				errorMessage = fmt.Sprintf("'%s' should be a regular file", trackerFilePath)

				vcs.On("Root").Return(repositoryFakeRoot, nil)
				fileReader.On("Stat", trackerFilePath).
					Return(&FakeFileInfo{FileMode: os.ModeDir}, nil)
			})

			It("forwards the error with context", func() {
				actual := tracker.TrackExecution(&configurationPath)

				Expect(actual.Error()).To(ContainSubstring(errorMessage))
			})
		})

		Context("with an error when reading configuration", func() {
			err := fmt.Errorf("configuration read error")

			BeforeEach(func() {
				vcs.On("Root").Return(repositoryFakeRoot, nil)
				fileReader.On("Stat", trackerFilePath).
					Return(&FakeFileInfo{FileMode: os.ModePerm}, nil)
				fileReader.On("Read", configurationPath).
					Return(nil, err)
			})

			It("forwards the error with context", func() {
				actual := tracker.TrackExecution(&configurationPath)

				Expect(actual).To(MatchError(err))
			})
		})

		Context("with an error when unmarshalling configuration", func() {
			BeforeEach(func() {
				vcs.On("Root").Return(repositoryFakeRoot, nil)
				fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: os.ModePerm}, nil)
				fileReader.On("Read", configurationPath).
					Return([]byte("not json"), nil)
			})

			It("forwards the error with context", func() {
				actual := tracker.TrackExecution(&configurationPath)

				Expect(actual.Error()).To(ContainSubstring("cannot unmarshal configuration"))
			})
		})

		Context("with an error when unmarshalling configuration", func() {
			err := fmt.Errorf("header read error")

			BeforeEach(func() {
				vcs.On("Root").Return(repositoryFakeRoot, nil)
				fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: os.ModePerm}, nil)
				fileReader.On("Read", configurationPath).Return([]byte(configuration), nil)
				fileReader.On("Read", headerPath).
					Return(nil, err)
			})

			It("forwards the error with context", func() {
				actual := tracker.TrackExecution(&configurationPath)

				Expect(actual).To(MatchError(err))
			})
		})
	})
})

func base64Encode(contents string) string {
	return base64.StdEncoding.EncodeToString([]byte(contents))
}
