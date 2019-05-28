package core_test

import (
	. "github.com/fbiville/headache/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Comment", func() {

	It("matches Hash comment style", func() {
		style := ParseCommentStyle("Hash")

		Expect(style.GetName()).To(Equal("Hash"))
		Expect(style.GetClosingString()).To(Equal(""))
		Expect(style.GetOpeningString()).To(Equal(""))
		Expect(style.GetString()).To(Equal("# "))
	})

	It("matches DashDash comment style", func() {
		style := ParseCommentStyle("DashDash")

		Expect(style.GetName()).To(Equal("DashDash"))
		Expect(style.GetClosingString()).To(Equal(""))
		Expect(style.GetOpeningString()).To(Equal(""))
		Expect(style.GetString()).To(Equal("-- "))
	})

	It("matches SemiColon comment style", func() {
		style := ParseCommentStyle("SemiColon")

		Expect(style.GetName()).To(Equal("SemiColon"))
		Expect(style.GetClosingString()).To(Equal(""))
		Expect(style.GetOpeningString()).To(Equal(""))
		Expect(style.GetString()).To(Equal("; "))
	})
})
