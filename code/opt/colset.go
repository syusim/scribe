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

func (c *ColSet) Len() int {
	return len(c.elems)
}

func (c *ColSet) Remove(col ColumnID) {
	c.ensure()
	delete(c.elems, col)
}

func (c *ColSet) Has(col ColumnID) bool {
	c.ensure()
	_, ok := c.elems[col]
	return ok
}

func (c *ColSet) ForEach(f func(c ColumnID)) {
	c.ensure()
	result := make([]ColumnID, 0)
	for c := range c.elems {
		result = append(result, c)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	for _, c := range result {
		f(c)
	}
}

func (c ColSet) SubsetOf(o ColSet) bool {
	c.ensure()
	o.ensure()
	for e := range c.elems {
		if _, ok := o.elems[e]; !ok {
			return false
		}
	}
	return true
}

func (c ColSet) Copy() ColSet {
	c.ensure()
	var res ColSet
	for e := range c.elems {
		res.Add(e)
	}
	return res
}

func (c ColSet) Equals(o ColSet) bool {
	c.ensure()
	o.ensure()
	return o.SubsetOf(c) && c.SubsetOf(o)
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
	c.ensure()
	// TODO: make this not suck!
	result := make([]ColumnID, 0)
	for c := range c.elems {
		result = append(result, c)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return fmt.Sprintf("%v", result)
}
