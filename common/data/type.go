package data

import "github.com/petersalex27/yew/api"

type dataType string

func (dt dataType) String() string {
	return string(dt)
}

func (dt dataType) Match(n api.Node) bool {
	return api.NodeTypeString(n) == dt.String()
}

func (xs List[a]) Type() api.NodeType {
	if xs.Len() == 0 {
		return dataType("empty list")
	}
	elemTypeString := api.NodeTypeString(xs.zeroElement())
	return dataType(elemTypeString + " list")
}

func (e inLeft[a, b]) Type() api.NodeType {
	return e.val.Type()
}

func (e inRight[a, b]) Type() api.NodeType {
	return e.val.Type()
}

func (e Err) Type() api.NodeType {
	return dataType("error")
}

func (n nothing[a]) Type() api.NodeType {
	return dataType("empty")
}

func (j just[a]) Type() api.NodeType {
	return dataType(api.NodeTypeString(j.unit))
}

func (ne NonEmpty[a]) Type() api.NodeType {
	s := api.NodeTypeString(ne.first)
	return dataType("non-empty " + s + " list")
}

func (s Solo[a]) Type() api.NodeType {
	return dataType(api.NodeTypeString(s.one))
}

func (p Pair[a, b]) Type() api.NodeType {
	fst := api.NodeTypeString(p.first)
	snd := api.NodeTypeString(p.second)
	return dataType(fst + " and " + snd + " pair")
}
