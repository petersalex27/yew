package data

import (
	"github.com/petersalex27/yew/api"
)

type (
	List[a api.Node] struct {
		elements []a
		api.Position
	}

	EmbedsList[a api.Node] interface {
		api.DescribableNode
		~struct{ List[a] }
	}
)

func (xs List[a]) Elements() []a {
	return xs.elements
}

func (List[a]) zeroElement() (_ a) { return }

// try to strengthen a list to a lifted non-empty list
func EStrengthen[ne EmbedsNonEmpty[a], a api.Node](xs List[a]) Maybe[ne] {
	return Bind(xs.Strengthen(), justLiftNonEmpty[ne])
}

func (xs List[a]) Strengthen() Maybe[NonEmpty[a]] {
	if xs.IsEmpty() {
		return Nothing[NonEmpty[a]](xs.Position)
	}
	if len(xs.elements) > 1 {
		rest := List[a]{elements: xs.elements[1:]}
		return Just(NonEmpty[a]{first: xs.elements[0], rest: rest, Position: xs.Position})
	}
	return Just(NonEmpty[a]{first: xs.elements[0], rest: Nil[a](), Position: xs.Position})
}

// only first element in cap is used if more than 0 provided
func Nil[a api.Node](cap ...int) List[a] {
	var effectiveCap int = 0
	if len(cap) > 0 {
		effectiveCap = cap[0]
	}
	return List[a]{elements: make([]a, 0, effectiveCap)}
}

func Makes[a api.Node](xs ...a) List[a] {
	return Nil[a](len(xs)).Append(xs...)
}

func EMakes[list EmbedsList[a], a api.Node](xs ...a) list {
	return list{Makes(xs...)}
}

func (xs List[a]) IsEmpty() bool {
	return len(xs.elements) == 0
}

func (xs List[a]) Len() int {
	return len(xs.elements)
}

func (xs List[a]) Snoc(e a) List[a] {
	if xs.elements == nil {
		xs.elements = make([]a, 0, 1)
	}
	xs.elements = append(xs.elements, e)
	xs.Position = xs.Update(e)
	return xs
}

func (xs List[a]) Append(es ...a) List[a] {
	xs.elements = append(xs.elements, es...)
	if len(es) < 1 {
		xs.Position = api.ZeroPosition()
	} else {
		xs.Position = api.WeakenRangeOver(es[0], es[1:]...)
	}
	return xs
}

func ListMap[a, b api.Node](f func(a) b) func(List[a]) List[b] {
	return func(xs List[a]) List[b] {
		ys := Nil[b](len(xs.elements))
		for _, x := range xs.elements {
			ys = ys.Snoc(f(x))
		}
		return ys
	}
}

func appendAll[a api.Node](xs ...a) List[a] {
	ys := List[a]{elements: xs}
	if len(xs) == 0 {
		ys.Position = api.ZeroPosition()
	} else if len(xs) == 1 {
		ys.Position = api.WeakenRangeOver(xs[0], xs[0])
	} else {
		ys.Position = api.WeakenRangeOver(xs[0], xs[1:]...)
	}
	return ys
}

func (xs List[a]) Head() Maybe[a] {
	if len(xs.elements) == 0 {
		return Nothing[a](xs)
	}
	return Just(xs.elements[0])
}

func (xs List[a]) Tail() Maybe[List[a]] {
	if len(xs.elements) == 1 {
		return Just(Nil[a]())
	} else if len(xs.elements) == 0 {
		return Nothing[List[a]](xs)
	}
	return Just(appendAll(xs.elements[1:]...))
}

// applies the folding function `f` from the end of the list to the beginning
//
// Example:
// 	```
//	FoldRight(subtract, 0)(Makes(1, 2, 3) 
//		= subtract(1, subtract(2, subtract(3, 0)))
//		= subtract(1, subtract(2, 3))
//		= subtract(1, -1)
//		= 2
//	```
//
// SEE: FoldLeft for a fold in the opposite direction
func FoldRight[a, b api.Node](f func(a, b) b, z b) func(List[a]) b {
	return func(xs List[a]) b {
		for i := xs.Len() - 1; i >= 0; i-- {
			z = f(xs.elements[i], z)
		}
		return z
	}
}

// applies the folding function `f` from the beginning of the list to the end
//
// Example:
// 	```
//	FoldLeft(subtract, 0)(Makes(1, 2, 3)
//		= subtract(subtract(subtract(0, 1), 2), 3)
//		= subtract(subtract(-1, 2), 3)
//		= subtract(-3, 3)
//		= -6
//	```
//
// SEE: FoldRight for a fold in the opposite direction
func FoldLeft[a, b api.Node](f func(b, a) b, z b) func(List[a]) b {
	return func(xs List[a]) b {
		for _, x := range xs.elements {
			z = f(z, x)
		}
		return z
	}
}