package types

import (
	"strings"
	"sync"
	"yew/demangler"
	err "yew/error"
	"yew/info"
	"yew/source"
	util "yew/utils"
)

type Superable interface {
	GetSuperType() Types
}

type inferenceRules struct {
	inferenceLock sync.Mutex
	rules         map[string]*[]Types
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
			name := ty.(Tau).name
			if inf.recursiveCheckForCycles(name, marked, visited) {
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
	if _, found := inf.rules[tau.name]; !found {
		return tau, false
	}

	var prev string = tau.name
	var next *[]Types
	var found bool
	stack := make([]string, 0, len(inf.rules)/2)
	searched := make(map[string]bool, len(inf.rules)/2)
	stack = append(stack, tau.name)
	for len(stack) > 0 {
		prev = stack[len(stack)-1]
		searched[prev] = true
		stack = stack[:len(stack)-1]
		next, found = inf.rules[prev]
		if !found {
			continue
		}
		for _, t := range *next {
			if t.GetTypeType() == TAU {
				x := t.(Tau).name
				if !searched[x] {
					stack = append(stack, t.(Tau).name)
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
		if nil == inf.rules[tau.name] {
			inf.rules[tau.name] = new([]Types)
			*inf.rules[tau.name] = make([]Types, 1)
			(*inf.rules[tau.name])[0] = to
			res, ok = to, true
		} else {
			// check if contained in rule already
			var exit bool = false
			for _, t := range *inf.rules[tau.name] {
				if t.Equals(to) {
					exit = true
					break
				}
			}
			if !exit {
				// is not already in rules, so append
				(*inf.rules[tau.name]) = append((*inf.rules[tau.name]), to)
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

func genVarTau(name string) Tau {
	return Tau{name: name, Loc: info.DefaultLoc()}
}

func GetNewTau() Tau {
	return genVarTau(demangler.TYPE.GetDemanglerPrefix())
}

func GetNewTaus(n int) (taus []Tau) {
	taus = make([]Tau, n)
	demangles := demangler.TYPE.GetDemanglerPrefixes(n)
	for i, d := range demangles {
		taus[i] = genVarTau(d)
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
	GetLocation() info.Location
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

func Var(s string) Tau { return genVarTau(s) }

func MakeTau(s string, l info.Loc) Tau { return Tau{name: s, Loc: l} }

type Int info.Loc
type Float info.Loc
type Char info.Loc
type Bool info.Loc
type Tau struct {
	name string
	Loc  info.Loc
}
type Function struct {
	Domain   Types
	Codomain Types
}
type Array struct {
	ElemType Types
	Loc      info.Loc
}
type Tuple []Types
type NamedTuple map[string]Types
type Class struct {
	Loc          info.Loc
	Name         string
	TypeVariable Tau
	Functions    map[string]Function
}
type Context struct {
	ClassName    Tau
	TypeVariable Tau
}
type ConstraintContext []Context
type Constraint struct {
	Context     ConstraintContext
	Constrained Types
	Loc         info.Loc
}
type Data struct {
	Name          string
	TypeVariables []Tau
	Constructors  map[string]Constructor
	Loc           info.Loc
}
type Constructor struct {
	Name    string
	Members Application
	Loc     info.Loc
}
type Error err.Error

func (e Error) ToError() err.Error {
	return err.Error(e)
}

func MakeContext(name string, variable string) Context {
	return Context{ClassName: Var(name), TypeVariable: Var(variable)}
}
func ConstrainType(constrained Types, cxts ...Context) Constraint {
	return Constraint{
		Context: cxts,
		Constrained: constrained,
	}
}

func (a Application) ValidClass() (bool, string, info.Locatable) {
	if len(a) != 2 {
		if len(a) < 2 {
			return false, "too few type variables, expected one", a
		}
		// get 3rd type in application
		return false, "too many type variables, expected one", a[2]
	}

	if a[0].GetTypeType() != TAU {
		return false, "expected class declaration", a[0]
	}

	if a[1].GetTypeType() != TAU {
		return false, "expected type variable", a[1]
	}
	return true, "", a
}

func (c Constraint) ConstrainApplication(a Application) (Class, bool, string, info.Locatable) {
	valid, msg, loc := a.ValidClass()
	if !valid {
		return Class{}, false, msg, loc
	}

	var out Class
	out.Name = a[0].(Tau).name
	out.TypeVariable = a[1].(Tau)

	for _, cxt := range c.Context {
		if !cxt.TypeVariable.Equals(out.TypeVariable) {
			return out, false, "illegal type parameter in type constraint's context", cxt.TypeVariable
		}
	}

	out.Loc = c.GetLocation().(info.Loc)
	return out, true, "", out
}

func (c Constraint) Constrain(f Function) Function {
	if f.Domain.GetTypeType() == FUNCTION {
		f.Domain = c.Constrain(f.Domain.(Function))
	} else if f.Domain.GetTypeType() == QUALIFIER {
		c2 := f.Domain.(Constraint)
		cxt := make(ConstraintContext, 0)
		cxt = append(cxt, c.Context...)
		cxt = append(cxt, c2.Context...)
		f.Domain = Constraint{
			Context:     cxt,
			Constrained: c2.Constrained,
		}
	} else {
		f.Domain = Constraint{
			Context:     c.Context,
			Constrained: f.Domain,
		}
	}

	if f.Codomain.GetTypeType() == FUNCTION {
		f.Codomain = c.Constrain(f.Codomain.(Function))
	} else if f.Codomain.GetTypeType() == QUALIFIER {
		c2 := f.Codomain.(Constraint)
		cxt := make(ConstraintContext, 0)
		cxt = append(cxt, c.Context...)
		cxt = append(cxt, c2.Context...)
		f.Codomain = Constraint{
			Context:     cxt,
			Constrained: c2.Constrained,
		}
	} else {
		f.Codomain = Constraint{
			Context:     c.Context,
			Constrained: f.Codomain,
		}
	}
	return f
}

// represents a type that has not yet been evaluated; by the end of the resolve type phase,
// there should be no application types left
type Application []Types

func (i Int) GetLocation() info.Location {
	return info.Loc(i)
}
func (b Bool) GetLocation() info.Location {
	return info.Loc(b)
}
func (c Char) GetLocation() info.Location {
	return info.Loc(c)
}
func (f Float) GetLocation() info.Location {
	return info.Loc(f)
}
func (tau Tau) GetLocation() info.Location {
	return tau.Loc
}
func (f Function) GetLocation() info.Location {
	return f.Domain.GetLocation()
}
func (a Array) GetLocation() info.Location {
	return a.Loc
}
func (tup Tuple) GetLocation() info.Location {
	if len(tup) == 0 {
		return info.DefaultLoc()
	}
	return tup[0].GetLocation()
}
func (e Error) GetLocation() info.Location {
	return e.ToError().GetLocation()
}
func (c ConstraintContext) GetLocation() info.Location {
	if len(c) == 0 {
		return info.DefaultLoc()
	}
	return c[0].ClassName.Loc
}
func (c Constraint) GetLocation() info.Location {
	return c.Loc
}
func (c Class) GetLocation() info.Location {
	return c.Loc
}
func (tup NamedTuple) GetLocation() info.Location {
	return info.DefaultLoc() // TODO
}
func (a Application) GetLocation() info.Location {
	if len(a) == 0 {
		return info.DefaultLoc()
	}
	return a[0].GetLocation()
}
func (d Data) GetLocation() info.Location {
	return d.Loc
}
func (c Constructor) GetLocation() info.Location {
	return c.Loc
}

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
	return string(tau.name)
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

	return left + " -> " + right
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
func (c ConstraintContext) ToString() string {
	if len(c) == 1 {
		return c[0].ClassName.ToString() + " " + c[0].TypeVariable.ToString()
	}

	ss := util.Fmap(c, func(t Context) string {
		return t.ClassName.ToString() + " " + t.TypeVariable.ToString()
	})
	return "(" + strings.Join(ss, ", ") + ")"
}
func (c Constraint) ToString() string {
	top := c.Context.ToString() + " => "
	if c.Constrained == nil {
		return top + "_"
	}
	return "(" + top + c.Constrained.ToString() + ")"
}
func (c Class) ToString() string {
	var tmp []string
	for k, f := range c.Functions {
		tmp = append(tmp, k+" "+f.ToString())
	}
	return c.Name + " " + c.TypeVariable.ToString() + "{" + strings.Join(tmp, "; ") + "}"
}
func (tup NamedTuple) ToString() string {
	var tmp []string
	for k, v := range tup {
		tmp = append(tmp, k+" "+v.ToString())
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

	return res[:len(res)-1]
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
func (c Constraint) Apply(t Types) Types {
	err.PrintBug()
	panic("")
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
func (c Constraint) GetTypeType() TypeType {
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
	return TAU == t.GetTypeType() && tau.name == t.(Tau).name
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
func (cxt Context) EqualsContext(cxt2 Context) bool {
	return cxt.ClassName.name == cxt2.ClassName.name &&
		cxt.TypeVariable.name == cxt2.TypeVariable.name
}
func (c1 ConstraintContext) Equals(c2 ConstraintContext) bool {
	if len(c1) != len(c2) {
		return false
	}

	for i, cxt := range c1 {
		if !cxt.EqualsContext(c2[i]) {
			return false
		}
	}
	return true
}
func (c Constraint) Equals(t Types) bool {
	if t.GetTypeType() != QUALIFIER {
		return false
	}
	c2, ok := t.(Constraint)
	if !ok {
		return false
	}

	return c.Context.Equals(c2.Context) &&
		c2.Constrained.Equals(c.Constrained)
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
	if tau.name == tau2.name {
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
	return Array{ElemType: a.ElemType.ReplaceTau(tau, t), Loc: a.Loc}
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
func (c Constraint) ReplaceTau(_ Tau, _ Types) Types {
	err.PrintBug()
	panic("")
}
func (q NamedTuple) ReplaceTau(_ Tau, _ Types) Types {
	err.PrintBug()
	panic("")
}

/*
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
*/
func (a Application) ReplaceTau(tau Tau, ty Types) Types {
	out := make(Application, len(a))
	for i, v := range a {
		out[i] = v.ReplaceTau(tau, ty)
	}
	return out
}
func (d Data) ReplaceTau(tau Tau, ty Types) Types {
	for _, v := range d.TypeVariables {
		if v.name == tau.name {
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
func (c Constraint) InferType(newType Types) Types {
	return c.Apply(newType)
}
func (tups NamedTuple) InferType(newType Types) Types {
	if TAU != newType.GetTypeType() {
		return Error{}
	}

	key := newType.(Tau).name
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
	return Constructor{Name: name.name, Members: make(Application, 0)}, true, err.Error{}
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
		Name:          name.name,
		Constructors:  make(map[string]Constructor, 0),
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
			tauString := tau.name
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
		out[i] = genVarTau(s[i])
	}
	return out
}

func MakeData2(name string, vars []string, cons []Constructor) Data {
	return Data{
		Name:          name,
		TypeVariables: MakeTaus(vars...),
		Constructors:  mapConstructors(cons),
	}
}

func MakeData(name string, vars []Tau, cons []Constructor) Data {
	return Data{
		Name:          name,
		TypeVariables: vars,
		Constructors:  mapConstructors(cons),
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

type _errorGenFn (func() err.UserMessage)

func typeErrorGen(message string) _errorGenFn {
	return (func() err.UserMessage {
		return err.CompileMessage(message, err.ERROR, err.TYPE, "", 0, 0, source.Source{""})
	})
}

type typeErrorType int

const (
	E_UNEXPECTED typeErrorType = iota
	E_EXPECTED_ARRAY
	E_EXPECTED_TUPLE
)

var TypeErrors = map[typeErrorType]_errorGenFn{
	E_UNEXPECTED:     typeErrorGen("unexpected type"),
	E_EXPECTED_ARRAY: typeErrorGen("expected array type"),
	E_EXPECTED_TUPLE: typeErrorGen("expected tuple type"),
}

// a 1 1; a <- (+)
// => (+) 1 1
// => 2
