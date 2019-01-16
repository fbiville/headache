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

package vcs_test

import (
	. "github.com/fbiville/headache/vcs"
	"github.com/fbiville/headache/vcs_mocks"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("VCS", func() {

	var (
		t          GinkgoTInterface
		controller *gomock.Controller
		vcs        Vcs
		vcsMock    *vcs_mocks.Vcs
	)

	BeforeEach(func() {
		t = GinkgoT()
		controller = gomock.NewController(t)
		vcs = new(vcs_mocks.Vcs)
		vcsMock = vcs.(*vcs_mocks.Vcs)

	})

	AfterEach(func() {
		vcsMock.AssertExpectations(t)
		controller.Finish()
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

	It("retrieves file history", func() {
		vcsMock.On("Log", "--format=%at", "--", "somefile.go").Return(`1537974554
1537973963
1537970000
1537846444
1537844925
1499817600
`, nil)

		history, err := GetFileHistory(vcs, "somefile.go", FakeTime{})

		Expect(err).To(BeNil())
		Expect(history.CreationYear).To(Equal(2017))
		Expect(history.LastEditionYear).To(Equal(2018))
	})

	It("returns current year for unversioned files", func() {
		vcsMock.On("Log", "--format=%at", "--", "somefile.go").Return(``, nil)
		fakeTime := FakeTime{timestamp: fakeNow}
		currentYear := fakeTime.Now().Year()

		history, err := GetFileHistory(vcs, "somefile.go", fakeTime)

		Expect(err).To(BeNil())
		Expect(history.CreationYear).To(Equal(currentYear))
		Expect(history.LastEditionYear).To(Equal(currentYear))
	})

	It("returns the current year for the last edition year for file committed only once", func() {
		vcsMock.On("Log", "--format=%at", "--", "somefile.go").Return(`405561600
`, nil)
		fakeTime := FakeTime{timestamp: fakeNow}
		currentYear := fakeTime.Now().Year()

		history, err := GetFileHistory(vcs, "somefile.go", fakeTime)

		Expect(err).To(BeNil())
		Expect(history.CreationYear).To(Equal(1982))
		Expect(history.LastEditionYear).To(Equal(currentYear))
	})
})

type FakeTime struct {
	timestamp int64
}

func (t FakeTime) Now() time.Time {
	return time.Unix(t.timestamp, 0)
}

const fakeNow = 510278400 // 4th of March, 1986
