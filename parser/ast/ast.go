package ast

import (
	"fmt"
	//err "yew/error"
	//. "yew/parser/node-type"
	"yew/parser/parser"
)

func printLines(xs []string) []string {
	for _, x := range xs {
		fmt.Printf("%s", x)
	}

	if xs[len(xs)-1] == " └─" {
		xs[len(xs)-1] = "   "
	} else if xs[len(xs)-1] == " ├─" {
		xs[len(xs)-1] = " │ "
	} // else stays the same

	return xs
}

func printSpaces(n int) {
	// ├ ┬ ┤ ┼ ┴ ─ │ └ ┘ ┌ ┐ ┄ ┆ 
	for i := 0; i < n - 1; i++ {
		fmt.Printf(" │ ")
	}
	if n > 0 {
		fmt.Printf(" ├─")
	}
}

func printEndSpaces(n int) {
	for i := 0; i < n - 1; i++ {
		fmt.Printf(" │ ")
	}
	if n > 0 {
		fmt.Printf(" └─")
	}
}

func EqualTest(a parser.Ast, b parser.Ast) bool {
	return a.Equal_test(b)
}

func PrintAst(a parser.Ast) {
	a.Print([]string{""})
}