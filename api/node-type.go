package api

type NodeType interface {
	// string representation
	String() string
	// check if the node's type matches the receiver
	Match(Node) bool
}

type ReservedNodeType byte

const (
	// empty node type--denotes the absence of a meaningful node type
	Empty ReservedNodeType = 0
	// atomic node type--denotes a single, basic node type
	Atomic ReservedNodeType = 1
	// compound node type--denotes a node type that is composed of other node types
	Compound ReservedNodeType = 2
)
