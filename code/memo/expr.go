package memo

import "github.com/justinj/scribe/code/lang"

type relExpr interface {
	lang.Expr
}

// TODO: this is actually a group, but not sure how to name them appropriately.
// Also not sure how to structure this right so that other packages and inspect it sanely.
// TODO/idea: what if every expr had its children directly, but also was a thing in a map to a metadata thing?
type RelGroup struct {
	// The logical expression.
	E  relExpr
	Es []relExpr

	Props Props
}

func (r *RelGroup) Unwrap() relExpr {
	return r.E
}

func (r *RelGroup) MemberCount() int {
	return 1
}

func (r *RelGroup) Member(i int) lang.Expr {
	return r.E
}
