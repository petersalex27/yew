//go:build debug
// +build debug

package parser

import (
	"fmt"
	"io"
	"github.com/petersalex27/yew/internal/types"
)

// information for debugging, embed inside parser
type debug_info_parser struct {
	// counter to be updated for each test iteration so conditional breakpoints can be set to stop for
	// some test iteration
	testCounter int
}

func (dip *debug_info_parser) debug_incTestCounter() { dip.testCounter++ }

func (dip *debug_info_parser) debug_resetCounter() { dip.testCounter = 0 }

func debug_log_reduce(w io.Writer, a, b, result termElem) {
	fmt.Fprintf(w, "red: (%v) (%v) = %v\n", a, b, result)
}

func debug_log_builtin(w io.Writer, name fmt.Stringer, term, typ types.Term) {
	fmt.Fprintf(w, "init builtin %s, %s : %s\n", name, term, typ)
}