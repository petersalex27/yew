package util

import (
	"io"

	"github.com/petersalex27/yew/api"
)

const (
	padding     string = "    "
	branch      string = "│   "
	childBranch string = "├── "
	finalChild  string = "└── "
)

func Reclassify[T api.Node](nodes []T) []api.Node {
	result := make([]api.Node, len(nodes))
	for i, node := range nodes {
		result[i] = node
	}
	return result
}

// recursively print the tree
func printTree_r(w io.Writer, n api.Node, lhs string) {
	name, children := Describe(n)

	w.Write([]byte(name))

	if len(children) == 0 {
		return
	}

	var position []byte
	for _, child := range children[:len(children)-1] {
		position = []byte("\n" + lhs + childBranch)
		w.Write(position)
		printTree_r(w, child, lhs+branch)
	}

	position = []byte("\n" + lhs + finalChild)
	w.Write(position)
	printTree_r(w, children[len(children)-1], lhs+padding)
}

func PrintTree(w io.Writer, n api.Node) {
	printTree_r(w, n, "")
}
