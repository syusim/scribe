package snippets

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/justinj/bitwise/datadriven"
)

func TestSnippets(t *testing.T) {
	datadriven.Walk(t, "testdata", func(t *testing.T, path string) {
		files := make(map[string]Block)
		datadriven.RunTest(t, path, func(d *datadriven.TestData) string {
			switch d.Cmd {
			case "load":
				var name string
				for _, a := range d.CmdArgs {
					if a.Key == "name" {
						if len(a.Vals) != 1 {
							t.Fatal("name needs one arg")
						}
						name = a.Vals[0]
					}
				}
				if name == "" {
					t.Fatal("load requires a name")
				}
				result, err := New(strings.NewReader(d.Input + "\n"))
				if err != nil {
					return fmt.Sprintf("error: %s\n", err.Error())
				}
				files[name] = result
				return ""
			case "extract":
				flagSet := make(FlagSet)
				var name string
				var section string
				for _, a := range d.CmdArgs {
					switch a.Key {
					case "flags":
						for _, f := range a.Vals {
							flagSet[f] = struct{}{}
						}
					case "name":
						name = a.Vals[0]
					case "section":
						section = a.Vals[0]
					}
				}
				var buf bytes.Buffer
				Extract(files[name], &buf, flagSet, section)
				return buf.String()
			case "render":
				flagSet := make(FlagSet)
				var name string
				for _, a := range d.CmdArgs {
					switch a.Key {
					case "flags":
						for _, f := range a.Vals {
							flagSet[f] = struct{}{}
						}
					case "name":
						name = a.Vals[0]
					}
				}
				var buf bytes.Buffer
				files[name].Render(&buf, flagSet)
				return buf.String()
			default:
				panic("unknown command")
			}
		})
	})
}
