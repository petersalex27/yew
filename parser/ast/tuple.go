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

type Tuple Sequence

var tupleRule = nodetype.NodeRule{
	Production: nodetype.TUPLE,
	Expression: []nodetype.NodeType{nodetype.SEQUENCE},
}

func (t Tuple) Make(p *parser.Parser) bool {
	if valid, e := p.Stack.Validate(tupleRule); !valid {
		e.Print()
		return false
	}
	tmp := p.Stack.Pop().(Sequence)
	p.Stack.Push(Tuple(tmp))
	return true
}
func (t Tuple) GetNodeType() nodetype.NodeType {
	return nodetype.TUPLE
}
func (t Tuple) Equal_test(a parser.Ast) bool {
	equal := a.GetNodeType() == nodetype.TUPLE
	if !equal {
		return false
	}
	t2 := a.(Tuple)
	for i := range t {
		if !t[i].Equal_test(t2[i]) {
			return false
		}
	}
	return true
}
func (t Tuple) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("Tuple\n")
	lines = append(lines, " ├─")
	for i := 0; i < len(t)-1; i++ {
		t[i].Print(lines)
	}

	lines[len(lines)-1] = " └─"
	if len(t) > 0 {
		t[len(t)-1].Print(lines)
	} else {
		printLines(lines)
		fmt.Printf("()\n")
	}
}
func (ls Tuple) ResolveNames(table *symbol.SymbolTable) bool {
	ok := true
	for _, l := range ls {
		if ok = l.ResolveNames(table); !ok {
			break
		}
	}
	return ok
}
func (t Tuple) ExpressionType() types.Types {
	ty := make(types.Tuple, len(t))
	for i, elem := range t {
		ty[i] = elem.ExpressionType()
	}
	return ty
}
func (t Tuple) DoTypeInference(newTypeInformation types.Types) types.Types {
	if newTypeInformation.GetTypeType() != types.TUPLE {
		return types.Error(types.TypeErrors[types.E_EXPECTED_TUPLE]().(err.Error))
	}
	tupType := newTypeInformation.(types.Tuple)
	if len(tupType) != len(t) {
		return types.Error(types.TypeErrors[types.E_UNEXPECTED]().(err.Error))
	}
	for i := range t {
		e := t[i].DoTypeInference(tupType[i])
		if e.GetTypeType() == types.ERROR {
			return e
		}
	}
	return newTypeInformation
}

func (t Tuple) FindStartToken() scan.Token {
	return t[0].FindStartToken()
}
