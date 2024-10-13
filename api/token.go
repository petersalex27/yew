package api

// a terminal node (leaf of AST), a token
type Token interface {
	// return the value of the token
	String() string
	// if token represents an error, return a non-nil error; otherwise return nil
	Error() error
	Node
}

type DescribableToken interface {
	Token
	Describable[Node]
}