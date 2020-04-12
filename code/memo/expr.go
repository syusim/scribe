package memo

import "github.com/justinj/scribe/code/opt"

type Expr interface {
}

// TODO: this is actually a group, but not sure how to name them appropriately.
// Also not sure how to structure this right so that other packages and inspect it sanely.
type RelExpr struct {
	E relExpr
}

// TODO: ?? this sucks, think about this more
func Wrap(e relExpr) RelExpr {
	return RelExpr{e}
}

type relExpr interface {
	Expr
}

type Scan struct {
	TableName string
	Cols      []opt.ColumnID
}

type Join struct {
	Left  RelExpr
	Right RelExpr
	On    ScalarExpr
}

type Project struct {
	Input RelExpr

	ColIDs      []opt.ColumnID
	Projections []ScalarExpr
}

type Select struct {
	Input RelExpr
	// TODO: unify terminology here: is it filter or predicate?
	Filter ScalarExpr
}
