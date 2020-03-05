package diff

import (
	d "github.com/andreyvit/diff"
	"strings"
)

func Diff(s1, s2 string) string {
	result := strings.Builder{}
	lines := d.LineDiffAsLines(s1, s2)
	for _, line := range lines {
		if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") {
			result.WriteString(line + "\n")
		}
	}
	return result.String()
}
