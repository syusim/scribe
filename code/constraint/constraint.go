package constraint

import (
	"bytes"

	"github.com/justinj/scribe/code/cat"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/scalar"
)

type Constraint struct {
	InclusiveStart bool
	Start          lang.Key

	InclusiveEnd bool
	End          lang.Key
}

func (c Constraint) Format(buf *bytes.Buffer) {
	if c.InclusiveStart {
		buf.WriteByte('[')
	} else {
		buf.WriteByte('(')
	}
	if c.Start != nil {
		c.Start.Format(buf)
	}
	buf.WriteString(" - ")
	if c.End != nil {
		c.End.Format(buf)
	}
	if c.InclusiveEnd {
		buf.WriteByte(']')
	} else {
		buf.WriteByte(')')
	}
}

func (c Constraint) IsUnconstrained() bool {
	return c.Start == nil && c.End == nil
}

func Generate(
	filters *scalar.Filters,
	colIdxs []lang.ColumnID,
	index *cat.Index,
) (c Constraint, remaining []scalar.Group) {
	// super baby mode. only handle an eq on the first col.
	f := filters.Filters
	if len(f) != 1 {
		return Constraint{}, filters.Filters
	}

	if eq, ok := lang.Unwrap(f[0]).(*scalar.Eq); ok {
		// lhs gotta be a var and rhs gotta be a const
		if col, ok := lang.Unwrap(eq.Left).(*scalar.ColRef); ok {
			if val, ok := lang.Unwrap(eq.Right).(*scalar.Constant); ok {
				// Check if the col is the first col of the index.
				ord := -1
				for i := range colIdxs {
					if colIdxs[i] == col.Id {
						ord = i
						break
					}
				}
				if ord == -1 {
					panic("didn't find col (this shouldn't happen)")
				}

				if len(index.Ordering) == 0 {
					// Can't generate anything! Probably shoulda bailed sooner.
					return Constraint{}, filters.Filters
				}

				if index.Ordering[0] == lang.ColOrdinal(ord) {
					// Now we're COOKIN!
					return Constraint{
						InclusiveStart: true,
						InclusiveEnd:   true,
						Start:          lang.Key{val.D},
						End:            lang.Key{val.D},
					}, []scalar.Group{}
				}
			}
		}
	}

	return Constraint{}, filters.Filters
}
