package parser

import (
	"strconv"
	"strings"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util/code"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/kind"
	"github.com/petersalex27/yew/internal/lexer"
)

func breakFlatAnnotation(annot flatAnnotation) (id api.Token, out []api.Node, warning *data.Err) {
	lex := lexer.Init(code.StringSource(annot.Break().String()))
	id = lex.Scan()
	if !token.Id.Match(id) {
		warning = new(data.Err)
		*warning = data.MkErr(lexer.ExpectedAnnotationId, annot)
		return id, nil, warning
	}

	for !lex.Eof() {
		tok := lex.Scan()
		if token.Error.Match(tok) {
			warning = new(data.Err)
			*warning = data.MkErr(tok.String(), annot)
			return id, nil, warning
		}
		out = append(out, lex.Scan())
	}

	return id, out, nil
}

func metaAnnots(p parser, annots annotations) {
	for _, annot := range annots.Elements() {
		var id api.Token
		var body []api.Node

		if flat, enc, isEnc := annot.Break(); !isEnc {
			var warning *data.Err
			id, body, warning = breakFlatAnnotation(flat)
			if warning != nil {
				p.report(parseError(p, *warning), false) // not fatal, a warning
				continue                                 // try the next annotation
			}
		} else {
			id = data.Cases(enc.Fst(), (lowerIdent).Break, (upperIdent).Break)
			body = enc.Snd().Elements()
		}

		resolveAnnot(p, id, body)
	}
}

func resolveAnnot(p parser, id api.Token, body []api.Node) {
	var warning *data.Err
	switch id.String() {
	case "infix", "infixl", "infixr":
		warning = parseInfixAnnotation(p, id, body)
	default:
		warning = new(data.Err)
		*warning = data.MkErr(UnknownMetaAnnotation, id)
	}

	if warning != nil {
		p.report(parseError(p, *warning), false) // not fatal, a warning
	}
}

func getInt64(s string) (int64, bool) {
	s = strings.Replace(s, "_", "", -1)
	i, err := strconv.ParseInt(s, 0, 64)
	return i, err == nil
}

func validatePrec(prec api.Token) (byte, bool) {
	if !token.IntValue.Match(prec) {
		return 0, false
	}

	res, ok := getInt64(prec.String())
	if !ok {
		return 0, false
	} else if res < 0 || res > 10 {
		return 0, false
	}

	return byte(res), true
}

func parseInfixAnnotation(p parser, id api.Token, body []api.Node) *data.Err {
	assoc := kind.None
	switch id.String() {
	case "infix":
		assoc = kind.Infix
	case "infixl":
		assoc = kind.Infixl
	case "infixr":
		assoc = kind.Infixr
	}

	if len(body) == 0 {
		warning := new(data.Err)
		*warning = data.MkErr(MissingPrecedence, id)
		return warning
	}

	precTok, ok := body[0].(api.Token)
	var prec byte
	if ok {
		prec, ok = validatePrec(precTok)
	}
	if !ok {
		warning := new(data.Err)
		*warning = data.MkErr(IllegalPrecedence, body[0])
		return warning
	}

	const expected = token.Infix
	for _, node := range body[1:] {
		if !expected.Match(node) {
			warning := new(data.Err)
			*warning = data.MkErr(ExpectedInfixId, node)
			return warning
		}
		tok, _ := node.(api.Token)
		declareFromName(p, tok, assoc, prec)
	}

	return nil
}
