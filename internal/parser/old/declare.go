package parser

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/petersalex27/yew/internal/common/table"
	"github.com/petersalex27/yew/internal/token"
	"github.com/petersalex27/yew/internal/types"
)

type declaration struct {
	implicit bool
	*termInfo
	available exports
}

type declMultiTable = table.MultiTable[fmt.Stringer, *declaration]

type declTable = table.Table[fmt.Stringer, *declaration]

type declare_meta struct {
	bp     uint8
	rAssoc bool
}

func fillInfo(decl *declaration, infixed bool, arity uint32, bp int8, rAssoc bool) {
	if arity == 0 {
		*decl.termInfo = termInfo{}
		return
	}

	if !infixed {
		rAssoc = false // regardless of truthiness, prefix IDs cannot be right associative--it doesn't make sense
	}

	*decl.termInfo = termInfo{bp, rAssoc, arity, infixed}
}

func createUseName(name string) string {
	name = strings.TrimPrefix(name, "(")
	name = strings.TrimSuffix(name, ")")

	return name
}

type Str string

func (s Str) String() string {
	return string(s)
}

type generateDecl = func(parser *Parser, typ types.Type, infixed bool, args ...uint8) bool

func generateAssignableTerm(name stringPos, arity uint32) types.Term {
	// create term (constant or lambda) for pattern matching and eventual translation
	s, e := name.Pos()
	C := types.Constant{C: name.String(), Start: s, End: e}
	// make lambda?
	if arity == 0 { // no
		// nothing to apply, just return
		return types.Var(name)
	}

	// yes, make lambda

	// generate free variables
	terms := make([]types.Term, 1, 1+arity)
	vars := make([]types.Variable, arity)
	terms[0] = C
	fmt.Fprintf(os.Stderr, "arity: %d\n", arity)
	for i := 0; i < len(vars); i++ {
		vars[i] = Var(fmt.Sprintf("x%d", i))
		terms = append(terms, vars[i])
	}
	// create application
	app := types.MakeApplication(types.Hole("a"), terms...)
	lambda := types.AutoAbstract(vars, app)
	// assign to environment
	return lambda
}

func addToEnvironment(parser *Parser, name stringPos, typ types.Type, arity uint32) (ok bool) {
	// name : typ
	if !parser.env.Declare(name, typ) {
		parser.transferEnvErrors()
		return false
	}

	// name := term, name : typ
	term := generateAssignableTerm(name, arity)
	if !parser.env.Assign(name, term) {
		parser.transferEnvErrors()
		return false
	}

	return true
}

type stringPos = interface {
	fmt.Stringer
	positioned
}

func generate_setType(name stringPos, decl *declaration) generateDecl {
	return func(parser *Parser, typ types.Type, infixed bool, args ...uint8) bool {
		arity := types.CalcArity(typ)
		var bp uint8 = 0
		rAssoc := false
		if len(args) > 0 {
			bp = args[0]
		} else if arity > 0 {
			bp = 10
		}

		if len(args) > 1 {
			rAssoc = args[1] != 0
		}

		if bp > math.MaxInt8 {
			panic("bug: illegal binding power, cap is 127 from an internal origin and 9 from a source-code origin")
		}

		fillInfo(decl, infixed, arity, int8(bp), rAssoc)

		return true
		//return addToEnvironment(parser, name, typ, arity)
	}
}

type named struct {
	name       string
	start, end int
}

func (n named) String() string {
	return n.name
}

func (n named) Pos() (int, int) {
	return n.start, n.end
}

func (parser *Parser) declareHelper(name string, start, end int, linking, implicit bool, export exports) (setType generateDecl, ok bool) {
	if linking {
		name = createUseName(name)
	}

	nm := Str(name)
	if _, found := parser.declarations.Find(nm); found {
		parser.error2(IllegalRedeclaration, start, end)
		ok = false
		return
	}

	ok = true
	decl := new(declaration)
	*decl = declaration{
		implicit: implicit,
		termInfo: new(termInfo),
		available: export,
	}

	parser.declarations.Map(nm, decl)
	setType = generate_setType(named{name, start, end}, decl)
	return
}

type extensionType = uint8

const (
	keywordType extensionType = iota
	varType
	markerType
)

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
	panic("TODO: implement") // TODO
	/*
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
	*/
}

func (parser *Parser) parseExtension(tokens []token.Token) (ok bool) {
	panic("TODO: implement") // TODO
	/*
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
	*/
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

// bp is arg[0], rAssoc is arg[1]
type setTypeFunc = generateDecl

type exports struct {
	*declTable
	*types.Locals
}

// remove all non-implicit names from exports
func (parser *Parser) processExports(es exports) exports {
	if es.declTable == nil {
		return es
	}

	for _, v := range es.declTable.All() {
		if !v.Value.implicit {
			es.declTable.Delete(v.Key)
			es.Locals.Table.Delete(v.Key)
		}
	}
	return es
}


func (parser *Parser) declare(name token.Token, implicit bool, export exports) (setType generateDecl, ok bool) {
	requireLink := name.Type == token.Infix
	var setTypeInit setTypeFunc
	export = parser.processExports(export)
	setTypeInit, ok = parser.declareHelper(name.Value, name.Start, name.End, false, implicit, export)
	if !ok {
		return
	}

	if requireLink {
		var setTypeSecond setTypeFunc
		setTypeSecond, ok = parser.declareHelper(name.Value, name.Start, name.End, true, implicit, export)
		if !ok {
			return
		}

		setType = func(parser *Parser, ty types.Type, infixed bool, args ...uint8) bool {
			if requireLink && !infixed {
				panic("bug: mismatch in fixedness identity (infix or prefix)")
			}

			ok := setTypeInit(parser, ty, false, args...)        // call first function
			ok = ok && setTypeSecond(parser, ty, infixed, args...) // call function for linked name
			return ok
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
