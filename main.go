package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/justinj/scribe/processors"
)

func main() {
	build()

	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	// log.Println("Listening on :3000...")
	// err := http.ListenAndServe(":3000", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func build() {
	in := "testbook/"
	out := "build/"

	os.RemoveAll(out)

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

	if err != nil {
		panic(err)
	}
}
