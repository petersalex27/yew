// =================================================================================================
// Alex Peters - February 27, 2024
// =================================================================================================
package types

import (
	"fmt"
	"io"
	"sync"

	"github.com/petersalex27/yew/common/stack"
	"github.com/petersalex27/yew/common/table"
	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/source"
)

// Type system S, 3-tuple
//
//	S = (C, A, R)
//
// where
//   - C is the set of constants,
//   - A is the set of axioms using C, and
//   - R is the set of rules using C
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

// standard type system
//
//	S = (C, A, R)
//
// where
//
//	C = {Type, Type 1}
//	A = {Axiom}
//	R = {
//		Type ~> Type : Type,
//		Type 1 ~> Type 1 : Type 1,
//		Type ~> Type 1 : Type 1,
//		Type 1 ~> Type : Type,
//	}
var std = System{
	// primitives
	constants: map[Sort]Sort{
		Type0: Type0,
		Type1: Type1,
	},
	// axioms
	axioms: map[Sort]Sort{
		Type0: Type1,
	},
	// rules for determining what sort something is
	// 		s ~> t : u
	rules: map[Sort]map[Sort]Sort{
		Type0: { // Type ~> ..
			Type0: Type0, // Type ~> Type : Type
			Type1: Type1, // Type ~> Type{1} : Type{1}
		},
		Type1: { // Type{1} ~> ..
			Type0: Type0, // Type{1} ~> Type : Type
			Type1: Type1, // Type{1} ~> Type{1} : Type{1}
		},
	},
	implicitKind: Type0,
}

// axiom:
//
//	------------- (Axiom)
//	Type : Type 1
func (sys System) Axiom(c Sort) (s Sort, t Sort, ok bool) {
	s, ok = sys.constants[c]
	if !ok {
		return
	}
	t, ok = sys.axioms[s]
	return
	//sys.axioms[]
}

// given sorts s and t, returns the sort u of
//
//	(s ~> t) : u
//
// where `~>` is the relation between any two types a, b such that
//
//	a -> b
//
// and
//
//	a : s, b : t
//
// then
//
//	s ~> t
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

type Replacement struct {
	Term Term
	Type Type
}

type holeGenerator struct {
	baseLower                        byte
	baseUpper                        byte
	iterationsLower, iterationsUpper uint
	ctrLower, ctrUpper               byte
}

func nextHole(ctr *byte, iterations *uint, base byte) Variable {
	const numLetters byte = 26
	hole := fmt.Sprintf("%c", base+*ctr)
	if *iterations > 0 {
		hole = hole + fmt.Sprint(*iterations)
	}

	*ctr = (*ctr + 1) % numLetters
	if *ctr == 0 {
		*iterations++
	}

	ty := dummyVarType(hole, false)
	ty.isHole = true
	return Variable{true, "?" + hole, nextEnvUid(), Erase, ty, 0, 0}
}

func (g *holeGenerator) nextTerm() Variable {
	return nextHole(&g.ctrLower, &g.iterationsLower, g.baseLower)
}

func (g *holeGenerator) nextKind() Variable {
	return nextHole(&g.ctrUpper, &g.iterationsUpper, g.baseUpper)
}

func (env *Environment) NextTermHole() Variable {
	g, _ := env.holes.Peek()
	return g.nextTerm()
}

func (env *Environment) NextKindHole() Variable {
	g, _ := env.holes.Peek()
	return g.nextKind()
}

func NewHoleGenerator() *holeGenerator {
	return &holeGenerator{
		baseLower:       'a',
		baseUpper:       'A',
		iterationsLower: 0,
		iterationsUpper: 0,
		ctrLower:        0,
		ctrUpper:        0,
	}
}

// Proof environment--symbolically, this is Γ
type Environment struct {
	allowGeneralization bool
	src      source.SourceCode
	imported *table.Table[Constant, *Environment]
	// this is more general than just type and data constructors--this holds all things that return a
	// term of a certain type
	constructors *table.MultiTable[Constant, *table.Table[fmt.Stringer, Replacement]]
	// typings for names
	//typings *table.Table[fmt.Stringer, Type]
	//
	replacements *table.MultiTable[fmt.Stringer, Replacement]
	// unification assignments
	unifications *table.MultiTable[Variable, Term]
	holes        *stack.Stack[*holeGenerator]
	// apply func for quick access (nil if not yet set)
	applyFunction *Lambda
	// type system
	system   System
	messages []errors.ErrorMessage
	envLock sync.Mutex
}

// Creates a quick way to access the apply function that is used for type inference; this should be
// set in the environment's replacement table, but isn't required. If it's not set in the table, an
// warning is reported. If this is intentional, call `SetApplyFunction_noWarning` instead.
//
// panics if the apply function is already set
//
// suggested name for `name` is `$`, but this is not required nor checked. This is simply convention
// and what documentation will use to refer to this function
func (env *Environment) SetApplyFunction(name strPos, f *Lambda) {
	if _, found := env.replacements.Find(name); !found {
		env.warning(Warn_ApplyFunctionNotInEnvironment, name)
	}
	env.SetApplyFunction_noWarning(f)
}

// See `SetApplyFunction`
func (env *Environment) SetApplyFunction_noWarning(f *Lambda) {
	if env.applyFunction != nil {
		panic("bug: apply function already set")
	}
	env.applyFunction = f
}

func (env *Environment) Messages() []errors.ErrorMessage {
	return env.messages
}

func (env *Environment) FlushMessages() []errors.ErrorMessage {
	messages := env.Messages()
	env.messages = []errors.ErrorMessage{}
	return messages
}

func (env *Environment) TypeOf(id fmt.Stringer) (ty Type, found bool) {
	replacement, ok := env.replacements.Find(id)
	found = ok
	if ok {
		ty = replacement.Type
	}
	return
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

// constructor must exist in the environment
func (env *Environment) mapConstructor(constructing Constant, constructorName strPos) bool {
	var constructor Replacement
	constructor, ok := env.replacements.Find(constructorName)
	if !ok {
		panic("bug: constructor not found; must declare constructor before recording it as a constructor")
	}

	var constructors *table.Table[fmt.Stringer, Replacement]
	constructors, ok = env.constructors.Find(constructing)
	if !ok {
		constructors = table.MakeTable[fmt.Stringer, Replacement](8)
	} else {
		// check if constructor already in map (in else b/c will never be true if this is the first constructor)
		if _, ok = constructors.Find(constructorName); ok {
			env.error(RedefinedConstructor, constructorName)
			return false
		}
	}

	constructors.Map(constructorName, constructor)
	env.constructors.Map(constructing, constructors)
	return true 
}

func NewEnvironment() *Environment {
	env := new(Environment)
	*env = Environment{
		allowGeneralization: true,
		imported:     table.MakeTable[Constant, *Environment](8),
		constructors: table.NewMultiTable[Constant, *table.Table[fmt.Stringer, Replacement]](16),
		replacements: table.NewMultiTable[fmt.Stringer, Replacement](16),
		unifications: table.NewMultiTable[Variable, Term](16),
		holes:        stack.NewStack[*holeGenerator](16),
		system:       std,
		messages:     []errors.ErrorMessage{},
	}
	env.holes.Push(NewHoleGenerator())
	return env
}

// lol
//
// generates a Pi type:
//
//	(𝚷x:A.y):u
func (env *Environment) PiHole(A Type) Pi {
	x := Hole("x")
	y := Hole("y")
	u := env.NextKindHole()
	x.Kind = A
	return Pi{binderVar: x, dependent: y, kind: u}
}

// check if v occurs in t. If it does, return true; else, return false.
// if t = v, then v is not in t, v is t
func (env *Environment) occurs(v Variable, t Term) (vOccursInT bool) {
	if t == nil {
		panic("bug: t is nil")
	}
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
func (env *Environment) union(v Variable, t Term, unify bool) bool {
	if t == nil {
		panic("bug: t is nil")
	}
	if env.occurs(v, t) {
		return false
	}

	// map v :-> t
	// if v2, isVar := t.(Variable); isVar && Equals(v, v2) {
	// 	return true // don't clutter the unifications table with trivial unifications
	// }
	if !unify {
		return true
	}

	env.unifications.Map(v, t)
	return true
}

// at this point it is known that both ta and tb belong to a unified kind
func (env *Environment) substitute(ta, tb Term, unify bool) (ok bool, doUnify bool) {
	ok, doUnify = true, true
	if v, isVar := ta.(Variable); isVar && v.isHole{
		ok, doUnify = env.union(v, tb, unify), false
	} else if v, isVar := tb.(Variable); isVar && v.isHole {
		ok, doUnify = env.union(v, ta, unify), false
	}

	return
}

func (env *Environment) splitUnify(ta, tb Term, unify bool) (ok bool) {
	// get constants and terms
	ca, termsA := Split(ta)
	cb, termsB := Split(tb)

	if ok = ca == cb; !ok {
		// error: mismatch constants
		env.mismatchUnifyingError(ta, tb)
		return false
	}

	if ok = len(termsA) == len(termsB); !ok {
		// error: patterns of different length, unification impossible
		env.impossibleUnificationBcOfLength(ta, tb, termsA, termsB)
		return false
	}

	// it. through all terms while stat is ok, unifying terms
	for i := 0; ok && i < len(termsA); i++ {
		a, b := termsA[i], termsB[i]
		ok = env.unifyAction(a, b, unify)
	}
	return ok
}

// tries to find typing for a reduced form of term `a`
//
// TODO: until termination/totality "checking" is implemented, this function has no
// defense against non-termination
func (env *Environment) FindUnified(a Term) (t Term) {
	ra := reduce(a)
	v, ok := ra.(Variable)
	if ok = ok && v.isHole; ok {
		t, ok = env.unifications.Find(v)
	}

	// for both failure to find and non-variable reduction result
	if !ok {
		t = ra
	}

	return
}

func (env *Environment) IsFree(v Variable) bool {
	_, found := env.unifications.Find(v)
	return !found
}

func (env *Environment) BeginClause() {
	env.unifications.Increase()
	env.constructors.Increase()
	env.replacements.Increase()
	env.holes.Push(NewHoleGenerator())
}

func (env *Environment) IncreaseWith(local Locals) {
	env.replacements.IncreaseWith(local.Table)
}

func (env *Environment) WithDecrease() {
	env.replacements.Decrease()
}

func (env *Environment) KillClause() table.Table[fmt.Stringer, Replacement] {
	env.unifications.Decrease()
	env.constructors.Decrease()
	out, _ := env.replacements.Decrease()
	env.holes.Pop()
	return out
}

func (env *Environment) DisableGeneralization() {
	env.allowGeneralization = false
}

type Locals struct {
	prefix string
	table.Table[fmt.Stringer, Replacement]
}

func (env *Environment) EndClause(prefix string) Locals {
	locals := Locals{prefix, env.KillClause()}
	return locals
}

func (env *Environment) Get(name fmt.Stringer) (t Term, A Type, found bool) {
	var r Replacement
	if r, found = env.replacements.Find(name); found {
		t, A = r.Term, r.Type
	}
	return
}

func (env *Environment) equateKinds(ka, kb Type) bool {
	ok := env.betaEquivalence(ka, kb)
	if !ok {
		env.equivalenceError(ka, kb, 0, 0) // TODO: not 0, 0
	}
	return ok
}

func (env *Environment) PrintUnifications(w io.Writer) {
	fmt.Fprintf(w, "\n==========================\nUnifications:\n==========================\n")
	env.unifications.Walk(
		func(v Variable, t Term) {
			fmt.Fprintf(w, "%v = %v\n", v, t)
		},
	)
	fmt.Fprintf(w, "==========================\n\n")
}

func (env *Environment) Unifiable(a, b Term) bool {
	return env.unifyAction(a, b, false)
}

func (env *Environment) unifyAction(a, b Term, unify bool) (ok bool) {
	ta := env.FindUnified(a)
	tb := env.FindUnified(b)

	var doSplitUnify bool
	ok, doSplitUnify = env.substitute(ta, tb, unify)
	if ok && doSplitUnify {
		ok = env.splitUnify(ta, tb, unify)
	}
	return
}

// unifies two terms a, b
func (env *Environment) Unify(a, b Term) (ok bool) {
	return env.unifyAction(a, b, true)
}

/*
threeIsThree : 3 = 3
threeIsThree = Refl
	Unify(3 = 3, x = x)
	x union 3, = union =, x union 3
*/
