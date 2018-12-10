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
	"fmt"
	. "github.com/fbiville/headache/versioning"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"testing"
)

func TestConfigurationInitWithLineCommentStyle(t *testing.T) {
	I := NewGomegaWithT(t)
	controller := gomock.NewController(t)
	defer controller.Finish()
	getChanges := func(Vcs, string) ([]FileChange, error) {
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
		}}, getChanges)

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`// Copyright {{.Year}} ACME Labs
//
// Some fictional license`))
	I.Expect(onlyPaths(configuration.vcsChanges)).To(Equal([]FileChange{{Path: "../fixtures/hello_world.txt"}}))
}

func TestConfigurationInitWithBlockCommentStyle(t *testing.T) {
	I := NewGomegaWithT(t)
	controller := gomock.NewController(t)
	defer controller.Finish()
	getChanges := func(Vcs, string) ([]FileChange, error) {
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
		}}, getChanges)

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`/*
 * Copyright {{.Year}} ACME Labs
 *
 * Some fictional license
 */`))
	I.Expect(onlyPaths(configuration.vcsChanges)).To(Equal([]FileChange{{Path: "../fixtures/hello_world_2017.txt"}}))
}

func TestHeaderDetectionRegexComputation(t *testing.T) {
	I := NewGomegaWithT(t)
	controller := gomock.NewController(t)
	defer controller.Finish()
	getChanges := func(Vcs, string) ([]FileChange, error) {
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
		}}, getChanges)

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`/*
 * Copyright {{.Year}} ACME Labs
 */`))
	regex := configuration.HeaderRegex
	I.Expect(regex.String()).To(Equal(`(?im)(?:\/\*\n)?(?:\/{2}| \*)[ \t]*\QCopyright \E.*\Q \E.*\Q\E[ \t\.]*\n?(?:(?:\/{2}| \*) ?\n)*(?: \*\/)?`))
	I.Expect(regex.MatchString(configuration.HeaderContents)).To(BeTrue(), "Regex should match contents")
	I.Expect(regex.MatchString("// Copyright 2018 ACME Labs")).To(BeTrue(), "Regex should match contents in different comment style")
	I.Expect(regex.MatchString(`/*
 * Copyright 2018-2042 ACME World corporation
 */`)).To(BeTrue(), "Regex should match contents with different data")
	I.Expect(onlyPaths(configuration.vcsChanges)).To(Equal([]FileChange{{Path: "../fixtures/hello_world_2017.txt"}}))

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
		}}, func(Vcs, string) ([]FileChange, error) {
		return nil, fmt.Errorf("should not be called")
	})

	I.Expect(configuration).To(BeNil())
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
