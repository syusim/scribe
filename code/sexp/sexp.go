package sexp

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
)

type Sexp interface {
	fmt.Stringer
}

type Atom string
type String string
type List []Sexp

func format(buf *bytes.Buffer, s Sexp) {
	switch a := s.(type) {
	case Atom:
		buf.WriteString(string(a))
	case String:
		fmt.Fprintf(buf, "%q", string(a))
	case List:
		buf.WriteByte('(')
		for i := range a {
			if i > 0 {
				buf.WriteByte(' ')
			}
			format(buf, a[i])
		}
		buf.WriteByte(')')
	}
}

func (a Atom) String() string {
	var buf bytes.Buffer
	format(&buf, a)
	return buf.String()
}

func (s String) String() string {
	var buf bytes.Buffer
	format(&buf, s)
	return buf.String()
}

func (l List) String() string {
	var buf bytes.Buffer
	format(&buf, l)
	return buf.String()
}

func Pretty(s Sexp) string {
	var buf bytes.Buffer
	depth := 0
	indent := func() {
		for i := 0; i < depth; i++ {
			buf.WriteByte(' ')
		}
	}

	var p func(s Sexp)
	p = func(s Sexp) {
		switch e := s.(type) {
		case List:
			depth++
			buf.WriteByte('(')
			for i := range e {
				if i > 0 {
					indent()
				}
				p(e[i])
				if i < len(e)-1 {
					buf.WriteByte('\n')
				}
			}
			buf.WriteByte(')')
			depth--
		default:
			buf.WriteString(s.String())
		}
	}

	p(s)
	return buf.String()
}

func isAtomChar(b byte) bool {
	switch b {
	case '(', ')', '[', ']', '{', '}':
		return false
	}
	return !unicode.IsSpace(rune(b))
}

func matching(c byte) byte {
	switch c {
	case '(':
		return ')'
	case '[':
		return ']'
	case '{':
		return '}'
	}

	panic("no closer")
}

func Parse(s string) (Sexp, error) {
	munch := func() {
		for len(s) > 0 && unicode.IsSpace(rune(s[0])) {
			s = s[1:]
		}
	}

	var next func() (Sexp, error)
	next = func() (Sexp, error) {
		switch s[0] {
		case '(', '[', '{':
			opener := s[0]
			closer := matching(s[0])
			s = s[1:]
			result := make(List, 0)
			munch()
			for len(s) > 0 && s[0] != closer {
				munch()
				n, err := next()
				if err != nil {
					return nil, err
				}
				result = append(result, n)
				munch()
				if len(s) == 0 {
					return nil, fmt.Errorf("unmatched '%c'", opener)
				}
				if s[0] == closer {
					s = s[1:]
					return result, nil
				}
			}
			// HACK: refactor this to be neater.
			if len(s) > 0 {
				s = s[1:]
			}
			return result, nil
		case ')', ']', '}':
			return nil, fmt.Errorf("unmatched '%c'", s[0])
		case '"':
			s = s[1:]
			i := 0
			for i < len(s) && s[i] != '"' {
				i++
				if s[i] == '\\' {
					if i+1 >= len(s) {
						panic("expected escape")
					}
					switch s[i+1] {
					case '"':
						i++
					default:
						panic("unknown escape sequence")
					}
				}
			}
			var r string
			r, s = s[:i], s[i+1:]
			return String(r), nil
		default:
			i := 0
			for i < len(s) && isAtomChar(s[i]) {
				i++
			}
			var r string
			r, s = s[:i], s[i:]
			return Atom(r), nil
		}
		return nil, fmt.Errorf("fail")
	}

	munch()
	return next()
}

func Nth(s Sexp, n int) Sexp {
	return s.(List)[n]
}

func Int(s Sexp) int {
	i, err := strconv.Atoi(string(s.(Atom)))
	if err != nil {
		panic(err)
	}
	return i
}
