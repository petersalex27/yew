package types

import (
	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/source"
)

const (
	CouldNotUnify     string = "could not unify terms"
	CouldNotEquate    string = "could not equate terms"
	KindMismatch      string = "kind mismatch"
	UnknownIdentifier string = "unknown identifier"
	TooLowOrder0      string = "expected a term belonging to at least '*'"
	TooLowOrder1      string = "expected a term belonging to at least '*1'"
)

// creates a syntax error from the arguments
func makeTypeError(msg string, path source.PathSpec, line, lineEnd, start, end int) errors.ErrorMessage {
	e := errors.MakeError("Type", msg, line, lineEnd, start, end)
	if path == nil {
		e.SourceName = "unknown"
	} else {
		e.SourceName = path.Path()
	}
	return e
}

// creates a syntax error from the arguments
func makeTypeWarning(msg string, path source.PathSpec, line, lineEnd, start, end int) errors.ErrorMessage {
	e := errors.MakeWarning("Type", msg, line, lineEnd, start, end)
	if path == nil {
		e.SourceName = "unknown"
	} else {
		e.SourceName = path.Path()
	}
	return e
}

func (env *Environment) error(msg string) {
	e := makeTypeError(msg, source.StdinSpec, 0, 0, 0, 0) // TODO: path and loc info
	env.messages = append(env.messages, e)
}

func (env *Environment) unifyingError(a, b Term, start, end int) {
	// TODO: use actual start and end, use path
	res := CouldNotUnify + " '" + a.String() + " = " + b.String() + "'"
	e := makeTypeError(res, source.StdinSpec, 0, 0, 0, 0)
	env.messages = append(env.messages, e)
}

func (env *Environment) equivalenceError(a, b Term, start, end int) {
	// TODO: use actual start and end, use path
	res := CouldNotEquate + " '" + a.String() + " = " + b.String() + "'"
	e := makeTypeError(res, source.StdinSpec, 0, 0, 0, 0)
	env.messages = append(env.messages, e)
}

func (env *Environment) unknownNameError(x Identifier, start, end int) {
	// TODO: use actual start and end, use path
	res := UnknownIdentifier + " '" + x.String() + "'"
	e := makeTypeError(res, source.StdinSpec, 0, 0, 0, 0)
	env.messages = append(env.messages, e)
}

func (env *Environment) orderPrecedesType0Error(a Term, start, end int) {
	// TODO: use actual start and end, use path
	res := TooLowOrder0 + "; " + a.String() + " : T, T <: *"
	e := makeTypeError(res, source.StdinSpec, 0, 0, 0, 0)
	env.messages = append(env.messages, e)
}
