package ast

import (
	//err "yew/error"
	scan "yew/lex"
	. "yew/parser/parser"
	types "yew/type"
	//"yew/symbol"
)

type Expression interface {
	Ast
	ExpressionType() types.Types
	DoTypeInference(newTypeInformation types.Types) types.Types
	FindStartToken() scan.Token
}