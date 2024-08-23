// =================================================================================================
// Alex Peters - March 03, 2024
// =================================================================================================
package parser

import (
	"github.com/petersalex27/yew/internal/types"
)

type (
	// identifier of some kind
	Ident struct {
		Name       string
		Start, End int
	}

	// maybe list, maybe list type--unclear
	AmbiguousList struct {
		Element    Term
		Start, End int
	}

	// could be many things, not clear yet
	//
	// stands for terms separated by commas, e.g.,
	//	x, x + 1, y, z
	Listing struct {
		Elements   []types.Term
		Start, End int
	}

	// a list value, e.g.,
	//	[x, x + 1, y, z]
	List Listing

	// maybe tuple, maybe tuple type--unclear
	AmbiguousTuple struct {
		Elements   []Term
		Start, End int
	}

	// a tuple value, e.g.,
	//	(x, x + 1, y, z)
	Tuple Listing

	Binding struct {
		Binder, Bound Term
		Start, End    int
	}

	namePattern interface {
		match(namePattern) bool
	}

	Term interface {
		Node
		String() string
	}

	Type interface {
		Term
		type_()
	}

	Expr interface {
		Node
		expr_()
	}

	NodeType uint8

	Node interface {
		Pos() (start, end int)
		NodeType() NodeType
	}
)

const (
	identType = iota
	intConstType
	charConstType
	floatConstType
	stringConstType
	functionType
	labeledFunctionType
	implicitFunctionType
	applicationType
	lambdaType
	typingType
	listExprType
	tupleType
	tupleExprType
	pairsType
	listingType
	implicitType
	datatypeType
)

func String(t types.Term) string {
	if t == nil {
		return "?"
	}
	return t.String()
}

func stringJoinTerms(ts []types.Term, sep string) string {
	if len(ts) == 0 {
		return ""
	}
	joined := String(ts[0])
	for _, t := range ts[1:] {
		joined += sep + String(t)
	}
	return joined
}

func (id Ident) String() string {
	return id.Name
}

func (t AmbiguousTuple) String() string {
	return "TODO"
}

func (AmbiguousList) String() string {
	return "TODO"
}

func (ps Tuple) String() string {
	return "(" + Listing(ps).String() + ")"
}

func (ls List) String() string {
	return "[" + Listing(ls).String() + "]"
}

func (ls Listing) String() string {
	return stringJoinTerms(ls.Elements, ", ")
}

func (Listing) NodeType() NodeType { return listingType }

func (ls Listing) Pos() (int, int) {
	return ls.Start, ls.End
}