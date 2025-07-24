package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/planklang/goplank/lexer"
	"github.com/planklang/goplank/parser/types"
	"strconv"
)

type AstType uint

const (
	AstTypeDefault AstType = 0
)

var (
	ErrUnknownValue    = errors.New("unknown value")
	ErrInternal        = errors.New("internal error")
	ErrInvalidModifier = errors.New("invalid modifier")
	ErrInvalidArgument = errors.New("invalid argument")
)

type Ast struct {
	Type AstType
	Body []Statement
}

func (a *Ast) Eval() error {
	for _, s := range a.Body {
		if err := s.Eval(); err != nil {
			return err
		}
	}
	return nil
}

func (a *Ast) String() string {
	m, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return ""
	}
	return string(m)
}

func Parse(lex []*lexer.Lexer) (*Ast, error) {
	tree := new(Ast)
	tree.Type = AstTypeDefault
	var stmt Statement
	var modif Modifier
	inModifier := false
	tuple := new(types.Tuple)
	for _, l := range lex {
		switch l.Type {
		case lexer.StatementDelimiterType:
			if stmt == nil {
				return nil, errors.Join(ErrInternal, fmt.Errorf("statement is nil but statement finished"))
			}
			var err error
			if inModifier {
				if modif == nil {
					return nil, errors.Join(ErrInternal, fmt.Errorf("modif is nil but statement finished"))
				}
				if err = modif.SetArgument(tuple); err != nil {
					return nil, err
				}
				err = stmt.AddModifier(modif)
			} else {
				err = stmt.SetArgument(tuple)
			}
			if err != nil {
				return nil, err
			}
			tree.Body = append(tree.Body, stmt)
		case lexer.KeywordType:
			switch l.Literal {
			case "axis":
				stmt = new(Axis)
			default:
				return nil, errors.Join(ErrUnknownValue, fmt.Errorf("unsupported keyword %s", l.Literal))
			}
		case lexer.ModifierDelimiterType:
			if stmt == nil {
				return nil, errors.Join(ErrInternal, fmt.Errorf("statement is nil but modifier started"))
			}
			if inModifier {
				if modif == nil {
					return nil, errors.Join(ErrInternal, fmt.Errorf("modif is nil but modifier finished"))
				}
				if err := modif.SetArgument(tuple); err != nil {
					return nil, err
				}
				if err := stmt.AddModifier(modif); err != nil {
					return nil, err
				}
			} else {
				if err := stmt.SetArgument(tuple); err != nil {
					return nil, err
				}
			}
			inModifier = true
			modif = nil
		case lexer.ModifierType:
			if inModifier && modif == nil {
				switch l.Literal {
				case "color":
					modif = &ModifierColor{}
				default:
					return nil, errors.Join(ErrUnknownValue, fmt.Errorf("unsupported identifier %s", l.Literal))
				}
			} else {
				return nil, errors.Join(ErrInternal, fmt.Errorf("identifier received but not in modifier or modif is not nil"))
			}
		default:
			err := parseLiteral(l, tuple)
			if err != nil {
				return nil, err
			}
		}
	}
	return tree, nil
}

func parseLiteral(lex *lexer.Lexer, tuple *types.Tuple) error {
	switch lex.Type {
	case lexer.LiteralType:
		tuple.AddValues(types.NewDefaultLiteral(lex.Literal))
	case lexer.VariableType:
		//TODO: handle
	case lexer.StringType:
		tuple.AddValues(types.String(lex.Literal))
	case lexer.IntType:
		i, err := strconv.ParseInt(lex.Literal, 10, 64)
		if err != nil {
			return err
		}
		tuple.AddValues(types.Int(i))
	case lexer.FloatType:
		f, err := strconv.ParseFloat(lex.Literal, 64)
		if err != nil {
			return err
		}
		tuple.AddValues(types.Float(f))
	case lexer.WeakDelimiterType:
		//TODO: handle
	default:
		return errors.Join(ErrUnknownValue, fmt.Errorf("unsupported lex types %s", lex.Type))
	}
	return nil
}
