package cat

import (
	"fmt"

	"github.com/justinj/scribe/code/ast"
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

func (c *Catalog) AddTable(def *ast.CreateTable) error {
	// TODO: validate this:
	// * don't allow repeated index names.
	tab := Table{
		Name: def.Name,
		cols: def.Columns,
	}

	idxs := def.Indexes

	if len(idxs) == 0 {
		idxs = []ast.IndexDef{
			{Name: "default"},
		}
	}

	tab.indexes = make([]Index, len(idxs))
	for i, idx := range idxs {
		tab.indexes[i].name = fmt.Sprintf(idx.Name)

		ords := make([]opt.ColOrdinal, len(idx.Cols))
		for j, idxCol := range idx.Cols {
			nextOrd := -1

			// TODO: use a better algorithm here (or don't. it's a free country)
			for k, c := range def.Columns {
				if c.Name == idxCol {
					nextOrd = k
					break
				}
			}

			if nextOrd == -1 {
				return fmt.Errorf("invalid index column %q", idxCol)
			}

			ords[j] = opt.ColOrdinal(nextOrd)
		}

		tab.indexes[i].ordering = ords
		tab.indexes[i].data = index.New(def.Data, ords)
	}

	c.tables = append(c.tables, tab)

	return nil
}
