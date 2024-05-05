package inf

import (
	"testing"

	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/testutil"
)

// tests variable inference rule
func TestVar(t *testing.T) {
	var v0 types.Variable[nameable.Testable]
	var ve0 expr.Variable[nameable.Testable]

	{
		// block prevents accidental use of cxt
		cxt := NewTestableContext()
		v0 = cxt.TypeContext.NewVar()
		ve0 = cxt.ExprContext.NewVar()
	}

	xName := nameable.MakeTestable("x")
	arrName := nameable.MakeTestable("Array")
	aName := nameable.MakeTestable("a")
	nName := nameable.MakeTestable("n")
	uintName := nameable.MakeTestable("Uint")

	x := expr.Const[nameable.Testable]{Name: xName}
	Array := types.MakeConst(arrName)                   // Array
	a := types.Var(aName)                               // a
	Array_a := types.Apply[nameable.Testable](Array, a) // Array a
	n := expr.Var(nName)                                // n
	Uint := types.MakeConst(uintName)                   // Uint
	n_Uint := types.Judgment(expr.Referable[nameable.Testable](n), types.Type[nameable.Testable](Uint))
	ve0_Uint := types.Judgment(expr.Referable[nameable.Testable](ve0), types.Type[nameable.Testable](Uint))
	var_n_Uint := types.Judgment[nameable.Testable, expr.Variable[nameable.Testable]](n, Uint)
	domain := []types.ExpressionJudgment[nameable.Testable, expr.Referable[nameable.Testable]]{n_Uint}
	domain2 := []types.ExpressionJudgment[nameable.Testable, expr.Referable[nameable.Testable]]{ve0_Uint}
	Array_a_n := types.Index(Array_a, domain...)       // (Array a; n)
	Array_a_ve0 := types.Index(Array_a, domain2...)    // (Array a; x0)
	mapval_n_Uint__Array_a := types.MakeDependentType( // mapval (n: Uint) . (Array a)
		[]types.TypeJudgment[nameable.Testable, expr.Variable[nameable.Testable]]{var_n_Uint},
		types.TypeFunction[nameable.Testable](types.Index(Array_a)),
	)
	Array_v0 := types.Apply[nameable.Testable](Array, v0)
	Array_v0_ve0 := types.Index(Array_v0, domain2...)

	tests := []struct {
		description string
		input       bridge.JudgmentAsExpression[nameable.Testable, expr.Const[nameable.Testable]]
		expect      Conclusion[nameable.Testable, expr.Const[nameable.Testable], types.Monotyped[nameable.Testable]]
	}{
		{
			"x: Array => x: Array",
			bridge.Judgment[nameable.Testable, expr.Const[nameable.Testable]](x, Array),
			Conclude[nameable.Testable](x, types.Monotyped[nameable.Testable](Array)),
		},
		{
			"x: a => x: a",
			bridge.Judgment[nameable.Testable, expr.Const[nameable.Testable]](x, a),
			Conclude[nameable.Testable](x, types.Monotyped[nameable.Testable](a)),
		},
		{
			"x: Array a => x: Array a",
			bridge.Judgment[nameable.Testable, expr.Const[nameable.Testable]](x, Array_a),
			Conclude[nameable.Testable](x, types.Monotyped[nameable.Testable](Array_a)),
		},
		{
			"x: forall a . a => x: $0",
			bridge.Judgment[nameable.Testable, expr.Const[nameable.Testable]](x, types.Forall(a).Bind(a)),
			Conclude[nameable.Testable](x, types.Monotyped[nameable.Testable](v0)),
		},
		{
			"x: forall a . Array a => x: Array $0",
			bridge.Judgment[nameable.Testable, expr.Const[nameable.Testable]](x, types.Forall(a).Bind(Array_a)),
			Conclude[nameable.Testable](x, types.Monotyped[nameable.Testable](Array_v0)),
		},
		{
			"x: (Array a; n) => x: (Array a; n)",
			bridge.Judgment[nameable.Testable, expr.Const[nameable.Testable]](x, Array_a_n),
			Conclude[nameable.Testable](x, types.Monotyped[nameable.Testable](Array_a_n)),
		},
		{
			"x: mapval (n: Uint) . (Array a) => x: (Array a; $0)",
			bridge.Judgment[nameable.Testable, expr.Const[nameable.Testable]](x, mapval_n_Uint__Array_a),
			Conclude[nameable.Testable](x, types.Monotyped[nameable.Testable](Array_a_ve0)),
		},
		{
			"x: forall a . mapval (n: Uint) . (Array a) => x: (Array $0; $e0)",
			bridge.Judgment[nameable.Testable, expr.Const[nameable.Testable]](x, types.Forall(a).Bind(mapval_n_Uint__Array_a)),
			Conclude[nameable.Testable](x, types.Monotyped[nameable.Testable](Array_v0_ve0)),
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.varBody(test.input)

		eq := types.JudgmentEquals[nameable.Testable, expr.Const[nameable.Testable], types.Type[nameable.Testable]](
			actual.Judgment().AsTypeJudgment(),
			test.expect.Judgment().AsTypeJudgment(),
		)
		if !eq {
			t.Fatal(
				testutil.
					Testing("equality", test.description).
					FailMessage(test.expect, actual, i))
		}
	}
}

func TestVarFail(t *testing.T) {
	xName := nameable.MakeTestable("x")
	x := expr.Const[nameable.Testable]{Name: xName}

	const expect Status = NameNotInContext

	cxt := NewTestableContext()
	actual := cxt.Var(x)
	if !actual.Status.Is(expect) {
		t.Fatal(testutil.Testing("status", "name not in context").FailMessage(expect, actual.Status))
	}
}

func TestApp(t *testing.T) {
	var v0 types.Variable[nameable.Testable]
	arrow := types.MakeInfixConst[nameable.Testable](nameable.MakeTestable("->"))

	{
		// block prevents accidental use of cxt
		cxt := NewTestableContext()
		v0 = cxt.TypeContext.NewVar()
	}

	// names
	xName := nameable.MakeTestable("x")
	yName := nameable.MakeTestable("y")
	tailName := nameable.MakeTestable("tail")
	take1Name := nameable.MakeTestable("take1")
	take2Name := nameable.MakeTestable("take2")
	zeroName := nameable.MakeTestable("0")
	uintName := nameable.MakeTestable("Uint")
	arrayName := nameable.MakeTestable("Array")
	bracketsName := nameable.MakeTestable("[]")
	aName := nameable.MakeTestable("a")
	bName := nameable.MakeTestable("b")
	nName := nameable.MakeTestable("n")
	qName := nameable.MakeTestable("q")
	wName := nameable.MakeTestable("w")
	zName := nameable.MakeTestable("z")
	succName := nameable.MakeTestable("Succ")

	// type constants
	Array := types.MakeConst(arrayName) // Array
	Uint := types.MakeConst(uintName)   // Uint

	// type vars
	a := types.Var(aName) // a
	b := types.Var(bName) // b

	// expr vars
	x := expr.Const[nameable.Testable]{Name: xName} // x
	y := expr.Const[nameable.Testable]{Name: yName} // y
	n := expr.Var(nName)                            // n
	q := expr.Var(qName)                            // q
	w := expr.Var(wName)                            // w
	z := expr.Var(zName)                            // z

	// expr constants
	tail := expr.Const[nameable.Testable]{Name: tailName}   // tail
	take1 := expr.Const[nameable.Testable]{Name: take1Name} // take1
	take2 := expr.Const[nameable.Testable]{Name: take2Name} // take2
	zero := expr.Const[nameable.Testable]{Name: zeroName}   // 0
	Succ := expr.Const[nameable.Testable]{Name: succName}   // Succ

	// expr applications
	arrayEnclose := types.MakeEnclosingConst[nameable.Testable](1, bracketsName) // [_]
	Array_a := types.Apply[nameable.Testable](arrayEnclose, a)                   // [a]
	Array_Uint := types.Apply[nameable.Testable](arrayEnclose, Uint)             // [Uint]

	// var: type
	n_Uint := types.Judgment(expr.Referable[nameable.Testable](n), types.Type[nameable.Testable](Uint)) // n: Uint
	q_Uint := types.Judgment(expr.Referable[nameable.Testable](q), types.Type[nameable.Testable](Uint)) // q: Uint
	w_b := types.Judgment(expr.Referable[nameable.Testable](w), types.Type[nameable.Testable](b))       // w: b
	z_a := types.Judgment(expr.Referable[nameable.Testable](z), types.Type[nameable.Testable](a))       // z: a

	// "n: Uint"/"0: Uint"
	en_Uint := bridge.Judgment(expr.Expression[nameable.Testable](n), types.Type[nameable.Testable](Uint))    // n: Uint
	e0_Uint := bridge.Judgment(expr.Expression[nameable.Testable](zero), types.Type[nameable.Testable](Uint)) // 0: Uint
	r0_Uint := types.Judgment(expr.Referable[nameable.Testable](zero), types.Type[nameable.Testable](Uint))   // 0: Uint

	// Succ1
	Succ_n := bridge.MakeData(Succ, en_Uint)                                                                         // Succ (n: Uint)
	Succ_0 := bridge.MakeData(Succ, e0_Uint)                                                                         // Succ (0: Uint)
	Succ_n_Uint := types.Judgment(expr.Referable[nameable.Testable](Succ_n), types.Type[nameable.Testable](Uint))    // (Succ n): Uint
	eSucc_n_Uint := bridge.Judgment(expr.Expression[nameable.Testable](Succ_n), types.Type[nameable.Testable](Uint)) // (Succ n): Uint
	Succ_0_Uint := types.Judgment(expr.Referable[nameable.Testable](Succ_0), types.Type[nameable.Testable](Uint))    // (Succ 0): Uint
	eSucc_0_Uint := bridge.Judgment(expr.Expression[nameable.Testable](Succ_0), types.Type[nameable.Testable](Uint)) // Succ (0: Uint)

	// Succ2
	Succ2_0 := bridge.MakeData(Succ, eSucc_0_Uint)                                                                  // Succ (Succ 0)
	Succ2_0_Uint := types.Judgment(expr.Referable[nameable.Testable](Succ2_0), types.Type[nameable.Testable](Uint)) // (Succ (Succ 0)): Uint
	Succ2_n := bridge.MakeData(Succ, eSucc_n_Uint)                                                                  // Succ (n: Uint)
	Succ2_n_Uint := types.Judgment(expr.Referable[nameable.Testable](Succ2_n), types.Type[nameable.Testable](Uint)) // (Succ n): Uint

	// domains
	domain := types.Indexes[nameable.Testable]{n_Uint}              // n: Uint
	domainSucc := types.Indexes[nameable.Testable]{Succ_n_Uint}     // (Succ n): Uint
	domainSucc2 := types.Indexes[nameable.Testable]{Succ2_n_Uint}   // (Succ (Succ n)): Uint
	domainSucc0 := types.Indexes[nameable.Testable]{Succ_0_Uint}    // (Succ 0): Uint
	domainSucc2_0 := types.Indexes[nameable.Testable]{Succ2_0_Uint} // (Succ 0): Uint
	domain0 := types.Indexes[nameable.Testable]{r0_Uint}            // 0: Uint
	domainQ := types.Indexes[nameable.Testable]{q_Uint}             // q: Uint
	domainW := types.Indexes[nameable.Testable]{w_b}                // w: b
	domainZ := types.Indexes[nameable.Testable]{z_a}                // z: a

	// [a; _]
	Array_a_n := types.Index(Array_a, domain...)            // [a; n]
	Array_a_Succ_n := types.Index(Array_a, domainSucc...)   // [a; Succ n]
	Array_a_Succ2_n := types.Index(Array_a, domainSucc2...) // [a; Succ n]

	// [Uint; Succ _]
	Array_Uint_Succ_n := types.Index(Array_Uint, domainSucc...)     // [Uint; Succ n]
	Array_Uint_Succ_0 := types.Index(Array_Uint, domainSucc0...)    // [Uint; Succ 0]
	Array_Uint_Succ2_0 := types.Index(Array_Uint, domainSucc2_0...) // [Uint; Succ (Succ 0)]

	// [Uint; 0]
	Array_Uint_0 := types.Index(Array_Uint, domain0...) // [Uint; 0]

	// [Uint; var]
	Array_Uint_q := types.Index(Array_Uint, domainQ...) // [Uint; q]
	Array_Uint_w := types.Index(Array_Uint, domainW...) // [Uint; w]
	Array_Uint_z := types.Index(Array_Uint, domainZ...) // [Uint; z]

	// dependent types
	tailFunc := types.Apply[nameable.Testable](arrow, Array_a_Succ_n, Array_a_n)        // [a; Succ n] -> [a; n]
	take1Func := types.Apply[nameable.Testable](arrow, Array_a_Succ2_n, Array_a_Succ_n) // [a; Succ (Succ n)] -> [a; Succ n]
	take2Func := types.Apply[nameable.Testable](arrow, Array_a_Succ2_n, Array_a_n)      // [a; Succ (Succ n)] -> [a; n]

	tests := []struct {
		description string
		input0      bridge.JudgmentAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
		input1      bridge.JudgmentAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
		findIn      types.Variable[nameable.Testable]
		findOut     types.Type[nameable.Testable]
		expect      Conclusion[nameable.Testable, expr.Application[nameable.Testable], types.Monotyped[nameable.Testable]]
	}{
		{
			"(x: a) (y: Array) => (x y): $0",
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](x, a),
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](y, Array),
			a, types.Apply[nameable.Testable](arrow, Array, v0), // a = Array -> $0
			Conclude[nameable.Testable](
				expr.Apply[nameable.Testable](x, y),
				types.Monotyped[nameable.Testable](v0),
			),
		},
		{
			"(y: b) (x: a) => (y x): $0",
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](y, b),
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](x, a),
			b, types.Apply[nameable.Testable](arrow, a, v0), // b = a -> $0
			Conclude[nameable.Testable](
				expr.Apply[nameable.Testable](y, x),
				types.Monotyped[nameable.Testable](v0),
			),
		},
		{
			"(tail: [a; Succ n] -> [a; n]) (x: [Uint; Succ 0]) => (tail x): [Uint; 0]",
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](tail, tailFunc),
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](x, Array_Uint_Succ_0),
			a, Uint,
			Conclude[nameable.Testable](
				expr.Apply[nameable.Testable](tail, x),
				types.Monotyped[nameable.Testable](Array_Uint_0),
			),
		},
		{
			"(take2: [a; Succ (Succ n)] -> [a; n]) (x: [Uint; Succ (Succ 0)]) => (take2 x): [Uint; 0]",
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](take2, take2Func),
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](x, Array_Uint_Succ2_0),
			a, Uint,
			Conclude[nameable.Testable](
				expr.Apply[nameable.Testable](take2, x),
				types.Monotyped[nameable.Testable](Array_Uint_0),
			),
		},
		{
			"(take1: [a; Succ (Succ n)] -> [a; Succ n]) (x: [Uint; Succ (Succ 0)]) => (take1 x): [Uint; Succ 0]",
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](take1, take1Func),
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](x, Array_Uint_Succ2_0),
			a, Uint,
			Conclude[nameable.Testable](
				expr.Apply[nameable.Testable](take1, x),
				types.Monotyped[nameable.Testable](Array_Uint_Succ_0),
			),
		},
		{
			"(take1: [a; Succ (Succ n)] -> [a; Succ n]) (x: [Uint; (q: Uint)]) => (take1 x): [Uint; Succ n]",
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](take1, take1Func),
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](x, Array_Uint_q),
			a, Uint,
			Conclude[nameable.Testable](
				expr.Apply[nameable.Testable](take1, x),
				types.Monotyped[nameable.Testable](Array_Uint_Succ_n),
			),
		},
		{
			"(take1: [a; Succ (Succ n)] -> [a; Succ n]) (x: [Uint; (w: b)]) => (take1 x): [Uint; Succ n]",
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](take1, take1Func),
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](x, Array_Uint_w),
			b, Uint,
			Conclude[nameable.Testable](
				expr.Apply[nameable.Testable](take1, x),
				types.Monotyped[nameable.Testable](Array_Uint_Succ_n),
			),
		},
		{
			"(take1: [a; Succ (Succ n)] -> [a; Succ n]) (x: [Uint; (z: a)]) => (take1 x): [Uint; Succ n]",
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](take1, take1Func),
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](x, Array_Uint_z),
			a, Uint,
			Conclude[nameable.Testable](
				expr.Apply[nameable.Testable](take1, x),
				types.Monotyped[nameable.Testable](Array_Uint_Succ_n),
			),
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.App(test.input0, test.input1)

		if cxt.HasErrors() {
			t.Fatal(
				testutil.
					Testing("errors", test.description).
					FailMessage(nil, cxt.GetReports(), i))
		}

		findOutActual := cxt.Find(test.findIn)
		eq := test.findOut.Equals(findOutActual)
		if !eq {
			t.Fatal(
				testutil.
					Testing("find", test.description).
					FailMessage(test.findOut, findOutActual, i))
		}

		eq = types.JudgmentEquals[nameable.Testable, expr.Application[nameable.Testable], types.Type[nameable.Testable]](
			actual.Judgment().AsTypeJudgment(),
			test.expect.Judgment().AsTypeJudgment(),
		)
		if !eq {
			t.Fatal(
				testutil.
					Testing("equality", test.description).
					FailMessage(test.expect, actual, i))
		}
	}
}

func TestUnifyStatus(t *testing.T) {
	// names
	zeroName := nameable.MakeTestable("0")
	oneName := nameable.MakeTestable("1")
	uintName := nameable.MakeTestable("Uint")
	bracketsName := nameable.MakeTestable("[]")
	aName := nameable.MakeTestable("a")
	bName := nameable.MakeTestable("b")
	nName := nameable.MakeTestable("n")
	myTypeName := nameable.MakeTestable("MyType")

	// type constants
	MyType := types.MakeConst(myTypeName) // MyType
	Uint := types.MakeConst(uintName)     // Uint

	// type vars
	a := types.Var(aName) // a
	b := types.Var(bName) // b

	// expr vars
	n := expr.Var(nName) // n

	// expr consts
	zero := expr.Const[nameable.Testable]{Name: zeroName} // 0
	one := expr.Const[nameable.Testable]{Name: oneName}   // 1

	// expr applications
	arrayEnclose := types.MakeEnclosingConst[nameable.Testable](1, bracketsName) // [_]
	Array_a := types.Apply[nameable.Testable](arrayEnclose, a)                   // [a]

	// var: type
	n_Uint := types.Judgment(expr.Referable[nameable.Testable](n), types.Type[nameable.Testable](Uint))     // n: Uint
	r0_Uint := types.Judgment(expr.Referable[nameable.Testable](zero), types.Type[nameable.Testable](Uint)) // 0: Uint
	r1_Uint := types.Judgment(expr.Referable[nameable.Testable](one), types.Type[nameable.Testable](Uint))  // 1: Uint

	// domains
	domain := types.Indexes[nameable.Testable]{n_Uint}   // n: Uint
	domain0 := types.Indexes[nameable.Testable]{r0_Uint} // 0: Uint
	domain1 := types.Indexes[nameable.Testable]{r1_Uint} // 1: Uint

	// [a; _]
	Array_a_n := types.Index(Array_a, domain...)  // [a; n]
	Array_a_0 := types.Index(Array_a, domain0...) // [a; 0]
	Array_a_1 := types.Index(Array_a, domain1...) // [a; 1]

	MyType_a_b := types.Apply[nameable.Testable](MyType, a, b)
	MyType_a := types.Apply[nameable.Testable](MyType, a)

	tests := []struct {
		desc        string
		left, right types.Monotyped[nameable.Testable]
		expect      Status
	}{
		{
			"bad constant",
			MyType_a, Array_a,
			ConstantMismatch,
		},
		{
			"bad kind constant",
			Array_a_0, Array_a_1,
			KindConstantMismatch,
		},
		{
			"bad param length",
			MyType_a_b, MyType_a,
			ParamLengthMismatch,
		},
		{
			"bad index length",
			Array_a_n, Array_a,
			IndexLengthMismatch,
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.Unify(test.left, test.right)

		if !test.expect.Is(actual) {
			t.Fatal(testutil.Testing("stat", test.desc).FailMessage(test.expect, actual, i))
		}
	}
}

func TestAbs(t *testing.T) {
	var v0 types.Variable[nameable.Testable]
	var ve0 expr.Variable[nameable.Testable]
	arrow := types.MakeInfixConst[nameable.Testable](nameable.MakeTestable("->"))

	{
		// block prevents accidental use of cxt
		cxt := NewTestableContext()
		v0 = cxt.TypeContext.NewVar()
		ve0 = cxt.ExprContext.NewVar()
	}

	xName := nameable.MakeTestable("x")
	yName := nameable.MakeTestable("y")
	arrName := nameable.MakeTestable("Array")
	aName := nameable.MakeTestable("a")

	x := expr.Const[nameable.Testable]{Name: xName}
	y := expr.Const[nameable.Testable]{Name: yName}
	Array := types.MakeConst(arrName)                   // Array
	a := types.Var(aName)                               // a
	Array_a := types.Apply[nameable.Testable](Array, a) // Array a

	tests := []struct {
		description string
		inputParam  nameable.Testable
		inputExpr   bridge.JudgmentAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
		expect      Conclusion[nameable.Testable, expr.Function[nameable.Testable], types.Monotyped[nameable.Testable]]
	}{
		{
			`x => y: Array => (\$0 -> y): $0 -> Array`,
			xName,
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](y, Array),
			Conclude[nameable.Testable](
				expr.Bind[nameable.Testable](ve0).In(y),
				types.Monotyped[nameable.Testable](types.Apply[nameable.Testable](arrow, v0, Array)),
			),
		},
		{
			`x => (x y): Array => (\$0 -> $0 y): $0 -> Array`,
			xName,
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](expr.Apply[nameable.Testable](x, y), Array),
			Conclude[nameable.Testable](
				expr.Bind[nameable.Testable](ve0).In(expr.Apply[nameable.Testable](ve0, y)),
				types.Monotyped[nameable.Testable](types.Apply[nameable.Testable](arrow, v0, Array)),
			),
		},
		{
			`x => (x y): a => (\$0 -> $0 y): $0 -> a`,
			xName,
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](expr.Apply[nameable.Testable](x, y), a),
			Conclude[nameable.Testable](
				expr.Bind[nameable.Testable](ve0).In(expr.Apply[nameable.Testable](ve0, y)),
				types.Monotyped[nameable.Testable](types.Apply[nameable.Testable](arrow, v0, a)),
			),
		},
		{
			`x => (x y): Array a => (\$0 -> $0 y): $0 -> Array a`,
			xName,
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](expr.Apply[nameable.Testable](x, y), Array_a),
			Conclude[nameable.Testable](
				expr.Bind[nameable.Testable](ve0).In(expr.Apply[nameable.Testable](ve0, y)),
				types.Monotyped[nameable.Testable](types.Apply[nameable.Testable](arrow, v0, Array_a)),
			),
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.Abs(test.inputParam)(test.inputExpr)

		eq := types.JudgmentEquals[nameable.Testable, expr.Function[nameable.Testable], types.Type[nameable.Testable]](
			actual.Judgment().AsTypeJudgment(),
			test.expect.Judgment().AsTypeJudgment(),
		)
		if !eq {
			t.Fatal(
				testutil.
					Testing("equality", test.description).
					FailMessage(test.expect, actual, i))
		}
	}
}

func TestGen(t *testing.T) {
	// var v0 types.Variable[nameable.Testable]

	// {
	// 	// block prevents accidental use of cxt
	// 	cxt := NewTestableContext()
	// 	v0 = cxt.typeContext.NewVar()
	// }

	arrName := nameable.MakeTestable("Array")
	aName := nameable.MakeTestable("a")
	nName := nameable.MakeTestable("n")
	uintName := nameable.MakeTestable("Uint")

	Array := types.MakeConst(arrName)                   // Array
	a := types.Var(aName)                               // a
	Array_a := types.Apply[nameable.Testable](Array, a) // Array a
	n := expr.Var(nName)                                // n
	Uint := types.MakeConst(uintName)                   // Uint
	n_Uint := types.Judgment(expr.Referable[nameable.Testable](n), types.Type[nameable.Testable](Uint))
	var_n_Uint := types.Judgment[nameable.Testable, expr.Variable[nameable.Testable]](n, Uint)
	domain := []types.ExpressionJudgment[nameable.Testable, expr.Referable[nameable.Testable]]{n_Uint}
	vs := []types.TypeJudgment[nameable.Testable, expr.Variable[nameable.Testable]]{var_n_Uint}
	Array_a_n := types.Index(Array_a, domain...) // (Array a; n)
	Array_n := types.Index(types.Apply[nameable.Testable](Array), domain...)

	tests := []struct {
		description string
		in          types.Monotyped[nameable.Testable]
		expect      types.Polytype[nameable.Testable]
	}{
		{
			"Array => forall _ . Array",
			Array,
			types.Forall[nameable.Testable]().Bind(Array),
		},
		{
			"a => forall a . a",
			a,
			types.Forall(a).Bind(a),
		},
		{
			"Array a => forall a . Array a",
			Array_a,
			types.Forall(a).Bind(Array_a),
		},
		{
			"Array; n => forall _ . mapval (n: Uint) . Array",
			Array_n,
			types.Forall[nameable.Testable]().Bind(types.MakeDependentType[nameable.Testable](vs, types.Apply[nameable.Testable](Array))),
		},
		{
			"Array a; n => forall a . mapval (n: Uint) . Array a",
			Array_a_n,
			types.Forall(a).Bind(types.MakeDependentType[nameable.Testable](vs, Array_a)),
		},
	}

	for i, test := range tests {
		cxt := NewContext[nameable.Testable]()
		actual := cxt.Gen(test.in)
		if !actual.Equals(test.expect) {
			t.Fatal(
				testutil.Testing("equality", test.description).
					FailMessage(test.expect, actual, i),
			)
		}
	}
}

func TestLet(t *testing.T) {
	arrow := types.MakeInfixConst[nameable.Testable](nameable.MakeTestable("->"))
	xName := nameable.MakeTestable("x")
	yName := nameable.MakeTestable("y")
	arrName := nameable.MakeTestable("Array")
	aName := nameable.MakeTestable("a")
	zeroName := nameable.MakeTestable("0")
	intName := nameable.MakeTestable("Int")

	x := expr.Const[nameable.Testable]{Name: xName}       // x (constant)
	y := expr.Const[nameable.Testable]{Name: yName}       // y (constant)
	zero := expr.Const[nameable.Testable]{Name: zeroName} // 0 (constant)
	yVar := expr.Var(yName)                               // y (variable)
	idFunc := expr.Bind[nameable.Testable](yVar).In(yVar) // (\y -> y)
	Int := types.MakeConst(intName)                       // Int
	Array := types.MakeConst(arrName)                     // Array
	a := types.Var(aName)                                 // a
	aToA := types.Apply[nameable.Testable](arrow, a, a)   // a -> a
	x_0 := expr.Apply[nameable.Testable](x, zero)         // (x 0)

	tests := []struct {
		description string
		inputParam  nameable.Testable
		inputAssign bridge.JudgmentAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
		inputExpr   bridge.JudgmentAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
		expect      Conclusion[nameable.Testable, expr.NameContext[nameable.Testable], types.Monotyped[nameable.Testable]]
	}{
		{
			`x, y: Array => x: Array => let x = y in x: Array`,
			xName,
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](y, Array),
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](x, Array),
			Conclude[nameable.Testable](
				expr.Let[nameable.Testable](x, y, x),
				types.Monotyped[nameable.Testable](Array),
			),
		},
		{
			`x, (\y -> y): a -> a => (x 0): Int => let x = (\y -> y) in x 0: Int`,
			xName,
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](idFunc, aToA),
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](x_0, Int),
			Conclude[nameable.Testable](
				expr.Let[nameable.Testable](x, idFunc, x_0),
				types.Monotyped[nameable.Testable](Int),
			),
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.Let(test.inputParam, test.inputAssign)(test.inputExpr)

		eq := types.JudgmentEquals[nameable.Testable, expr.NameContext[nameable.Testable], types.Type[nameable.Testable]](
			actual.judgment.AsTypeJudgment(),
			test.expect.judgment.AsTypeJudgment(),
		)
		if !eq {
			t.Fatal(
				testutil.
					Testing("equality", test.description).
					FailMessage(test.expect, actual, i))
		}
	}
}

func TestRec(t *testing.T) {
	arrow := types.MakeInfixConst[nameable.Testable](nameable.MakeTestable("->"))
	xName := nameable.MakeTestable("x")
	yName := nameable.MakeTestable("y")
	fName := nameable.MakeTestable("f")
	gName := nameable.MakeTestable("g")
	v1Name := nameable.MakeTestable("v1")
	v2Name := nameable.MakeTestable("v2")
	v3Name := nameable.MakeTestable("v3")
	addName := nameable.MakeTestable("add")
	arrName := nameable.MakeTestable("Array")
	//aName := nameable.MakeTestable("a")
	twoName := nameable.MakeTestable("2")
	intName := nameable.MakeTestable("Int")

	x := expr.Const[nameable.Testable]{Name: xName}          // x (constant)
	y := expr.Const[nameable.Testable]{Name: yName}          // y (constant)
	f := expr.Const[nameable.Testable]{Name: fName}          // f (constant)
	g := expr.Const[nameable.Testable]{Name: gName}          // g (constant)
	xVar := expr.Var(xName)                                  // x (variable)
	yVar := expr.Var(yName)                                  // y (variable)
	addConst := expr.Const[nameable.Testable]{Name: addName} // add (constant)
	two := expr.Const[nameable.Testable]{Name: twoName}      // 2 (constant)
	add := expr.Bind(yVar, xVar).In(                         // (\y x -> f (add x y))
		expr.Apply[nameable.Testable](
			f,
			expr.Apply[nameable.Testable](addConst, xVar, yVar),
		),
	)
	add2 := add.Apply(two)                                                   // (\x -> f (add x 2))
	composeG_x := expr.Bind(xVar).In(expr.Apply[nameable.Testable](g, xVar)) // (\x -> g x)
	fOfTwo := expr.Apply[nameable.Testable](f, two)                          // f 2

	v1 := types.Var(v1Name)
	v2 := types.Var(v2Name)
	v3 := types.Var(v3Name)
	// idFunc := expr.Bind[nameable.Testable](yVar).In(yVar) // (\y -> y)
	Int := types.MakeConst(intName) // Int
	v1_to_v2 := types.Apply[nameable.Testable](arrow, v1, v2)
	Int_to_v3 := types.Apply[nameable.Testable](arrow, Int, v3)
	Array := types.MakeConst(arrName) // Array
	// a := types.Var(aName)                                 // a
	// aToA := types.Apply[nameable.Testable](arrow, a, a)   // a -> a
	// x_0 := expr.Apply[nameable.Testable](x, zero)         // (x 0)

	tests := []struct {
		description string
		inputParams []nameable.Testable
		inputAssign []TypeJudgment[nameable.Testable]
		inputExpr   bridge.JudgmentAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
		expect      Conclusion[nameable.Testable, expr.RecIn[nameable.Testable], types.Monotyped[nameable.Testable]]
	}{
		{
			`x => y: Array => x: Array => rec x = y in x: Array`,
			[]nameable.Testable{xName},
			[]TypeJudgment[nameable.Testable]{
				bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](y, Array),
			},
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](x, Array),
			Conclude[nameable.Testable](
				expr.Rec[nameable.Testable](expr.Declare(x.Name).Instantiate(y))(x),
				types.Monotyped[nameable.Testable](Array),
			),
		},
		{
			`f, g => (\x -> g x): v1->v2, (\x -> f (add x 2)): Int->v3 => f 0: v3 => rec f x = (g x) and g x = f (add x 2) in f 2: v3`,
			[]nameable.Testable{fName, gName},
			[]TypeJudgment[nameable.Testable]{
				bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](composeG_x, v1_to_v2),
				bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](add2, Int_to_v3),
			},
			bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](fOfTwo, v3),
			Conclude[nameable.Testable](
				expr.Rec[nameable.Testable](
					expr.Declare(fName).Instantiate(composeG_x),
					expr.Declare(gName).Instantiate(add2),
				)(fOfTwo),
				types.Monotyped[nameable.Testable](v3),
			),
		},
	}

	for i, test := range tests {
		cxt := NewTestableContext()
		actual := cxt.Rec(
			test.inputParams,
		)(
			test.inputAssign,
		)(
			test.inputExpr,
		)

		eq := types.JudgmentEquals[nameable.Testable, expr.RecIn[nameable.Testable], types.Type[nameable.Testable]](
			actual.judgment.AsTypeJudgment(),
			test.expect.judgment.AsTypeJudgment(),
		)
		if !eq {
			t.Fatal(
				testutil.
					Testing("equality", test.description).
					FailMessage(test.expect, actual, i))
		}
	}
}

func TestFind(t *testing.T) {
	intName := nameable.MakeTestable("Int")
	myTypeName := nameable.MakeTestable("MyType")
	aName, bName := nameable.MakeTestable("a"), nameable.MakeTestable("b")

	Int := types.MakeConst(intName)
	a, b := types.Var(aName), types.Var(bName)
	MyType := types.MakeConst(myTypeName)
	MyType_a := types.Apply[nameable.Testable](MyType, a)
	MyType_b := types.Apply[nameable.Testable](MyType, b)

	tests := []struct {
		desc string
		targ types.Variable[nameable.Testable]
		sub  types.Monotyped[nameable.Testable]
	}{
		{
			"a = Int",
			a, Int,
		},
		{
			"a = b",
			a, b,
		},
		{
			"a = MyType b",
			a, MyType_b,
		},
		{
			"a = MyType a",
			a, MyType_a,
		},
	}

	for i, test := range tests {
		expectBefore, expectAfter := test.targ, test.sub
		cxt := NewContext[nameable.Testable]()

		// test find before substitution added
		beforeSub := cxt.Find(test.targ)
		if !beforeSub.Equals(expectBefore) {
			t.Fatal(
				testutil.
					Testing("find before sub. added", test.desc).
					FailMessage(expectBefore, beforeSub, i))
		}

		// now add substitution
		cxt.typeSubs.Add(test.targ.GetReferred(), test.sub)

		// test find after substitution added
		afterSub := cxt.Find(test.targ)
		if !afterSub.Equals(expectAfter) {
			t.Fatal(
				testutil.
					Testing("find after sub. added", test.desc).
					FailMessage(expectAfter, afterSub, i))
		}
	}
}

func TestUnify(t *testing.T) {
	type expected struct {
		inTable bool
		in, out types.Monotyped[nameable.Testable]
	}
	intName := nameable.MakeTestable("Int")
	myTypeName := nameable.MakeTestable("MyType")
	myOtherTypeName := nameable.MakeTestable("MyOtherType")
	aName, bName := nameable.MakeTestable("a"), nameable.MakeTestable("b")

	Int := types.MakeConst(intName)
	a, b := types.Var(aName), types.Var(bName)
	MyType := types.MakeConst(myTypeName)
	MyOtherType := types.MakeConst(myOtherTypeName)
	MyType_a_b := types.Apply[nameable.Testable](MyType, a, b)
	MyOtherType_a := types.Apply[nameable.Testable](MyOtherType, a)
	MyOtherType_b := types.Apply[nameable.Testable](MyOtherType, b)
	MyType_a := types.Apply[nameable.Testable](MyType, a)
	MyType_b := types.Apply[nameable.Testable](MyType, b)

	tests := []struct {
		desc        string
		left, right types.Monotyped[nameable.Testable]
		expectStat  Status
		expect      []expected
	}{
		{
			"Unify(a, b)",
			a, b,
			Ok,
			[]expected{
				{true, a, b},
				{false, b, b},
			},
		},
		{
			"Unify(a, Int)",
			a, Int,
			Ok,
			[]expected{
				{true, a, Int},
				{false, Int, Int},
			},
		},
		{
			"Unify(Int, Int)",
			Int, Int,
			Ok,
			[]expected{
				{false, Int, Int},
			},
		},
		{
			"Unify(a, MyType b)",
			a, MyType_b,
			Ok,
			[]expected{
				{true, a, MyType_b},
				{false, MyType_b, MyType_b},
			},
		},
		{
			"Unify(MyType b, MyType b)",
			MyType_b, MyType_b,
			Ok,
			[]expected{
				{false, MyType, MyType},
				{true, b, b},
			},
		},
		{
			"Unify(MyType a, MyType b)",
			MyType_a, MyType_b,
			Ok,
			[]expected{
				{false, MyType, MyType},
				{true, a, b},
				{false, b, b},
			},
		},
		{
			"Unify(MyOtherType b, MyType b)",
			MyOtherType_b, MyType_b,
			ConstantMismatch,
			[]expected{
				{false, MyOtherType, MyOtherType},
				{false, MyType, MyType},
				{false, b, b},
			},
		},
		{
			"Unify(MyOtherType a, MyType b)",
			MyOtherType_a, MyType_b,
			ConstantMismatch,
			[]expected{
				{false, MyOtherType, MyOtherType},
				{false, MyType, MyType},
				{false, a, a},
				{false, b, b},
			},
		},
		{
			"Unify(MyType a b, MyType b)",
			MyType_a_b, MyType_b,
			ParamLengthMismatch,
			[]expected{
				{false, MyType, MyType},
				{false, a, a},
				{false, b, b},
			},
		},
		{
			"Unify(a, MyType a)",
			a, MyType_a,
			OccursCheckFailed,
			[]expected{
				{false, MyType, MyType},
				{false, a, a},
			},
		},
		{
			"Unify(a, MyType a b)",
			a, MyType_a_b,
			OccursCheckFailed,
			[]expected{
				{false, MyType, MyType},
				{false, a, a},
				{false, b, b},
			},
		},
		{
			"Unify(b, MyType a b)",
			b, MyType_a_b,
			OccursCheckFailed,
			[]expected{
				{false, MyType, MyType},
				{false, a, a},
				{false, b, b},
			},
		},
	}

	for i, test := range tests {
		cxt := NewContext[nameable.Testable]()

		stat := cxt.Unify(test.left, test.right)
		if stat != test.expectStat {
			t.Fatal(
				testutil.
					Testing("stat", test.desc).
					FailMessage(test.expectStat, stat, i))
		}

		for j, expect := range test.expect {
			// check if expected value for whether in sub. table
			_, inTable := cxt.typeSubs.Get(expect.in.GetReferred())
			if inTable != expect.inTable {
				t.Fatal(
					testutil.
						Testing("found in sub. table", test.desc).
						FailMessage(expect.inTable, inTable, i, j))
			}

			// check if expected result for find
			out := cxt.Find(expect.in)
			if !out.Equals(expect.out) {
				t.Fatal(
					testutil.
						Testing("find return value", test.desc).
						FailMessage(expect.out, out, i, j))
			}
		}
	}
}

// prove:
//
//	let x = (\y -> y) in x 0: Int
func TestProofValidation(t *testing.T) {
	// prove:
	//	let x = (\y -> y) in x 0: Int
	//
	// full proof:
	//	𝚪 = {0: Int, (λy.y): a -> a}:
	//
	//		  [ x: forall a. a -> a ]¹    Inst(forall a. a -> a)
	//		  -------------------------------------------------- [Var]
	//		                      x: v -> v                       0: Int    t0, Int = v
	//		                      ----------------------------------------------------- [App]
	//		                                              x 0: t0
	//		                                              -------- [Id]
	//		  (λy.y): a -> a                              x 0: Int
	//		1 ---------------------------------------------------- [Let]
	//		               let x = (λy.y) in x 0: Int

	arrow := types.MakeInfixConst[nameable.Testable](nameable.MakeTestable("->"))
	xName := nameable.MakeTestable("x")
	yName := nameable.MakeTestable("y")
	aName := nameable.MakeTestable("a")
	zeroName := nameable.MakeTestable("0")
	intName := nameable.MakeTestable("Int")

	x := expr.Const[nameable.Testable]{Name: xName}       // x (constant)
	zero := expr.Const[nameable.Testable]{Name: zeroName} // 0 (constant)
	yVar := expr.Var(yName)                               // y (variable)
	idFunc := expr.Bind[nameable.Testable](yVar).In(yVar) // (\y -> y)
	Int := types.MakeConst(intName)                       // Int
	a := types.Var(aName)                                 // a
	aToA := types.Apply[nameable.Testable](arrow, a, a)   // a -> a
	x_0 := expr.Apply[nameable.Testable](x, zero)         // (x 0)

	letExpr := expr.Let[nameable.Testable](x, idFunc, x_0)

	cxt := NewTestableContext()

	Γ := struct {
		id, zero bridge.JudgmentAsExpression[nameable.Testable, expr.Expression[nameable.Testable]]
	}{
		// (λy.y): a -> a
		id: bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](idFunc, aToA),
		// 0: Int
		zero: bridge.Judgment[nameable.Testable, expr.Expression[nameable.Testable]](zero, Int),
	}

	// step 1: add context (i.e., `x: Gen(a -> a)`) and first premise for let expression
	step := 0
	discharge_x_assumption := cxt.Let(xName, Γ.id) // returns function that discharges assumption

	// step 2: do Var rule
	step++
	var_x := cxt.Var(x)
	if var_x.NotOk() {
		t.Fatal(testutil.Testing("x assumption get").FailMessage(Ok.String(), var_x.Status.String(), step))
	}

	// step 3: do App rule
	step++
	app_x_0_conclusion := cxt.App(var_x.judgment, Γ.zero)

	if cxt.HasErrors() {
		t.Fatal(testutil.Testing("app rule errors").FailMessage(nil, cxt.GetReports(), step))
	}

	{
		actualExpr, actualType := app_x_0_conclusion.judgment.GetExpressionAndType()
		expectExpr, expectType := x_0, Int

		if !expectExpr.StrictEquals(actualExpr) {
			t.Fatal(testutil.Testing("app rule expression result").FailMessage(expectExpr, actualExpr, step))
		}

		if !expectType.Equals(actualType) {
			t.Fatal(testutil.Testing("app rule type result").FailMessage(expectType, actualType, step))
		}
	}

	// step 4: discharge assumption and introduce let expression w/ type Int
	step++
	conclusion := discharge_x_assumption(app_x_0_conclusion.judgment)
	{
		actualExpr, actualType := conclusion.judgment.GetExpressionAndType()
		expectExpr, expectType := letExpr, Int

		if !expectExpr.StrictEquals(actualExpr) {
			t.Fatal(testutil.Testing("let rule expression result").FailMessage(expectExpr, actualExpr, step))
		}

		if !expectType.Equals(actualType) {
			t.Fatal(testutil.Testing("let rule type result").FailMessage(expectType, actualType, step))
		}
	}
}

// prove:
//
//	tail (0::(0::[])): [Uint; Succ 0]
func TestProof2Validation(t *testing.T) {
	arrow := types.MakeInfixConst[nameable.Testable](nameable.MakeTestable("->"))
	nName := nameable.MakeTestable("n")
	aName := nameable.MakeTestable("a")
	zeroName := nameable.MakeTestable("0")
	uintName := nameable.MakeTestable("Uint")
	tailName := nameable.MakeTestable("tail")
	consName := nameable.MakeTestable("::")
	emptyName := nameable.MakeTestable("[]")
	succName := nameable.MakeTestable("Succ")
	bracketsName := nameable.MakeTestable("[]")

	a := types.Var(aName)
	Uint := types.MakeConst(uintName)

	n := expr.Var(nName)
	zero := expr.Const[nameable.Testable]{Name: zeroName}   // 0
	tail := expr.Const[nameable.Testable]{Name: tailName}   // tail
	cons := expr.Const[nameable.Testable]{Name: consName}   // (::)
	empty := expr.Const[nameable.Testable]{Name: emptyName} // []
	Succ := expr.Const[nameable.Testable]{Name: succName}   // Succ

	arrayEnclose := types.MakeEnclosingConst[nameable.Testable](1, bracketsName) // [_]
	Array_a := types.Apply[nameable.Testable](arrayEnclose, a)                   // [a]
	Array_Uint := types.Apply[nameable.Testable](arrayEnclose, Uint)             // [Uint]

	// var: type
	n_Uint := types.Judgment(expr.Referable[nameable.Testable](n), types.Type[nameable.Testable](Uint)) // n: Uint
	var_n_Uint := types.Judgment(n, types.Type[nameable.Testable](Uint))

	// "n: Uint"/"0: Uint"
	en_Uint := bridge.Judgment(expr.Expression[nameable.Testable](n), types.Type[nameable.Testable](Uint))    // n: Uint
	r0_Uint := types.Judgment(expr.Referable[nameable.Testable](zero), types.Type[nameable.Testable](Uint))   // 0: Uint
	e0_Uint := bridge.Judgment(expr.Expression[nameable.Testable](zero), types.Type[nameable.Testable](Uint)) // 0: Uint

	// Succ
	Succ_n := bridge.MakeData(Succ, en_Uint)                                                                      // Succ (n: Uint)
	Succ_n_Uint := types.Judgment(expr.Referable[nameable.Testable](Succ_n), types.Type[nameable.Testable](Uint)) // (Succ n): Uint
	Succ_0 := bridge.MakeData(Succ, e0_Uint)                                                                      // Succ (0: Uint)
	Succ_0_Uint := types.Judgment(expr.Referable[nameable.Testable](Succ_0), types.Type[nameable.Testable](Uint)) // (Succ 0): Uint

	// domains
	domain := types.Indexes[nameable.Testable]{n_Uint}          // n: Uint
	domainSucc := types.Indexes[nameable.Testable]{Succ_n_Uint} // (Succ n): Uint
	domain0 := types.Indexes[nameable.Testable]{r0_Uint}        // 0: Uint
	domainSucc0 := types.Indexes[nameable.Testable]{Succ_0_Uint}

	// [a; _]
	Array_a_n := types.Index(Array_a, domain...)                 // [a; n]
	Array_a_0 := types.Index(Array_a, domain0...)                // [a; 0]
	Array_a_Succ_n := types.Index(Array_a, domainSucc...)        // [a; Succ n]
	Array_Uint_Succ_0 := types.Index(Array_Uint, domainSucc0...) // [Uint; Succ 0]

	tailFunc_free :=
		types.Apply[nameable.Testable](arrow, Array_a_Succ_n, Array_a_n) // [a; Succ n] -> [a; n]
	tailFunc_mapped := types.Map(var_n_Uint).To(tailFunc_free) // mapval (n: Uint) . [a; Succ n] -> [a; n]
	tailFunc := types.Forall(a).Bind(tailFunc_mapped)          // forall a . mapval (n: Uint) . [a; Succ n] -> [a; n]

	consFunc_free := types.Apply[nameable.Testable](
		arrow, a,
		types.Apply[nameable.Testable](
			arrow, Array_a_n, Array_a_Succ_n)) // a -> [a; n] -> [a; Succ n]
	consFunc_mapped := types.Map(var_n_Uint).To(consFunc_free) // mapval (n: Uint) . a -> [a; n] -> [a; Succ n]
	consFunc := types.Forall(a).Bind(consFunc_mapped)          // forall a . mapval (n: Uint) . a -> [a; n] -> [a; Succ n]
	emptyArrTy := types.Forall(a).Bind(Array_a_0)              // forall a . [a; 0]

	expect := types.TypedJudge[nameable.Testable](
		// tail ((::) 0 ((::) 0 [])) == tail (0::0::[])
		expr.Apply[nameable.Testable](
			tail,
			// (::) 0 ((::) 0 [])
			expr.Apply[nameable.Testable](
				// (::) 0
				expr.Apply[nameable.Testable](
					cons, zero,
				),
				// (::) 0 []
				expr.Apply[nameable.Testable](
					expr.Apply[nameable.Testable](
						cons, zero,
					),
					empty,
				),
			),
		),
		// [Uint; Succ 0]
		Array_Uint_Succ_0,
	)

	cxt := NewTestableContext()

	//	𝚪 = {
	//		0: Uint,
	//		TailTy = forall a. mapval (n: Uint). [a; Succ n] -> [a; n],
	//		tail: TailTy,
	//		ConsTy = forall a. mapval (n: Uint). a -> [a; n] -> [a; Succ n],
	//		(::): ConsTy,
	//		[]: forall a. [a; 0]
	//	}

	// add context
	cxt.Shadow(zero, Uint)        // 0
	cxt.Shadow(tail, tailFunc)    // tail: forall a. mapval (n: Uint). [a; Succ n] -> [a; n]
	cxt.Shadow(cons, consFunc)    // (::): forall a. mapval (n: Uint). a -> [a; n] -> [a; Succ n]
	cxt.Shadow(empty, emptyArrTy) // []: forall a. [a; 0]

	// start by building (0::(0::[]))

	// declare vars needed to build (0::[]), i.e., (::), 0, []

	// step 1: declare (::)
	step := 0
	varCons0 := cxt.Var(cons) // (::): $0 -> [$0; $e0] -> [$0; Succ $e0]

	if varCons0.Status.NotOk() {
		t.Fatal(testutil.Testing("conclusion stat").FailMessage(Ok, varCons0.Status, step))
	}

	//fmt.Printf("step %d: %v\n", step+1, varCons0)

	// step 2: declare 0
	step++
	varZero0 := cxt.Var(zero)

	if varZero0.Status.NotOk() {
		t.Fatal(testutil.Testing("conclusion stat").FailMessage(Ok, varZero0.Status, step))
	}

	//fmt.Printf("step %d: %v\n", step+1, varZero0)

	// step 3: declare []
	step++
	varEmpty := cxt.Var(empty)

	if varEmpty.Status.NotOk() {
		t.Fatal(testutil.Testing("conclusion stat").FailMessage(Ok, varEmpty.Status, step))
	}

	//fmt.Printf("step %d: %v\n", step+1, varEmpty)

	// step 4: apply ((::) 0)
	step++
	app0 := cxt.App(varCons0.judgment, varZero0.judgment)

	if app0.Status.NotOk() {
		t.Fatal(testutil.Testing("conclusion stat").FailMessage(Ok, app0.Status, step))
	}

	//fmt.Printf("step %d: %v\n", step+1, app0)

	// step 5: apply ((::) 0) []
	step++
	app1 := cxt.App(app0.judgment, varEmpty.judgment)

	if app1.Status.NotOk() {
		t.Fatal(testutil.Testing("conclusion stat").FailMessage(Ok, app1.Status, step))
	}

	//fmt.Printf("step %d: %v\n", step+1, app1)

	// step 6: declare (::)
	step++
	varCons1 := cxt.Var(cons)

	if varCons1.Status.NotOk() {
		t.Fatal(testutil.Testing("conclusion stat").FailMessage(Ok, varCons1.Status, step))
	}

	//fmt.Printf("step %d: %v\n", step+1, varCons1)

	// step 7: declare 0
	step++
	varZero1 := cxt.Var(zero)

	if varZero1.Status.NotOk() {
		t.Fatal(testutil.Testing("conclusion stat").FailMessage(Ok, varZero1.Status, step))
	}

	//fmt.Printf("step %d: %v\n", step+1, varZero1)

	// step 8: apply ((::) 0)
	step++
	app2 := cxt.App(varCons1.judgment, varZero1.judgment)

	if app2.Status.NotOk() {
		t.Fatal(testutil.Testing("conclusion stat").FailMessage(Ok, app2.Status, step))
	}

	//fmt.Printf("step %d: %v\n", step+1, app2)

	// step 9: apply ((::) 0 (((::) 0) []))
	step++
	app3 := cxt.App(app2.judgment, app1.judgment)

	if app3.Status.NotOk() {
		t.Fatal(testutil.Testing("app3 stat").FailMessage(Ok, app3.Status, step))
	}

	//fmt.Printf("step %d: %v\n", step+1, app3)

	// step 10: declare tail
	step++
	varTail := cxt.Var(tail)

	if varTail.Status.NotOk() {
		t.Fatal(testutil.Testing("varTail stat").FailMessage(Ok, varTail.Status, step))
	}

	//fmt.Printf("step %d: %v\n", step+1, varTail)

	// step 11: apply tail ((::) 0 (((::) 0) []))
	step++
	conclusion := cxt.App(varTail.judgment, app3.judgment)

	if conclusion.Status.NotOk() {
		t.Fatal(testutil.Testing("conclusion stat").FailMessage(Ok, conclusion.Status, step))
	}

	if cxt.HasErrors() {
		t.Fatal(testutil.Testing("errors").FailMessage(nil, cxt.GetReports(), step))
	}

	//fmt.Printf("step %d: %v\n", step+1, conclusion)

	if !JudgmentsEqual[nameable.Testable](conclusion.judgment, expect) {
		t.Fatal(testutil.Testing("conclusion equality").FailMessage(expect, conclusion.judgment, step))
	}
}
