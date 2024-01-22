// +build debug

package lexer

func debug_validateLineAndChar(lex *Lexer) {
	okLine := lex.ValidLine()
	okChar := lex.ValidChar()
	if !(okLine && okChar) {
		panic("line ok: ", okLine, "; char ok: ", okChar)
	}
}