package parser

func (a app[T]) apply(t T) app[T] {
	return app[T]{elems: append(a.elems, t), position: a.reposition(t)}
}