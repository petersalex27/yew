package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/common/data"
)

// annotatable things
func (def *def) annotate(as data.Maybe[annotations])              { def.annotations = as }
func (specDef *specDef) annotate(as data.Maybe[annotations])      { specDef.annotations = as }
func (specInst *specInst) annotate(as data.Maybe[annotations])    { specInst.annotations = as }
func (typeDef *typeDef) annotate(as data.Maybe[annotations])      { typeDef.annotations = as }
func (typeAlias *typeAlias) annotate(as data.Maybe[annotations])  { typeAlias.annotations = as }
func (typing *typing) annotate(as data.Maybe[annotations])        { typing.annotations = as }
func (syntax *syntax) annotate(as data.Maybe[annotations])        { syntax.annotations = as }
func (cons *typeConstructor) annotate(as data.Maybe[annotations]) { cons.annotations = as }
func (module *module) annotate(as data.Maybe[annotations])        { module.annotations = as }

// parses annotations enclosed in brackets, like '[@inline]'
//
// the annotation's brackets are not included in the result, but any other brackets are
func parseOptionalEnclosedAnnotation(p Parser) data.Either[data.Ers, data.Maybe[enclosedAnnotation]] {
	if !token.LeftBracketAt.Match(p.current()) {
		n := data.Nothing[enclosedAnnotation](p.current())
		return data.Ok(n)
	}
	p.advance() // consume '[@'
	p.dropNewlines()
	mId := parseIdent(p)
	if mId.IsNothing() { // must have an id if we have an opening annotation bracket
		return data.Fail[data.Maybe[enclosedAnnotation]](ExpectedId, p.current())
	}

	id, _ := mId.Break()

	openBrackets := 1
	tokens := data.Nil[api.Node]()
	// parse the enclosed annotation until the closing bracket
	for {
		p.dropNewlines()
		cur := p.current()
		if token.LeftBracket.Match(cur) {
			openBrackets++
		} else if token.RightBracket.Match(cur) {
			openBrackets--
		}

		if openBrackets == 0 {
			p.advance() // consume ']' and don't add it to the data.List
			break
		} else if token.EndOfTokens.Match(cur) {
			return data.Fail[data.Maybe[enclosedAnnotation]](UnexpectedEOF, cur)
		}

		// add token to the data.List, including enclosed brackets
		tokens = tokens.Snoc(cur)
		p.advance() // consume token
	}

	annot := data.EMakePair[enclosedAnnotation](id, tokens)
	return data.Ok(data.Just(annot))
}

// parses things like '--@inline'
func parseOptionalFlatAnnotation(p Parser) data.Maybe[flatAnnotation] {
	cur := p.current()
	if !token.FlatAnnotation.Match(cur) {
		return data.Nothing[flatAnnotation](cur)
	}
	p.advance()
	return data.Just(data.EOne[flatAnnotation](cur))
}

// parses a single annotation
func parseAnnotation(p Parser) data.Either[data.Ers, data.Maybe[annotation]] {
	if matchCurrent(token.FlatAnnotation)(p) {
		unit, just := parseOptionalFlatAnnotation(p).Break()
		if !just {
			return data.Inr[data.Ers](data.Nothing[annotation](p))
		}
		// lift the result into an either after generalizing the annotation
		return data.Inr[data.Ers](data.Just(data.Inl[enclosedAnnotation](unit)))
	}

	// parse enclosed annotation
	es, res, isRes := parseOptionalEnclosedAnnotation(p).Break()
	if !isRes {
		return data.Inl[data.Maybe[annotation]](es)
	} else if unit, just := res.Break(); !just {
		return data.Inr[data.Ers](data.Nothing[annotation](res))
	} else {
		// lift the result into an data.Either after generalizing the annotation
		return data.Inr[data.Ers](data.Just(data.Inr[flatAnnotation](unit)))
	}
}

func annotationIteration(p Parser) data.Either[data.Ers, data.Maybe[annotation]] {
	res := parseAnnotation(p)
	if !res.IsLeft() {
		p.dropNewlines()
	}
	return res
}

// parses a block of annotations
func parseAnnotations(p Parser) data.Either[data.Ers, data.Maybe[annotations]] {
	cur := p.current()
	if !token.FlatAnnotation.Match(cur) && !token.LeftBracketAt.Match(cur) {
		return data.Ok(data.Nothing[annotations](p)) // data.Ok result, data.Just no annotations found
	}

	var annots data.NonEmpty[annotation]
	has1stAnnot := false
	for {
		lhs, rhs, isRight := annotationIteration(p).Break()
		if !isRight {
			return data.Inl[data.Maybe[annotations]](lhs)
		} else if rhs.IsNothing() {
			break
		}

		unit, _ := rhs.Break()
		if !has1stAnnot {
			has1stAnnot = true
			annots = data.Singleton(unit)
		} else {
			annots = annots.Snoc(unit)
		}
	}
	if !has1stAnnot {
		return data.Ok(data.Nothing[annotations](p))
	}
	return data.Ok(data.Just(annotations{annots}))
}
