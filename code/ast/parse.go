package ast

import (
	"fmt"
	"strconv"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/sexp"
)

func ParseExpr(s string) (Expr, error) {
	t, err := sexp.Parse(s)
	if err != nil {
		return nil, err
	}
	return parseExpr(t)
}

func parseExpr(s sexp.Sexp) (Expr, error) {
	switch e := s.(type) {
	case sexp.Atom:
		switch e {
		case "true":
			return lang.DBool(true), nil
		case "false":
			return lang.DBool(false), nil
		}

		if '0' <= e[0] && e[0] <= '9' {
			n, err := strconv.Atoi(string(e))
			if err != nil {
				return nil, err
			}
			return lang.DInt(n), nil
		}
		return ColumnReference(e), nil
	case sexp.String:
		return lang.DString(e), nil
	case sexp.List:
		if len(e) == 0 {
			return nil, fmt.Errorf("empty list is not an expression")
		}
		head, ok := e[0].(sexp.Atom)
		if !ok {
			return nil, fmt.Errorf("head must be function name")
		}
		op, err := lang.GetFunc(string(head))
		if err != nil {
			return nil, err
		}
		args := make([]Expr, len(e)-1)
		for i, arg := range e[1:] {
			parsed, err := parseExpr(arg)
			if err != nil {
				return nil, err
			}
			args[i] = parsed
		}
		return &ScalarFunc{op, args}, nil
	}
	panic(fmt.Sprintf("unexpected type: %T", s))
}

func parseExprList(s sexp.Sexp) ([]Expr, error) {
	l, ok := s.(sexp.List)
	if !ok {
		return nil, fmt.Errorf("expected expression list")
	}
	result := make([]Expr, len(l))
	for i := range l {
		next, err := parseExpr(l[i])
		if err != nil {
			return nil, err
		}
		result[i] = next
	}
	return result, nil
}

func ParseRelExpr(s string) (RelExpr, error) {
	t, err := sexp.Parse(s)
	if err != nil {
		return nil, err
	}
	return parseRelExpr(t)
}

func parseRelExpr(s sexp.Sexp) (RelExpr, error) {
	switch e := s.(type) {
	case sexp.Atom:
		// An atom in a query position is a table reference.
		return &TableRef{string(e)}, nil
	case sexp.String:
		return nil, fmt.Errorf("expected relational expression, found string (%q)", e)
	case sexp.List:
		if len(e) == 0 {
			return nil, fmt.Errorf("empty list is not a relational expression")
		}
		head, ok := e[0].(sexp.Atom)
		if !ok {
			return nil, fmt.Errorf("head of operator must be atom")
		}
		switch head {
		case "join":
			return parseJoin(e)
		case "select":
			return parseSelect(e)
		case "project":
			return parseProject(e)
		case "as":
			return parseAs(e)
		default:
			return nil, fmt.Errorf("unrecognized relational operator %s", e[0])
		}
	default:
		panic(fmt.Sprintf("unhandled type %T", s))
	}
}

func parseJoin(l sexp.List) (RelExpr, error) {
	if len(l) != 4 {
		return nil, fmt.Errorf("join takes 3 arguments")
	}
	left, err := parseRelExpr(l[1])
	if err != nil {
		return nil, err
	}
	right, err := parseRelExpr(l[2])
	if err != nil {
		return nil, err
	}
	on, err := parseExpr(l[3])
	if err != nil {
		return nil, err
	}
	return &Join{left, right, on}, nil
}

func parseSelect(l sexp.List) (RelExpr, error) {
	if len(l) != 3 {
		return nil, fmt.Errorf("select takes 2 arguments")
	}
	in, err := parseRelExpr(l[1])
	if err != nil {
		return nil, err
	}
	pred, err := parseExpr(l[2])
	if err != nil {
		return nil, err
	}
	return &Select{in, pred}, nil
}

func parseProject(l sexp.List) (RelExpr, error) {
	if len(l) != 3 {
		return nil, fmt.Errorf("project takes 2 arguments")
	}
	in, err := parseRelExpr(l[1])
	if err != nil {
		return nil, err
	}

	l, ok := l[2].(sexp.List)
	if !ok {
		return nil, fmt.Errorf("expected expression list")
	}
	projs := make([]Expr, len(l))
	aliases := make([]string, len(l))
	for i := range l {
		// Special case: bare column references pass through their name.
		if c, ok := l[i].(sexp.Atom); ok {
			aliases[i] = string(c)
		}

		next, err := parseExpr(l[i])
		if err != nil {
			return nil, err
		}
		projs[i] = next
	}

	return &Project{
		in,
		projs,
		aliases,
	}, nil
}

func parseAtomList(s sexp.Sexp) ([]string, error) {
	l, ok := s.(sexp.List)
	if !ok {
		return nil, fmt.Errorf("expected atom list, got %s", s)
	}

	result := make([]string, len(l))
	for i, a := range l {
		n, ok := a.(sexp.Atom)
		if !ok {
			return nil, fmt.Errorf("expected atom list, got %s", s)
		}
		result[i] = string(n)
	}
	return result, nil
}

func parseAs(l sexp.List) (RelExpr, error) {
	if len(l) != 3 && len(l) != 4 {
		return nil, fmt.Errorf("as takes 2 or 3 arguments")
	}
	in, err := parseRelExpr(l[1])
	if err != nil {
		return nil, err
	}
	name, ok := l[2].(sexp.Atom)
	if !ok {
		return nil, fmt.Errorf("as name must be atom")
	}
	var names []string
	if len(l) == 4 {
		names, err = parseAtomList(l[3])
		if err != nil {
			return nil, err
		}
	}
	return &As{in, string(name), names}, nil
}

func ParseStatement(s string) (Statement, error) {
	t, err := sexp.Parse(s)
	if err != nil {
		return nil, err
	}
	return parseStatement(t)
}

func parseType(s sexp.Sexp) (lang.Type, error) {
	n, ok := s.(sexp.Atom)
	if !ok {
		return 0, fmt.Errorf("expected atom, got %s", s)
	}
	switch string(n) {
	case "int":
		return lang.Int, nil
	case "string":
		return lang.String, nil
	case "bool":
		return lang.Bool, nil
	default:
		return 0, fmt.Errorf("invalid type %q", n)
	}
}

func parseColumnDef(s sexp.Sexp) (lang.Column, error) {
	l, ok := s.(sexp.List)
	if !ok || len(l) != 2 {
		return lang.Column{}, fmt.Errorf("expected column def, got %s", s)
	}
	name, ok := l[0].(sexp.Atom)
	if !ok {
		return lang.Column{}, fmt.Errorf("expected column name, got %s", l[0])
	}
	typ, err := parseType(l[1])
	if err != nil {
		return lang.Column{}, err
	}
	return lang.Column{string(name), typ}, nil
}

func parseDatum(s sexp.Sexp) (lang.Datum, error) {
	str, ok := s.(sexp.String)
	if ok {
		return lang.DString(string(str)), nil
	}
	a, ok := s.(sexp.Atom)
	if !ok {
		return nil, fmt.Errorf("expected datum, got %s", a)
	}
	if a == "true" {
		return lang.DBool(true), nil
	}
	if a == "false" {
		return lang.DBool(false), nil
	}
	d, err := strconv.Atoi(string(a))
	if err != nil {
		return nil, err
	}
	return lang.DInt(d), nil
}

func parseRow(s sexp.Sexp) (lang.Row, error) {
	r, ok := s.(sexp.List)
	if !ok {
		return nil, fmt.Errorf("expected error, got %s", s)
	}
	result := make(lang.Row, len(r))
	for i, d := range r {
		parsed, err := parseDatum(d)
		if err != nil {
			return nil, err
		}
		result[i] = parsed
	}
	return result, nil
}

func parseStatement(s sexp.Sexp) (Statement, error) {
	l, ok := s.(sexp.List)
	if !ok {
		return nil, fmt.Errorf("can't parse %s as a statement", s)
	}
	if len(l) == 0 {
		return nil, fmt.Errorf("%s isn't a statement", s)
	}
	h, ok := l[0].(sexp.Atom)
	if !ok {
		return nil, fmt.Errorf("head must be atom, was %s", s)
	}
	switch string(h) {
	case "run":
		if len(l) != 2 {
			return nil, fmt.Errorf("query takes one argument")
		}
		input, err := parseRelExpr(l[1])
		if err != nil {
			return nil, err
		}
		return &RunQuery{input}, nil
	case "create-table":
		if len(l) < 3 || len(l) > 4 {
			return nil, fmt.Errorf("create-table takes 3 or 4 arguments")
		}

		name, ok := l[1].(sexp.Atom)
		if !ok {
			return nil, fmt.Errorf("table name must be atom, was %s", l[1])
		}

		cols, ok := l[2].(sexp.List)
		if !ok {
			return nil, fmt.Errorf("expected list of cols, got %s", l[2])
		}

		defs := make([]lang.Column, len(cols))
		for i, c := range cols {
			def, err := parseColumnDef(c)
			if err != nil {
				return nil, err
			}
			defs[i] = def
		}

		var rows []lang.Row
		if len(l) > 3 {
			rowList, ok := l[3].(sexp.List)
			if !ok {
				return nil, fmt.Errorf("expected list of rows, got %s", l[3])
			}

			rows = make([]lang.Row, len(rowList))
			for i, r := range rowList {
				row, err := parseRow(r)
				if err != nil {
					return nil, err
				}
				rows[i] = row
			}
		}

		return &CreateTable{
			Name:    string(name),
			Columns: defs,
			Data:    rows,
		}, nil
	default:
		return nil, fmt.Errorf("unknown statement %s", h)
	}
}
