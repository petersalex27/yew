package kind

import (
	"io"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/stringer"
)

type (
	Type interface {
		// opposite of `Apply` method
		//
		// simply call `Apply` when allowed in most cases
		//
		// return nil if the type cannot be applied, otherwise return the type after the application
		Pipe(into Type) Type

		// usually called from the arg calling `a.Pipe(r)` to let the receiver `r` know it's allowed to apply `a` to itself
		Apply(a Type) App

		// return the type as a non-empty string IFF it is a constant--do not return variable names!
		Constant() string

		// break off the name from the type (if one exists), returning it and whatever remains
		//
		// reasonable examples:
		//	- variable is stored as just a name: return the name and nil
		//	- constant is stored as a string: return the string and nil (Constant() will differentiate this from a variable)
		//	- function is stored as a pair: return "->" and a list of two elements [lhs, rhs]
		//	- application is stored as a pair: return "" and a list of two elements [lsh, rhs]
		Break() (head string, tail []Type)

		String() string

		Typeable

		api.Positioned
	}

	Typeable interface {
		AsType(Solver) Type
	}

	Term interface {
		Abstract(x Variable) (Lambda, Type)
		Kind(s Solver) Type
		Apply(a Term) Term
	}

	Solver interface {
		Find(Typeable) Type
		Occurs(string, Type) bool
		Union(string, Type) bool
		NewVar() Type
		RecordError(errorMessage string, pos api.Positioned)
		Declare(api.Token, byte, byte)
		RecordNames()
		EndRecording() map[string]Symbol
		Print(w io.Writer)
	}

	Variable struct {
		x string
		api.Position
	}

	Lambda struct {
		domain Symbol
		target Term
	}

	Pi struct {
	}
)

type App struct {
	ts []Type
	api.Position
}

func Apply(a Type, b Type) App {
	// DO NOT just call `a.Apply(b)` here--this could cause an infinite loop
	return App{[]Type{a, b}, api.WeakenRangeOver(a, b)}
}

func FlatApply(as []Type) App {
	if len(as) < 2 {
		panic("illegal call to FlatApply: must provide at least two types")
	}
	return App{as, api.WeakenRangeOver(as[0], as[1:]...)}
}

func (as App) Pipe(into Type) Type {
	return into.Apply(as)
}

func (as App) AsType(Solver) Type { return as }

func (as App) Apply(b Type) App {
	as.ts = append(as.ts, b)
	as.Position = as.Position.Update(b)
	return as
}

func (as App) Constant() string {
	return ""
}

func (as App) Break() (head string, tail []Type) {
	if len(as.ts) == 0 {
		return "", []Type{}
	}
	if as.ts[0].Constant() != "" {
		return as.ts[0].Constant(), as.ts[1:]
	}

	return "", as.ts
}

func (as App) String() string {
	_, t := as.Break()
	return stringer.Join(t, " ")
}
