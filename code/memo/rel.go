package memo

import (
	"github.com/justinj/scribe/code/constraint"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

// This file implements the various relational logical operators.

type Scan struct {
	TableName string
	Cols      []opt.ColumnID

	// Index is the ordinal of the index in the given table.
	Index      int
	Constraint constraint.Constraint
}

func (s *Scan) ChildCount() int {
	return 0
}

func (s *Scan) Child(i int) lang.Group {
	panic("no children")
}

type Join struct {
	Left  *RelGroup
	Right *RelGroup
	On    scalar.Group
}

func (j *Join) ChildCount() int {
	return 3
}

func (j *Join) Child(i int) lang.Group {
	switch i {
	case 0:
		return j.Left
	case 1:
		return j.Right
	case 2:
		return j.On
	default:
		panic("out of bounds")
	}
}

type Project struct {
	Input *RelGroup

	ColIDs          []opt.ColumnID
	Projections     []scalar.Group
	PassthroughCols opt.ColSet
}

func (p *Project) ChildCount() int {
	return 1 + len(p.Projections)
}

func (p *Project) Child(i int) lang.Group {
	switch i {
	case 0:
		return p.Input
	default:
		return p.Projections[i-1]
	}
}

type Select struct {
	Input *RelGroup
	// TODO: unify terminology here: is it filter or predicate?
	Filter scalar.Group
}

func (s *Select) ChildCount() int {
	return 2
}

func (s *Select) Child(i int) lang.Group {
	switch i {
	case 0:
		return s.Input
	case 1:
		return s.Filter
	default:
		panic("out of bounds")
	}
}

// TODO: Can we kill Root and JUST have sort? Order-by introduces a sort at the top which maybe gets eliminated?
type Root struct {
	Input *RelGroup

	Ordering opt.Ordering
}

func (r *Root) ChildCount() int {
	return 1
}

func (r *Root) Child(i int) lang.Group {
	switch i {
	case 0:
		return r.Input
	default:
		panic("out of bounds")
	}
}

type Sort struct {
	Input *RelGroup

	Ordering opt.Ordering
}

func (r *Sort) ChildCount() int {
	return 1
}

func (r *Sort) Child(i int) lang.Group {
	switch i {
	case 0:
		return r.Input
	default:
		panic("out of bounds")
	}
}
