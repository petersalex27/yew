// =================================================================================================
// Alex Peters - March 03, 2024
// =================================================================================================
package parser

func (ident Ident) Pos() (start, end int) { 
	return ident.Start, ident.End
}

func (ident Ident) NodeType() NodeType {
	return identType
}