package ast

import (
	"fmt"
	"os"
	"strconv"
	err "yew/error"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

// gives a name to an anonymous function
type Function struct {
	dec      Id     // identifies function
	function Lambda // actual function
	//res Expression
}

func (f Function) GetSymbol() symbol.Symbolic {
	return f
}

func (f Function) GetIdToken() scan.IdToken {
	return f.dec.token
}
func (f Function) GetType() types.Types {
	return f.function.ExpressionType()
}
func (f Function) SetType(t types.Types) symbol.Symbolic {
	f.dec.ty = t
	return f
}
func (f Function) IsDefined() bool {
	return true
}

type StackMarker struct{}

func (StackMarker) ResolveNames(*symbol.SymbolTable) bool {
	err.PrintBug()
	panic("")
}
func (StackMarker) Make(*Parser) bool {
	err.PrintBug()
	panic("")
}
func (StackMarker) GetNodeType() NodeType {
	return STACK_MARKER
}
func (StackMarker) Equal_test(Ast) bool {
	err.PrintBug()
	panic("")
}
func (StackMarker) Print([]string) {
	err.PrintBug()
	panic("")
}

func (a Application) split() (Expression, Expression) {
	return a.left, a.right
}

func unrollFunctionType(f types.Function, i int) []types.Types {
	if f.Codomain.GetTypeType() != types.FUNCTION {
		tmp := make([]types.Types, i+2)
		tmp[i] = f.Domain
		tmp[i+1] = f.Codomain
		return tmp
	} else {
		tmp := unrollFunctionType(f.Codomain.(types.Function), i+1)
		tmp[i] = f.Domain
		return tmp
	}
}

func declareParam(table *symbol.SymbolTable, i int, expr Expression) (Parameter, bool) {
	var ty types.Types
	var id Id

	if TYPE_ANNOTATION == expr.GetNodeType() {
		ty = expr.(ExpressionTypeAnnotation).expressionType
		expr = expr.(ExpressionTypeAnnotation).expression
	} else {
		ty = types.GetNewTau()
	}

	if IDENTIFIER == expr.GetNodeType() {
		id = expr.(Id)
	} else {
		id = MakeId(scan.MakeIdToken("_p"+strconv.Itoa(i), 0, 0))
	}

	e, decd := table.DeclareLocal(symbol.MakeSymbol(id.token), ty)
	if !decd {
		e.ToError().Print()
		return Parameter{}, false
	}
	return Parameter{paramIndex: i, pattern: ExpressionTypeAnnotation{
		expression:     expr,
		expressionType: ty,
	}}, true
}
func buildAnnotated(p *Parser, i int) bool {
	app := p.Stack.Pop().(Application)
	lam := p.Stack.Pop().(Lambda)
	anot := p.Stack.Pop().(ExpressionTypeAnnotation)
	app2, right := app.split()
	anot.expression = right
	param, ok := declareParam(p.Table, i, anot)
	if !ok {
		return false
	}
	p.Stack.Push(Lambda{binder: param, bound: lam})
	p.Stack.Push(app2)
	return true
}
func buildUnannotated(p *Parser, i int) bool {
	app := p.Stack.Pop().(Application)
	lam := p.Stack.Pop().(Lambda)
	app2, right := app.split()
	param, ok := declareParam(p.Table, i, right)
	if !ok {
		return false
	}
	p.Stack.Push(Lambda{binder: param, bound: lam})
	p.Stack.Push(app2)
	return true
}
var progressLambdas = []NodeRule{
	{NodeType(IN_PROGRESS__ | LAMBDA), /* ::= */ []NodeType{LAMBDA, APPLICATION}},
	{NodeType(IN_PROGRESS__ | LAMBDA), /* ::= */ []NodeType{TYPE_ANNOTATION, LAMBDA, APPLICATION}},
}

func validateUnannotated(stack *AstStack) bool {
	valid, _ := stack.Validate(progressLambdas[0])
	return valid
}
func validateAnnotated(stack *AstStack) bool {
	valid, _ := stack.Validate(progressLambdas[1])
	return valid
}

func createInitial(
	parseFunctionBody func(*Parser) bool,
	validNodeTypes []NodeType,
	updateRight func(*Parser, int, Expression) (Parameter, bool),
) func(p *Parser) bool {
	return func(p *Parser) bool {
		rule := NodeRule{
			Production: NodeType(IN_PROGRESS__ | FUNCTION), 
			Expression: validNodeTypes,
		}
		if valid, e := p.Stack.Validate(rule); !valid {
			e.Print()
			return false
		}
		app, right := p.Stack.Pop().(Application).split()
		retType := p.Stack.Pop().(ExpressionTypeAnnotation)
		param, ok := updateRight(p, 0, right)
		if !ok {
			return false
		}
		p.Stack.Push(param)
		ok = parseFunctionBody(p)
		if !ok {
			return false
		}
		retType.expression = p.Stack.Pop().(Expression)
		p.Stack.Push(retType)
		ok = Lambda{}.Make(p) // create function from last param to body
		if !ok {
			return false
		}
		p.Stack.Push(app)
		return true
	}
}

type unrollAction struct {
	initial func(*Parser) bool
	build func(*Parser, int) bool
	validate func(*AstStack) bool
}

func buildFunctionDeclaration(p *Parser, action unrollAction, start int) bool {
	for i := start; action.validate(p.Stack); i++ {
		if !action.build(p, i) {
			return false
		}
	}
	return true
}

var functionRule2 = NodeRule{
	Production: IN_PROGRESS__ | FUNCTION, /* ::= */ 
	Expression: []NodeType{LAMBDA, IDENTIFIER},
}
func unrollApplication(p *Parser, action unrollAction, functionName Id, fnType types.Types) bool {
	// first, find function name and add it to the symbol table (might be used 
	// 	inside function definition, so it needs to be declared prior to 
	// 	resolving names inside the function definition)
	e, decd := p.Table.DeclareLocal(symbol.MakeSymbol(functionName.token), fnType)
	if !decd {
		e.ToError().Print()
		return false
	}

	p.Table.AddScope(symbol.NewScopeTable())

	var ok bool = false
	var fn Function
	for { // does not loop! this is just here so `break` can be used
		ok = action.initial(p) && 
				buildFunctionDeclaration(p, action, 1)
		if !ok {
			break
		}
		valid, e :=	p.Stack.Validate(functionRule2)
		if !valid {
			ok = valid
			e.Print()
			break
		}
		dec := p.Stack.Pop().(Id) // this should match functionName
		body := p.Stack.Pop().(Lambda)
		fn = Function{dec: dec, function: body}
		p.Stack.Push(fn)
		break
	}

	p.Table.RemoveScope()

	// define function (previously declared)
	if ok {
		e, ok = p.Table.DefineSymbol(fn) // define function
		if !ok {
			e.ToError().Print()
		}
	}
	return ok
}

// breaks apart function's type annotation 
func pushFunctionTypeAnnotation(stack *AstStack, tyAnnot types.Function) {
	if tyAnnot.Codomain.GetTypeType() == types.FUNCTION {
		stack.Push(ExpressionTypeAnnotation{expressionType: tyAnnot.Domain})
		pushFunctionTypeAnnotation(stack, tyAnnot.Codomain.(types.Function))
	} else {
		stack.Push(ExpressionTypeAnnotation{expressionType: tyAnnot.Domain})
		stack.Push(ExpressionTypeAnnotation{expressionType: tyAnnot.Codomain})
	}
}

func DeclareFunction(p *Parser, functionName Id, parseFunctionBody func(*Parser) bool) bool {
	if valid, _ := p.Stack.TryValidate([]NodeType{TYPE_ANNOTATION}); valid {
		annot := p.Stack.Pop().(ExpressionTypeAnnotation)
		if annot.expression.GetNodeType() != APPLICATION {
			fmt.Fprintf(os.Stderr, "Error: TODO--expected APPLICATION\n")
			return false // TODO: error message
		}
		if annot.expressionType.GetTypeType() != types.FUNCTION {
			fmt.Fprintf(os.Stderr, "Error: TODO--expected FUNCTION\n")
			return false // TODO: error message
		}
		tyAnnot := annot.expressionType.(types.Function)
		pushFunctionTypeAnnotation(p.Stack, tyAnnot)
		app := annot.expression.(Application)
		p.Stack.Push(app)
		initialAnnotated := createInitial(
			parseFunctionBody,
			[]NodeType{TYPE_ANNOTATION, APPLICATION},
			func(p *Parser, i int, e Expression) (Parameter, bool) {
				ex := p.Stack.Pop().(ExpressionTypeAnnotation)
				ex.expression = e
				return declareParam(p.Table, i, ex)
			})
		action := unrollAction{initialAnnotated, buildAnnotated, validateAnnotated}
		return unrollApplication(p, action, functionName, tyAnnot)
	} else if valid, _ := p.Stack.TryValidate([]NodeType{APPLICATION}); valid {
		ty := ExpressionTypeAnnotation{
			expression: EmptyExpression{}, 
			expressionType: types.GetNewTau(), // return type
		}
		tmp := p.Stack.Pop()
		p.Stack.Push(ty)
		p.Stack.Push(tmp)
		initialUnannotated := createInitial(
			parseFunctionBody,
			[]NodeType{TYPE_ANNOTATION, APPLICATION},
			func(p *Parser, i int, e Expression) (Parameter, bool) {
				return declareParam(p.Table, i, e)
			})
		action := unrollAction{initialUnannotated, buildUnannotated, validateUnannotated}
		return unrollApplication(p, action, functionName, types.GetNewTau())
	} else {
		fmt.Fprintf(os.Stderr, "Error: TODO--could not declare function\n")
		return false // TODO: error message
	}
}

func (f Function) ResolveNames(table *symbol.SymbolTable) bool {
	return f.dec.ResolveNames(table) && f.function.ResolveNames(table)
}

func MakeFunction(dec Id, function Lambda) Function {
	return Function{dec: dec, function: function}
}

func (f Function) GetNodeType() NodeType { return FUNCTION }

func (f Function) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == FUNCTION
	f2 := a.(Function)
	return equal &&
		f.dec.Equal_test(f2.dec) &&
		f.function.Equal_test(f2.function)
}

func (f Function) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("Function\n")
	lines = append(lines, " ├─")
	f.dec.Print(lines)
	lines[len(lines)-1] = " └─"
	f.function.Print(lines)
}

// Function ::= Declaration Anonymous-Function
var functionRule = NodeRule{
	FUNCTION, /* ::= */ []NodeType{DECLARATION, LAMBDA},
}

func (f Function) Make(p *Parser) bool {
	valid, e := p.Stack.Validate(functionRule)
	if !valid {
		e.Print()
		return false
	}

	lambda := p.Stack.Pop().(Lambda)
	dec := p.Stack.Pop().(Declaration)
	f.dec = MakeId(dec.token)
	f.function = lambda
	p.Stack.Push(f)
	return true
}

func (f Function) StackLogString() string {
	return fmt.Sprintf("%s; %s", f.GetNodeType().ToString(), Id(f.dec).token.ToString())
}
