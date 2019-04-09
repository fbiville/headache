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
		Expect(style.GetClosingSymbol().Value).To(Equal(""))
		Expect(style.GetOpeningSymbol().Value).To(Equal(""))
		Expect(style.GetContinuationSymbol().Value).To(Equal("# "))
	})
})
