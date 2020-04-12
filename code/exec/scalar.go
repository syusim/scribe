package exec

import (
	"fmt"

	"github.com/justinj/scribe/code/lang"
)

type ScalarExpr interface {
	Eval(binding lang.Row) (lang.Datum, error)
}

type ColRef struct {
	idx int
}

func (c *ColRef) Eval(binding lang.Row) (lang.Datum, error) {
	// TODO: panic with a sane error message if oob
	return binding[c.idx], nil
}

type FuncInvocation struct {
	op   lang.Func
	args []ScalarExpr
}

func (f *FuncInvocation) Eval(binding lang.Row) (lang.Datum, error) {
	evaledArgs := make([]lang.Datum, len(f.args))
	for i, a := range f.args {
		d, err := a.Eval(binding)
		if err != nil {
			return nil, err
		}
		evaledArgs[i] = d
	}
	switch f.op {
	case lang.Eq:
		return lang.DBool(lang.Compare(evaledArgs[0], evaledArgs[1]) == lang.EQ), nil
	default:
		panic(fmt.Sprintf("unhandled operator %v", f.op))
	}
}
