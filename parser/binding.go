package parser

// binding (term = term)

func (binding Binding) Pos() (start, end int) {
	return binding.Start, binding.End
}

func (Binding) NodeType() NodeType {
	return bindingType
}
