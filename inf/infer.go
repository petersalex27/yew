// =============================================================================
// Author-Date: Alex Peters - 2023
//
// Content:
// contains type inference rules
//
// Notes: -
// =============================================================================
package inf

import "github.com/petersalex27/yew/types"

// check if variable v occurs in monotype t. If it does, return true; else, return false.
// if t = v, then v is not in t, v is t
func (cxt *Context) occurs(v types.Variable, t types.Monotype) (vOccursInT bool) {
	if IsVariable(t) {
		return false
	}

	us := types.FreeVariables(t)
	for _, u := range us {
		if v.String() == u.String() {
			return true
		}
	}

	return false
}

func getVars(e types.Matchable) []types.Var {
	m := map[types.Var]types.Monotype{}
	e.Vars(m, map[types.Variable][]types.Monotype{}, nil)
	vs := make([]types.Var, 0, len(m))
	for v := range m {
		vs = append(vs, v)
	}
	return vs
}

// check if variable v occurs in expression e. If it does, return true; else, return false.
// if e = v, then v is not in e: v is e
func (cxt *Context) kindOccurs(v types.Var, e types.Matchable) (vOccursInE bool) {
	if _, isVar := e.(types.Var); isVar {
		return false
	}

	vs := getVars(e)
	for _, u := range vs {
		if v.String() == u.String() {
			return true
		}
	}

	return false
}

// declares, for types variable v, monotype t:
//
//	v = t
//
// if v ∈ t, then union returns `OccursCheckFailed`; else, skipUnify is returned
func (cxt *Context) union(v types.Variable, t types.Monotype) Status {
	if cxt.occurs(v, t) {
		return OccursCheckFailed
	}

	// if d, ok := t.(types.DependentTypeInstance[T]); ok {
	// 	t = cxt.reindex(d)
	// }

	cxt.typeSubs.Map(v, t)
	return skipUnify
}

// generalizes a monotype into a dependent type
func DependentGeneralization(ty types.Type) types.DependentType {
	if m, ok := ty.(types.Monotype); ok {
		return types.Π()(m)
	}
	if _, ok := ty.(types.Polytype); ok {
		panic("illegal argument")
	}
	return ty.(types.DependentType)
}

// generalizes a type: binds all free variables w/in monotype
func (cxt *Context) Gen(ty types.Type) types.Polytype {
	if t, ok := ty.(types.Monotype); ok {
		// DependentGeneralization(`(t a0 .. aK; x0 .. xN)`) = `mapval (x0: X0) .. (xN: XN) . (t a0 .. aK)`
		dep := DependentGeneralization(t)
		// (t a0 .. aK; x0 .. xN) -> a0 .. aK
		vs := types.FreeVariables(t)
		// forall a0 .. aK . mapval (x0: X0) .. (xN: XN) . (t a0 .. aK)
		return types.Forall(vs...)(dep)
	}
	return ty.(types.Polytype)
}

func IsVariable(ty types.Monotype) bool {
	_, ok := ty.(types.Variable)
	return ok
}

func checkStatus(c0, c1 string, ms0, ms1 []types.Monotype, isA, isB []types.Matchable) Status {
	if c0 != c1 {
		return ConstantMismatch
	}

	if len(ms0) != len(ms1) {
		return ParamLengthMismatch
	}

	if len(isA) != len(isB) {
		return IndexLengthMismatch
	}

	return Ok
}

func Split(m types.Monotype) (c string, params []types.Monotype, indexes []types.Matchable) {
	st, isTypeFunc := m.(types.IndexedDependentType)
	if !isTypeFunc {
		c, params, indexes = m.String(), nil, nil
		if app, ok := m.(types.Application); ok {
			params = app.Args
		}
	} else {
		app := st.Indexer
		indexes = st.Indexes
		c = app.String()
		params = app.Args
	}
	return
}

func SplitKind(kind types.Matchable) (c string, mems []types.Matchable) {
	if app, ok := kind.(types.App); ok {
		mems = make([]types.Matchable, 0, len(app))
		for _, p := range app {
			mems = append(mems, p)
		}
	}

	c = kind.Name()
	return
}

func (cxt *Context) kindUnion(v types.Var, e types.Matchable) Status {
	if cxt.kindOccurs(v, e) {
		return OccursCheckFailed
	}

	cxt.expressionSubs.Map(v, e)
	return skipUnify
}

func checkKindStatus(ca, cb string, memsOfA, memsOfB []types.Matchable) Status {
	if ca != cb {
		return KindConstantMismatch
	}

	if len(memsOfA) != len(memsOfB) {
		return MemsLengthMismatch
	}

	return Ok
}

func (cxt otherwiseDo) otherwiseUnifyKind(a, b types.Matchable) Status {
	// function pre-condition: substitution already happened or occurs-check
	// failed?
	if cxt.stat.NotOk() {
		return fixSkip(cxt.stat)
	}

	// get constants and mems
	ca, memsOfA := SplitKind(a)
	cb, memsOfB := SplitKind(b)

	// check if alright to use in loop
	stat := checkKindStatus(ca, cb, memsOfA, memsOfB)

	// it. through all mems while stat is ok, unifying mems
	for i := 0; stat.IsOk() && i < len(memsOfA); i++ {
		ma, mb := memsOfA[i], memsOfB[i]
		stat = cxt.UnifyKind(ma, mb)
	}

	return stat
}

// tries to creates a substitution from a variable to a monotype
func (cxt *Context) substituteKind(ea, eb types.Matchable) otherwiseDo {
	stat := Ok

	if v, ok := ea.(types.Var); ok {
		stat = cxt.kindUnion(v, eb)
	} else if v, ok := eb.(types.Var); ok {
		stat = cxt.kindUnion(v, ea)
	}

	return otherwiseDo{stat, cxt}
}

func (cxt *Context) UnifyKind(a, b types.Matchable) Status {
	ea := cxt.FindKind(a)
	eb := cxt.FindKind(b)

	return cxt.substituteKind(ea, eb).otherwiseUnifyKind(ea, eb)
}

func (cxt *Context) UnifyIndex(indexOfA, indexOfB types.Matchable) Status {
	// split type and expression
	//ea, ta := indexOfA.AsTypeJudgment().GetExpressionAndType()
	//eb, tb := indexOfB.AsTypeJudgment().GetExpressionAndType()

	// type assertions monotype
	//ma := ta.(types.Monotype)
	//mb := tb.(types.Monotype)
	// type assertions referable
	//ra := ea.(expr)
	//rb := eb.(expr.Referable)

	// unify types
	//stat := cxt.Unify(ma, mb)

	// if type union Ok, then unify kinds
	//if stat.IsOk() {
	//	stat = cxt.UnifyKind(ra, rb)
	//}

	//return stat
	return cxt.UnifyKind(indexOfA, indexOfB)
}

// returns Ok iff (stat.IsOk() || stat.Is(skipUnify))
func fixSkip(stat Status) Status {
	if stat.Is(skipUnify) {
		return Ok
	}
	return stat
}

func (cxt otherwiseDo) otherwiseUnify(a, b types.Monotype) Status {
	// function pre-condition: substitution already happened or occurs-check
	// failed?
	if cxt.stat.NotOk() {
		return fixSkip(cxt.stat)
	}

	// get constants, params, and indexes
	ca, paramsOfA, indexesOfA := Split(a)
	cb, paramsOfB, indexesOfB := Split(b)

	// check if alright to use in loop
	stat := checkStatus(ca, cb, paramsOfA, paramsOfB, indexesOfA, indexesOfB)

	// it. through all params while stat is ok, unifying params
	for i := 0; stat.IsOk() && i < len(paramsOfA); i++ {
		pa, pb := paramsOfA[i], paramsOfB[i]
		stat = cxt.Unify(pa, pb)
	}

	// it. through all indexes while stat is ok, unifying indexes
	for i := 0; stat.IsOk() && i < len(indexesOfA); i++ {
		ia, ib := indexesOfA[i], indexesOfB[i]
		stat = cxt.UnifyIndex(ia, ib)
	}

	return stat
}

type otherwiseDo struct {
	stat Status
	*Context
}

// tries to creates a substitution from a variable to a monotype
func (cxt *Context) substitute(ta, tb types.Monotype) otherwiseDo {
	stat := Ok

	if v, ok := ta.(types.Variable); ok {
		stat = cxt.union(v, tb)
	} else if v, ok := tb.(types.Variable); ok {
		stat = cxt.union(v, ta)
	}

	return otherwiseDo{stat, cxt}
}

func ()

// unifies two monotypes a, b
func (cxt *Context) Unify(a, b types.Monotype) Status {
	ta := cxt.Find(a)
	tb := cxt.Find(b)

	return cxt.substitute(ta, tb).otherwiseUnify(ta, tb)
}
