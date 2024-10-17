package typ

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
)

type nodeType int32

const (
	// shift to avoid overlap with token types (range=[0x0, 0xff]) and node types (range=[0x100, 0xfff])
	_start_ nodeType = (1 + iota) << 0x10
	EmptyType_
	PairType_
	ListType_
	NonEmptyType_

	// end of control types
	_last_
)

// ASSUMPTION: USER_DEFINED_START must not be 0 (this is a tested assumption, though; see the corresponding test file in the token package)

const MinProperNodeTypeValue nodeType = nodeType(token.USER_DEFINED_START)
const MaxProperNodeTypeValue nodeType = 0xff

// actual types

const (
	Access nodeType = MinProperNodeTypeValue + iota
	Annotations
	AppType
	Body
	BodyElem
	CaseArm
	CaseArms
	CaseExpr
	Char
	ConstrainedType
	Constrainer
	Constraint
	Def
	DefBody
	DefaultExpr
	Deriving
	EnclosedAnnotation
	EnclosedType
	Error
	ExprApp
	FlatAnnotation
	Footer
	ForallBinders
	ForallType
	FunctionType
	Header
	Hole
	ImplicitType
	ImplicitTyping
	ImportPathIdent
	ImportStatement
	Importing
	Impossible
	InnerTypeTerms
	InnerTyping
	LambdaAbstraction
	LambdaBinders
	LetBinding
	LetExpr
	Literal
	LowerIdent
	MainElem
	Meta
	Module
	Modality
	Name
	PackageImport
	PatternApp
	PatternEnclosed
	PatternImplicitArg
	RawString
	SpecBody
	SpecDef
	SpecHead
	SpecInst
	Syntax
	SyntaxRawKeyword
	SyntaxRule
	SyntaxRuleIdent
	TypeAlias
	TypeConstructor
	TypeDef
	Typing
	UnitType
	UpperIdent
	Visibility
	WhereClause
	Wildcard
	WithClause
	WithClauseArm
	WithClauseArms
	_last_node_type_ // last node type =======================================
	YewSource        = _last_node_type_
	UnknownNodeType  = _last_node_type_ + 1 // unknown marker =================
)

func (nt nodeType) Match(n api.Node) bool {
	return nt.String() == n.Type().String()
}

func (nt nodeType) String() string {
	switch nt {
	case UnknownNodeType:
		return "?unknown"
	case Access:
		return "access"
	case Error:
		return "error"
	case Annotations:
		return "annotations"
	case AppType:
		return "type application"
	case Body:
		return "body"
	case BodyElem:
		return "body element"
	case CaseArm:
		return "case arm"
	case CaseArms:
		return "case arms"
	case CaseExpr:
		return "case expression"
	case Char:
		return "char literal"
	case ConstrainedType:
		return "constrained type"
	case Constrainer:
		return "constrainer"
	case Constraint:
		return "constraint"
	case Def:
		return "definition"
	case DefBody:
		return "definition body"
	case DefaultExpr:
		return "default expression"
	case Deriving:
		return "deriving clause"
	case EnclosedAnnotation:
		return "enclosed annotation"
	case EnclosedType:
		return "enclosed type"
	case ExprApp:
		return "expression application"
	case FlatAnnotation:
		return "flat annotation"
	case Footer:
		return "footer section"
	case ForallBinders:
		return "forall binders"
	case ForallType:
		return "forall type"
	case FunctionType:
		return "function type"
	case Header:
		return "header section"
	case Hole:
		return "hole"
	case ImplicitType:
		return "implicit type"
	case ImplicitTyping:
		return "implicit typing"
	case ImportStatement:
		return "import statement"
	case Importing:
		return "imports"
	case Impossible:
		return "impossible definition body"
	case InnerTypeTerms:
		return "inner type terms"
	case InnerTyping:
		return "inner typing"
	case LambdaAbstraction:
		return "lambda abstraction"
	case LambdaBinders:
		return "lambda binders"
	case LetBinding:
		return "let bindings"
	case LetExpr:
		return "let expression"
	case Literal:
		return "literal"
	case LowerIdent:
		return "lower identifier"
	case MainElem:
		return "main element"
	case Meta:
		return "meta section"
	case Module:
		return "module"
	case Modality:
		return "modality"
	case Name:
		return "name"
	case ImportPathIdent:
		return "import path identifier"
	case PackageImport:
		return "package import"
	case PatternApp:
		return "pattern application"
	case PatternEnclosed:
		return "enclosed pattern"
	case PatternImplicitArg:
		return "implicit argument pattern"
	case RawString:
		return "raw string literal"
	case SpecBody:
		return "spec body"
	case SpecDef:
		return "spec definition"
	case SpecHead:
		return "spec head"
	case SpecInst:
		return "inst definition"
	case Syntax:
		return "syntax definition"
	case SyntaxRawKeyword:
		return "syntax keyword"
	case SyntaxRule:
		return "syntax rule"
	case SyntaxRuleIdent:
		return "syntax rule ident"
	case TypeAlias:
		return "type alias definition"
	case TypeConstructor:
		return "type constructor"
	case TypeDef:
		return "type definition"
	case Typing:
		return "typing"
	case UnitType:
		return "unit type"
	case UpperIdent:
		return "upper identifier"
	case Visibility:
		return "visibility"
	case WhereClause:
		return "where clause"
	case Wildcard:
		return "wildcard"
	case WithClause:
		return "with clause"
	case WithClauseArm:
		return "with clause arm"
	case WithClauseArms:
		return "with clause arms"
	case YewSource:
		return "yew source"
	}
	return UnknownNodeType.String()
}
