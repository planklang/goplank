package lexer

import (
	"fmt"
	"slices"
	"strings"
)

type LexType string

const (
	KeywordType   LexType = "keyword"
	LiteralType   LexType = "literal"
	DelimiterType LexType = "delimiter"

	ImplicitDelimiter = "implicit"
)

var (
	keywords   = []string{"plot", "default", "overwrite", "ow", "axis"}
	delimiters = []string{";;", "|"}
)

type Lexer struct {
	Type    LexType
	Literal string
}

func (lex *Lexer) String() string {
	return fmt.Sprintf("%s(%s)", lex.Type, lex.Literal)
}

func Lex(content string) []*Lexer {
	var lexs []*Lexer
	lines := strings.Split(content, "\n")
	delimiterAdded := true
	for _, line := range lines {
		i := 0
		words := strings.Fields(line)
		for i < len(words) && words[i][0] != '#' { // skip comments
			word := words[i]
			isDelimiter := slices.Contains(delimiters, word)
			if !delimiterAdded && i == 0 && !isDelimiter { // implicit delimiter
				lexs = append(lexs, &Lexer{DelimiterType, ImplicitDelimiter})
				delimiterAdded = false
			}
			if slices.Contains(keywords, word) {
				lexs = append(lexs, &Lexer{KeywordType, word})
			} else if isDelimiter {
				lexs = append(lexs, &Lexer{DelimiterType, word})
				delimiterAdded = true
			} else {
				lexs = append(lexs, &Lexer{LiteralType, word})
			}
			i++
		}
		delimiterAdded = false
	}
	return lexs
}
