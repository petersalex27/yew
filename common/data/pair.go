package data

import "github.com/petersalex27/yew/api"

type (
	Pair[a, b api.Node] struct {
		first  a
		second b
		api.Position
	}

	EmbedsPair[a, b api.Node] interface{
		api.DescribableNode
		~struct{ Pair[a, b] }
	}
)

func MakePair[a, b api.Node](first a, second b) Pair[a, b] {
	p := api.WeakenRangeOver(first.GetPos(), second.GetPos())
	return Pair[a, b]{first: first, second: second, Position: p}
}

// constructs a `pair` node from a `solo` node and a node
func Two[a, b api.Node](first Solo[a], second b) Pair[a, b] {
	return Pair[a, b]{first: first.one, second: second, Position: api.WeakenRangeOver[api.Node](first, second)}
}

// returns left and right nodes of a `pair` node
func (p Pair[a, b]) Split() (a, b) {
	return p.first, p.second
}

// returns the first node of a `pair` node
func (p Pair[a, b]) Fst() a { return p.first }

// returns the second node of a `pair` node
func (p Pair[a, b]) Snd() b { return p.second }

// make an embedded `pair` node
func EMakePair[c EmbedsPair[a, b], a, b api.Node](x a, y b) c {
	return c{MakePair(x, y)}
}

func ETwo[pair EmbedsPair[a, b], a, b api.Node](x Solo[a], y b) pair {
	return pair{Two(x, y)}
}
