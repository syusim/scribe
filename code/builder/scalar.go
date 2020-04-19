package builder

import (
	"fmt"

	"github.com/justinj/scribe/code/ast"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/scalar"
)

func (b *builder) BuildScalar(e ast.Expr, scope *scope) (scalar.Group, error) {
	switch a := e.(type) {
	case ast.ColumnReference:
		id, typ, ok := scope.resolve(string(a))
		if !ok {
			return nil, fmt.Errorf("no column named %q", a)
		}
		return b.memo.ColRef(id, typ), nil
	case *ast.ScalarFunc:
		args := make([]scalar.Group, len(a.Args))
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
		case lang.And:
			if len(args) != 2 {
				return nil, fmt.Errorf("and takes 2 args")
			}
			for _, arg := range args {
				if arg.Type() != lang.Bool {
					return nil, fmt.Errorf("args to %s must be int", a.Op)
				}
			}
			return b.memo.And(args[0], args[1]), nil
		case lang.Plus:
			if len(args) != 2 {
				return nil, fmt.Errorf("+ takes 2 args")
			}
			for _, arg := range args {
				if arg.Type() != lang.Int {
					return nil, fmt.Errorf("args to %s must be int", a.Op)
				}
			}
			return b.memo.Plus(args[0], args[1]), nil
		case lang.Times:
			if len(args) != 2 {
				return nil, fmt.Errorf("* takes 2 args")
			}
			for _, arg := range args {
				if arg.Type() != lang.Int {
					return nil, fmt.Errorf("args to %s must be int", a.Op)
				}
			}
			return b.memo.Times(args[0], args[1]), nil
		case lang.Minus:
			for _, arg := range args {
				if arg.Type() != lang.Int {
					return nil, fmt.Errorf("args to %s must be int", a.Op)
				}
			}
		}
		return b.memo.Func(a.Op, args), nil
	case lang.Datum:
		return b.memo.Constant(a), nil
	default:
		panic(fmt.Sprintf("unhandled: %T", e))
	}
}
