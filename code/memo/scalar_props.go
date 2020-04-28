package memo

import (
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/scalar"
)

type ScalarProps struct {
	FreeVars lang.ColSet
}

func (m *Memo) GetScalarProps(g scalar.Group) ScalarProps {
	e := lang.Unwrap(g)
	if p, ok := m.scalarProps[g]; ok {
		return p
	}

	props := ScalarProps{
		FreeVars: lang.ColSet{},
	}

	switch e := e.(type) {
	case *scalar.ColRef:
		props.FreeVars = lang.SetFromCols(e.Id)
	default:
		for i, n := 0, e.ChildCount(); i < n; i++ {
			props.FreeVars.UnionWith(m.GetScalarProps(e.Child(i).(scalar.Group)).FreeVars)
		}
	}

	// TODO: wait, are we re-storing this elsewhere?

	return props
}
