package memo

import (
	"fmt"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/scalar"
)

// TODO: fix this. make an op enum. code gen it
// ughhhh this code suckksssss
func (m *Memo) Render(e lang.Expr, args []lang.Group) lang.Group {
	if e, ok := e.(scalar.Group); ok {
		return e
	}
	var res lang.Expr
	switch e := e.(type) {
	case *Select:
		res = &Select{args[0].(*RelGroup), args[1].(scalar.Group)}
	case *Join:
		res = &Join{args[0].(*RelGroup), args[1].(*RelGroup), args[2].(scalar.Group)}
	case *Project:
		projs := make([]scalar.Group, len(args)-1)
		for i := range projs {
			projs[i] = args[i+1].(scalar.Group)
		}
		res = &Project{
			Input:           args[0].(*RelGroup),
			ColIDs:          e.ColIDs,
			Projections:     projs,
			PassthroughCols: e.PassthroughCols,
		}
	case *Scan:
		cpy := *e
		res = &cpy
	case *Sort:
		cpy := *e
		cpy.Input = args[0].(*RelGroup)
		res = &cpy
	case *Root:
		cpy := *e
		cpy.Input = args[0].(*RelGroup)
		res = &cpy
	default:
		panic(fmt.Sprintf("unhandled: %T", e))
	}

	return &RelGroup{
		Es: []relExpr{res},
		// TODO: do we care if we have the props here? we have them
		// but that's not really the point.
	}
}

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
