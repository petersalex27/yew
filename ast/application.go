package ast

import (
	"fmt"
	"yew/symbol"
	"yew/type"
)

type Application struct {
	left Expression
	right Expression
} 

func (app Application) ExpressionType() types.Types {
	return app.left.ExpressionType().Apply(app.right.ExpressionType())
}
func (app Application) ResolveNames(*symbol.SymbolTable) {
	// TODO
}
func (app Application) DoTypeInference(newTypeInformation types.Types) types.Types {
	// TODO
	panic("")
}

func (app Application) GetNodeType() NodeType { return APPLICATION }

func (app Application) Make(stack *AstStack) bool {
	valid := stack.Validate([]NodeType{FUNCTION, EXPRESSION})
	valid = stack.Validate([]NodeType{EXPRESSION, EXPRESSION})
	if !valid {
		return false
	}
	app.right = stack.Pop().(Expression)
	app.left = stack.Pop().(Expression)
	stack.Push(app)
	return true
}

func (app Application) equal_test(a Ast) bool {
	equal := a.GetNodeType() == APPLICATION
	app2 := a.(Application)
	return equal &&
			app2.left.equal_test(app.left) &&
			app2.right.equal_test(app.right)
}
func (app Application) print(n int) {
	printSpaces(n)
	fmt.Printf("Application\n")
	app.left.print(n + 1)
	app.right.print(n + 1)
}
func MakeApplication(left Expression, right Expression) Application {
	return Application{left: left, right: right}
}