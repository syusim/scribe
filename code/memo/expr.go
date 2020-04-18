package memo

import "github.com/justinj/scribe/code/lang"

type relExpr interface {
	lang.Expr
}

// TODO: this is actually a group, but not sure how to name them appropriately.
// Also not sure how to structure this right so that other packages and inspect it sanely.
// TODO/idea: what if every expr had its children directly, but also was a thing in a map to a metadata thing?
type RelExpr struct {
	// The logical expression.
	E relExpr

	// Physical implementations

	Props Props
}

// TODO: these types seem extremely wonky.
// should Child/ChildCount just be methods on
// RelExpr which then defer to a big switch for
// the op? regardless it seems they shouldn't
// BOTH be like this.
func (r RelExpr) ChildCount() int {
	return r.E.ChildCount()
}

func (r RelExpr) Child(i int) lang.Expr {
	return r.E.Child(i)
}
