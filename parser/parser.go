package parser

import (
	"errors"
	"fmt"
	"github.com/planklang/goplank/lexer"
	"github.com/planklang/goplank/parser/types"
	"strconv"
)

func Parse(lex *lexer.TokenList) (*Ast, error) {
	// top-level = [ figure, [{ figure-delimiter, [figure] }] ];

	tree := new(Ast)
	tree.Type = AstTypeDefault

	if lex.Empty() {
		return tree, nil
	}

	for lex.Next() {
		fig, err := parseFigure(lex)
		if err != nil {
			return nil, err
		}
		tree.Body = append(tree.Body, fig)

		if !lex.Next() {
			return tree, nil
		}
		if lex.Current().Type != lexer.FigureDelimiterType {
			return nil, errors.Join(ErrUnexpectedToken, fmt.Errorf("expected figure delimiter, got %s", lex.Current()))
		}
	}

	return tree, nil
}

func parseFigure(lex *lexer.TokenList) (*Figure, error) {
	// figure = [ statement, [{ statement-delimiter, [ statement ] }] ];

	fig := new(Figure)

	if lex.Empty() {
		return fig, nil
	}

	for lex.Next() {
		stmt, err := parseStatement(lex)
		if err != nil {
			return fig, err
		}
		fig.Stmts = append(fig.Stmts, stmt)

		if !lex.Next() {
			return fig, nil
		}
		if lex.Current().Type != lexer.StatementDelimiterType {
			return nil, errors.Join(ErrUnexpectedToken, fmt.Errorf("expected statement delimiter, got %s", lex.Current()))
		}
	}

	return fig, nil
}

func parseStatement(lex *lexer.TokenList) (*Statement, error) {
	// statement = keyword, [ arguments ], [{ property-delimiter, property }];

	if lex.Current().Type != lexer.KeywordType {
		return nil, errors.Join(ErrUnexpectedToken, fmt.Errorf("expected keyword, got %s", lex.Current()))
	}

	stmt := new(Statement)
	stmt.Keyword = lex.Current().Literal

	if !lex.Next() {
		return stmt, nil
	}

	if lex.Current().Type != lexer.ModifierDelimiterType {
		args, err := parseTuple(lex)
		if err != nil {
			return nil, err
		}
		stmt.Arguments = args
	}

	var mods []*Modifier

	for !lex.Empty() {
		if lex.Current().Type != lexer.ModifierDelimiterType {
			return stmt, nil
		}

		if !lex.Next() {
			return nil, errors.Join(lexer.ErrInvalidExpression, fmt.Errorf("expected modifier definition after modifier delimiter"))
		}

		mod, err := parseProperty(lex)
		if err != nil {
			return nil, err
		}
		mods = append(mods, mod)
	}

	stmt.Modifiers = mods
	return stmt, nil
}
func parseProperty(lex *lexer.TokenList) (*Modifier, error) {
	// property = ? identifier ?, [ arguments ]

	if lex.Current().Type != lexer.IdentifierType {
		return nil, errors.Join(ErrUnexpectedToken, fmt.Errorf("expected modifier name, got %v", lex.Current()))
	}

	mod := new(Modifier)
	mod.Name = lex.Current().Literal

	if !lex.Next() {
		return mod, nil
	}

	args, err := parseTuple(lex)
	if err != nil {
		return nil, err
	}
	mod.Arguments = args
	return mod, nil
}

func parseTuple(lex *lexer.TokenList) (*types.Tuple, error) {
	return new(types.Tuple), nil
}

func parseLiteral(lex *lexer.Lexer, tuple *types.Tuple) error {
	switch lex.Type {
	case lexer.IdentifierType:
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
