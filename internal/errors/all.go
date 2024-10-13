package errors

import (
	"errors"
	"fmt"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util"
)

func windowError(s api.SourceCode, typ string, msg string, start, end int) error {
	window := util.Window(s, start, end)
	line, char := util.CalcLocation(s, start, false)
	e := fmt.Sprintf("[%d:%d] Error (%s): %s\n%s", line, char, typ, msg, window)
	return errors.New(e)
}

func windowWarning(s api.SourceCode, msg string, start, end int) error {
	window := util.Window(s, start, end)
	line, char := util.CalcLocation(s, start, false)
	e := fmt.Sprintf("[%d:%d] Warning: %s\n%s", line, char, msg, window)
	return errors.New(e)
}

func Syntax(s api.SourceCode, msg string, start, end int) error {
	return windowError(s, "Syntax", msg, start, end)
}

func Warning(s api.SourceCode, msg string, start, end int) error {
	return windowWarning(s, msg, start, end)
}

func Type(s api.SourceCode, msg string, start, end int) error {
	return windowError(s, "Type", msg, start, end)
}

func Lexical(s api.SourceCode, msg string, start, end int) error {
	return windowError(s, "Lexical", msg, start, end)
}

func OS(msg string) error {
	return errors.New(fmt.Sprintf("Error (OS): %s", msg))
}
