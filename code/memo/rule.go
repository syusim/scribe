package memo

import (
	"github.com/justinj/scribe/code/lang"
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
