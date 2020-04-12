package builder

import (
	"fmt"

	"github.com/justinj/scribe/code/ast"
	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/opt"
)

type colInfo struct {
	name string
	id   opt.ColumnID
	typ  lang.Type
}

type builder struct {
	cat  *cat.Catalog
	cols []colInfo
}

func (b *builder) addCol(name string, typ lang.Type) opt.ColumnID {
	id := opt.ColumnID(len(b.cols) + 1)
	b.cols = append(b.cols, colInfo{
		name: name,
		id:   id,
		typ:  typ,
	})
	return id
}

func New(cat *cat.Catalog) *builder {
	return &builder{
		cat: cat,
	}
}

// TODO: pull this out into a file and decompose
func (b *builder) buildScalar(e ast.Expr, scope *scope) (memo.ScalarExpr, error) {
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
			s, err := b.buildScalar(a.Args[i], scope)
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

// TODO: extract each arm
func (b *builder) Build(e ast.RelExpr) (memo.RelExpr, *scope, error) {
	switch a := e.(type) {
	case *ast.TableRef:
		tab, ok := b.cat.TableByName(a.Name)
		if !ok {
			return memo.RelExpr{}, nil, fmt.Errorf("no table named %q", a.Name)
		}

		cols := make([]opt.ColumnID, tab.ColumnCount())
		s := newScope()
		for i, n := 0, tab.ColumnCount(); i < n; i++ {
			col := tab.Column(i)
			id := b.addCol(col.Name, col.Type)
			s.addCol(col.Name, id, col.Type)
			cols[i] = id
		}

		// TODO: look it up in the catalog.
		return memo.Wrap(&memo.Scan{
			TableName: a.Name,
			Cols:      cols,
		}), s, nil
	case *ast.Select:
		input, s, err := b.Build(a.Input)
		if err != nil {
			return memo.RelExpr{}, nil, nil
		}
		filter, err := b.buildScalar(a.Predicate, s)
		if err != nil {
			return memo.RelExpr{}, nil, err
		}
		return memo.Wrap(&memo.Select{
			Input:  input,
			Filter: filter,
		}), s, nil
	case *ast.Join:
		left, leftScope, err := b.Build(a.Left)
		if err != nil {
			return memo.RelExpr{}, nil, nil
		}

		right, rightScope, err := b.Build(a.Right)
		if err != nil {
			return memo.RelExpr{}, nil, nil
		}

		s := appendScopes(leftScope, rightScope)

		on, err := b.buildScalar(a.On, s)
		if err != nil {
			return memo.RelExpr{}, nil, nil
		}

		return memo.Wrap(&memo.Join{
			Left:  left,
			Right: right,
			On:    on,
		}), s, nil
	case *ast.Project:
		in, inScope, err := b.Build(a.Input)
		if err != nil {
			return memo.RelExpr{}, nil, nil
		}

		exprs := make([]memo.ScalarExpr, len(a.Exprs))
		outCols := make([]opt.ColumnID, len(exprs))

		outScope := newScope()

		for i, e := range a.Exprs {
			proj, err := b.buildScalar(e, inScope)
			if err != nil {
				return memo.RelExpr{}, nil, err
			}
			exprs[i] = proj
			outCols[i] = b.addCol(a.Aliases[i], proj.Type())
			outScope.addCol(a.Aliases[i], outCols[i], proj.Type())
		}

		return memo.Wrap(&memo.Project{
			Input: in,

			ColIDs:      outCols,
			Projections: exprs,
		}), outScope, nil
	default:
		panic(fmt.Sprintf("unhandled: %T", e))
	}
}
