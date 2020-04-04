package compiler

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday"
)

type document struct {
	c *Compiler

	sourcePath string
	parsed     DocMeta
}

func (d document) Target() string {
	return fmt.Sprintf("%s/%s.html", d.c.outDir, d.sourcePath[len(d.c.inDir):])
}

type Compiler struct {
	outDir     string
	inDir      string
	snippetDir string

	documents []document
	h         *highlighter
	corpus    *corpus
	running   bool
}

func New(
	outDir string,
	inDir string,
	snippetDir string,
) *Compiler {
	cor, err := buildCorpus(snippetDir)
	if err != nil {
		panic(err)
	}

	return &Compiler{
		outDir:     outDir,
		inDir:      inDir,
		snippetDir: snippetDir,

		h:      newHighlighter(),
		corpus: cor,
	}
}

func (c *Compiler) Run() chan<- Event {
	if c.running {
		panic("double-run")
	}
	c.running = true

	ch := make(chan Event)
	go func() {
		for {
			switch e := (<-ch).(type) {
			case FileAdded:
				if err := c.AddFile(e.Path); err != nil {
					panic(err)
				}
			case stop:
				return
			}
		}
	}()

	return ch
}

func (c *Compiler) AddFile(s string) error {
	parsed, err := c.ParseDocument(s)
	if err != nil {
		return err
	}

	doc := document{
		c:          c,
		sourcePath: s,
		parsed:     parsed,
	}

	c.documents = append(c.documents, doc)

	if err := c.buildDoc(doc); err != nil {
		return err
	}

	return nil
}

func span(buf *bytes.Buffer, classes string, f func()) {
	fmt.Fprintf(buf, `<span class="%s">`, classes)
	f()
	buf.WriteString(`</span>`)
}

func (c *Compiler) pageHTML(d DocMeta) string {
	var buf bytes.Buffer

	for _, section := range d.Contents {
		switch s := section.(type) {
		case Prose:
			buf.Write(
				blackfriday.MarkdownCommon([]byte(s.Contents)),
			)
		case Code:
			// I'm sure some CSS nerd will tell me this is a bad use of span
			span(&buf, "greyout top-code", func() {
				c.h.highlight(&buf, s.Pre)
			})

			c.h.highlight(&buf, s.Code)

			span(&buf, "greyout bottom-code", func() {
				c.h.highlight(&buf, s.Post)
			})
		case Error:
			fmt.Fprintf(&buf, `<div class="error"><div class="head">Error:</div>%s</div>`, s.Msg)
		}
	}

	return buf.String()
}

func (c *Compiler) buildDoc(d document) error {
	out := d.Target()

	body := template.HTML(c.pageHTML(d.parsed))

	rendered, err := renderHTML(tmplData{
		Css:     template.CSS(c.h.css()),
		Content: body,
	})
	if err != nil {
		return err
	}

	if err := func() error {
		dir, _ := filepath.Split(out)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}

		f, err := os.Create(out)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, rendered)
		return err
	}(); err != nil {
		return err
	}

	return nil
}

func (c *Compiler) Build() error {
	os.RemoveAll(c.outDir)

	for _, d := range c.documents {
		if err := c.buildDoc(d); err != nil {
			return err
		}
	}

	return nil
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

	if err := tmpl.Execute(&out, d); err != nil {
		return nil, err
	}

	return &out, nil
}
