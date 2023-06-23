package bp

type BindingPower int

const (
	None BindingPower = iota
	PatternMatch
	ExpressionAnotation
	Mapping
	Disjunctive
	Conjunctive
	Equitable
	Ordered
	Additive
	Multiplicative
	Unary
	Power
	Postfix
	Compose
	Applicative
	Group
)
