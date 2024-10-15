package data

import "github.com/petersalex27/yew/api"

type (
	NonEmpty[a api.Node] struct {
		first a
		rest  List[a]
		api.Position
	}

	EmbedsNonEmpty[a api.Node] interface {
		api.DescribableNode
		~struct{ NonEmpty[a] }
	}
)

func (ne NonEmpty[a]) Append(es ...a) NonEmpty[a] {
	ne.rest = ne.rest.Append(es...)
	ne.Position = ne.Update(ne.rest)
	return ne
}

// uses the fact that the constructors for NonEmpty[a] all allocate memory for the rest.elements slice
func (ne NonEmpty[a]) Un_constructed() bool {
	return ne.rest.elements == nil
}

func (ne NonEmpty[a]) Snoc(e a) NonEmpty[a] {
	if ne.rest.elements == nil {
		panic("illegal receiver: NonEmpty[a] 'ne' is, um, empty ...")
	}
	ne.Position = ne.Update(e)
	ne.rest = ne.rest.Snoc(e)
	return ne
}

func Singleton[a api.Node](first a) NonEmpty[a] {
	ne := NonEmpty[a]{first: first, rest: Nil[a]()}
	ne.Position = ne.Update(first)
	return ne
}

func MakeApp[app EmbedsPair[a, NonEmpty[a]], a api.Node](x a, y a, zs ...a) app {
	return EMakePair[app](x, Construct(y, zs...))
}

func NonEmptyToAppLikePair[app EmbedsPair[a, NonEmpty[a]], a api.Node](ne NonEmpty[a]) (out Maybe[app]) {
	if ne.Len() < 2 {
		// need at least two elements
		return Nothing[app](ne)
	}
	appValue := MakeApp[app](ne.first, ne.rest.elements[0], ne.rest.elements[1:]...)
	return Just(appValue)
}

func EConstruct[ene EmbedsNonEmpty[a], a api.Node](first a, rest ...a) ene {
	return ene{Construct(first, rest...)}
}

func Construct[a api.Node](first a, rest ...a) NonEmpty[a] {
	var ne NonEmpty[a]
	ne.first = first
	ne.rest = Nil[a](len(rest))
	ne.Position = ne.Update(first)
	ne.rest = ne.rest.Append(rest...)
	return ne
}

func (xs NonEmpty[a]) Map(f func(a) a) NonEmpty[a] {
	xs.first = f(xs.first)
	xs.rest = xs.rest.Map(f)
	return xs
}

func MapNonEmpty[a, b api.Node](f func(a) b) func(xs NonEmpty[a]) NonEmpty[b] {
	listMap := MapList(f)
	return func(xs NonEmpty[a]) NonEmpty[b] {
		return NonEmpty[b]{first: f(xs.first), rest: listMap(xs.rest)}
	}
}

func (xs NonEmpty[a]) Head() a { return xs.first }

func (xs NonEmpty[a]) Tail() List[a] { return xs.rest }

func (xs NonEmpty[a]) Len() int { return 1 + xs.rest.Len() }

func (xs NonEmpty[a]) Elements() []a {
	return append([]a{xs.first}, xs.rest.elements...)
}