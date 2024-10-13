package api

type Type interface {
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
}

type App []Type

func (as App) Pipe(into Type) Type {
	return into.Apply(as)
}

func (as App) Apply(b Type) App {
	return append(as, b)
}

func (as App) Constant() string {
	return ""
}

func (as App) Break() (head string, tail []Type) {
	if len(as) == 0 {
		return "", []Type{}
	}
	if as[0].Constant() != "" {
		return as[0].Constant(), as[1:]
	}

	return "", as
}