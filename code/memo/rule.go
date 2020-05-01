package memo

import (
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/scalar"
)

type rule func(m *Memo, args []interface{}) lang.Group

// Plus Rules.

func FoldZeroPlus(m *Memo, args []interface{}) lang.Group {
	if eqConst(args[0].(scalar.Group), lang.DInt(0)) {
		return args[1].(lang.Group)
	}
	return nil
}

func FoldPlusZero(m *Memo, args []interface{}) lang.Group {
	if eqConst(args[1].(scalar.Group), lang.DInt(0)) {
		return args[0].(lang.Group)
	}
	return nil
}

func AssociatePlus(m *Memo, args []interface{}) lang.Group {
	if left, ok := args[0].(*scalar.Plus); ok {
		return m.Plus(
			left.Left,
			m.Plus(
				left.Right,
				args[1].(scalar.Group),
			),
		)
	}
	return nil
}

// Times Rules.

func FoldTimesOne(m *Memo, args []interface{}) lang.Group {
	if eqConst(args[0].(scalar.Group), lang.DInt(1)) {
		return args[1].(lang.Group)
	}
	return nil
}

func FoldOneTimes(m *Memo, args []interface{}) lang.Group {
	if eqConst(args[1].(scalar.Group), lang.DInt(1)) {
		return args[0].(lang.Group)
	}
	return nil
}

func AssociateTimes(m *Memo, args []interface{}) lang.Group {
	if left, ok := args[0].(*scalar.Times); ok {
		return m.Times(
			left.Left,
			m.Times(
				left.Right,
				args[1].(scalar.Group),
			),
		)
	}
	return nil
}

// Join Rules.

func WrapJoinConditionInFilters(m *Memo, args []interface{}) lang.Group {
	left, right, on := args[0].(*RelGroup), args[1].(*RelGroup), args[2].(scalar.Group)
	if _, ok := on.(*scalar.Filters); !ok {
		return m.Join(
			left, right, m.Filters([]scalar.Group{on}),
		)
	}
	return nil
}

func UnfoldJoinCondition(m *Memo, args []interface{}) lang.Group {
	left, right, on := args[0].(*RelGroup), args[1].(*RelGroup), args[2].(*scalar.Filters)
	newFilters := unfoldFilters(on.Filters)

	if newFilters != nil {
		return m.Join(left, right, m.Filters(newFilters))
	}

	return nil
}

func PushFilterIntoJoinLeft(m *Memo, args []interface{}) lang.Group {
	left, right, on := args[0].(*RelGroup), args[1].(*RelGroup), args[2].(*scalar.Filters)
	bound, unbound := extractBoundUnbound(m, on.Filters, left.Props.OutputCols)

	if len(bound) > 0 {
		return m.Join(
			m.Select(
				left,
				m.Filters(bound),
			),
			right,
			m.Filters(unbound),
		)
	}
	return nil
}

func PushFilterIntoJoinRight(m *Memo, args []interface{}) lang.Group {
	left, right, on := args[0].(*RelGroup), args[1].(*RelGroup), args[2].(*scalar.Filters)
	bound, unbound := extractBoundUnbound(m, on.Filters, right.Props.OutputCols)

	if len(bound) > 0 {
		return m.Join(
			left,
			m.Select(
				right,
				m.Filters(bound),
			),
			m.Filters(unbound),
		)
	}
	return nil
}

// Select Rules.

func WrapSelectConditionInFilters(m *Memo, args []interface{}) lang.Group {
	input, filter := args[0].(*RelGroup), args[1].(scalar.Group)
	if _, ok := filter.(*scalar.Filters); !ok {
		return m.Select(
			input, m.Filters([]scalar.Group{filter}),
		)
	}
	return nil
}

func UnfoldSelectCondition(m *Memo, args []interface{}) lang.Group {
	input, filter := args[0].(*RelGroup), args[1].(*scalar.Filters)
	newFilters := unfoldFilters(filter.Filters)

	if newFilters != nil {
		return m.Select(input, m.Filters(newFilters))
	}

	return nil
}

func SimplifySelectFilters(m *Memo, args []interface{}) lang.Group {
	input, filter := args[0].(*RelGroup), args[1].(*scalar.Filters)

	var newFilters []scalar.Group
	for i, f := range filter.Filters {
		if eqConst(f, lang.DBool(true)) {
			newFilters = make([]scalar.Group, i)
			copy(newFilters, filter.Filters)
		} else if newFilters != nil {
			newFilters = append(newFilters, f)
		}
	}
	if newFilters != nil {
		return m.Select(
			input,
			m.Filters(newFilters),
		)
	}

	return nil
}

func EliminateSelect(m *Memo, args []interface{}) lang.Group {
	input, filter := args[0].(*RelGroup), args[1].(*scalar.Filters)
	if len(filter.Filters) == 0 {
		return input
	}

	return nil
}

func MergeSelects(m *Memo, args []interface{}) lang.Group {
	input, filter := args[0].(*RelGroup), args[1].(*scalar.Filters)
	if s, ok := input.Unwrap().(*Select); ok {
		innerInput, innerFilter := s.Input, s.Filter.(*scalar.Filters)
		return m.Select(
			innerInput,
			m.Filters(concat(
				filter,
				innerFilter,
			)),
		)
	}

	return nil
}

// Project Rules.

func EliminateProject(m *Memo, args []interface{}) lang.Group {
	input, _, projections, passthrough := args[0].(*RelGroup), args[1].([]lang.ColumnID), args[2].([]scalar.Group), args[3].(lang.ColSet)

	if len(projections) > 0 {
		return nil
	}

	if input.Props.OutputCols.Equals(passthrough) {
		return input
	}

	return nil
}

func MergeProjectProject(m *Memo, args []interface{}) lang.Group {
	input, colIDs, projections, passthrough := args[0].(*RelGroup), args[1].([]lang.ColumnID), args[2].([]scalar.Group), args[3].(lang.ColSet)

	if p, ok := input.Unwrap().(*Project); ok {
		// passthrough is the same as before, except we need to get
		// rid of the things that we're now computing outselves.
		newPassthrough := passthrough.Copy()
		var toInclude []int
		for i, c := range p.ColIDs {
			newPassthrough.Remove(c)
			if passthrough.Has(c) {
				toInclude = append(toInclude, i)
			}
		}

		// projections are also the same as before, but with:
		// * columns we were passing through before being computed now
		//   and
		// * input columns being inlined into our projections.
		newProjections := make([]scalar.Group, len(projections), len(projections)+len(toInclude))

		for i := range newProjections {
			newProjections[i] = inlineIn(m, projections[i], p.Projections, p.ColIDs)
		}

		newColIDs := make([]lang.ColumnID, len(colIDs), len(colIDs)+len(toInclude))
		copy(newColIDs, colIDs)
		for _, idx := range toInclude {
			newProjections = append(newProjections, p.Projections[idx])
			newColIDs = append(newColIDs, p.ColIDs[idx])
		}

		return m.Project(
			p.Input,
			newColIDs,
			newProjections,
			newPassthrough,
		)
	}

	return nil
}
