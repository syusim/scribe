package main

import (
	"flag"
	"fmt"
)

// To compile an expression, we must generate:
// * the struct definition (in either scalar or memo)
// * the interning func (in memo)
// * the the memo construction func? (not yet i guess)

var target = flag.String("o", "", "where to put the output file")

type field struct {
	Name string
	Typ  string
}

type scalarExprDef struct {
	Name   string
	Typ    string
	Fields []field
}

var scalarExprs = []scalarExprDef{
	{
		Name: "Plus",
		Typ:  "lang.DInt",
		Fields: []field{
			{"Left", "Expr"},
			{"Right", "Expr"},
		},
	}, {
		Name: "Minus",
		Typ:  "lang.DInt",
		Fields: []field{
			{"Left", "Expr"},
			{"Right", "Expr"},
		},
	}, {
		Name: "Times",
		Typ:  "lang.DInt",
		Fields: []field{
			{"Left", "Expr"},
			{"Right", "Expr"},
		},
	},
}

func main() {
	flag.Parse()
	fmt.Println("vim-go")
}
