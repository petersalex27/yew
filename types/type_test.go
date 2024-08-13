package types

import (
	"testing"
)

type named string

func (n named) Pos() (start, end int) {
	return 0, 0
}

func (n named) String() string {
	return string(n)
}

func TestTypes(t *testing.T) {
	t.Cleanup(_clearUidCounter)

	A := MakeVar(named("A"), nextEnvUid(), Unrestricted, nil)
	u := Type1
	x := typingChain(named("x"), Unrestricted, A, u)
	B := typingChain(named("B"), Unrestricted, u)
	B.Kind = u
	expected := Pi{false, x, B, u, 0, 0}
	env := NewEnvironment()

	intro, ok := env.Prod(x)
	if !ok {
		t.Fatal("Prod rule failed")
	}
	pi, pass := intro(B)
	if !pass {
		t.Fatal("Prod rule failed on second judgment")
	}
	if !Equals(pi, expected) {
		t.Fatalf("not equal (%v): got %v", expected, pi)
	}
}

func TestApp(t *testing.T) {
	t.Cleanup(_clearUidCounter)

	A := MakeVar(named("A"), nextEnvUid(), Unrestricted, nil)
	u := Type0
	x := typingChain(named("x"), Unrestricted, A, u)
	x0 := typingChain(named("x0"), Unrestricted, A, u)
	B := typingChain(named("B"), Unrestricted, u)
	expected := Pi{false, x0, B, u, 0, 0}

	env := NewEnvironment()
	intro, ok := env.Prod(x)
	if !ok {
		t.Fatal("Prod rule failed")
	}
	pi, pass := intro(B)
	if !pass {
		t.Fatal("Prod rule failed on second judgment")
	}
	if !Equals(pi, expected) {
		t.Fatalf("not equal (%v): got %v", expected, pi)
	}

	f := DummyVar("f")
	b := DummyVar("b")

	var lam Lambda

	lam, pi, ok = env.Abs(f, A)(b, B)(pi)
	if !ok {
		t.Fatal("Abstraction failed")
	}

	res, ty, applied := env.App(lam, pi, x)
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
