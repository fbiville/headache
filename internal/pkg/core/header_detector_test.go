/*
 * Copyright 2019 Florent Biville (@fbiville)
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

package core_test

import (
	"fmt"
	"github.com/fbiville/headache/internal/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"regexp"
	"strings"
)

var _ = Describe("Header detector", func() {

	Context("with a simple header and source file without... source", func() {

		const file = `// some multi-line header
// with some text`

		It("should detect it", func() {
			regex, err := core.ComputeHeaderDetectionRegex(
				[]string{"some multi-line header", "with some text"},
				map[string]string{})

			Expect(err).NotTo(HaveOccurred())
			Expect(matchLeftMostPositions(regex, file)).To(Equal([]int{0, 43}))
		})
	})

	Context("with a simple header and source file", func() {

		const file = `// some multi-line header
// with some text
hello
world`

		It("should detect it", func() {
			regex, err := core.ComputeHeaderDetectionRegex(
				[]string{"some multi-line header", "with some text"},
				map[string]string{})

			Expect(err).NotTo(HaveOccurred())
			Expect(matchLeftMostPositions(regex, file)).To(Equal([]int{0, 44}))
		})
	})

	Context("with whitespace variations", func() {

		const file = `    /*     
    *     some multi-line header       
 * with some text
      */         
hello
world`

		It("should detect it", func() {
			regex, err := core.ComputeHeaderDetectionRegex(
				[]string{"some multi-line header", "with some text"},
				map[string]string{})

			Expect(err).NotTo(HaveOccurred())
			Expect(matchLeftMostPositions(regex, file)).To(Equal([]int{0, 87}))
		})
	})

	Context("with newline variations", func() {

		const file = `
# 
# some multi-line header

# 
# with some text
hello
world`

		It("should detect it", func() {
			regex, err := core.ComputeHeaderDetectionRegex(
				[]string{"some multi-line header", "with some text"},
				map[string]string{})

			Expect(err).NotTo(HaveOccurred())
			Expect(matchLeftMostPositions(regex, file)).To(Equal([]int{0, 50}))
		})
	})

	Context("with newline and whitespace variations", func() {


		const file = `    /**
 * 

    *     some multi-line header     

 * with some text
      */         
hello
world`

		It("should detect it", func() {
			regex, err := core.ComputeHeaderDetectionRegex(
				[]string{"some multi-line header", "with some text"},
				map[string]string{})

			Expect(err).NotTo(HaveOccurred())
			Expect(matchLeftMostPositions(regex, file)).To(Equal([]int{0, 87}))
		})
	})

	Context("with punctuation and whitespace variations in the source file", func() {


		const file = `// some multi-line header.!:
// with some text ,?;
hello
world`

		It("should detect it", func() {
			regex, err := core.ComputeHeaderDetectionRegex(
				[]string{"some multi-line header", "with some text"},
				map[string]string{})

			Expect(err).NotTo(HaveOccurred())
			Expect(matchLeftMostPositions(regex, file)).To(Equal([]int{0, 51}))
		})
	})

	Context("with punctuation and whitespace variations in the license header template", func() {


		const file = `// some multi-line header
// with some text
hello
world`

		It("should detect it", func() {
			regex, err := core.ComputeHeaderDetectionRegex(
				[]string{"some. multi-line.  header!", "with  some text?"},
				map[string]string{})

			Expect(err).NotTo(HaveOccurred())
			Expect(matchLeftMostPositions(regex, file)).To(Equal([]int{0, 44}))
		})
	})

	Context("with whitespaces variations with comment style symbols including whitespaces", func() {


		const file = `/* 
 * 
*some multi-line header
 *with some text
*/
hello
world`

		It("should detect it", func() {
			regex, err := core.ComputeHeaderDetectionRegex(
				[]string{"some multi-line header", "with some text"},
				map[string]string{})

			Expect(err).NotTo(HaveOccurred())
			Expect(matchLeftMostPositions(regex, file)).To(Equal([]int{0, 51}))
		})
	})

	Context("with a realistic header", func() {

		const template = `Copyright {{.YearRange}} {{.Owner}}

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
`

		const file = `/*
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

package main

import "fmt"

func main() {
	fmt.Println("Hello world")
}
`
		styles := core.SupportedStyles()

		It("matches the header opening comment line", func() {
			regex := fmt.Sprintf("%s%s", core.Flags(), core.OpeningLine(styles))

			Expect(matchLeftMostPositions(regex, file)).To(Equal([]int{0, 3}))
		})

		It("matches the header intermediate comment line", func() {
			regex := fmt.Sprintf("%s%s", core.Flags(), core.MatchingLine(
				"    http://www.apache.org/licenses/LICENSE-2.0",
				styles))

			Expect(matchLeftMostPositions(regex, file)).To(Equal([]int{233, 283}))
		})

		It("computes a regex to match existing header", func() {
			templateLines := strings.Split(template, "\n")
			templateParameters := map[string]string{
				"Owner":     "Florent Biville (@fbiville)",
				"YearRange": "{{.YearRange}}",
			}

			regex, err := core.ComputeHeaderDetectionRegex(templateLines, templateParameters)

			Expect(err).NotTo(HaveOccurred())
			Expect(matchLeftMostPositions(regex, file)).To(Equal([]int{0, 610}))
		})
	})
})

func matchLeftMostPositions(regex, contents string) []int {
	return regexp.MustCompile(regex).FindStringIndex(contents)
}
