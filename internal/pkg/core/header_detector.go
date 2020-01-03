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

package core

import (
	"fmt"
	tpl "html/template"
	"regexp"
	"sort"
	"strings"
)

func ComputeHeaderDetectionRegex(lines []string, data map[string]string) (string, error) {
	unprocessedRegex := strings.Join(computeRegex(lines), "")
	processedRegex, err := injectDataRegex(unprocessedRegex, data)
	return processedRegex, err
}

func computeRegex(lines []string) []string {
	styles := extractValues(SupportedStyleCatalog())
	result := make([]string, 0)
	result = append(result, Flags())
	result = append(result, OpeningLine(styles))
	for _, line := range lines {
		if line == "" {
			continue
		}
		result = append(result, commentedEmptyLine(styles))
		result = append(result, MatchingLine(line, styles))
	}
	result = append(result, commentedEmptyLine(styles))
	result = append(result, ClosingLine(styles))
	return result
}

// visible for testing
func Flags() string {
	return "(?im)"
}

// visible for testing
func OpeningLine(styles []CommentStyle) string {
	openingLine := fmt.Sprintf(`([\t\v\f\r ]*%s[\t\v\f\r ]*\n)?`, combineRegexes(styles, func(style CommentStyle) string { return style.GetOpeningString() }))
	return openingLine
}

// visible for testing
func MatchingLine(line string, styles []CommentStyle) string {
	openingStyleSymbolRegex := combineRegexes(styles, func(style CommentStyle) string { return style.GetString() })
	normalizedLine := normalizePunctuation(line)
	middleLine := fmt.Sprintf(`[\t\v\f\r ]*%s?[\t\v\f\r ]*\Q%s\E[,.;:?!\t\v\f\r ]*\n?`, openingStyleSymbolRegex, normalizedLine)
	builder := strings.Builder{}
	builder.WriteString(middleLine)
	return builder.String()
}

func normalizePunctuation(line string) string {
	ignore := `\E.?\Q`
	normalizedLine := ""
	// we could use a(n only) slightly better heuristic with a regex matching all dots
	// not prefixed by a digit or a couple of { and 0-n spaces
	// but no support for negative lookbehind in Golang regex engine so here we go :(
	for k, v := range line {
		if v == ',' || v == ';' || v == ':' || v == '?' || v == '!' {
			normalizedLine += ignore
		} else if v == '.' && (k == len(line)-1 || line[k:k+2] == ". ") {
			// dots in template expressions must and in numerical expressions should be preserved
			// hence the poor heuristic of the following space or dot as last line's character
			normalizedLine += ignore
		} else {
			normalizedLine += string(v)
		}
	}
	result := strings.NewReplacer(
		`\t\t`, `\t`,
		`\v\v`, `\v`,
		`\f\f`, `\f`,
		`\r\r`, `\r`,
		"  ", " ",
	).Replace(normalizedLine)
	return strings.NewReplacer(
		"\t", `\E\t+\Q`,
		"\v", `\E\v+\Q`,
		"\f", `\E\f+\Q`,
		"\r", `\E\r+\Q`,
		" ", `\E +\Q`,
	).Replace(result)
}

// visible for testing
func ClosingLine(styles []CommentStyle) string {
	closingLine := fmt.Sprintf(`(?:[\t\v\f\r ]*%s[\t\v\f\r ]*)?`, combineRegexes(styles, func(style CommentStyle) string { return style.GetClosingString() }))
	return closingLine
}

func commentedEmptyLine(styles []CommentStyle) string {
	emptyLines := combineRegexes(styles, func(style CommentStyle) string { return style.GetString() })
	return fmt.Sprintf(`(?:%s?\n)*`, emptyLines)
}

func combineRegexes(styles []CommentStyle, getLine func(CommentStyle) string) string {
	regexes := make([]string, 0)
	for _, style := range styles {
		commentSymbol := getLine(style)
		if line := commentSymbol; line != "" {
			regex := escape(line)
			if strings.HasSuffix(commentSymbol, " ") {
				// right spaces may be formatted away
				// make the right space optional
				regex += "?"
			}
			regexes = append(regexes, regex)
		}
	}
	return fmt.Sprintf("(?:%s)", strings.Join(regexes, "|"))
}

func escape(str string) string {
	return strings.Replace(regexp.QuoteMeta(str), "/", `\/`, -1)
}

func injectDataRegex(result string, data map[string]string) (string, error) {
	template, err := tpl.New("header-regex").Parse(result)
	if err != nil {
		return "", err
	}
	builder := &strings.Builder{}
	templateParameters := regexValues(data)
	err = template.Execute(builder, templateParameters)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func regexValues(data map[string]string) map[string]string {
	result := make(map[string]string)
	for k := range data {
		result[k] = `\E.*\Q`
	}
	return result
}

func extractValues(commentStyles map[string]CommentStyle) []CommentStyle {
	result := make([]CommentStyle, len(commentStyles))
	i := 0
	for _, v := range commentStyles {
		result[i] = v
		i++
	}
	sort.SliceStable(result, CommentStyleSorter(result))
	return result
}
