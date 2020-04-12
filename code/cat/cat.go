package cat

import (
	"fmt"

	"github.com/justinj/scribe/code/index"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
)

type Index struct {
	name     string
	ordering []opt.ColOrdinal
	data     *index.T
}

func (i *Index) Scan() *index.Iterator {
	return i.data.Iter()
}

func (i *Index) ScanGE(key lang.Key) *index.Iterator {
	return i.data.SeekGE(key)
}

type Table struct {
	Name    string
	cols    []lang.Column
	indexes []Index
}

func (t *Table) ColumnCount() int {
	return len(t.cols)
}

func (t *Table) Column(i int) *lang.Column {
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

func (c *Catalog) TableByName(s string) (*Table, bool) {
	for i := range c.tables {
		if c.tables[i].Name == s {
			return &c.tables[i], true
		}
	}
	return nil, false
}

func New() *Catalog {
	return &Catalog{}
}

func (c *Catalog) AddTable(
	name string,
	cols []lang.Column,
	data []lang.Row,
	// TODO: have a way for these to have names.
	indexes [][]opt.ColOrdinal,
) {
	tab := Table{
		Name: name,
		cols: cols,
	}

	tab.indexes = make([]Index, len(indexes))
	for i := range indexes {
		// TODO: use the names of the relevant columns?
		tab.indexes[i].name = fmt.Sprintf("%s_idx_%d", name, i+1)
		tab.indexes[i].ordering = indexes[i]
		tab.indexes[i].data = index.New(data, indexes[i])
	}

	c.tables = append(c.tables, tab)
}
