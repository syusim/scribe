package builder

import (
	"fmt"

	"github.com/justinj/scribe/code/ast"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
)

func (b *builder) BuildScalar(e ast.Expr, scope *scope) (memo.ScalarExpr, error) {
	switch a := e.(type) {
	case ast.ColumnReference:
		id, typ, ok := scope.resolve(string(a))
		if !ok {
			return nil, fmt.Errorf("no column named %q", a)
		}
		return &memo.ColRef{
			Id:  id,
			Typ: typ,
		}, nil
	case *ast.ScalarFunc:
		args := make([]memo.ScalarExpr, len(a.Args))
		// TODO: i'm very inconsistent in my use of the two
		// arg-syntax, fix that up.
		for i := range a.Args {
			s, err := b.BuildScalar(a.Args[i], scope)
			if err != nil {
				return nil, err
			}
			args[i] = s
		}

		switch a.Op {
		case lang.Eq:
			if len(args) != 2 {
				return nil, fmt.Errorf("= takes two args, got %v", args)
			}
			if args[0].Type() != args[1].Type() {
				return nil, fmt.Errorf("arguments to = must be same type, got (= %v %v)", args[0].Type(), args[1].Type())
			}
		case lang.Plus, lang.Minus, lang.Times:
			for _, arg := range args {
				if arg.Type() != lang.Int {
					return nil, fmt.Errorf("args to %s must be int", a.Op)
				}
			}
		}
		return &memo.Func{
			Op:   a.Op,
			Args: args,
		}, nil
	case lang.Datum:
		return &memo.Constant{
			D: a,
		}, nil
	default:
		panic(fmt.Sprintf("unhandled: %T", e))
	}
}
