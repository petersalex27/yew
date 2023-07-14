package ast

import (
	"fmt"
	scan "yew/lex"
	nodetype "yew/parser/node-type"
	"yew/parser/parser"
)

type Annotation struct {
	id Id
	args []parser.Ast
}

func MakeAnnotation(id Id, args ...parser.Ast) Annotation {
	return Annotation{id: id, args: args}
}

func (a Annotation) Make(p *parser.Parser) bool {
	ok, _ := p.Stack.TryValidate([]nodetype.NodeType{nodetype.ANNOTATION}) 
	if ok {
		return true // nothing to do, annot has no args
	}
	
	if valid, e := p.Stack.Validate(annotationRule); !valid {
		e(p.Input).Print()
		return false
	}

	thing := p.Stack.Pop()
	annot := p.Stack.Pop().(Annotation)
	if len(annot.args) == 0 {
		annot.args = make([]parser.Ast, 0, 1)
	}
	annot.args = append(annot.args, thing)
	p.Stack.Push(annot)
	return true
}

func (a Annotation) GetNodeType() nodetype.NodeType {
	return nodetype.ANNOTATION
}

func (a Annotation) Equal_test(ast parser.Ast) bool {
	if ast.GetNodeType() != nodetype.ANNOTATION {
		return false
	}
	a2, ok := ast.(Annotation)
	if !ok {
		return false 
	}
	if !a.id.Equal_test(a2.id) {
		return false
	}
	if len(a.args) != len(a2.args) {
		return false
	}
	for i := range a.args {
		if !a.args[i].Equal_test(a2.args[i]) {
			return false
		}
	}
	return true
}

func (a Annotation) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("Annotation\n")
	lines = append(lines, " ├─")
	if len(a.args) == 0 {
		lines[len(lines)-1] = " └─"
	}
	a.id.Print(lines)
	for i := 0; i < len(a.args)-1; i++ {
		a.args[i].Print(lines)
	}
	if len(a.args) > 0 {
		lines[len(lines)-1] = " └─"
		a.args[len(a.args)-1].Print(lines)
	}
}

func (a Annotation) ResolveNames(p *parser.Parser) bool {
	return true // TODO
}

func (a Annotation) FindStartToken() scan.Token {
	return a.id.FindStartToken()
}