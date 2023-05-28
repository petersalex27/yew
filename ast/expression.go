package ast

import (
	types "yew/type"
	"yew/symbol"
)

type Expression interface {
	Ast
	ExpressionType() types.Types
	ResolveNames(*symbol.SymbolTable)
	DoTypeInference(newTypeInformation types.Types) types.Types
}