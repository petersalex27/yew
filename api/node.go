package api

import "github.com/petersalex27/yew/api/util/fun/pair"

// nodes implement this but also tokens should as well
type Describable[T any] interface {
	Describe() (name string, children []T)
}

// a node that can be described by its name and children
//
// this is primarily used for printing a tree
type DescribableNode interface {
	Node
	Describable[Node]
}

type Node interface {
	// return the position (inclusive start index and exclusive end index) of the node in the source code
	Positioned
	// return the type of the token
	//
	// NOTE: this **must** work even on nil receivers!! this is vital!!
	Type() NodeType
}

func NodeTypeString(n Node) string {
	if n == nil {
		return "empty"
	}
	return TypeString(n.Type())
}

func TypeString(n NodeType) string {
	if n == nil {
		return "empty"
	}
	return n.String()
}

type Positioned interface {
	Pos() (int, int)
	GetPos() Position
}

type Position struct{ data pair.Data[int, int] }

func MakePosition(start, end int) Position { return Position{pair.Make(start, end)} }

func (p Position) Update(p2 Positioned) Position {
	return Position{pair.WeakenRange(p.data, p2.GetPos().data)}
}

func (p Position) ZeroNegatives() Position {
	return Position{pair.MaxPositive(p.data, ZeroPosition().data)}
}

// WeakenRangeOver returns a new position that is the weakest range that includes all the given positions
// (not including non-positive values if possible--will never contain negative values, though)
func WeakenRangeOver[a Positioned](x a, ys ...a) Position {
	acc := pair.MaxPositive(x.GetPos().data, ZeroPosition().data) // zero out negative values
	for _, y := range ys {
		acc = pair.WeakenRange(acc, y.GetPos().data)
	}
	return Position{acc}
}

// ZeroPosition returns a position with both values set to 0
func ZeroPosition() Position { return Position{pair.Make(0, 0)} } // technically, can just use Position{} but this is clearer

func (p Position) Pos() (int, int) { return p.data.Both() }

func (p Position) GetPos() Position { return p }
