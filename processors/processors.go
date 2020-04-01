package processors

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/justinj/scribe/snippets"
	"github.com/russross/blackfriday"
)

type corpus struct {
	files map[string]snippets.Block
	tags  map[string]string
}

func buildCorpus(dir string) (*corpus, error) {
	c := &corpus{
		files: make(map[string]snippets.Block),
		tags:  make(map[string]string),
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		b, err := snippets.New(f)
		if err != nil {
			return err
		}
		tags, err := snippets.Tags(b)
		if err != nil {
			return err
		}

		c.files[path] = b
		for tag, _ := range tags {
			if _, ok := c.tags[tag]; ok {
				return fmt.Errorf("duplicate tag: %q", tag)
			}
			c.tags[tag] = path
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}

type lineKind int

const (
	_ lineKind = iota
	bareLine
	snippetRefLine
)

type line struct {
	kind lineKind
	lex  string
}

func tokenize(in io.Reader) []line {
	result := make([]line, 0)
	r := bufio.NewReader(in)

	var buf bytes.Buffer

	flush := func() {
		if buf.Len() > 0 {
			result = append(result, line{
				kind: bareLine,
				lex:  buf.String(),
			})
			buf.Reset()
		}
	}

	for l, err := r.ReadString('\n'); err == nil; l, err = r.ReadString('\n') {
		if strings.HasPrefix(l, "% ") {
			flush()
			result = append(result, line{
				kind: snippetRefLine,
				lex:  strings.Trim(l[2:], " \n"),
			})
		} else {
			buf.WriteString(l)
		}
	}

	flush()

	return result
}

func unindent(s string) string {
	smallestIndent := len(s)
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		for i := 0; i < len(l); i++ {
			if !unicode.IsSpace(rune(l[i])) {
				if smallestIndent > i {
					smallestIndent = i
					break
				}
			}
		}
		if smallestIndent == 0 {
			return s
		}
	}

	var out bytes.Buffer
	for i, l := range lines {
		if len(l) <= smallestIndent {
			// l is probably just a blank line.
			out.WriteString(l)
		} else {
			out.WriteString(l[smallestIndent:])
		}
		if i < len(lines)-1 {
			out.WriteByte('\n')
		}
	}

	return out.String()
}

func Build(
	in io.Reader,
	snippetDir string,
) (io.Reader, error) {
	c, err := buildCorpus(snippetDir)
	if err != nil {
		return nil, err
	}
	toks := tokenize(in)

	// CHROMA STUFF
	// TODO: wrap this so it's nicer to use for my purposes
	lexer := lexers.Get("go")
	style := styles.Get("monokai")
	formatter := html.New(html.WithClasses(true))

	var buf bytes.Buffer

	buf.WriteString("<style>")
	buf.WriteString("html { tab-size: 4; -moz-tab-size: 4; }")
	if err := formatter.WriteCSS(&buf, style); err != nil {
		return nil, err
	}
	buf.WriteString("pre { margin-top: 0; margin-bottom: 0; }")
	buf.WriteString(".greyout { filter: brightness(50%); }")
	buf.WriteString("</style>")

	seenFlags := make(snippets.FlagSet)

	for _, tok := range toks {
		switch tok.kind {
		case bareLine:
			buf.Write(
				blackfriday.MarkdownBasic([]byte(tok.lex)),
			)
		case snippetRefLine:
			referenced := tok.lex
			seenFlags[tok.lex] = struct{}{}
			referencedFile, ok := c.tags[referenced]
			if !ok {
				panic(fmt.Sprintf("bad snippet: %q", referenced))
			}
			file := c.files[referencedFile]
			var b bytes.Buffer
			extracted := snippets.ExtractCtx(file, seenFlags, referenced)
			b.WriteString(extracted.Pre)
			b.WriteString("++++\n")
			b.WriteString(extracted.Contents)
			b.WriteString("++++\n")
			b.WriteString(extracted.Post)
			code := unindent(b.String())
			// TODO: ugh, we need to unindent this properly
			sections := strings.Split(code, "++++\n")
			pre, mid, post := sections[0], sections[1], sections[2]

			it, err := lexer.Tokenise(nil, pre)
			if err != nil {
				return nil, err
			}

			// I'm sure some CSS nerd will tell me this is a bad use of span
			buf.WriteString(`<span class='greyout'>`)
			err = formatter.Format(&buf, style, it)
			buf.WriteString(`</span>`)

			it, err = lexer.Tokenise(nil, mid)
			if err != nil {
				return nil, err
			}

			err = formatter.Format(&buf, style, it)

			it, err = lexer.Tokenise(nil, post)
			if err != nil {
				return nil, err
			}

			buf.WriteString(`<span class='greyout'>`)
			err = formatter.Format(&buf, style, it)
			buf.WriteString(`</span>`)
		}
	}

	return &buf, nil
}

func main() {
	fmt.Println("vim-go")
}
