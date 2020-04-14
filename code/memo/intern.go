package memo

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
)

func hashField(buf *bytes.Buffer, f interface{}) {
	switch e := f.(type) {
	case string:
		buf.WriteString(e)
	case *RelExpr:
		fmt.Fprintf(buf, "%p", e)
	case opt.ColumnID:
		fmt.Fprintf(buf, "%d", e)
	case []opt.ColumnID:
		for i, c := range e {
			if i > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(buf, "%d", c)
		}
	case ScalarExpr:
		fmt.Fprintf(buf, "%p", e)
	case []ScalarExpr:
		for _, c := range e {
			fmt.Fprintf(buf, "%p", c)
		}
	case fmt.Stringer:
		buf.WriteString(e.String())
	default:
		panic(fmt.Sprintf("unhandled field type: %T", f))
	}
}

func hash(x interface{}) string {
	switch e := x.(type) {
	case lang.Datum:
		return e.String()
	default:
		var buf bytes.Buffer
		typ := reflect.TypeOf(x)
		buf.WriteString(typ.Name())
		v := reflect.ValueOf(x)
		for i, n := 0, v.NumField(); i < n; i++ {
			buf.WriteByte('/')
			hashField(&buf, v.Field(i).Interface())
		}
		return buf.String()
	}
}

// TODO: codegen these
func (m *Memo) internColRef(x ColRef) *ColRef {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*ColRef)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internConstant(x Constant) *Constant {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*Constant)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internAnd(x And) *And {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*And)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internPlus(x Plus) *Plus {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*Plus)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internFunc(x Func) *Func {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*Func)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internScan(x Scan) *RelExpr {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*RelExpr)
	}
	p := &RelExpr{E: &x}
	buildProps(p)
	m.hashes[h] = p
	return p
}

func (m *Memo) internSelect(x Select) *RelExpr {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*RelExpr)
	}
	p := &RelExpr{E: &x}
	buildProps(p)
	m.hashes[h] = p
	return p
}

func (m *Memo) internProject(x Project) *RelExpr {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*RelExpr)
	}
	p := &RelExpr{E: &x}
	buildProps(p)
	m.hashes[h] = p
	return p
}

func (m *Memo) internJoin(x Join) *RelExpr {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*RelExpr)
	}
	p := &RelExpr{E: &x}
	buildProps(p)
	m.hashes[h] = p
	return p
}
