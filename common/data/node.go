package data

import "github.com/petersalex27/yew/api"

func (o Solo[a]) Describe() (string, []api.Node) {
	if d, ok := api.Node(o.one).(api.DescribableNode); ok {
		return d.Describe()
	}
	return "one", []api.Node{o.one}
}

func (o Solo[T]) Children() []api.Node {
	return []api.Node{o.one}
}

func (p Pair[T, U]) Describe() (string, []api.Node) {
	return p.Type().String(), p.Children()
}

func (p Pair[T, U]) Children() []api.Node {
	return []api.Node{p.first, p.second}
}

func (n NonEmpty[a]) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}

func (n NonEmpty[T]) Children() []api.Node {
	if len(n.rest.elements) == 0 {
		return []api.Node{n.first}
	}
	return append([]api.Node{n.first}, n.rest.Children()...)
}

func (l List[T]) Describe() (string, []api.Node) {
	return l.Type().String(), l.Children()
}

func (l List[T]) Children() []api.Node {
	out := make([]api.Node, len(l.elements))
	for i, el := range l.elements {
		out[i] = el
	}
	return out
}

func (m just[a]) Describe() (string, []api.Node) {
	if d, ok := api.Node(m.unit).(api.DescribableNode); ok {
		return d.Describe()
	}
	return "just", m.Children()
}

func (m nothing[a]) Describe() (string, []api.Node) {
	return "empty", m.Children()
}

func (m just[a]) Children() []api.Node {
	return []api.Node{m.unit}
}

func (m nothing[a]) Children() []api.Node {
	return []api.Node{}
}

func (lhs inLeft[a, b]) Describe() (string, []api.Node) {
	if d, ok := api.Node(lhs.val).(api.DescribableNode); ok {
		return d.Describe()
	}
	return "in-left", []api.Node{lhs.val}
}

func (rhs inRight[a, b]) Describe() (string, []api.Node) {
	if d, ok := api.Node(rhs.val).(api.DescribableNode); ok {
		return d.Describe()
	}
	return "in-right", []api.Node{rhs.val}
}

func (lhs inLeft[a, b]) Children() []api.Node {
	return []api.Node{lhs.val}
}

func (rhs inRight[a, b]) Children() []api.Node {
	return []api.Node{rhs.val}
}