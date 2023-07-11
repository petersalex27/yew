package ast

import (
	"fmt"
	"os"
	"unicode"
	err "yew/error"
	scan "yew/lex"

	//"yew/parser/ast"
	errorgen "yew/parser/error-gen"
	nodetype "yew/parser/node-type"
	"yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type Type struct {
	ty types.Types
}

type Constructor struct {
	parent types.Data
	self types.Constructor
	idToken scan.IdToken
}

func RegisterConstructors(p *parser.Parser, dat types.Data) bool {
	for name, constructor := range dat.Constructors {
		cons := Constructor{
			self: constructor, 
			idToken: scan.MakeIdToken(name, constructor.Loc.GetLine(), constructor.Loc.GetChar()),
		}
		e, ok := p.Table.DeclareLocal(cons, dat)
		if !ok {
			e.ToError().Print()
			return false
		}
	}
	return true
}

func (c Constructor) GetIdToken() scan.IdToken {
	return c.idToken
}
func (c Constructor) GetType() types.Types {
	return c.self
}
func (c Constructor) SetType(dat types.Types) symbol.Symbolic {
	if dat.GetTypeType() != types.DATA {
		err.PrintBug()
		panic("")
	}
	c.parent = dat.(types.Data)
	return c
}
func (c Constructor) IsDefined() bool {
	return true 
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
		id := p.Stack.Pop().(Id)
		tau := types.MakeTau(id.GetName(), scan.ToLoc(id.token))
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
			id := right.(Id)
			tau := types.MakeTau(id.GetName(), scan.ToLoc(id.token))
			out = append(types.Application{tau}, out...)
		}
	}

	if left.GetNodeType() != nodetype.IDENTIFIER {
		out = append(types.Application{types.Error{}}, out...) // should be caught later
	} else {
		id := left.(Id)
		tau := types.MakeTau(id.GetName(), scan.ToLoc(id.token))
		out = append(types.Application{tau}, out...)
	}
	ty := MakeType(out)
	p.Stack.Push(ty)
	return true
}

func (t Type) MakeApplication(p *parser.Parser) bool {
	valid, e := p.Stack.Validate(typeAppRule)
	if !valid {
		e(p.Input).Print()
		return false
	}

	tyTail := p.Stack.Pop().(Type)
	tyHead := p.Stack.Pop().(Type)
	if tyHead.ty.GetTypeType() == types.APPLICATION {
		head := tyHead.ty.(types.Application)
		app := make(types.Application, 0, len(head) + 1)
		app = append(app, head...)
		app = append(app, tyTail.ty)
		p.Stack.Push(Type{app})
	} else {
		app := types.Application{tyHead.ty, tyTail.ty}
		p.Stack.Push(Type{app})
	}
	return true
}

func (t Type) MakeData(p *parser.Parser) bool {
	/*ok := applicationToTauApplication(p)
	if !ok {
		return false
	}*/

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

func GrabConstructorName(from types.Types) (types.Constructor, errorgen.GenerateErrorFunction) {
	if from.GetTypeType() != types.TAU {
		return types.Constructor{}, errorgen.UnexpectedType.Generate()
	}
	name := from.(types.Tau)
	if !unicode.IsUpper(rune(name.ToString()[0])) {
		return types.Constructor{}, errorgen.ExpectedTypeIdentifierNotVar.Generate()
	}
	return types.Constructor{Name: name.ToString(), Members: make(types.Application, 0)}, nil
}

func ToConstructor(from types.Types) (types.Constructor, errorgen.GenerateErrorFunction) {
	tt := from.GetTypeType()
	if tt == types.TAU {
		return GrabConstructorName(from)
	} else if tt == types.APPLICATION {
		head, tail := from.(types.Application).Split()
		c, e := GrabConstructorName(head)
		if e != nil {
			return c, e
		}

		if tail.GetTypeType() == types.APPLICATION {
			c.Members = tail.(types.Application)
		} else {
			c.Members = types.Application{tail}
		}
		return c, nil
	}

	return types.Constructor{}, errorgen.UnexpectedType.Generate()
}

func (t Type) AddConstructor(p *parser.Parser) bool {
	return makeBinaryType(p, func(left types.Types, right types.Types) types.Types {
		if left.GetTypeType() != types.DATA {
			return types.Error(types.TypeErrors[types.E_UNEXPECTED]().(err.Error))
		}
		dat := left.(types.Data)
		cons, e := ToConstructor(right)
		if nil != e {
			tok := scan.MakeIdToken("", right.GetLocation().GetLine(), right.GetLocation().GetChar())
			return types.Error(e(tok, p.Input))
		}
		consName := cons.Name
		_, found := dat.Constructors[consName]
		if found {
			tok := scan.MakeIdToken("", right.GetLocation().GetLine(), right.GetLocation().GetChar())
			return types.Error(errorgen.RedeclaredConstructor.Generate()(tok, p.Input))
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

func (t Type) ResolveNames(p *parser.Parser) bool { return true /* TODO: ??? is this right */ }

func (t Type) FindStartToken()scan.Token {
	return scan.MakeBlankToken()
}

