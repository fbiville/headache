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

package core

import (
	"fmt"
	tpl "html/template"
	"regexp"
	"sort"
	"strings"
)

func ComputeHeaderDetectionRegex(lines []string, data map[string]string) (string, error) {
	return injectDataRegex(strings.Join(computeRegex(lines), ""), data)
}

func computeRegex(lines []string) []string {
	styles := extractValues(SupportedStyles())
	emptyCommentedLine := func(style CommentStyle) string {
		return style.GetString()
	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf(`(?im)\n*(?:%s\n)?\n*`, combineRegexes(styles,
		func(style CommentStyle) string {
			return style.GetOpeningString()
		})))
	for _, line := range lines {
		result = append(result, fmt.Sprintf(`\n*(?:(?:%s) ?\n)*\n*`, combineRegexes(styles, emptyCommentedLine)))
		result = append(result, fmt.Sprintf(`\n*(?:%s)[ \t]*\Q%s\E[ \t\.]*\n*`, combineRegexes(styles,
			func(style CommentStyle) string {
				return style.GetString()
			}),
			line))
	}
	result = append(result, fmt.Sprintf(`\n*(?:(?:%s) ?\n)*\n*`, combineRegexes(styles, emptyCommentedLine)))
	result = append(result, fmt.Sprintf(`\n*(?:%s)?\n*`, combineRegexes(styles,
		func(style CommentStyle) string {
			return style.GetClosingString()
		})))
	return result
}

func combineRegexes(styles []CommentStyle, getLine func(CommentStyle) string) string {
	regexes := make([]string, 0)
	for _, style := range styles {
		if line := getLine(style); line != "" {
			regexes = append(regexes, escape(line))
		}
	}
	return strings.Join(regexes, "|")
}

func escape(str string) string {
	return strings.TrimRight(strings.Replace(regexp.QuoteMeta(str), "/", `\/`, -1), " ")
}

func injectDataRegex(result string, data map[string]string) (string, error) {
	template, err := tpl.New("header-regex").Parse(result)
	if err != nil {
		return "", err
	}
	builder := &strings.Builder{}
	err = template.Execute(builder, regexValues(&data))
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}

func regexValues(data *map[string]string) *map[string]string {
	for k := range *data {
		(*data)[k] = "\\E.*\\Q"
	}
	return data
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
