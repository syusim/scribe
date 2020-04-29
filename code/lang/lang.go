package lang

import (
	"bytes"
	"fmt"
)

type Expr interface {
	ChildCount() int
	Child(i int) Group
}

type Group interface {
	MemberCount() int
	Member(i int) Expr
}

func Unwrap(g Group) Expr {
	return g.Member(0)
}

type Type int

const (
	_ Type = iota
	Int
	String
	Bool
)

func (t Type) Format(buf *bytes.Buffer) {
	switch t {
	case Int:
		buf.WriteString("int")
	case String:
		buf.WriteString("string")
	case Bool:
		buf.WriteString("bool")
	}
}

func (t Type) String() string {
	var buf bytes.Buffer
	t.Format(&buf)
	return buf.String()
}

type Column struct {
	Name string
	Type Type
}

func (c *Column) Format(buf *bytes.Buffer) {
	buf.WriteByte('[')
	buf.WriteString(c.Name)
	buf.WriteByte(' ')
	c.Type.Format(buf)
	buf.WriteByte(']')
}

type Func int

const (
	_ Func = iota
	Eq
	Ne

	And
	Or
	Not

	Plus
	Minus
	Times
)

// NOTE: I think this stuff should be added in after.
func GetFunc(s string) (Func, error) {
	switch s {
	case "=":
		return Eq, nil
	case "!=":
		return Ne, nil
	case "and":
		return And, nil
	case "or":
		return Or, nil
	case "not":
		return Not, nil
	case "+":
		return Plus, nil
	case "-":
		return Minus, nil
	case "*":
		return Times, nil
	default:
		return 0, fmt.Errorf("unknown function %q", s)
	}
}

func (f Func) String() string {
	switch f {
	case Eq:
		return "="
	case Ne:
		return "!="
	case And:
		return "and"
	case Or:
		return "or"
	case Not:
		return "not"
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Times:
		return "*"
	}
	panic(fmt.Sprintf("unknown Func %d", f))
}

//(relational-types
type Relation struct {
	ColNames []string
	Rows     []Row
} //)

//(col-ordinal-type
type ColOrdinal int //)

type ColumnID int

type Ordering []ColumnID

func (o Ordering) Cols() ColSet {
	var cols ColSet
	for _, c := range o {
		cols.Add(c)
	}
	return cols
}

// TODO: merge the lang and opt packages
func RowCompare(a, b Row, ord []ColOrdinal) CmpResult {
	for _, idx := range ord {
		cmp := Compare(a[idx], b[idx])
		if cmp != EQ {
			return cmp
		}
	}
	return EQ
}

func RowCompare2(a, b Row, aOrd, bOrd []ColOrdinal) CmpResult {
	for i := range aOrd {
		cmp := Compare(a[aOrd[i]], b[bOrd[i]])
		if cmp != EQ {
			return cmp
		}
	}
	return EQ
}

func KeyCompare(a Row, b Key, ord []ColOrdinal) CmpResult {
	for i, idx := range ord {
		if i >= len(b) {
			return EQ
		}
		cmp := Compare(a[idx], b[i])
		if cmp != EQ {
			return cmp
		}
	}
	return EQ
}

////(relation.string
//func (t Relation) String() string {
//	widest := make([]int, len(t.ColNames))

//	for i, n := range t.ColNames {
//		if widest[i] < len(n) {
//			widest[i] = len(n)
//		}
//	}

//	for i := range t.Rows {
//		for j := range t.Rows[i] {
//			l := len(t.Rows[i][j])
//			if widest[j] < l {
//				widest[j] = l
//			}
//		}
//	}

//	var buf bytes.Buffer
//	for i, n := range t.ColNames {
//		if i > 0 {
//			buf.WriteString(" | ")
//		}
//		for k := 0; k < (widest[i]-len(n))/2; k++ {
//			buf.WriteByte(' ')
//		}
//		buf.WriteString(n)
//		for k := 0; k < (widest[i]-len(n))/2; k++ {
//			buf.WriteByte(' ')
//		}
//	}
//	buf.WriteByte('\n')
//	for i := range widest {
//		if i > 0 {
//			buf.WriteString("-+-")
//		}
//		for j := 0; j < widest[i]; j++ {
//			buf.WriteByte('-')
//		}
//	}
//	buf.WriteByte('\n')
//	for i := range t.Rows {
//		for j := range t.Rows[i] {
//			d := t.Rows[i][j]
//			if j > 0 {
//				buf.WriteString(" | ")
//			}
//			buf.WriteString(d)
//			for k := 0; k < widest[j]-len(d); k++ {
//				buf.WriteByte(' ')
//			}
//		}
//		buf.WriteByte('\n')
//	}

//	return buf.String()
//} //)
