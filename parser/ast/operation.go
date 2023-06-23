package ast

import (
	builtin "yew/builtin"
	"yew/ir"
	"yew/parser/parser"
	types "yew/type"
)

type Operation interface {
	AsFunction(p *parser.Parser) Function
}

func buildQualified(class types.Class) func(types.Tau) func(types.Types) types.Qualifier {
	return func(typeVariable types.Tau) func(types.Types) types.Qualifier {
		return func(t types.Types) types.Qualifier {
			return types.Qualifier{
				Class:        class,
				TypeVariable: typeVariable,
				Qualified:    t,
			}
		}
	}
}

var arith = buildQualified(builtin.Number)("a")
var relate = buildQualified(builtin.Equalable)("a")
var order = buildQualified(builtin.Orderable)("a")
var list = buildQualified(builtin.Listable)("a")
var factorial = types.Function{Domain: types.Int{}, Codomain: types.Int{}}

func (u UOpType) Compile(builder *ir.IrBuilder) {

}
