package data

import "github.com/petersalex27/yew/api"

type (
	Maybe[a api.Node] interface {
		api.DescribableNode
		IsNothing() bool
		Break() (unit a, isJust bool)
		Children() []api.Node
		Update(p api.Positioned) Maybe[a]
		maybe()
	}

	nothing[a api.Node] struct {
		api.Position
	}

	just[a api.Node] struct {
		unit a
		api.Position
	}

	EmbedsMaybe[a api.Node] interface {
		api.DescribableNode
		~struct{ Maybe[a] }
	}
)

func IsNothing[a api.Node](m Maybe[a]) bool { return m.IsNothing() }

// Positive constructor for Maybe
func Just[a api.Node](x a) Maybe[a] {
	return just[a]{unit: x, Position: x.GetPos()}
}

// Empty constructor for Maybe
//
// position is optional, but may also be multiple positions; in the latter case, the non-zero range
// over the positions is used with a bias towards positive-only positions
func Nothing[a api.Node](pos ...api.Positioned) Maybe[a] {
	if len(pos) > 0 {
		pos := api.WeakenRangeOver(pos[0], pos[1:]...)
		return nothing[a]{Position: pos}
	}
	return nothing[a]{Position: api.ZeroPosition()}
}

func (n nothing[a]) maybe() {}

func (n nothing[a]) Update(p api.Positioned) Maybe[a] {
	n.Position = api.WeakenRangeOver(p.GetPos(), n.Position)
	return n
}

func (n nothing[a]) IsNothing() bool { return true }

func (n nothing[a]) Break() (unit a, isJust bool) { return }

func (j just[a]) maybe() {}

func (j just[a]) Update(p api.Positioned) Maybe[a] {
	j.Position = api.WeakenRangeOver(p.GetPos(), j.Position)
	return j
}

func (j just[a]) IsNothing() bool { return false }

func (j just[a]) Break() (unit a, isJust bool) {
	return j.unit, true
}

func Bind[b, a api.Node](m Maybe[a], f func(a) Maybe[b]) Maybe[b] {
	if x, just := m.Break(); !just {
		return Nothing[b](m)
	} else {
		return f(x)
	}
}

func MapMaybe[a, b api.Node](f func(a) b) func(Maybe[a]) Maybe[b] {
	return func(m Maybe[a]) Maybe[b] {
		if x, just := m.Break(); just {
			return Just(f(x))
		}
		return Nothing[b](m)
	}
}
