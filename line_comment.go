package header

import (
	"fmt"
)

type CommentStyle int

const (
	SlashStar CommentStyle = iota
	SlashSlash
)

func (style *CommentStyle) opening() bool {
	return *style == SlashStar
}

func (style *CommentStyle) closing() bool {
	return *style == SlashStar
}

func (style *CommentStyle) open() string {
	switch *style {
	case SlashStar:
		return "/* "
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
	return fmt.Sprintf("%s%s", style.commentSymbol(), line)
}

func (style *CommentStyle) commentSymbol() string {
	switch *style {
	case SlashSlash:
		return "// "
	case SlashStar:
		return " * "
	default:
		panic(unknownCommentError())
	}
}

func unknownCommentError() string {
	return "Unexpected comment style, must be one of: SlashSlash, SlashStar"
}

func unknownSurroundingCommentError() string {
	return "Unexpected surrounding comment style, must be one of: SlashStar"
}
