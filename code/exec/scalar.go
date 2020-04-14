package exec

import (
	"fmt"

	"github.com/justinj/scribe/code/lang"
)

type ScalarExpr interface {
	Eval(binding lang.Row) (lang.Datum, error)
}

type ColRef struct {
	Idx int
}

func (c *ColRef) Eval(binding lang.Row) (lang.Datum, error) {
	// TODO: panic with a sane error message if oob
	return binding[c.Idx], nil
}

type FuncInvocation struct {
	Op   lang.Func
	Args []ScalarExpr
}

func (f *FuncInvocation) Eval(binding lang.Row) (lang.Datum, error) {
	evaledArgs := make([]lang.Datum, len(f.Args))
	for i, a := range f.Args {
		d, err := a.Eval(binding)
		if err != nil {
			return nil, err
		}
		evaledArgs[i] = d
	}
	switch f.Op {
	case lang.Eq:
		return lang.DBool(lang.Compare(evaledArgs[0], evaledArgs[1]) == lang.EQ), nil
	case lang.Plus:
		sum := 0
		for _, x := range evaledArgs {
			sum += int(x.(lang.DInt))
		}
		return lang.DInt(sum), nil
	case lang.And:
		for _, x := range evaledArgs {
			if x == lang.DBool(false) {
				return lang.DBool(false), nil
			}
		}
		return lang.DBool(true), nil
	default:
		panic(fmt.Sprintf("unhandled operator %v", f.Op))
	}
}
