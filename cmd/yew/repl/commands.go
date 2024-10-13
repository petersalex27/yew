package repl

import (
	"fmt"

	"github.com/petersalex27/yew/api/util"
	"github.com/petersalex27/yew/internal/parser"
)

func expose(p *parser.ParserState) []error {
	tokens, err := util.Tokenize((*p.ReferenceScanner()), nil)
	if err != nil {
		return []error{(*err).Error()}
	}

	for _, token := range tokens {
		fmt.Printf("yew> %s\n", util.ExposeToken(token))
	}
	return nil // TODO: finish
}
