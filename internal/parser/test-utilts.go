//go:build test
// +build test

package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/source"
)

var (
	// tokens (as api.Token interface)
	id_x_tok         api.Token = token.Id.MakeValued("x")               // x
	id_MyId_tok      api.Token = token.Id.MakeValued("MyId")            // MyId
	underscoreTok    api.Token = token.Underscore.Make()                // _
	floatValTok      api.Token = token.FloatValue.MakeValued("1.0")     // 1.0
	charValTok       api.Token = token.CharValue.MakeValued("a")        // 'a'
	stringValTok     api.Token = token.StringValue.MakeValued("abc")    // "abc"
	rawStringValTok  api.Token = token.RawStringValue.MakeValued("abc") // `abc`
	importPathTok    api.Token = token.ImportPath.MakeValued("a/b/c")   // "a/b/c"
	integerValTok    api.Token = token.IntValue.MakeValued("1")         // 1
	nilListTok       api.Token = token.EmptyBracketEnclosure.Make()     // []
	unitTypeTok      api.Token = token.EmptyParenEnclosure.Make()       // ()
	hole_x_tok       api.Token = token.Hole.MakeValued("?x")            // ?x
	hole_MyId_tok    api.Token = token.Hole.MakeValued("?MyId")         // ?MyId
	id_dollar_tok    api.Token = token.Id.MakeValued("$")               // $
	infix_dollar_tok api.Token = token.Infix.MakeValued("($)")          // ($)
	__               api.Token = token.Underscore.Make()                // just a placeholder
	equalTok         api.Token = token.Equal.Make()                     // =
	eraseTok         api.Token = token.Erase.Make()                     // erase
	onceTok          api.Token = token.Once.Make()                      // once
	annotOpenTok     api.Token = token.LeftBracketAt.Make()
	id_test_tok      api.Token = token.Id.MakeValued("test")
	lbracketToken    api.Token = token.LeftBracket.Make()
	rbracket         api.Token = token.RightBracket.Make()
	annot            api.Token = token.Token{Typ: token.FlatAnnotation, Value: "--@test"}

	// tokens
	backslash   = token.Backslash.Make()  // \
	comma       = token.Comma.Make()      // ,
	thickArrow  = token.ThickArrow.Make() // =>
	lbrace      = token.LeftBrace.Make()  // {
	rbrace      = token.RightBrace.Make() // }
	lparen      = token.LeftParen.Make()  // (
	rparen      = token.RightParen.Make() // )
	newline     = token.Newline.Make()    // \n
	colon       = token.Colon.Make()      // :
	equal       = token.Equal.Make()      // =
	colonEqual  = token.ColonEqual.Make() // :=
	arrow       = token.Arrow.Make()      // ->
	bar         = token.Bar.Make()        // |
	let         = token.Let.Make()        // let
	in          = token.In.Make()         // in
	caseTok     = token.Case.Make()       // case
	with        = token.With.Make()       // with
	as          = token.As.Make()         // as
	of          = token.Of.Make()         // of
	derivingTok = token.Deriving.Make()   // deriving
	requiring   = token.Requiring.Make()  // requiring
	spec        = token.Spec.Make()       // spec
	inst        = token.Inst.Make()       // inst
	where       = token.Where.Make()      // where
	alias       = token.Alias.Make()      // alias
	using       = token.Using.Make()      // using
	syntaxTok   = token.Syntax.Make()     // syntax
	public      = token.Public.Make()     // public
	open        = token.Open.Make()       // open
	forall      = token.Forall.Make()     // forall
	moduleTok   = token.Module.Make()     // module
	importTok   = token.Import.Make()     // import

	// nodes
	name_x                = data.EOne[name](id_x_tok)
	name_eq               = data.EOne[name](equalTok)                                                                    // x
	name_MyId             = data.EOne[name](id_MyId_tok)                                                                 // MyId
	name_dollar           = data.EOne[name](id_dollar_tok)                                                               // $
	name_infix_dollar     = data.EOne[name](infix_dollar_tok)                                                            // ($)
	holeNode              = data.EOne[hole](hole_x_tok)                                                                  // ?x
	holePatName           = data.Inl[name](holeNode)                                                                     // ?x
	nilList               = data.Inr[hole](data.EOne[name](nilListTok))                                                  // []
	literalNode           = data.EOne[literal](integerValTok)                                                            // 1
	x_as_lower            = data.EOne[lowerIdent](id_x_tok)                                                              // x
	MyId_as_upper         = data.EOne[upperIdent](id_MyId_tok)                                                           // MyId
	lowerId               = data.Inl[upperIdent](x_as_lower)                                                             // x
	upperId               = data.Inr[lowerIdent](MyId_as_upper)                                                          // MyId
	patternName_x         = data.Inr[hole](name_x)                                                                       // x
	patternName_eq        = data.Inr[hole](name_eq)                                                                      // =
	patternAtomNode       = data.Inr[literal](patternName_x)                                                             // x
	enc                   = data.Singleton[pattern](name_x)                                                              // x
	pattern_x_x           = data.Construct[pattern](name_x, name_x)                                                      // x, x OR x x
	exprAtomNode          = data.Inl[lambdaAbstraction](patternAtomNode)                                                 // x
	exprNode              = expr(name_x)                                                                                 // x
	_exprGroup            = data.Construct[expr](exprNode, exprNode)                                                     // x x
	_patternGroup         = data.Construct[pattern](name_x, name_x)                                                      // x x
	exprAppNode           = expr(data.EMakePair[exprApp](exprNode, data.Singleton(exprNode)))                            // x x
	exprAppNode2          = expr(data.EMakePair[exprApp](exprNode, _exprGroup))                                          // x x x
	patternNode           = pattern(name_x)                                                                              // x
	patternAppNode        = pattern(data.EMakePair[patternApp](patternNode, data.Singleton(patternNode)))                // x x
	patternAppNode2       = pattern(data.EMakePair[patternApp](patternNode, _patternGroup))                              // x x x
	encPattern            = pattern(patternEnclosed{NonEmpty: enc})                                                      // ( x )
	encPattern2           = pattern(patternEnclosed{NonEmpty: pattern_x_x})                                              // (x, x)
	encPatternImplicit    = pattern(patternEnclosed{NonEmpty: enc, implicit: true})                                      // {x}
	encPattern2Implicit   = pattern(patternEnclosed{NonEmpty: pattern_x_x, implicit: true})                              // {x, x}
	wildcardNode          = wildcard{data.One(underscoreTok)}                                                            // _
	lowerBinder           = data.Inl[pattern](lowerId)                                                                   // x
	lambdaBinderNode      = data.EInl[lambdaBinder](lowerBinder)                                                         // x
	lambdaBinders1        = data.EConstruct[lambdaBinders](lambdaBinderNode)                                             // x
	lambdaBinders2        = data.EConstruct[lambdaBinders](lambdaBinderNode, lambdaBinderNode)                           // x, x
	lambdaAbs1            = data.EMakePair[lambdaAbstraction](lambdaBinders1, exprNode)                                  // \x => x
	lambdaAbs2            = data.EMakePair[lambdaAbstraction](lambdaBinders2, exprNode)                                  // \x, x => x
	typ_x                 = typ(name_x)                                                                                  // x
	enclosedExpr          = expr(name_x)                                                                                 // x
	_noAnnots             = data.Nothing[annotations]()                                                                  //
	_noVis                = data.Nothing[visibility]()                                                                   //
	typingPair            = data.MakePair(name_x, typ_x)                                                                 // x : x
	typingNode            = typing{annotations: _noAnnots, visibility: _noVis, typing: typingPair}                       // x : x
	bindingGroupMem_def   = bindingGroupMember(data.Inl[typingMember](data.MakePair(lowerBinder, exprNode)))             // x := x
	bindingGroupMem_typ   = bindingGroupMember(data.Inr[binderMember](data.MakePair(typingNode, data.Nothing[expr]())))  // x : x
	bindingGroupMem_typ_2 = bindingGroupMember(data.Inr[binderMember](data.MakePair(typingNode, data.Just(exprNode))))   // x : x := x
	letBinding_b          = data.EConstruct[letBinding](bindingGroupMem_def)                                             // x := x
	letBinding_t          = data.EConstruct[letBinding](bindingGroupMem_typ)                                             // x : x
	letBinding_a          = data.EConstruct[letBinding](bindingGroupMem_typ_2)                                           // x : x := x
	letBinding_bt         = data.EConstruct[letBinding](bindingGroupMem_def, bindingGroupMem_typ)                        // x := x; x : x
	letExprNode           = data.EMakePair[letExpr](letBinding_b, exprNode)                                              // let x := x in x
	defBodyPossible_x     = data.EMakePair[defBodyPossible](data.Inr[withClause](exprNode), data.Nothing[whereClause]()) // x
	defBody_x             = data.EInr[defBody](defBodyPossible_x)                                                        // x
	caseArmNode           = data.EMakePair[caseArm](pattern(name_x), defBody_x)                                          // x => x
	caseExprNode          = data.EMakePair[caseExpr](pattern(name_x), data.EConstruct[caseArms](caseArmNode))            // case x of x => x

	// annotation nodes

	test_id                = data.Inl[upperIdent](data.EOne[lowerIdent](id_test_tok))                                           // test
	flatAnnot              = data.EOne[flatAnnotation](annot)                                                                   // --@test
	annotSimple            = data.EMakePair[enclosedAnnotation](test_id, data.Nil[api.Node]())                                  // [@test]
	annotSomeContent       = data.EMakePair[enclosedAnnotation](test_id, data.Nil[api.Node](1).Snoc(id_test_tok))               // [@test test]
	annotWithInnerBrackets = data.EMakePair[enclosedAnnotation](test_id, data.Nil[api.Node](2).Append(lbracketToken, rbracket)) // [@test[]]
	annotation_enclosed    = annotation(data.Inr[flatAnnotation](annotSimple))                                                  // [@test]
	annotation_flat        = annotation(data.Inl[enclosedAnnotation](flatAnnot))                                                // --@test
	annotationBlock1       = data.EConstruct[annotations](annotation_flat)                                                      // --@test
	annotationBlock2       = annotations{data.Construct[annotation](annotation_flat).Snoc(annotation_enclosed)}                 // --@test\n[@test]

	// type nodes

	eraseMultiplicity                 = data.EOne[modality](eraseTok)                                                   // erase
	onceMultiplicity                  = data.EOne[modality](onceTok)                                                    // once
	innerTypeTermsSeq                 = data.EConstruct[innerTypeTerms](typ_x, typ_x)                                   // x, x
	innerTypeTermsPairNode            = data.MakePair(data.EConstruct[innerTypeTerms](typ_x), typ_x)                    // x : x
	innerTypeTermsSeqPairNode         = data.MakePair(innerTypeTermsSeq, typ_x)                                         // x, x : x
	innerTypingNode                   = innerTyping{mode: data.Nothing[modality](), typing: innerTypeTermsPairNode}     // x : x
	onceInnerTypingNode               = innerTyping{mode: data.Just(onceMultiplicity), typing: innerTypeTermsPairNode}  // once x : x
	eraseInnerTypingNode              = innerTyping{mode: data.Just(eraseMultiplicity), typing: innerTypeTermsPairNode} // erase x : x
	innerTypingSeqNode                = innerTyping{mode: data.Nothing[modality](), typing: innerTypeTermsSeqPairNode}  // x, x : x
	defaultExprNode                   = data.EOne[defaultExpr](exprNode)                                                // x
	implicitTypingNode_default        = data.EMakePair[implicitTyping](innerTypingNode, defaultExprNode)                // x : x := x
	implicitTypingNode                = innerTypingNode                                                                 // x : x
	implicitTypingNodeSeq_default     = data.EMakePair[implicitTyping](innerTypingSeqNode, defaultExprNode)             // x, x : x := x
	implicitTypingNodeSeq             = innerTypingSeqNode                                                              // x, x : x
	enclosedTypeNode                  = enclosedType{false, typ_x}                                                      // ( x )
	enclosedTypeSeqNode               = enclosedType{false, innerTypeTermsSeq}                                          // ( x, x )
	enclosedTypingNode                = enclosedType{false, innerTypingNode}                                            // ( x : x )
	enclosedOnceTypingNode            = enclosedType{false, onceInnerTypingNode}                                        // ( once x : x )
	enclosedEraseTypingNode           = enclosedType{false, eraseInnerTypingNode}                                       // ( erase x : x )
	enclosedTypingSeqNode             = enclosedType{false, innerTypingSeqNode}                                         // ( x, x : x )
	implicitEnclosedTypeNode          = enclosedType{true, typ_x}                                                       // { x }
	implicitEnclosedTypingNode        = enclosedType{true, innerTypingNode}                                             // { x : x }
	implicitEnclosedOnceTypingNode    = enclosedType{true, onceInnerTypingNode}                                         // { once x : x }
	implicitEnclosedEraseTypingNode   = enclosedType{true, eraseInnerTypingNode}                                        // { erase x : x }
	implicitEnclosedTypingNode_def    = enclosedType{true, implicitTypingNode_default}                                  // { x : x := x }
	implicitEnclosedTypingSeqNode_def = enclosedType{true, implicitTypingNodeSeq_default}                               // { x, x : x := x }
	implicitEnclosedTypeSeqNode       = enclosedType{true, innerTypeTermsSeq}                                           // { x, x }
	implicitEnclosedTypingSeqNode     = enclosedType{true, implicitTypingNodeSeq}                                       // { x, x : x }
	forallTypeNode                    = data.EMakePair[forallType](data.EConstruct[forallBinders](lowerId), typ_x)      // forall x in x

	// header nodes

	module_x                = data.EOne[module](x_as_lower)                                                                  // module x
	abc_path                = data.EOne[importPathIdent](importPathTok)                                                      // "a/b/c"
	hideAbcNode             = data.EMakePair[packageImport](abc_path, data.Nothing[selections]())                            // "a/b/c" using _
	allSelections           = data.Inr[lowerIdent](data.Nothing[data.NonEmpty[name]]())                                      // <no token representation>
	pkgImport_abc           = data.EMakePair[packageImport](abc_path, data.Just[selections](allSelections))                  // "a/b/c"
	imports_abc             = data.EConstruct[importing](pkgImport_abc)                                                      // import "a/b/c"
	importStmtNode          = data.EMakePair[importStatement](data.Nothing[annotations](), imports_abc)                      // import "a/b/c"
	annotationsNode         = data.Just(data.EConstruct[annotations](annotation_enclosed))                                   // [@test]
	annotatedImportStmtNode = data.EMakePair[importStatement](annotationsNode, imports_abc)                                  // [@test] import "a/b/c"
	headerNode              = data.EMakePair[header](data.Just(module_x), data.Nil[importStatement](1).Snoc(importStmtNode)) // module x\nimport "a/b/c"
	as_x                    = data.Inl[data.Maybe[data.NonEmpty[name]]](x_as_lower)                                          // as x
	using_x_x               = data.Inr[lowerIdent](data.Just(data.Construct[name](name_x, name_x)))                          // using (x, x)

	// body nodes

	bodyElemNode = typingNode                                        // x : x
	bodyNode     = body{data.Nil[bodyElement](1).Snoc(bodyElemNode)} // x : x
)

// a very simple function that creates a test source from a list of tokens
//
// this just concatenates the tokens into a string, separating them with ' ' unless the token is
// a newline
func createTestSourceCodeFromTokens(tokens []api.Token) *source.SourceCode {
	b := &strings.Builder{}

	for _, t := range tokens {
		b.WriteString(t.String())
		if !token.Newline.Match(t) {
			b.WriteByte(' ')
		}
	}

	src := util.StringSource(b.String())
	out := source.SourceCode{}.Set(src).(source.SourceCode)
	return &out
}

func initTestParser(input []api.Token) *ParserState {
	scanner := &testScanner{
		tokens: input,
		counter: 0,
		SourceCode: createTestSourceCodeFromTokens(input),
	}

	return &ParserState{state: state{scanner: scanner, tokens: input}}
}

func runResultTest[a api.DescribableNode](p Parser, t *testing.T, want a, fut func(p Parser) data.Either[data.Ers, a]) {
	es, actual, isActual := fut(p).Break()
	if !isActual {
		t.Errorf("expected \n%s\n, got \n%s\n", sprintTree(want), sprintTree(es))
	} else if !equals(actual, want) {
		t.Errorf("expected \n%v\n, got \n%v\n", sprintTree(want), sprintTree(actual))
	}
}

// generate a test for a result-output function (a result-output function is one that returns `either[ers, a]` for some `a`)
func resultOutputFUT[a api.DescribableNode](input []api.Token, want a, fut func(p Parser) data.Either[data.Ers, a]) func(*testing.T) {
	return func(t *testing.T) {
		runResultTest(initTestParser(input), t, want, fut)
	}
}

func resultOutputFUT_endCheck[a api.DescribableNode](input []api.Token, want a, fut func(p Parser) data.Either[data.Ers, a], endsAt int) func(*testing.T) {
	return func(t *testing.T) {
		p := initTestParser(input)

		runResultTest(p, t, want, fut)

		if endsAt < 0 {
			endsAt = len(input) + 1 + endsAt
		}
		if p.state.tokenCounter != endsAt {
			t.Errorf("expected parser's token counter to be at %d, got %d", endsAt, p.state.tokenCounter)
		}
	}
}

func runMaybeOutputTest[a api.DescribableNode](p Parser, t *testing.T, want data.Maybe[a], fut func(p Parser) (*data.Ers, data.Maybe[a])) {
	es, mActual := fut(p)

	if es != nil {
		t.Errorf("expected \n%s\n, got \n%s\n", sprintTree(want), sprintTree(*es))
	} else if !equals(mActual, want) {
		t.Errorf("expected \n%v\n, got \n%v\n", sprintTree(want), sprintTree(mActual))
	}
}

// generate a test for a maybe-output function (a maybe-output function is one that returns `(*ers, maybe[a])` for some `a`)
func maybeOutputFUT[a api.DescribableNode](input []api.Token, want data.Maybe[a], fut func(p Parser) (*data.Ers, data.Maybe[a])) func(t *testing.T) {
	return func(t *testing.T) {
		runMaybeOutputTest(initTestParser(input), t, want, fut)
	}
}

func maybeOutputFUT_endCheck[a api.DescribableNode](input []api.Token, want data.Maybe[a], fut func(p Parser) (*data.Ers, data.Maybe[a]), endsAt int) func(t *testing.T) {
	return func(t *testing.T) {
		p := initTestParser(input)

		runMaybeOutputTest(p, t, want, fut)

		if endsAt < 0 {
			endsAt = len(input) + 1 + endsAt
		}
		if p.state.tokenCounter != endsAt {
			t.Errorf("expected parser's token counter to be at %d, got %d", endsAt, p.state.tokenCounter)
		}
	}
}

func sprintTree(n api.DescribableNode) string {
	b := &strings.Builder{}
	util.PrintTree(b, n)
	return b.String()
}

func equals[n1, n2 api.DescribableNode](x n1, y n2) bool {
	if !x.Type().Match(y) {
		return false
	}
	nameX, childrenX := x.Describe()
	nameY, childrenY := y.Describe()
	if nameX != nameY {
		return false
	}
	if len(childrenX) != len(childrenY) {
		return false
	}

	for i := range childrenX {
		cx, okCx := childrenX[i].(api.DescribableNode)
		cy, okCy := childrenY[i].(api.DescribableNode)
		if !(okCx && okCy) {
			// panic: one or more children cannot be described
			out := "expected DescribableNode"
			if !okCx {
				out = fmt.Sprintf("%s, got %T", out, childrenX[i])
			}
			if !okCy {
				out = fmt.Sprintf("%s, got %T", out, childrenY[i])
			}
			panic(out)
		}
		if !equals(cx, cy) {
			return false
		}
	}
	return true
}
