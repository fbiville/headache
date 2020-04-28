/*
 * Copyright 2019 Florent Biville (@fbiville)
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
	"encoding/json"
	"io/ioutil"
	"sort"
	"strings"

	. "github.com/fbiville/headache/internal/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Comment styles", func() {

	DescribeTable("are properly defined",
		func(name, openingStr, closingStr, str string) {
			style := ParseCommentStyle(name)

			Expect(style.GetName()).To(Equal(name))
			Expect(style.GetClosingString()).To(Equal(closingStr))
			Expect(style.GetOpeningString()).To(Equal(openingStr))
			Expect(style.GetString()).To(Equal(str))
		},
		Entry("matches SlashStar comment style", "SlashStar", "/*", " */", " * "),
		Entry("matches SlashSlash comment style", "SlashSlash", "", "", "// "),
		Entry("matches Hash comment style", "Hash", "", "", "# "),
		Entry("matches DashDash comment style", "DashDash", "", "", "-- "),
		Entry("matches SemiColon comment style", "SemiColon", "", "", "; "),
		Entry("matches REM comment style", "REM", "", "", "REM "),
		Entry("matches SlashStarStar comment style", "SlashStarStar", "/**", " */", " * "),
	)

	It("include only the following", func() {
		catalog := namesOf(SupportedStyleCatalog())

		Expect(catalog).To(Equal([]string{
			"DashDash",
			"Hash",
			"REM",
			"SemiColon",
			"SlashSlash",
			"SlashStar",
			"SlashStarStar",
			"XML",
		}))
		Expect(lowerAll(catalog)).To(Equal(sortedStylesInSchema("../../../docs/schema.json")),
			"Expected all declared styles to be included in JSON schema")
	})

	It("are parsed case-insensitively", func() {
		style := ParseCommentStyle("slashslash")

		Expect(style.GetName()).To(Equal("SlashSlash"))
	})
})

func namesOf(styles map[string]CommentStyle) []string {
	result := make([]string, len(styles))
	i := 0
	for key, style := range styles {
		Expect(key).To(Equal(style.GetName()), "expected indexed style name to be style name")
		result[i] = style.GetName()
		i++
	}
	sort.Strings(result)
	return result
}

func sortedStylesInSchema(schemaFileLocation string) []string {
	bytes, err := ioutil.ReadFile(schemaFileLocation)
	Expect(err).NotTo(HaveOccurred())
	var schema HeadacheSchema
	err = json.Unmarshal(bytes, &schema)
	return schema.SortedStyleNames()
}

type HeadacheSchema struct {
	Properties HeadacheProperties `json:"properties"`
}

type HeadacheProperties struct {
	Style CommentStyleProperty `json:"style"`
}

type CommentStyleProperty struct {
	Names []string `json:"enum"`
}

func (schema *HeadacheSchema) SortedStyleNames() []string {
	result := schema.Properties.Style.Names
	sort.Strings(result)
	return result
}

func lowerAll(catalog []string) []string {
	for i, value := range catalog {
		catalog[i] = strings.ToLower(value)
	}
	return catalog
}
