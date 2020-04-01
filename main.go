package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/justinj/scribe/processors"
)

func main() {
	in := "testbook/"
	out := "build/"

	if err := os.RemoveAll(out); err != nil {
		// If it doesn't exist, that's fine.
	}

	err := filepath.Walk(in, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}

		res, err := processors.Build(f, "code/")
		if err != nil {
			return err
		}

		postPath := path[len(in):]
		newPath := filepath.Join(out, postPath) + ".html"

		dir, _ := filepath.Split(newPath)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}

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
