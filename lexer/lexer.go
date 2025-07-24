package lexer

import (
	"errors"
	"fmt"
	"github.com/planklang/goplank/errorshelper"
	"slices"
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
	IntType                LexType = "int"
	FloatType              LexType = "float"

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
		parenthesisCounter := 0
		squareBracketsCounter := 0
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
				ls, err := parseLiteral(&i, words, &parenthesisCounter, &squareBracketsCounter)
				if err != nil {
					fmt.Println(genErrorMessage(err, i, words, ln)) // i-1 because every error here leads to i == len(words)
					return nil, err
				}
				lexs = append(lexs, ls...)
			}
			i++
		}
		if parenthesisCounter != 0 {
			err := errors.Join(ErrInvalidExpression, fmt.Errorf("missing )"))
			fmt.Println(genErrorMessage(err, i-1, words, ln))
			return nil, err
		}
		if squareBracketsCounter != 0 {
			err := errors.Join(ErrInvalidExpression, fmt.Errorf("missing ]"))
			fmt.Println(genErrorMessage(err, i-1, words, ln))
			return nil, err
		}
		delimiterAdded = false
		inProperty = false
		identifierAdded = false
	}
	return lexs, nil
}

func parseLiteral(i *int, words []string, parenthesisCounter *int, squareBracketsCounter *int) ([]*Lexer, error) {
	word := words[*i]
	f := word[0]
	if ok, dec := isDigit(word); ok {
		if f == '.' {
			word = "0" + word
		}
		var typ LexType
		if dec {
			typ = FloatType
		} else {
			typ = IntType
		}
		return []*Lexer{{typ, word}}, nil
	}
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

	for _, c := range word {
		if slices.Contains(weakDelimiters, string(c)) {
			switch c {
			case '(':
				*parenthesisCounter++
			case ')':
				if *parenthesisCounter == 0 {
					return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("missing ("))
				}
				*parenthesisCounter--
			case '[':
				*squareBracketsCounter++
			case ']':
				if *squareBracketsCounter == 0 {
					return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("missing ["))
				}
				*squareBracketsCounter--
			}
			fnUpdate(WeakDelimiterType)
		} else if !acceptContent {
			return nil, errors.Join(ErrInvalidExpression, fmt.Errorf("cannot parse %s", word))
		} else if ok, dec := isDigit(string(c)); ok && (!dec || !isDecimal) {
			if !isDecimal {
				isDecimal = dec
			}
			if content == "" {
				content += "0" // turns .5 into 0.5
			}
			if isDecimal {
				if precType == IntType {
					precType = FloatType
				}
				fnUpdate(FloatType)
			} else {
				fnUpdate(IntType)
			}
		} else {
			fnUpdate(LiteralType)
		}
		content += string(c)
	}

	return append(lexs, &Lexer{precType, content}), nil
}

func genErrorMessage(err error, i int, words []string, line int) string {
	return errorshelper.GenErrorMessage("Parsing error!", err, i, words, line)
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

func isDigit(word string) (bool, bool) {
	isDecimal := false
	for _, c := range word {
		if c == '.' {
			if !isDecimal {
				isDecimal = true
			} else {
				return false, false
			}
		} else if int(c) < int('0') || int(c) > int('9') {
			return false, false
		}
	}
	return true, isDecimal
}
