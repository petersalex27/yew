// =============================================================================
// Author-Date: Alex Peters - November 19, 2023
//
// Content: methods associated w/ Context's constructor table
//
// Notes: -
// =============================================================================
package inf

import (
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/fun"
	"github.com/petersalex27/yew-packages/types"
)

// removes type from context
func (cxt *Context[N]) RemoveType(name N) {
	cxt.consTable.Remove(name)
}

func (cxt *Context[N]) makeConsJudge(sigma types.Polytype[N]) consJudge[N] {
	const wildcardConstructorName string = "_"
	out := consJudge[N]{sigma, make(constructorMapType[N])}
	wildcardConstructor := expr.Bind[N]().In(cxt.ExprContext.NewVar())
	// (\-> $v): σ
	judgment := types.TypedJudge[N](wildcardConstructor, sigma)
	out.constructors[wildcardConstructorName] = judgment
	return out
}

// gets the Yew type referred to by `typeName`
func (cxt *Context[N]) GetReferredType(typeName N) (sigma types.Type[N], stat Status) {
	constructors, found := cxt.consTable.Get(typeName)
	if !found {
		stat = TypeNotDefined
		cxt.appendReport(makeNameReport("Get Referred Type", TypeNotDefined, expr.MakeConst(typeName)))
		return
	}

	sigma, stat = constructors.GetType(), Ok
	return
}

// adds type to context iff type DNE in context
func (cxt *Context[N]) AddType(name N, ty types.Type[N]) Status {
	// attempt to look up existing type
	_, alreadyDefined := cxt.consTable.Get(name)
	if alreadyDefined {
		cxt.appendReport(makeTypeReport("AddType", TypeRedef, ty))
		return TypeRedef
	}

	ty = cxt.Gen(ty)

	constructors := cxt.makeConsJudge(ty.(types.Polytype[N]))

	cxt.consTable.Add(name, constructors)
	return Ok
}

// tries to get constructors for a type with the name that's the value of `typeName`
func (cxt *Context[N]) GetConstructorsForType(typeName N) (constructors consJudge[N], defined bool) {
	constructors, defined = cxt.consTable.Get(typeName)
	return
}

// assumes `returnType` closes (for both val vars and type vars) the type it binds
func (cxt *Context[N]) createFunction(params []types.Monotyped[N], returnType types.Polytype[N]) types.Polytype[N] {
	// binders that close openess of bound type
	binders := returnType.GetBinders()
	// type bound by polytype binders
	unboundReturnType := returnType.GetBound()
	// montype that is possibly open
	freeReturnType := types.GetDependent(unboundReturnType)
	// variables that close open dependencies for freeType
	valDependees := types.GetDependees(unboundReturnType)

	function := fun.FoldRight(
		freeReturnType,
		params,
		func(left, right types.Monotyped[N]) types.Monotyped[N] {
			return cxt.TypeContext.Function(left, right)
		},
	).(types.TypeFunction[N])

	// now, re-close type
	closedDeps := types.Map(valDependees...).To(function) // close free val vars
	return types.Forall(binders...).Bind(closedDeps)      // close open type vars
}

func (cxt *Context[N]) buildConstructor(data bridge.Data[N], sigma types.Polytype[N]) types.TypedJudgment[N, expr.Function[N], types.Polytype[N]] {
	// grab type of each member in data
	memberTypes := fun.FMap(
		data.Members,
		func(member bridge.JudgmentAsExpression[N, expr.Expression[N]]) types.Monotyped[N] {
			t, _ := member.TypeAndExpr()
			return t.(types.Monotyped[N])
		},
	)

	// create constructor params
	constructorParams := cxt.ExprContext.NumNewVars(len(memberTypes))
	// create constructor
	constructor := expr.Bind(constructorParams...).In(data)
	// create kind (constructor's type) for constructor
	constructorType := cxt.createFunction(memberTypes, sigma)

	return types.TypedJudge[N](constructor, constructorType)
}

// creates and adds a constructor of `data` for the type referred to by `typeName`
func (cxt *Context[N]) AddConstructorFor(typeName N, data bridge.Data[N]) Status {
	constructors, typeDefined := cxt.GetConstructorsForType(typeName)
	if !typeDefined {
		cxt.appendReport(makeNameReport("Add Constructor", TypeNotDefined, expr.MakeConst(typeName), data.GetTag()))
		return TypeNotDefined
	}

	_, found := constructors.Find(data.GetTag().Name)
	if found {
		cxt.appendReport(makeNameReport("Add Constructor", TypeNotDefined, expr.MakeConst(typeName), data.GetTag()))
		return ConstructorRedef
	}

	// get Yew type of type with the name that's the value of `typeName`
	yewPolytype := constructors.GetType()
	constructor := cxt.buildConstructor(data, yewPolytype)
	key := data.GetTag().Name.GetName()
	constructors.constructors[key] = constructor
	return Ok
}
