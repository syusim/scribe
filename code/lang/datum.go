package lang

import (
	"bytes"
	"fmt"
)

type Datum interface {
	fmt.Stringer

	Format(buf *bytes.Buffer)

	// this should actually be the expr thing
	Type() Type

	// To meet the scalar.Group interface.
	Eval(binding Row) (Datum, error)
}

type DInt int

func (d DInt) Format(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "%d", d)
}

func (d DInt) String() string {
	var buf bytes.Buffer
	d.Format(&buf)
	return buf.String()
}

func (d DInt) Child(i int) Group {
	panic("no children")
}

func (d DInt) ChildCount() int {
	return 0
}

func (d DInt) MemberCount() int {
	return 1
}

func (d DInt) Member(i int) Expr {
	return d
}

func (d DInt) Eval(_ Row) (Datum, error) {
	return d, nil
}

func (d DInt) Type() Type {
	return Int
}

type DString string

func (d DString) Format(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "%q", string(d))
}

func (d DString) String() string {
	var buf bytes.Buffer
	d.Format(&buf)
	return buf.String()
}

func (d DString) Child(i int) Group {
	panic("no children")
}

func (d DString) ChildCount() int {
	return 0
}

func (d DString) MemberCount() int {
	return 1
}

func (d DString) Member(i int) Expr {
	return d
}

func (d DString) Eval(_ Row) (Datum, error) {
	return d, nil
}

func (d DString) Type() Type {
	return String
}

type DBool bool

func (d DBool) Format(buf *bytes.Buffer) {
	if d {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
}

func (d DBool) String() string {
	var buf bytes.Buffer
	d.Format(&buf)
	return buf.String()
}

func (d DBool) Child(i int) Group {
	panic("no children")
}

func (d DBool) ChildCount() int {
	return 0
}

func (d DBool) MemberCount() int {
	return 1
}

func (d DBool) Member(i int) Expr {
	return d
}

func (d DBool) Eval(_ Row) (Datum, error) {
	return d, nil
}

func (d DBool) Type() Type {
	return Bool
}

type CmpResult int

const (
	_ CmpResult = iota
	LT
	EQ
	GT
)

func incompatible(a, b Datum) {
	panic(fmt.Sprintf("cannot compare %T with %T", a, b))
}

func Compare(a, b Datum) CmpResult {
	switch x := a.(type) {
	case DInt:
		y, ok := b.(DInt)
		if !ok {
			incompatible(a, b)
		}
		if x < y {
			return LT
		} else if x == y {
			return EQ
		} else {
			return GT
		}
	case DString:
		y, ok := b.(DString)
		if !ok {
			incompatible(a, b)
		}
		if x < y {
			return LT
		} else if x == y {
			return EQ
		} else {
			return GT
		}
	case DBool:
		// false < true
		y, ok := b.(DBool)
		if !ok {
			incompatible(a, b)
		}
		if x == y {
			return EQ
		} else if !x {
			return LT
		} else {
			return GT
		}
	default:
		panic(fmt.Sprintf("Compare not implemented for %T", a))
	}
}
