package value

import (
	"strconv"
	strings "strings"
	"yew/type"
	util "yew/utils"
)

type ValueType int

const (
	INT = iota
	CHAR 
	BOOL
	FLOAT
	ARRAY
	STRUCT
	TUPLE
	FUNCTION
)

type Value interface {
	util.Stringable
	GetValueType() ValueType
	GetType() types.Types
}
type Array struct {
	elements []Value
	elementType types.Types
}
type Tuple []Value
type Bool bool
type Char byte
type Int int64
type Float float64
type parameter struct {
	name string
	paramType types.Types
	codeSubIndex int
}
type Function struct {
	name string
	demangle string
	params []parameter
	codeIndex int
}
type NamedTuple struct {
	params []parameter
}

type ValuedType struct {
	value Value
}
func (v ValuedType) GetSuperType() types.Types {
	return v.value.GetType()
}
func (v ValuedType) ToString() string { return v.value.ToString() }
func (v ValuedType) GetTypeType() types.TypeType { return types.VALUE }
func (v ValuedType) Equals(t types.Types) bool {
	if t.GetTypeType() != types.VALUE {
		return false
	}
	return types.CanBeReplaced(v, t) && // checks that super types are equal
		v.value.ToString() == t.ToString() // checks that their representation is equal
}
func (v ValuedType) Apply(t types.Types) types.Types {
	if v.value.GetType().GetTypeType() == types.FUNCTION {
		return v.value.GetType().Apply(t)
	}
	return types.Application([]types.Types{v, t})
}
func (v ValuedType) ReplaceTau(tau types.Tau, t types.Types) types.Types {
	// TODO: this might be wrong--especially if value itself must be updated (not just type 
	// 	returned from GetType)
	return v.value.GetType().ReplaceTau(tau, t)
}
func (v ValuedType) InferType(from types.Types) types.Types {
	// TODO: this is problem wrong
	return v.value.GetType().InferType(from)
}

func (a Array) ToString() string {
	return "[" + 
		strings.Join(util.Fmap(a.elements, func(v Value) string {
			return v.ToString()
		}), ", ") + 
		"]"
}
func (t Tuple) ToString() string {
	return "(" + 
		strings.Join(util.Fmap(t, func(v Value) string {
			return v.ToString()
		}), ", ") + 
		")"
}
func (b Bool) ToString() string {
	if b {
		return "True"
	} else {
		return "False"
	}
}
func (c Char) ToString() string {
	return string(c)
}
func (i Int) ToString() string {
	return strconv.Itoa(int(i))
}
func (f Float) ToString() string {
	out := strconv.FormatFloat(float64(f), 'f', -1, 64)
	return out
}
func (f Function) ToString() string {
	return f.name
}
func (n NamedTuple) ToString() string {
	return "(" + 
		strings.Join(util.Fmap(n.params, func(p parameter) string {
			return p.name
		}), ", ") + 
		")"
}

func (a Array) GetValueType() ValueType {
	return ARRAY
}
func (t Tuple) GetValueType() ValueType {
	return TUPLE
}
func (b Bool) GetValueType() ValueType {
	return BOOL
}
func (c Char) GetValueType() ValueType {
	return CHAR
}
func (i Int) GetValueType() ValueType {
	return INT
}
func (f Float) GetValueType() ValueType {
	return FLOAT
}
func (f Function) GetValueType() ValueType {
	return FUNCTION
}
func (n NamedTuple) GetValueType() ValueType {
	return TUPLE
}

func (a Array) GetType() types.Types {
	return a.elementType
}
func (t Tuple) GetType() types.Types {
	return types.Tuple(util.Fmap(t, func(v Value) types.Types {
		return v.GetType()
	}))
}
func (b Bool) GetType() types.Types {
	return types.Bool{}
}
func (c Char) GetType() types.Types {
	return types.Char{}
}
func (i Int) GetType() types.Types {
	return types.Int{}
}
func (f Float) GetType() types.Types {
	return types.Float{}
}
func (f Function) GetType() types.Types {
	return util.FoldRight(
		f.params[:len(f.params) - 2], 
		types.Function{
			Domain: f.params[len(f.params) - 2].paramType, 
			Codomain: f.params[len(f.params) - 1].paramType,
		}, 
		func(base types.Function, p parameter) types.Function {
			return types.Function{Domain: p.paramType, Codomain: base}
		})
}
func (n NamedTuple) GetType() types.Types {
	return util.FoldLeft(
		n.params[1:],
		types.Tuple{n.params[0].paramType},
		func(base types.Tuple, p parameter) types.Tuple {
			return append(base, p.paramType)
		},
	)
}

func MakeArray[T types.Types, K Value](elems []K) Array {
	var t T
	value := make([]Value, len(elems))
	for i, e := range elems {
		value[i] = e
	}
	return Array{elements: value, elementType: t}
}