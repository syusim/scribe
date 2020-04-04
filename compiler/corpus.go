package compiler

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/justinj/scribe/snippets"
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
		defer f.Close()
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

func (c *corpus) getSnip(seenFlags snippets.FlagSet, name string) (pre, mid, post string, ok bool) {
	referencedFile, ok := c.tags[name]
	if !ok {
		return "", "", "", false
	}
	file := c.files[referencedFile]
	// TODO: unindent stuff?
	extracted := snippets.ExtractCtx(file, seenFlags, name)
	return extracted.Pre, extracted.Contents, extracted.Post, true
}
