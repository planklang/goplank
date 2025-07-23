package lexer

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type LexType string

const (
	KeywordType            LexType = "keyword"
	LiteralType            LexType = "literal"
	WeakDelimiterType      LexType = "weak_delimiter"
	PropertyDelimiterType  LexType = "property_delimiter"
	FigureDelimiterType    LexType = "figure_delimiter"
	StatementDelimiterType LexType = "statement_delimiter"
	IdentifierType         LexType = "identifier"
	VariableType           LexType = "variable"
	StringType             LexType = "string"
	NumberType             LexType = "number"

	ImplicitDelimiter = "implicit"
	FigureDelimiter   = "---"
)

var (
	keywords            = []string{"plot", "default", "overwrite", "ow", "axis"}
	propertyDelimiters  = []string{"|"}
	statementDelimiters = []string{";;"}
	weakDelimiters      = []string{"(", ")", "[", "]"}

	ErrStatementExcepted = fmt.Errorf("statement excepted")
	ErrInvalidExpression = fmt.Errorf("invalid expression")
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
	inProperty := false
	identifierAdded := false
	for ln, line := range lines {
		i := 0
		words := strings.Fields(line)
		for i < len(words) && words[i][0] != '#' { // skip comments
			word := words[i]
			isDelim, typ := isDelimiter(word)
			if !delimiterAdded && i == 0 && !isDelim { // implicit delimiter
				inStatement = false
				inProperty = false
				identifierAdded = false
				delimiterAdded = false
				lexs = append(lexs, &Lexer{StatementDelimiterType, ImplicitDelimiter})
			}
			if isDelim {
				inStatement = false
				inProperty = false
				identifierAdded = false
				delimiterAdded = true
				if typ == FigureDelimiterType {
					lexs = append(lexs, &Lexer{typ, FigureDelimiter})
				} else {
					lexs = append(lexs, &Lexer{typ, word})
					if typ == PropertyDelimiterType {
						inProperty = true
						inStatement = true
					}
				}
			} else if slices.Contains(keywords, word) {
				lexs = append(lexs, &Lexer{KeywordType, word})
				inStatement = true
			} else if !inStatement {
				fmt.Println(genErrorMessage(ErrStatementExcepted, i, words, ln))
				return nil, ErrStatementExcepted
			} else if !identifierAdded && inProperty {
				lexs = append(lexs, &Lexer{IdentifierType, word})
				identifierAdded = true
			} else {
				ls, err := parseLiteral(&i, words)
				if err != nil {
					fmt.Println(genErrorMessage(err, i, words, ln)) // i-1 because every error here leads to i == len(words)
					return nil, err
				}
				lexs = append(lexs, ls...)
			}
			i++
		}
		delimiterAdded = false
		inProperty = false
		identifierAdded = false
	}
	return lexs, nil
}

func parseLiteral(i *int, words []string) ([]*Lexer, error) {
	word := words[*i]
	if isDigit(word) {
		return []*Lexer{{NumberType, word}}, nil
	}
	f := word[0]
	switch f {
	case '$':
		if len(word) == 1 {
			return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("$ is reserved to call variables"))
		}
		return []*Lexer{{VariableType, word[1:]}}, nil
	case '"', '\'', '`':
		s := ""
		finished := false
		j := 1
		for *i < len(words) && !finished {
			c := words[*i]
			for j < len(c) && !finished {
				s += string(c[j])
				if j < len(c)-1 {
					finished = c[j+1] == f
				}
				j++
			}
			s += " "
			j = 0
			*i++
		}
		if !finished {
			*i--
			return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("string is not finished"))
		}
		return []*Lexer{{StringType, s[:len(s)-1]}}, nil
	}

	var lexs []*Lexer

	var precType LexType
	content := ""
	isDecimal := false
	acceptContent := true

	fnUpdate := func(newType LexType) {
		if precType == newType {
			return
		}
		if precType != "" {
			if precType != WeakDelimiterType {
				acceptContent = false
			}
			lexs = append(lexs, &Lexer{precType, content})
		}
		content = ""
		precType = newType
	}

	for _, c := range []rune(word) {
		if slices.Contains(weakDelimiters, string(c)) {
			fnUpdate(WeakDelimiterType)
		} else if !acceptContent {
			return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("cannot parse %s", word))
		} else if isDigit(string(c)) || (c == '.' && !isDecimal) {
			if !isDecimal {
				isDecimal = c == '.'
			}
			if content == "" {
				content += "0" // turns .5 into 0.5
			}
			fnUpdate(NumberType)
		} else {
			fnUpdate(LiteralType)
		}
		content += string(c)
	}

	return append(lexs, &Lexer{precType, content}), nil
}

func genErrorMessage(err error, i int, words []string, line int) string {
	s := ""
	for j := range i {
		s += words[j] + " "
	}
	l1 := len(s)
	s += words[i]
	for j := range len(words) - i - 1 {
		s += " " + words[j+i+1]
	}
	l2 := len(s)
	title := " Parsing error! "
	displayLine := fmt.Sprintf(" (line %d)", line+1)
	after := ""
	size := l2 - len(title) + len(displayLine)
	if size > 0 {
		for n := range size {
			if n%2 == 0 {
				title = "=" + title
			} else {
				title += "="
			}
		}
	}
	for range len(title) {
		after += "="
	}
	s += "\n"
	if i == len(words)-1 {
		for range l2 - 1 {
			s += "-"
		}
		s += "^"
	} else {
		for range l1 {
			s += "-"
		}
		s += "^"
		for range l2 - l1 - 1 {
			s += "-"
		}
	}
	return fmt.Sprintf("%s\n%s%s\n\n%s\n%s\n\n", title, s, displayLine, err.Error(), after)
}

func isDelimiter(word string) (bool, LexType) {
	if slices.Contains(statementDelimiters, word) {
		return true, StatementDelimiterType
	} else if slices.Contains(propertyDelimiters, word) {
		return true, PropertyDelimiterType
	}
	if len(word) >= 3 && word[:3] == "---" && strings.Count(word, "-") == len(word) {
		return true, FigureDelimiterType
	}
	return false, ""
}

func isDigit(word string) bool {
	_, err := strconv.ParseFloat(word, 64)
	return err == nil
}
