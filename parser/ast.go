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

type (
	// x y
	Application struct {
		Left, Right ExprNode
	}

	// [1, 2, 3]
	Array struct {
		Start, End int
		Elems      []ExprNode
	}

	// a -> b
	ArrowType struct {
		Left, Right TypeNode
	}

	// x
	Id struct {
		token.Token
	}

	// import math
	Import struct {
		Start, End int
		ImportName token.Token
		LookupName token.Token
	}

	// module main where { ...
	Module struct {
		Start, End int
		ModuleName token.Token
	}

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
	start, _ = a.Left.Pos()
	_, end = a.Right.Pos()
	return
}

func (a Array) Pos() (start, end int) {
	return a.Start, a.End
}

func (a ArrowType) Pos() (start, end int) {
	start, _ = a.Left.Pos()
	_, end = a.Right.Pos()
	return
}

func (id Id) Pos() (start, end int) {
	return id.Start, id.End
}

func (im Import) Pos() (start, end int) {
	return im.Start, im.End
}

func (module Module) Pos() (start, end int) {
	return module.Start, module.End
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
func (a Application) exprNode() {}
func (a Array) exprNode()       {}
func (id Id) exprNode()         {}

// types
func (a ArrowType) typeNode()    {}
func (c TypeConstant) typeNode() {}
func (v TypeVariable) typeNode() {}
