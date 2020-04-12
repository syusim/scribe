package opt

import "github.com/cockroachdb/cockroach/pkg/sql/opt"

// TODO: optimize this!
type ColSet struct {
	elems map[opt.ColumnID]struct{}
}

func MakeColSet() *ColSet {
	return &ColSet{
		elems: make(map[opt.ColumnID]struct{}),
	}
}

func (c *ColSet) Add(col opt.ColumnID) {
	c.elems[col] = struct{}{}
}

func (c *ColSet) Has(col opt.ColumnID) bool {
	_, ok := c.elems[col]
	return ok
}
