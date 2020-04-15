package memo

import (
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

func freeCols(s scalar.Expr) opt.ColSet {
	var cols opt.ColSet
	computeFreeCols(cols, s)
	return cols
}

// func extractBoundConditions(r *RelExpr, e []ScalarExpr) (bool, []ScalarExpr, []ScalarExpr) {
// 	// buh
// 	cols := freeCols(e)
// 	return cols.SubsetOf(r.Props.OutputCols)
// }

func eqConst(s scalar.Expr, d lang.Datum) bool {
	if c, ok := s.(*scalar.Constant); ok {
		if c.D.Type() != d.Type() {
			return false
		}
		return lang.Compare(c.D, d) == lang.EQ
	}
	return false
}
