package main

import (
	"bufio"
	. "github.com/onsi/gomega"
	"regexp"
	"strings"
	"testing"
)

func TestHeaderWrite(t *testing.T) {
	I := NewGomegaWithT(t)
	stringBuilder := strings.Builder{}
	writer := bufio.NewWriter(&stringBuilder)
	newHeader := `// some multi-line header
// with some text`
	regex, _ := computeDetectionRegex([]string{"some multi-line header", "with some text"}, map[string]string{})
	configuration := configuration{
		HeaderRegex:    regexp.MustCompile(regex),
		HeaderContents: newHeader,
		Includes:       []string{"fixtures/hello_world.txt"},
		writer:         writer,
	}

	InsertHeader(&configuration)
	writer.Flush()

	I.Expect(stringBuilder.String()).To(Equal(`// some multi-line header
// with some text

hello
world`))
}

func TestHeaderDoesNotWriteTwice(t *testing.T) {
	I := NewGomegaWithT(t)
	stringBuilder := strings.Builder{}
	writer := bufio.NewWriter(&stringBuilder)
	newHeader := `// some multi-line header
// with some text`
	regex, _ := computeDetectionRegex([]string{"some multi-line header", "with some text"}, map[string]string{})
	configuration := configuration{
		HeaderRegex:    regexp.MustCompile(regex),
		HeaderContents: newHeader,
		Includes:       []string{"fixtures/*_world_with_header.txt"},
		writer:         writer,
	}

	InsertHeader(&configuration)
	writer.Flush()

	I.Expect(stringBuilder.String()).To(Equal(`// some multi-line header
// with some text

hello
world`), "it should rewrite the file as is")
}

func TestHeaderWriteWithExcludes(t *testing.T) {
	I := NewGomegaWithT(t)
	stringBuilder := strings.Builder{}
	writer := bufio.NewWriter(&stringBuilder)
	newHeader := `// some header`
	regex, _ := computeDetectionRegex([]string{"some header"}, map[string]string{})
	configuration := configuration{
		HeaderRegex:    regexp.MustCompile(regex),
		HeaderContents: newHeader,
		Includes:       []string{"fixtures/*_world.txt"},
		Excludes:       []string{"fixtures/hello_*.txt"},
		writer:         writer,
	}

	InsertHeader(&configuration)
	writer.Flush()

	I.Expect(stringBuilder.String()).To(Equal(`// some header

bonjour
le world`))
}

func TestHeaderWithRecursiveGlobs(t *testing.T) {
	I := NewGomegaWithT(t)
	stringBuilder := strings.Builder{}
	writer := bufio.NewWriter(&stringBuilder)
	newHeader := `// some header`
	regex, _ := computeDetectionRegex([]string{"some header"}, map[string]string{})
	configuration := configuration{
		HeaderRegex:    regexp.MustCompile(regex),
		HeaderContents: newHeader,
		Includes:       []string{"fixtures/**/inception.txt"},
		Excludes:       []string{"fixtures/**/ignored.txt"},
		writer:         writer,
	}

	InsertHeader(&configuration)
	writer.Flush()

	I.Expect(stringBuilder.String()).To(Equal(`// some header

a dream
within a dream`))
}

func TestHeaderCommentUpdate(t *testing.T) {
	I := NewGomegaWithT(t)
	stringBuilder := strings.Builder{}
	writer := bufio.NewWriter(&stringBuilder)
	newHeader := `/*
 * some multi-line header
 * with some text
 */`
	regex, _ := computeDetectionRegex([]string{"some multi-line header", "with some text"}, map[string]string{})
	configuration := configuration{
		HeaderRegex:    regexp.MustCompile(regex),
		HeaderContents: newHeader,
		Includes:       []string{"fixtures/*_world_with_header.txt"},
		writer:         writer,
	}

	InsertHeader(&configuration)
	writer.Flush()

	I.Expect(stringBuilder.String()).To(Equal(`/*
 * some multi-line header
 * with some text
 */

hello
world`), "it should rewrite the file with slashstar style")
}

func TestHeaderDataUpdate(t *testing.T) {
	I := NewGomegaWithT(t)
	stringBuilder := strings.Builder{}
	writer := bufio.NewWriter(&stringBuilder)
	newHeader := `// some multi-line header 2018-2020
// with some text from Pairing Corp`
	regex, _ := computeDetectionRegex([]string{"some multi-line header {{.Year}}", "with some text from {{.Company}}"},
		map[string]string{
			"Year":    "2017",
			"Company": "Soloing Inc.",
		})
	configuration := configuration{
		HeaderRegex:    regexp.MustCompile(regex),
		HeaderContents: newHeader,
		Includes:       []string{"fixtures/hello_world_with_parameterized_header.txt"},
		writer:         writer,
	}

	InsertHeader(&configuration)
	writer.Flush()

	s := stringBuilder.String()
	I.Expect(s).To(Equal(`// some multi-line header 2018-2020
// with some text from Pairing Corp

hello
world`), "it should rewrite the file with slashstar style")
}
