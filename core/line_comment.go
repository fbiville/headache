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
)

type CommentStyle int

const (
	SlashStar CommentStyle = iota
	SlashSlash
)

func newCommentStyle(str string) CommentStyle {
	switch str {
	case "SlashStar":
		return SlashStar
	case "SlashSlash":
		return SlashSlash
	default:
		panic(unknownStyleError())
	}
}

func (style *CommentStyle) opening() bool {
	return *style == SlashStar
}

func (style *CommentStyle) closing() bool {
	return *style == SlashStar
}

func (style *CommentStyle) open() string {
	switch *style {
	case SlashStar:
		return "/*"
	default:
		panic(unknownSurroundingCommentError())
	}
}

func (style *CommentStyle) close() string {
	switch *style {
	case SlashStar:
		return " */"
	default:
		panic(unknownSurroundingCommentError())
	}
}

func (style *CommentStyle) apply(line string) string {
	if line == "" {
		return style.commentSymbol()
	}
	return fmt.Sprintf("%s %s", style.commentSymbol(), line)
}

func (style *CommentStyle) commentSymbol() string {
	switch *style {
	case SlashSlash:
		return "//"
	case SlashStar:
		return " *"
	default:
		panic(unknownStyleError())
	}
}

func unknownStyleError() string {
	return "Unexpected comment style, must be one of: SlashSlash, SlashStar"
}

func unknownSurroundingCommentError() string {
	return "Unexpected surrounding comment style, must be one of: SlashStar"
}
