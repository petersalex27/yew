package types

import (
	"fmt"
	"os"
	"testing"
)

var appValid = Application{Var("C"), Var("a")} // len == 2, tau @ 0, tau @ 1
var appInvalid1 = Application{}
var appInvalid2 = Application{Var("C")}
var appInvalid3 = Application{Var("C"), Var("a"), Var("b")}
var appInvalid4 = Application{Int{}, Var("a")}
var appInvalid5 = Application{Var("C"), Int{}}

func TestValidClass(t *testing.T) {
	/*
		len(a) < 2
			return false, "too few type variables, expected one", a
		len(a) > 2
			return false, "too many type variables, expected one", a[2]

		a[0].GetTypeType() != TAU
			return false, "expected class declaration", a[0]

		if a[1].GetTypeType() != TAU
			return false, "expected type variable", a[1]
	
		return true, "", a
	*/
	
	expectations := []struct{
		expect bool
		app Application
		msg string
		loc Types
	}{
		{true, appValid, "", appValid},
		{false, appInvalid1, "too few type variables, expected one", appInvalid1}, // len == 0
		{false, appInvalid2, "too few type variables, expected one", appInvalid2}, // len < 2
		{false, appInvalid3, "too many type variables, expected one", appInvalid3[2]}, // len > 2
		{false, appInvalid4, "expected class declaration", appInvalid4[0]}, // tau !@ 0
		{false, appInvalid5, "expected type variable", appInvalid5[1]}, // tau !@ 1
	}

	for _, expected := range expectations {
		actual, msg, loc := expected.app.ValidClass()
		if actual != expected.expect {
			fmt.Fprintf(os.Stderr, "Expected:\n%t\nActual:\n%t\n",
				expected.expect, actual)
			t.FailNow()
		}
		if msg != expected.msg {
			fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n",
				expected.msg, msg)
			t.FailNow()
		}
		if ty, ok := loc.(Types); ok {
			if !ty.Equals(expected.loc) {
				fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n",
					expected.loc.ToString(), ty.ToString())
				t.FailNow()
			}
		} else {
			fmt.Fprintf(os.Stderr, "Expected:\n%v\nActual:\n%v\n",
				expected.loc, loc)
			t.FailNow()
		}
	}
}

func TestConstrainApplication(t *testing.T) {
	var1 := Var("MyClass")
	var2 := Var("a")
	constraint := Constraint{Context: ConstraintContext{MakeContext("Num", "a")},}
	appValid := Application{var1, var2}
	classGood := Class{Name: "MyClass", TypeVariable: var2}
	classBad := Class{}
	expectations := []struct{
		app Application
		class Class
		expect bool
		msg string
		loc Types
	}{
		{ app: appValid, class: classGood, expect: true, msg: "", loc: classGood, },
		{appInvalid1, classBad, false, "too few type variables, expected one", appInvalid1},
		{appInvalid2, classBad, false, "too few type variables, expected one", appInvalid2},
		{appInvalid3, classBad, false, "too many type variables, expected one", appInvalid3[2]},
		{appInvalid4, classBad, false, "expected class declaration", appInvalid4[0]},
		{appInvalid5, classBad, false, "expected type variable", appInvalid5[1]},
	}
	
	// valid tests 
	for _, expected := range expectations {
		class, actual, msg, loc := constraint.ConstrainApplication(expected.app)
		if !class.Equals(expected.class) {
			fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n",
				expected.class.ToString(), class.ToString())
		}
		if actual != expected.expect {
			fmt.Fprintf(os.Stderr, "Expected:\n%t\nActual:\n%t\n",
				expected.expect, actual)
			t.FailNow()
		}
		if msg != expected.msg {
			fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n",
				expected.msg, msg)
			t.FailNow()
		}
		if ty, ok := loc.(Types); ok {
			if !ty.Equals(expected.loc) {
				fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n",
					expected.loc.ToString(), ty.ToString())
				t.FailNow()
			}
		} else {
			fmt.Fprintf(os.Stderr, "Expected:\n%v\nActual:\n%v\n",
				expected.loc, loc)
			t.FailNow()
		}
	}

	// invalid tests
}