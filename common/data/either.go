package data

import "github.com/petersalex27/yew/api"

type (
	Either[a, b api.Node] interface {
		api.DescribableNode
		Break() (left a, right b, isRight bool)
		IsLeft() bool
		either()
		Update(api.Positioned) Either[a, b]
		Children() []api.Node
	}

	EmbedsEither[a, b api.Node] interface{
		api.DescribableNode
		~struct{ Either[a, b] }
	}

	inLeft[a, b api.Node] struct {
		val a
		api.Position
	}

	inRight[a, b api.Node] struct {
		val b
		api.Position
	}
)

func (e inLeft[a, b]) Update(p api.Positioned) Either[a, b] {	
	e.Position = e.Position.Update(p)
	return e
}

func (e inRight[a, b]) Update(p api.Positioned) Either[a, b]{
	e.Position = e.Position.Update(p)
	return e
}

func IsLeft[a, b api.Node](e Either[a, b]) bool { return e.IsLeft() }

func IsRight[a, b api.Node](e Either[a, b]) bool { return !e.IsLeft() }

func __PassErs[b api.Node](e Ers) Either[Ers, b] { return Inl[b](e) }

func Inl[b, a api.Node](x a) Either[a, b] { return inLeft[a, b]{x, x.GetPos()} }

func Inr[a, b api.Node](x b) Either[a, b] { return inRight[a, b]{x, x.GetPos()} }

func (inLeft[a, b]) IsLeft() bool { return true }

func (lhs inLeft[a, b]) Break() (left a, right b, isRight bool) {
	left = lhs.val
	return left, right, false
}

func (inLeft[a, b]) either() {}

func (inRight[a, b]) IsLeft() bool { return false }

func (rhs inRight[a, b]) Break() (left a, right b, isRight bool) {
	right = rhs.val
	return left, right, true
}

func (inRight[a, b]) either() {}

// result in receiver is left alone, new result is returned
//
// if left is not empty, it is appended to with the error message and position
// so, it may modify the receiver's copy as well
func __Fail[a api.Node](msg string, positioned api.Positioned) Either[Ers, a] {
	es := Nil[Err](1).Snoc(MkErr(msg, positioned))
	return Inl[a](es)
}

func Ok[a api.Node](x a) Either[Ers, a] { return Inr[Ers](x) }

func SeqResult[a api.Node](msg string) func(ma Maybe[a]) Either[Ers, a] {
	return func(ma Maybe[a]) Either[Ers, a] {
		unit, isJust := ma.Break()
		if !isJust{
			return Fail[a](msg, ma)
		}
		return Ok(unit)
	}
}

func EInl[c EmbedsEither[a, b], a, b api.Node](x a) c {
	return c{Inl[b](x)}
}

func EInr[c EmbedsEither[a, b], a, b api.Node](x b) c {
	return c{Inr[a](x)}
}

func Cases[a, b api.Node, c any](eab Either[a, b], f func(a) c, g func(b) c) c {
	lhs, rhs, isRhs := eab.Break()
	if isRhs {
		return g(rhs)
	}
	return f(lhs)
}


func EitherMap[a, b, d, c api.Node](f func(a) c, g func(b) d) func(Either[a, b]) Either[c, d] {
	return func(eab Either[a, b]) Either[c, d] {
		if lhs, rhs, isRhs := eab.Break(); isRhs {
			return Inr[c](g(rhs))
		} else {
			return Inl[d](f(lhs))
		}
	}
}