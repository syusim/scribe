package ast

import "testing"

// TODO: test error cases.

func TestParseExpr(t *testing.T) {
	testCases := []string{
		"1",
		"2",
		"3",
		"foo",
		`"foo"`,
		"(+ 1 2)",
		"(+ 1 2 3)",
		"(+ 1 2 (+ 4 5))",
		"(+ (+ 3 4) 2 (+ 4 5))",
		"(- 1 2)",
		"(* 1 2)",
		`(+ "hello" "world")`,
		`(+ foo bar)`,
		"true",
		"false",
	}

	for _, tc := range testCases {
		e, err := ParseExpr(tc)
		if err != nil {
			t.Fatalf("failed to parse %q: %s", tc, err)
		}

		fmted := ExprStr(e)
		if tc != fmted {
			t.Errorf("in: %q, out: %q", tc, fmted)
		}
	}
}

func TestParseQuery(t *testing.T) {
	testCases := []string{
		"table-name",
		"(join foo bar true)",
		"(join foo bar (= a b))",
		"(join foo (join bar baz true) (= a b))",
		// "(select table-name 1 (+ 2 3) foo)",
	}

	for _, tc := range testCases {
		e, err := ParseQuery(tc)
		if err != nil {
			t.Fatalf("failed to parse %q: %s", tc, err)
		}

		fmted := ExprStr(e)
		if tc != fmted {
			t.Errorf("in: %q, out: %q", tc, fmted)
		}
	}
}
