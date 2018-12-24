/*
 * Copyright 2018 Florent Biville (@fbiville)
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
	. "github.com/onsi/gomega"
	"testing"
)

func TestDiff(t *testing.T) {
	I := NewGomegaWithT(t)

	result, err := Diff(`"foo\nbar\nbaz"`, `"foo\nfighters\nbaz"`)

	I.Expect(err).To(BeNil())
	I.Expect(result).To(Equal(`2c2
< bar
---
> fighters
`))
}

func TestDiffSameString(t *testing.T) {
	I := NewGomegaWithT(t)

	result, err := Diff(`foo`, `foo`)

	I.Expect(err).To(BeNil())
	I.Expect(result).To(Equal(""))
}
