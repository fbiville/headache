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
