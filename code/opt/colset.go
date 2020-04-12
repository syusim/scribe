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

func (c *ColSet) Add(col ColumnID) {
	c.elems[col] = struct{}{}
}

func (c *ColSet) Has(col ColumnID) bool {
	_, ok := c.elems[col]
	return ok
}
