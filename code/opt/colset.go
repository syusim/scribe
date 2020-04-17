package opt

import (
	"fmt"
	"sort"
)

// TODO: optimize this!
type ColSet struct {
	elems map[ColumnID]struct{}
}

func (c *ColSet) ensure() {
	if c.elems == nil {
		c.elems = make(map[ColumnID]struct{})
	}
}

func SetFromCols(cols ...ColumnID) ColSet {
	var s ColSet
	for _, col := range cols {
		s.Add(col)
	}
	return s
}

func (c *ColSet) Add(col ColumnID) {
	c.ensure()
	c.elems[col] = struct{}{}
}

func (c *ColSet) Has(col ColumnID) bool {
	c.ensure()
	_, ok := c.elems[col]
	return ok
}

func (c *ColSet) ForEach(f func(c ColumnID)) {
	c.ensure()
	for k := range c.elems {
		f(k)
	}
}

func (c *ColSet) SubsetOf(o ColSet) bool {
	c.ensure()
	o.ensure()
	for e := range c.elems {
		if _, ok := o.elems[e]; !ok {
			return false
		}
	}
	return true
}

func (c *ColSet) UnionWith(o ColSet) {
	if c == nil {
		return
	}
	c.ensure()
	o.ensure()
	if c.elems == nil {
		c.elems = make(map[ColumnID]struct{})
	}
	for e := range o.elems {
		c.Add(e)
	}
}

func (c *ColSet) String() string {
	// TODO: make this not suck!
	result := make([]ColumnID, 0)
	for c := range c.elems {
		result = append(result, c)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return fmt.Sprintf("%v", result)
}
