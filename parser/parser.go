package parser

import (
	"errors"
	"fmt"
	"github.com/planklang/goplank/lexer"
	"github.com/planklang/goplank/parser/types"
	"strconv"
)

var (
	ErrUnknownValue      = errors.New("unknown value")
	ErrInternal          = errors.New("internal error")
	ErrInvalidModifier   = errors.New("invalid modifier")
	ErrInvalidArgument   = errors.New("invalid argument")
	ErrUnexpectedToken   = errors.New("unexpected token")
	ErrDelimiterExcepted = errors.Join(ErrUnexpectedToken, errors.New("delimiter excepted"))
)

func Parse(lex *lexer.TokenList) (*Ast, error) {
	// top-level = [ figure, [{ figure-delimiter, [figure] }] ];

	tree := new(Ast)
	tree.Type = AstTypeDefault

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
			return nil, errors.Join(ErrDelimiterExcepted, fmt.Errorf("expected figure delimiter, not %s", lex.Current()))
		}
	}

	return tree, nil
}

func parseFigure(lex *lexer.TokenList) (*Figure, error) {
	// figure = [ statement, [{ statement-delimiter, [ statement ] }] ];

	fig := new(Figure)

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
			return nil, errors.Join(ErrDelimiterExcepted, fmt.Errorf("expected statement delimiter, not %s", lex.Current()))
		}
	}

	return fig, nil
}

func parseStatement(lex *lexer.TokenList) (*Statement, error) {
	// statement = keyword, [ arguments ], [{ property-delimiter, property }];

	if lex.Current().Type != lexer.KeywordType {
		return nil, errors.Join(ErrUnexpectedToken, fmt.Errorf("expected keyword, not %s", lex.Current()))
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

	for lex.Current().Type == lexer.ModifierDelimiterType {
		if !lex.Next() {
			return nil, errors.Join(lexer.ErrInvalidExpression, fmt.Errorf("expected modifier definition after modifier delimiter"))
		}

		mod, err := parseProperty(lex)
		if err != nil {
			return nil, err
		}
		stmt.Modifiers = append(stmt.Modifiers, mod)
	}

	return stmt, nil
}

func parseProperty(lex *lexer.TokenList) (*Modifier, error) {
	// property = ? identifier ?, [ arguments ]

	if lex.Current().Type != lexer.IdentifierType {
		return nil, errors.Join(ErrUnexpectedToken, fmt.Errorf("expected modifier name, not %v", lex.Current()))
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
	tuple := new(types.Tuple)

	for lex.Current().Type != lexer.StatementDelimiterType &&
		lex.Current().Type != lexer.ModifierDelimiterType &&
		lex.Current().Type != lexer.FigureDelimiterType { // do not call [TokenList.Next] because tuple does not require anything
		lit, err := parseLiteral(lex.Current())
		if err != nil {
			return nil, err
		}
		tuple.AddValues(lit)

		if !lex.Next() { // call [TokenList.Next] here because parseLiteral never call [TokenList.Next] for simple literal
			return tuple, nil
		}
	}

	return tuple, nil
}

func parseLiteral(lex *lexer.Lexer) (types.Value, error) {
	switch lex.Type {
	case lexer.IdentifierType:
		return types.NewDefaultLiteral(lex.Literal), nil
	case lexer.VariableType:
		//TODO: handle
	case lexer.StringType:
		return types.String(lex.Literal), nil
	case lexer.IntType:
		i, err := strconv.ParseInt(lex.Literal, 10, 64)
		if err != nil {
			return nil, err
		}
		return types.Int(i), nil
	case lexer.FloatType:
		f, err := strconv.ParseFloat(lex.Literal, 64)
		if err != nil {
			return nil, err
		}
		return types.Float(f), nil
	case lexer.WeakDelimiterType:
		//TODO: handle
	}
	return nil, errors.Join(ErrUnknownValue, fmt.Errorf("unsupported lex types %s", lex.Type))
}
