package builder

import (
	"fmt"

	"github.com/justinj/scribe/code/ast"
	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/constraint"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

type colInfo struct {
	name string
	id   opt.ColumnID
	typ  lang.Type
}

type builder struct {
	cat *cat.Catalog
	// TODO: extract this into a metadata struct,
	// since it needs to be accessed from elsewhere
	// (say, the memo formatter)
	cols []colInfo

	memo *memo.Memo
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

// TODO: just refer to the memo's handle onto catalog?
func New(cat *cat.Catalog, memo *memo.Memo) *builder {
	return &builder{
		cat:  cat,
		memo: memo,
	}
}

func (b *builder) Build(e ast.RelExpr) (*memo.RelGroup, *scope, error) {
	// This is legal at the root, and nowhere else.
	if o, ok := e.(*ast.OrderBy); ok {
		m, s, err := b.build(o.Input)
		if err != nil {
			return nil, nil, err
		}

		var ord opt.Ordering
		for _, col := range o.ColNames {
			c, _, ok := s.resolve(col)
			if !ok {
				return nil, nil, fmt.Errorf("no col named %q", col)
			}
			ord = append(ord, c)
		}

		return b.memo.Root(m, ord), s, nil
	}
	return b.build(e)
}

// TODO: extract each arm
func (b *builder) build(e ast.RelExpr) (*memo.RelGroup, *scope, error) {
	switch a := e.(type) {
	case *ast.TableRef:
		tab, ok := b.cat.TableByName(a.Name)
		if !ok {
			return nil, nil, fmt.Errorf("no table named %q", a.Name)
		}

		cols := make([]opt.ColumnID, tab.ColumnCount())
		s := newScope()
		for i, n := 0, tab.ColumnCount(); i < n; i++ {
			col := tab.Column(i)
			id := b.addCol(col.Name, col.Type)
			s.addCol(col.Name, id, col.Type)
			cols[i] = id
		}

		// TODO: look it up in the catalog to ensure it exists.
		// TODO: make a logical version of this?
		return b.memo.Scan(a.Name, cols, 0, constraint.Constraint{}), s, nil
	case *ast.Select:
		input, s, err := b.build(a.Input)
		if err != nil {
			return nil, nil, err
		}
		filter, err := b.BuildScalar(a.Predicate, s)
		if err != nil {
			return nil, nil, err
		}
		return b.memo.Select(input, filter), s, nil
	case *ast.Join:
		left, leftScope, err := b.build(a.Left)
		if err != nil {
			return nil, nil, err
		}

		right, rightScope, err := b.build(a.Right)
		if err != nil {
			return nil, nil, err
		}

		s := appendScopes(leftScope, rightScope)

		on, err := b.BuildScalar(a.On, s)
		if err != nil {
			return nil, nil, err
		}

		return b.memo.Join(left, right, on), s, nil
	case *ast.Project:
		in, inScope, err := b.build(a.Input)
		if err != nil {
			return nil, nil, err
		}

		exprs := make([]scalar.Group, 0, len(a.Exprs))
		outCols := make([]opt.ColumnID, 0, len(exprs))

		outScope := newScope()
		var passthrough opt.ColSet

		for i, e := range a.Exprs {
			proj, err := b.BuildScalar(e, inScope)
			if err != nil {
				return nil, nil, err
			}
			// Sneak a peek!
			if v, ok := proj.(*scalar.ColRef); ok {
				passthrough.Add(v.Id)
				outScope.addCol(a.Aliases[i], v.Id, proj.Type())
			} else {
				exprs = append(exprs, proj)
				newCol := b.addCol(a.Aliases[i], proj.Type())
				outCols = append(outCols, newCol)
				outScope.addCol(a.Aliases[i], newCol, proj.Type())
			}
		}

		return b.memo.Project(in, outCols, exprs, passthrough), outScope, nil
	case *ast.As:
		expr, inScope, err := b.build(a.Input)
		if err != nil {
			return nil, nil, err
		}
		outScope := newScope()
		if len(a.ColNames) > len(inScope.cols) {
			return nil, nil, fmt.Errorf("too many cols!")
		}
		for i := range a.ColNames {
			outScope.addCol(a.ColNames[i], inScope.cols[i].id, inScope.cols[i].typ)
		}
		return expr, outScope, nil
	default:
		panic(fmt.Sprintf("unhandled: %T", e))
	}
}
