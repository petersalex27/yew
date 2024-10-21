package symbol

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/internal/kind"
)

func Pi(domain kind.Symbol, target kind.Type) pi {
	return pi{
		domain,
		target,
		api.WeakenRangeOver[api.Positioned](domain, target),
	}
}

func Implicit(domain kind.Symbol, target kind.Type) implicitPi {
	return implicitPi{Pi(domain, target)}
}

func Const(name api.Token) constant {
	return constant{name.String(), name.GetPos()}
}

func Var(name api.Token) variable {
	return variable{name.String(), name.GetPos()}
}

func Var2(name string, pos api.Position) variable {
	return variable{name, pos}
}