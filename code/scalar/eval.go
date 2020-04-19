package scalar

import (
	"fmt"

	"github.com/justinj/scribe/code/lang"
)

func Eval(g Group, binding lang.Row) (lang.Datum, error) {
	e := lang.Unwrap(g)
	switch e := e.(type) {
	case *Constant:
		return e.D, nil
	case *ColRef:
		panic("cannot Eval ColRef, make sure this expression has been execbuilt!")
	case *ExecColRef:
		// TODO: panic with a sane error message if oob
		return binding[e.Idx], nil
	case *Plus:
		left, err := Eval(e.Left, binding)
		if err != nil {
			return nil, err
		}
		right, err := Eval(e.Right, binding)
		if err != nil {
			return nil, err
		}
		return lang.DInt(left.(lang.DInt) + right.(lang.DInt)), nil
	case *And:
		left, err := Eval(e.Left, binding)
		if err != nil {
			return nil, err
		}
		if left != lang.DBool(true) {
			return lang.DBool(false), nil
		}
		return Eval(e.Right, binding)
	case *Eq:
		left, err := Eval(e.Left, binding)
		if err != nil {
			return nil, err
		}
		right, err := Eval(e.Right, binding)
		if err != nil {
			return nil, err
		}
		if lang.Compare(left, right) == lang.EQ {
			return lang.DBool(true), nil
		}
		return lang.DBool(false), nil
	case lang.Datum:
		return e, nil
	default:
		panic(fmt.Sprintf("unhandled Eval: %T", e))
	}
}
