package data

import "github.com/petersalex27/yew/api"

type (
	Solo[a api.Node] struct {
		one a
		api.Position
	}

	EmbedsSolo[a api.Node] interface {
		api.DescribableNode
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

func SoloMap[a, b api.Node](f func(a) b) func(Solo[a]) Solo[b] {
	return func(s Solo[a]) Solo[b] {
		return One(f(s.one))
	}
}