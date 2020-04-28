package logictest

import (
	"fmt"
	"testing"

	"github.com/cockroachdb/datadriven"
	"github.com/justinj/scribe/code/executor"
)

// TODO: either vendor datadriven, rewrite it, or reference the cockroach one.
func TestLogic(t *testing.T) {
	datadriven.Walk(t, "testdata", func(t *testing.T, path string) {
		e := executor.New()
		datadriven.RunTest(t, path, func(t *testing.T, td *datadriven.TestData) string {
			r, err := e.Run(td.Input)
			if err != nil {
				return fmt.Sprintf("error: %s\n", err)
			}
			return fmt.Sprintf("%s\n", r.Msg)
		})
	})
}
