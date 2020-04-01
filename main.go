package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/justinj/scribe/processors"
)

func main() {
	p, err := processors.MakeSnippetProcessor("code/")
	if err != nil {
		panic(err)
	}

	in := "testbook/"
	out := "build/"

	if err := os.RemoveAll(out); err != nil {
		// If it doesn't exist, that's fine.
	}

	err = filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}

		res, err := p(f)
		if err != nil {
			return err
		}

		postPath := path[len(in):]
		newPath := filepath.Join(out, postPath)

		dir, _ := filepath.Split(newPath)
		os.MkdirAll(dir, 0700)

		out, err := os.Create(newPath)
		if err != nil {
			return err
		}
		defer out.Close()
		io.Copy(out, res)

		return nil
	})
	_ = out

	if err != nil {
		panic(err)
	}

}
