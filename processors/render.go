package processors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"

	"github.com/justinj/scribe/snippets"
	"github.com/russross/blackfriday"
)

func span(buf *bytes.Buffer, classes string, f func()) {
	fmt.Fprintf(buf, `<span class="%s">`, classes)
	f()
	buf.WriteString(`</span>`)
}

type docMeta struct {
	Title string
}

func Build(
	in io.Reader,
	snippetDir string,
) (io.Reader, error) {
	c, err := buildCorpus(snippetDir)
	if err != nil {
		return nil, err
	}

	var meta docMeta

	h := newHighlighter()

	seenFlags := make(snippets.FlagSet)

	var buf bytes.Buffer

	for _, tok := range tokenize(in) {
		switch tok.kind {
		case metaLine:
			json.Unmarshal([]byte(tok.lex), &meta)
		case bareLine:
			buf.Write(
				blackfriday.MarkdownCommon([]byte(tok.lex)),
			)
		case snippetRefLine:
			seenFlags[tok.lex] = struct{}{}
			pre, mid, post := c.getSnip(seenFlags, tok.lex)

			// I'm sure some CSS nerd will tell me this is a bad use of span
			span(&buf, "greyout top-code", func() {
				h.highlight(&buf, pre)
			})

			h.highlight(&buf, mid)

			span(&buf, "greyout bottom-code", func() {
				h.highlight(&buf, post)
			})
		}
	}

	return renderHTML(tmplData{
		Css:     template.CSS(h.css()),
		Content: template.HTML(buf.String()),
	})
}

type tmplData struct {
	Css     template.CSS
	Content template.HTML
}

func renderHTML(d tmplData) (io.Reader, error) {
	t, err := ioutil.ReadFile("templates/page.html")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("post").Parse(string(t))
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer

	tmpl.Execute(&out, d)

	return &out, nil
}
