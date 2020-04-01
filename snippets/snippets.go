package snippets

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"
)

type FlagSet map[string]struct{}

type Block interface {
	Render(*bytes.Buffer, FlagSet)
}

type pair struct {
	l Block
	r Block
}

func (b *pair) Render(buf *bytes.Buffer, f FlagSet) {
	b.l.Render(buf, f)
	b.r.Render(buf, f)
}

type literal struct {
	contents string
}

func (b *literal) Render(buf *bytes.Buffer, _ FlagSet) {
	buf.WriteString(b.contents)
}

type fence struct {
	tag      string
	contents Block
	anti     bool
}

func (b *fence) Render(buf *bytes.Buffer, f FlagSet) {
	if _, ok := f[b.tag]; ok != b.anti {
		b.contents.Render(buf, f)
	}
}

// TODO: just make this return a string, that's how it's used now.
func Extract(b Block, buf *bytes.Buffer, f FlagSet, tag string) {
	switch e := b.(type) {
	case *pair:
		Extract(e.l, buf, f, tag)
		Extract(e.r, buf, f, tag)
	case *fence:
		if e.tag == tag {
			// We found it! Render according to the flagSet.
			e.contents.Render(buf, f)
		} else {
			// This is not the block we were looking for. Keep digging.
			Extract(e.contents, buf, f, tag)
		}
	}
}

func Tags(b Block) (map[string]struct{}, error) {
	m := make(map[string]struct{})
	err := tags(b, m)
	return m, err
}

func tags(b Block, m map[string]struct{}) error {
	switch e := b.(type) {
	case *pair:
		if err := tags(e.l, m); err != nil {
			return err
		}
		if err := tags(e.r, m); err != nil {
			return err
		}
	case *fence:
		if _, ok := m[e.tag]; ok {
			return fmt.Errorf("duplicate tag: %q", e.tag)
		}
		m[e.tag] = struct{}{}
		return tags(e.contents, m)
	}
	return nil
}

func concat(l, r Block) Block {
	// We could maybe do this "optimization," if we're going to re-render a tree
	// a lot. Kind of annoyingly quadratic, though.
	// ll, ok1 := l.(*literal)
	// rl, ok2 := r.(*literal)
	// if ok1 && ok2 {
	// 	return &literal{
	// 		contents: ll.contents + rl.contents,
	// 	}
	// }

	return &pair{l, r}
}

type tokStream struct {
	buffered []token
	r        *bufio.Reader

	// level is the number of comments to strip off the beginning of each line.
	level int
}

type kind int

const (
	noKind kind = iota
	lineKind
	startBlockKind
	endBlockKind
	startAntiBlockKind
	endAntiBlockKind
)

type token struct {
	k   kind
	lex string
}

var openRegexp = regexp.MustCompile(`^\s*//\((\S+)\s*$`)

// These two cases could possibly be combined.
var closeRegexpAppended = regexp.MustCompile(`^(.*\S)\s*//\)\s*$`)
var closeRegexp = regexp.MustCompile(`^\s*//\)\s*$`)
var openAntiRegexp = regexp.MustCompile(`^\s*//\[(\S+)\s*$`)
var closeAntiRegexp = regexp.MustCompile(`^\s*\]\s*$`)

func (t *tokStream) populate() bool {
	if len(t.buffered) == 0 {
		line, err := t.r.ReadString('\n')
		if err != nil {
			return false
		}

		start := 0
		for start < len(line) && unicode.IsSpace(rune(line[start])) {
			start++
		}
		pre := line[:start]
		line = line[start:]
		for i := 0; i < t.level; i++ {
			if len(line) < 2 || line[:2] != "//" {
				panic("expected // indentation")
			}
			line = line[2:]
		}
		line = pre + line

		matches := openRegexp.FindStringSubmatch(line)
		if len(matches) > 0 {
			t.buffered = append(t.buffered, token{
				k:   startBlockKind,
				lex: matches[1],
			})
			return true
		}

		matches = closeRegexpAppended.FindStringSubmatch(line)
		if len(matches) > 0 {
			t.buffered = append(t.buffered,
				token{
					k:   lineKind,
					lex: strings.TrimRight(matches[1], " \t") + "\n",
				},
				token{
					k:   endBlockKind,
					lex: "",
				})
			return true
		}

		if closeRegexp.MatchString(line) {
			t.buffered = append(t.buffered, token{
				k:   endBlockKind,
				lex: "",
			})
			return true
		}

		matches = openAntiRegexp.FindStringSubmatch(line)
		if len(matches) > 0 {
			t.buffered = append(t.buffered, token{
				k:   startAntiBlockKind,
				lex: matches[1],
			})
			return true
		}

		if closeAntiRegexp.MatchString(line) {
			t.buffered = append(t.buffered, token{k: endAntiBlockKind})
			return true
		}

		t.buffered = append(t.buffered, token{
			k:   lineKind,
			lex: line,
		})
		return true
	}

	return true
}

func (t *tokStream) Next() (tok token, ok bool) {
	if !t.populate() {
		return token{}, false
	}
	var n token
	n, t.buffered = t.buffered[0], t.buffered[1:]
	return n, true
}

func (t *tokStream) Peek() (tok token, ok bool) {
	if !t.populate() {
		return token{}, false
	}
	return t.buffered[0], true
}

func closer(in kind) kind {
	switch in {
	case startBlockKind:
		return endBlockKind
	case startAntiBlockKind:
		return endAntiBlockKind
	}
	panic("no closer")
}

func build(in *tokStream) (Block, error) {
	tok, ok := in.Next()
	if !ok {
		return &literal{""}, io.EOF
	}
	switch tok.k {
	case lineKind:
		return &literal{tok.lex}, nil
	case startBlockKind, startAntiBlockKind:
		if tok.k == startAntiBlockKind {
			// We need to strip off a layer of comments at the beginning.
			in.level++
		}
		var result Block
		result = &literal{""}

		r, _ := in.Peek()
		if r.k == closer(tok.k) {
			_, _ = in.Next()
			return result, nil
		}

		for n, err := build(in); err == nil; n, err = build(in) {
			result = concat(result, n)
			r, _ := in.Peek()
			if r.k == closer(tok.k) {
				in.level--
				// Skip over it.
				_, _ = in.Next()
				return &fence{
					contents: result,
					tag:      tok.lex,
					anti:     tok.k == startAntiBlockKind,
				}, nil
			}
		}
	case endBlockKind, endAntiBlockKind:
		panic(fmt.Sprintf("hit end block: %v", tok))
	}

	panic(fmt.Sprintf("no good chief: %v", tok.k))
}

func New(in io.Reader) (Block, error) {
	r := &tokStream{r: bufio.NewReader(in)}
	var result Block
	result = &literal{""}

	for n, err := build(r); err == nil; n, err = build(r) {
		result = concat(result, n)
	}
	return result, nil
}
