package symbol

import "github.com/petersalex27/yew/api"

type (
	variable string

	Pi struct {
		domain sym
		target api.Type
	}

	Constant struct {
		// X part of X.c version of a constant's name, used to disambiguate constants with the same
		// name in different namespaces
		namespace string
		// c part of X.c version of a constant's name, the standard way of referring to it
		name string
	}
)

func (s sym) Pipe(into api.Type) api.Type {
	return into.Apply(s)
}

func (s sym) Apply(a api.Type) api.App {
	return api.App{s, a}
}

// if the symbol's kind is a constant, return the name of the constant
func (s sym) Constant() string {
	return s.typ.Constant()
}

func (s sym) Break() (head string, tail []api.Type) {
	return s.typ.Break()
}

func (v variable) Pipe(into api.Type) api.Type {
	return into.Apply(v)
}

func (r variable) Apply(a api.Type) api.App {
	return api.App{r, a}
}

// always returns ""--free variables are not constants
func (v variable) Constant() string {
	return "" // free variables are not constants
}

// returns the name of the variable and an empty list
func (v variable) Break() (head string, tail []api.Type) {
	return string(v), []api.Type{}
}

func (p Pi) Pipe(into api.Type) api.Type {
	return into.Apply(p)
}

func (p Pi) Apply(a api.Type) api.App {
	return api.App{p, a}
}

// always returns ""--pi types are not constants
func (p Pi) Constant() string {
	return ""
}

// returns the domain and target of the pi type and `->` as the head
func (p Pi) Break() (head string, tail []api.Type) {
	return "->", []api.Type{p.domain, p.target}
}

func (c Constant) Pipe(into api.Type) api.Type {
	return into.Apply(c)
}

func (c Constant) Apply(a api.Type) api.App {
	return api.App{c, a}
}

// returns the name of the constant (with the namespace prepended)
func (c Constant) Constant() string {
	return c.namespace + "." + c.name
}

// returns the qualified name of the constant and an empty list
func (c Constant) Break() (head string, tail []api.Type) {
	return c.Constant(), []api.Type{}
}