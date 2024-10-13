package data

import (
	"testing"

	"github.com/petersalex27/yew/api"
)

type dummyNode struct{}

func (dummyNode) GetPos() api.Position { return api.Position{} }

func (dummyNode) Pos() (int, int) { return 0, 0 }

func (dummyNode) Type() api.NodeType { return dataType("_dummy_node_type_") }

func _[e EmbedsMaybe[a], a api.Node]() api.DescribableNode        { return e{} }
func _[e EmbedsList[a], a api.Node]() api.DescribableNode         { return e{} }
func _[e EmbedsNonEmpty[a], a api.Node]() api.DescribableNode     { return e{} }
func _[e EmbedsErr]() api.DescribableNode                         { return e{} }
func _[e EmbedsSolo[a], a api.Node]() api.DescribableNode         { return e{} }
func _[e EmbedsPair[a, b], a, b api.Node]() api.DescribableNode   { return e{} }
func _[e EmbedsEither[a, b], a, b api.Node]() api.DescribableNode { return e{} }

func Test_assert_DescribableNode(t *testing.T) {
	var (
		_ api.DescribableNode = nothing[dummyNode]{}
		_ api.DescribableNode = just[dummyNode]{}
		_ api.DescribableNode = inLeft[dummyNode, dummyNode]{}
		_ api.DescribableNode = inRight[dummyNode, dummyNode]{}
		_ api.DescribableNode = Solo[dummyNode]{}
		_ api.DescribableNode = Pair[dummyNode, dummyNode]{}
		_ api.DescribableNode = Err{}
		_ api.DescribableNode = NonEmpty[dummyNode]{}
		_ api.DescribableNode = List[dummyNode]{}
		_ api.DescribableNode = Maybe[dummyNode](nil)
		_ api.DescribableNode = Either[dummyNode, dummyNode](nil)
	)
	// yippee!
}
