package ast

import (
	"fmt"
	"os"
	err "yew/error"
	scan "yew/lex"

	//"yew/parser/ast"
	nodetype "yew/parser/node-type"
	"yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type Type struct {
	ty types.Types
}

func MakeType(ty types.Types) Type {
	return Type{ty}
}

func (t Type) GetType() types.Types {
	return t.ty
}

func makeBinaryType(p *parser.Parser, createType func (types.Types, types.Types) types.Types) bool {
	valid, e := p.Stack.Validate(binaryTypeRule)
	if !valid {
		e(p.Input).Print()
		return false
	}
	right := p.Stack.Pop().(Type).ty
	left := p.Stack.Pop().(Type).ty
	ty := createType(left, right)
	if ty.GetTypeType() == types.ERROR {
		if ty.(types.Error) != (types.Error{}) {
			ty.(types.Error).ToError().Print()
		}
		return false
	}
	p.Stack.Push(Type{ty})
	return true
}

func (t Type) MakeFunction(p *parser.Parser) bool {
	return makeBinaryType(p, func(left types.Types, right types.Types) types.Types {
		return types.Function{Domain: left, Codomain: right}
	})
}

func (t Type) MakeTuple(p *parser.Parser) bool {
	return makeBinaryType(p, func(left types.Types, right types.Types) types.Types {
		if left.GetTypeType() == types.TUPLE {
			return types.Tuple(append(left.(types.Tuple), right))
		}
		return types.Tuple{left, right}
	})
}

func applicationToTauApplication(p *parser.Parser) bool {
	ok, _ := p.Stack.TryValidate([]nodetype.NodeType{nodetype.IDENTIFIER})
	if ok {
		tau := types.Tau(p.Stack.Pop().(Id).GetName())
		p.Stack.Push(MakeType(tau))
		//println("here", string(tau))
		return true
	}

	valid, e := p.Stack.Validate(nodetype.NodeRule{
		Production: nodetype.TYPE,
		Expression: []nodetype.NodeType{nodetype.APPLICATION},
	})
	if !valid {
		e(p.Input).Print()
		return false
	}
	app := p.Stack.Pop().(Application)
	var left Expression
	left = app
	out := make(types.Application, 0)
	for {
		var right Expression
		if left.GetNodeType() != nodetype.APPLICATION {
			break
		}
		left, right = left.(Application).split()
		if right.GetNodeType() != nodetype.IDENTIFIER {
			out = append(types.Application{types.Error{}}, out...) // should be caught later
		} else {
			out = append(types.Application{types.Tau(right.(Id).GetName())}, out...)
		}
	}

	if left.GetNodeType() != nodetype.IDENTIFIER {
		out = append(types.Application{types.Error{}}, out...) // should be caught later
	} else {
		out = append(types.Application{types.Tau(left.(Id).GetName())}, out...)
	}
	ty := MakeType(out)
	p.Stack.Push(ty)
	return true
}

func (t Type) MakeData(p *parser.Parser) bool {
	ok := applicationToTauApplication(p)
	if !ok {
		return false
	}

	valid, e := p.Stack.Validate(justTypeRule)
	if !valid {
		e(p.Input).Print()
		return false
	}

	dat, ok := types.ToData(p.Stack.Pop().(Type).ty)
	if !ok {
		return false
	}
	p.Stack.Push(Type{dat})
	return true
}

func (t Type) AddConstructor(p *parser.Parser) bool {
	return makeBinaryType(p, func(left types.Types, right types.Types) types.Types {
		if left.GetTypeType() != types.DATA {
			return types.Error(types.TypeErrors[types.E_UNEXPECTED]().(err.Error))
		}
		dat := left.(types.Data)
		cons, ok := types.ToConstructor(right)
		if !ok {
			return types.Error{}
		}
		consName := cons.Name
		_, found := dat.Constructors[consName]
		if found {
			panic("TODO: print better error--redeclared constructor\n")
		}
		dat.Constructors[consName] = cons
		return dat
	})
}

func (t Type) MakeArray(p *parser.Parser) bool {
	valid, e := p.Stack.Validate(justTypeRule)
	if !valid {
		e(p.Input).Print()
		return false
	}
	ty := p.Stack.Pop().(Type)
	ty.ty = types.Array{ElemType: ty.ty}
	p.Stack.Push(ty)
	return true
}

func (t Type) Make(*parser.Parser) bool {
	fmt.Fprintf(os.Stderr, 
		"use Make<What>(*parser.Parser) instead of Make(*parser.Parser)\n")
	err.PrintBug()
	panic("")
}
func (t Type) GetNodeType() nodetype.NodeType {
	return nodetype.TYPE
}
func (t Type) Equal_test(a parser.Ast) bool {
	ok := a.GetNodeType() == nodetype.TYPE
	if !ok {
		return false
	}
	return t.ty.Equals(a.(Type).ty)
}
func (t Type) Print(lines []string) {
	printLines(lines)
	fmt.Printf("Type == %s\n", t.ty.ToString())
}
func (t Type) ResolveNames(table *symbol.SymbolTable) bool {
	return true // TODO?
}

func (t Type) FindStartToken()scan.Token {
	return scan.MakeBlankToken()
}

