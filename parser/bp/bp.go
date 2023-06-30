package bp

type BindingPower int

const (
	None BindingPower = iota
	Sequencer
	PatternMatch
	ExpressionAnotation
	Constraint
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
