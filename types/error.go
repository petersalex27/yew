package types

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/source"
)

const (
	CouldNotUnify                          string = "could not unify terms"
	CouldNotEquate                         string = "could not equate terms"
	KindMismatch                           string = "kind mismatch"
	UnknownIdentifier                      string = "unknown identifier"
	Redefinition                           string = "variable redefined"
	TooLowOrder0                           string = "expected a term belonging to at least 'Type'"
	TooLowOrder1                           string = "expected a term belonging to at least 'Type 1'"
	NonFunctionError                       string = "expected function"
	FailedGeneralization                   string = "failed to generalize type '%v' to '%v -> %v'"
	IllegalSort                            string = "illegal sort"
	IllegalBinder                          string = "illegal binder"
	ExpectedType                           string = "expected type"
	Untypable                              string = "constant is cannot be typed, it is a just tag"
	VarCannotReduceToWHNF                  string = "variable at expression's head prevents the reduction to weak head normal form"
	IllegalConstant                        string = "illegal constant" // TODO: can more be said?
	ExpectedTypeUniverse                   string = "expected type universe 'Type'"
	ExpectedVariable                       string = "expected variable"
	ExpectedProductType                    string = "expected product type"
	AlreadyDeclared                        string = "name already declared"
	RedefinedConstructor                   string = "constructor redefined"
	ExpectedProductTypeAfterSpecialization string = "expected product type after specialization"
	ExpectedImplicitTermInProduct          string = "expected implicit term in product type"
	ExpectedImplicitTerm                   string = "expected implicit term"
	ExpectedLambdaAfterSpecialization      string = "expected function term after specialization"
	ExpectedTypeOfTermToMatchProdTerm      string = "expected type of term to match corresponding type in product type"
	UnexpectedTyping                       string = "unexpected typing"

	Warn_ApplyFunctionNotInEnvironment string = "the builtin apply function is not in the environment; this is likely unintended"
)

func multiplicityPreventsUse(v Variable) string {
	return fmt.Sprintf("multiplicity of '%v' prevents it from being available in this context", v)
}

func unexpectedType(Ts ...Term) string {
	ss := make([]string, len(Ts))
	for i, T := range Ts {
		ss[i] = TypingString(T)
	}
	mid := " "
	if len(ss) > 1 {
		mid = "s" + mid
	}
	return UnexpectedTyping + mid + strings.Join(ss, ", ")
}

func expectedTypeOfTermToMatchProdTerm(T, U Type) string {
	return fmt.Sprintf("expected type of term '%v' to match corresponding type '%v' in product type", T, U)
}

func typeNConsNotLegal(Type_n Universe) string {
	return fmt.Sprintf("illegal explicit constructor for type universe 'Type %d'; explicit type constructors are only permitted for 'Type' (Type 0)", Type_n)
}

func failedGeneralizationError(T, U Type) string {
	return fmt.Sprintf(FailedGeneralization, T, U, T)
}

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

func (env *Environment) warning(msg string, elem strPos) {
	start, end := elem.Pos()
	l1, l2, c1, c2 := env.src.CalcLocationRange(start, end)
	e := makeTypeWarning(msg, env.src.Path, l1, l2, c1, c2)
	env.messages = append(env.messages, e)
}

func illegalConstructor(Z Constant, T Type) string {
	return fmt.Sprintf("illegal constructor type '%v' for type '%v'", T, Z)
}

func invalidType(T Type) string {
	return fmt.Sprintf("invalid type '%v' for term", T)
}

func (env *Environment) error(msg string, elem strPos) {
	start, end := elem.Pos()
	l1, l2, c1, c2 := env.src.CalcLocationRange(start, end)
	e := makeTypeError(msg, env.src.Path, l1, l2, c1, c2)
	env.messages = append(env.messages, e)
}

func (env *Environment) mismatchUnifyingError(a, b Term) {
	start, end := calcStartEnd(a, b)
	env.unifyingError(a, b, start, end)
}

func (env *Environment) impossibleUnificationBcOfLength(a, b Term, aTerms, bTerms []Term) {
	lenA, lenB := len(aTerms), len(bTerms)
	smaller := lenA
	terms := bTerms
	term := b
	format := "nothing to unify with '%v' in '%v'"
	if lenA > lenB {
		smaller = lenB
		terms = aTerms
		term = a
	}
	msg := fmt.Sprintf(format, joinStringed(terms[smaller:], " "), term)
	start, end := calcStartEnd(a, b)
	l1, l2, c1, c2 := env.src.CalcLocationRange(start, end)
	e := makeTypeError(msg, env.src.Path, l1, l2, c1, c2)
	env.messages = append(env.messages, e)
}

func (env *Environment) unifyingError(a, b Term, start, end int) {
	// TODO: use actual start and end, use path
	msg := CouldNotUnify + " '" + a.String() + " = " + b.String() + "'"
	l1, l2, c1, c2 := env.src.CalcLocationRange(start, end)
	e := makeTypeError(msg, env.src.Path, l1, l2, c1, c2)
	env.messages = append(env.messages, e)
}

func (env *Environment) equivalenceError(a, b Term, start, end int) {
	// TODO: use actual start and end, use path
	msg := CouldNotEquate + " '" + a.String() + " = " + b.String() + "'"
	l1, l2, c1, c2 := env.src.CalcLocationRange(start, end)
	e := makeTypeError(msg, env.src.Path, l1, l2, c1, c2)
	env.messages = append(env.messages, e)
}

// reports an error that no rule exists from `s ~> t` where `A : s` and `B : t`
func (env *Environment) ruleError(A Term, s Term, B Term, t Term) {
	msg := fmt.Sprintf("no rule exists '%v ~> %v' where '%v : %v' and '%v : %v'", s, t, A, s, B, t)
	start, end := calcStartEnd(A, B)
	l1, l2, c1, c2 := env.src.CalcLocationRange(start, end)
	e := makeTypeError(msg, env.src.Path, l1, l2, c1, c2)
	env.messages = append(env.messages, e)
}

// reports an error that no rule exists from `s ~> t` where `A : s` and `B : t`
func (env *Environment) ruleError2(s Term, t Term) {
	msg := fmt.Sprintf("no rule exists '%v ~> %v'", s, t)
	start, end := calcStartEnd(s, t)
	l1, l2, c1, c2 := env.src.CalcLocationRange(start, end)
	e := makeTypeError(msg, env.src.Path, l1, l2, c1, c2)
	env.messages = append(env.messages, e)
}

func (env *Environment) unknownNameError(x fmt.Stringer, start, end int) {
	// TODO: use actual start and end, use path
	res := UnknownIdentifier + " '" + x.String() + "'"
	e := makeTypeError(res, env.src.Path, 0, 0, 0, 0)
	env.messages = append(env.messages, e)
}

func (env *Environment) orderPrecedesType0Error(a Term, start, end int) {
	res := TooLowOrder0 + "; " + a.String() + " : T, T <: Type"
	l1, l2, c1, c2 := env.src.CalcLocationRange(a.Pos())
	e := makeTypeError(res, env.src.Path, l1, l2, c1, c2)
	env.messages = append(env.messages, e)
}
