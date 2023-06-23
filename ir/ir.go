package ir

import (
	"strconv"
	"strings"
	types "yew/type"
	err "yew/error"
	"yew/utils"
)

type IrType interface {
	util.Stringable
	ToYewType() types.Types
}

type Int int
type Float struct {}
type Double struct {}
type Struct []IrType
type Array struct {
	size uint64
	elementType IrType
}
type Void struct {}
type Pointer struct {to IrType}
type Empty struct {}
type Function struct {

}

func (e Empty) ToString() string { return "" }
func (i Int) ToString() string {
	return "i" + strconv.Itoa(int(i))
}
func (f Float) ToString() string {
	return "float"
}
func (d Double) ToString() string {
	return "double"
}
func (s Struct) ToString() string {
	var builder strings.Builder
	builder.WriteByte('{')
	length := len(s)
	for i, t := range s {
		builder.WriteString(t.ToString())
		if i + 1 < length {
			builder.WriteByte(';')
		} else {
			builder.WriteByte('}')
		}
	}
	return builder.String()
}
func (a Array) ToString() string {
	// [<size> x <ty>]
	return "[" + strconv.Itoa(int(a.size)) + " x " +  a.elementType.ToString() + "]"
}
func (v Void) ToString() string {
	return "void"
}
func (pointer Pointer) ToString() string {
	return "*" + pointer.to.ToString()
}
func (e Empty) ToYewType() types.Types {
	return types.Tuple{}
}
func (i Int) ToYewType() types.Types {
	if i == 1 {
		return types.Bool{}
	} else if i <= 8 {
		return types.Char{}
	} else if i <= 64 {
		return types.Int{}
	} else {
		err.PrintBug()
		panic("")
	}
}
func (f Float) ToYewType() types.Types {
	return types.Float{}
}
func (d Double) ToYewType() types.Types {
	return types.Float{}
}
func (s Struct) ToYewType() types.Types {
	out := make(types.Tuple, len(s))
	for i, t := range s {
		out[i] = t.ToYewType()
	}
	return out
}
func (a Array) ToYewType() types.Types {
	return types.Array{ElemType: a.elementType.ToYewType()}
}

type IrBuilder struct {
	builder strings.Builder
}

