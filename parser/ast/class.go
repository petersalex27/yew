package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	errorgen "yew/parser/error-gen"
	nodetype "yew/parser/node-type"
	"yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type Class struct {
	name Id
	typeParameter types.Tau
	functions map[string]types.Function
}

func (c Class) GetClassName() string {
	return c.name.GetName()
}

type classEntry struct{
	class Class
	instances map[string]map[string]Function
}
type ClassTable struct {
	table *map[string]classEntry
}
func (ClassTable) InitClassTable() parser.ClassTable_ {
	tab := ClassTable{
		table: new(map[string]classEntry),
	}
	*tab.table = make(map[string]classEntry)
	return tab
}

func (tab ClassTable) Lookup(className string) (class parser.Class_, found bool) {
	classEntry, ok := (*tab.table)[className]
	found = ok
	if ok {
		class = classEntry.class
	}

	return
}

func (tab ClassTable) getClass(className string) (entry classEntry, errorFn errorgen.GenerateErrorFunction) {
	var found bool
	entry, found = (*tab.table)[className]
	if !found {
		errorFn = errorgen.UndefinedClass.Generate() // class not defined
	}
	return
}

func (tab ClassTable) GetClass(p *parser.Parser, class parser.Class_) (entry classEntry, found bool) {
	var errorFn errorgen.GenerateErrorFunction
	entry, errorFn = tab.getClass(class.GetClassName())
	if errorFn != nil {
		errorFn(class.(Class).FindStartToken(), p.Input).Print()
		found = false
	} else {
		found = true
	}
	return
}

func (entry classEntry) checkUninstantiated(instanceString string) errorgen.GenerateErrorFunction {
	_, found := entry.instances[instanceString]
	if found {
		// instance of class already exists
		return errorgen.RedeclaredClassInstance.Generate()
	}
	return nil
}

func (entry classEntry) CheckUninstantiated(p *parser.Parser, instanceString string, inst types.Types) bool {
	errFn := entry.checkUninstantiated(instanceString)
	if nil != errFn {
		// instance of class already exists
		loc := inst.GetLocation()
		dummyToken := scan.MakeOtherToken(scan.TYPE_ID, loc.GetLine(), loc.GetChar())
		errFn(dummyToken, p.Input).Print()
		return false
	}
	return true
}

func (entry classEntry) confirmFunctionDeclared(fn Function) errorgen.GenerateErrorFunction {
	_, found := entry.class.functions[fn.dec.GetName()]
	if !found {
		return errorgen.FunctionNotInClass.Generate()
	}
	return nil
}

func (entry classEntry) ConfirmFunctionDeclared(p *parser.Parser, fn Function) bool {
	errFn := entry.confirmFunctionDeclared(fn)
	if errFn != nil {
		errFn(fn.FindStartToken(), p.Input).Print()
		return false
	}
	return true
}

func confirmUniqueDefinition(instances map[string]Function, fn Function) errorgen.GenerateErrorFunction {
	_, found := instances[fn.dec.GetName()]
	if found {
		return errorgen.FunctionInstanceRedefined.Generate()
	}
	return nil
}

func ConfirmUniqueDefinition(p *parser.Parser, instances map[string]Function, fn Function) bool {
	errFn := confirmUniqueDefinition(instances, fn)
	if errFn != nil {
		errFn(fn.FindStartToken(), p.Input).Print()
		return false
	}
	return true
}

func (tab ClassTable) DeclareClass(p *parser.Parser, newClass parser.Class_) bool {
	class := newClass.(Class)
	_, found := tab.Lookup(class.name.GetName())
	if found {
		// class redeclared
		errorgen.RedeclaredClass.
			Generate()(class.FindStartToken(), p.Input).Print()
		return false
	}

	(*tab.table)[class.GetClassName()] = classEntry{
		class: class,
		instances: make(map[string]map[string]Function),
	}
	return true
}
func (tab ClassTable) DeclareInstance(p *parser.Parser, class parser.Class_, instance types.Types) bool {
	className := class.GetClassName()
	
	// search for class
	entry, found := tab.GetClass(p, class)
	if !found {
		return false
	}

	// confirm that class has not been instantiated yet for the given type 
	instanceString := instance.ToString()
	uninst := entry.CheckUninstantiated(p, instanceString, instance)
	if !uninst {
		return false
	}

	// create map for new instance
	newInstance := make(map[string]Function, len(entry.class.functions))
	entry.instances[instanceString] = newInstance

	// add map
	(*tab.table)[className] = entry

	return true
}
func (tab ClassTable) DefineInstanceFunction(p *parser.Parser, class parser.Class_, instance types.Types, function parser.InstanceFunction_) bool {
	className := class.GetClassName()
	// search for class
	entry, found := tab.GetClass(p, class)
	if !found {
		return false
	}

	// confirm that class instantce has been declared for the given type 
	instanceString := instance.ToString()
	classInstance, foundInst := entry.instances[instanceString]
	if !foundInst {
		// this should never happen
		err.PrintBug()
		panic("")
	}	

	fn := function.(Function) // should always be true
	fnName := fn.dec.GetName()

	// confirm function is declared in class being instantiated
	found = entry.ConfirmFunctionDeclared(p, fn)
	if !found {
		return false
	}

	if !ConfirmUniqueDefinition(p, classInstance, fn) {
		return false
	}

	// add function
	classInstance[fnName] = fn
	entry.instances[instanceString] = classInstance
	(*tab.table)[className] = entry
	return true
}

func (c Class) SetTypeParameter(param types.Tau) Class {
	c.typeParameter = param
	return c
}

func MakeClass(name Id, typeParameter types.Tau, fns map[string]types.Function) Class {
	return Class{
		name: name,
		typeParameter: typeParameter,
		functions: fns,
	}
}

func InitClass(name Id) Class {
	return Class{name: name, functions: make(map[string]types.Function)}
}

func (c Class) GetSymbol() symbol.Symbolic {
	return symbol.MakeSymbol(c.name.token)
}

func (c Class) GetIdToken() scan.IdToken {
	return c.name.token
}

func (c Class) GetType() types.Types {
	return types.Class{
		Loc: scan.ToLoc(c.name.FindStartToken()),
		Name: c.name.GetName(),
		TypeVariable: c.typeParameter,
		Functions: c.functions,
	}
}

func (c Class) SetType(ty types.Types) symbol.Symbolic {
	if ty.GetTypeType() != types.CLASS {
		err.PrintBug()
		panic("")
	}
	cFns := ty.(types.Class).Functions
	c.functions = cFns
	return c
}
func (c Class) IsDefined() bool {
	return true
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

func (c Class) ResolveNames(p *parser.Parser) bool {
	panic("TODO") // TODO
}

func (c Class) FindStartToken() scan.Token {
	return c.name.token
}