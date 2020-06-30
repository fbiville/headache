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

package core

import (
	"log"
	"strings"

	styles "github.com/fbiville/headache/internal/pkg/core/comment_styles"
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

func SupportedStyleCatalog() map[string]CommentStyle {
	commentStyles := SupportedStyles()
	result := make(map[string]CommentStyle, len(commentStyles))
	for _, style := range commentStyles {
		result[style.GetName()] = style
	}
	return result
}

func SupportedStyles() []CommentStyle {
	return []CommentStyle{
		styles.SlashStar{},
		styles.SlashSlash{},
		styles.Hash{},
		styles.DashDash{},
		styles.SemiColon{},
		styles.Rem{},
		styles.SlashStarStar{},
		styles.Xml{},
		styles.SingleQuote{},
	}
}

func ParseCommentStyle(name string) CommentStyle {
	commentStyles := SupportedStyleCatalog()
	for styleName, style := range commentStyles {
		if strings.ToLower(styleName) == strings.ToLower(name) {
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
