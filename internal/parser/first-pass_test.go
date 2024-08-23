package parser

import "testing"

func Test_CompileTime_ImplementSyntacticElem(*testing.T) {
	// tests won't compile if rhs is not a SyntacticElem
	var _ SyntacticElem = DeclarationElem{}
	var _ SyntacticElem = TypeElem{}
	var _ SyntacticElem = BindingElem{}
	var _ SyntacticElem = WhereClause{}
	var _ SyntacticElem = LetBindingElem{}
	var _ SyntacticElem = TokensElem{}
	var _ SyntacticElem = TypeConstructorElem{}
	var _ SyntacticElem = DataTypeElem{}
	var _ SyntacticElem = SpecElem{}
	var _ SyntacticElem = InstanceElem{}
	var _ SyntacticElem = AnnotationElem{}
}
