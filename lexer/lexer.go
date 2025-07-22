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
	FigureDelimiter   = "---"
)

var (
	keywords   = []string{"plot", "default", "overwrite", "ow", "axis"}
	delimiters = []string{";;", "|"}

	ErrStatementExcepted = fmt.Errorf("statement excepted")
)

type Lexer struct {
	Type    LexType
	Literal string
}

func (lex *Lexer) String() string {
	return fmt.Sprintf("%s(%s)", lex.Type, lex.Literal)
}

func Lex(content string) ([]*Lexer, error) {
	var lexs []*Lexer
	lines := strings.Split(content, "\n")
	delimiterAdded := true
	inStatement := false
	for _, line := range lines {
		i := 0
		words := strings.Fields(line)
		for i < len(words) && words[i][0] != '#' { // skip comments
			word := words[i]
			isDel, isDelFig := isDelimiter(word)
			if !delimiterAdded && i == 0 && !isDel { // implicit delimiter
				lexs = append(lexs, &Lexer{DelimiterType, ImplicitDelimiter})
				delimiterAdded = false
				inStatement = false
			}
			if slices.Contains(keywords, word) {
				lexs = append(lexs, &Lexer{KeywordType, word})
				inStatement = true
			} else if isDel {
				inStatement = false
				if isDelFig {
					lexs = append(lexs, &Lexer{DelimiterType, FigureDelimiter})
				} else {
					lexs = append(lexs, &Lexer{DelimiterType, word})
					inStatement = word == "|"
				}
				delimiterAdded = true
			} else if !inStatement {
				return nil, ErrStatementExcepted
			} else {
				lexs = append(lexs, &Lexer{LiteralType, word})
			}
			i++
		}
		delimiterAdded = false
	}
	return lexs, nil
}

func isDelimiter(word string) (bool, bool) {
	if slices.Contains(delimiters, word) {
		return true, false
	}
	if len(word) >= 3 && word[:3] == "---" && strings.Count(word, "-") == len(word) {
		return true, true
	}
	return false, false
}
