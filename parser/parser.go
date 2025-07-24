package parser

import (
	"errors"
	"fmt"
	"github.com/planklang/goplank/lexer"
	"github.com/planklang/goplank/parser/types"
	"strconv"
)

func Parse(lex lexer.TokenList) (*Ast, error) {
	// top-level = [ figure, [{ figure-delimiter, [figure] }] ];

	tree := new(Ast)
	tree.Type = AstTypeDefault

	var fig *Figure
	var err error

	if lex.Empty() {
		return tree, nil
	}

	fig, lex, err = parseFigure(lex)
	if err != nil {
		return nil, err
	}
	tree.Body = []*Figure{fig}

	for !lex.Empty() {
		if lex.Current().Type != lexer.FigureDelimiterType {
			return nil, errors.Join(ErrUnexpectedToken, fmt.Errorf("expected figure delimiter, got %v", lex.Current()))
		}

		_, ok := lex.Next()
		if !ok {
			return tree, nil
		}

		for lex.Current().Type == lexer.FigureDelimiterType {
			_, ok := lex.Next()
			if !ok {
				return tree, nil
			}
		}

		fig, lex, err = parseFigure(lex)
		if err != nil {
			return nil, err
		}
		tree.Body = append(tree.Body, fig)
	}

	return tree, nil
}

func parseFigure(lex lexer.TokenList) (*Figure, lexer.TokenList, error) {
	// figure = [ statement, [{ statement-delimiter, [ statement ] }] ];

	fig := new(Figure)

	if lex.Empty() {
		return fig, lex, nil
	}

	var stmt *Statement
	var err error

	stmt, lex, err = parseStatement(lex)
	if err != nil {
		return nil, lex, err
	}
	fig.Stmts = []*Statement{stmt}

	for !lex.Empty() {
		if lex.Current().Type != lexer.StatementDelimiterType {
			return nil, lex, errors.Join(ErrUnexpectedToken, fmt.Errorf("expected statement delimiter, got %v", lex.Current()))
		}

		_, ok := lex.Next()
		if !ok {
			return fig, lex, nil
		}

		for lex.Current().Type == lexer.StatementDelimiterType {
			_, ok := lex.Next()
			if !ok {
				return fig, lex, nil
			}
		}

		stmt, lex, err = parseStatement(lex)
		if err != nil {
			return fig, lex, err
		}
		fig.Stmts = append(fig.Stmts, stmt)
	}

	return fig, lex, nil
}

func parseStatement(lex lexer.TokenList) (*Statement, lexer.TokenList, error) {
	// statement = keyword, [ arguments ], [{ property-delimiter, property }];

	if lex.Current().Type != lexer.KeywordType {
		return nil, lex, errors.Join(ErrUnexpectedToken, fmt.Errorf("expected keyword, got %v", lex.Current()))
	}

	return new(Statement), lex, nil
}

func parseTuple(lex lexer.TokenList) (*types.Tuple, lexer.TokenList, error) {
	return new(types.Tuple), lex, nil
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
