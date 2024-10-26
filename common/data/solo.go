package data

import "github.com/petersalex27/yew/api"

type (
	Solo[a api.Node] struct {
		one a
		api.Position
	}

	EmbedsSolo[a api.Node] interface {
		api.DescribableNode
		Break() a
		~struct{ Solo[a] }
	}
)

// make an embedded `solo` node
func EOne[solo EmbedsSolo[a], a api.Node](x a) solo {
	return solo{One(x)}
}

// constructs a `solo` node
func One[a api.Node](node a) Solo[a] {
	return Solo[a]{one: node, Position: node.GetPos()}
}

func MapSolo[a, b api.Node](f func(a) b) func(Solo[a]) Solo[b] {
	return func(s Solo[a]) Solo[b] {
		return One(f(s.one))
	}
}

func EmbMapSolo[solo EmbedsSolo[b], a, b api.Node](f func(a) b) func(Solo[a]) solo {
	return func(s Solo[a]) solo { return solo{MapSolo(f)(s)} }
}

func (s Solo[a]) Break() a {
	return s.one
}
