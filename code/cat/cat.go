package cat

import (
	"fmt"

	"github.com/justinj/scribe/code/index"
	"github.com/justinj/scribe/code/opt"
)

type Column struct {
	Name string
	Type opt.Type
}

type Index struct {
	name     string
	ordering []opt.ColOrdinal
	data     *index.T
}

func (i *Index) Scan(key opt.Key) *index.Iterator {
	return i.data.SeekGE(key)
}

type Table struct {
	Name    string
	cols    []Column
	indexes []Index
}

func (t *Table) ColumnCount() int {
	return len(t.cols)
}

func (t *Table) Column(i int) *Column {
	return &t.cols[i]
}

func (t *Table) IndexCount() int {
	return len(t.indexes)
}

func (t *Table) Index(i int) *Index {
	return &t.indexes[i]
}

type Catalog struct {
	tables []Table
}

func (c *Catalog) TableCount() int {
	return len(c.tables)
}

func (c *Catalog) Table(i int) *Table {
	return &c.tables[i]
}

func New() *Catalog {
	return &Catalog{}
}

func (c *Catalog) AddTable(
	name string,
	cols []Column,
	data opt.Relation,
	// TODO: have a way for these to have names.
	indexes [][]opt.ColOrdinal,
) {
	tab := Table{
		Name: name,
		cols: cols,
	}

	idxs := make([]Index, len(indexes))
	for i := range indexes {
		// TODO: use the names of the relevant columns?
		idxs[i].name = fmt.Sprintf("%s_idx_%d", name, i+1)
		idxs[i].ordering = indexes[i]
		idxs[i].data = index.New(data, indexes[i])
	}

	c.tables = append(c.tables, tab)
}
