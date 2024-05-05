// =================================================================================================
// Alex Peters - February 29, 2024
//
// =================================================================================================
package types

import "github.com/petersalex27/yew/common"

// specifically equals relation on a variable and a term
type Equality struct {
	Variable
	Term
}

func (env *Environment) SpecializeImplicit(ty Type) (t Type) {
	for pi, ok := ty.(Pi); ok && pi.implicit; pi, ok = t.(Pi) {
		v := Var("?")
		t = pi.betaReduce(v)
	}
	return t
}

// apply an implicit argument explicitly
func (env *Environment) AppToPi(f Lambda, F Pi, a Term, Ap Type) (fa Term, Fa Type, ok bool) {
	A := F.binderType
	if ok = betaEquivalence(A, Ap); !ok {
		return // can't apply term 'a' because it is not a term of type A
	}
	fa = f.betaReduce(a)
	Fa = F.betaReduce(a)
	return fa, Fa, ok
}

// generalizes a type to a product type and attempts unification
//
// 	- H is the type hypothesized to be a product type
// 	- A is the type applied to the product type
func (env *Environment) Generalize(H Type, A Type) (Hp Pi, ok bool) {
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
	return
}

// Application rule (App)
//
//	𝚪 ⊢ f:->>(x:A->B)   𝚪 ⊢ a:A'   A =β A'
//	-------------------------------------- (App)
//	        𝚪 ⊢ (f a): B[x:=a]
func (env *Environment) App(f Lambda, F Type, a Term, Ap Type) (_ Term, _ Type, ok bool) {
	// Hypothesized product
	var H Pi
 	if H, ok = env.Generalize(F, Ap); !ok {
		return
	}

	return env.AppToPi(f, H, a, Ap)
}

func (env *Environment) makeDischarger(x Variable, A Type) (discharge func()) {
	// check if variable is shadowed
	if term, reset := env.typings.Find(x); reset { 
		// restore to shadowed variable instead of deleting
		return func() { env.typings.Map(x, term) }
	}
	// delete variable, it doesn't shadow anything
	return func() { env.typings.Delete(x) }
}

func (env *Environment) assume(x Variable, A Type) (discharge func()) {
	discharge = env.makeDischarger(x, A)
	env.typings.Map(x, A)
	return
}

// Abstraction rule (Abs)
//
//	𝚪,x:A ⊢ b:B    𝚪 ⊢ (x:A->B):t
//	------------------------------ (Abs)
//	    𝚪 ⊢ (\x => b): (x:A->B)
//
// requiring the product type on the first call ensures that it was derived without the assumption
// [x:A] created within
func (env *Environment) Abs(x Variable, A Type, P Pi) (derive func(b Term, B Type) (_ Lambda, _ Pi, ok bool)) {
	discharge := env.assume(x, A) // create assumption

	derive = func(b Term, B Type) (_ Lambda, _ Pi, ok bool) {
		defer discharge() // discharge assumption

		// now, given P, derive lambda and associated type product type
		binderTypeMatch := Equals(P.binderType, A)
		resTypeMatch := Equals(P.dependent, B)
		ok = binderTypeMatch && resTypeMatch
		if !ok {
			return
		}

		typing := VarTyping{x, A}
		lam := Lambda{binder: typing, bound: b}
		return lam, P, true
	}
	return
}


type piIntro = func(B Type, KindOfB Type) (pi Pi, s Sort, ok bool)

func (env *Environment) isImplicitlyBound(A Term) (a Variable, implicit bool) {
	if a, implicit = A.(Variable); !implicit {
		return
	}
	_, found := env.typings.Find(a)
	implicit = !found // not found in context, so it must not be bound; i.e., it's free
	return 
}

func trivialBinding(bound Type, t Sort) Pi {
	if pi, ok := bound.(Pi); ok {
		return pi
	}
	// trivial implicit binding--binds nothing
	return Pi{
		implicit: true,
		binderVar: __,
		binderType: __,
		dependent: bound,
		kind: t,
	}
}

// func (env *Environment) guessSortHelper(A Type) Sort {
// 	if p, ok := A.(Pi); ok {
// 		return p.kind
// 	}

// 	if id, ok := A.(Identifier); ok {
// 		term, found := env.typings.Find(id)
// 		if !found {
// 			return EVar("?")
// 		}
// 		if p, ok := term.(Pi); ok {
// 			return p.k
// 		}
// 	}
// }

// func (env *Environment) guessSort(A Type) Sort {
// 	// A will either be a term of * or *1 since it's a type
// 	if p, ok := A.(Pi); ok {
// 		return p.kind
// 	}
	
// 	c, ok := A.(Constant)
// 	if !ok {
// 		return env.system.implicitKind
// 	}
// 	var s Sort
// 	s, ok = env.system.axioms[c]
// 	if !ok {
// 		return env.system.implicitKind
// 	}
// 	return s
// }

// noop when a:A is already bound
// func (env *Environment) bindImplicitImplicitly(a Variable, A Type) func(bound Type, t Sort) Pi {

// 	Av, implicit := env.isImplicitlyBound(A)
// 	if !implicit {
// 		return trivialBinding
// 	}

// 	sA := env.guessSort(A)

// 	if a != __ {
// 		// important! do not allow inferred implicitly bound name to be available in function
// 		//	- only explicit implicitly bound names should be available
// 		s := env.system.implicitKind
// 		discharge := env.assume(Av, s)
// 		return func(bound Type, t Sort) Pi {
// 			defer discharge()
// 			u := env.system.rules[s][t]
// 			a.mult = Erase // make multiplicity 0--this is an implicit implicit binding of `a`
// 			return Pi{ // s -> (s -> t) : 
// 				implicit: true,
// 				binderVar: Av,
// 				binderType: env.system.implicitKind,
// 				dependent: Pi{ // s -> t : u
// 					implicit: true,
// 					binderVar: a,
// 					binderType: Av,
// 					dependent: bound,
// 					kind: s,
// 				},
// 				kind: u,
// 			}
// 		}
// 	}

// 	s := env.system.implicitKind
// 	discharge := env.assume(Av, s)
// 	return func(bound Type, t Sort) Pi {
// 		defer discharge()
// 		u := env.system.rules[s][t]
// 		return Pi{
// 			implicit: true,
// 			binderVar: a,
// 			binderType: env.system.implicitKind,
// 			dependent: bound,
// 			kind: u,
// 		}
// 	}
// 	discharge := env.assume(a, A)
// }

func (env *Environment) implicitToExplicitTypeBinding(a Variable) func(bound Type, t Sort) Pi {
	s := env.system.implicitKind
	discharge := env.assume(a, s)
	return func(bound Type, t Sort) Pi {
		defer discharge()
		u := env.system.rules[s][t]
		return Pi{
			implicit: true,
			binderVar: a,
			binderType: env.system.implicitKind,
			dependent: bound,
			kind: u,
		}
	}
}

// just like product rule, but with "_" as the term of "A"
func (env *Environment) Prod2(A Type, KindOfA Type) (intro piIntro, ok bool) {
	return env.Prod(__, A, KindOfA)
}

// Product rule (Prod):
//
//	𝚪 ⊢ A:->>s   𝚪,x:A ⊢ B:->>t   s~>t:u
//	------------------------------------ (Prod)
//	         𝚪 ⊢ (x:A)->B:u
func (env *Environment) Prod(x Variable, A Type, KindOfA Type) (intro piIntro, ok bool) {
	ka := env.Red(KindOfA)
	var s, t Sort
	if s, ok = ka.(Sort); !ok {
		return
	}

	doBind := common.Const[*Pi]

	if a, implicit := env.isImplicitlyBound(A); implicit {
		bind := env.implicitToExplicitTypeBinding(a)
		doBind = func(x *Pi) func(any) *Pi {
			return func(any) *Pi { *x = bind(*x, x.kind); return x }
		}

		// if _, implicit := env.isImplicitlyBound(x); implicit {
		// 	bind := env.implicitToExplicitTypeBinding(x)
		// 	doTypeBind := doBind
		// 	doBind = func(x *Pi) func(any) *Pi {
		// 		return func(any) *Pi {
		// 			*x = bin
		// 		}
		// 	}
		// }
	}

	discharge := env.assume(x, A)

	intro = func(B, KindOfB Type) (pi Pi, u Sort, ok bool) {
		// wraps result in implicit binding and discharges any assumptions
		defer doBind(&pi)(nil) 
		// discharges `x : A`
		defer discharge()
		// determine kind of dependent product type (and if valid)
		kb := env.Red(KindOfB)
		if t, ok = kb.(Sort); !ok {
			return
		}
		// see if there's a rule `s ~> t : u`
		if u, ok = env.system.Rule(s, t); !ok {
			return
		}
		pi = Pi{
			binderVar:  x,
			binderType: A,
			dependent:  B,
			kind:       u,
		}
		return
	}

	return
}

func (env *Environment) ProdImplicit(x Variable, A Type, KindOfA Type) (intro piIntro, ok bool) {
	f0, ok0 := env.Prod(x, A, KindOfA)
	if !ok0 {
		return
	}

	ok = ok0

	intro = func(B, KindOfB Type) (pi Pi, s Sort, ok bool) {
		pi, s, ok = f0(B, KindOfB)
		pi.implicit = true
		return
	}

	return intro, ok
}

//
//	C ∉ 𝚪    𝚪 ⊢ (x:A)->B:u
//	----------------------- (Con)
//	   𝚪 ⊢ C : (x:A)->B
// func (env *Environment) Con(c Constant, pi Pi) (ok bool) {
// 	_, found := env.typings.Find(c)
// 	if ok = !found; !ok {
// 		return // error: redeclared
// 	}

// 	env.declare()
// }

func (env *Environment) Var(x Variable, A Type) (_ Variable, _ Type, ok bool) {
	var t Term
	var B Type
	// check if variable exists in environment
	t, ok = env.typings.Find(x)
	if !ok {
		env.unknownNameError(x, 0, 0) // TODO: not 0, 0
		return
	}
	// check if variable's type matches type `A`
	B, ok = t.(Type)
	if !ok {
		env.orderPrecedesType0Error(t, 0, 0)
		return
	}
	B = env.SpecializeImplicit(B)
	return x, A, env.Unify(A, B) // try to unify
}

func (env *Environment) Red(t Type) Type {
	return reduceToWHNF(t)
}

func (env *Environment) declare(id Identifier, ty Type) (ok bool) {
	// is re-declared?
	if _, found := env.typings.Find(id); found {
		ok = false
		return
	}

	env.typings.Map(id, ty)
	return
}

func (env *Environment) Declare(typing Typing) (ok bool) {
	var id Identifier
	id, ok = typing.Term.(Identifier)
	if !ok {
		return
	} 
	return env.declare(id, typing.Kind)
}

func (env *Environment) declarations(ts []Typing) (ok bool) {
	for _, typing := range ts {
		if ok = env.Declare(typing); !ok {
			return
		}
	}
	return
}

// func (env *Environment) define(eq Equality) (ok bool) {
// 	// make sure variable declared
// 	if _, ok = env.Find(eq.Variable); !ok {
// 		return
// 	}

// 	// is re-defined?
// 	if _, found := env.defs.Find(eq.Variable); found {
// 		// TODO: allow for pattern matching defs
// 		// TODO: do NOT allow for re-definitions
// 		//	- a variable is redefined if it is already defined AND the previous thing defined was not of the same name
// 		ok = false
// 		return
// 	}

// 	// declare for the first time in environment
// 	env.Map(eq.Variable, eq.Term)
// 	return
// }

// func (env *Environment) definitions(eqs []Equality) (ok bool) {
// 	for _, eq := range eqs {
// 		if ok = env.define(eq); !ok {
// 			return
// 		}
// 	}
// 	return
// }

// func (env *Environment) Let()

// creates a where-clause that allows for mutually recursive definitions and declarations
//
//	 𝚪,x1:A1,..,xN:AN ⊢ x1=e1:A1,..,xN=eN:AN    𝚪 ⊢ m:T
//	---------------------------------------------------- (Where)
//	𝚪 ⊢ m:T where {x1:A1,..,xN:AN, x1=e1:A1,..,xN=eN:AN}
// func (env *Environment) Where(ts ...Typing) (
// 	assign func(defs ...Equality) (
// 		bind func(Term, Type) (TermBind, bool),
// 		ok bool,
// 	),
// 	whereEnv *EmbeddedEnvironment,
// ) {

// 	whereEnv = NewEmbeddedEnvironment(env)
// 	ok := whereEnv.declarations(ts)
// 	if !ok {
// 		whereEnv = nil
// 		return
// 	}

// 	// create assignments
// 	assign = func(defs ...Equality) (bind func(Term, Type) (TermBind, bool), ok bool) {
// 		ok = whereEnv.definitions(defs)
// 		if !ok {
// 			return
// 		}

// 		// bind environment to term and type m and t resp.
// 		bind = func(m Term, t Type) (TermBind, bool) {
// 			term := TermBind{
// 				bound:     m,
// 				boundType: t,
// 				binding:   whereEnv,
// 			}
// 			return term, true
// 		}
// 		return
// 	}
// 	return
// }
