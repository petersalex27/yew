//go:build test
// +build test

package parser

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/kind"
	"github.com/petersalex27/yew/internal/source"
	"github.com/petersalex27/yew/internal/symbol"
)

// there are obviously lots of ridiculous definitions here, but they are all representative of syntactically valid yew code
//
// most of it will fail type checking though, assuming reasonable definitions paired with the tokens

var (
	// tokens (as api.Token interface)
	id_x_tok         api.Token = token.Id.MakeValued("x")                   // x
	id_MyId_tok      api.Token = token.Id.MakeValued("MyId")                // MyId
	underscoreTok    api.Token = token.Underscore.Make()                    // _
	floatValTok      api.Token = token.FloatValue.MakeValued("1.0")         // 1.0
	charValTok       api.Token = token.CharValue.MakeValued("a")            // 'a'
	stringValTok     api.Token = token.StringValue.MakeValued("abc")        // "abc"
	rawStringValTok  api.Token = token.RawStringValue.MakeValued("abc")     // `abc`
	importPathTok    api.Token = token.ImportPath.MakeValued("a/b/c")       // "a/b/c"
	integerValTok    api.Token = token.IntValue.MakeValued("1")             // 1
	nilListTok       api.Token = token.EmptyBracketEnclosure.Make()         // []
	unitTypeTok      api.Token = token.EmptyParenEnclosure.Make()           // ()
	hole_x_tok       api.Token = token.Hole.MakeValued("?x")                // ?x
	hole_MyId_tok    api.Token = token.Hole.MakeValued("?MyId")             // ?MyId
	id_dollar_tok    api.Token = token.Id.MakeValued("$")                   // $
	infix_dollar_tok api.Token = token.Infix.MakeValued("$")                // ($)
	infix_MyId_tok   api.Token = token.Infix.MakeValued("MyId")             // (MyId)
	infix_x_tok      api.Token = token.Infix.MakeValued("x")                // (x)
	method_run_tok   api.Token = token.MethodSymbol.MakeValued("run")       // (.run)
	raw_my_tok       api.Token = token.RawStringValue.MakeValued("my")      // `my`
	__               api.Token = token.Underscore.Make()                    // just a placeholder
	equalTok         api.Token = token.Equal.Make()                         // =
	eraseTok         api.Token = token.Erase.Make()                         // erase
	onceTok          api.Token = token.Once.Make()                          // once
	annotOpenTok     api.Token = token.LeftBracketAt.Make()                 // [@
	id_test_tok      api.Token = token.Id.MakeValued("test")                // test
	lbracket         api.Token = token.LeftBracket.Make()                   // [
	rbracket         api.Token = token.RightBracket.Make()                  // ]
	annot            api.Token = token.FlatAnnotation.MakeValued("--@test") // doesn't matter, but I think the value should technically just be "test"
	impossibleTok    api.Token = token.Impossible.Make()                    // impossible

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
	auto        = token.Auto.Make()       // auto
	derivingTok = token.Deriving.Make()   // deriving
	requiring   = token.Requiring.Make()  // requiring
	spec        = token.Spec.Make()       // spec
	inst        = token.Inst.Make()       // inst
	from        = token.From.Make()       // from
	where       = token.Where.Make()      // where
	alias       = token.Alias.Make()      // alias
	using       = token.Using.Make()      // using
	syntaxTok   = token.Syntax.Make()     // syntax
	public      = token.Public.Make()     // public
	open        = token.Open.Make()       // open
	forall      = token.Forall.Make()     // forall
	moduleTok   = token.Module.Make()     // module
	importTok   = token.Import.Make()     // import
	dot         = token.Dot.Make()        // .

	// nodes
	name_x                 = data.EOne[name](id_x_tok)
	name_eq                = data.EOne[name](equalTok)                                                                    // x
	name_MyId              = data.EOne[name](id_MyId_tok)                                                                 // MyId
	name_infix_MyId        = data.EOne[name](infix_MyId_tok)                                                              // (MyId)
	name_dollar            = data.EOne[name](id_dollar_tok)                                                               // $
	name_infix_dollar      = data.EOne[name](infix_dollar_tok)                                                            // ($)
	name_method_run        = data.EOne[name](method_run_tok)                                                              // (.run)
	holeNode               = data.EOne[hole](hole_x_tok)                                                                  // ?x
	holePatName            = data.Inl[name](holeNode)                                                                     // ?x
	nilList                = data.Inr[hole](data.EOne[name](nilListTok))                                                  // []
	literalNode            = data.EOne[literal](integerValTok)                                                            // 1
	x_as_lower             = data.EOne[lowerIdent](id_x_tok)                                                              // x
	MyId_as_upper          = data.EOne[upperIdent](id_MyId_tok)                                                           // MyId
	lowerId                = data.Inl[upperIdent](x_as_lower)                                                             // x
	upperId                = data.Inr[lowerIdent](MyId_as_upper)                                                          // MyId
	patternName_x          = data.Inr[hole](name_x)                                                                       // x
	patternName_eq         = data.Inr[hole](name_eq)                                                                      // =
	patternAtomNode        = data.Inr[literal](patternName_x)                                                             // x
	enc                    = data.Singleton[pattern](name_x)                                                              // x
	pattern_x_x            = data.Construct[pattern](name_x, name_x)                                                      // x, x OR x x
	exprAtomNode           = data.Inl[lambdaAbstraction](patternAtomNode)                                                 // x
	exprNode               = expr(name_x)                                                                                 // x
	_exprGroup             = data.Construct(exprNode, exprNode)                                                           // x x
	_patternGroup          = data.Construct[pattern](name_x, name_x)                                                      // x x
	exprAppNode            = expr(data.EMakePair[exprApp](exprNode, data.Singleton(exprNode)))                            // x x
	exprAppNode2           = expr(data.EMakePair[exprApp](exprNode, _exprGroup))                                          // x x x
	exprAppAccess          = expr(data.EMakePair[exprApp](exprNode, data.Singleton[expr](_singleAcc)))                    // x.x
	exprAppAccessDouble    = expr(data.EMakePair[exprApp](exprNode, data.Construct[expr](_accGroup)))                     // x.x.x
	patternNode            = pattern(name_x)                                                                              // x
	patternAppNode         = pattern(data.EMakePair[patternApp](patternNode, data.Singleton(patternNode)))                // x x
	patternAppNode2        = pattern(data.EMakePair[patternApp](patternNode, _patternGroup))                              // x x x
	patternAppAccess       = pattern(data.EMakePair[patternApp](patternNode, data.Singleton[pattern](_singleAcc)))        // x.x
	patternAppAccessDouble = pattern(data.EMakePair[patternApp](patternNode, data.Construct[pattern](_accGroup)))         // x.x.x
	encPattern             = pattern(patternEnclosed{NonEmpty: enc})                                                      // ( x )
	encPattern2            = pattern(patternEnclosed{NonEmpty: pattern_x_x})                                              // (x, x)
	encPatternImplicit     = pattern(patternEnclosed{NonEmpty: enc, implicit: true})                                      // {x}
	encPattern2Implicit    = pattern(patternEnclosed{NonEmpty: pattern_x_x, implicit: true})                              // {x, x}
	wildcardNode           = wildcard{data.One(underscoreTok)}                                                            // _
	lowerBinder            = data.Inl[pattern](lowerId)                                                                   // x
	lambdaBinderNode       = data.EInl[lambdaBinder](lowerBinder)                                                         // x
	lambdaBinders1         = data.EConstruct[lambdaBinders](lambdaBinderNode)                                             // x
	lambdaBinders2         = data.EConstruct[lambdaBinders](lambdaBinderNode, lambdaBinderNode)                           // x, x
	lambdaAbs1             = makeAbstraction(xMap, lambdaBinders1, exprNode)                                              // \x => x
	lambdaAbs2             = makeAbstraction(xMap, lambdaBinders2, exprNode)                                              // \x, x => x
	typ_x                  = typ(name_x)                                                                                  // x
	enclosedExpr           = expr(name_x)                                                                                 // x
	_noAnnots              = data.Nothing[annotations]()                                                                  //
	_noVis                 = data.Nothing[visibility]()                                                                   //
	typingPair             = data.MakePair(name_x, typ_x)                                                                 // x : x
	typingNode             = typing{annotations: _noAnnots, visibility: _noVis, typing: typingPair}                       // x : x
	annotTypingNode        = typing{annotations: data.Just(annotationBlock1), visibility: _noVis, typing: typingPair}     // --@test\nx : x
	upperTypingNode        = typing{annotations: _noAnnots, visibility: _noVis, typing: data.MakePair(name_MyId, typ_x)}  // MyId : x
	bindingGroupMem_def    = bindingGroupMember(data.Inl[typingMember](data.MakePair(lowerBinder, exprNode)))             // x := x
	bindingGroupMem_typ    = bindingGroupMember(data.Inr[binderMember](data.MakePair(typingNode, data.Nothing[expr]())))  // x : x
	bindingGroupMem_typ_2  = bindingGroupMember(data.Inr[binderMember](data.MakePair(typingNode, data.Just(exprNode))))   // x : x := x
	letBinding_b           = data.EConstruct[letBinding](bindingGroupMem_def)                                             // x := x
	letBinding_t           = data.EConstruct[letBinding](bindingGroupMem_typ)                                             // x : x
	letBinding_a           = data.EConstruct[letBinding](bindingGroupMem_typ_2)                                           // x : x := x
	letBinding_bt          = data.EConstruct[letBinding](bindingGroupMem_def, bindingGroupMem_typ)                        // x := x; x : x
	letExprNode            = data.EMakePair[letExpr](letBinding_b, exprNode)                                              // let x := x in x
	defBodyPossible_x      = data.EMakePair[defBodyPossible](data.Inr[withClause](exprNode), data.Nothing[whereClause]()) // x
	defBody_x              = data.EInr[defBody](defBodyPossible_x)                                                        // x
	caseArmNode            = makeCaseArm(xMap, pattern(name_x), defBody_x)                                                // x => x
	caseExprNode           = data.EMakePair[caseExpr](pattern(name_x), data.EConstruct[caseArms](caseArmNode))            // case x of x => x

	// annotation nodes

	test_id                = data.Inl[upperIdent](data.EOne[lowerIdent](id_test_tok))                                      // test
	flatAnnot              = data.EOne[flatAnnotation](annot)                                                              // --@test
	annotSimple            = data.EMakePair[enclosedAnnotation](test_id, data.Nil[api.Node]())                             // [@test]
	annotSomeContent       = data.EMakePair[enclosedAnnotation](test_id, data.Nil[api.Node](1).Snoc(id_test_tok))          // [@test test]
	annotWithInnerBrackets = data.EMakePair[enclosedAnnotation](test_id, data.Nil[api.Node](2).Append(lbracket, rbracket)) // [@test[]]
	annotation_enclosed    = annotation(data.Inr[flatAnnotation](annotSimple))                                             // [@test]
	annotation_flat        = annotation(data.Inl[enclosedAnnotation](flatAnnot))                                           // --@test
	annotationBlock1       = data.EConstruct[annotations](annotation_flat)                                                 // --@test
	annotationBlock2       = annotations{data.Construct[annotation](annotation_flat).Snoc(annotation_enclosed)}            // --@test\n[@test]

	// type nodes

	_accGroup                         = data.EConstruct[accessList](name_x, name_x)                                // .x.x
	_singleAcc                        = data.EConstruct[accessList](name_x)                                        // .x
	appTypeNode                       = typ(data.EMakePair[appType](typ_x, data.Singleton(typ_x)))                 // x x
	appTypeAccessNode                 = typ(data.EMakePair[appType](typ_x, data.Singleton[typ](_singleAcc)))       // x.x
	appTypeAccessDoubleNode           = typ(data.EMakePair[appType](typ_x, data.Singleton[typ](_accGroup)))        // x.x.x
	eraseMultiplicity                 = data.EOne[modality](eraseTok)                                              // erase
	onceMultiplicity                  = data.EOne[modality](onceTok)                                               // once
	innerTypeTermsSeq                 = data.EConstruct[innerTypeSeq](typ_x, typ_x)                                // x, x
	innerTypeTermsPairNode            = data.MakePair(data.EConstruct[innerTypeSeq](typ_x), typ_x)                 // x : x
	innerTypeTermsSeqPairNode         = data.MakePair(innerTypeTermsSeq, typ_x)                                    // x, x : x
	innerNamesNoMode                  = data.Inr[data.Pair[modality, name]](data.Construct(name_x))                // x
	innerNamesSeq                     = data.Inr[data.Pair[modality, name]](data.Construct(name_x, name_x))        // x, x
	innerNamesErase                   = data.Inl[data.NonEmpty[name]](data.MakePair(eraseMultiplicity, name_x))    // erase x
	innerNamesOnce                    = data.Inl[data.NonEmpty[name]](data.MakePair(onceMultiplicity, name_x))     // once x
	innerTypingNode                   = data.EMakePair[innerTyping](innerNamesNoMode, typ_x)                       // x : x
	onceInnerTypingNode               = data.EMakePair[innerTyping](innerNamesOnce, typ_x)                         // once x : x
	eraseInnerTypingNode              = data.EMakePair[innerTyping](innerNamesErase, typ_x)                        // erase x : x
	innerTypingSeqNode                = data.EMakePair[innerTyping](innerNamesSeq, typ_x)                          // x, x : x
	defaultExprNode                   = data.EOne[defaultExpr](exprNode)                                           // x
	implicitTypingNode_default        = data.EMakePair[implicitTyping](innerTypingNode, defaultExprNode)           // x : x := x
	implicitTypingNode                = innerTypingNode                                                            // x : x
	implicitTypingNodeSeq_default     = data.EMakePair[implicitTyping](innerTypingSeqNode, defaultExprNode)        // x, x : x := x
	implicitTypingNodeSeq             = innerTypingSeqNode                                                         // x, x : x
	enclosedTypeNode                  = enclosedType{false, typ_x}                                                 // ( x )
	enclosedTypeSeqNode               = enclosedType{false, innerTypeTermsSeq}                                     // ( x, x )
	enclosedTypingNode                = enclosedType{false, innerTypingNode}                                       // ( x : x )
	enclosedOnceTypingNode            = enclosedType{false, onceInnerTypingNode}                                   // ( once x : x )
	enclosedEraseTypingNode           = enclosedType{false, eraseInnerTypingNode}                                  // ( erase x : x )
	enclosedTypingSeqNode             = enclosedType{false, innerTypingSeqNode}                                    // ( x, x : x )
	implicitEnclosedTypeNode          = enclosedType{true, typ_x}                                                  // { x }
	implicitEnclosedTypingNode        = enclosedType{true, innerTypingNode}                                        // { x : x }
	implicitEnclosedOnceTypingNode    = enclosedType{true, onceInnerTypingNode}                                    // { once x : x }
	implicitEnclosedEraseTypingNode   = enclosedType{true, eraseInnerTypingNode}                                   // { erase x : x }
	implicitEnclosedTypingNode_def    = enclosedType{true, implicitTypingNode_default}                             // { x : x := x }
	implicitEnclosedTypingSeqNode_def = enclosedType{true, implicitTypingNodeSeq_default}                          // { x, x : x := x }
	implicitEnclosedTypeSeqNode       = enclosedType{true, innerTypeTermsSeq}                                      // { x, x }
	implicitEnclosedTypingSeqNode     = enclosedType{true, implicitTypingNodeSeq}                                  // { x, x : x }
	forallTypeNode                    = data.EMakePair[forallType](data.EConstruct[forallBinders](lowerId), typ_x) // forall x in x

	// header nodes

	emptyAnnots             = data.Nothing[annotations]()                                                                    // <no token representation>
	module_x                = module{annotations: emptyAnnots, name: data.One(x_as_lower)}                                   // module x
	module_annot_x          = module{annotations: data.Just(annotationBlock1), name: data.One(x_as_lower)}                   // --@test\nmodule x
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

	// spec inst/def nodes

	constrainerParam_x = data.MakePair[ident, data.Maybe[typ]](lowerId, data.Nothing[typ]())                                // x
	constrainerNode    = data.EMakePair[constrainer](MyId_as_upper, data.EConstruct[constrainerParams](constrainerParam_x)) // MyId x
	vConstraintNode    = data.EConstruct[constraint](data.MakePair(data.Nil[upperIdent](), constrainerNode))                // MyId x
	vConstraint2Node   = data.EConstruct[constraint](data.MakePair(data.Makes(MyId_as_upper), constrainerNode))             // MyId, MyId x
	specDefBodyNode    = data.EConstruct[specBody](data.Inr[def](typingNode))                                               // where x : x
	specInstBodyNode   = data.EConstruct[specBody](data.Inl[typing](defNode))                                               // where x = x
	requiringNode      = data.Construct(defNode)                                                                            // requiring x = x
	dependencyNode     = pattern(name_x)                                                                                    // from x
	specInstTargeted   = makeSpecInst(specHeadNode, data.Just(constrainerNode), specInstBodyNode)                           // spec MyId x = MyId x where x = x

	// body nodes

	emptyVisibility         = data.Nothing[visibility]()                                                                         // <no token representation>
	defBodyNode             = data.EInr[defBody](defBodyPossible_x)                                                              // x
	defBodyImpossible       = data.EInl[defBody](data.EOne[impossible](impossibleTok))                                           // impossible
	defNode                 = def{emptyAnnots, name_x, defBodyNode, api.ZeroPosition()}                                          // x = x
	whereClauseNode         = makeWhereClause(xMap, data.Construct(defNode.asMainElement()))                                     // where x = x
	defBodyPossible_where   = data.EMakePair[defBodyPossible](data.Inr[withClause](exprNode), data.Just(whereClauseNode))        // x where x = x
	defBodyWhereNode        = data.EInr[defBody](defBodyPossible_where)                                                          // x where x = x
	defImpossibleNode       = def{emptyAnnots, name_x, defBodyImpossible, api.ZeroPosition()}                                    // x impossible
	specHeadNode            = data.EMakePair[specHead](data.Nothing[constraint](), constrainerNode)                              // MyId x
	specHeadConstrNode      = data.EMakePair[specHead](data.Just(vConstraintNode), constrainerNode)                              // MyId x => MyId x
	specDefNode             = makeSpecDef(specHeadNode, data.Nothing[pattern](), specDefBodyNode, data.Nothing[specRequiring]()) // spec MyId x where x : x
	specInstNode            = makeSpecInst(specHeadNode, data.Nothing[constrainer](), specInstBodyNode)                          // spec MyId x where x = x
	rawStringNode           = data.EOne[rawString](raw_my_tok)                                                                   // `my`
	rawKeyNode              = data.EOne[syntaxRawKeyword](rawStringNode)                                                         // `my`
	rawSym                  = data.Inr[syntaxRuleIdent](rawKeyNode)                                                              // `my`
	bindingIdSymNode        = data.Inl[syntaxRawKeyword](makeBindingSyntaxRuleIdent(lowerId))                                    // {x}
	idSymNode               = data.Inl[syntaxRawKeyword](makeStdSyntaxRuleIdent(lowerId))                                        // x
	syntaxRuleNode          = data.EConstruct[syntaxRule](rawSym, idSymNode)                                                     // `my` x
	syntaxNode              = makeSyntax(syntaxRuleNode, expr(name_x))                                                           // syntax `my` x = x
	aliasNode               = makeAlias(name_MyId, name_MyId)                                                                    // alias MyId = MyId
	typeConsNode            = makeCons(name_MyId, typ_x)                                                                         // MyId : x
	_consGroup              = data.Construct(typeConsNode)                                                                       // MyId : x
	typeDefNode             = makeTypeDef(upperTypingNode, data.Inl[impossible](_consGroup), data.Nothing[deriving]())           // MyId : x where MyId : x
	body_typing             = body{data.Makes(typingNode.asBodyElement())}                                                       // x : x
	body_def                = body{data.Makes(defNode.asBodyElement())}                                                          // x = x
	body_defImpossible      = body{data.Makes(defImpossibleNode.asBodyElement())}                                                // x impossible
	body_specDef            = body{data.Makes(specDefNode.asBodyElement())}                                                      // spec MyId x where x : x
	body_specInst           = body{data.Makes(specInstNode.asBodyElement())}                                                     // inst MyId x where x = x
	body_syntax             = body{data.Makes(syntaxNode.asBodyElement())}                                                       // syntax `my` x = x
	body_alias              = body{data.Makes(aliasNode.asBodyElement())}                                                        // alias MyId = MyId
	body_typeDef            = body{data.Makes(typeDefNode.asBodyElement())}                                                      // type x = x
	body_annotTyping        = body{data.Makes(annotTypingNode.asBodyElement())}                                                  // --@test\nx : x
	singleConsNode          = data.Construct(typeConsNode)                                                                       // MyId : x
	multiConsNode           = data.Construct(typeConsNode, typeConsNode)                                                         // MyId, MyId : x
	impossibleNode          = data.EOne[impossible](impossibleTok)                                                               // impossible
	derivingNode            = data.EConstruct[deriving](constrainerNode)                                                         // deriving MyId x
	derivingNode2           = data.EConstruct[deriving](constrainerNode, constrainerNode)                                        // deriving (MyId x, MyId x)
	typeDefNodeWithDeriving = makeTypeDef(upperTypingNode, data.Inl[impossible](singleConsNode), data.Just(derivingNode))        // MyId : x where MyId : x deriving MyId x
	withArmLhsNode          = data.Inl[data.Pair[pattern, pattern]](pattern(name_x))                                             // x
	withArmLhsVRNode        = makeWithArmLhsRefined(pattern(name_x), pattern(name_x))                                            // x | x
	withClauseArmNode       = makeWithClauseArm(xMap, withArmLhsNode, defBodyNode)                                               // x => x
	withClauseVRNode        = makeWithClauseArm(xMap, withArmLhsVRNode, defBodyNode)                                             // x | x => x
	withClauseArmsNode      = data.EConstruct[withClauseArms](withClauseArmNode)                                                 // x => x
	withClauseNode          = makeWithClause(pattern(name_x), withClauseArmsNode)                                                // with x of x => x
	xMap                    = map[string]kind.Symbol{"x": kind.Sym2("x", symbol.Var2("~a", api.ZeroPosition()))}                 // "x" |-> (x : ~a)
	withClauseArmNodeWhere  = makeWithClauseArm(xMap, withArmLhsNode, defBodyWhereNode)                                          // x => x where x = x
	withClauseArmsNodeWhere = data.EConstruct[withClauseArms](withClauseArmNodeWhere)                                            // x => x where x = x
	withClauseNodeWhere     = makeWithClause(pattern(name_x), withClauseArmsNodeWhere)                                           // with x of x => x where x = x
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
		tokens:     input,
		counter:    0,
		SourceCode: createTestSourceCodeFromTokens(input),
	}

	return &ParserState{state: state{scanner: scanner, tokens: input, solver: symbol.Solver(symbol.New())}}
}

func runResultTest[a api.DescribableNode](p parser, t *testing.T, want a, fut func(p parser) data.Either[data.Ers, a]) {
	es, actual, isActual := fut(p).Break()
	if !isActual {
		printErrors(parseErrors(p, es)...)
		t.Errorf("expected \n%s\n, got \n%s\n", sprintTree(want), sprintTree(es))
	} else if !equals(actual, want) {
		t.Errorf("expected \n%v\n, got \n%v\n", sprintTree(want), sprintTree(actual))
	}
}

func eitherOutputFUT[a, b api.DescribableNode](input []api.Token, want data.Either[a, b], fut func(p parser) data.Either[a, b]) func(*testing.T) {
	return func(t *testing.T) {
		p := initTestParser(input)
		got := fut(p)
		if !equals(got, want) {
			t.Errorf("expected \n%s\n, got \n%s\n", sprintTree(want), sprintTree(got))
		}
	}
}

// generate a test for a result-output function (a result-output function is one that returns `either[ers, a]` for some `a`)
func resultOutputFUT[a api.DescribableNode](input []api.Token, want a, fut func(p parser) data.Either[data.Ers, a]) func(*testing.T) {
	return func(t *testing.T) {
		runResultTest(initTestParser(input), t, want, fut)
	}
}

func resultOutputFUT_endCheck[a api.DescribableNode](input []api.Token, want a, fut func(p parser) data.Either[data.Ers, a], endsAt int) func(*testing.T) {
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

func runMaybeOutputTest[a api.DescribableNode](p parser, t *testing.T, want data.Maybe[a], fut func(p parser) (*data.Ers, data.Maybe[a])) {
	es, mActual := fut(p)

	if es != nil {
		printErrors(parseErrors(p, *es)...)
		t.Errorf("expected \n%s\n, got \n%s\n", sprintTree(want), sprintTree(*es))
	} else if !equals(mActual, want) {
		t.Errorf("expected \n%v\n, got \n%v\n", sprintTree(want), sprintTree(mActual))
	}
}

// generate a test for a maybe-output function (a maybe-output function is one that returns `(*ers, maybe[a])` for some `a`)
func maybeOutputFUT[a api.DescribableNode](input []api.Token, want data.Maybe[a], fut func(p parser) (*data.Ers, data.Maybe[a])) func(t *testing.T) {
	return func(t *testing.T) {
		runMaybeOutputTest(initTestParser(input), t, want, fut)
	}
}

func maybeOutputFUT_endCheck[a api.DescribableNode](input []api.Token, want data.Maybe[a], fut func(p parser) (*data.Ers, data.Maybe[a]), endsAt int) func(t *testing.T) {
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
	type entry interface {
		equals(te2 tableEntry) bool
		symbols() tableEntry
	}

	if !x.Type().Match(y) {
		return false
	}

	if te, ok := api.DescribableNode(x).(entry); ok {
		te2, ok := api.DescribableNode(y).(entry)
		if !ok {
			return false
		}

		ok = te.equals(te2.symbols())
		if !ok {
			t1 := map[string]symbol.Symbol(te.symbols())
			t2 := map[string]symbol.Symbol(te2.symbols())
			fmt.Fprintf(os.Stderr, "expected table\n%v\n, got \n%v\n", t1, t2)
			return false
		}
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

func (te tableEntry) equals(te2 tableEntry) bool {
	if len(te) != len(te2) {
		return false
	}

	for k, v := range te {
		if v2, ok := te2[k]; !ok {
			return false
		} else if v.String() != v2.String() { // TODO: this is definitely not the best way to compare symbols
			return false
		}
	}

	return true
}

func (te tableEntry) symbols() tableEntry {
	return te
}

func printErrors(es ...error) {
	for _, e := range es {
		fmt.Fprintf(os.Stderr, "%s\n", e.Error())
	}
}
