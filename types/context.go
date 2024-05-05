// =================================================================================================
// Alex Peters - February 27, 2024
// =================================================================================================
package types

import (
	"sync"
	
	"github.com/petersalex27/yew/common/table"
	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/source"
)

// Type system
type System struct {
	// set of primitives
	constants map[Sort]Sort
	// set of axioms
	axioms map[Sort]Sort
	// the default kind of an implicit variable
	implicitKind Sort
	// rules using the constants for determining what sort something is
	rules map[Sort]map[Sort]Sort
}

const (
	Monotype  Constant = "*"
	Polytype  Constant = "**"
	Monotype2 Constant = "*1"
	Polytype2 Constant = "**1"
)

var std = System{
	constants: map[Sort]Sort{
		Type0: Type0,
		Type1: Type1,
	},
	axioms: map[Sort]Sort{
		Type0: Type1,
	},
	rules: map[Sort]map[Sort]Sort{
		Type0: {
			Type0: Type0, // * ~> * : *
			Type1: Type1, // * ~> *1 : *1
		},
		Type1: {
			Type0: Type0, // *1 ~> * : *
			Type1: Type1, // *1 ~> *1 : *1
		},
	},
	implicitKind: Type0,
}

var stdPoly = System{
	constants: map[Sort]Sort{
		Monotype:  Monotype,
		Polytype:  Polytype,
		Monotype2: Monotype2,
		Polytype2: Polytype2,
	},
	axioms: map[Sort]Sort{
		Monotype: Monotype2, // * : #
		Polytype: Polytype2, // ** : ##
	},
	rules: map[Sort]map[Sort]Sort{
		Monotype: {
			Monotype:  Monotype,
			Polytype:  Polytype,
			Monotype2: Monotype2,
			// Polytype2: Polytype2, ??? // TODO: ???
		},
		Polytype: {
			Monotype: Polytype,
			Polytype: Polytype,
		},
		Monotype2: {
			Monotype:  Polytype,
			Polytype:  Polytype,
			Monotype2: Monotype2,
		},
	},
	implicitKind: Monotype,
}

func (sys System) Axiom(c Sort) (s Sort, t Sort, ok bool) {
	s, ok = sys.constants[c]
	if !ok {
		return
	}
	t, ok = sys.axioms[s]
	return
	//sys.axioms[]
}

func (sys System) Rule(s, t Sort) (u Sort, ok bool) {
	var m map[Sort]Sort
	m, ok = sys.rules[s]
	if !ok {
		return
	}
	u, ok = m[t]
	return
}

type TermOf struct {
	Term
	Type
}

type Assumptions interface {
	// functions for `Table`s from "github.com/petersalex27/yew/common"
	Find(k Identifier) (Term, bool)
	Delete(k Identifier)
	Map(k Identifier, t Term)
	Parent() (parent Assumptions, exists bool)
}

type Identifier interface {
	Term
	_identifier_()
}

type definitions struct {
	// give access to unifications
	typings *table.Table[Identifier, Term]
	// TODO: need a way to get access to the explicitly implicitly bound
}

// Proof environment--symbolically, this is Γ
type Environment struct {
	src      source.SourceCode
	imported *table.Table[Constant, *Environment]
	// typings for names
	typings *table.Table[Identifier, Term]
	// unification assignments
	unifications *table.Table[Variable, Term]
	defs *table.Table[Variable, Term]
	//defs        *common.Table[Variable, Term]
	// type system
	system   System
	messages []errors.ErrorMessage
}

func (env *Environment) TypeOf(id Identifier) (ty Term, found bool) {
	return env.typings.Find(id)
}

var uidCounter uint = 0
var mu sync.Mutex

func _clearUidCounter() {
	mu.Lock()
	defer mu.Unlock()

	uidCounter = 0
}

func nextEnvUid() uint {
	mu.Lock()
	defer mu.Unlock()

	res := uidCounter
	uidCounter++
	return res
}

type EmbeddedEnvironment struct {
	uid               uint
	parentEnvironment Assumptions
	Environment
}

func NewEmbeddedEnvironment(parent Assumptions) *EmbeddedEnvironment {
	return &EmbeddedEnvironment{
		uid:               nextEnvUid(),
		parentEnvironment: parent,
		Environment:       *NewEnvironment(),
	}
}

// returns parent and true--second return value should never be false
func (env *EmbeddedEnvironment) Parent() (Assumptions, bool) {
	return env.parentEnvironment, true
}

func NewEnvironment() *Environment {
	return &Environment{
		imported:    table.MakeTable[Constant, *Environment](8),
		typings:     table.MakeTable[Identifier, Term](16),
		unifications: table.MakeTable[Variable, Term](16),
		system:      std,
		messages:    []errors.ErrorMessage{},
	}
}

// lol
//
// generates a Pi type:
//
//	(𝚷x:A.y):u
func (env *Environment) PiHole(A Type) Pi {
	uid := nextEnvUid()
	x := demVar("?x", uid)
	y := demVar("?y", uid)
	u := demVar("?Kind", uid)
	return Pi{binderVar: x, binderType: A, dependent: y, kind: u}
}

// returns _, false
func (env *Environment) Parent() (_ Assumptions, exists bool) {
	exists = false
	return
}

// check if v occurs in t. If it does, return true; else, return false.
// if t = v, then v is not in t, v is t
func (env *Environment) occurs(v Variable, t Term) (vOccursInT bool) {
	if _, ok := t.(Variable); ok {
		return false // v and t are identical--a variable does not occur in itself
	}

	return t.Locate(v)
}

// declares, for types variable v, monotype t:
//
//	v = t
//
// if v ∈ t, then union returns `OccursCheckFailed`; else, skipUnify is returned
func (env *Environment) union(v Variable, t Term) bool {
	if env.occurs(v, t) {
		return false
	}

	// map v :-> t
	env.unifications.Map(v, t)
	return true
}

// at this point it is known that both ta and tb belong to a unified kind
func (env *Environment) substitute(ta, tb Term) (ok bool, doUnify bool) {
	ok, doUnify = true, true
	if v, isVar := ta.(Variable); isVar {
		ok, doUnify = env.union(v, tb), false
	} else if v, isVar := tb.(Variable); isVar {
		ok, doUnify = env.union(v, ta), false
	}

	return
}

func (env *Environment) splitUnify(ta, tb Term) (ok bool) {
	// get constants, params, and indexes
	ca, termsA := Split(ta)
	cb, termsB := Split(tb)

	if ca != cb {
		// error: mismatch constants
		return false
	}

	if len(termsA) != len(termsB) {
		// error: patterns of different length, unification impossible
		return false
	}

	// it. through all params while stat is ok, unifying params
	for i := 0; ok && i < len(termsA); i++ {
		a, b := termsA[i], termsB[i]
		ok = env.Unify(a, b)
	}
	return
}

// creates a hole for the kind of `a`
//
//	a : ?type.#
func (env *Environment) HoleForKind(a Term) Typing {
	// TODO: keep track of what terms belong to the newly created hole?
	hole := Var("?type.")
	return Typing{Term: a, Kind: hole}
}

// tries to find typing for a reduced form of term `a`
//
// TODO: until termination/totality "checking" is implemented, this function has no
// defense against non-termination
func (env *Environment) FindUnified(a Term) (t Term) {
	ra := reduce(a)
	v, ok := ra.(Variable)
	if ok {
		t, ok = env.unifications.Find(v)
	}

	// for both failure to find and non-variable reduction result
	if !ok {
		t = ra
	}

	return
}

// returns kind of term or creates a hole
func (env *Environment) KindOf(id Identifier) (kind Term) {
	var found bool
	if kind, found = env.TypeOf(id); !found {
		kind = Var("?") // create a hole
	}
	return kind
}

func (env *Environment) equateKinds(ka, kb Type) bool {
	ok := betaEquivalence(ka, kb)
	if !ok {
		env.equivalenceError(ka, kb, 0, 0) // TODO: not 0, 0
	}
	return ok
}

// unifies two terms a, b
func (env *Environment) Unify(a, b Term) (ok bool) {
	ta := env.FindUnified(a)
	tb := env.FindUnified(b)

	var doSplitUnify bool
	ok, doSplitUnify = env.substitute(ta, tb)
	if ok && doSplitUnify {
		ok = env.splitUnify(ta, tb)
	}
	return
}

/*
threeIsThree : 3 = 3
threeIsThree = Refl
	Unify(3 = 3, x = x)
	x union 3, = union =, x union 3
*/
