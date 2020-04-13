package memo

import "github.com/justinj/scribe/code/opt"

// This file implements the various relational logical operators.

type Scan struct {
	TableName string
	Cols      []opt.ColumnID
}

func (s *Scan) ChildCount() int {
	return 0
}

func (s *Scan) Child(i int) Expr {
	panic("no children")
}

type Join struct {
	Left  RelExpr
	Right RelExpr
	On    ScalarExpr
}

func (j *Join) ChildCount() int {
	return 3
}

func (j *Join) Child(i int) Expr {
	switch i {
	case 0:
		return j.Left
	case 1:
		return j.Right
	case 2:
		return j.On
	default:
		panic("out of bound")
	}
}

type Project struct {
	Input RelExpr

	ColIDs      []opt.ColumnID
	Projections []ScalarExpr
}

func (p *Project) ChildCount() int {
	return 1 + len(p.Projections)
}

func (p *Project) Child(i int) Expr {
	switch i {
	case 0:
		return p.Input
	default:
		return p.Projections[i-1]
	}
}

type Select struct {
	Input RelExpr
	// TODO: unify terminology here: is it filter or predicate?
	Filter ScalarExpr
}

func (s *Select) ChildCount() int {
	return 2
}

func (s *Select) Child(i int) Expr {
	switch i {
	case 0:
		return s.Input
	case 1:
		return s.Filter
	default:
		panic("out of bounds")
	}
}
