package parser

import (
	"testing"

	"github.com/petersalex27/yew/internal/errors"
	"github.com/petersalex27/yew/internal/source"
	"github.com/petersalex27/yew/internal/token"
)

// Test is kinda half-assed. Need to not just test against strings--need to test actual structures
func TestStandardActions(t *testing.T) {
	// 	Nat : * where
	//		Zero
	//		Succ Nat
	p := Init(source.SourceCode{})

	add_ := token.Infix.MakeValued("(+)")
	mul_ := token.Infix.MakeValued("(*)")
	pow_ := token.Infix.MakeValued("(**)")
	add := token.Id.MakeValued("+")
	mul := token.Id.MakeValued("*")
	pow := token.Id.MakeValued("**")
	fun := token.Id.MakeValued("fun")

	x := token.Id.MakeValued("x")
	y := token.Id.MakeValued("y")
	z := token.Id.MakeValued("z")

	backslash := token.Backslash.Make()
	thickArrow := token.ThickArrow.Make()
	comma := token.Comma.Make()

	lparen := token.LeftParen.Make()
	rparen := token.RightParen.Make()

	Int := token.Id.MakeValued("Int")

	arrow := token.Arrow.Make()

	// .. : Int -> Int
	var IntToInt TypeElem = TypeElem{
		Type: []token.Token{Int, arrow, Int},
	}
	// .. : Int -> Int -> Int
	var IntToIntToInt TypeElem = TypeElem{
		Type: []token.Token{Int, arrow, Int, arrow, Int},
	}

	// (+) : Int -> Int -> Int
	addDec := DeclarationElem{
		Vis: Open,
		Name: add_,
		rAssoc: false,
		bp: 6,
		Typing: IntToIntToInt,
	}
	if !addDec.Parse(p) {
		t.Fatalf("failed to declare (+)")
	}

	// (*) : Int -> Int -> Int
	mulDec := DeclarationElem{
		Vis: Open,
		Name: mul_,
		rAssoc: false,
		bp: 7,
		Typing: IntToIntToInt,
	}
	if !mulDec.Parse(p) {
		t.Fatalf("failed to declare (*)")
	}

	// (**) : Int -> Int -> Int
	powDec := DeclarationElem{
		Vis: Open,
		Name: pow_,
		rAssoc: true,
		bp: 9,
		Typing: IntToIntToInt,
	}
	if !powDec.Parse(p) {
		t.Fatalf("failed to declare (**)")
	}

	funDec := DeclarationElem{
		Vis: Open,
		Name: fun,
		rAssoc: false,
		bp: 10,
		Typing: IntToInt,
	}
	if !funDec.Parse(p) {
		t.Fatalf("failed to declare fun")
	}

	tests := []struct {
		tokens   []token.Token
		expected string
	}{
		{
			[]token.Token{x, pow, y, pow, z},
			`(**) x (**) y z`,
		},
		{
			[]token.Token{x, add, y, mul, z},
			`(+) x (*) y z`,
		},
		{
			[]token.Token{x, add, y, add, z},
			`(+) (+) x y z`,
		},
		{
			[]token.Token{fun, x, add, y, add, z},
			`(+) (+) fun x y z`,
		},
		{
			// x + (y + z)
			[]token.Token{x, add, lparen, y, add, z, rparen},
			`(+) x (+) y z`,
		},
		{
			// fun (x + y) + z
			[]token.Token{fun, lparen, x, add, y, rparen, add, z},
			`(+) fun (+) x y z`,
		},
		{
			[]token.Token{lparen, add, rparen, x, y},
			`(+) x y`,
		},
		{
			[]token.Token{backslash, x, thickArrow, x},
			`\x => x`,
		},
		{
			[]token.Token{lparen, backslash, x, thickArrow, x, rparen},
			`\x => x`,
		},
		{
			[]token.Token{backslash, x, thickArrow, x, add, x},
			`\x => (+) x x`,
		},
		{
			[]token.Token{backslash, x, thickArrow, backslash, y, thickArrow, x},
			`\x, y => x`,
		},
		{
			[]token.Token{backslash, x, comma, y, thickArrow, x},
			`\x, y => x`,
		},
	}

	for _, test := range tests {
		term, ok := p.Process(standardActions, test.tokens)
		if !ok {
			errors.PrintErrors(p.FlushMessages())
			t.Fatalf("testing: %v\nparsing failed with above messages", test)
		}

		actual := term.String()
		if test.expected != actual {
			t.Fatalf("expected:\n%s\ngot:\n%s", test.expected, actual)
		}
	}
}

// Test is kinda half-assed. Need to not just test against strings--need to test actual structures
func TestTypeActions(t *testing.T) {
	p := Init(source.SourceCode{})

	arrow := token.Arrow.Make()

	x := token.Id.MakeValued("x")
	a := token.Id.MakeValued("a")
	b := token.Id.MakeValued("b")
	Type := token.Id.MakeValued("Type")
	zero := token.IntValue.MakeValued("0")
	one := token.IntValue.MakeValued("1")
	eq := token.Equal.Make()

	//comma := token.Comma.Make()
	colon := token.Colon.Make()

	lparen := token.LeftParen.Make()
	rparen := token.RightParen.Make()
	lbrace := token.LeftBrace.Make()
	rbrace := token.RightBrace.Make()

	tests := []struct {
		tokens   []token.Token
		expected string
	}{
		{
			[]token.Token{a, arrow, b},
			`a -> b`,
		},
		{
			[]token.Token{lparen, x, colon, a, rparen, arrow, b},
			`(x : a) -> b`,
		},
		{
			[]token.Token{lparen, b, colon, Type, arrow, Type, rparen, arrow, b, a, arrow, a},
			`(b : Type -> Type) -> b a -> a`,
		},
		{
			[]token.Token{zero, eq, one, arrow, a},
			`= 0 1 -> a`,
		},
		{
			// {x : a} -> x = x
			[]token.Token{lbrace, x, colon, a, rbrace, arrow, x, eq, x},
			`{x : a} -> = x x`,
		},
	}

	p.parsingTypeSig = true
	p.env.DisableGeneralization()
	for _, test := range tests {
		term, ok := p.Process(typingActions, test.tokens)

		if !ok {
			errors.PrintErrors(p.FlushMessages())
			t.Fatalf("parsing failed with above messages")
		}

		actual := term.String()
		if test.expected != actual {
			t.Fatalf("expected:\n%s\ngot:\n%s", test.expected, actual)
		}

		p.debug_incTestCounter()
	}
}
