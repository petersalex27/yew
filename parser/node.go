// =================================================================================================
// Alex Peters - March 03, 2024
// =================================================================================================
package parser

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew/types"
)

type (
	// has many uses, e.g.,
	// 	- temporary markers on the parse stack
	//	- user defined syntax keywords
	Key struct {
		Name       string
		Start, End int
	}

	Marker struct {
		nodeType   NodeType
		Start, End int
	}

	// identifier of some kind
	Ident struct {
		Name       string
		Start, End int
	}

	// wildcard symbol:
	//	_
	Wildcard struct {
		Start, End int
	}

	// TODO: remove?
	Affixed struct {
		Parts      []namePattern
		Start, End int
	}

	// literal integer characters (unsigned, technically), e.g.,
	//	100245
	IntConst struct {
		int        types.IntConst
		Start, End int
	}

	// literal characters, e.g.,
	//	'a'
	CharConst struct {
		char       types.CharConst
		Start, End int
	}

	// literal floating point numbers, e.g.,
	//	3.14159
	// and
	//	31.4159e-1
	FloatConst struct {
		float      types.FloatConst
		Start, End int
	}

	// literal strings, e.g,
	//	"hello, world"
	StringConst struct {
		string     types.StringConst
		Start, End int
	}

	ConstrainedType struct {
		Constraint  Tuple
		Constrained Type
	}

	// function types of any order, e.g.,
	//	Int -> Int
	// and
	//	* -> Uint -> *
	FunctionType struct {
		Left, Right     Term
		availableInDefs *declTable
		Start, End      int
	}

	// applications that don't fit into other node categories, e.g.,
	//	add 1 2
	// but also infix applications, e.g.,
	//	1 + 2
	Application struct {
		Left, Right Term
		Start, End  int
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
		Elements   []Term
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

	// lambda abstraction, e.g.,
	//	\x, y => e
	Lambda struct {
		Binders    []Ident
		Bound      Term
		Start, End int
	}

	Binding struct {
		Binder, Bound Term
		Start, End    int
	}

	namePattern interface {
		match(namePattern) bool
	}

	// a typing, e.g.,
	//	x : Int
	// but also more complex typings and type/data constructors, e.g.,
	//	_=_ : a -> b -> *
	// and
	//	Refl : x = x
	Typing struct {
		multiplicity types.Multiplicity
		Term
		Type       Term
		Start, End int
	}

	// some kind of term enclosed by parens
	EnclosedTerm struct {
		Term
		Start, End int
	}

	Implicit EnclosedTerm

	Term interface {
		Node
		String() string
		Translate(parser *Parser) types.Term
		term_()
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

func (m Marker) NodeType() NodeType { return m.nodeType }

func (m Marker) String() string {
	switch m.nodeType {
	case closeParenType:
		return "(...)"
	case closeBracketType:
		return "[...]"
	case closeBraceType:
		return "{...}"
	default:
		return "_marker_"
	}
}

func (m Marker) Pos() (int, int) {
	return m.Start, m.End
}

func (Wildcard) String() string { return "_" }

func calcArity(term Term) (arity uint) {
	var f FunctionType
	var ok bool = true
	for f, ok = term.(FunctionType); ok; f, ok = term.(FunctionType) {
		arity++
		term = f.Right
	}
	return
}

// func takesArgs(term Term) (ok bool) {
// 	_, ok = term.(FunctionType)
// 	return
// }

// func takeArg(ft FunctionType) (argTyp Term, retTyp Term) {
// 	return ft.Left, ft.Right
// }

const (
	identType NodeType = iota
	wildcardType
	affixedType
	intConstType
	charConstType
	floatConstType
	stringConstType
	functionType
	applicationType
	lambdaType
	bindingType
	typingType
	listType
	listExprType
	tupleType
	tupleExprType
	pairsType
	listingType
	constrainedTypeType
	implicitType
	closeParenType
	closeBracketType
	closeBraceType

	syntaxExtensionType
)

func (Ident) type_()           {}
func (FunctionType) type_()    {}
func (Application) type_()     {}
func (AmbiguousTuple) type_()  {}
func (AmbiguousList) type_()   {}
func (ConstrainedType) type_() {}
func (Implicit) type_()        {}

func (Ident) expr_()       {}
func (Application) expr_() {}
func (Lambda) expr_()      {}
func (StringConst) expr_() {}
func (FloatConst) expr_()  {}
func (CharConst) expr_()   {}
func (IntConst) expr_()    {}

func (Ident) term_()           {}
func (Application) term_()     {}
func (FunctionType) term_()    {}
func (Lambda) term_()          {}
func (StringConst) term_()     {}
func (FloatConst) term_()      {}
func (CharConst) term_()       {}
func (IntConst) term_()        {}
func (AmbiguousTuple) term_()  {}
func (AmbiguousList) term_()   {}
func (Tuple) term_()           {}
func (List) term_()            {}
func (Key) term_()             {}
func (EnclosedTerm) term_()    {}
func (Listing) term_()         {}
func (ConstrainedType) term_() {}
func (Implicit) term_()        {}
func (Marker) term_()          {}
func (Wildcard) term_()        {}

func stringJoinTerms[T Term](ts []T, sep string) string {
	var b strings.Builder
	switch len(ts) {
	case 0:
		return ""
	}

	b.WriteString(ts[0].String())
	for _, elem := range ts[1:] {
		b.WriteString(sep)
		b.WriteString(elem.String())
	}
	return b.String()
}

func String(t Term) string {
	if t == nil {
		return "_?_"
	}
	return t.String()
}

func (i Implicit) String() string {
	return "{" + String(i.Term) + "}"
}

func (t Typing) String() string {
	left, right := String(t.Term), String(t.Type)
	return left + " : " + right
}

func (id Ident) String() string {
	return id.Name
}

func (a Application) String() string {
	left, right := String(a.Left), String(a.Right)
	return left + " " + right
}

func (f FunctionType) String() string {
	left, right := String(f.Left), String(f.Right)
	return left + " -> " + right
}

func (l Lambda) String() string {
	return fmt.Sprintf("\\%s => %s", stringJoinTerms(l.Binders, ", "), String(l.Bound))
}

func (s StringConst) String() string {
	return s.string.String()
}

func (f FloatConst) String() string {
	return f.float.String()
}

func (c CharConst) String() string {
	return c.char.String()
}

func (i IntConst) String() string {
	return i.int.String()
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

func (c ConstrainedType) String() string {
	arrow, constrainers := " => ", ""
	switch len(c.Constraint.Elements) {
	case 0:
		arrow = ""
	case 1:
		constrainers = c.Constraint.Elements[0].String()
	default:
		constrainers = c.Constraint.String()
	}
	return constrainers + arrow + c.Constrained.String()
}

func (e EnclosedTerm) String() string {
	enclosed := String(e.Term)
	return "(" + enclosed + ")"
}

func (key Key) String() string {
	return key.Name
}

func (Listing) NodeType() NodeType { return listingType }

func (ConstrainedType) NodeType() NodeType { return constrainedTypeType }

func (Key) NodeType() NodeType {
	return syntaxExtensionType
}

func (Lambda) NodeType() NodeType { return lambdaType }

func (Implicit) NodeType() NodeType { return implicitType }

func (k Key) Pos() (int, int) {
	return k.Start, k.End
}

func (lambda Lambda) Pos() (int, int) {
	return lambda.Start, lambda.End
}

func (e EnclosedTerm) Pos() (int, int) {
	return e.Start, e.End
}

func (c ConstrainedType) Pos() (int, int) {
	start := c.Constraint.Start
	_, end := c.Constrained.Pos()
	return start, end
}

func (ls Listing) Pos() (int, int) {
	return ls.Start, ls.End
}

func (i Implicit) Pos() (int, int) {
	return i.Start, i.End
}
