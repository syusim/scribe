package memo

import "github.com/justinj/scribe/code/opt"

func (m *Memo) Scan(tableName string, cols []opt.ColumnID) RelExpr {
	return RelExpr{
		&Scan{
			TableName: tableName,
			Cols:      cols,
		},
	}
}

func (m *Memo) Join(left, right RelExpr, on ScalarExpr) RelExpr {
	return RelExpr{
		&Join{
			Left:  left,
			Right: right,
			On:    on,
		},
	}
}

// TODO: standardize on xxxIDs vs. xxxIds
func (m *Memo) Project(
	input RelExpr,
	colIDs []opt.ColumnID,
	projections []ScalarExpr,
) RelExpr {
	return RelExpr{
		&Project{
			Input:       input,
			ColIDs:      colIDs,
			Projections: projections,
		},
	}
}

func (m *Memo) Select(input RelExpr, filter ScalarExpr) RelExpr {
	return RelExpr{
		&Select{
			Input:  input,
			Filter: filter,
		},
	}
}
