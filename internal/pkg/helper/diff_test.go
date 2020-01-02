/*
 * Copyright 2018-2019 Florent Biville (@fbiville)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package helper_test

import (
	. "github.com/fbiville/headache/internal/pkg/helper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diff", func() {

	It("just works", func() {
		result, err := Diff(`"foo\nbar\nbaz"`, `"foo\nfighters\nbaz"`)

		Expect(err).To(BeNil())
		Expect(result).To(Equal(`2c2
< bar
---
> fighters
`))
	})

	It("produces an empty string if the contents are the same", func() {
		result, err := Diff("foo", "foo")

		Expect(err).To(BeNil())
		Expect(result).To(Equal(""))
	})
})
