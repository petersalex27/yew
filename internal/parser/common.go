// =================================================================================================
// Alex Peters - May 05, 2024
// =================================================================================================
package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/petersalex27/yew/token"
	"github.com/petersalex27/yew/types"
)

type positioned interface {
	Pos() (start, end int)
}

func Var (s string) types.Variable {
	return types.Var(token.Id.MakeValued(s))
}

// transfers messages if there are any, otherwise no-op
func (parser *Parser) transferEnvErrors() {
	if msgs := parser.env.FlushMessages(); len(msgs) > 0 {
		parser.messages = append(parser.messages, msgs...)
		parser.panicking = true
	}
}

func (export exports) String() string {
	var b *strings.Builder
	w := &walkingWriter{false, b,}
	bothNil := true
	if export.declTable != nil {
		bothNil = false
		export.declTable.Walk(printDeclarationsWalker(w))
	}
	if export.Locals != nil {
		bothNil = false
		export.Locals.Walk(printLocalsWalker(w))
	}
	if bothNil {
		return ""
	}
	return b.String()
}

func (info termInfo) String() string {
	return fmt.Sprintf("info: bp=%d, rAssoc=%t, arity=%d, infixed=%t", info.bp, info.rAssoc, info.arity, info.infixed)
}

type walkingWriter struct {
	unset bool
	io.Writer
}

func printLocalsWalker(w *walkingWriter) func(s fmt.Stringer, r types.Replacement) {
	return func(s fmt.Stringer, r types.Replacement) {
		if !w.unset {
			w.unset = true
		} else {
			fmt.Fprintf(w, "\n")
		}

		fmt.Fprintf(w, "%v := %v : %v", s, r.Term, r.Type)
	}
}

func printDeclarationsWalker(w *walkingWriter) func(s fmt.Stringer, d *declaration) {
	return func(s fmt.Stringer, d *declaration) {
		if !w.unset {
			w.unset = true
		} else {
			fmt.Fprintf(w, "\n")
		}

		if d.implicit {
			fmt.Fprintf(w, "{%v} ", s)
		} else {
			fmt.Fprintf(w, "%v ", s)
		}

		if d.termInfo != nil {
			fmt.Fprintf(w, "[%v]", *d.termInfo)
		} else {
			fmt.Fprintf(w, "[info: _]")
		}

		res := d.available.String()
		if res != "" {
			res = strings.ReplaceAll(res, "\n", "\n\t")
			fmt.Fprintf(w, "\n\texports:\n\t%s", res)
		}
	}
}

func (parser *Parser) printDeclarations(w io.Writer) {
	/*
	==========================
	Declarations:
	==========================
	(+) [info: bp=6, rAssoc=false, arity=2, infixed=false]
	+ [info: bp=6, rAssoc=false, arity=2, infixed=true]
	TypeOf [info: bp=10, rAssoc=false, arity=1, infixed=false]
	  exports:
	  {a} [info: bp=0, rAssoc=false, arity=0, infixed=false]
		a := a : Type
	...
	==========================
	*/
	walker := &walkingWriter{
		false,
		w,
	}
	fmt.Fprintf(w, "\n==========================\nDeclarations:\n==========================\n")
	parser.declarations.Walk(printDeclarationsWalker(walker))
	fmt.Fprintf(w, "\n==========================\n\n")
}
