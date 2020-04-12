package memo

import "github.com/justinj/scribe/code/opt"

type Expr interface {
}

// TODO: this is actually a group, but not sure how to name them appropriately.
type RelExpr struct {
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
