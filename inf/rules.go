package inf

import (
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
)

// [Var] rule:
//
//			x: σ ∈ 𝚪    t = Inst(σ)
//	   ----------------------- [Var]
//	         𝚪 ⊢ x: t
func (cxt *Context[N]) varBody(x bridge.JudgmentAsExpression[N, expr.Const[N]]) Conclusion[N, expr.Const[N], types.Monotyped[N]] {
	var t types.Monotyped[N]

	tmp, xConst := x.TypeAndExpr()

	// grab polytype
	sigma, ok := tmp.(types.Polytype[N])
	if !ok { // still technically a polytype, just one w/ no zero binders, so make that explicit
		// all types that aren't polytypes, are dependent types, so assertion will pass
		dep, _ := tmp.(types.DependentTyped[N])
		sigma = types.Forall[N]().Bind(dep)
	}

	// replace all bound (including kind-) variables with free variables
	t = cxt.Inst(sigma)
	// return judgment `x: t`
	return Conclude[N](xConst, t)
}

// This is just the "Var" rule but for builtin primitives
func (cxt *Context[N]) Primitive(x bridge.Prim[N]) Conclusion[N, bridge.Prim[N], types.Monotyped[N]] {
	t := x.Val.GetType()
	return Conclude[N](x, t)
}

// [Var] rule:
//
//			x: σ ∈ 𝚪    t = Inst(σ)
//	   ----------------------- [Var]
//	         𝚪 ⊢ x: t
func (cxt *Context[N]) Var(x expr.Const[N]) Conclusion[N, expr.Const[N], types.Monotyped[N]] {
	xJudge, found := cxt.Get(x)
	if !found {
		// `x` is not in the context
		cxt.appendReport(makeNameReport("Var", NameNotInContext, x))
		return CannotConclude[N, expr.Const[N], types.Monotyped[N]](NameNotInContext)
	}

	return cxt.varBody(xJudge)
}

// [App] rule:
//
//			𝚪 ⊢ e0: t0    𝚪 ⊢ e1: t1    t2 = newvar    t0 = t1 -> t2
//	   -------------------------------------------------------- [App]
//			                     𝚪 ⊢ (e0 e1): t2
//
// applies j0 and j1 resulting in a type t2 and the implication that
//
//	t0 = t1 -> t2
//
// the *magic* of this rule comes from the new equation which provides more
// information about type t0
//
// curry-howard: conditional elim
func (cxt *Context[N]) App(j0, j1 TypeJudgment[N]) Conclusion[N, expr.Application[N], types.Monotyped[N]] {
	// split judgments into types and expressions
	e0, tmp0 := j0.GetExpressionAndType()
	e1, tmp1 := j1.GetExpressionAndType()
	// get monotypes
	t0 := tmp0.(types.Monotyped[N])
	t1 := tmp1.(types.Monotyped[N])
	// premise `t2 = newvar`
	t2 := cxt.TypeContext.NewVar()
	// create monotype `t1 -> t2`
	t1_to_t2 := cxt.TypeContext.Function(t1, t2)
	// premise `t0 = t1 -> t2`
	stat := cxt.Unify(t0, t1_to_t2)
	if stat.NotOk() {
		terms := []TypeJudgment[N]{j0, j1}
		report := makeReport("App", stat, terms...)
		cxt.appendReport(report)
		return CannotConclude[N, expr.Application[N], types.Monotyped[N]](stat)
	}
	// "(e0 e1)" in result of rule
	appliedExpression := expr.Apply(e0, e1)
	// (e0 e1): t2
	return Conclude[N](appliedExpression, cxt.GetSub(t2))
}

// [Abs] rule:
//
//	t0 = newvar    𝚪, param: t0 ⊢ e: t1
//	-----------------------------------
//	    𝚪 ⊢ (λparam . e): t0 -> t1
//
// notice that the second param adds context and the third premise no longer
// has that context
//
// curry-howard: conditional intro
func (cxt *Context[N]) Abs(param N) func(TypeJudgment[N]) Conclusion[N, expr.Function[N], types.Monotyped[N]] {
	// first, add context (this is the first premise)
	paramConst := expr.Const[N]{Name: param}
	t0 := cxt.TypeContext.NewVar()
	// grow context w/ type judgment `param: t0`
	cxt.Shadow(paramConst, t0)

	// now, return function to allow second premise of Abs when needed
	return func(j TypeJudgment[N]) Conclusion[N, expr.Function[N], types.Monotyped[N]] {
		// remove context added
		cxt.Remove(paramConst)

		// split judgment
		e, tmp1 := j.GetExpressionAndType()
		t1 := tmp1.(types.Monotyped[N])

		// create function body by converting param-name to param-var in e
		v := cxt.ExprContext.NewVar()
		e = e.BodyAbstract(v, paramConst)

		// actual function creation, finish abstraction of `e`
		f := expr.Bind(v).In(e)

		// create function type
		var fnType types.Monotyped[N] = cxt.TypeContext.Function(t0, t1)

		// last line of rule: `(λparam . e): t0 -> t1`
		return Conclude[N](f, fnType)
	}
}

type letAssumptionDischarge[N nameable.Nameable] func(TypeJudgment[N]) Conclusion[N, expr.NameContext[N], types.Monotyped[N]]

// [Let] rule:
//
//	𝚪 ⊢ e0: t     𝚪, name: Gen(t) ⊢ e1: t1
//	-------------------------------------- [Let]
//	     𝚪 ⊢ let name = e0 in e1: t1
//
// notice that the second param adds context and the third premise no longer
// has that context
//
// This rule allows for a kind of polymorphism:
//
//	𝚪 = {0: Int, (λy.y): a -> a}:
//
//		            [ x: ∀a.a->a ]¹   Inst(∀a.a->a)        0: Int   Int=Inst(Int)
//		            ------------------------------- [Var]  ---------------------- [Var]
//		                                x: v->v                    0: Int        v->v = Int->t0
//		                                ------------------------------------------------------- [App]
//		  (λy.y): a->a   a->a=Inst(a->a)                       (x 0): t0
//		  ------------------------------ [Var]                 ---------- [Id]
//		           (λy.y): a->a                                (x 0): Int
//		         1 ------------------------------------------------------ [Let]
//		                          let x = (λy.y) in x 0: Int
func (cxt *Context) Let(name string, e0 types.Matchable, t0 types.Monotype) letAssumptionDischarge[N] {
	nameConst := expr.Const[N]{Name: name}
	generalized_t0 := cxt.Gen(t0)
	cxt.Shadow(nameConst, generalized_t0)

	return func(j1 TypeJudgment[N]) Conclusion[N, expr.NameContext[N], types.Monotyped[N]] {
		cxt.Remove(nameConst)

		e1, t1 := j1.GetExpressionAndType()
		mono := t1.(types.Monotyped[N])
		let := expr.Let(nameConst, e0, e1)
		return Conclude[N](let, mono)
	}
}

// [Rec] rule:
//
//	𝚪,𝚪ʹ ⊢ e1: t1   ...   𝚪,𝚪ʹ ⊢ eN: tN    𝚪,𝚪ʹʹ ⊢ e0: t0
//	----------------------------------------------------- [Rec]
//	    𝚪 ⊢ rec v1 = e1 and ... and vN = eN in e0: t0
//	where
//	    𝚪ʹ = v1: t1, ..., vN: tN
//	    𝚪ʹʹ = v1: Gen(t1), ..., vN: Gen(tN)
func (cxt *Context[N]) Rec(names []N) func(js []TypeJudgment[N]) func(tj TypeJudgment[N]) Conclusion[N, expr.RecIn[N], types.Monotyped[N]] {
	// non-zero length slice of names
	if len(names) < 1 {
		cxt.appendReport(makeReport[N]("Rec", RecArgsLengthMismatch))
		return func(js []TypeJudgment[N]) func(tj TypeJudgment[N]) Conclusion[N, expr.RecIn[N], types.Monotyped[N]] {
			return func(tj TypeJudgment[N]) Conclusion[N, expr.RecIn[N], types.Monotyped[N]] {
				return CannotConclude[N, expr.RecIn[N], types.Monotyped[N]](RecArgsLengthMismatch)
			}
		}
	}

	vs := cxt.TypeContext.NumNewVars(len(names))
	defs := make([]expr.Def[N], len(names))
	// add 𝚪ʹ to context
	for i, name := range names {
		defs[i] = expr.Declare(name)
		c := defs[i].GetName()
		cxt.Shadow(c, vs[i])
	}

	// function for discharging 𝚪ʹ or 𝚪ʹʹ
	removeNames := func() {
		for _, def := range defs {
			cxt.Remove(def.GetName())
		}
	}

	return func(js []TypeJudgment[N]) func(tj TypeJudgment[N]) Conclusion[N, expr.RecIn[N], types.Monotyped[N]] {
		removeNames() // discharge 𝚪ʹ

		if len(js) != len(names) {
			// report error and return fail fn
			cxt.appendReport(makeReport("Rec", RecArgsLengthMismatch, js...))
			return func(TypeJudgment[N]) Conclusion[N, expr.RecIn[N], types.Monotyped[N]] {
				return CannotConclude[N, expr.RecIn[N], types.Monotyped[N]](RecArgsLengthMismatch)
			}
		}

		// add 𝚪ʹʹ to context
		for i, def := range defs {
			e, t := js[i].GetExpressionAndType()
			m := t.(types.Monotyped[N])
			defs[i] = def.Instantiate(e)
			sigma := cxt.Gen(m) // generalize
			cxt.Shadow(def.GetName(), sigma)
		}

		return func(tj TypeJudgment[N]) Conclusion[N, expr.RecIn[N], types.Monotyped[N]] {
			removeNames() // discharge 𝚪ʹʹ

			e0, t0 := tj.GetExpressionAndType()
			mono := t0.(types.Monotyped[N])
			rec := expr.Rec(defs...)(e0)
			return Conclude[N](rec, mono)
		}
	}
}

func (ecxt *ExportableContext[N]) exportConstructors(cxt *Context[N], typeName N, src consJudge[N], constructorNames []N) bool {
	for _, name := range constructorNames {
		constructor, found := src.Find(name)
		if !found {
			cxt.appendReport(makeNameReport("Export Constructor", UndefinedConstructor, expr.MakeConst(name)))
			return false
		}

		// add constructor
		tab, _ := ecxt.consTable.Get(typeName)
		tab.constructors[name.GetName()] = constructor
	}
	return true
}

func (cxt *Context[N]) exportTypes(typeNames []N, constructorNames [][]N) *ExportableContext[N] {
	out := NewExportableContext[N]()

	for i, typeName := range typeNames {
		cj, ok := cxt.consTable.Get(typeName)
		if !ok {
			cxt.appendReport(makeNameReport("Export Type", UndefinedType, expr.MakeConst(typeName)))
			return nil
		}
		if !out.exportConstructors(cxt, typeName, cj, constructorNames[i]) {
			return nil
		}
	}
	return out
}

func (ecxt *ExportableContext[N]) exportNames(cxt *Context[N], names []N) *ExportableContext[N] {
	for _, name := range names {
		sym, ok := cxt.syms.Get(name)
		if !ok {
			nameConst := expr.MakeConst(name)
			cxt.appendReport(makeNameReport("Export Functions", UndefinedFunction, nameConst))
			return nil
		}
		export, exported := sym.Export()
		if !exported {
			nameConst := expr.MakeConst(name)
			cxt.appendReport(makeNameReport("Export Functions", AmbiguousFunction, nameConst))
			return nil
		}
		ecxt.export(name, *export)
	}
	return ecxt
}

// [Export] rule:
//
//	module M (x0, .., xN)    𝚪 ⊢ x0: σ0   ...   𝚪 ⊢ xN: σN
//	------------------------------------------------------ [Export]
//	              M = { x0: σ0, .., xN: σN }
func Export[N nameable.Nameable](name N, nameMaker func(string) N, names, typeNames []N, constructorNames [][]N) (*Context[N], func() *ExportableContext[N]) {
	// precondition
	if len(typeNames) != len(constructorNames) {
		panic("illegal arguments: len(typeNames) != len(constructorNames)")
	}

	// initialize context
	cxt := NewContext[N]()
	cxt.ExprContext = cxt.ExprContext.SetNameMaker(nameMaker)
	cxt.TypeContext = cxt.TypeContext.SetNameMaker(nameMaker)

	v := cxt.TypeContext.NewVar()
	sigma := cxt.Gen(v) // simga = Gen(v) = forall v . v

	// add all names as any type
	for _, name := range names {
		c := expr.MakeConst(name)
		cxt.Add(c, sigma)
	}

	return cxt, func() *ExportableContext[N] {
		// now, for each type name, export it and its constructors
		ecxt := cxt.exportTypes(typeNames, constructorNames)
		if ecxt == nil {
			return nil
		}
		ecxt.name = name

		return ecxt.exportNames(cxt, names)
	}
}

// [Import] rule:
//
//	M = 𝚪∗    𝚪, 𝚪∗ ⊢ e: t
//	---------------------- [Import]
//	 𝚪 ⊢ import M in e: t
func (cxt *Context[N]) Import(qualification QualificationType, moduleName N, as N) {
	// export := lookup(moduleName)
	// cxt.import(as, export)
	// TODO
}
