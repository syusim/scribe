package opt

// TODO: optimize this!
type ColMap struct {
	elems map[ColumnID]int
}

func (c *ColMap) Set(col ColumnID, i int) {
	if c.elems == nil {
		c.elems = make(map[ColumnID]int)
	}
	c.elems[col] = i
}

func (c *ColMap) Get(col ColumnID) (int, bool) {
	if c.elems == nil {
		return 0, false
	}
	i, ok := c.elems[col]
	return i, ok
}

func (c *ColMap) ForEach(f func(from ColumnID, to int)) {
	for k, v := range c.elems {
		f(k, v)
	}
}

func (c *ColMap) Len() int {
	return len(c.elems)
}
