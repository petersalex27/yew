// =================================================================================================
// Alex Peters - February 23, 2024
// =================================================================================================
package infer


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

// declares, for types variable v, monotype t:
//
//	v = t
//
// if v ∈ t, then union returns `OccursCheckFailed`; else, skipUnify is returned
func (cxt *Context) union(v types.Variable, t types.Monotype) types.Trilean {
	if cxt.occurs(v, t) {
		return types.False
	}

	// if d, ok := t.(types.DependentTypeInstance[T]); ok {
	// 	t = cxt.reindex(d)
	// }

	cxt.typeSubs.Map(v, t)
	return types.Undecided
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

func checkStatus(c0, c1 string, ms0, ms1 []types.Monotype, isA, isB []types.Matchable) bool {
	if c0 != c1 {
		return false
	}

	if len(ms0) != len(ms1) {
		return false
	}

	if len(isA) != len(isB) {
		return false
	}

	return true
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

// unifies two monotypes a, b
func (cxt *Context) Unify(a, b types.Monotype) Status {
	ta := cxt.Find(a)
	tb := cxt.Find(b)

	return cxt.substitute(ta, tb).otherwiseUnify(ta, tb)
}
