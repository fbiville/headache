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
	"fmt"
	"github.com/fbiville/headache/core"
	. "github.com/fbiville/headache/fs"
	"github.com/fbiville/headache/fs_mocks"
	"github.com/fbiville/headache/vcs_mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strings"
	"time"
)

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
			currentHeaderFile     string
			currentHeaderContents string
			currentData           map[string]string
			currentConfiguration  *core.Configuration

			fakeRepositoryRoot string
			trackerFilePath    string
		)

		BeforeEach(func() {
			configurationPath := "/path/to/headache.json"
			currentHeaderFile = "current-header-file"
			currentHeaderContents = "some\nheader"
			currentData = map[string]string{"foo": "bar"}
			currentConfiguration = &core.Configuration{
				HeaderFile:   currentHeaderFile,
				TemplateData: currentData,
				Path:         &configurationPath,
			}
			fakeRepositoryRoot = "/path/to"
			trackerFilePath = fakeRepositoryRoot + "/.headache-run"
		})

		Describe("for the first execution ever of headache", func() {

			Context("without any errors", func() {
				BeforeEach(func() {
					fileReader.On("Read", currentHeaderFile).
						Return([]byte(currentHeaderContents), nil)
					vcs.On("Root").
						Return(fakeRepositoryRoot, nil)
					fileReader.On("Stat", trackerFilePath).
						Return(nil, os.ErrNotExist)
				})

				It("returns the current contents as current and former, and returns no revision", func() {
					versionedTemplate, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

					Expect(err).NotTo(HaveOccurred())
					Expect(versionedTemplate.Revision).To(BeEmpty())
					Expect(versionedTemplate.Current).To(Equal(versionedTemplate.Previous))
					Expect(versionedTemplate.Current.Data).To(Equal(currentData))
					Expect(strings.Join(versionedTemplate.Current.Lines, "\n")).To(Equal(currentHeaderContents))
				})
			})

			Context("with an error when reading current header file", func() {
				expectedErr := fmt.Errorf("current header error")

				BeforeEach(func() {
					fileReader.On("Read", currentHeaderFile).
						Return(nil, expectedErr)
				})

				It("forwards the error", func() {
					_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

					Expect(err).To(MatchError(expectedErr))
				})
			})

			Context("with an error proceeding to repository root", func() {
				expectedErr := fmt.Errorf("root error")

				BeforeEach(func() {
					fileReader.On("Read", currentHeaderFile).
						Return([]byte(currentHeaderContents), nil)
					vcs.On("Root").
						Return("", expectedErr)
				})

				It("forwards the error", func() {
					_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

					Expect(err).To(MatchError(expectedErr))
				})
			})

			Context("with an error \"stat'ing\" execution tracker file", func() {
				expectedErr := fmt.Errorf("tracker file error")

				BeforeEach(func() {
					fileReader.On("Read", currentHeaderFile).
						Return([]byte(currentHeaderContents), nil)
					vcs.On("Root").
						Return(fakeRepositoryRoot, nil)
					fileReader.On("Stat", trackerFilePath).
						Return(nil, expectedErr)
				})

				It("forwards the error", func() {
					_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

					Expect(err).To(MatchError(expectedErr))
				})
			})

			Context("with an invalid execution tracker path", func() {
				BeforeEach(func() {
					fileReader.On("Read", currentHeaderFile).
						Return([]byte(currentHeaderContents), nil)
					vcs.On("Root").
						Return(fakeRepositoryRoot, nil)
					fileReader.On("Stat", trackerFilePath).
						Return(&FakeFileInfo{FileMode: os.ModeDir}, nil)
				})

				It("forwards the error", func() {
					_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

					Expect(err).To(MatchError(fmt.Sprintf("'%s' should be a regular file", trackerFilePath)))
				})
			})
		})

		Describe("after an execution of headache", func() {

			previousHeaderFile := "previous-header-file"
			previousHeaderContents := "some\nprevious\nheader"
			lastExecutionRevision := "b7c7db75695d8ffff37d556780a616ffbf5c2696"
			configurationAtRevision := fmt.Sprintf(`{
  "headerFile": "%s",
  "data": {"some": "thing"}
}`, previousHeaderFile)

			BeforeEach(func() {
				fileReader.On("Read", currentHeaderFile).Return([]byte(currentHeaderContents), nil)
				vcs.On("Root").Return(fakeRepositoryRoot, nil)
				fileReader.On("Stat", trackerFilePath).Return(&FakeFileInfo{FileMode: 0777}, nil)
			})

			Describe("when the former configuration path has not been serialized (backward compat)", func() {

				Context("with no errors", func() {

					BeforeEach(func() {
						fileReader.On("Read", trackerFilePath).
							Return([]byte("# Generated by headache | 1547741491 -- commit me!"), nil)
						vcs.On("LatestRevision", trackerFilePath).
							Return(lastExecutionRevision, nil)
						vcs.On("ShowContentAtRevision", *currentConfiguration.Path, lastExecutionRevision).
							Return(configurationAtRevision, nil)
						fileReader.On("Read", previousHeaderFile).
							Return([]byte(previousHeaderContents), nil)
					})

					It("reads the current configuration path at the last revision", func() {
						versionedTemplate, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

						Expect(err).NotTo(HaveOccurred())
						Expect(versionedTemplate.Revision).To(Equal(lastExecutionRevision))
						Expect(versionedTemplate.Current.Data).To(Equal(currentData))
						Expect(strings.Join(versionedTemplate.Current.Lines, "\n")).To(Equal(currentHeaderContents))
						Expect(versionedTemplate.Previous.Data).To(Equal(map[string]string{"some": "thing"}))
						Expect(strings.Join(versionedTemplate.Previous.Lines, "\n")).To(Equal(previousHeaderContents))
					})
				})

				Context("with an error when reading the tracker file", func() {
					expectedErr := fmt.Errorf("tracker file error")

					BeforeEach(func() {
						fileReader.On("Read", trackerFilePath).
							Return(nil, expectedErr)

					})

					It("forwards the error", func() {
						_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

						Expect(err).To(MatchError(expectedErr))
					})
				})

				Context("with an error when getting the last execution revision", func() {
					BeforeEach(func() {
						fileReader.On("Read", trackerFilePath).
							Return([]byte("# Generated by headache | 1547741491 -- commit me!"), nil)
						vcs.On("LatestRevision", trackerFilePath).
							Return("", fmt.Errorf("some error"))
					})

					It("forwards the error", func() {
						_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

						Expect(err).To(MatchError("could not detect previous execution's revision"))
					})
				})

				Context("with an error when getting the tracker file contents at last execution revision", func() {
					expectedErr := fmt.Errorf("tracker file at last revision error")

					BeforeEach(func() {
						fileReader.On("Read", trackerFilePath).
							Return([]byte("# Generated by headache | 1547741491 -- commit me!"), nil)
						vcs.On("LatestRevision", trackerFilePath).
							Return(lastExecutionRevision, nil)
						vcs.On("ShowContentAtRevision", *currentConfiguration.Path, lastExecutionRevision).
							Return("", expectedErr)
					})

					It("forwards the error", func() {
						_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

						Expect(err).To(MatchError(expectedErr))
					})
				})

				Context("with an error when reading the header file configured at the last revision", func() {
					expectedErr := fmt.Errorf("previous header file error")

					BeforeEach(func() {
						fileReader.On("Read", trackerFilePath).
							Return([]byte("# Generated by headache | 1547741491 -- commit me!"), nil)
						vcs.On("LatestRevision", trackerFilePath).
							Return(lastExecutionRevision, nil)
						vcs.On("ShowContentAtRevision", *currentConfiguration.Path, lastExecutionRevision).
							Return(configurationAtRevision, nil)
						fileReader.On("Read", previousHeaderFile).
							Return(nil, expectedErr)
					})

					It("forwards the error", func() {
						_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

						Expect(err).To(MatchError(expectedErr))
					})
				})
			})

			Describe("when the former configuration path has been serialized (legacy configuration)", func() {
				Context("with no errors", func() {

					serializedConfigurationPath := "/path/to/former/headache.json"

					BeforeEach(func() {
						fileReader.On("Read", trackerFilePath).
							Return([]byte("# Generated by headache | 1547741491 -- commit me!\nconfiguration:"+serializedConfigurationPath), nil)
						vcs.On("LatestRevision", trackerFilePath).
							Return(lastExecutionRevision, nil)
						vcs.On("ShowContentAtRevision", serializedConfigurationPath, lastExecutionRevision).
							Return(configurationAtRevision, nil)
						fileReader.On("Read", previousHeaderFile).
							Return([]byte(previousHeaderContents), nil)
					})

					It("reads the former configuration path at the last revision", func() {
						versionedTemplate, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

						Expect(err).NotTo(HaveOccurred())
						Expect(versionedTemplate.Revision).To(Equal(lastExecutionRevision))
						Expect(versionedTemplate.Current.Data).To(Equal(currentData))
						Expect(strings.Join(versionedTemplate.Current.Lines, "\n")).To(Equal(currentHeaderContents))
						Expect(versionedTemplate.Previous.Data).To(Equal(map[string]string{"some": "thing"}))
						Expect(strings.Join(versionedTemplate.Previous.Lines, "\n")).To(Equal(previousHeaderContents))
					})
				})
			})

			Describe("when the former configuration and header template have been serialized", func() {

				previousConfiguration := `{
  "headerFile": "./license-header.txt",
  "style": "SlashStar",
  "includes": ["**/*.go"],
  "excludes": ["vendor/**/*", "*_mocks/**/*"],
  "data": {
    "Owner": "Florent Biville (@fbiville)"
  }
}`

				Context("with no errors", func() {
					BeforeEach(func() {
						fileReader.On("Read", trackerFilePath).
							Return([]byte(fmt.Sprintf("# Generated by headache | 1547741491 -- commit me!\nencoded_configuration:%s\nencoded_header:%s",
								base64Encode(previousConfiguration),
								base64Encode(previousHeaderContents))), nil)
						vcs.On("LatestRevision", trackerFilePath).
							Return(lastExecutionRevision, nil)
					})

					It("reads the encoded data", func() {
						versionedTemplate, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

						Expect(err).NotTo(HaveOccurred())
						Expect(versionedTemplate.Revision).To(Equal(lastExecutionRevision))
						Expect(versionedTemplate.Current.Data).To(Equal(currentData))
						Expect(strings.Join(versionedTemplate.Current.Lines, "\n")).To(Equal(currentHeaderContents))
						Expect(versionedTemplate.Previous.Data).To(Equal(map[string]string{"Owner": "Florent Biville (@fbiville)"}))
						Expect(strings.Join(versionedTemplate.Previous.Lines, "\n")).To(Equal(previousHeaderContents))
					})
				})

				Context("with an error when reading misencoded configuration", func() {
					BeforeEach(func() {
						fileReader.On("Read", trackerFilePath).
							Return([]byte(fmt.Sprintf("# Generated by headache | 1547741491 -- commit me!\nencoded_configuration:%s\nencoded_header:%s",
								"not base 64 encoded",
								base64Encode(previousHeaderContents))), nil)
						vcs.On("LatestRevision", trackerFilePath).
							Return(lastExecutionRevision, nil)
					})

					It("forwards the error", func() {
						_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

						Expect(err).To(MatchError("could not decode encoded configuration: " +
							"illegal base64 data at input byte 3"))
					})
				})

				Context("with an error when reading encoded malformed configuration", func() {
					BeforeEach(func() {
						fileReader.On("Read", trackerFilePath).
							Return([]byte(fmt.Sprintf("# Generated by headache | 1547741491 -- commit me!\nencoded_configuration:%s\nencoded_header:%s",
								base64Encode("not JSON"),
								base64Encode(previousHeaderContents))), nil)
						vcs.On("LatestRevision", trackerFilePath).
							Return(lastExecutionRevision, nil)
					})

					It("forwards the error", func() {
						_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

						Expect(err).To(MatchError("could not unmarshal decoded configuration: " +
							"invalid character 'o' in literal null (expecting 'u')"))
					})
				})

				Context("with an error when not finding encoded header template", func() {
					BeforeEach(func() {
						fileReader.On("Read", trackerFilePath).
							Return([]byte(fmt.Sprintf("# Generated by headache | 1547741491 -- commit me!\nencoded_configuration:%s",
								base64Encode(previousConfiguration))), nil)
						vcs.On("LatestRevision", trackerFilePath).
							Return(lastExecutionRevision, nil)
					})

					It("forwards the error", func() {
						_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

						Expect(err).To(MatchError("cannot retrieve encoded header template"))
					})
				})

				Context("with an error when reading misencoded header template", func() {
					BeforeEach(func() {
						fileReader.On("Read", trackerFilePath).
							Return([]byte(fmt.Sprintf("# Generated by headache | 1547741491 -- commit me!\nencoded_configuration:%s\nencoded_header:%s",
								base64Encode(previousConfiguration),
								"not base 64")), nil)
						vcs.On("LatestRevision", trackerFilePath).
							Return(lastExecutionRevision, nil)
					})

					It("forwards the error", func() {
						_, err := tracker.RetrieveVersionedTemplate(currentConfiguration)

						Expect(err).To(MatchError("could not decode encoded header template: " +
							"illegal base64 data at input byte 3"))
					})
				})
			})
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

type FixedClock struct{}

func (*FixedClock) Now() time.Time {
	return time.Unix(42, 42)
}
