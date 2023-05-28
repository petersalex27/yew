package types

import (
	"strings"
	"sync"
	"yew/demangler"
	err "yew/error"
	util "yew/utils"
)

type Superable interface {
	GetSuperType() Types
}

type inferenceRules struct {
	inferenceLock sync.Mutex
	rules map[string]*[]Types
}

var inferences = &inferenceRules{ 
	rules: make(map[string]*[]Types),
}

func (inf *inferenceRules) recursiveCheckForCycles(tau string, marked *map[string]bool, visited *map[string]bool) bool {
	if (*marked)[tau] {
		return true
	}

	if (*visited)[tau] {
		return false
	}

	(*marked)[tau] = true
	(*visited)[tau] = true
	for _, ty := range *inf.rules[tau] {
		if ty.GetTypeType() == TAU {
			if inf.recursiveCheckForCycles(string(ty.(Tau)), marked, visited) {
				return true
			}
		}
	} 

	(*marked)[tau] = false
	
	return false
}

// DO NOT CALL UNLESS LOCK IS HELD!!!
func (inf *inferenceRules) checkForCycles() string {
	// stack := make([]string, 0, len(inf.rules) / 2) // `len / 2` is somewhat arbitrary
	marked := make(map[string]bool, len(inf.rules))
	visited := make(map[string]bool, len(inf.rules))
	
	for tau := range inf.rules {
		if inf.recursiveCheckForCycles(tau, &marked, &visited) {
			return tau
		}
	}

	return "" // return empty string, i.e., empty cycle
}

// DO NOT CALL UNLESS LOCK IS HELD!!
func (inf *inferenceRules) tryToStrengthen(tau Tau) (ty Types, stronger bool) {
	if _, found := inf.rules[string(tau)]; !found {
		return tau, false
	}

	var prev string = string(tau)
	var next *[]Types
	var found bool
	stack := make([]string, 0, len(inf.rules) / 2)
	searched := make(map[string]bool, len(inf.rules) / 2)
	stack = append(stack, string(tau))
	for ; len(stack) > 0; {
		prev = stack[len(stack) - 1]
		searched[prev] = true
		stack = stack[:len(stack) - 1]
		next, found = inf.rules[prev]
		if !found {
			continue
		}
		for _, t := range *next {
			if t.GetTypeType() == TAU {
				x := string(t.(Tau))
				if !searched[x] {
					stack = append(stack, string(t.(Tau)))
				}
			} else {
				return t, true
			}
		}
	}

	return tau, false
} 

// false on error
func (inf *inferenceRules) addRule(tau Tau, to Types) (res Types, ok bool) {
	inf.inferenceLock.Lock()
	// check if already has stronger rule
	var stronger bool 
	res, stronger = inf.tryToStrengthen(tau)
	if stronger {
		if to.GetTypeType() == TAU {
			// add rule for type that was attempted to be added
			res, ok = inf.addRule(to.(Tau), res)
		} else {
			// check if inferences match
			ok = to.Equals(res)
		}
	} else {
		// rule must be added
		if nil == inf.rules[string(tau)] {
			inf.rules[string(tau)] = new([]Types)
			*inf.rules[string(tau)] = make([]Types, 1)
			(*inf.rules[string(tau)])[0] = to
			res, ok = to, true
		} else {
			// check if contained in rule already
			var exit bool = false
			for _, t := range (*inf.rules[string(tau)]) {
				if t.Equals(to) {
					exit = true 
					break
				}
			}
			if !exit {
				// is not already in rules, so append
				(*inf.rules[string(tau)]) = append((*inf.rules[string(tau)]), to)
				res, ok = to, ok
			}
		}
	}
	inf.inferenceLock.Unlock()
	return
}

func Strengthen(dest *Types, src Types) Types {
	if TAU == (*dest).GetTypeType() {
		*dest = src
	}
	return *dest
}

type TypeType int

func GetNewTau() Tau {
	return Tau(demangler.TYPE.GetDemanglerPrefix())
} 

func GetNewTaus(n int) (taus []Tau) {
	taus = make([]Tau, n)
	demangles := demangler.TYPE.GetDemanglerPrefixes(n)
	for i, d := range demangles {
		taus[i] = Tau(d)
	}
	return
}

const (
	INT TypeType = iota
	BOOL
	CHAR
	FLOAT
	ARRAY
	FUNCTION
	TUPLE
	TAU
	QUALIFIER
	CLASS 
	DICTIONARY
	APPLICATION
	ERROR
	VALUE
)

type Types interface {
	util.Stringable
	GetTypeType() TypeType
	// checks for equality
	Equals(Types) bool
	// (λT.e) U [T in e ::= U]
	Apply(Types) Types
	// receiver == tau -> receiver ::= t
	//  - inferences also get updated
	ReplaceTau(tau Tau, t Types) Types
	// receiver ::= from
	InferType(from Types) Types
}

func CanBeReplaced(t Types, replacement Types) bool {
	if t.GetTypeType() == VALUE {
		t = t.(Superable).GetSuperType()
	}
	if replacement.GetTypeType() == VALUE {
		replacement = replacement.(Superable).GetSuperType()
	}

	return t.Equals(replacement)
}

type Int struct {}
type Float struct {}
type Char struct {}
type Bool struct {}
type Tau string
type Function struct {
	Domain Types
	Codomain Types
}
type Array struct {ElemType Types}
type Tuple []Types
type NamedTuple map[string]Types
type Class struct {
	Name string
	TypeVariable Tau
	Functions map[string]Function
}
type Qualifier struct {
	Class Class
	TypeVariable Tau
	Qualified Types
}
type Error err.Error
func (e Error) ToError() err.Error {
	return err.Error(e)
}

// represents a type that has not yet been evaluated; by the end of the resolve type phase,
// there should be no application types left
type Application []Types

func (i Int) ToString() string {
	return "Int"
}
func (b Bool) ToString() string {
	return "Bool"
}
func (c Char) ToString() string {
	return "Char"
}
func (f Float) ToString() string {
	return "Float"
}
func (tau Tau) ToString() string {
	return string(tau)
}
func (f Function) ToString() string {
	var left string
	var right string 

	if FUNCTION == f.Domain.GetTypeType() {
		left = "(" + f.Domain.ToString() + ")"
	} else {
		left = f.Domain.ToString()
	}

	if FUNCTION == f.Codomain.GetTypeType() {
		right = "(" + f.Codomain.ToString() + ")"
	} else {
		right = f.Codomain.ToString()
	}

	return  left + " -> " + right 
}
func (a Array) ToString() string {
	return "[" + a.ElemType.ToString() + "]"
}
func (tup Tuple) ToString() string {
	var ts []Types = tup[:]
	ss := util.Fmap(ts, func(t Types) string {
		return t.ToString()
	})
	return "(" + strings.Join(ss, ", ") + ")"
}
func (e Error) ToString() string {
	return err.Error(e).ToString()
}
func (q Qualifier) ToString() string {
	return q.Class.Name + " " + string(q.TypeVariable)
}
func (c Class) ToString() string {
	var tmp []string
	for k, f := range c.Functions {
		tmp = append(tmp, k + " " + f.ToString())
	}
	return c.Name + " " + c.TypeVariable.ToString() + "{" + strings.Join(tmp, "; ") + "}"
}
func (tup NamedTuple) ToString() string {
	var tmp []string
	for k, v := range tup {
		tmp = append(tmp, k + " " + v.ToString())
	}
	return "(" + strings.Join(tmp, ", ") + ")"
}
func (a Application) ToString() string {
	if len(a) == 0 {
		return ""
	}
	var builder strings.Builder
	for _, t := range a {
		builder.WriteString(t.ToString())
		builder.WriteByte(' ')
	}
	res := builder.String()

	return res[:len(res) - 1]
}

func (i Int) Apply(t Types) Types {
	return Application([]Types{i, t})
}
func (b Bool) Apply(t Types) Types {
	return Application([]Types{b, t})
}
func (c Char) Apply(t Types) Types {
	return Application([]Types{c, t})
}
func (f Float) Apply(t Types) Types {
	return Application([]Types{f, t})
}
func (tau Tau) Apply(t Types) Types {
	return Application([]Types{tau, t})
}
func (f Function) Apply(t Types) Types {
	if f.Domain.Equals(t) {
		return f.Codomain
	} else if TAU == f.Domain.GetTypeType() {
		return f.Codomain.ReplaceTau(f.Domain.(Tau), t)
	}
	return Application([]Types{f, t})
}
func (a Array) Apply(t Types) Types {
	return Application([]Types{a, t})
}
func (tup Tuple) Apply(t Types) Types {
	return Application([]Types{tup, t})
}
func (tup NamedTuple) Apply(t Types) Types {
	if TAU == t.GetTypeType() {
		ty, found := tup[t.ToString()]
		if found {
			return ty
		}
	}
	return Application([]Types{tup, t})
}
func (e Error) Apply(t Types) Types {
	return e
}
func (c Class) Apply(t Types) Types {
	var out NamedTuple = make(NamedTuple)

	for m, f := range c.Functions {
		out[m] = f.ReplaceTauInFunction(c.TypeVariable, t)
	}
	return out
}
func (q Qualifier) Apply(t Types) Types {
	return q.Qualified.ReplaceTau(q.TypeVariable, t)
}
func (a Application) Apply(t Types) Types {
	return Application([]Types{a, t})
}

func (i Int) GetTypeType() TypeType {
	return INT
}
func (b Bool) GetTypeType() TypeType {
	return BOOL
}
func (c Char) GetTypeType() TypeType {
	return CHAR
}
func (f Float) GetTypeType() TypeType {
	return FLOAT
}
func (tau Tau) GetTypeType() TypeType {
	return TAU
}
func (f Function) GetTypeType() TypeType {
	return FUNCTION
}
func (a Array) GetTypeType() TypeType {
	return ARRAY
}
func (tup Tuple) GetTypeType() TypeType {
	return TUPLE
}
func (e Error) GetTypeType() TypeType {
	return ERROR
}
func (q Qualifier) GetTypeType() TypeType {
	return QUALIFIER
}
func (tup NamedTuple) GetTypeType() TypeType {
	return DICTIONARY
}
func (c Class) GetTypeType() TypeType {
	return CLASS
}
func (a Application) GetTypeType() TypeType {
	return APPLICATION
}

func (i Int) Equals(t Types) bool {
	return INT == t.GetTypeType()
}
func (b Bool) Equals(t Types) bool {
	return BOOL == t.GetTypeType()
}
func (c Char) Equals(t Types) bool {
	return CHAR == t.GetTypeType()
}
func (f Float) Equals(t Types) bool {
	return FLOAT == t.GetTypeType()
}
func (tau Tau) Equals(t Types) bool {
	return TAU == t.GetTypeType() && string(tau) == string(t.(Tau))
}
func (f Function) Equals(t Types) bool {
	return FUNCTION == t.GetTypeType() && 
		f.Domain.Equals(t.(Function).Domain) &&
		f.Codomain.Equals(t.(Function).Codomain)
}
func (a Array) Equals(t Types) bool {
	return ARRAY == t.GetTypeType() && a.ElemType.Equals(t.(Array).ElemType)
}
func (tup Tuple) Equals(t Types) bool {
	if TUPLE != t.GetTypeType() {
		return false
	}
	if len(t.(Tuple)) != len(tup) {
		return false
	}

	tup2, _ := t.(Tuple)
	for i, v := range tup {
		if !v.Equals(tup2[i]) {
			return false
		}
	}
	return true
}
func (e Error) Equals(t Types) bool {
	return false
}
func (q Qualifier) Equals(t Types) bool {
	return false
}
func (c Class) Equals(t Types) bool {
	if CLASS != t.GetTypeType() {
		return false
	}

	c2 := t.(Class)
	if c2.Name != c.Name || !c2.TypeVariable.Equals(c.TypeVariable) || len(c2.Functions) != len(c.Functions) {
		return false
	}

	for m, f := range c.Functions {
		g, found := c2.Functions[m]
		if !found || !g.Equals(f) {
			return false
		}
	}
	return true
}
func (tup NamedTuple) Equals(t Types) bool {
	if DICTIONARY != t.GetTypeType() {
		return false
	}

	tup2 := t.(NamedTuple)
	if len(tup) != len(tup2) {
		return false
	} 
	for k, v := range tup {
		v2, found := tup2[k]
		if !found || !v2.Equals(v) {
			return false
		}
	}
	return true
}
func (a Application) Equals(t Types) bool {
	if APPLICATION != t.GetTypeType() {
		return false
	}

	a2 := t.(Application)
	if len(a) != len(a2) {
		return false
	}
	for i, ty := range a {
		if !ty.Equals(a2[i]) {
			return false
		}
	}
	return true
}

func (i Int) ReplaceTau(tau Tau, t Types) Types {
	return i
}
func (b Bool) ReplaceTau(tau Tau, t Types) Types {
	return b
}
func (c Char) ReplaceTau(tau Tau, t Types) Types {
	return c
}
func (f Float) ReplaceTau(tau Tau, t Types) Types {
	return f
}
func (tau Tau) ReplaceTau(tau2 Tau, t Types) Types {
	if string(tau) == string(tau2) {
		return t
	}
	return tau
}
func (f Function) ReplaceTauInFunction(tau Tau, t Types) Function {
	var out Function
	out.Domain = f.Domain.ReplaceTau(tau, t)
	out.Codomain = f.Codomain.ReplaceTau(tau, t)
	return f
}
func (f Function) ReplaceTau(tau Tau, t Types) Types {
	return f.ReplaceTauInFunction(tau, t)
}
func (a Array) ReplaceTau(tau Tau, t Types) Types {
	return Array{a.ElemType.ReplaceTau(tau, t)}
}
func (tup Tuple) ReplaceTau(tau Tau, t Types) Types {
	out := make([]Types, len(tup))
	for _, v := range tup {
		out = append(out, v.ReplaceTau(tau, t))
	}
	return Tuple(out)
}
func (e Error) ReplaceTau(_ Tau, _ Types) Types {
	return e
}
func (c Class) ReplaceTau(tau Tau, t Types) Types {
	err.PrintBug() // this should never happen
	panic("")
}
func (q Qualifier) ReplaceTau(_ Tau, _ Types) Types {
	err.PrintBug()
	panic("")
}
func (q NamedTuple) ReplaceTau(_ Tau, _ Types) Types {
	err.PrintBug()
	panic("")
}
func (a Application) ReplaceTau(tau Tau, t Types) Types {
	// replace tau, then try application if applicable
	tryApplication := 
			FUNCTION == t.GetTypeType() &&
			len(a) >= 2 &&
			a[0].Equals(tau)
	newApplication := Application(make([]Types, len(a)))
	for i, ty := range a {
		newApplication[i] = ty.ReplaceTau(tau, t)
	}

	if tryApplication {
		// loop only loops when application is successful and there is at least one element that
		// hasn't been the left-hand side of an application
		for i := 0; i < len(newApplication) - 1; i++ {
			if FUNCTION != newApplication[i].GetTypeType() {
				// cannot do another application
				return newApplication[i:]
			}

			res := newApplication[i].Apply(newApplication[i + 1])
			if APPLICATION == res.GetTypeType() {
				// application did nothing
				return newApplication[i:]
			}
			// application was successful
			newApplication[i + 1] = res
		}

		// there must be only one element left; this follows from the loop's restrictions above
		return newApplication[len(newApplication) - 1]
	}
	return newApplication
}

func DoTypeInference(t Types, newType Types) Types {
	if TAU == newType.GetTypeType() {
		res, ok := inferences.addRule(newType.(Tau), t)
		if ok {
			return res
		}
	} else if t.Equals(newType) {
		return t
	}

	return Error{}
}

func (i Int) InferType(newType Types) Types {
	return DoTypeInference(i, newType)
}
func (b Bool) InferType(newType Types) Types {
	return DoTypeInference(b, newType)
}
func (c Char) InferType(newType Types) Types {
	return DoTypeInference(c, newType)
}
func (f Float) InferType(newType Types) Types {
	return DoTypeInference(f, newType)
}
func (tau Tau) InferType(newType Types) Types {
	res, ok := inferences.addRule(tau, newType)
	if !ok {
		return Error{}
	}
	return res
}
func (f Function) InferType(newType Types) Types { return f.Apply(newType) }
func (a Array) InferType(newType Types) Types {
	if ARRAY != newType.GetTypeType() {
		return Error{}
	}

	if ERROR == a.ElemType.InferType(newType.(Array).ElemType).GetTypeType() {
		return Error{}
	}
	return a
}
func (tup Tuple) InferType(newType Types) Types {
	return DoTypeInference(tup, newType)
}
func (e Error) InferType(newType Types) Types {
	return e
}
func (q Qualifier) InferType(newType Types) Types {
	return q.Apply(newType)
}
func (tups NamedTuple) InferType(newType Types) Types {
	if TAU != newType.GetTypeType() {
		return Error{}
	}
	
	key := string(newType.(Tau))
	for dom, tup := range tups {
		if dom == key {
			return tup
		}
	}

	return Error{}
}
func (c Class) InferType(newType Types) Types {
	err.PrintBug()
	panic("")
}
func (as Application) InferType(newType Types) Types {
	if APPLICATION != newType.GetTypeType() {
		return Error{}
	}

	as2 := newType.(Application)
	if len(as2) != len(as) {
		return Error{} 
	}

	for i, a := range as {
		t := DoTypeInference(a, as2[i])
		if ERROR == t.GetTypeType() {
			return Error{}
		}
	}
	return as
}

type _errorGenFn (func () err.UserMessage)

func typeErrorGen(message string) _errorGenFn {
	return (func () err.UserMessage {
		return err.CompileMessage(message, err.ERROR, err.TYPE, "", 0, 0, 0, "")
	})
}

type typeErrorType int
const (
	E_UNEXPECTED typeErrorType = iota
)
var TypeErrors = map[typeErrorType] _errorGenFn {
	E_UNEXPECTED: typeErrorGen("unexpected type"),
}

// a 1 1; a <- (+)
// => (+) 1 1
// => 2 