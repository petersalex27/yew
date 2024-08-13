// =================================================================================================
// Alex Peters - February 29, 2024
//
// =================================================================================================
package types

import (
	"fmt"
	"github.com/petersalex27/yew/common"
)

func (env *Environment) Rule(s, t Type) (u Sort, ok bool) {
	if s, ok := s.(Universe); ok {
		if t, ok := t.(Universe); ok {
			u, ok := env.system.Rule(s, t)
			if ok {
				debug_log_Rule(s, t, u)
				return u, true
			}
			env.ruleError2(s, t)
			return nil, false
		}
	}
	// if both aren't universes, then return hole
	u, ok = env.NextKindHole(), true
	debug_log_Rule(s, t, u)
	return
}

// a special kind of abstraction that creates an implicit term in head position
func (env *Environment) Constrain(constraint Type, constrained Type) (Pi, bool) {
	var x Variable
	var A Type
	var ok bool
	if x, ok = constraint.(Variable); ok {
		A = GetKind(&x)
	} else {
		A = constraint
		x = AsTyping(constraint)
	}

	kA := GetKind(&A)
	kc := GetKind(&constrained)
	s, ok := env.Rule(kA, kc)
	if !ok {
		return Pi{}, false
	}

	x.Kind = A
	start, end := calcStartEnd(x, A, constrained)
	pi := Pi{
		implicit:  true,
		binderVar: x,
		dependent: constrained,
		kind:      s,
		Start:     start,
		End:       end,
	}

	debug_log_Constrain(constraint, constrained, pi)
	return pi, true
}

// assumes that ty is the correct type for the term tr
// func (env *Environment) SpecializeImplicit(tr Term, ty Type) (Term, Type) {
// 	for pi, ok := ty.(Pi); ok && pi.implicit; {
// 		v := env.NextTermHole()
// 		justType := pi.binderVar.mult == Erase
// 		ty = pi.betaReduce(v)
// 		if !justType {
// 			tr = tr.(Lambda).betaReduce(v)
// 		}

// 		if ok := SetKind(&tr, ty); !ok {
// 			panic("bug: failed to set kind")
// 		}

// 		pi, ok = ty.(Pi)
// 	}
// 	return tr, ty
// }

func (env *Environment) appToPi(f Lambda, F Pi, a Term, Ap Type) (fa Term, Fa Type, ok bool) {
	F2, ok2 := env.Red(F).(Pi)
	if ok = ok2; !ok {
		env.error(ExpectedProductType, F)
		return
	}
	F = F2
	A := F.binderVar.Kind
	if ok = env.Unify(A, Ap); !ok {
		start, end := calcStartEnd(A, Ap)
		env.equivalenceError(A, Ap, start, end)
		return // can't apply term 'a' because it is not a term of type A
	}
	Ap_u := env.getUnified(Ap)
	if !SetKind(&a, Ap_u) {
		panic("bug: failed to set kind")
	}
	// TODO: need a Unify somewhere?
	fa = f.betaReduce(a)
	Fa = F.betaReduce(a)
	if !SetKind(&fa, Fa) {
		panic("bug: failed to set kind")
	}
	return fa, Fa, ok
}

// generalizes a type to a product type and attempts unification
//
//   - H is the type hypothesized to be a product type
//   - A is the type applied to the product type
func (env *Environment) Generalize(H Type, A Type) (Hp Pi, ok bool) {
	if !env.allowGeneralization { // this is primarily for testing purposes
		if Hp, ok = H.(Pi); !ok {
			env.error(ExpectedProductType, H)
			return Hp, false
		}
		return Hp, true
	}

	Hr := env.Red(H)
	Hp, ok = Hr.(Pi)
	if ok {
		return
	}

	Hp = env.PiHole(A)
	if ok = env.Unify(Hr, Hp); !ok { // TODO: should it unify rF = pF, or F = pF???
		// gen. Hr can't be unified with a prod type--this happens when Hr isn't a var
		return
	}
	debug_log_Gen(H, A, Hp)
	return
}

func (env *Environment) generateApplyFunction() Lambda {
	f := DummyVar("f")
	x := DummyVar("x")
	a := Hole("a")
	b := Hole("b")

	kA := GetKind(&a)
	_a := AsTyping(a)
	kB := GetKind(&b)
	var s, t Sort
	s, _ = env.Rule(kA, kB)
	t, _ = env.Rule(s, s)

	a_to_b := Pi{
		implicit:  false,
		binderVar: _a,
		dependent: b,
		kind:      s,
	}

	forall := Forall{
		variables: []Variable{a, b},
		body: Pi{
			implicit: false,
			binderVar: AsTyping(a_to_b),
			dependent: a_to_b,
			kind: t,
		},
	}

	// (\f, x => f x) : forall A, B . (A -> B) -> A -> B
	return Lambda{
		binder: f,
		bound: Lambda{
			binder: x,
			bound:  MakeApplication(b, f, x),
			Type:   a_to_b,
		},
		Type: forall,
	}
}

// Gets (or generates) an apply function `$`. If `$` is not set for quick access, it is generated
// and set for future quick access.
//
// NOTE: if `$` is set, it cannot be set after. So, if a specific `$` is needed, it must be set
// prior to calling this function.
func (env *Environment) getApplyFunction() Lambda {
	if env.applyFunction == nil {
		env.applyFunction = new(Lambda)
		*env.applyFunction = env.generateApplyFunction()
	}
	return *env.applyFunction
}

// findAbstraction attempts to find an associated lambda abstraction.
//
//   - if x is a lambda, it is returned
//   - if x is a variable, it is looked up in replacements
//   - if x is found and is a lambda, it is returned
//   - if x is found and is not a lambda, an error is returned
//   - if x is not found, it is treated as a free variable and applied to ...
//   - `$ x` := `(\a, b => a b) x` := `(\b => x b)`
func (env *Environment) findFunction(x Term) (Lambda, bool) {
	lam, ok := x.(Lambda)
	if ok {
		return lam, ok
	}

	var v Variable
	if v, ok = x.(Variable); !ok {
		env.error(NonFunctionError, x)
		return Lambda{}, false
	}

	var replacement Replacement
	replacement, ok = env.replacements.Find(v)
	if !ok {
		ok = true
		lam = env.getApplyFunction()
		term := lam.betaReduce(v)
		lam = term.(Lambda) // always true b/c of how apply function is generated
	} else {
		lam, ok = replacement.Term.(Lambda)
		if !ok {
			env.error(NonFunctionError, x)
			return lam, false
		}
	}

	return lam, true
}

type ForallIntro = func(ty Type) (forall Forall, ok bool)

// Forall rule (Forall): captures all free occurrences of each x in xs
//
//		  𝚪,x1,..,xN ⊢ t:T
//	------------------------ (Forall)
//	𝚪 ⊢ forall x1,..,xN . t : T
func (env *Environment) Forall(xs []Variable) ForallIntro {
	return func(ty Type) (forall Forall, ok bool) {
		// create a new forall type
		forall = Forall{variables: xs, body: ty}
		debug_log_Forall(xs, ty, forall)
		ok = true
		return
	}
}

func (env *Environment) Apply(f Term, F Type, a Term) (_ Term, _ Type, ok bool) {
	fun, ok := env.findFunction(f)
	if !ok {
		return nil, nil, false
	}

	return env.App(fun, F, a)
}

func (env *Environment) ImplicitApp(f Lambda, F Pi, a Term, Ap Type) (Term, Type, bool) {
	if !F.implicit {
		env.error(ExpectedImplicitTermInProduct, F.binderVar)
		return nil, nil, false
	}

	term, ty, good := env.appToPi(f, F, a, Ap)
	if !good {
		return nil, nil, false
	}
	debug_log_IApp(f, F, a, Ap, term, ty)
	return term, ty, true
}

// Application rule (App)
//
//	𝚪 ⊢ f:->>(x:A->B)   𝚪 ⊢ a:A'   A =β A'
//	-------------------------------------- (App)
//	        𝚪 ⊢ (f a): B[x:=a]
func (env *Environment) App(f Lambda, F Type, a Term) (Term, Type, bool) {
	Ap := GetKind(&a)
	// Hypothesized product
	var H Pi
	if _, isFn := F.(functionType); !isFn {
		var ok bool
		H, ok = env.Generalize(F, Ap)
		if !ok {
			return nil, nil, false
		}
	} else {
		F = env.specialize(F)
		var ok bool
		if H, ok = F.(Pi); !ok {
			env.error(ExpectedProductType, F)
			return nil, nil, false
		}
	}

	term, ty, good := env.appToPi(f, H, a, Ap)
	if !good {
		return nil, nil, false
	}
	debug_log_App(f, F, a, Ap, H, term, ty)
	return term, ty, true
}

func (env *Environment) makeDischarger(x Variable) (discharge func()) {
	// check if variable is shadowed
	if replacement, reset := env.replacements.Find(x); reset {
		// restore to shadowed variable instead of deleting
		return func() { env.replacements.Map(x, replacement) }
	}
	// delete variable, it doesn't shadow anything
	return func() { env.replacements.Delete(x) }
}

func (env *Environment) assume(x Variable, A Type) (discharge func()) {
	discharge = env.makeDischarger(x)
	env.replacements.Map(x, Replacement{x, A})
	return
}

// Abstraction rule (IAbs)
//
//	𝚪,x:A ⊢ b:B    𝚪 ⊢ ({x:A}->B):t
//	------------------------------ (IAbs)
//	    𝚪 ⊢ (\{x} => b): ({x:A}->B)
//
// requiring the product type on the first call ensures that it was derived without the assumption
// [x:A] created within
func (env *Environment) ImplicitAbs(x Variable, A Type, P Pi) (derive func(b Term, B Type) (Lambda, Pi, bool)) {
	A = env.getUnified(A)
	discharge := env.assume(x, A) // create assumption

	derive = func(b Term, B Type) (Lambda, Pi, bool) {
		var lam Lambda
		defer discharge() // discharge assumption

		// now, given P, derive lambda and associated type product type
		if !env.Unifiable(P.binderVar.Kind, A) {
			env.error(expectedTypeOfTermToMatchProdTerm(P.binderVar.Kind, A), A)
			return lam, P, false
		}
		if !env.Unifiable(P.dependent, B) {
			env.error(expectedTypeOfTermToMatchProdTerm(P.dependent, B), B)
			return lam, P, false
		}

		start, end := calcStartEnd(x, b)
		lam = Lambda{binder: x, bound: b, Type: P, implicit: true, Start: start, End: end}
		debug_log_IAbs(x, A, b, B, P, lam)
		return lam, P, true
	}
	return
}

func (env *Environment) getUnifiedTerm(a Term) Term {
	v, isVar := a.(Variable)
	if !isVar || !v.isHole {
		return a
	}
	tmp, found := env.unifications.Find(v)
	if !found {
		return a
	}
	return tmp
}

func (env *Environment) getUnified(A Type) Type {
	a, isVar := A.(Variable)
	if !isVar || !a.isHole {
		return A
	}
	tmp, found := env.unifications.Find(a)
	if !found {
		return A
	}
	if tmp2, ok := tmp.(Type); ok {
		return tmp2
	}

	return A
}

type providePiFunc = func(Pi) (Lambda, Pi, bool)

type AbsSecondFunc = func(b Term, B Type) providePiFunc

func (env *Environment) unifyingSubstitution(x Variable, b Term) (Variable, Term, bool) {
	_ = GetKind(&x) // make sure kind is set
	_ = GetKind(&b) // make sure kind is set
	vs := b.CollectVariables(make(map[string]Variable))
	if u, found := vs[x.x]; found {
		if !env.Unify(x.Kind, u.Kind) {
			s, e := calcStartEnd(x.Kind, u.Kind)
			env.unifyingError(x.Kind, u.Kind, s, e)
			return x, b, false
		}
		x.Kind = env.getUnified(x.Kind)
	}
	return x, b, true
}

// Abstraction rule (Abs)
//
//	𝚪,x:A ⊢ b:B    𝚪 ⊢ (x:A->B):t
//	------------------------------ (Abs)
//	    𝚪 ⊢ (\x => b): (x:A->B)
func (env *Environment) Abs(x Variable, A Type) (derive AbsSecondFunc) {
	//fmt.Fprintf(os.Stderr, "%v : %v\n", x, A)
	A = env.getUnified(A)
	discharge := env.assume(x, A) // create assumption

	derive = func(b Term, B Type) providePiFunc {
		//fmt.Fprintf(os.Stderr, "%v : %v\n", b, B)
		B = env.getUnified(B)
		m := b.CollectVariables(make(map[string]Variable))
		if v, found := m[x.x]; found {
			if !env.Unify(v.Kind, A) {
				env.unifyingError(v.Kind, A, v.Start, v.End)
				return nil
			}
			x.Kind = env.getUnified(A)
			A = x.Kind
		}
		defer discharge() // discharge assumption now that term is derived
		// return a function that derives the lambda term
		return func(P Pi) (Lambda, Pi, bool) {
			//fmt.Fprintf(os.Stderr, "%v\n", P)
			// now, given P, derive lambda and associated type product type

			// get type (and avoid nil pointer error)
			P.binderVar.Kind = env.getUnified(GetKind(&P.binderVar))
			P.dependent = env.getUnified(P.dependent)

			if !env.Unifiable(P.binderVar.Kind, A) {
				env.error(expectedTypeOfTermToMatchProdTerm(P.binderVar.Kind, A), A)
				return Lambda{}, Pi{}, false
			}
			if !env.Unifiable(P.dependent, B) {
				env.error(expectedTypeOfTermToMatchProdTerm(P.dependent, B), B)
				return Lambda{}, Pi{}, false
			}

			start, end := calcStartEnd(x, b)
			var ok bool
			x, b, ok = env.unifyingSubstitution(x, b)
			if !ok {
				return Lambda{}, P, false
			}
			lam := Lambda{binder: x, bound: b, Type: P, Start: start, End: end}
			debug_log_Abs(x, A, b, B, P, lam)
			return lam, P, true
		}
	}
	return derive
}

func (env *Environment) isImplicitlyBound(A Term) (a Variable, implicit bool) {
	if a, implicit = A.(Variable); !implicit {
		return
	}
	_, found := env.replacements.Find(a)
	implicit = !found // not found in context, so it must not be bound; i.e., it's free
	return
}

func (env *Environment) implicitToExplicitTypeBinding(a Variable) func(bound Type, t Sort) Pi {
	s := env.system.implicitKind
	a.mult = Erase
	discharge := env.assume(a, s)
	return func(bound Type, t Sort) Pi {
		defer discharge()
		u := env.system.rules[s][t]
		return Pi{
			implicit:  true,
			binderVar: a,
			dependent: bound,
			kind:      u,
		}
	}
}

// given a type `A`, returns `erase _ : A`
func AsTyping(A Type) Variable {
	start, end := A.Pos()
	v := Wildcard()
	v.demangler = nextEnvUid()
	v.isHole = false
	v.Kind = A
	v.Start, v.End = start, end
	return v
}

type PiIntro = func(B Type) (pi Pi, ok bool)

func (env *Environment) Product(term Term) (intro PiIntro, ok bool) {
	var x Variable
	if x, ok = term.(Variable); !ok {
		env.error(ExpectedVariable, term)
		return
	}
	return env.Prod(x)
}

// Product rule (Prod):
//
//	𝚪 ⊢ A:->>s   𝚪,x:A ⊢ B:->>t   s~>t:u
//	------------------------------------ (Prod)
//	         𝚪 ⊢ {x:A}->B:u
func (env *Environment) ImplicitProd(x Variable) (intro PiIntro, ok bool) {
	var A Type
	var s Sort
	if x, A, s, ok = reducedTypingTriple(env, x); !ok {
		return nil, false
	}

	doIntro := env.generatePiIntro(x, A, s, true)

	intro = func(B Type) (pi Pi, ok bool) {
		if pi, ok = doIntro(B); !ok {
			return pi, false
		}
		debug_log_IProd(A, s, x, B, pi.kind, pi)
		return pi, true
	}

	return intro, true
}

func (env *Environment) Prod_NoGeneralization(x Variable) (intro PiIntro, ok bool) {
	var A Type
	var s Sort
	if x, A, s, ok = reducedTypingTriple(env, x); !ok {
		return nil, false
	}

	doIntro := env.generatePiIntro(x, A, s, false)

	intro = func(B Type) (pi Pi, ok bool) {
		if pi, ok = doIntro(B); !ok {
			return pi, false
		}
		debug_log_Prod(A, s, x, B, pi.kind, pi)
		return pi, true
	}

	return intro, true
}

// assumes `x : A` and returns a function that generates a dependent product type
func (env *Environment) generatePiIntro(x Variable, A Type, s Sort, implicit bool) PiIntro {
	discharge := env.assume(x, A)
	return func(B Type) (pi Pi, ok bool) {
		// discharges `x : A`
		defer discharge()
		// determine kind of dependent product type (and if valid)
		var t Sort
		if B, t, ok = reducedTypingDouble[Type, Sort](env, B); !ok {
			return
		}

		// see if there's a rule `s ~> t : u`
		var u Sort
		if u, ok = env.Rule(s, t); !ok {
			return pi, false
		}
		start, end := calcStartEnd(x, A, B, t)
		pi = Pi{
			implicit:  implicit,
			binderVar: x,
			dependent: B,
			kind:      u,
			Start:     start,
			End:       end,
		}
		return pi, true
	}
}

func reducedTypingTriple[T Term](env *Environment, arg T) (a T, A Type, s Sort, ok bool) {
	a = arg
	A = GetKind(&a)
	A = env.Red(A)
	kA := GetKind(&A)
	if s, ok = kA.(Sort); !ok {
		env.error(unexpectedType(a, A, kA), a)
		return a, A, s, false
	}
	return a, A, s, true
}

func reducedTypingDouble[T Term, S Type](env *Environment, arg T) (A T, s S, ok bool) {
	A = arg
	kA := GetKind(&A)
	kA = env.Red(kA)
	if s, ok = kA.(S); !ok {
		env.error(unexpectedType(A, kA), A)
		return A, s, false
	}
	return A, s, true
}

// Product rule (Prod):
//
//	𝚪 ⊢ A:->>s   𝚪,x:A ⊢ B:->>t   s~>t:u
//	------------------------------------ (Prod)
//	         𝚪 ⊢ (x:A)->B:u
func (env *Environment) Prod(x Variable) (intro PiIntro, ok bool) {
	var A Type
	var s Sort
	if x, A, s, ok = reducedTypingTriple(env, x); !ok {
		return nil, false
	}

	doBind := common.Const[*Pi]

	if a, implicit := env.isImplicitlyBound(A); implicit && env.allowGeneralization {
		bind := env.implicitToExplicitTypeBinding(a)
		doBind = func(x *Pi) func(any) *Pi {
			return func(any) *Pi { *x = bind(*x, x.kind); return x }
		}
	}

	doIntro := env.generatePiIntro(x, A, s, false)

	intro = func(B Type) (pi Pi, ok bool) {
		// wraps result in implicit binding and discharges any assumptions
		defer doBind(&pi)(nil)
		if pi, ok = doIntro(B); !ok {
			return pi, false
		}

		debug_log_Prod(A, s, x, B, pi.kind, pi)
		return pi, true
	}

	return intro, true
}

func matchHead(head Constant, actual Term) bool {
	return GetConstant(actual).C == head.C
}

func (env *Environment) con_end_clause(expectedTailHead Constant, c Constant, t Type, vs []Term) (term Term, ok bool) {
	if ok = matchHead(expectedTailHead, t); !ok {
		env.error(illegalConstructor(expectedTailHead, t), t)
		return nil, false
	}
	term = MakeCApplication(t, c, vs...)
	return term, ok
}

func (env *Environment) con_helper_helper(expectedTailHead Constant, c Constant, pi Pi, vs []Term) (Lambda, bool) {
	v := DummyVar(fmt.Sprintf("x%d", len(vs)))
	v.mult = pi.binderVar.mult
	vs = append(vs, v)
	lam := Lambda{
		binder:   v,
		Type:     pi,
		implicit: pi.implicit,
	}
	var ok bool
	if piInner, isPi := pi.dependent.(Pi); isPi {
		lam.bound, ok = env.con_helper2(expectedTailHead, c, piInner, vs)
	} else {
		lam.bound, ok = env.con_end_clause(expectedTailHead, c, pi.dependent, vs)
	}
	return lam, ok
}

func (env *Environment) con_helper2(expectedTailHead Constant, c Constant, t Type, ts []Term) (Term, bool) {
	if pi, ok := t.(Pi); ok {
		if pi.implicit && pi.binderVar.mult == Erase {
			return env.con_helper2(expectedTailHead, c, pi.dependent, ts)
		}
		lam, ok := env.con_helper_helper(expectedTailHead, c, pi, ts)
		return lam, ok
	}
	return env.con_end_clause(expectedTailHead, c, t, ts)
}

func (env *Environment) con_helper(expectedTailHead Constant, c Constant, t Type) (Term, bool) {
	return env.con_helper2(expectedTailHead, c, t, []Term{})
}

// Constructor rule (Con):
//
//	free(`C`)    𝚪, Z .. where .. ⊢ (x:A)->(y:B)-> .. -> Z : u
//	---------------------------------------------------------- (Con)
//	 𝚪, Z .. where .. ⊢ C := (\x, y, .., z => `C` x y .. z) : (x:A)->B-> .. -> Z
func (env *Environment) Con(TypeConsKind Universe, Z Constant, C Constant, t Type) (constructor Term, ok bool) {
	// attempt to unify type constructor's kind and the kind of the data constructor
	//
	// the only time these would fail to match are ...
	//		- type constructor is not fully applied in the data constructor's final term
	//		- the final term is not the applied type constructor
	// in the final case this will be caught by the `con_end_clause` function even if
	// it has the same kind as the type constructor
	tmp, w := t.GetKind()
	t = tmp.(Type)
	if ok = env.Unify(w, TypeConsKind); !ok {
		return
	}

	constructor, ok = env.con_helper(Z, C, t)
	if !ok {
		return
	}
	ok = env.Declare(C, t)
	ok = ok && env.Assign(C, constructor)
	ok = ok && env.mapConstructor(Z, C)
	if !ok {
		return
	}
	debug_log_Con(Z, TypeConsKind, C, t, constructor)
	return constructor, ok
}

type TypeConIntro = func(Cs []Constant, Ts []Type) (ok bool)

func (env *Environment) generateTypeConIntro(Z Constant, TypeConsConstructs Universe) TypeConIntro {
	return func(Cs []Constant, Ts []Type) (ok bool) {
		if len(Cs) != len(Ts) {
			panic("bug: mismatched number of constructors and types")
		}

		for i := range Cs {
			if _, ok = env.Con(TypeConsConstructs, Z, Cs[i], Ts[i]); !ok {
				return
			}
		}
		return
	}
}

// Type Constructor rule (TypeCon):
//
// Z : A
//
//   - U is the last term in the type constructor
//   - A is the type constructor
//   - Z is the name of the type constructor
func (env *Environment) TypeCon(U Type, Z Constant, A Type) (intro TypeConIntro, ok bool) {
	// check that U is a universe
	if _, ok = U.(Universe); !ok {
		env.error(ExpectedTypeUniverse, U)
		return
	}

	// make sure A has a kind
	tmp, kA := A.GetKind()
	A = tmp.(Type)

	// check that type kind of the type constructor is ...
	//		- a universe
	//		- and that it is either Type or Type 1
	// keep in mind that A.GetKind() returns something that is not
	// explicitly set by the user, but is inferred by the system via a rule
	if s, ok := kA.(Sort); !ok || !s.Known() {
		env.error(IllegalConstant, kA)
		return nil, false
	}

	var TypeCon_Type Universe
	if TypeCon_Type, ok = kA.(Universe); !ok {
		env.error(ExpectedTypeUniverse, kA)
		return nil, false
	}

	// Type_n is explicitly set by the user, it's the last term in the type constructor
	Type_n := U.(Universe)

	if Type_n != Universe(0) {
		env.error(typeNConsNotLegal(Type_n), A)
		return nil, false
	}

	s, e := U.Pos()
	// create constant version of type universe
	u := Constant{Type_n.String(), s, e}

	// Z constructs a type of kind u
	_, ok = env.Con(TypeCon_Type, u, Z, A)
	if !ok {
		return nil, false
	}
	intro = env.generateTypeConIntro(Z, Type_n)
	return intro, ok
}

// x : A ∈ 𝚪
// ---------- (Var)
// 𝚪 ⊢ x : A
func (env *Environment) Var(x Variable, A Type) (_ Variable, _ Type, ok bool) {
	// check if variable exists in environment
	_, ok = env.replacements.Find(x)
	if !ok {
		start, end := x.Pos()
		env.unknownNameError(x, start, end)
		return
	}
	debug_log_Var(x, A)
	return x, A, ok
}

func (env *Environment) Red(A Type) (rA Type) {
	if c, ok := A.(Constant); ok {
		if r, ok := env.replacements.Find(c); ok {
			A = r.Term.(Type)
		}
	}
	if v, ok := A.(Variable); ok {
		if r, ok := env.replacements.Find(v); ok {
			A = r.Term.(Type)
		}
	}
	rA = env.reduceToWHNF(A).(Type)
	debug_log_Red(A, rA)
	return rA
}

func (env *Environment) declare(id strPos, ty Type) (ok bool) {
	// is re-declared?
	if _, found := env.replacements.Find(id); found {
		env.error(AlreadyDeclared, id)
		return false
	}

	// initially has no term
	env.replacements.Map(id, Replacement{__, ty})
	return true
}

// Assign a replacement term to a name. The name must already be declared
func (env *Environment) Assign(name strPos, term Term) (ok bool) {
	var r Replacement
	// check if variable declared
	if r, ok = env.replacements.Find(name); !ok {
		start, end := name.Pos()
		env.unknownNameError(name, start, end)
		return
	}

	if ok = SetKind(&term, r.Type); !ok {
		return
	}
	r = Replacement{term, r.Type}
	// now has both term and type in environment
	env.replacements.Map(name, r)
	return true
}

func (env *Environment) Declare(name strPos, T Type) (ok bool) {
	return env.declare(name, T)
}
