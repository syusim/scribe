package lang

import (
	"bytes"
	"fmt"
)

type Datum interface {
}

type DInt int

func (d DInt) Format(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "%d", d)
}

type DString string

func (d DString) Format(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "%q", string(d))
}

type DBool bool

func (d DBool) Format(buf *bytes.Buffer) {
	if d {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
}
