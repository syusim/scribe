package memo

import (
	"fmt"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/scalar"
)

// TODO: this should be codegenned
func (m *Memo) Walk(in lang.Group, f func(e lang.Group) lang.Group) lang.Group {
	switch e := in.(type) {
	// All the 0-child ops.
	case *scalar.ColRef, *scalar.ExecColRef,
		*scalar.Constant:
		return f(in)
	case *scalar.Plus:
		left := m.Walk(e.Left, f).(scalar.Group)
		right := m.Walk(e.Right, f).(scalar.Group)
		return f(m.Plus(left, right))
	case *scalar.Times:
		left := m.Walk(e.Left, f).(scalar.Group)
		right := m.Walk(e.Right, f).(scalar.Group)
		return f(m.Times(left, right))
	case *scalar.Eq:
		left := m.Walk(e.Left, f).(scalar.Group)
		right := m.Walk(e.Right, f).(scalar.Group)
		return f(m.Eq(left, right))
	case *scalar.And:
		left := m.Walk(e.Left, f).(scalar.Group)
		right := m.Walk(e.Right, f).(scalar.Group)
		return f(m.And(left, right))
	case *scalar.Func:
		args := make([]scalar.Group, len(e.Args))
		for i := range e.Args {
			args[i] = m.Walk(e.Args[i], f).(scalar.Group)
		}
		return f(m.Func(e.Op, args))
	default:
		panic(fmt.Sprintf("unhandled: %T", in))
	}
}
