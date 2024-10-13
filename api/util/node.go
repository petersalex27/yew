package util

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew/api"
)

// ExposeNodeSurface returns a string representation of the node's exposed top level data (no children)
//
// If the node has children, the children will be represented as "[...]"; otherwise, they will be represented as "[]"
//
// Example:
//	"Node{line: 1, column: 1, name: foo, children: [...]}"
func ExposeNodeSurface(node api.Node) string {
	name, children := Describe(node)
	line, col := node.Pos()
	childrenStr := "[]"
	if len(children) > 0 {
		childrenStr = "[...]"
	}
	return fmt.Sprintf("Node{line: %d, col: %d, name: %s, children: %s}", line, col, name, childrenStr)
}

func exposeNode_r(node api.Node, b *strings.Builder) {
	name, children := Describe(node)
	line, col := node.Pos()
	b.WriteString(fmt.Sprintf("Node{line: %d, col: %d, name: %s, children: [", line, col, name))
	for i, child := range children {
		if i > 0 {
			b.WriteString(", ")
		}
		exposeNode_r(child, b)
	}
	b.WriteString("]}")
}

// ExposeNode returns a string representation of the node's exposed data
//
// Example:
//	"Node{line: 1, column: 1, name: foo, children: [Node{line: 1, column: 5, name: bar, children: []}]}"
func ExposeNode(node api.Node) string {
	b := &strings.Builder{}
	exposeNode_r(node, b)
	return b.String()
}