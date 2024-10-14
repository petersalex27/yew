package logic

type Value int8

const (
	Neg  Value = -1
	Nil  Value = 0
	Unit Value = 1

	// alias for Neg
	False Value = Neg
	// alias for Nil
	Unknown Value = Nil
	// alias for Unit
	True Value = Unit
)

func Not(a Value) Value {
	return -a
}

func (a Value) And(b Value) Value {
	return min(a, b)
}

func (a Value) Or(b Value) Value {
	return max(a, b)
}

func (a Value) Xor(b Value) Value {
	lhs := a.And(Not(b))
	rhs := Not(a).And(b)
	return lhs.Or(rhs)
}

// Material implication for Kleene logic
//
// a.Implies(b) == Not(a).Or(b)
//
// Truth table:
//  ```
//  a.Implies(b)
//  ┌───╥───┬───┬───┐
//  │a\b║ - ┆ 0 ┆ + │
//  ╞═══╬═══╪═══╪═══╡
//  │ - ║ + │ + │ + │
//  ├┄┄┄╫───┼───┼───┤
//  │ 0 ║ 0 │ 0 │ + │
//  ├┄┄┄╫───┼───┼───┤
//  │ + ║ - │ 0 │ + │
//  └┄┄┄╨───┴───┴───┘
//  ```
func (a Value) Implies(b Value) Value {
	return Not(a).Or(b)
}

// Material implication for Łukasiewicz logic
//
// Truth table:
//  ```
//  a.ImpliesLuk(b)
//  ┌───╥───┬───┬───┐
//  │a\b║ - ┆ 0 ┆ + │
//  ╞═══╬═══╪═══╪═══╡
//  │ - ║ + │ + │ + │
//  ├┄┄┄╫───┼───┼───┤
//  │ 0 ║ 0 │ + │ + │
//  ├┄┄┄╫───┼───┼───┤
//  │ + ║ - │ 0 │ + │
//  └┄┄┄╨───┴───┴───┘
//  ```
func (a Value) ImpliesLuk(b Value) Value {
	return min(Unit, Unit-a+b)
}