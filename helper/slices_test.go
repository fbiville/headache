/*
 * Copyright 2018-2019 Florent Biville (@fbiville)
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

package helper_test

import (
	. "github.com/fbiville/headache/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing/quick"
)

var _ = Describe("String slices", func() {

	var t GinkgoTInterface

	BeforeEach(func() {
		t = GinkgoT()
	})

	It("compares slices", func() {
		Expect(SliceEqual(nil, nil)).To(BeTrue())
		Expect(SliceEqual(nil, []string{})).To(BeFalse())
		Expect(SliceEqual([]string{}, nil)).To(BeFalse())
		Expect(SliceEqual([]string{}, []string{})).To(BeTrue())
		Expect(SliceEqual([]string{"a"}, []string{})).To(BeFalse())
		Expect(SliceEqual([]string{"a"}, []string{"b"})).To(BeFalse())
		Expect(SliceEqual([]string{"a", "c"}, []string{"b", "c"})).To(BeFalse())
		Expect(SliceEqual([]string{"a", "c"}, []string{"a", "c"})).To(BeTrue())
		Expect(SliceEqual([]string{"a", "c"}, []string{"c", "a"})).To(BeFalse())
	})

	It("prepends strings to it", func() {
		f := func(head string, tail []string) bool {
			result := PrependString(head, tail)
			return len(result) == 1+len(tail) && result[0] == head && SliceEqual(result[1:], tail)
		}
		if err := quick.Check(f, nil); err != nil {
			t.Error(err)
		}
	})

})
