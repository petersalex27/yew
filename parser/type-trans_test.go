package parser

import (
	"fmt"
	"os"
	"testing"

	//"github.com/llir/llvm/ir/types"
	"github.com/petersalex27/yew/source"
	"github.com/petersalex27/yew/token"
)

func TestTranslateSimplyTypedType(t *testing.T) {
	// 	Nat : * where
	//		Zero
	//		Succ Nat
	p := Init(source.SourceCode{})
	add_ := token.Affixed.MakeValued("_+_")
	mul_ := token.Affixed.MakeValued("_*_")
	pow_ := token.Affixed.MakeValued("_**_")
	add := token.Id.MakeValued("+")
	mul := token.Id.MakeValued("*")
	pow := token.Id.MakeValued("**")
	fun := token.Id.MakeValued("fun")

	x := token.ImplicitId.MakeValued("x")
	y := token.ImplicitId.MakeValued("y")
	z := token.ImplicitId.MakeValued("z")

	bslash := token.Backslash.Make()
	tarrow := token.ThickArrow.Make()
	comma := token.Comma.Make()

	lparen := token.LeftParen.Make()
	rparen := token.RightParen.Make()
	Int := Ident{Name: "Int"}
	iii := FunctionType{
		Left:  Int,
		Right: FunctionType{Left: Int, Right: Int},
	}
	i2 := FunctionType{
		Left: Int,
		Right: Int,
	}
	setter, _ := p.declare(add_)
	setter(iii, 6)
	setter, _ = p.declare(mul_)
	setter(iii, 7)
	setter, _ = p.declare(pow_)
	setter(iii, 9, 1)
	setter, _ = p.declare(fun)
	setter(i2, 10)

	term, ok := p.Process(standardActions, []token.Token{x, add, y, mul, z})
	if !ok {
		t.Fatal(p.FlushMessages())
	}
	fmt.Fprintf(os.Stderr, "x + y * z :=> %v\n", term)

	term, ok = p.Process(standardActions, []token.Token{x, pow, y, pow, z})
	if !ok {
		t.Fatal(p.FlushMessages())
	}
	fmt.Fprintf(os.Stderr, "x ** y ** z :=> %v\n", term)

	term, ok = p.Process(standardActions, []token.Token{x, add, y, add, z})
	if !ok {
		t.Fatal(p.FlushMessages())
	}
	fmt.Fprintf(os.Stderr, "x + y + z :=> %v\n", term)

	term, ok = p.Process(standardActions, []token.Token{fun, x, add, y, add, z})
	if !ok {
		t.Fatal(p.FlushMessages())
	}
	fmt.Fprintf(os.Stderr, "fun x + y + z :=> %v\n", term)

	term, ok = p.Process(standardActions, []token.Token{fun, lparen, x, add, y, rparen, add, z})
	if !ok {
		t.Fatal(p.FlushMessages())
	}
	fmt.Fprintf(os.Stderr, "fun (x + y) + z :=> %v\n", term)

	term, ok = p.Process(standardActions, []token.Token{lparen, add, rparen, x, y})
	if !ok {
		t.Fatal(p.FlushMessages())
	}
	fmt.Fprintf(os.Stderr, "(+) x y :=> %v\n", term)

	term, ok = p.Process(standardActions, []token.Token{bslash, x, tarrow, x})
	if !ok {
		t.Fatal(p.FlushMessages())
	}
	fmt.Fprintf(os.Stderr, `\x => x :=> %v%s`, term, "\n")

	term, ok = p.Process(standardActions, []token.Token{lparen, bslash, x, tarrow, x, rparen})
	if !ok {
		t.Fatal(p.FlushMessages())
	}
	fmt.Fprintf(os.Stderr, `(\x => x) :=> %v%s`, term, "\n")

	term, ok = p.Process(standardActions, []token.Token{bslash, x, tarrow, x, add, x})
	if !ok {
		t.Fatal(p.FlushMessages())
	}
	fmt.Fprintf(os.Stderr, `(\x => x + x) :=> %v%s`, term, "\n")

	term, ok = p.Process(standardActions, []token.Token{bslash, x, comma, y, tarrow, x})
	if !ok {
		t.Fatal(p.FlushMessages())
	}
	fmt.Fprintf(os.Stderr, `\x, y => x :=> %v%s`, term, "\n")

	term, ok = p.Process(standardActions, []token.Token{bslash, x, tarrow, bslash, y, tarrow, x})
	if !ok {
		t.Fatal(p.FlushMessages())
	}
	fmt.Fprintf(os.Stderr, `\x => \y => x :=> %v%s`, term, "\n")
}