// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package parser

import (
	"strings"
	"unicode"

	"github.com/petersalex27/yew/token"
)

func badLambdaBinder(s string) bool {
	if len(s) == 0 {
		return true
	}
	return unicode.IsLower(rune(s[0]))
}

func (parser *Parser) getAfterAbstractionVar(lambda *Lambda, data *actionData) (ok, again bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	again = tok.Type == token.Comma
	if ok = again || tok.Type == token.ThickArrow; !ok {
		parser.error2(ExpectedCommaOrThickArrow, tok.Start, tok.End)
	} else {
		lambda.End = tok.End // keep updating end point so errors can use this value if needed
	}
	return
}

func (parser *Parser) getAbstractionVar(lambda *Lambda, data *actionData) (ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	switch tok.Type {
	case token.Underscore:
		break // this is okay: wildcard
	case token.Id, token.ImplicitId:
		ok = strings.ToLower(tok.Value) == tok.Value && badLambdaBinder(tok.Value)
		if !ok {
			parser.error2(BadIdent, tok.Start, tok.End)
			return
		}
	default:
		ok = false
		parser.error2(ExpectedIdentifier, tok.Start, tok.End)
		return
	}

	id := makeIdent(tok)
	found := parser.findTermInTop(tok)
	if found {
		// okay, convert shadowed name to wildcard and throw a warning
		parser.warnShadowedLambdaBinder(lambda, id.Name)
	}
	lambda.Binders = append(lambda.Binders, id)
	return true
}

// creates a lambda abstraction
//
//	(\v0, .., vN => e) : T0 -> .. -> TN -> Te
func abstractionAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	parser.declarations.Increase()
	lambda := Lambda{
		Binders: make([]Ident, 0, 4),
		Start:   tok.Start,
	}

	ok = parser.getAbstractionVar(&lambda, data)
	arity := uint(0)
	again := true
	for again && ok {
		// if loop makes it to this point, increase number of arguments lambda takes
		arity++
		ok, again = parser.getAfterAbstractionVar(&lambda, data)
		if !ok || !again {
			break
		}
		ok = parser.getAbstractionVar(&lambda, data)
	}

	if ok {
		term = termElem{lambda, termInfo{10, false, arity}}
	}

	parser.shift(term)
	parser.terms.Save()

	term, ok = parser.terminalAction(data) // require something after '=>'
	// resolution will finish lambda abstraction
	return
}

// noop if binder with Name `name` is not found; else, updates shadowed binder and writes a warning
func (parser *Parser) warnShadowedLambdaBinder(lambda *Lambda, name string) {
	var id *Ident = nil
	for i, v := range lambda.Binders {
		if v.Name == name {
			lambda.Binders[i].Name = "_"
			id = &v
			break
		}
	}

	if id != nil {
		parser.warning2(NonBindingVariable_warn, id.Start, id.End)
	}
}
