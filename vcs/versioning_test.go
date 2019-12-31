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

package vcs_test

import (
	. "github.com/fbiville/headache/vcs"
	"github.com/fbiville/headache/vcs_mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("VCS", func() {

	var (
		t       GinkgoTInterface
		vcs     Vcs
		vcsMock *vcs_mocks.Vcs
	)

	BeforeEach(func() {
		t = GinkgoT()
		vcs = new(vcs_mocks.Vcs)
		vcsMock = vcs.(*vcs_mocks.Vcs)

	})

	AfterEach(func() {
		vcsMock.AssertExpectations(t)
	})

	It("retrieves committed changes", func() {
		vcsMock.On("Diff", "--name-status", "origin/master..HEAD").Return(`M	.gitignore
M	configuration.go
D	header.go
D	header_test.go
R099	line_comment.go	core/line_comment.go
A	license-header.txt
`, nil)

		changes, err := GetCommittedChanges(vcs, "origin/master")

		Expect(err).To(BeNil())
		Expect(changes).To(Equal([]FileChange{
			{Path: ".gitignore"},
			{Path: "configuration.go"},
			{Path: "core/line_comment.go"},
			{Path: "license-header.txt"},
		}))
	})

	It("retrieves uncommitted files", func() {
		vcsMock.On("Status", "--porcelain").Return(` M Gopkg.lock
 D main.go
?? build.sh
?? git.go
`, nil)

		changes, err := GetUncommittedChanges(vcs)

		Expect(err).To(BeNil())
		Expect(changes).To(Equal([]FileChange{
			{Path: "Gopkg.lock"},
			{Path: "build.sh"},
			{Path: "git.go"},
		}))
	})

	It("retrieves no changes when everything is committed", func() {
		vcsMock.On("Status", "--porcelain").Return(` M Gopkg.lock
 D main.go
?? build.sh
?? git.go
`, nil)

		changes, err := GetUncommittedChanges(vcs)

		Expect(err).To(BeNil())
		Expect(changes).To(Equal([]FileChange{
			{Path: "Gopkg.lock"},
			{Path: "build.sh"},
			{Path: "git.go"},
		}))
	})

	Describe("retrieves file history", func() {

		var (
			logArguments []interface{}
			fakeTime FakeTime
		)

		BeforeEach(func() {
			logArguments = []interface{}{"--follow", "--name-status", "--format=%at", "--"}
			fakeTime = FakeTime{timestamp: fakeNow}
		})

		It("retrieves the first and last commit years", func() {
			vcsMock.On("Log", append(logArguments, "somefile.go")...).Return(`1537974554

M	somefile.go
1537844925

M	somefile.go
1499817600

A	cmd/commands/ginkgo_suite_test.go
`, nil)

			history, err := GetFileHistory(vcs, "somefile.go", FakeTime{})

			Expect(err).To(BeNil())
			Expect(history.CreationYear).To(Equal(2017))
			Expect(history.LastEditionYear).To(Equal(2018))
		})

		It("returns current year for unversioned files", func() {
			vcsMock.On("Log", append(logArguments, "somefile.go")...).Return(``, nil)
			currentYear := fakeTime.Now().Year()

			history, err := GetFileHistory(vcs, "somefile.go", fakeTime)

			Expect(err).To(BeNil())
			Expect(history.CreationYear).To(Equal(currentYear))
			Expect(history.LastEditionYear).To(Equal(currentYear))
		})

		It("returns the commit year for both creation and last edition year when file has been committed only once", func() {
			vcsMock.On("Log", append(logArguments, "somefile.go")...).Return(`405561600

A	somefile.go
`, nil)
			history, err := GetFileHistory(vcs, "somefile.go", fakeTime)

			Expect(err).To(BeNil())
			Expect(history.CreationYear).To(Equal(1982))
			Expect(history.LastEditionYear).To(Equal(1982))
		})

		It("ignores commits which are pure renames", func() {
			vcsMock.On("Log", append(logArguments, "pkg/core/ginkgo_suite_test.go")...).
				Return(`1551657600

R100	cmd/commands/ginkgo_suite_test.go	pkg/core/ginkgo_suite_test.go
1531499156

A	cmd/commands/ginkgo_suite_test.go
`, nil)
			history, err := GetFileHistory(vcs, "pkg/core/ginkgo_suite_test.go", fakeTime)

			Expect(err).To(BeNil())
			Expect(history.CreationYear).To(Equal(2018))
			Expect(history.LastEditionYear).To(Equal(2018))
		})

		It("ignores commits which are pure copies", func() {
			vcsMock.On("Log", append(logArguments, "pkg/core/ginkgo_suite_test.go")...).
				Return(`1551657600

C100	cmd/commands/ginkgo_suite_test.go	pkg/core/ginkgo_suite_test.go
1531499156

A	cmd/commands/ginkgo_suite_test.go
`, nil)
			history, err := GetFileHistory(vcs, "pkg/core/ginkgo_suite_test.go", fakeTime)

			Expect(err).To(BeNil())
			Expect(history.CreationYear).To(Equal(2018))
			Expect(history.LastEditionYear).To(Equal(2018))
		})

		It("should fail on invalid output", func() {
			vcsMock.On("Log", append(logArguments, "somefile.go")...).Return(`wat
saywat
`, nil)

			_, err := GetFileHistory(vcs, "somefile.go", fakeTime)

			Expect(err).To(MatchError("could not parse timestamp (line 1) of file \"somefile.go\" history. Full commit log below\nwat\nsaywat\n"))
		})

	})
})

type FakeTime struct {
	timestamp int64
}

func (t FakeTime) Now() time.Time {
	return time.Unix(t.timestamp, 0)
}

const fakeNow = 510278400 // 4th of March, 1986
