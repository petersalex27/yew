package parser

// functions, applications, type ident, typing

// =================================================================================================
// Function Type
// =================================================================================================

// position in source
func (f FunctionType) Pos() (start, end int) {
	return f.Start, f.End
}

// node type
func (FunctionType) NodeType() NodeType {
	return functionType
}

// =================================================================================================
// Application Type
// =================================================================================================

func (a Application) Pos() (start, end int) {
	return a.Start, a.End
}

func (Application) NodeType() NodeType {
	return applicationType
}

// =================================================================================================
// Typing
// =================================================================================================

func (t Typing) Pos() (start, end int) {
	return t.Start, t.End
}

func (Typing) NodeType() NodeType {
	return typingType
}