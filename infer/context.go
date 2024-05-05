// =================================================================================================
// Alex Peters - February 23, 2024
// =================================================================================================
package infer

import (
	"github.com/petersalex27/yew/common"
	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/types"
)

// Type a = Data a Int
//
// Data = (\x y -> Data x y): a -> Int -> Type a
// type consJudge struct {
// 	forType      types.Polytype
// 	constructors constructorMapType
// }

// type constructorMapType[N nameable.Nameable] map[string]types.TypedJudgment[N, expr.Function, types.Polytype]

// // tries to find constructor named `constructorName` w/in construtor map receiver
// func (constructors consJudge) Find(constructorName N) (constructor types.TypedJudgment[N, expr.Function, types.Polytype], found bool) {
// 	if constructors.constructors == nil {
// 		found = false
// 	} else {
// 		constructor, found = constructors.constructors[constructorName.GetName()]
// 	}
// 	return
// }

// func (constructors consJudge) GetType() types.Polytype {
// 	if constructors.constructors == nil {
// 		panic("bug: constructor map is uninitialized")
// 	}

// 	return constructors.forType
// }

type Context struct {
	reports        []errors.ErrorMessage //[]errorReport
	typeSubs       *common.Table[types.Variable, types.Monotype]
	expressionSubs *common.Table[types.Var, types.Matchable]
	ung            *common.UniqueNameGenerator
	//exprSubs    *table.Table[expr.Referable]
	//consTable   *table.Table[consJudge]
	//syms        *table.Table[Symbol]
	//TypeContext *types.Context
	//ExprContext *expr.Context
}

// convenience method for a type judgment with a new, free type variable; i.e., for an expression e,
//
//	e: newvar
// func (cxt *Context) Judge(e expr.Expression) TypeJudgment {
// 	var newvar types.Type = cxt.TypeContext.NewVar()
// 	return bridge.Judgment(e, newvar)
// }

// type ExportableContext[N nameable.Nameable] struct {
// 	name      N
// 	consTable *table.Table[consJudge]
// 	syms      *table.Table[Symbol]
// }

// func (ecxt *ExportableContext) export(name N, sym Symbol) Status {
// 	// re-exported?
// 	_, ok := ecxt.syms.Get(name)
// 	if ok {
// 		return IllegalShadow
// 	}

// 	// add symbol to table
// 	ecxt.syms.Add(name, sym)
// 	return Ok
// }

// func newConsAndSymsTables[N nameable.Nameable]() (*table.Table[consJudge], *table.Table[Symbol]) {
// 	return table.NewTable[consJudge](), table.NewTable[Symbol]()
// }

// func NewExportableContext[N nameable.Nameable]() *ExportableContext {
// 	cxt := new(ExportableContext)
// 	cxt.consTable, cxt.syms = newConsAndSymsTables()
// 	return cxt
// }

// creates new inf context
func NewContext() *Context {
	cxt := new(Context)
	cxt.typeSubs = common.MakeTable[types.Variable, types.Monotype](8)
	cxt.ung = common.InitUniqueNameGenerator("$", "")
	cxt.expressionSubs = common.MakeTable[types.Var, types.Matchable](8)
	//cxt.consTable, cxt.syms = newConsAndSymsTables()
	//cxt.ExprContext = expr.NewContext()
	//cxt.TypeContext = types.NewContext()
	cxt.reports = []errors.ErrorMessage{}
	return cxt
}

func (cxt *Context) Inst(sigma types.Polytype) types.Monotype {
	var t types.IndexedDependentType = sigma.Bound
	typeVars := sigma.Binders

	// create new type variables
	bindings := make(map[types.Variable]types.Monotype, len(typeVars))
	for _, tv := range typeVars {
		v := cxt.ung.Generate()
		bindings[tv] = types.Variable(v)
	}

	// if d, ok := t.(types.DependentType); ok {
	// 	// replace all bound expression variables w/ new expression variables
	// 	t = d.FreeIndex(cxt.ExprContext)
	// }

	// replace all bound variables w/ newly created type variables
	return t.Bind(bindings)
}

// func NewTestableContext() *Context[nameable.Testable] {
// 	cxt := NewContext[nameable.Testable]()
// 	cxt.TypeContext = cxt.TypeContext.SetNameMaker(nameable.MakeTestable)
// 	cxt.ExprContext = cxt.ExprContext.SetNameMaker(nameable.MakeTestable)
// 	return cxt
// }

// // applies kind and type substitutions to expression and type of judgment respectively
// func (cxt *Context) judgmentSubstitution(judge bridge.JudgmentAsExpression[N, expr.Expression]) bridge.JudgmentAsExpression[N, expr.Expression] {
// 	referable, monotype := GetExpressionAndType[N, expr.Referable, types.Monotype](judge)

// 	var kindSubResult expr.Expression = cxt.GetKindSub(referable)
// 	var typeSubResult types.Type = cxt.GetSub(monotype)

// 	return bridge.Judgment(kindSubResult, typeSubResult)
// }

// applies kind substitutions to `postFindKind`
//
// ASSUMPTION: `postFindKind` is
//
//	cxt.findKindSub(someKind) = postFindKind
// func (cxt *Context) applyKindSubstitutions(postFindKind parser.ExprNode) parser.ExprNode {
// 	data, isData := postFindKind.Substitute(cxt)
// 	if !isData {
// 		return postFindKind
// 	}

// 	memsSubResult := fun.FMap(data.Members, cxt.judgmentSubstitution)
// 	return bridge.MakeData(data.GetTag(), memsSubResult...)
// }

// returns the result of applying all applicable substitutions to `rawKind`.
//
// For example, given substitutions
//
//	> Sub = { n ⟼ 0, k ⟼ Succ n },
//
// and given an input of
//
//	> Succ k
//
// return
//
//	> Succ (Succ 0)
func (cxt *Context) GetKindSub(rawKind types.Matchable) (kind types.Matchable) {
	kind, _ = cxt.findKindSub(rawKind) // returns rawKind if no sub exists
	return kind.Substitute(cxt.expressionSubs)
}

func (cxt *Context) GetSub(m types.Monotype) (out types.Monotype) {
	var found bool

	out, found = cxt.findSub(m)

	if !found {
		out = m
	}

	if function, ok := out.(types.IndexedDependentType); ok {
		out = function.
			SubstitutePattern(cxt.expressionSubs).
			Substitute(cxt.typeSubs)
	}

	return out
}

// first return value is base substitution for `m` (or `m` itself when second return value is false)
//
// second return value is true iff `m` is a variable and `m` has a registered substitution
func (cxt *Context) findSub(m types.Monotype) (out types.Monotype, found bool) {
	found = false
	if nm, ok := m.(types.Variable); ok {
		out, found = cxt.typeSubs.Find(nm)
	}

	if !found {
		out = m
	}

	return
}

// first return value is base substitution for `e` (or `e` itself when second return value is false)
//
// second return value is true iff `e` is a variable and `e` has a registered substitution
func (cxt *Context) findKindSub(e types.Matchable) (out types.Matchable, found bool) {
	found = false
	if v, ok := e.(types.Var); ok {
		out, found = cxt.expressionSubs.Find(v)
	}

	if !found {
		out = e
	}

	return
}

// returns representative for type equiv. class
func (cxt *Context) Find(m types.Monotype) (representative types.Monotype) {
	representative, _ = cxt.findSub(m)
	return
}

// returns representative for kind equiv. class
func (cxt *Context) FindKind(e expr.Referable) (representative expr.Referable) {
	representative, _ = cxt.findKindSub(e)
	return
}
