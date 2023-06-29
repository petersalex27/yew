package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type Application struct {
	left  Expression
	right Expression
}

func (app Application) Split() (Expression, Expression) {
	return app.split()
}

func (app Application) ExpressionType() types.Types {
	return app.left.ExpressionType().Apply(app.right.ExpressionType())
}
func (app Application) ResolveNames(table *symbol.SymbolTable) bool {
	return app.left.ResolveNames(table) && app.right.ResolveNames(table)
	// TODO
}
func (app Application) DoTypeInference(newTypeInformation types.Types) types.Types {
	// TODO
	panic("")
}

func (app Application) GetNodeType() NodeType { return APPLICATION }

func (app Application) FindStartToken() scan.Token {
	return app.left.FindStartToken()
}

func (app Application) Make(p *Parser) bool {
	valid, _ := p.Stack.TryValidate(appRule1.Expression)
	if valid {
		app.right = p.Stack.Pop().(Expression)
		app.left = p.Stack.Pop().(Function).function
	} else {
		var e func (scan.InputStream) err.Error
		valid, e = p.Stack.Validate(appRule2)
		if !valid {
			e(p.Input).Print()
			return false
		}
		app.right = p.Stack.Pop().(Expression)
		app.left = p.Stack.Pop().(Expression)
	}

	p.Stack.Push(app)
	return true
}

func (app Application) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == APPLICATION
	app2, ok := a.(Application)
	if !ok {
		return false
	}
	return equal &&
		app2.left.Equal_test(app.left) &&
		app2.right.Equal_test(app.right)
}
func (app Application) Print(lines []string) {
	next := make([]string, len(lines))
	next = append(next, lines...)
	next = printLines(next)
	fmt.Printf("Application\n")
	next = append(next, " ├─")
	app.left.Print(next)
	next[len(next)-1] = " └─"
	app.right.Print(next)
}
func MakeApplication(left Expression, right Expression) Application {
	return Application{left: left, right: right}
}
