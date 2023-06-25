package errorgen

import "yew/lex"
import err "yew/error"

func GenerateSyntaxError(message string) func(scan.Token, scan.InputStream) err.Error {
	return func(t scan.Token, i scan.InputStream) err.Error {
		loc := t.GetLocation()
		return err.CompileMessage(
			message, err.ERROR, err.SYNTAX, (i).GetPath(), loc.GetLine(), loc.GetChar(),
			t.GetSourceIndex(), (i).GetSource()).(err.Error)
	}
}