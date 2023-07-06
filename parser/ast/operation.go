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

var arith = buildQualified(types.Var("Number"))(types.Var("a"))
var relate = buildQualified(types.Var("Equitable"))(types.Var("a"))
var order = buildQualified(types.Var("Orderable"))(types.Var("a"))
var list = buildQualified(types.Var("Listable"))(types.Var("a"))
var factorial = types.Function{Domain: types.Int{}, Codomain: types.Int{}}

func (u UOpType) Compile(builder *ir.IrBuilder) {

}
