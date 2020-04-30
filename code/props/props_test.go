package props

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/datadriven"
)

func TestProps(t *testing.T) {
	datadriven.Walk(t, "testdata", func(t *testing.T, path string) {
		datadriven.RunTest(t, path, func(t *testing.T, td *datadriven.TestData) string {
			switch td.Cmd {
			case "fd":
				return "fd!\n"
			default:
				panic(fmt.Sprintf("unhandled: %q", td.Cmd))
			}
		})
	})
}
