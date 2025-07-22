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
			isDel, isDelFig := isDelimiter(word)
			if !delimiterAdded && i == 0 && !isDel { // implicit delimiter
				lexs = append(lexs, &Lexer{DelimiterType, ImplicitDelimiter})
				delimiterAdded = false
			}
			if slices.Contains(keywords, word) {
				lexs = append(lexs, &Lexer{KeywordType, word})
			} else if isDel {
				if isDelFig {
					lexs = append(lexs, &Lexer{DelimiterType, FigureDelimiter})
				} else {
					lexs = append(lexs, &Lexer{DelimiterType, word})
				}
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

func isDelimiter(word string) (bool, bool) {
	if slices.Contains(delimiters, word) {
		return true, false
	}
	if len(word) >= 3 && word[:3] == "---" && strings.Count(word, "-") == len(word) {
		return true, true
	}
	return false, false
}
