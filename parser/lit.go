// =================================================================================================
// Alex Peters - March 03, 2024
// =================================================================================================
package parser

// int, char, float, string

// =================================================================================================
// Int
// =================================================================================================

// position in source
func (c IntConst) Pos() (start, end int) {
	return c.Start, c.End
}

// node type
func (IntConst) NodeType() NodeType {
	return intConstType
}

// =================================================================================================
// Char
// =================================================================================================

// position in source
func (c CharConst) Pos() (start, end int) {
	return c.Start, c.End
}

// node type
func (CharConst) NodeType() NodeType {
	return charConstType
}

// =================================================================================================
// Floating point
// =================================================================================================

// position in source
func (c FloatConst) Pos() (start, end int) {
	return c.Start, c.End
}

// node type
func (FloatConst) NodeType() NodeType {
	return floatConstType
}

// =================================================================================================
// String
// =================================================================================================

// position in source
func (c StringConst) Pos() (start, end int) {
	return c.Start, c.End
}

// node type
func (StringConst) NodeType() NodeType {
	return stringConstType
}