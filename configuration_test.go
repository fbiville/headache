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

package main_test

import (
	. "github.com/fbiville/header"
	. "github.com/onsi/gomega"
	"testing"
)

func TestConfigurationInitWithLineCommentStyle(t *testing.T) {
	I := NewGomegaWithT(t)

	configuration, err := ParseConfiguration(Configuration{
		HeaderFile:   "fixtures/license.txt",
		CommentStyle: "SlashSlash",
		Includes:     []string{"*.txt"},
		Excludes:     []string{},
		TemplateData: map[string]string{
			"Year":  "2018",
			"Owner": "ACME Labs",
		}})

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`// Copyright 2018 ACME Labs
//
// Some fictional license`))
	I.Expect(configuration.Includes).To(Equal([]string{"*.txt"}))
	I.Expect(configuration.Excludes).To(BeEmpty())
}

func TestConfigurationInitWithBlockCommentStyle(t *testing.T) {
	I := NewGomegaWithT(t)

	configuration, err := ParseConfiguration(Configuration{
		HeaderFile:   "fixtures/license.txt",
		CommentStyle: "SlashStar",
		Includes:     []string{"*.txt"},
		Excludes:     []string{},
		TemplateData: map[string]string{
			"Year":  "2018",
			"Owner": "ACME Labs",
		}})

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`/*
 * Copyright 2018 ACME Labs
 *
 * Some fictional license
 */`))
	I.Expect(configuration.Includes).To(Equal([]string{"*.txt"}))
	I.Expect(configuration.Excludes).To(BeEmpty())
}

func TestHeaderDetectionRegexComputation(t *testing.T) {
	I := NewGomegaWithT(t)
	configuration, err := ParseConfiguration(Configuration{
		HeaderFile:   "fixtures/short-license.txt",
		CommentStyle: "SlashStar",
		Includes:     []string{"*.txt"},
		Excludes:     []string{},
		TemplateData: map[string]string{
			"Year":  "2018",
			"Owner": "ACME Labs",
		}})

	I.Expect(err).To(BeNil())
	I.Expect(configuration.HeaderContents).To(Equal(`/*
 * Copyright 2018 ACME Labs
 */`))
	regex := configuration.HeaderRegex
	I.Expect(regex.String()).To(Equal("(?m)(?:\\/\\*\n)?(?:\\/{2}| \\*) ?\\QCopyright \\E.*\\Q \\E.*\\Q\\E\n?(?: \\*\\/)?"))
	I.Expect(regex.MatchString(configuration.HeaderContents)).To(BeTrue(), "Regex should match contents")
	I.Expect(regex.MatchString("// Copyright 2018 ACME Labs")).To(BeTrue(), "Regex should match contents in different comment style")
	I.Expect(regex.MatchString(`/*
 * Copyright 2018-2042 ACME World corporation
 */`)).To(BeTrue(), "Regex should match contents with different data")

}
