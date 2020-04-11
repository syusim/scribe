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

func TestParseRelExpr(t *testing.T) {
	testCases := []string{
		"table-name",
		"(join foo bar true)",
		"(join foo bar (= a b))",
		"(join foo (join bar baz true) (= a b))",
		"(select table-name (and (= a 2) (= b 3)))",
		"(as foo bar)",
		"(as foo bar [u v w])",
	}

	for _, tc := range testCases {
		e, err := ParseRelExpr(tc)
		if err != nil {
			t.Fatalf("failed to parse %q: %s", tc, err)
		}

		fmted := ExprStr(e)
		if tc != fmted {
			t.Errorf("in: %q, out: %q", tc, fmted)
		}
	}
}

func TestParseStatement(t *testing.T) {
	testCases := []string{
		"(run (join foo bar true))",
		"(run bar)",
		"(create-table bar [[x int] [y int] [z int]] [[10 20 30] [40 50 60]])",
		"(create-table bar [[x int] [y int] [z int]] [])",
	}

	for _, tc := range testCases {
		e, err := ParseStatement(tc)
		if err != nil {
			t.Fatalf("failed to parse %q: %s", tc, err)
		}

		fmted := ExprStr(e)
		if tc != fmted {
			t.Errorf("in: %q, out: %q", tc, fmted)
		}
	}
}
