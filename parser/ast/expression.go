package ast

import (
	types "yew/type"
	. "yew/parser/parser"
	//"yew/symbol"
)

type Expression interface {
	Ast
	ExpressionType() types.Types
	DoTypeInference(newTypeInformation types.Types) types.Types
}