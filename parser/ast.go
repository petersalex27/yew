// =================================================================================================
// Alex Peters - January 24, 2024
//
// Type definitions and methods for Ast nodes
// =================================================================================================

package parser

import "github.com/petersalex27/yew/token"

type Node interface {
	Pos() (start, end int)
}

// all expressions implement this
type ExprNode interface {
	Node
	exprNode()
}

// all types implement this
type TypeNode interface {
	Node
	typeNode()
}

type ListLike = struct {
	Start, End int
	Elems      []ExprNode
}

type (
	// x y
	Application ListLike

	// a -> b
	ArrowType struct {
		Left, Right TypeNode
	}

	// ()
	BoringKind struct {
		Start, End int
	}

	// ()
	BoringType struct {
		Start, End int
	}

	Constant struct {
		token.Token
	}

	// x
	Ident struct {
		token.Token
	}

	// import math
	Import struct {
		Start, End int
		ImportName token.Token
		LookupName token.Token
	}

	// Just
	KindIdent struct {
		token.Token
	}

	// (\x -> x)
	Lambda struct {
		Start, End int
		Binders    []Ident
		Bound      ExprNode
	}

	// [1, 2, 3]
	List ListLike

	// module main where { ...
	Module struct {
		Start, End int
		ModuleName token.Token
	}

	// (e)
	ParenExpr struct {
		Start, End int
		ExprNode
	}

	TupleKind ListLike

	// Int
	TypeConstant struct {
		token.Token
	}

	// x: Int
	TypeJudgment struct {
		ExprNode
		TypeNode
	}

	// a
	TypeVariable struct {
		token.Token
	}
)

func (a Application) Pos() (start, end int) {
	return a.Start, a.End
}

func (a ArrowType) Pos() (start, end int) {
	start, _ = a.Left.Pos()
	_, end = a.Right.Pos()
	return
}

func (b BoringKind) Pos() (start, end int) {
	return b.Start, b.End
}

func (b BoringType) Pos() (start, end int) {
	return b.Start, b.End
}

func (c Constant) Pos() (start, end int) {
	return c.Start, c.End
}

func (id Ident) Pos() (start, end int) {
	return id.Start, id.End
}

func (im Import) Pos() (start, end int) {
	return im.Start, im.End
}

func (kind KindIdent) Pos() (start, end int) {
	return kind.Start, kind.End
}

func (lambda Lambda) Pos() (start, end int) {
	return lambda.Start, lambda.End
}

func (ls List) Pos() (start, end int) {
	return ls.Start, ls.End
}

func (module Module) Pos() (start, end int) {
	return module.Start, module.End
}

func (paren ParenExpr) Pos() (start, end int) {
	return paren.Start, paren.End
}

func (t TupleKind) Pos() (start, end int) {
	return t.Start, t.End
}

func (c TypeConstant) Pos() (start, end int) {
	return c.Start, c.End
}

func (t TypeJudgment) Pos() (start, end int) {
	start, _ = t.ExprNode.Pos()
	_, end = t.TypeNode.Pos()
	return
}

func (v TypeVariable) Pos() (start, end int) {
	return v.Start, v.End
}

// expressions
func (Application) exprNode() {}
func (BoringKind) exprNode()  {}
func (Constant) exprNode()    {}
func (Ident) exprNode()       {}
func (KindIdent) exprNode()   {}
func (Lambda) exprNode()      {}
func (List) exprNode()        {}
func (ParenExpr) exprNode()   {}
func (TupleKind) exprNode()   {}

// types
func (ArrowType) typeNode()    {}
func (BoringType) typeNode()   {}
func (TypeConstant) typeNode() {}
func (TypeVariable) typeNode() {}
