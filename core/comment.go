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
	"log"
	"regexp"
	"strings"
)

type CommentStyle interface {
	GetName() string
	GetOpeningSymbol() *CommentSymbol
	GetContinuationSymbol() *CommentSymbol
	GetClosingSymbol() *CommentSymbol
}

func EmptyCommentSymbol() *CommentSymbol {
	return &CommentSymbol{Value: ""}
}

func NewCommentSymbol(value string) *CommentSymbol {
	return &CommentSymbol{Value: value}
}

func NewOptionalCommentSymbol(value string) *CommentSymbol {
	return &CommentSymbol{Value: value, Optional: true}
}

type CommentSymbol struct {
	Value    string
	Optional bool
}

type SlashStar struct{}

func (SlashStar) GetName() string {
	return "SlashStar"
}
func (SlashStar) GetOpeningSymbol() *CommentSymbol {
	return NewCommentSymbol("/*")
}
func (SlashStar) GetContinuationSymbol() *CommentSymbol {
	return NewOptionalCommentSymbol(" * ")
}
func (SlashStar) GetClosingSymbol() *CommentSymbol {
	return NewCommentSymbol(" */")
}

type SlashSlash struct{}

func (SlashSlash) GetName() string {
	return "SlashSlash"
}
func (SlashSlash) GetOpeningSymbol() *CommentSymbol {
	return EmptyCommentSymbol()
}
func (SlashSlash) GetContinuationSymbol() *CommentSymbol {
	return NewCommentSymbol("// ")
}
func (SlashSlash) GetClosingSymbol() *CommentSymbol {
	return EmptyCommentSymbol()
}

type Hash struct{}

func (Hash) GetName() string {
	return "Hash"
}
func (Hash) GetOpeningSymbol() *CommentSymbol {
	return EmptyCommentSymbol()
}
func (Hash) GetContinuationSymbol() *CommentSymbol {
	return NewCommentSymbol("# ")
}
func (Hash) GetClosingSymbol() *CommentSymbol {
	return EmptyCommentSymbol()
}

func ParseCommentStyle(name string) CommentStyle {
	styles := supportedStyles()
	keys := extractKeys(styles)
	for _, key := range keys {
		if name == key {
			return styles[key]
		}
	}
	log.Fatalf("headache configuration error, unexpected comment style\n\tmust be one of: " + strings.Join(keys, ","))
	return nil
}

func ComputeDetectionRegex(lines []string, data map[string]string) (string, error) {
	regex := computeRegex(lines)
	return injectDataRegex(strings.Join(regex, ""), data)
}

func computeRegex(lines []string) []string {
	styles := extractValues(supportedStyles())
	emptyCommentedLine := func(style CommentStyle) *CommentSymbol {
		return style.GetContinuationSymbol()
	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf(`(?im)\n*(?:%s\n)?\n*`, combineRegexes(styles,
		func(style CommentStyle) *CommentSymbol {
			return style.GetOpeningSymbol()
		})))
	for _, line := range lines {
		result = append(result, fmt.Sprintf(`\n*(?:(?:%s) ?\n)*\n*`, combineRegexes(styles, emptyCommentedLine)))
		result = append(result, fmt.Sprintf(`\n*(?:%s)[ \t]*\Q%s\E[ \t\.]*\n*`, combineRegexes(styles,
			func(style CommentStyle) *CommentSymbol {
				return style.GetContinuationSymbol()
			}),
			line))
	}
	result = append(result, fmt.Sprintf(`\n*(?:(?:%s) ?\n)*\n*`, combineRegexes(styles, emptyCommentedLine)))
	result = append(result, fmt.Sprintf(`\n*(?:%s)?\n*`, combineRegexes(styles,
		func(style CommentStyle) *CommentSymbol {
			return style.GetClosingSymbol()
		})))
	return result
}

func combineRegexes(styles []CommentStyle, getSymbol func(CommentStyle) *CommentSymbol) string {
	regexes := make([]string, 0)
	for _, style := range styles {
		if symbol := getSymbol(style); symbol.Value != "" {
			regexes = append(regexes, escape(symbol))
		}
	}
	return strings.Join(regexes, "|")
}

func escape(symbol *CommentSymbol) string {
	value := symbol.Value
	result := strings.TrimRight(strings.Replace(regexp.QuoteMeta(value), "/", `\/`, -1), " ")
	if symbol.Optional {
		return fmt.Sprintf("(?:%s)?", value)
	}
	return result
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

func supportedStyles() map[string]CommentStyle {
	return map[string]CommentStyle{
		"SlashStar":  SlashStar{},
		"SlashSlash": SlashSlash{},
		"Hash":       Hash{},
	}
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

func extractValues(myMap map[string]CommentStyle) []CommentStyle {
	values := make([]CommentStyle, len(myMap))
	i := 0
	for _, v := range myMap {
		values[i] = v
		i++
	}
	return values
}

type CommentStyleBuilder func() CommentStyle
