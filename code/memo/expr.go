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
}

type Cross struct {
	Left  RelExpr
	Right RelExpr
}

type Project struct {
	Input RelExpr

	Columns opt.ColSet
}

type Select struct {
	Input  RelExpr
	Filter ScalarExpr

	Columns opt.ColSet
}
