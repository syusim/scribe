package memo

import (
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
)

func freeCols(s ScalarExpr) opt.ColSet {
	var cols opt.ColSet
	computeFreeCols(cols, s)
	return cols
}

// func extractBoundConditions(r *RelExpr, e []ScalarExpr) (bool, []ScalarExpr, []ScalarExpr) {
// 	// buh
// 	cols := freeCols(e)
// 	return cols.SubsetOf(r.Props.OutputCols)
// }

func eqConst(s ScalarExpr, d lang.Datum) bool {
	if c, ok := s.(*Constant); ok {
		if c.D.Type() != d.Type() {
			return false
		}
		return lang.Compare(c.D, d) == lang.EQ
	}
	return false
}
