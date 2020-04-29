package memo

import (
	"github.com/justinj/scribe/code/constraint"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/scalar"
)

func (m *Memo) Scan(tableName string, cols []lang.ColumnID, indexId int, constraint constraint.Constraint) *RelGroup {
	return m.internScan(Scan{
		TableName:  tableName,
		Cols:       cols,
		Index:      indexId,
		Constraint: constraint,
	})
}

func (m *Memo) matchRules(args []interface{}, rules []rule) lang.Group {
	for _, r := range rules {
		if n := r(m, args); n != nil {
			return n
		}
	}
	return nil
}

func (m *Memo) Join(left, right *RelGroup, on scalar.Group) *RelGroup {
	if e := m.matchRules([]interface{}{left, right, on}, []rule{
		WrapJoinConditionInFilters,
		UnfoldJoinCondition,
		PushFilterIntoJoinLeft,
		PushFilterIntoJoinRight,
	}); e != nil {
		return e.(*RelGroup)
	}

	return m.internJoin(Join{
		Left:  left,
		Right: right,
		On:    on,
	})
}

func (m *Memo) HashJoin(build, probe *RelGroup, leftCols, rightCols []lang.ColumnID) *RelGroup {
	return m.internHashJoin(HashJoin{
		Build:     build,
		Probe:     probe,
		LeftCols:  leftCols,
		RightCols: rightCols,
	})
}

// TODO: standardize on xxxIDs vs. xxxIds
func (m *Memo) Project(
	input *RelGroup,
	colIDs []lang.ColumnID,
	projections []scalar.Group,
	passthrough lang.ColSet,
) *RelGroup {
	if e := m.matchRules([]interface{}{input, colIDs, projections, passthrough}, []rule{
		EliminateProject,
		MergeProjectProject,
	}); e != nil {
		return e.(*RelGroup)
	}

	return m.internProject(Project{
		Input:           input,
		ColIDs:          colIDs,
		Projections:     projections,
		PassthroughCols: passthrough,
	})
}

func (m *Memo) Select(input *RelGroup, filter scalar.Group) *RelGroup {
	if e := m.matchRules([]interface{}{input, filter}, []rule{
		WrapSelectConditionInFilters,
		UnfoldSelectCondition,
		// TODO: Can this be its own rule?
		SimplifySelectFilters,
		EliminateSelect,
	}); e != nil {
		return e.(*RelGroup)
	}

	// TODO: make this a real rule
	// MergeSelectJoin
	if j, ok := input.Unwrap().(*Join); ok {
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

func (m *Memo) Root(input *RelGroup, ordering lang.Ordering) *RelGroup {
	return m.internRoot(Root{
		Input:    input,
		Ordering: ordering,
	})
}

func (m *Memo) Sort(input *RelGroup, ordering lang.Ordering) *RelGroup {
	return m.internSort(Sort{
		Input:    input,
		Ordering: ordering,
	})
}

func (m *Memo) Constant(d lang.Datum) scalar.Group {
	return m.internConstant(scalar.Constant{d})
}

func (m *Memo) ColRef(id lang.ColumnID, typ lang.Type) scalar.Group {
	return m.internColRef(scalar.ColRef{id, typ})
}

func (m *Memo) Plus(left, right scalar.Group) scalar.Group {
	if e := m.matchRules([]interface{}{left, right}, []rule{
		FoldZeroPlus,
		FoldPlusZero,
		AssociatePlus,

		// Goofy rules to simplify one very specific case.
		SimplifyPlusPlus,
		SimplifyPlusTimes,
	}); e != nil {
		return e.(scalar.Group)
	}

	return m.internPlus(scalar.Plus{left, right})
}

func (m *Memo) Times(left, right scalar.Group) scalar.Group {
	if e := m.matchRules([]interface{}{left, right}, []rule{}); e != nil {
		return e.(scalar.Group)
	}

	return m.internTimes(scalar.Times{left, right})
}

func (m *Memo) And(left, right scalar.Group) scalar.Group {
	return m.internAnd(scalar.And{left, right})
}

func (m *Memo) Eq(left, right scalar.Group) scalar.Group {
	return m.internEq(scalar.Eq{left, right})
}

func (m *Memo) Filters(args []scalar.Group) scalar.Group {
	return m.internFilters(scalar.Filters{args})
}

func (m *Memo) Func(op lang.Func, args []scalar.Group) scalar.Group {
	return m.internFunc(scalar.Func{op, args})
}
