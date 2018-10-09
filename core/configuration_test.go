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

package core

import (
	"github.com/fbiville/headache/mocks"
	. "github.com/fbiville/headache/versioning"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"testing"
)

var (
	vcs Vcs
)

func TestConfigurationInitWithLineCommentStyle(t *testing.T) {
	I := NewGomegaWithT(t)
	controller := gomock.NewController(t)
	defer controller.Finish()
	vcs = new(mocks.Vcs)
	getChanges := func(Vcs, string, string, bool) ([]FileChange, error) {
		return []FileChange{
			{Path: "../fixtures/hello_world.txt"},
			{Path: "../fixtures/short-license.txt"}}, nil
	}

	configuration, err := parseConfiguration(Configuration{
		HeaderFile:   "../fixtures/license.txt",
		CommentStyle: "SlashSlash",
		Includes:     []string{"../fixtures/hello_*.txt"},
		Excludes:     []string{},
		TemplateData: map[string]string{
			"Owner": "ACME Labs",
		}}, RegularRunMode, nil, vcs, getChanges)

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`// Copyright {{.Year}} ACME Labs
//
// Some fictional license`))
	I.Expect(configuration.vcsChanges).To(Equal([]FileChange{{Path: "../fixtures/hello_world.txt"}}))
}

func TestConfigurationInitWithBlockCommentStyle(t *testing.T) {
	I := NewGomegaWithT(t)
	controller := gomock.NewController(t)
	defer controller.Finish()
	vcs = new(mocks.Vcs)
	getChanges := func(Vcs, string, string, bool) ([]FileChange, error) {
		return []FileChange{
			{Path: "../fixtures/hello_world_2017.txt"},
			{Path: "../fixtures/license.txt"}}, nil
	}

	configuration, err := parseConfiguration(Configuration{
		HeaderFile:   "../fixtures/license.txt",
		CommentStyle: "SlashStar",
		Includes:     []string{"../fixtures/*2017.txt"},
		Excludes:     []string{},
		TemplateData: map[string]string{
			"Owner": "ACME Labs",
		}}, RegularRunMode, nil, vcs, getChanges)

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`/*
 * Copyright {{.Year}} ACME Labs
 *
 * Some fictional license
 */`))
	I.Expect(configuration.vcsChanges).To(Equal([]FileChange{{Path: "../fixtures/hello_world_2017.txt"}}))
}

func TestHeaderDetectionRegexComputation(t *testing.T) {
	I := NewGomegaWithT(t)
	controller := gomock.NewController(t)
	defer controller.Finish()
	vcs = new(mocks.Vcs)
	getChanges := func(Vcs, string, string, bool) ([]FileChange, error) {
		return []FileChange{
			{Path: "../fixtures/hello_world_2017.txt"},
			{Path: "../fixtures/license.txt"}}, nil
	}

	configuration, err := parseConfiguration(Configuration{
		HeaderFile:   "../fixtures/short-license.txt",
		CommentStyle: "SlashStar",
		Includes:     []string{"../fixtures/*_2017.txt"},
		Excludes:     []string{},
		TemplateData: map[string]string{
			"Owner": "ACME Labs",
		}}, RegularRunMode, nil, vcs, getChanges)

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`/*
 * Copyright {{.Year}} ACME Labs
 */`))
	regex := configuration.HeaderRegex
	I.Expect(regex.String()).To(Equal("(?m)(?:\\/\\*\n)?(?:\\/{2}| \\*) ?\\QCopyright \\E.*\\Q \\E.*\\Q\\E\n?(?: \\*\\/)?"))
	I.Expect(regex.MatchString(configuration.HeaderContents)).To(BeTrue(), "Regex should match contents")
	I.Expect(regex.MatchString("// Copyright 2018 ACME Labs")).To(BeTrue(), "Regex should match contents in different comment style")
	I.Expect(regex.MatchString(`/*
 * Copyright 2018-2042 ACME World corporation
 */`)).To(BeTrue(), "Regex should match contents with different data")
	I.Expect(configuration.vcsChanges).To(Equal([]FileChange{{Path: "../fixtures/hello_world_2017.txt"}}))

}

func TestMatch(t *testing.T) {
	I := NewGomegaWithT(t)
	includes := []string{"../fixtures/*.txt"}
	excludes := []string{"../fixtures/*_with_header.txt"}

	I.Expect(match("../fixtures/bonjour_world.txt", includes, excludes)).To(BeTrue())
	I.Expect(match("../fixtures/bonjour_world.go", includes, excludes)).To(BeFalse())
	I.Expect(match("../fixtures/hello_world_with_header.txt", includes, excludes)).To(BeFalse())
}

func TestMatchOnlyFiles(t *testing.T) {
	I := NewGomegaWithT(t)

	I.Expect(match("../fixtures", []string{"../fixture*"}, []string{})).To(BeFalse())
}

func TestFailOnReservedYearParameter(t *testing.T) {
	I := NewGomegaWithT(t)

	configuration, err := parseConfiguration(Configuration{
		HeaderFile:   "../fixtures/short-license.txt",
		CommentStyle: "SlashStar",
		Includes:     []string{"../fixtures/*_2017.txt"},
		Excludes:     []string{},
		TemplateData: map[string]string{
			"Year": "2042",
		}}, RegularRunMode, nil, vcs, func(Vcs, string, string, bool) ([]FileChange, error) {
		panic("should not be called!")
	})

	I.Expect(configuration).To(BeNil())
	I.Expect(err).To(MatchError("Year is a reserved parameter and is automatically computed.\n" +
		"Please remove it from your configuration"))
}

func TestInitDryRunMode(t *testing.T) {
	I := NewGomegaWithT(t)
	controller := gomock.NewController(t)
	defer controller.Finish()
	vcs = new(mocks.Vcs)
	vcsMock := vcs.(*mocks.Vcs)
	vcsMock.On("Log", mock.AnythingOfType("[]string")).Return("", nil)
	inputConfig := Configuration{
		HeaderFile:   "../fixtures/short-license.txt",
		CommentStyle: "SlashStar",
		Includes:     []string{"../fixtures/*world.txt"},
		Excludes:     []string{},
		TemplateData: map[string]string{}}

	configuration, err := parseConfiguration(inputConfig,
		DryRunInitMode,
		nil,
		vcs,
		func(Vcs, string, string, bool) ([]FileChange, error) { panic("nope!") })

	I.Expect(err).To(BeNil())
	changes := paths(configuration.vcsChanges)
	I.Expect(len(changes)).To(Equal(3))
	I.Expect(changes).To(ContainElement("../fixtures/hello_world.txt"))
	I.Expect(changes).To(ContainElement("../fixtures/hello_ignored_world.txt"))
	I.Expect(changes).To(ContainElement("../fixtures/bonjour_world.txt"))
}

func paths(changes []FileChange) []string {
	result := make([]string, len(changes))
	for i, change := range changes {
		result[i] = change.Path
	}
	return result
}


