package ast

import "yew/ir"

type InstructionName int
const (
	Add InstructionName = iota
	FAdd
	Sub
	FSub
	Mul
	FMul
	Div
	FDiv
)

type Instruction struct {
	
	ty ir.IrType

}