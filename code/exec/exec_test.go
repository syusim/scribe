package exec

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/datadriven"
	"github.com/justinj/scribe/code/ast"
	"github.com/justinj/scribe/code/lang"
	"github.com/justinj/scribe/code/sexp"
)

type arg struct {
	key, val string
}

func parseArgs(s string) []arg {
	args := make([]arg, 0)
	for _, l := range strings.Split(s, "\n") {
		vs := strings.Split(l, "=")
		args = append(args, arg{
			key: strings.TrimSpace(vs[0]),
			val: strings.TrimSpace(vs[1]),
		})
	}

	return args
}

func printResult(node Node) string {
	rows := Spool(node)
	// TODO: pull this out, it also exists in the executor
	var buf bytes.Buffer
	for i, row := range rows {
		buf.WriteByte('[')
		for j, d := range row {
			if j > 0 {
				buf.WriteByte(' ')
			}
			d.Format(&buf)
		}
		buf.WriteByte(']')
		if i != len(rows)-1 {
			buf.WriteByte('\n')
		}
	}
	if buf.Len() > 0 {
		buf.WriteByte('\n')
	}
	return buf.String()
}

func TestExec(t *testing.T) {
	datadriven.Walk(t, "testdata", func(t *testing.T, path string) {
		rowSets := make(map[string][]lang.Row)
		datadriven.RunTest(t, path, func(t *testing.T, td *datadriven.TestData) string {
			switch td.Cmd {
			case "load":
				s, err := sexp.Parse(td.Input)
				if err != nil {
					return fmt.Sprintf("error: %s", err)
				}
				rows, err := ast.ParseRowList(s)
				if err != nil {
					return fmt.Sprintf("error: %s", err)
				}

				var name string
				for _, a := range td.CmdArgs {
					if a.Key == "name" {
						name = a.Vals[0]
					}
				}
				if name == "" {
					t.Fatal("need name for row set")
				}

				rowSets[name] = rows

				return ""

			case "cross":
				args := parseArgs(td.Input)
				var left Node
				var right Node
				for _, a := range args {
					switch a.key {
					case "left":
						left = Constant(rowSets[a.val])
					case "right":
						right = Constant(rowSets[a.val])
					}
				}

				return printResult(Cross(left, right))

			case "hash":
				args := parseArgs(td.Input)
				var left Node
				var right Node
				var leftIdxs []lang.ColOrdinal
				var rightIdxs []lang.ColOrdinal
				for _, a := range args {
					switch a.key {
					case "left":
						left = Constant(rowSets[a.val])
					case "right":
						right = Constant(rowSets[a.val])
					case "eq":
						s, err := sexp.Parse(a.val)
						if err != nil {
							t.Fatal(err)
						}
						ls := s.(sexp.List)
						for _, pair := range ls {
							leftIdxs = append(leftIdxs, lang.ColOrdinal(sexp.Int(sexp.Nth(pair, 0))))
							rightIdxs = append(rightIdxs, lang.ColOrdinal(sexp.Int(sexp.Nth(pair, 0))))
						}
					}
				}

				return printResult(Hash(left, right, leftIdxs, rightIdxs))

			case "merge":
				args := parseArgs(td.Input)
				var left Node
				var right Node
				var leftIdxs []lang.ColOrdinal
				var rightIdxs []lang.ColOrdinal
				for _, a := range args {
					switch a.key {
					case "left":
						left = Constant(rowSets[a.val])
					case "right":
						right = Constant(rowSets[a.val])
					case "eq":
						s, err := sexp.Parse(a.val)
						if err != nil {
							t.Fatal(err)
						}
						ls := s.(sexp.List)
						for _, pair := range ls {
							leftIdxs = append(leftIdxs, lang.ColOrdinal(sexp.Int(sexp.Nth(pair, 0))))
							rightIdxs = append(rightIdxs, lang.ColOrdinal(sexp.Int(sexp.Nth(pair, 0))))
						}
					}
				}

				return printResult(Merge(left, right, leftIdxs, rightIdxs))

			default:
				panic(fmt.Sprintf("unhandled: %q", td.Cmd))
			}
		})
	})
}
