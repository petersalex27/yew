package ir

import "strings"

type Operand struct {
	OpdType IrType
	Name string
}

type Instruction struct {
	name string
	resType IrType
	operands []Operand
}

func genUnary(name string) (func (ty IrType) (func (o Operand) Instruction)) {
	return func(ty IrType) func(o Operand) Instruction {
		return func(o Operand) Instruction {
			return Instruction{name: name, resType: ty, operands: []Operand{o}}
		}
	}
}
func genBinary(name string) (func (ty IrType) (func (o1 Operand) (func (o2 Operand) Instruction))) {
	return func(ty IrType) func(o1 Operand) func(o2 Operand) Instruction {
		return func(o1 Operand) func(o2 Operand) Instruction {
			return func(o2 Operand) Instruction {
				return Instruction{name: name, resType: ty, operands: []Operand{o1, o2}}
			}
		}
	}
}

var AddInt = genBinary("add")(Int(64))
var AddDouble = genBinary("fadd")(Double{})
var SubInt = genBinary("sub")(Int(64))
var SubDouble = genBinary("fsub")(Double{})
var MulInt = genBinary("mul")(Int(64))
var MulDouble = genBinary("fmul")(Double{})
var DivInt = genBinary("sdiv")(Int(64))
var DivDouble = genBinary("fdiv")(Double{})
var NegInt = genBinary("sub")(Int(64))(Operand{Int(64), "0"})
var NegDouble = genBinary("fsub")(Double{})(Operand{Double{}, "0.0"})

func (instr Instruction) ToString() string {
	var builder strings.Builder
	builder.WriteString(instr.name)
	builder.WriteByte(' ')
	builder.WriteString(instr.resType.ToString())
	opdsLen := len(instr.operands)
	for i/*, o*/ := range instr.operands {
		//builder.WriteString(o.opdType.ToString())
		builder.WriteByte(' ')
	//	builder.WriteString(o.name)
		if i + 1 < opdsLen {
			builder.WriteByte(',')
		}
	}
	return builder.String()
}