package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	nodetype "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type List Sequence

var listRule = nodetype.NodeRule{
	Production: nodetype.LIST,
	Expression: []nodetype.NodeType{nodetype.SEQUENCE},
}

func (ls List) Make(p *Parser) bool {
	if valid, e := p.Stack.Validate(listRule); !valid {
		e.Print()
		return false
	}
	tmp := p.Stack.Pop().(Sequence)
	p.Stack.Push(List(tmp))
	return true
}
func (ls List) GetNodeType() nodetype.NodeType {
	return nodetype.LIST
}
func (ls List) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == nodetype.LIST
	if !equal {
		return false
	}
	ls2 := a.(List)
	for i := range ls {
		if !ls[i].Equal_test(ls2[i]) {
			return false
		}
	}
	return true
}
func (ls List) Print(lns []string) {
	lines := make([]string, len(lns))
	lines = append(lines, lns...)
	lines = printLines(lines)
	fmt.Printf("List\n")
	lines = append(lines, " ├─")
	for i := 0; i < len(ls)-1; i++ {
		ls[i].Print(lines)
	}

	lines[len(lines)-1] = " └─"
	if len(ls) > 0 {
		ls[len(ls)-1].Print(lines)
	} else {
		printLines(lines)
		fmt.Printf("[]\n")
	}
}
func (ls List) ResolveNames(table *symbol.SymbolTable) bool {
	ok := true
	for _, l := range ls {
		if ok = l.ResolveNames(table); !ok {
			break
		}
	}
	return ok
}
func (ls List) ExpressionType() types.Types {
	if len(ls) == 0 {
		return types.Array{ElemType: types.GetNewTau()}
	}
	return types.Array{ElemType: ls[0].ExpressionType()}
}
func (ls List) DoTypeInference(newTypeInformation types.Types) types.Types {
	if newTypeInformation.GetTypeType() != types.ARRAY {
		return types.Error(types.TypeErrors[types.E_EXPECTED_ARRAY]().(err.Error))
	}
	elemType := newTypeInformation.(types.Array).ElemType
	for _, elem := range ls {
		e := elem.DoTypeInference(elemType)
		if e.GetTypeType() == types.ERROR {
			return e
		}
	}
	return newTypeInformation
}

func (ls List) FindStartToken() scan.Token {
	return ls[0].FindStartToken()
}
