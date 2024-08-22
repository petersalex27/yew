// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package parser

import (
	"os"
	"strings"
	"unicode"

	"github.com/petersalex27/yew/common/stack"
	"github.com/petersalex27/yew/token"
	"github.com/petersalex27/yew/types"
)

func goodLambdaBinder(s string) bool {
	return len(s) != 0 && unicode.IsLower(rune(s[0]))
}

func (parser *Parser) getAfterAbstractionVar(data *actionData) (end int, ok, again bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	again = tok.Type == token.Comma
	if ok = again || tok.Type == token.ThickArrow; !ok {
		parser.error2(ExpectedCommaOrThickArrow, tok.Start, tok.End)
	} else {
		end = tok.End // keep updating end point so errors can use this value if needed
	}
	return
}

/*
// (?a : ?A) -> (?b : ?B)
	var intro types.PiIntro
	kA := types.GetKind(&x.A)
	intro, ok = parser.env.Prod(types.Wildcard(), x.A, kA)
 	if !ok {
		parser.transferEnvErrors()
		return
	}

	P, _, yes := intro(B, B.A)
	if ok = yes; !ok {
		parser.transferEnvErrors()
		return
	}
*/

func (parser *Parser) getAbstractionVar(data *actionData) (derive types.AbsSecondFunc, x types.Variable, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	switch tok.Type {
	case token.Underscore:
		break // this is okay: wildcard
	case token.Id:
		ok = strings.ToLower(tok.Value) == tok.Value && goodLambdaBinder(tok.Value)
		if !ok {
			parser.error2(BadIdent, tok.Start, tok.End)
			return
		}
	default:
		ok = false
		parser.error2(ExpectedIdentifier, tok.Start, tok.End)
		return
	}

	//id := makeIdent(tok)
	found := parser.findTermInTop(tok)
	if found {
		// okay, convert shadowed name to wildcard and throw a warning
		//parser.warnShadowedLambdaBinder(lambda, id.Name)
	}
	x = types.Var(tok)
	A := types.GetKind(&x)

	derive = parser.env.Abs(x, A)
	return
}

// creates a lambda abstraction
//
//	(\v0, .., vN => e) : T0 -> .. -> TN -> Te
func abstractionAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	//parser.declarations.Increase()
	type deriveData struct {
		derive types.AbsSecondFunc
		x      types.Variable
	}

	// parse abstraction binders
	var d deriveData

	d.derive, d.x, ok = parser.getAbstractionVar(data)
	deriveStack := stack.NewStack[deriveData](4)
	if ok {
		deriveStack.Push(d)
	} else {
		return termElem{}, false
	}

	arity := uint32(0)
	again := true
	var end int
	for again && ok {
		// if loop makes it to this point, increase number of arguments lambda takes
		arity++
		end, ok, again = parser.getAfterAbstractionVar(data)
		if !ok || !again {
			break
		}
		d.derive, d.x, ok = parser.getAbstractionVar(data)
		if ok {
			deriveStack.Push(d)
		}
	}

	if !ok {
		return termElem{}, false
	}

	var body termElem
	// process the bound expression (the part after '\ .. =>')
	if body, ok = parser.process(-1, data); !ok {
		return
	}

	// create lambda abstraction
	var b types.Term = body.Term
	B := types.GetKind(&b)
	var lambda types.Lambda
	for !deriveStack.Empty() {
		d, _ = deriveStack.Pop()
		
		// (?a : ?A) -> (?b : ?B)
		var intro types.PiIntro
		{
			A := types.AsTyping(d.x.Kind)
			intro, ok = parser.env.Prod_NoGeneralization(A)
			if !ok {
				parser.transferEnvErrors()
				return
			}
		}

		var P types.Pi
		if P, ok = intro(B); !ok {
			parser.transferEnvErrors()
			return
		}

		b, B, ok = d.derive(b, B)(P)
		if !ok {
			parser.transferEnvErrors()
			return
		}
	}

	arity += types.CalcArity(body.Term)   // incorporate arity of expression into arity of lambda 
	if ok = types.SetKind(&b, B); !ok {
		parser.transferEnvErrors()
		return
	}
	lambda = b.(types.Lambda)
	term = termElem{lambdaType, lambda, termInfo{10, false, arity, false}, tok.Start, end}

	debug_log_reduce(os.Stderr, termElem{}, body, term)

	return term, true
}
