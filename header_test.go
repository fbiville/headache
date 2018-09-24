package header

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
		HeaderContents: `// some header
// with some text`,
		Includes: []string{"fixtures/*_world.txt"},
		writer:   writer,
	}

	Insert(&configuration)
	writer.Flush()

	I.Expect(stringBuilder.String()).To(Equal(`// some header
// with some text
hello
world`))
}

func TestHeaderDoesNotWriteTwice(t *testing.T) {
	I := NewGomegaWithT(t)
	stringBuilder := strings.Builder{}
	writer := bufio.NewWriter(&stringBuilder)
	configuration := configuration{
		HeaderContents: `// some header
// with some text`,
		Includes: []string{"fixtures/*_world_with_header.txt"},
		writer:   writer,
	}

	Insert(&configuration)
	writer.Flush()

	I.Expect(stringBuilder.String()).To(Equal(``))
}
