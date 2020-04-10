package sexp

import "testing"

func TestParse(t *testing.T) {
	tcs := []string{
		"1",
		"(1)",
		"((1))",
		"((1) (2))",
		"()",
		`"   "`,
	}

	for _, tc := range tcs {
		r, err := Parse(tc)
		if err != nil {
			t.Error(err)
			continue
		}
		s := r.String()
		if tc != s {
			t.Errorf("not equal: %q vs %q", tc, s)
		}
	}

	tcs = []string{
		"  1  ",
		"( 1 )",
		"((1))",
		"( [ 1 ] )",
		"((1) (2))",
		"((1)   (2))  ",
		"(first (list (+ 1 2) 3))",
		`("hello" "world")`,
		"[]",
	}

	for _, tc := range tcs {
		r, err := Parse(tc)
		if err != nil {
			t.Error(err)
			continue
		}
		s1 := r.String()
		r2, err := Parse(s1)
		if err != nil {
			t.Error(err)
			continue
		}
		s2 := r2.String()
		if s1 != s2 {
			t.Errorf("not equal: %q vs %q", s1, s2)
		}
	}
}

func TestPretty(t *testing.T) {
	tcs := []struct {
		in  string
		out string
	}{
		{"1", "1"},
		{"(1)", "(1)"},
		{"(1 2 3)", "(1\n 2\n 3)"},
	}

	for _, tc := range tcs {
		s, _ := Parse(tc.in)
		r := Pretty(s)
		if tc.out != r {
			t.Errorf("not equal: %q vs %q", tc.out, r)
		}
	}
}
