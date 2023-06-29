package ast

import (
	"fmt"
	scan "yew/lex"
	errorgen "yew/parser/error-gen"
	nodetype "yew/parser/node-type"
	"yew/parser/parser"
	"yew/symbol"
	err "yew/error"
	types "yew/type"
)

type Pattern struct {
	Expression Expression
	Matchers []Lambda
}

var statementError = errorgen.GenerateSyntaxError("unexpected statement inside pattern")
var expressionError = errorgen.GenerateSyntaxError("expected anonymous function")
var emptyPatternError = errorgen.GenerateSyntaxError("cannot have an empty pattern")
var expectedFunctionType = errorgen.GenerateTypeError("expected a function type")
var failedTypeInference = errorgen.GenerateTypeError("could not inference type") // TODO: say why!!

func printStatementErrors(p *parser.Parser, statements []Statement) {
	for _, statement := range statements {
		tok := statement.GetSymbol().GetIdToken()
		statementError(tok, p.Input).Print()
	}
}

func makeFromSequence(p *parser.Parser, patStartToken scan.Token, seq Sequence) ([]Lambda, bool) {
	if seq == nil || len(seq) == 0 {
		emptyPatternError(patStartToken, p.Input).Print()
		return []Lambda{}, false
	}

	mat := make([]Lambda, len(seq))
	for i, s := range seq {
		if s.GetNodeType() != nodetype.LAMBDA {
			expressionError(s.FindStartToken(), p.Input).Print()
			return []Lambda{}, false
		}
		
		mat[i] = s.(Lambda)
	}
	return mat, true
}

// two type classes are valid for an annotation: (1) function type, (2) tau type.
func handlePatternTypeAnnot(p *parser.Parser, a ExpressionTypeAnnotation) (Lambda, bool) {
	// In the case of (1), just split the function annotation and wrap the binder
	// in an annotation of the domain type and wrap the bound in an annotation of the
	// codomain type. In the case of (2), create two new taus (call them t1 and t2), 
	// add a new rule for the tau from the original annotation for (t1 -> t2); then, 
	// like case (1), wrap the lambda's parts in an annotation
	if a.expression.GetNodeType() != nodetype.LAMBDA {
		expressionError(a.FindStartToken(), p.Input).Print()
		return Lambda{}, false
	}

	lambda := a.expression.(Lambda)

	var functionType types.Function
	if a.expressionType.GetTypeType() == types.FUNCTION {
		// case (1)
		functionType = a.expressionType.(types.Function)
	} else if a.expressionType.GetTypeType() == types.TAU {
		// case (2)
		// type from annotation
		tau := a.expressionType.(types.Tau)

		// create new functionType
		taus := types.GetNewTaus(2)
		functionType.Domain = taus[0]
		functionType.Codomain = taus[1]
		
		// add new rule: tau = functionType
		typeType := types.DoTypeInference(tau, functionType).GetTypeType()
		if typeType == types.ERROR { // this shouldn't happen
			err.PrintBug()
			return Lambda{}, false
		}
	} else {
		expectedFunctionType(a.FindStartToken(), p.Input).Print()
		return Lambda{}, false
	}

	ty := types.DoTypeInference(lambda.binder.pattern.expressionType, functionType.Domain)
	lambda.binder.pattern.expressionType = ty
	if ty.GetTypeType() == types.ERROR {
		failedTypeInference(lambda.binder.FindStartToken(), p.Input).Print()
		return Lambda{}, false
	}

	var annot ExpressionTypeAnnotation
	if lambda.bound.GetNodeType() == nodetype.TYPE_ANNOTATION {
		annot = lambda.bound.(ExpressionTypeAnnotation)
		ty := types.DoTypeInference(annot.expressionType, functionType.Codomain)
		annot.expressionType = ty
		if ty.GetTypeType() == types.ERROR {
			failedTypeInference(lambda.bound.FindStartToken(), p.Input).Print()
			return Lambda{}, false
		}
	} else {
		annot.expression = lambda.bound
		annot.expressionType = functionType.Codomain
	}

	lambda.bound = annot
	return lambda, true
}

func makePatternFromAstArray(p *parser.Parser, patStartToken scan.Token, as []parser.Ast) ([]Lambda, bool) {
	if len(as) == 0 {
		emptyPatternError(patStartToken, p.Input).Print()
		return []Lambda{}, false
	}

	matches := make([]Lambda, len(as))
	for i, a := range as {
		var lam Lambda
		nodeType := a.GetNodeType()
		if nodeType != nodetype.LAMBDA {
			if nodeType != nodetype.TYPE_ANNOTATION {
				expressionError(a.FindStartToken(), p.Input).Print()
				return []Lambda{}, false
			}
			annot := a.(ExpressionTypeAnnotation)
			lamTmp, ok := handlePatternTypeAnnot(p, annot)
			if !ok {
				return matches, false
			}
			lam = lamTmp
		} else {
			lam = a.(Lambda)
		}
		matches[i] = lam
	}
	return matches, true
}

func toAstArr(seq Sequence) []parser.Ast {
	as := make([]parser.Ast, len(seq))
	for i := range seq {
		as[i] = seq[i]
	}
	return as
}

func (pat Pattern) MakePattern(p *parser.Parser, patStartToken scan.Token) bool {
	var arr []parser.Ast

	// pattern ::= expression lambda | expression sequence | expression program
	if valid, _ := p.Stack.TryValidate(patternRule.Expression); valid {
		lam := p.Stack.Pop().(Lambda)
		expr := p.Stack.Pop().(Expression)
		pat.Expression = expr
		pat.Matchers = []Lambda{lam}
	} else { // pattern ::= expression sequence | expression program
		if valid, _ := p.Stack.TryValidate(patternRule2.Expression); valid {
			// patterm ::= expression sequence
			seq := p.Stack.Pop().(Sequence)
			arr = toAstArr(seq)
		} else {
			// pattern ::= expression program
			valid, e := p.Stack.Validate(patternRule3)
			if !valid {
				e(p.Input).Print()
				return false
			}

			arr = p.Stack.Pop().(Program)
		}

		expr := p.Stack.Pop().(Expression)

		mat, ok := makePatternFromAstArray(p, patStartToken, arr)
		if !ok {
			return false
		}
		pat = Pattern{Expression: expr, Matchers: mat}
	}

	p.Stack.Push(pat)
	return true
}

func (pat Pattern) Make(p *parser.Parser) bool {
	err.PrintBug()
	panic("")
}
func (pat Pattern) GetNodeType() nodetype.NodeType {
	return nodetype.PATTERN
}
func (pat Pattern) Equal_test(a parser.Ast) bool {
	if a.GetNodeType() != nodetype.PATTERN {
		return false
	}
	pat2, ok := a.(Pattern)
	if !ok {
		return false
	}
	for i := range pat.Matchers {
		if !pat.Matchers[i].Equal_test(pat2.Matchers[i]) {
			return false
		}
	}
	return true
}

func (pat Pattern) Print(lines []string) {
	copy := make([]string, 0, len(lines))
	copy = append(copy, lines...)
	copy = printLines(copy)
	fmt.Printf("Pattern\n")
	copy = append(copy, " ├─")
	pat.Expression.Print(copy)
	for i := 0; i < len(pat.Matchers)-1; i++ {
		pat.Matchers[i].Print(copy)
	}
	if len(pat.Matchers) > 0 {
		copy[len(copy)-1] = " └─"
		pat.Matchers[len(pat.Matchers)-1].Print(copy)
	}
}
func (pat Pattern) ResolveNames(table *symbol.SymbolTable) bool {
	// TODO
	panic("TODO")
}
func (pat Pattern) ExpressionType() types.Types {
	if nil == pat.Matchers || len(pat.Matchers) == 0 {
		return types.Tuple{} // empty type
	}
	return pat.Matchers[0].ExpressionType()
}
func (pat Pattern) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("TODO") // TODO
}
func (pat Pattern) FindStartToken() scan.Token {
	return pat.Matchers[0].FindStartToken()
}
