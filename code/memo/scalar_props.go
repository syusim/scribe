package memo

import (
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

type ScalarProps struct {
	FreeVars opt.ColSet
}

func (m *Memo) GetScalarProps(e scalar.Expr) ScalarProps {
	if p, ok := m.scalarProps[e]; ok {
		return p
	}

	props := ScalarProps{
		FreeVars: opt.ColSet{},
	}

	switch e := e.(type) {
	case *scalar.ColRef:
		props.FreeVars = opt.SetFromCols(e.Id)
	default:
		for i, n := 0, e.ChildCount(); i < n; i++ {
			props.FreeVars.UnionWith(m.GetScalarProps(e.Child(i).(scalar.Expr)).FreeVars)
		}
	}

	return props
}
