//go:build debug
// +build debug

package lexer

import "fmt"

func debug_validateLineAndChar(lex *Lexer) {
	okLine := lex.ValidLine()
	okChar := lex.ValidChar()
	if !(okLine && okChar) {
		panic(fmt.Sprintf("line ok: %t; char ok: %t", okLine, okChar))
	}
}
