package memo

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/opt"
	"github.com/justinj/scribe/code/scalar"
)

// TODO: this is subtle broken: the separator needs to be escaped
// in all these instances.
func hashField(buf *bytes.Buffer, f interface{}) {
	switch e := f.(type) {
	case int:
		fmt.Fprintf(buf, "%d", e)
	case string:
		buf.WriteString(e)
	case *RelGroup:
		fmt.Fprintf(buf, "%p", e)
	case opt.ColumnID:
		fmt.Fprintf(buf, "%d", e)
	case opt.ColSet:
		// TODO: need to make a real thing here, but this is safe
		// because Go guarantees order in stringified form.
		fmt.Fprintf(buf, "%s", e.String())
	case opt.Ordering:
		for i, c := range e {
			if i > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(buf, "%d", c)
		}
	case []opt.ColumnID:
		for i, c := range e {
			if i > 0 {
				buf.WriteByte(',')
			}
			fmt.Fprintf(buf, "%d", c)
		}
	case scalar.Group:
		fmt.Fprintf(buf, "%p", e)
	case []scalar.Group:
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
func (m *Memo) internColRef(x scalar.ColRef) *scalar.ColRef {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*scalar.ColRef)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internConstant(x scalar.Constant) *scalar.Constant {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*scalar.Constant)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internAnd(x scalar.And) *scalar.And {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*scalar.And)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internFilters(x scalar.Filters) *scalar.Filters {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*scalar.Filters)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internPlus(x scalar.Plus) *scalar.Plus {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*scalar.Plus)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internTimes(x scalar.Times) *scalar.Times {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*scalar.Times)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internFunc(x scalar.Func) *scalar.Func {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*scalar.Func)
	}
	p := &x
	m.hashes[h] = p
	return p
}

func (m *Memo) internScan(x Scan) *RelGroup {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*RelGroup)
	}
	p := &RelGroup{E: &x}
	buildProps(p)
	m.hashes[h] = p
	return p
}

func (m *Memo) internSelect(x Select) *RelGroup {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*RelGroup)
	}
	p := &RelGroup{E: &x}
	buildProps(p)
	m.hashes[h] = p
	return p
}

func (m *Memo) internProject(x Project) *RelGroup {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*RelGroup)
	}
	p := &RelGroup{E: &x}
	buildProps(p)
	m.hashes[h] = p
	return p
}

func (m *Memo) internJoin(x Join) *RelGroup {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*RelGroup)
	}
	p := &RelGroup{E: &x}
	buildProps(p)
	m.hashes[h] = p
	return p
}

func (m *Memo) internRoot(x Root) *RelGroup {
	h := hash(x)
	if v, ok := m.hashes[h]; ok {
		return v.(*RelGroup)
	}
	p := &RelGroup{E: &x}
	buildProps(p)
	m.hashes[h] = p
	return p
}
