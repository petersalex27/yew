package parser

// ListType, List, TupleType, Tuple, Pairs

func (ls AmbiguousList) Pos() (start, end int) {
	return ls.Start, ls.End
}

func (List) NodeType() NodeType {
	return listExprType
}

func (ls List) Pos() (start, end int) {
	return ls.Start, ls.End
}

func (AmbiguousTuple) NodeType() NodeType {
	return tupleType
}

func (t AmbiguousTuple) Pos() (start, end int) {
	return t.Start, t.End
}

func (Tuple) NodeType() NodeType {
	return pairsType
}

func (ps Tuple) Pos() (start, end int) {
	return ps.Start, ps.End
}
