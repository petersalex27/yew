package errors

import (
	"errors"
	"fmt"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util/code"
)

func windowError(s api.SourceCode, typ string, msg string, start, end int) error {
	window := code.Window(s, start, end)
	line, char := code.CalcLocation(s, start, false)
	e := fmt.Sprintf("[%d:%d] error (%s): %s\n%s", line, char, typ, msg, window)
	return errors.New(e)
}

func windowWarning(s api.SourceCode, msg string, start, end int) error {
	window := code.Window(s, start, end)
	line, char := code.CalcLocation(s, start, false)
	e := fmt.Sprintf("[%d:%d] warning: %s\n%s", line, char, msg, window)
	return errors.New(e)
}

func Syntax(s api.SourceCode, msg string, start, end int) error {
	return windowError(s, "syntax", msg, start, end)
}

func Warning(s api.SourceCode, msg string, start, end int) error {
	return windowWarning(s, msg, start, end)
}

func Type(s api.SourceCode, msg string, start, end int) error {
	return windowError(s, "type", msg, start, end)
}

func Lexical(s api.SourceCode, msg string, start, end int) error {
	return windowError(s, "lexical", msg, start, end)
}

func Name(s api.SourceCode, msg string, start, end int) error {
	return windowError(s, "name", msg, start, end)
}

func OS(msg string) error {
	return fmt.Errorf("error (os): %s", msg)
}
