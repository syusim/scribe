package fd

import (
	"bytes"
	"fmt"

	"github.com/justinj/scribe/code/lang"
)

type FD struct {
	Lhs lang.ColSet
	Rhs lang.ColumnID
}

func (fd *FD) String() string {
	return fmt.Sprintf("%s->%d", &fd.Lhs, fd.Rhs)
}

type FDSet struct {
	fds []FD
}

var None = &FDSet{}

func (fd *FDSet) cpy() *FDSet {
	newFDs := make([]FD, len(fd.fds))
	copy(newFDs, fd.fds)
	return &FDSet{newFDs}
}

func (fd *FDSet) Imply(fds []FD) *FDSet {
	newFDs := fd.cpy()
	newFDs.fds = append(newFDs.fds, fds...)
	return newFDs
}

func (fd *FDSet) Implies(lhs lang.ColSet, rhs lang.ColSet) bool {
	return rhs.SubsetOf(fd.closure(lhs))
}

func (fd *FDSet) closure(s lang.ColSet) lang.ColSet {
	result := s.Copy()
	for i := 0; i < len(fd.fds); i++ {
		e := fd.fds[i]
		if e.Lhs.SubsetOf(result) {
			if !result.Has(e.Rhs) {
				result.Add(e.Rhs)
				i = -1
			}
		}
	}
	return result
}

func (fd *FDSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('(')
	for i, f := range fd.fds {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(f.String())
	}
	buf.WriteByte(')')
	return buf.String()
}
