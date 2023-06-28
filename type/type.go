package types

import (
	"strings"
	"sync"
	"yew/demangler"
	err "yew/error"
	"yew/source"
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
	DATA
	CONSTRUCTOR
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
type Data struct {
	Name string
	TypeVariables []Tau
	Constructors map[string]Constructor
}
type Constructor struct {
	Name string
	Members Application
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
func (d Data) ToString() string {
	var builder strings.Builder
	builder.WriteString(d.Name)
	builder.WriteByte(' ')
	for _, v := range d.TypeVariables {
		builder.WriteString(v.ToString())
		builder.WriteByte(' ')
	}
	builder.WriteString("::")
	for k, c := range d.Constructors {
		builder.WriteByte(' ')
		builder.WriteString(k)
		builder.WriteByte(' ')
		builder.WriteString(c.Members.ToString())
		builder.WriteString(" |")
	}
	res := builder.String()
	return res[:len(res)-2]
}
func (c Constructor) ToString() string {
	return c.Name + " " + c.Members.ToString()
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
func (d Data) Apply(t Types) Types {
	if t.GetTypeType() == TAU {
		c, found := d.Constructors[t.(Tau).ToString()]
		if found {
			return c
		}
	}
	return Application{d, t}
}
func (c Constructor) Apply(t Types) Types {
	return Application{c, t}
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
func (d Data) GetTypeType() TypeType {
	return DATA
}
func (c Constructor) GetTypeType() TypeType {
	return CONSTRUCTOR
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
func (d Data) Equals(t Types) bool {
	if DATA != t.GetTypeType() {
		return false 
	}

	d2 := t.(Data)
	if d.Name != d2.Name {
		return false
	}

	if len(d.TypeVariables) != len(d2.TypeVariables) {
		return false
	}

	if len(d.Constructors) != len(d2.Constructors) {
		return false
	}

	for i := range d.TypeVariables {
		if !d.TypeVariables[i].Equals(d2.TypeVariables[i]) {
			return false
		}
	}

	for key, val := range d.Constructors {
		val2, found := d2.Constructors[key]
		if !found {
			return false
		}
		if !val.Equals(val2) {
			return false
		}
	}

	return true
}
func (c Constructor) Equals(t Types) bool {
	if t.GetTypeType() != CONSTRUCTOR {
		return false
	}
	c2 := t.(Constructor)
	return c.Name == c2.Name && c.Members.Equals(c2.Members)
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
/*func (a Application) ReplaceTau(tau Tau, t Types) Types {
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
}*/
func (a Application) ReplaceTau(tau Tau, ty Types) Types {
	out := make(Application, len(a))
	for i, v := range a {
		out[i] = v.ReplaceTau(tau, ty)
	}
	return out
}
func (d Data) ReplaceTau(tau Tau, ty Types) Types {
	for _, v := range d.TypeVariables {
		if string(v) == string(tau) {
			for key, con := range d.Constructors {
				d.Constructors[key] = con.ReplaceTau(tau, ty).(Constructor)
			}

			break
		}
	}

	return d
}

func (c Constructor) ReplaceTau(tau Tau, ty Types) Types {
	c.Members = c.Members.ReplaceTau(tau, ty).(Application)
	return c
}

func DoTypeInference(t Types, newType Types) Types {
	if TAU == newType.GetTypeType() {
		res, ok := inferences.addRule(newType.(Tau), t)
		if ok {
			return res
		}
	} else if TAU == t.GetTypeType() {
		res, ok := inferences.addRule(t.(Tau), newType)
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
func (Constructor) InferType(Types) Types {
	err.PrintBug()
	panic("")
}
func (Data) InferType(Types) Types {
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

func (a Application) Head() Types {
	if len(a) == 0 {
		return a
	}
	return a[0]
}
func (a Application) Tail() Types {
	if len(a) == 0 {
		return Tuple{}
	} else if len(a) == 1 {
		return Tuple{}
	} else if len(a) == 2 {
		return a[1]
	}

	return a[1:]
}

func (a Application) Split() (left Types, right Types) {
	return a.Head(), a.Tail()
}

func GrabConstructorName(from Types) (Constructor, bool, err.Error) {
	if from.GetTypeType() != TAU {
		return Constructor{}, false, TypeErrors[E_UNEXPECTED]().(err.Error)
	}
	name := from.(Tau)
	return Constructor{Name: string(name), Members: make(Application, 0)}, true, err.Error{}
}

func ToConstructor(from Types) (Constructor, bool) {
	tt := from.GetTypeType()
	if tt == TAU {
		c, ok, e := GrabConstructorName(from)
		if !ok {
			e.Print()
		}
		return c, ok
	} else if tt == APPLICATION {
		head, tail := from.(Application).Split()
		c, ok, e := GrabConstructorName(head)
		if !ok {
			e.Print()
			return c, ok
		}

		if tail.GetTypeType() == APPLICATION {
			c.Members = tail.(Application)
		} else {
			c.Members = Application{tail}
		}
		return c, ok
	}

	TypeErrors[E_UNEXPECTED]().Print()
	return Constructor{}, false
}

func GrabDataName(from Types) (Data, bool, err.Error) {
	if from.GetTypeType() != TAU {
		return Data{}, false, TypeErrors[E_UNEXPECTED]().(err.Error)
	}
	name := from.(Tau)
	//println("name =", string(name), "len(", len(string(name)), ")")
	return Data{
		Name: string(name), 
		Constructors: make(map[string]Constructor, 0),
		TypeVariables: make([]Tau, 0),
	}, true, err.Error{}
}

func redeclareTypeVarError(name string, vars Application, i int) err.Error {
	var builder strings.Builder
	builder.WriteString("type variable redeclared:\n")
	spaces, _ := builder.WriteString(name)
	builder.WriteByte(' ')
	spaces = spaces + 1

	for j, v := range vars {
		tmp, _ := builder.WriteString(v.ToString())
		builder.WriteByte(' ')
		tmp = tmp + 1

		if j < i {
			spaces = spaces + tmp
		}
	}
	for j := 0; j < spaces; j++ {
		builder.WriteByte(' ')
	}
	builder.WriteString("^-- here")
 	return typeErrorGen(builder.String())().(err.Error)
}

func addTypeVariables(d Data, vars Types) (Data, bool) {
	if vars.GetTypeType() == APPLICATION {
		maybeVars := vars.(Application)
		d.TypeVariables = make([]Tau, 0, len(maybeVars))
		tVarMap := make(map[string]Tau, len(maybeVars))

		for i, m := range maybeVars {
			if m.GetTypeType() != TAU {
				TypeErrors[E_UNEXPECTED]().Print()
				return d, false
			}
			tau := m.(Tau)
			tauString := string(tau)
			_, redeclared := tVarMap[tauString]
			if redeclared {
				// type variable redeclared
				redeclareTypeVarError(d.Name, maybeVars, i).Print()
				return d, false
			}
			// else add var to tracking map
			tVarMap[tauString] = tau
			// add to type variables
			d.TypeVariables = append(d.TypeVariables, tau)
		}
	} else if vars.GetTypeType() == TAU {
		d.TypeVariables = append(d.TypeVariables, vars.(Tau))
	} else {
		TypeErrors[E_UNEXPECTED]().Print()
		return d, false
	}

	return d, true
}

func MakeConstructor(name string, members Application) Constructor {
	return Constructor{Name: name, Members: members}
}

// should only be used after verifying all constructors in `cs` have unique names 
func mapConstructors(cs []Constructor) map[string]Constructor {
	out := make(map[string]Constructor, len(cs))
	for _, c := range cs {
		_, found := out[c.Name]
		if found {
			err.PrintBug()
			panic("")
		}
		out[c.Name] = c
	}
	return out
}

func MakeTaus(s ...string) []Tau {
	out := make([]Tau, len(s))
	for i := range s {
		out[i] = Tau(s[i])
	}
	return out
}

func MakeData2(name string, vars []string, cons []Constructor) Data {
	return Data{
		Name: name,
		TypeVariables: MakeTaus(vars...),
		Constructors: mapConstructors(cons),
	}
}

func MakeData(name string, vars []Tau, cons []Constructor) Data {
	return Data{
		Name: name,
		TypeVariables: vars, 
		Constructors: mapConstructors(cons),
	}
}

func ToData(from Types) (Data, bool) {
	tt := from.GetTypeType()
	if tt == TAU {
		d, ok, e := GrabDataName(from)
		if !ok {
			e.Print()
		}
		return d, ok
	} else if tt == APPLICATION {
		head, tail := from.(Application).Split()
		d, ok, e := GrabDataName(head)
		if !ok {
			e.Print()
			return d, ok
		}

		return addTypeVariables(d, tail)
	}

	TypeErrors[E_UNEXPECTED]().Print()
	return Data{}, false
}

type _errorGenFn (func () err.UserMessage)

func typeErrorGen(message string) _errorGenFn {
	return (func () err.UserMessage {
		return err.CompileMessage(message, err.ERROR, err.TYPE, "", 0, 0, source.Source{""})
	})
}

type typeErrorType int
const (
	E_UNEXPECTED typeErrorType = iota
	E_EXPECTED_ARRAY
	E_EXPECTED_TUPLE
)
var TypeErrors = map[typeErrorType] _errorGenFn {
	E_UNEXPECTED: typeErrorGen("unexpected type"),
	E_EXPECTED_ARRAY: typeErrorGen("expected array type"),
	E_EXPECTED_TUPLE: typeErrorGen("expected tuple type"),
}

// a 1 1; a <- (+)
// => (+) 1 1
// => 2 