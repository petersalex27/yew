package ast

import (
	"yew/ir"
	"yew/parser/parser"
	types "yew/type"
)

type Operation interface {
	AsFunction(p *parser.Parser) Function
}

func buildQualified(class types.Tau) func(types.Tau) func(types.Function) types.Constraint {
	return func(typeVariable types.Tau) func(types.Function) types.Constraint {
		return func(t types.Function) types.Constraint {
			return types.Constraint{
				Context: types.ConstraintContext{
					types.Context{
						ClassName: class,
						TypeVariable: typeVariable,
					},
				},
				Constrained: t,
			}
		}
	}
}

var arith = buildQualified("Number")("a")
var relate = buildQualified("Equitable")("a")
var order = buildQualified("Orderable")("a")
var list = buildQualified("Listable")("a")
var factorial = types.Function{Domain: types.Int{}, Codomain: types.Int{}}

func (u UOpType) Compile(builder *ir.IrBuilder) {

}
