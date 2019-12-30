package core_test

import (
	. "github.com/fbiville/headache/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"sort"
)

var _ = Describe("Comment styles", func() {

	It("includes the following", func() {
		styles := SupportedStyleCatalog()

		Expect(namesOf(styles)).To(Equal([]string{
			"DashDash",
			"Hash",
			"REM",
			"SemiColon",
			"SlashSlash",
			"SlashStar",
			"SlashStarStar",
		}))
	})

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
