package render

import (
	"fmt"

	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/exec"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/memo"
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

func (b *builder) buildScalar(e scalar.Group, m lang.ColMap) (exec.ScalarExpr, error) {
	return exec.ScalarExpr(b.memo.Walk(e, func(in lang.Group) lang.Group {
		if ref, ok := in.(*scalar.ColRef); ok {
			idx, _ := m.Get(ref.Id)
			return &scalar.ExecColRef{
				Idx: idx,
				Typ: ref.Typ,
			}
		}
		return in
	}).(scalar.Group)), nil
}

func (b *builder) Build(e *memo.RelGroup, outCols []lang.ColumnID) (exec.Node, error) {
	n, m, err := b.build(e)
	if err != nil {
		return nil, err
	}

	out := make([]exec.ScalarExpr, len(outCols))
	for i := range outCols {
		idx, _ := m.Get(outCols[i])
		out[i] = &scalar.ExecColRef{
			Idx: idx,
			// TODO: do we need a type?
			Typ: 0,
		}
	}

	return exec.Project(n, out), nil
}

func (b *builder) build(e *memo.RelGroup) (exec.Node, lang.ColMap, error) {
	switch o := e.Unwrap().(type) {
	case *memo.Scan:
		tab, ok := b.cat.TableByName(o.TableName)
		if !ok {
			// This should have been verified already.
			panic(fmt.Sprintf("table %q not found", o.TableName))
		}
		if tab.IndexCount() <= o.Index {
			// TODO: this should be ensured to be impossible.
			// TODO: have a better error here.
			panic("invalid index")
		}

		idx := tab.Index(o.Index)

		var m lang.ColMap
		for i, id := range o.Cols {
			m.Set(id, i)
		}

		return exec.Scan(idx, o.Constraint), m, nil
	case *memo.Select:
		in, m, err := b.build(o.Input)
		if err != nil {
			return nil, lang.ColMap{}, err
		}

		// TODO: one unified scalar repr 2020
		var pred exec.ScalarExpr = lang.DBool(true)
		for _, p := range o.Filter.(*scalar.Filters).Filters {
			next, err := b.buildScalar(p, m)
			if err != nil {
				return nil, lang.ColMap{}, err
			}
			pred = &scalar.And{pred, next}
		}

		return exec.Select(in, pred), m, nil

	case *memo.HashJoin:
		build, leftMap, err := b.build(o.Build)
		if err != nil {
			return nil, lang.ColMap{}, err
		}

		probe, rightMap, err := b.build(o.Probe)
		if err != nil {
			return nil, lang.ColMap{}, err
		}

		// TODO: is there a neater way to do this?
		// We're just combining them.
		var m lang.ColMap
		leftMap.ForEach(func(from lang.ColumnID, to int) {
			m.Set(from, to)
		})
		rightMap.ForEach(func(from lang.ColumnID, to int) {
			m.Set(from, to+leftMap.Len())
		})

		leftIdxs := make([]lang.ColOrdinal, len(o.LeftCols))
		for i := range leftIdxs {
			idx, _ := leftMap.Get(o.LeftCols[i])
			leftIdxs[i] = lang.ColOrdinal(idx)
		}
		rightIdxs := make([]lang.ColOrdinal, len(o.RightCols))
		for i := range rightIdxs {
			idx, _ := rightMap.Get(o.RightCols[i])
			rightIdxs[i] = lang.ColOrdinal(idx)
		}

		return exec.Hash(build, probe, leftIdxs, rightIdxs), m, nil
	case *memo.Join:
		left, leftMap, err := b.build(o.Left)
		if err != nil {
			return nil, lang.ColMap{}, err
		}

		right, rightMap, err := b.build(o.Right)
		if err != nil {
			return nil, lang.ColMap{}, err
		}

		// TODO: is there a neater way to do this?
		// We're just combining them.
		var m lang.ColMap
		leftMap.ForEach(func(from lang.ColumnID, to int) {
			m.Set(from, to)
		})
		rightMap.ForEach(func(from lang.ColumnID, to int) {
			m.Set(from, to+leftMap.Len())
		})

		var pred exec.ScalarExpr = lang.DBool(true)
		for _, p := range o.On.(*scalar.Filters).Filters {
			next, err := b.buildScalar(p, m)
			if err != nil {
				return nil, lang.ColMap{}, err
			}
			pred = &scalar.And{pred, next}
		}

		// TODO: make a real join operator!
		return exec.Select(
			exec.Cross(left, right),
			pred,
		), m, nil
	case *memo.Project:
		in, m, err := b.build(o.Input)
		if err != nil {
			return nil, lang.ColMap{}, err
		}

		outMap := lang.ColMap{}

		exprs := make([]exec.ScalarExpr, 0, len(o.Projections))
		for i := range o.Projections {
			p, err := b.buildScalar(o.Projections[i], m)
			if err != nil {
				return nil, lang.ColMap{}, err
			}
			exprs = append(exprs, p)
			outMap.Set(o.ColIDs[i], i)
		}

		o.PassthroughCols.ForEach(func(c lang.ColumnID) {
			idx, _ := m.Get(c)
			// Just synthesize a col ref.
			exprs = append(exprs, &scalar.ExecColRef{
				// TODO: do we need a type?
				// TODO: standardize Typ vs Type
				Typ: 0,
				Idx: idx,
			})
			outMap.Set(c, len(exprs)-1)
		})

		return exec.Project(in, exprs), outMap, nil
	case *memo.Root:
		return b.build(o.Input)

	default:
		panic(fmt.Sprintf("unhandled: %T", e.Unwrap()))
	}
}
