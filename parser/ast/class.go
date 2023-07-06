package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	nodetype "yew/parser/node-type"
	"yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type Class struct {
	name Id
	functions map[string]types.Function
}

func MakeClass(name Id, fns map[string]types.Function) Class {
	return Class{
		name: name,
		functions: fns,
	}
}

func InitClass(name Id) Class {
	return Class{name: name, functions: make(map[string]types.Function)}
}

func (c Class) GetSymbol() symbol.Symbolic {
	return symbol.MakeSymbol(c.name.token)
}

func constructClass(p *parser.Parser) (bool, err.Error) {
	ok, e := p.Stack.Validate(classRule)
	if !ok {
		return false, e(p.Input)
	}

	annot := p.Stack.Pop().(ExpressionTypeAnnotation)
	class := p.Stack.Pop().(Class)

	if annot.expression.GetNodeType() != nodetype.IDENTIFIER {
		eLoc := p.Input.MakeErrorLocation(annot.expression.FindStartToken())
		e := err.SyntaxError("expected a function declaration", eLoc)
		return false, e
	}

	id := annot.expression.(Id)
	ty := annot.expressionType

	if ty.GetTypeType() != types.FUNCTION {
		eLoc := p.Input.MakeErrorLocation(ty)
		e := err.TypeError("unexpected type, expected a function type", eLoc)
		return false, e
	}

	_, found := class.functions[id.GetName()]
	if found {
		e := err.NameError(
			"illegal redefinition of " + id.GetName() + 
			" in the " + class.name.GetName() + " class",
			p.Input.MakeErrorLocation(id.token),
		)
		return false, e
	}

	if p.HasConstraint {
		// constraint conflicts have already been checked for
		ty = p.ClassConstraint.Constrain(ty.(types.Function))
	}

	class.functions[id.GetName()] = ty.(types.Function)
	p.Stack.Push(class)
	return true, err.Error{}
} 

func (Class) Make(p *parser.Parser) bool {
	ok, e := constructClass(p)
	if !ok {
		e.Print()
	}
	return ok
}

func (c Class) GetNodeType() nodetype.NodeType {
	return nodetype.CLASS_DEFINITION
}

func (c Class) Equal_test(ast parser.Ast) bool {
	if ast.GetNodeType() != nodetype.CLASS_DEFINITION {
		return false
	}
	c2 := ast.(Class)
	if !c.name.Equal_test(c2.name) {
		return false
	}
	if len(c.functions) != len(c2.functions) {
		return false 
	}

	for k, v := range c.functions {
		v2, found := c2.functions[k]
		if !found {
			return false
		}
		if !v.Equals(v2) {
			return false
		}
	}
	return true
}

func (c Class) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("Class\n")
	lines = append(lines, " ├─")
	if len(c.functions) == 0 {
		lines[len(lines)-1] = " └─"
	}
	c.name.Print(lines)

	i := 0
	ln := len(c.functions)

	for k, v := range c.functions {
		i++
		if i == ln {
			lines[len(lines)-1] = " └─"
		}
		printLines(lines)
		fmt.Printf("%s :: %s\n", k, v.ToString())
	}
}

func (c Class) ResolveNames(table *symbol.SymbolTable) bool {
	panic("TODO") // TODO
}

func (c Class) FindStartToken() scan.Token {
	return c.name.token
}