package types

import (
	"fmt"
	"testing"
)

func TestTypes(t *testing.T) {
	t.Cleanup(_clearUidCounter)

	A := Var("A")
	B := Var("B")
	x := Var("x")
	u := Type1
	expected := Pi{
		binderVar:  x,
		binderType: A,
		dependent:  B,
		kind:       u,
	}
	env := NewEnvironment()
	env.assume(A, Type1)
	env.assume(B, Type1)

	f, ok := env.Prod(x, A, u)
	if !ok {
		t.Fatal("Prod rule failed")
	}
	pi, c, pass := f(B, u)
	if !pass {
		t.Fatal("Prod rule failed on second judgment")
	}
	if !Equals(pi, expected) {
		t.Fatalf("not equal (%v): got %v", expected, pi)
	}
	if !Equals(c, u) {
		t.Fatalf("not equal (%v): got %v", u, c)
	}
	fmt.Printf("type: %v\n", pi)
}

func TestApp(t *testing.T) {
	t.Cleanup(_clearUidCounter)

	A := Var("A")
	B := Var("B")
	x := Var("x")
	u := Type0
	expected := Pi{
		binderVar:  x,
		binderType: A,
		dependent:  B,
		kind:       u,
	}

	env := NewEnvironment()
	env.assume(A, u)
	env.assume(B, u)
	createProd, ok := env.Prod(x, A, u)
	if !ok {
		t.Fatal("Prod rule failed")
	}
	pi, c, pass := createProd(B, u)
	if !pass {
		t.Fatal("Prod rule failed on second judgment")
	}
	if !Equals(pi, expected) {
		t.Fatalf("not equal (%v): got %v", expected, pi)
	}
	if !Equals(c, u) {
		t.Fatalf("not equal (%v): got %v", u, c)
	}
	fmt.Printf("type: %v\n", pi)

	f := Var("f")
	b := Constant("b")

	var lam Lambda

	defineFunc := env.Abs(f, A, pi)
	lam, pi, ok = defineFunc(b, B)
	if !ok {
		t.Fatal("Abstraction failed")
	}

	res, ty, applied := env.App(lam, pi, x, A)
	if !applied {
		t.Fatal("Application failed")
	}
	if !Equals(ty, B) {
		t.Fatalf("not equal (%v): got %v", B, ty)
	}
	if res.String() != "b" {
		t.Fatalf("not equal (b): got %v", res)
	}
}

func TestRefl(t *testing.T) {
	// t.Cleanup(_clearUidCounter)

	// Equal := Constant("Equal")
	// //Refl := Constant("Refl")
	// env := NewEnvironment()
	
	// // create constructor for Equal
	// {
	// 	B := Var("b")
	// 	A := Var("a")
	// 	// {a : *} -> a -> {b : *} -> b -> *

	// 	expected := ImplicitBind(A, Type0).To(
	// 		Bind(__, A).To(
	// 			ImplicitBind(B, Type0).To(
	// 				Bind(__, B).To(Type0)(Type1),
	// 				)(Type1))(Type1))(Type1)

	// 	createProd, ok := env.Prod(__, A, Type0)
	// 	if !ok {
	// 		t.Fatal("Prod rule failed")
	// 	}
	// 	var innerPi, innerKind Type
	// 	if createInnerProd, ok2 := env.Prod(__, B, Type0); ok2 {
	// 		innerPi, innerKind, ok = createInnerProd(Type0, Type1)
	// 		if !ok {
	// 			t.Fatal("inner product abstraction failed")
	// 		}
	// 	} else {
	// 		t.Fatal("inner product assumptions failed")
	// 	}

	// 	equalConstructor, _, good := createProd(innerPi, innerKind)
	// 	if !good {
	// 		t.Fatalf("failed to create Equal constructor, Equal : a -> b -> *")
	// 	}

	// 	if !Equals(expected, equalConstructor) {
	// 		t.Fatalf("%v != %v", expected, equalConstructor)
	// 	}

	// 	fmt.Printf("type: %v\n", equalConstructor)
	// 	env.Map(Equal, equalConstructor)
	// }

	// // create data constructor Refl
	// {
	// 	x := Var("x")
	// 	a := Var("a")
	// 	expected := Pi{
	// 		implicit: true,
	// 		binderVar:  a,
	// 		binderType: Type0,
	// 		dependent: Pi{
	// 			implicit: true,
	// 			binderVar: x,
	// 			binderType: a,
	// 			dependent:  Application{[]Term{Equal, x, x}, Type0},
	// 			kind:       Type0,
	// 		},
	// 		kind: Type0,
	// 	}

	// 	createProd0, ok0 := env.ProdImplicit(a, Type0, Type1)
	// 	if !ok0 {
	// 		t.Fatal("Prod rule failed!")
	// 	}
	// 	createProd1, ok1 := env.ProdImplicit(x, a, Type0)
	// 	if !ok1 {
	// 		t.Fatal("Prod rule failed!")
	// 	}	
	// }
}

// func TestWhere(t *testing.T) {
// 	t.Cleanup(_clearUidCounter)

// 	A := Constant("A")
// 	B := Constant("B")
// 	x := Var("x")
// 	u := Constant("*1")
// 	f := Var("f")
// 	// m : Module
// 	m := Var("m")
// 	Module := Constant("Module")
// 	// (x:A) -> B
// 	pi := Pi{
// 		binderVar:  x,
// 		binderType: A,
// 		dependent:  B,
// 		kind:       u,
// 	}
// 	env := NewEnvironment()
// 	assign, eEnv := env.Where(Typing{f, pi})
// 	if eEnv == nil {
// 		t.Fatal("Where-init rule failed")
// 	}

// 	var bind func(Term, Type) (TermBind, bool)
// 	var ok bool
// 	bind, ok = assign(Equality{f, x})
// 	if !ok {
// 		t.Fatalf("Where-assign rule failed")
// 	}

// 	//var term TermBind
// 	_, ok = bind(m, Module)
// 	if !ok {
// 		t.Fatalf("Where-bind rule failed")
// 	}
// }
