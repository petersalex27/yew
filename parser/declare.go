package parser

import (
	"math"
	"strings"

	"github.com/petersalex27/yew/token"
)

type Definition struct {
}

type Declaration struct {
	implicit bool
	name     Ident
	typing   Term
	termInfo
}

type declare_meta struct {
	bp     uint8
	rAssoc bool
}

func (dm *declare_meta) LoadMetaData(parser *Parser, start, end int, args ...Term) (ok bool) {
	if ok = len(args) <= 2; !ok {
		start, _ = args[2].Pos()
		_, end = args[len(args)-1].Pos()
		parser.error2(UnexpectedMetaArgs, start, end)
		return
	}

	switch len(args) {
	case 2:
		arg := args[1]
		start, end = arg.Pos()
		if arg.NodeType() != intConstType {
			parser.error2(ExpectedUint, start, end)
			return
		}
		if ok = arg.(IntConst).int.X.IsUint64(); !ok {
			parser.error2(ExpectedUintRange1_9, start, end)
			return
		}
		res := uint8(arg.(IntConst).int.X.Uint64())
		dm.bp = res
		fallthrough
	case 1:
		arg := args[0]
		start, end = arg.Pos()
		if arg.NodeType() != identType {
			parser.error2(ExpectedLeftRightNone, start, end)
			return
		}
		res := arg.(Ident).Name
		if ok = "Left" == res || "None" == res; ok {
			dm.rAssoc = false
		} else if ok = "Right" == res; ok {
			dm.rAssoc = true
		} else {
			ok = false
			parser.error2(ExpectedLeftRightNone, start, end)
			return
		}
	case 0:
		dm.bp = 1
		dm.rAssoc = false
	}
	return
}

func (decl Declaration) makeTerm() termElem {
	return termElem{
		Term:     decl.name,
		termInfo: decl.termInfo,
	}
}

func fillInfo(decl *Declaration, prefixed bool, arity uint, bp int8, rAssoc bool) {
	if arity == 0 {
		decl.termInfo = termInfo{}
		return
	}

	if prefixed { // prefixed
		decl.termInfo = termInfo{bp: 10, rAssoc: false, arity: arity} // TODO: need to get info from annotation if one exists
		return
	}

	if arity > 2 {
		panic("TODO: allow affixed IDs with arity > 2?")
	}
	decl.termInfo = termInfo{bp: bp, rAssoc: rAssoc, arity: arity} // TODO: need to get info from annotation if one exists
}

func createUseName(name string, start, end int) Ident {
	name = strings.TrimPrefix(name, "_")
	name = strings.TrimSuffix(name, "_")
	if strings.Contains(name, "_") {
		panic("TODO: allow affixed IDs with arity > 2") // TODO
	}

	return Ident{Name: name, Start: start, End: end}
}

type Str string

func (s Str) String() string {
	return string(s)
}

func generate_setType(unprocessedName string, decl *Declaration) func(typ Term, args ...uint8) {
	return func(typ Term, args ...uint8) {
		prefixed := !strings.Contains(unprocessedName, "_")
		decl.typing = typ
		arity := calcArity(typ)
		var bp uint8 = 0
		rAssoc := false
		if len(args) > 0 {
			bp = args[0]
		}
		if len(args) > 1 {
			rAssoc = args[1] != 0
		}

		if bp > math.MaxInt8 {
			panic("bug: illegal binding power, cap is 127 from an internal origin and 9 from a source-code origin")
		}
		fillInfo(decl, prefixed, arity, int8(bp), rAssoc)
	}
}

func (parser *Parser) declareHelper(name string, start, end int, linking bool) (setType func(Term, ...uint8), ok bool) {
	unprocessed := name
	if linking {
		name = createUseName(name, start, end).Name
	}

	nm := Str(name)
	if _, found := parser.declarations.Find(nm); found {
		parser.error2(IllegalRedeclaration, start, end)
		ok = false
		return
	}

	ok = true
	decl := new(Declaration)
	decl.name = Ident{Name: name, Start: start, End: end}

	decl.termInfo = termInfo{} // set as default info for now
	parser.declarations.Map(nm, decl)
	setType = generate_setType(unprocessed, decl)
	return
}

type extensionType = uint8

const (
	keywordType extensionType = iota
	varType
	markerType
)

func MakeKeywordExt(keyword token.Token) extensionElem {
	key := Key{keyword.Value, keyword.Start, keyword.End}
	return extensionElem{typ: keywordType, Term: key}
}

func MakeVarExt(v token.Token) extensionElem {
	variable := Key{v.Value, v.Start, v.End}
	return extensionElem{typ: varType, Term: variable}
}

func MakeMarkerExt(mark token.Token) extensionElem {
	marker := Key{mark.Value, mark.Start, mark.End}
	return extensionElem{typ: markerType, Term: marker}
}

type extensionElem struct {
	typ uint8
	Term
}

type extensionFunc = func(Extension, *Parser, [][]token.Token) (termElem, bool)

type Extension struct {
	pattern   []extensionElem
	replace   extensionFunc
	verifiers []func(*Parser, []token.Token) (termElem, bool)
	fn        termElem
}

func calcParts(pattern []extensionElem) (parts []string) {
	parts = make([]string, 0, 8) // 8 is somewhat arbitrary
	for _, pat := range pattern {
		if pat.typ == varType {
			continue
		}
		parts = append(parts, pat.Term.String())
	}
	return parts[:len(parts)-1]
}

func parseExtVar(parser *Parser, ptr *uint, tokens []token.Token) (term termElem, ok bool) {
	if uint(len(tokens)) <= *ptr { // assumes len != 0
		tok := tokens[len(tokens)-1]
		parser.error2(ExpectedIdentifier, tok.Start, tok.End)
		return termElem{}, false
	}

	// get variable name
	*ptr++
	tok := tokens[*ptr-1]
	if ok = tok.Type != token.Id && tok.Type != token.ImplicitId; !ok {
		// not an id
		parser.error2(ExpectedIdentifier, tok.Start, tok.End)
		return
	} else if ok = validTypeIdent(tok.Value); !ok {
		// not a camel case id
		parser.error2(RequireCamelCase, tok.Start, tok.End)
		return
	}

	// create term
	key := Key{tok.Value, tok.Start, tok.End}
	term = termElem{Term: key, termInfo: termInfo{}}

	if uint(len(tokens)) <= *ptr {
		tok := tokens[len(tokens)-1]
		parser.error2(UnexpectedFinalTok, tok.Start, tok.End)
		return termElem{}, false
	}

	// get closing bracket
	*ptr++
	rbracket := tokens[*ptr-1]
	if ok = rbracket.Type == token.RightBracket; !ok {
		parser.error2(ExpectedRBracket, rbracket.Start, rbracket.End)
		return
	}

	return term, true
}

func (parser *Parser) parseExtension(tokens []token.Token) (ok bool) {
	ext := Extension{}
	ext.pattern = []extensionElem{}
	var i uint
	for i = 0; i < uint(len(tokens)) && tokens[i].Type != token.Equal; {
		i++
		switch tok := tokens[i-1]; tok.Type {
		case token.StringValue:
			ext.pattern = append(ext.pattern, MakeKeywordExt(tok))
		case token.Id:
			ext.pattern = append(ext.pattern, MakeMarkerExt(tok))
		case token.LeftBracket:
			var t termElem
			t, ok = parseExtVar(parser, &i, tokens)
			if !ok {
				return
			}
			ext.pattern = append(ext.pattern, extensionElem{varType, t.Term})
		default:
			ok = false
			parser.error2(UnexpectedToken, tok.Start, tok.End)
			return
		}
	}

	i++ // move past equals

	ext.fn, ok = parser.Process(standardActions, tokens[i:])

	// TODO: register extension
	return
}

func (parser *Parser) makeExtension(pattern []extensionElem, expression []token.Token) (ext Extension, ok bool) {
	ext.pattern = pattern
	ext.fn, ok = parser.Process(standardActions, expression)
	if !ok {
		return
	}

	// // each keywords/marker (part) should have an index in tokens, even if there are no tokens for that keyword
	// ext.replace = func(self Extension, parser *Parser, tokens [][]token.Token) (term termElem, ok bool) {
	// 	if len(tokens) != len(parts) {
	// 		panic("bad, bad lemon bad") // TODO: actual error message
	// 	}

	// 	term = self.fn
	// 	i := 0
	// 	for i, tokens := range tokens {

	// 		var result termElem
	// 		result, ok = parser.process(standardActions, tokens)
	// 		if !ok {
	// 			return
	// 		}
	// 		term = term.reduce(result)
	// 	}
	// 	return term, true
	// }
	panic("TODO: implement") // TODO
}

// var if_then_else = Extension{
// 	pattern: []ExtensionElem{
// 		MakeMarkerExt("if"),
// 		MakeVarExt("c"),
// 		MakeMarkerExt("then"),
// 		MakeMarkerExt("t"),
// 		MakeMarkerExt("else"),
// 		MakeVarExt("f"),
// 	},
// 	replace: func(self Extension, parser *Parser, tokens [][]token.Token) (term termElem, ok bool) {
// 		if len(tokens) != 3 {
// 			panic("bad, bad lemon bad") // TODO: actual error message
// 		}

// 		term = self.fn
// 		for _, tokens := range tokens {
// 			var result termElem
// 			result, ok = parser.process(standardActions, tokens)
// 			if !ok {
// 				return
// 			}
// 			term = term.reduce(result)
// 		}
// 		return term, true
// 	},
// }

// func (parser *Parser) createExtension()

func (parser *Parser) declare(name token.Token) (setType func(Term, ...uint8), ok bool) {
	requireLink := name.Type == token.Affixed
	var setTypeInit func(Term, ...uint8)
	setTypeInit, ok = parser.declareHelper(name.Value, name.Start, name.End, false)
	if !ok {
		return
	}

	if requireLink {
		var setTypeSecond func(Term, ...uint8)
		setTypeSecond, ok = parser.declareHelper(name.Value, name.Start, name.End, true)
		if !ok {
			return
		}

		setType = func(ty Term, args ...uint8) {
			setTypeInit(ty, args...)   // call first function
			setTypeSecond(ty, args...) // call function for linked name
		}
	} else {
		setType = setTypeInit
	}
	return
}

func (parser *Parser) findTermInTop(name token.Token) (found bool) {
	_, found = parser.declarations.PeekFind(name)
	return
}

func (parser *Parser) lookupTermFromId(id Ident) (decl Declaration, found bool) {
	tok := token.Id.MakeValued(id.Name)
	return parser.lookupTerm(tok)
}

func (parser *Parser) lookupTerm(name token.Token) (decl Declaration, found bool) {
	d, ok := parser.declarations.Find(name)
	found = ok
	if !found {
		return
	}

	decl = *d
	// change value to current occurrence, leaving one in map the same
	decl.name.Start, decl.name.End = name.Start, name.End
	return
}
