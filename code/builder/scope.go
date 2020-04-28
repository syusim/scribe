package builder

import (
	"github.com/justinj/scribe/code/lang"
)

// TODO: I guess this needs the table names as well.
type col struct {
	name string
	id   lang.ColumnID

	// this seems bad: duplicated, but convenient.
	// TODO: move this to be only stored in one place
	typ lang.Type
}

type scope struct {
	cols []col
}

func newScope() *scope {
	return &scope{
		cols: make([]col, 0),
	}
}

func (s *scope) addCol(name string, id lang.ColumnID, typ lang.Type) {
	s.cols = append(s.cols, col{name, id, typ})
}

func (s *scope) resolve(name string) (lang.ColumnID, lang.Type, bool) {
	for _, c := range s.cols {
		if c.name == name {
			return c.id, c.typ, true
		}
	}
	return 0, 0, false
}

func appendScopes(a, b *scope) *scope {
	// TODO: make this abstract?
	return &scope{
		cols: append(append(make([]col, 0), a.cols...), b.cols...),
	}
}

func (s *scope) OutCols() []lang.ColumnID {
	cols := make([]lang.ColumnID, len(s.cols))
	for i := range s.cols {
		cols[i] = s.cols[i].id
	}
	return cols
}
