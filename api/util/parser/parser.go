package parser

import (
	"fmt"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util/code"
	"github.com/petersalex27/yew/api/util/node"
)

// returns a string representation of exposed api
func exposeError(err error) string {
	// line, char := err.Location()
	// return fmt.Sprintf(
	// 	"Error{SourceCode: %q, Location: (%d, %d), Type: %s, Message: %q, ConfigName: %s}",
	// 	err.SourceCode(), line, char, err.Type(), err.Message(), err.ConfigName())
	return "Error{" + err.Error() + "}"
}

func InitParser[T api.Parser](s api.Scanner) (out T) {
	out = out.Clear().(T)
	if s == nil {
		s = code.DummyScanner()
	}
	*out.ReferenceScanner() = s
	return
}

func ExposeParser(p api.Parser) string {
	var exposedScanner string
	scannerRef := p.ReferenceScanner()
	if scannerRef == nil {
		exposedScanner = "nil"
	} else {
		exposedScanner = code.ExposeScanner(*scannerRef)
	}

	es := node.ExposeList(exposeError, p.Errors(), ", ")
	return fmt.Sprintf("Parser{Scanner: %s, Ast: %s, Errors: %s}", exposedScanner, node.ExposeNode(p.Ast()), es)
}