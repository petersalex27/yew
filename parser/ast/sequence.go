package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
	types "yew/type"
	util "yew/utils"
)

type Sequence []Expression

func (s Sequence) ResolveNames(table *symbol.SymbolTable) bool {
	for _, z := range s {
		if !z.ResolveNames(table) {
			return false
		}
	}
	return true
}
func (s Sequence) GetNodeType() NodeType { return SEQUENCE }
func (s Sequence) Make(p *Parser) bool {
	// convert any statements into empty expressions
	valid, _ := p.Stack.TryValidate([]NodeType{STATEMENT})
	if valid {
		if !(Program{}.Make(p)) { // makes a statement with an empty expression
			return false
		}
		// programs are expressions
	}

	valid, _ = p.Stack.TryValidate([]NodeType{SEQUENCE, EXPRESSION})
	if valid {
		// continue a sequence
		ex := p.Stack.Pop().(Expression)
		seq := p.Stack.Pop().(Sequence)
		seq = append(seq, ex)
		p.Stack.Push(seq)
		return true
	}

	var e err.Error
	valid, e = p.Stack.Validate(NodeRule{SEQUENCE, []NodeType{EXPRESSION}})
	if !valid {
		e.Print()
		return false
	}
	// create new sequence
	ex := p.Stack.Pop().(Expression)
	p.Stack.Push(append(make(Sequence, 0, 1), ex))
	return true
}
func (s Sequence) ExpressionType() types.Types {
	e, found := util.Tail(s)
	if !found {
		return types.Tuple{} // empty tuple, i.e., ()
	}
	return e.ExpressionType()
}
func (s Sequence) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("TODO: implement") // TODO
}

func (s Sequence) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == SEQUENCE
	s2 := a.(Sequence)
	equal = equal && len(s2) == len(s)
	if !equal {
		return false
	}

	for i, z := range s2 {
		equal = equal && s[i].Equal_test(z)
	}
	return equal
}
func (s Sequence) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("Sequence\n")
	lines = append(lines, " ├─")
	for i := 0; i < len(s)-1; i++ {
		s[i].Print(lines)
	}
	if len(s) > 0 {
		lines[len(lines)-1] = " └─"
		s[len(s)-1].Print(lines)
	}
}

func (s Sequence) FindStartToken() scan.Token {
	return s[0].FindStartToken()
}
