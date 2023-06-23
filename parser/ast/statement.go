package ast

import (
	"yew/symbol"
	. "yew/parser/parser"
)

type Statement interface {
	Ast
	GetSymbol() symbol.Symbolic
}