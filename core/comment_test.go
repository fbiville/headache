package core_test

import (
	. "github.com/fbiville/headache/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Comment", func() {

	It("matches Dash comment style", func() {
		style := ParseCommentStyle("Hash")

		Expect(style.GetName()).To(Equal("Hash"))
		Expect(style.GetClosingString()).To(Equal(""))
		Expect(style.GetOpeningString()).To(Equal(""))
		Expect(style.GetString()).To(Equal("# "))
	})
})
