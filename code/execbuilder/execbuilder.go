package execbuilder

import (
	"fmt"

	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/exec"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/opt"
)

type builder struct {
	cat *cat.Catalog
}

func New(cat *cat.Catalog) *builder {
	return &builder{
		cat: cat,
	}
}

// TODO: we should get rid of this, just walk the tree, keep
// everything the same, but replace the absolute column references
// with ordinal references. bingo bango
func (b *builder) buildScalar(e memo.ScalarExpr, m opt.ColMap) (exec.ScalarExpr, error) {
	switch s := e.(type) {
	case *memo.Constant:
		return s.D, nil
	case *memo.ColRef:
		i, ok := m.Get(s.Id)
		if !ok {
			panic(fmt.Sprintf("no column with id %d", s.Id))
		}
		return &exec.ColRef{
			Idx: i,
		}, nil
	case *memo.Func:
		args := make([]exec.ScalarExpr, len(s.Args))
		for i := range s.Args {
			a, err := b.buildScalar(s.Args[i], m)
			if err != nil {
				return nil, err
			}
			args[i] = a
		}
		return &exec.FuncInvocation{
			Op:   s.Op,
			Args: args,
		}, nil

	default:
		panic(fmt.Sprintf("unhandled: %T", s))
	}
}

func (b *builder) Build(e memo.RelExpr) (exec.Node, opt.ColMap, error) {
	switch o := e.E.(type) {
	case *memo.Scan:
		tab, ok := b.cat.TableByName(o.TableName)
		if !ok {
			// This should have been verified already.
			panic(fmt.Sprintf("table %q not found", o.TableName))
		}
		if tab.IndexCount() == 0 {
			// TODO: this should be ensured to be impossible.
			// TODO: have a better error here.
			panic("no indexes buddy!")
		}

		// TODO: pass in which one to use
		idx := tab.Index(0)
		iter := idx.Scan()

		var m opt.ColMap
		for i, id := range o.Cols {
			m.Set(id, i)
		}

		return exec.Scan(iter), m, nil
	case *memo.Select:
		in, m, err := b.Build(o.Input)
		if err != nil {
			return nil, opt.ColMap{}, err
		}

		pred, err := b.buildScalar(o.Filter, m)
		if err != nil {
			return nil, opt.ColMap{}, err
		}

		return exec.Select(in, pred), m, nil

	case *memo.Join:
		left, leftMap, err := b.Build(o.Left)
		if err != nil {
			return nil, opt.ColMap{}, err
		}

		right, rightMap, err := b.Build(o.Right)
		if err != nil {
			return nil, opt.ColMap{}, err
		}

		// TODO: is there a neater way to do this?
		// We're just combining them.
		var m opt.ColMap
		leftMap.ForEach(func(from opt.ColumnID, to int) {
			m.Set(from, to)
		})
		rightMap.ForEach(func(from opt.ColumnID, to int) {
			m.Set(from, to+leftMap.Len())
		})

		on, err := b.buildScalar(o.On, m)
		if err != nil {
			return nil, opt.ColMap{}, err
		}

		// TODO: make a real join operator!
		return exec.Select(
			exec.Cross(left, right),
			on,
		), m, nil
	case *memo.Project:
		in, m, err := b.Build(o.Input)
		if err != nil {
			return nil, opt.ColMap{}, err
		}

		outMap := opt.ColMap{}

		exprs := make([]exec.ScalarExpr, len(o.Projections))
		for i := range o.Projections {
			p, err := b.buildScalar(o.Projections[i], m)
			if err != nil {
				return nil, opt.ColMap{}, err
			}
			exprs[i] = p
			outMap.Set(o.ColIDs[i], i)
		}

		return exec.Project(in, exprs), outMap, nil

	default:
		panic(fmt.Sprintf("unhandled: %T", e.E))
	}
}
