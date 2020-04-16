package opt

// TODO: optimize this!
type ColSet struct {
	elems map[ColumnID]struct{}
}

// TODO: get rid of this, use zero value as empty
func MakeColSet() *ColSet {
	return &ColSet{
		elems: make(map[ColumnID]struct{}),
	}
}

func SetFromCols(cols ...ColumnID) *ColSet {
	s := MakeColSet()
	for _, col := range cols {
		s.Add(col)
	}
	return s
}

func (c *ColSet) Add(col ColumnID) {
	if c.elems == nil {
		c.elems = make(map[ColumnID]struct{})
	}
	c.elems[col] = struct{}{}
}

func (c *ColSet) Has(col ColumnID) bool {
	_, ok := c.elems[col]
	return ok
}

func (c *ColSet) ForEach(f func(c ColumnID)) {
	for k := range c.elems {
		f(k)
	}
}

func (c *ColSet) SubsetOf(o ColSet) bool {
	for e := range c.elems {
		if _, ok := o.elems[e]; !ok {
			return false
		}
	}
	return true
}

func (c *ColSet) UnionWith(o *ColSet) {
	if o == nil || o.elems == nil {
		return
	}
	if c.elems == nil {
		c.elems = make(map[ColumnID]struct{})
	}
	for e := range o.elems {
		c.Add(e)
	}
}
