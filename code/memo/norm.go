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

func (m *Memo) Join(left, right *RelExpr, on []scalar.Expr) *RelExpr {
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
) *RelExpr {
	return m.internProject(Project{
		Input:       input,
		ColIDs:      colIDs,
		Projections: projections,
	})
}

func (m *Memo) Select(input *RelExpr, filter []scalar.Expr) *RelExpr {
	// MergeSelectJoin
	if j, ok := input.E.(*Join); ok {
		newFilter := make([]scalar.Expr, len(filter)+len(j.On))
		for i := range filter {
			newFilter[i] = filter[i]
		}
		for i := range j.On {
			newFilter[i+len(filter)] = j.On[i]
		}
		return m.Join(
			j.Left,
			j.Right,
			newFilter,
		)
	}

	return m.internSelect(Select{
		Input:  input,
		Filter: filter,
	})
}

func (m *Memo) Constant(d lang.Datum) scalar.Expr {
	return m.internConstant(scalar.Constant{d})
}

func (m *Memo) ColRef(id opt.ColumnID, typ lang.Type) scalar.Expr {
	return m.internColRef(scalar.ColRef{id, typ})
}

func (m *Memo) Plus(left, right scalar.Expr) scalar.Expr {
	// FoldZeroPlus
	if eqConst(left, lang.DInt(0)) {
		return right
	}

	// FoldPlusZero
	if eqConst(right, lang.DInt(0)) {
		return left
	}

	// AssociatePlus
	if l, ok := left.(*scalar.Plus); ok {
		return m.Plus(
			l.Left,
			m.Plus(
				l.Right,
				right,
			),
		)
	}

	return m.internPlus(scalar.Plus{left, right})
}

func (m *Memo) And(left, right scalar.Expr) scalar.Expr {
	// AssociateAnd
	if l, ok := left.(*scalar.And); ok {
		return m.And(
			l.Left,
			m.And(
				l.Right,
				right,
			),
		)
	}

	return m.internAnd(scalar.And{left, right})
}

func (m *Memo) Func(op lang.Func, args []scalar.Expr) scalar.Expr {
	return m.internFunc(scalar.Func{op, args})
}
