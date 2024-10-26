package symbol

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/internal/kind"
)

type (
	variable struct {
		x string
		api.Position
	}

	pi struct {
		domain kind.Symbol
		target kind.Type
		api.Position
	}

	implicitPi struct {
		pi
	}

	constant struct {
		name string
		api.Position
	}
)

func (v variable) AsType(kind.Solver) kind.Type { return v }

// always returns ""--free variables are not constants
func (v variable) Constant() string {
	return "" // free variables are not constants
}

// returns the name of the variable and an empty list
func (v variable) Break() (head string, tail []kind.Type) {
	return v.x, []kind.Type{}
}

func (v variable) String() string {
	return v.x
}

func (p pi) AsType(kind.Solver) kind.Type { return p }

// always returns ""--pi types are not constants
func (p pi) Constant() string {
	return ""
}

// returns the domain and target of the pi type and `->` as the head
func (p pi) Break() (head string, tail []kind.Type) {
	return "->", []kind.Type{p.domain.Typing(), p.target}
}

func (p pi) String() string {
	domain := p.domain.String()
	target := p.target.String()
	if _, domPi := p.domain.Typing().(pi); domPi {
		domain = "(" + domain + ")"
	}

	return domain + " -> " + target
}

func (c constant) AsType(kind.Solver) kind.Type { return c }

// returns the name of the constant (with the namespace prepended)
func (c constant) Constant() string {
	return c.name
}

// returns the qualified name of the constant and an empty list
func (c constant) Break() (head string, tail []kind.Type) {
	return c.Constant(), []kind.Type{}
}

func (c constant) String() string {
	return c.name
}
