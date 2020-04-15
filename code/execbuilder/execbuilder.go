package execbuilder

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/exec"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

type builder struct {
	cat  *cat.Catalog
	memo *memo.Memo
}

func New(cat *cat.Catalog, memo *memo.Memo) *builder {
	return &builder{
		cat:  cat,
		memo: memo,
	}
}

func (b *builder) buildScalar(e scalar.Expr, m opt.ColMap) (exec.ScalarExpr, error) {
	return exec.ScalarExpr(b.memo.Walk(e, func(in lang.Expr) lang.Expr {
		if ref, ok := in.(*scalar.ColRef); ok {
			idx, _ := m.Get(ref.Id)
			return &scalar.ExecColRef{
				Idx: idx,
				Typ: ref.Typ,
			}
		}
		return in
	}).(scalar.Expr)), nil
}

func (b *builder) Build(e *memo.RelExpr) (exec.Node, opt.ColMap, error) {
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

		// TODO: one unified scalar repr 2020
		var pred exec.ScalarExpr = lang.DBool(true)
		for _, p := range o.Filter {
			next, err := b.buildScalar(p, m)
			if err != nil {
				return nil, opt.ColMap{}, err
			}
			pred = &scalar.And{pred, next}
		}

		spew.Dump(pred)

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

		var pred exec.ScalarExpr = lang.DBool(true)
		for _, p := range o.On {
			next, err := b.buildScalar(p, m)
			if err != nil {
				return nil, opt.ColMap{}, err
			}
			pred = &scalar.And{pred, next}
		}

		// TODO: make a real join operator!
		return exec.Select(
			exec.Cross(left, right),
			pred,
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
