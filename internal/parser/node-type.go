package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/internal/common"
	t "github.com/petersalex27/yew/internal/parser/typ"
)

func (n access) Type() api.NodeType               { return t.Access }
func (n annotations) Type() api.NodeType          { return t.Annotations }
func (n appType) Type() api.NodeType              { return t.AppType }
func (n body) Type() api.NodeType                 { return t.Body }
func (n enclosedAnnotation) Type() api.NodeType   { return t.EnclosedAnnotation }
func (n caseArm) Type() api.NodeType              { return t.CaseArm }
func (n caseArms) Type() api.NodeType             { return t.CaseArms }
func (n caseExpr) Type() api.NodeType             { return t.CaseExpr }
func (n constrainedType) Type() api.NodeType      { return t.ConstrainedType }
func (n constrainer) Type() api.NodeType          { return t.Constrainer }
func (n constraintUnverified) Type() api.NodeType { return t.Constraint }
func (n constraintVerified) Type() api.NodeType   { return t.Constraint }
func (n def) Type() api.NodeType                  { return t.Def }
func (n defBody) Type() api.NodeType              { return n.Either.Type() }
func (n defBodyPossible) Type() api.NodeType      { return t.DefBody }
func (n defaultExpr) Type() api.NodeType          { return t.DefaultExpr }
func (n deriving) Type() api.NodeType             { return t.Deriving }
func (n derivingBody) Type() api.NodeType         { return t.DerivingBody }
func (n enclosedType) Type() api.NodeType {
	return ifThenElse(n.implicit, t.ImplicitType, t.EnclosedType)
}
func (n exprApp) Type() api.NodeType           { return t.ExprApp }
func (n flatAnnotation) Type() api.NodeType    { return t.FlatAnnotation }
func (n footer) Type() api.NodeType            { return t.Footer }
func (n forallBinders) Type() api.NodeType     { return t.ForallBinders }
func (n forallType) Type() api.NodeType        { return t.ForallType }
func (n functionType) Type() api.NodeType      { return t.FunctionType }
func (n header) Type() api.NodeType            { return t.Header }
func (n hole) Type() api.NodeType              { return t.Hole }
func (n implicitTyping) Type() api.NodeType    { return t.ImplicitTyping }
func (n importStatement) Type() api.NodeType   { return t.ImportStatement }
func (n importing) Type() api.NodeType         { return t.Importing }
func (n importPathIdent) Type() api.NodeType   { return t.ImportPathIdent }
func (n impossible) Type() api.NodeType        { return t.Impossible }
func (n innerTypeTerms) Type() api.NodeType    { return t.InnerTypeTerms }
func (n innerTyping) Type() api.NodeType       { return t.InnerTyping }
func (n lambdaAbstraction) Type() api.NodeType { return t.LambdaAbstraction }
func (n lambdaBinders) Type() api.NodeType     { return t.LambdaBinders }
func (n letBinding) Type() api.NodeType        { return t.LetBinding }
func (n letExpr) Type() api.NodeType           { return t.LetExpr }
func (n literal) Type() api.NodeType           { return t.Literal }
func (n lowerIdent) Type() api.NodeType        { return t.LowerIdent }

//func (n meta) Type() api.NodeType              { return t.Meta }

func (n modality) Type() api.NodeType          { return t.Modality }
func (n module) Type() api.NodeType            { return t.Module }
func (n name) Type() api.NodeType {
	// try to get most specific type for name
	chs := n.Solo.Children()
	if len(chs) != 1 {
		return t.Name
	}
	c := chs[0].(api.Token)
	if common.Is_camelCase2(c) {
		return t.LowerIdent
	} else if common.Is_PascalCase2(c) {
		return t.UpperIdent
	}
	return t.Name
}
func (n packageImport) Type() api.NodeType { return t.PackageImport }
func (n patternApp) Type() api.NodeType    { return t.PatternApp }
func (n patternEnclosed) Type() api.NodeType {
	return ifThenElse(n.implicit, t.PatternImplicitArg, t.PatternEnclosed)
}
func (n rawString) Type() api.NodeType           { return t.RawString }
func (n specBody) Type() api.NodeType            { return t.SpecBody }
func (n specHead) Type() api.NodeType            { return t.SpecHead }
func (n specDef) Type() api.NodeType             { return t.SpecDef }
func (n specInst) Type() api.NodeType            { return t.SpecInst }
func (n syntax) Type() api.NodeType              { return t.Syntax }
func (n syntaxRule) Type() api.NodeType          { return t.SyntaxRule }
func (n syntaxRuleIdent) Type() api.NodeType     { return t.SyntaxRuleIdent }
func (n typeAlias) Type() api.NodeType           { return t.TypeAlias }
func (n typeConstructor) Type() api.NodeType     { return t.TypeConstructor }
func (n typeDef) Type() api.NodeType             { return t.TypeDef }
func (n typing) Type() api.NodeType              { return t.Typing }
func (n unitType) Type() api.NodeType            { return t.UnitType }
func (n upperIdent) Type() api.NodeType          { return t.UpperIdent }
func (n visibility) Type() api.NodeType          { return t.Visibility }
func (n whereBody) Type() api.NodeType           { return t.WhereBody }
func (n whereClause) Type() api.NodeType         { return t.WhereClause }
func (n wildcard) Type() api.NodeType            { return t.Wildcard }
func (n withClause) Type() api.NodeType          { return t.WithClause }
func (n withClauseArm) Type() api.NodeType       { return t.WithClauseArm }
func (n withClauseArms) Type() api.NodeType      { return t.WithClauseArms }
func (n yewSource) Type() api.NodeType           { return t.YewSource }

func (n annotations) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n appType) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n body) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n enclosedAnnotation) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n caseArm) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n caseArms) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n caseExpr) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n constrainedType) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n constrainer) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n constraintUnverified) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n constraintVerified) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n def) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{n.annotations, n.pattern, n.defBody}
}
func (n defBody) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n defBodyPossible) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n defaultExpr) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n deriving) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n derivingBody) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n enclosedType) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{n.typ}
}
func (n exprApp) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n flatAnnotation) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n footer) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n forallBinders) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n forallType) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n functionType) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n header) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n hole) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n implicitTyping) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n importStatement) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n importing) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n importPathIdent) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n impossible) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n innerTypeTerms) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n innerTyping) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{n.mode, n.typing}
}
func (n lambdaAbstraction) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n lambdaBinders) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n letBinding) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n letExpr) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n literal) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n lowerIdent) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}

// func (n meta) Describe() (string, []api.Node) {
// 	return n.Type().String(), n.Children()
// }

func (n modality) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n module) Describe() (string, []api.Node) {
	return n.Type().String(), append([]api.Node{n.annotations}, n.name.Children()...)
}
func (n name) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n packageImport) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n patternApp) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n patternEnclosed) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n rawString) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n specBody) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n specHead) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n specDef) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{n.annotations, n.visibility, n.specHead, n.dependency, n.specBody, n.requiring}
}
func (n specInst) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{n.annotations, n.visibility, n.head, n.target, n.body}
}
func (n syntax) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{n.annotations, n.visibility, n.rule}
}
func (n syntaxRule) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n syntaxRuleIdent) Describe() (string, []api.Node) {
	s := "binding "
	if !n.binding {
		s = ""
	}
	return s + n.Type().String(), n.id.Children()
}
func (n typeAlias) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{n.annotations, n.visibility, n.alias}
}
func (n typeConstructor) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{n.annotations, n.constructor}
}
func (n typeDef) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{n.annotations, n.visibility, n.typedef, n.deriving}
}
func (n typing) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{n.annotations, n.visibility, n.typing}
}
func (n unitType) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n upperIdent) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n visibility) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n whereBody) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n whereClause) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n wildcard) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n withClause) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n withClauseArm) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n withClauseArms) Describe() (string, []api.Node) {
	return n.Type().String(), n.Children()
}
func (n yewSource) Describe() (string, []api.Node) {
	return n.Type().String(), []api.Node{ /*n.meta,*/ n.header, n.body, n.footer}
}
