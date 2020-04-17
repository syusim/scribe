package memo

import (
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

type rule func(m *Memo, args []interface{}) lang.Expr

// Plus Rules.

func FoldZeroPlus(m *Memo, args []interface{}) lang.Expr {
	if eqConst(args[0].(scalar.Expr), lang.DInt(0)) {
		return args[1].(lang.Expr)
	}
	return nil
}

func FoldPlusZero(m *Memo, args []interface{}) lang.Expr {
	if eqConst(args[1].(scalar.Expr), lang.DInt(0)) {
		return args[0].(lang.Expr)
	}
	return nil
}

func AssociatePlus(m *Memo, args []interface{}) lang.Expr {
	if left, ok := args[0].(*scalar.Plus); ok {
		return m.Plus(
			left.Left,
			m.Plus(
				left.Right,
				args[1].(scalar.Expr),
			),
		)
	}
	return nil
}

func SimplifyPlusPlus(m *Memo, args []interface{}) lang.Expr {
	left, right := args[0].(scalar.Expr), args[1].(scalar.Expr)
	if left == right {
		return m.Times(
			m.Constant(lang.DInt(2)),
			left,
		)
	}
	return nil
}

func SimplifyPlusTimes(m *Memo, args []interface{}) lang.Expr {
	left, right := args[0].(scalar.Expr), args[1].(scalar.Expr)
	if t, ok := right.(*scalar.Times); ok {
		if t.Right == left {
			if c, ok := t.Left.(*scalar.Constant); ok {
				return m.Times(
					m.Constant(lang.DInt(c.D.(lang.DInt)+1)),
					t.Right,
				)
			}
		}
	}
	return nil
}

// Join Rules.

func WrapJoinConditionInFilters(m *Memo, args []interface{}) lang.Expr {
	left, right, on := args[0].(*RelExpr), args[1].(*RelExpr), args[2].(scalar.Expr)
	if _, ok := on.(*scalar.Filters); !ok {
		return m.Join(
			left, right, m.Filters([]scalar.Expr{on}),
		)
	}
	return nil
}

func UnfoldJoinCondition(m *Memo, args []interface{}) lang.Expr {
	left, right, on := args[0].(*RelExpr), args[1].(*RelExpr), args[2].(*scalar.Filters)
	newFilters := unfoldFilters(on.Filters)

	if newFilters != nil {
		return m.Join(left, right, m.Filters(newFilters))
	}

	return nil
}

func PushFilterIntoJoinLeft(m *Memo, args []interface{}) lang.Expr {
	left, right, on := args[0].(*RelExpr), args[1].(*RelExpr), args[2].(*scalar.Filters)
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

func PushFilterIntoJoinRight(m *Memo, args []interface{}) lang.Expr {
	left, right, on := args[0].(*RelExpr), args[1].(*RelExpr), args[2].(*scalar.Filters)
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

func WrapSelectConditionInFilters(m *Memo, args []interface{}) lang.Expr {
	input, filter := args[0].(*RelExpr), args[1].(scalar.Expr)
	if _, ok := filter.(*scalar.Filters); !ok {
		return m.Select(
			input, m.Filters([]scalar.Expr{filter}),
		)
	}
	return nil
}

func UnfoldSelectCondition(m *Memo, args []interface{}) lang.Expr {
	input, filter := args[0].(*RelExpr), args[1].(*scalar.Filters)
	newFilters := unfoldFilters(filter.Filters)

	if newFilters != nil {
		return m.Select(input, m.Filters(newFilters))
	}

	return nil
}

func SimplifySelectFilters(m *Memo, args []interface{}) lang.Expr {
	input, filter := args[0].(*RelExpr), args[1].(*scalar.Filters)

	var newFilters []scalar.Expr
	for i, f := range filter.Filters {
		if eqConst(f, lang.DBool(true)) {
			newFilters = make([]scalar.Expr, i)
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

func EliminateSelect(m *Memo, args []interface{}) lang.Expr {
	input, filter := args[0].(*RelExpr), args[1].(*scalar.Filters)
	if len(filter.Filters) == 0 {
		return input
	}

	return nil
}

// Project Rules.

func EliminateProject(m *Memo, args []interface{}) lang.Expr {
	input, _, projections, passthrough := args[0].(*RelExpr), args[1].([]opt.ColumnID), args[2].([]scalar.Expr), args[3].(opt.ColSet)

	if len(projections) > 0 {
		return nil
	}

	if input.Props.OutputCols.Equals(passthrough) {
		return input
	}

	return nil
}

func MergeProjectProject(m *Memo, args []interface{}) lang.Expr {
	input, colIDs, projections, passthrough := args[0].(*RelExpr), args[1].([]opt.ColumnID), args[2].([]scalar.Expr), args[3].(opt.ColSet)

	if p, ok := input.E.(*Project); ok {
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
		newProjections := make([]scalar.Expr, len(projections), len(projections)+len(toInclude))

		for i := range newProjections {
			newProjections[i] = inlineIn(m, projections[i], p.Projections, p.ColIDs)
		}

		newColIDs := make([]opt.ColumnID, len(colIDs), len(colIDs)+len(toInclude))
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
