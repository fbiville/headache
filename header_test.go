package main

import (
	"bufio"
	. "github.com/onsi/gomega"
	"strings"
	"testing"
)

func TestHeaderWrite(t *testing.T) {
	I := NewGomegaWithT(t)
	stringBuilder := strings.Builder{}
	writer := bufio.NewWriter(&stringBuilder)
	configuration := configuration{
		HeaderContents: `// some multi-line header
// with some text`,
		Includes: []string{"fixtures/hello_world.txt"},
		writer:   writer,
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
	configuration := configuration{
		HeaderContents: `// some multi-line header
// with some text`,
		Includes: []string{"fixtures/*_world_with_header.txt"},
		writer:   writer,
	}

	InsertHeader(&configuration)
	writer.Flush()

	I.Expect(stringBuilder.String()).To(Equal(``))
}


func TestHeaderWriteWithExcludes(t *testing.T) {
	I := NewGomegaWithT(t)
	stringBuilder := strings.Builder{}
	writer := bufio.NewWriter(&stringBuilder)
	configuration := configuration{
		HeaderContents: `// some header`,
		Includes: []string{"fixtures/*_world.txt"},
		Excludes: []string{"fixtures/hello_*.txt"},
		writer:   writer,
	}

	InsertHeader(&configuration)
	writer.Flush()

	I.Expect(stringBuilder.String()).To(Equal(`// some header

bonjour
le world`))
}
