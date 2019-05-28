package core_test

import (
	. "github.com/fbiville/headache/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Comment", func() {

	DescribeTable("matches comment style",
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
	)
})
