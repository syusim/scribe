package memo

import "github.com/justinj/scribe/code/opt"

// TODO: pull this into a pkg/file
type Props struct {
	OutputCols opt.ColSet
}

func buildProps(r *RelExpr) {
	switch e := r.E.(type) {
	case *Scan:
		for _, c := range e.Cols {
			r.Props.OutputCols.Add(c)
		}
	case *Select:
		r.Props.OutputCols = e.Input.Props.OutputCols
	case *Project:
		for _, c := range e.ColIDs {
			r.Props.OutputCols.Add(c)
		}
	case *Join:
		e.Left.Props.OutputCols.ForEach(func(c opt.ColumnID) {
			r.Props.OutputCols.Add(c)
		})
		e.Right.Props.OutputCols.ForEach(func(c opt.ColumnID) {
			r.Props.OutputCols.Add(c)
		})
	}
}

func computeFreeCols(o opt.ColSet, s ScalarExpr) {
	if c, ok := s.(*ColRef); ok {
		o.Add(c.Id)
	} else {
		for i, n := 0, s.ChildCount(); i < n; i++ {
			computeFreeCols(o, s.Child(i).(ScalarExpr))
		}
	}
}
