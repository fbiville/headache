package helper_test

import (
	. "github.com/fbiville/headache/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Maps utilities", func() {

	It("extracts keys", func() {
		Expect(Keys(nil)).To(BeEmpty())
		Expect(Keys(map[string]string{})).To(BeEmpty())
		Expect(Keys(map[string]string{"foo": "_"})).To(ConsistOf("foo"))
	})

	It("extracts keys in alphabetical order", func() {
		keys := Keys(map[string]string{"foo": "_", "baz": "_"})

		Expect(keys).To(HaveLen(2))
		Expect(keys[0]).To(Equal("baz"))
		Expect(keys[1]).To(Equal("foo"))
	})
})
