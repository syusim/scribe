package memo

import (
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

// TODO: pull this into a pkg/file
type Props struct {
	OutputCols opt.ColSet
}

func buildProps(r *RelGroup) {
	switch e := r.Unwrap().(type) {
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
		r.Props.OutputCols.UnionWith(e.PassthroughCols)
	case *Join:
		e.Left.Props.OutputCols.ForEach(func(c opt.ColumnID) {
			r.Props.OutputCols.Add(c)
		})
		e.Right.Props.OutputCols.ForEach(func(c opt.ColumnID) {
			r.Props.OutputCols.Add(c)
		})
	}
}

func computeFreeCols(o opt.ColSet, s scalar.Group) {
	if c, ok := s.(*scalar.ColRef); ok {
		o.Add(c.Id)
	} else {
		for i, n := 0, s.ChildCount(); i < n; i++ {
			computeFreeCols(o, s.Child(i).(scalar.Group))
		}
	}
}
