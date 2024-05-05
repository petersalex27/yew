// =================================================================================================
// Alex Peters - March 03, 2024
// =================================================================================================
package parser

// wildcard

// position in source
func (w Wildcard) Pos() (start, end int) {
	return w.Start, w.End
}

// node type
func (Wildcard) NodeType() NodeType {
	return wildcardType
}