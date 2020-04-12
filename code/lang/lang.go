package lang

import (
	"bytes"
	"fmt"
)

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
