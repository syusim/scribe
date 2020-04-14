package memo

import (
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
)

func (m *Memo) Scan(tableName string, cols []opt.ColumnID) *RelExpr {
	return m.internScan(Scan{
		TableName: tableName,
		Cols:      cols,
	})
}

func (m *Memo) Join(left, right *RelExpr, on ScalarExpr) *RelExpr {
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
	projections []ScalarExpr,
) *RelExpr {
	return m.internProject(Project{
		Input:       input,
		ColIDs:      colIDs,
		Projections: projections,
	})
}

func (m *Memo) Select(input *RelExpr, filter ScalarExpr) *RelExpr {
	return m.internSelect(Select{
		Input:  input,
		Filter: filter,
	})
}

func (m *Memo) Constant(d lang.Datum) ScalarExpr {
	return m.internConstant(Constant{d})
}

func (m *Memo) ColRef(id opt.ColumnID, typ lang.Type) ScalarExpr {
	return m.internColRef(ColRef{id, typ})
}
