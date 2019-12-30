package core

import (
	styles "github.com/fbiville/headache/core/comment_styles"
	"log"
	"strings"
)

type CommentStyle interface {
	GetName() string
	GetOpeningString() string
	GetString() string
	GetClosingString() string
}

func CommentStyleSorter(styles []CommentStyle) func(i, j int) bool {
	return func(i, j int) bool {
		return styles[i].GetName() < styles[j].GetName()
	}
}

func SupportedStyles() map[string]CommentStyle {
	commentStyles := []CommentStyle{
		styles.SlashStar{},
		styles.SlashSlash{},
		styles.Hash{},
		styles.DashDash{},
		styles.SemiColon{},
		styles.Rem{},
		styles.SlashStarStar{},
	}
	result := make(map[string]CommentStyle, len(commentStyles))
	for _, style := range commentStyles {
		result[style.GetName()] = style
	}
	return result
}

func ParseCommentStyle(name string) CommentStyle {
	commentStyles := SupportedStyles()
	for styleName, style := range commentStyles {
		if styleName == name {
			return style
		}
	}
	log.Fatalf("headache configuration error, unexpected comment style\n\tmust be one of: " +
		strings.Join(extractKeys(commentStyles), ","))
	return nil
}

func ApplyComments(lines []string, style CommentStyle) ([]string, error) {
	result := make([]string, 0)
	if openingLine := style.GetOpeningString(); openingLine != "" {
		result = append(result, openingLine)
	}
	for _, line := range lines {
		result = append(result, prependLine(style, line))
	}
	if closingLine := style.GetClosingString(); closingLine != "" {
		result = append(result, closingLine)
	}
	return result, nil
}

func prependLine(style CommentStyle, line string) string {
	comment := style.GetString()
	if line == "" {
		return strings.TrimRight(comment, " ")
	}
	return comment + line
}

func extractKeys(myMap map[string]CommentStyle) []string {
	keys := make([]string, len(myMap))
	i := 0
	for k := range myMap {
		keys[i] = k
		i++
	}
	return keys
}
