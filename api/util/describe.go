package util

import (
	"fmt"

	"github.com/petersalex27/yew/api"
)

type unknownNode struct{}

func (unknownNode) Describe() (string, []api.Node) {
	return "?unknown", nil
}

func (unknownNode) Pos() (int, int) { return 1, 1 }

func (unknownNode) GetPos() api.Position { return api.ZeroPosition() }

func (unknownNode) Type() api.NodeType { return nil }

func Describe(n api.Node) (name string, children []api.Node) {
	if n, ok := n.(api.DescribableNode); ok {
		return n.Describe()
	}

	if n, ok := n.(interface{ String() string }); ok {
		return n.String(), nil
	}

	return fmt.Sprintf("%T", n), []api.Node{unknownNode{}}
}
