package memo

import (
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

func (m *Memo) Scan(tableName string, cols []opt.ColumnID) *RelExpr {
	return m.internScan(Scan{
		TableName: tableName,
		Cols:      cols,
	})
}

func (m *Memo) matchRules(args []interface{}, rules []rule) lang.Expr {
	for _, r := range rules {
		if n := r(m, args); n != nil {
			return n
		}
	}
	return nil
}

func (m *Memo) Join(left, right *RelExpr, on scalar.Expr) *RelExpr {
	if e := m.matchRules([]interface{}{left, right, on}, []rule{
		WrapJoinConditionInFilters,
		UnfoldJoinCondition,
		PushFilterIntoJoinLeft,
		PushFilterIntoJoinRight,
	}); e != nil {
		return e.(*RelExpr)
	}

	return m.internJoin(Join{
		Left:  left,
		Right: right,
		On:    on,
	})
}

// TODO: standardize on xxxIDs vs. xxxIds
func (m *Memo) Project(
	input *RelExpr,
	colIDs []opt.ColumnID,
	projections []scalar.Expr,
	passthrough opt.ColSet,
) *RelExpr {
	if e := m.matchRules([]interface{}{input, colIDs, projections, passthrough}, []rule{
		EliminateProject,
		MergeProjectProject,
	}); e != nil {
		return e.(*RelExpr)
	}

	return m.internProject(Project{
		Input:           input,
		ColIDs:          colIDs,
		Projections:     projections,
		PassthroughCols: passthrough,
	})
}

func (m *Memo) Select(input *RelExpr, filter scalar.Expr) *RelExpr {
	if e := m.matchRules([]interface{}{input, filter}, []rule{
		WrapSelectConditionInFilters,
		UnfoldSelectCondition,
		// TODO: Can this be its own rule?
		SimplifySelectFilters,
		EliminateSelect,
	}); e != nil {
		return e.(*RelExpr)
	}

	// TODO: make this a real rule
	// MergeSelectJoin
	if j, ok := input.E.(*Join); ok {
		return m.Join(
			j.Left,
			j.Right,
			m.Filters(concat(j.On.(*scalar.Filters), filter.(*scalar.Filters))),
		)
	}

	return m.internSelect(Select{
		Input:  input,
		Filter: filter,
	})
}

func (m *Memo) Root(input *RelExpr, ordering opt.Ordering) *RelExpr {
	return m.internRoot(Root{
		Input:    input,
		Ordering: ordering,
	})
}

func (m *Memo) Constant(d lang.Datum) scalar.Expr {
	return m.internConstant(scalar.Constant{d})
}

func (m *Memo) ColRef(id opt.ColumnID, typ lang.Type) scalar.Expr {
	return m.internColRef(scalar.ColRef{id, typ})
}

func (m *Memo) Plus(left, right scalar.Expr) scalar.Expr {
	if e := m.matchRules([]interface{}{left, right}, []rule{
		FoldZeroPlus,
		FoldPlusZero,
		AssociatePlus,

		// Goofy rules to simplify one very specific case.
		SimplifyPlusPlus,
		SimplifyPlusTimes,
	}); e != nil {
		return e.(scalar.Expr)
	}

	return m.internPlus(scalar.Plus{left, right})
}

func (m *Memo) Times(left, right scalar.Expr) scalar.Expr {
	if e := m.matchRules([]interface{}{left, right}, []rule{}); e != nil {
		return e.(scalar.Expr)
	}

	return m.internTimes(scalar.Times{left, right})
}

func (m *Memo) And(left, right scalar.Expr) scalar.Expr {
	return m.internAnd(scalar.And{left, right})
}

func (m *Memo) Filters(args []scalar.Expr) scalar.Expr {
	return m.internFilters(scalar.Filters{args})
}

func (m *Memo) Func(op lang.Func, args []scalar.Expr) scalar.Expr {
	return m.internFunc(scalar.Func{op, args})
}
