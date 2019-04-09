package core_test

import (
	. "github.com/fbiville/headache/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("Comment styles", func() {

	It("include Dash (#)", func() {
		style := ParseCommentStyle("Hash")

		Expect(style.GetName()).To(Equal("Hash"))
		Expect(style.GetOpeningSymbol().Value).To(Equal(""))
		Expect(style.GetContinuationSymbol().Value).To(Equal("# "))
		Expect(style.GetClosingSymbol().Value).To(Equal(""))
	})

	Describe("includes SlashStar (/* */)", func() {
		const licenseText = `
/*
Copyright 2018 The Knative authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
`
		var slashStarStyle CommentStyle

		BeforeEach(func() {
			slashStarStyle = ParseCommentStyle("SlashStar")
		})

		It("defining the relevant symbols", func() {
			Expect(slashStarStyle.GetName()).To(Equal("SlashStar"))
			Expect(slashStarStyle.GetOpeningSymbol().Value).To(Equal("/*"))
			Expect(slashStarStyle.GetContinuationSymbol()).To(Equal(&CommentSymbol{Value: " * ", Optional: true}))
			Expect(slashStarStyle.GetClosingSymbol().Value).To(Equal(" */"))
		})

		It("matches SlashStar headers without continuation symbols", func() {
			regex, err := ComputeDetectionRegex(strings.Split(licenseText, "\n"), map[string]string{})

			Expect(err).NotTo(HaveOccurred())
			Expect(licenseText).To(MatchRegexp(regex))
		})
	})

})
