package lang

import (
	"fmt"
)

type Type int

const (
	_ Type = iota
	Int
	String
	Bool
)

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