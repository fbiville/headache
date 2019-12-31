/*
 * Copyright 2019 Florent Biville (@fbiville)
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
	styles "github.com/fbiville/headache/core/comment_styles"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Template parser", func() {

	var (
		template          core.HeaderTemplate
		yearRangeTemplate core.HeaderTemplate
		legacyTemplate    core.HeaderTemplate
	)

	BeforeEach(func() {
		template = core.HeaderTemplate{
			Lines: []string{"Copyright (c) {{.StartYear}} -- {{.EndYear}} {{.Author}}"},
			Data:  map[string]string{"Author": "Florent"},
		}
		legacyTemplate = core.HeaderTemplate{
			Lines: []string{"Copyright (c) {{.Year}} {{.Author}}"},
			Data:  map[string]string{"Author": "Florent"},
		}
		yearRangeTemplate = core.HeaderTemplate{
			Lines: []string{"Copyright (c) {{.YearRange}} {{.Author}}"},
			Data:  map[string]string{"Author": "Florent"},
		}
	})

	It("preserves the start and end year parameter for later substitution", func() {
		versionedTemplate := &core.VersionedHeaderTemplate{
			Previous: &template,
			Current:  &template,
			Revision: "",
		}

		result, err := core.ParseTemplate(versionedTemplate, styles.Hash{})

		Expect(err).NotTo(HaveOccurred())
		Expect(result.ActualContent).To(Equal("# Copyright (c) {{.StartYear}} -- {{.EndYear}} Florent"))
	})

	It("preserves the year range parameter for later substitution", func() {
		versionedTemplate := &core.VersionedHeaderTemplate{
			Previous: &yearRangeTemplate,
			Current:  &yearRangeTemplate,
			Revision: "",
		}

		result, err := core.ParseTemplate(versionedTemplate, styles.Hash{})

		Expect(err).NotTo(HaveOccurred())
		Expect(result.ActualContent).To(Equal("# Copyright (c) {{.YearRange}} Florent"))
	})

	It("replaces the legacy year range parameter with the newer parameters for later substitution", func() {
		versionedTemplate := &core.VersionedHeaderTemplate{
			Previous: &legacyTemplate,
			Current:  &legacyTemplate,
			Revision: "",
		}

		result, err := core.ParseTemplate(versionedTemplate, styles.Hash{})

		Expect(err).NotTo(HaveOccurred())
		Expect(result.ActualContent).To(Equal("# Copyright (c) {{.YearRange}} Florent"))
	})

	It("computes a regex that detect headers with newline differences", func() {
		versionedTemplate := &core.VersionedHeaderTemplate{
			Previous: &core.HeaderTemplate{Lines: []string{"hello", "world"}, Data: map[string]string{}},
			Current:  &core.HeaderTemplate{Lines: []string{"hello", "world"}, Data: map[string]string{}},
			Revision: "",
		}

		result, err := core.ParseTemplate(versionedTemplate, styles.SlashSlash{})

		Expect(err).NotTo(HaveOccurred())
		Expect(result.ActualContent).To(Equal("// hello\n// world"))
		Expect(result.DetectionRegex.MatchString("// hello\n// world")).To(BeTrue(), "matches exact contents")
		Expect(result.DetectionRegex.MatchString("\n// hello\n\n\n\n// world\n\n")).To(BeTrue(), "matches with extra newlines")
		Expect(result.DetectionRegex.MatchString("\n// hello\n// \n// \n\n// world\n// \n\n")).To(BeTrue(), "matches with extra commented empty lines")
	})
})
