package lexer

import (
	"fmt"
	"slices"
	"strings"
)

type LexType string

const (
	keywordType   LexType = "keyword"
	literalType   LexType = "literal"
	delimiterType LexType = "delimiter"
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
	return fmt.Sprintf("%s(%s)", lex.Literal, lex.Type)
}

func Lex(content string) []*Lexer {
	var lexs []*Lexer
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		i := 0
		words := strings.Fields(line)
		for i < len(words) && words[i][0] != '#' { // skip comments
			if slices.Contains(keywords, words[i]) {
				lexs = append(lexs, &Lexer{keywordType, words[i]})
			} else if slices.Contains(delimiters, words[i]) {
				lexs = append(lexs, &Lexer{delimiterType, words[i]})
			} else {
				lexs = append(lexs, &Lexer{literalType, words[i]})
			}
			i++
		}
	}
	return lexs
}
