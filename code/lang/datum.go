package lang

import (
	"bytes"
	"fmt"
)

type Datum interface {
	Format(buf *bytes.Buffer)

	// To meet the exec.ScalarExpr interface.
	Eval(binding Row) (Datum, error)
}

type DInt int

func (d DInt) Format(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "%d", d)
}

func (d DInt) Eval(_ Row) (Datum, error) {
	return d, nil
}

type DString string

func (d DString) Format(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "%q", string(d))
}

func (d DString) Eval(_ Row) (Datum, error) {
	return d, nil
}

type DBool bool

func (d DBool) Format(buf *bytes.Buffer) {
	if d {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
}

func (d DBool) Eval(_ Row) (Datum, error) {
	return d, nil
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
