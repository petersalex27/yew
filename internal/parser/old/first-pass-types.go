package parser

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew/internal/token"
)

type (
	SyntacticElem interface {
		Parse(*Parser) (ok bool)
		Pos() (start, end int)
		String() string
	}

	SymbolicElem interface {
		SyntacticElem
		GetSyntaxClass() SyntaxClass
	}

	ClauseElem interface {
		SyntacticElem
		_c_()
	}

	ExpressionElem interface {
		SyntacticElem
		parseExpressionPart(*Parser, *actionData) (end, ok bool)
		_e_()
	}
)

type Visibility uint8

const (
	Private Visibility = 0b00 // default visibility
	Public  Visibility = 0b01
	Open    Visibility = 0b10
)

type ModuleElement struct {
	Visibility
	Declaration DeclarationElem
}

func (v Visibility) String() string {
	switch v {
	case Public:
		return "public"
	case Open:
		return "open"
	default:
		return ""
	}
}

type (
	ModuleElem struct {
		Name       token.Token
		Imports    ImportTable
		Elems      []SyntacticElem
		Start, End int
	}

	TypeElem struct {
		Constraint []token.Token
		Type       []token.Token
		Start, End int
	}

	declarationFlags uint8

	localDefinition struct {
		Expression []ExpressionElem
		Clause     ClauseElem
		Start, End int
	}

	DeclarationElem struct {
		Vis    Visibility
		flags  declarationFlags
		Name   token.Token
		rAssoc bool
		bp     int8
		Typing TypeElem
		// non-nil if binding is the special case of let-binding that allows both declaration and
		// associated binding in same language construct
		//
		// Example:
		//	thisIsThree : 3 = 3
		//	threeIsThree = let x : 3 = 3 := Refl in x
		*localDefinition
		Start, End int
	}

	BindingElem struct {
		Head  []token.Token
		Start int
		localDefinition
		// End = localDefinition.End
	}

	RefinedBinding struct {
		BindingElem
		NewView []token.Token
	}

	// Example:
	// 	f x y with g x of (
	//    f x 0 | 1 = x
	//    _ | _ = y
	//	)
	ViewBinding struct {
		Head  []token.Token
		View  []token.Token
		Start int
		Arms  []RefinedBinding
	}

	WhereClause struct {
		// should only have declarations, bindings, and data types
		Elems      []SymbolicElem
		Start, End int
	}

	MutualBlockElem struct {
		// should only have declarations, bindings, and data types
		Traits       []SpecElem
		Types        []DataTypeElem
		Instances    []InstanceElem
		Declarations []DeclarationElem
		Bindings     []BindingElem
		atTop        bool
		Start, End   int
	}

	// just the let binding part (doesn't include 'in' and onwards). Example:
	//	let x = 3
	// or
	//	let
	//	  x = 3
	//	  y = x + 1
	LetBindingElem WhereClause

	TypeConstructorElem struct {
		DeclarationElem
		//Name   token.Token
		//Args   []token.Token
		//Typing TypeElem
		// Start = Name.Start
		// End = Typing.End
	}

	// a type constructor paired with zero or more data constructors for that type
	//
	//	List : * -> Uint -> * where
	//	  Nil : List a 0
	//	  Cons : a -> List a n -> List a (n+1)
	//
	//	Nat : * where
	//	  Zero : Nat
	//	  Succ : Nat -> Nat
	//
	// `Nat` could also be defined as ...
	//
	//	Nat : * where
	//	  Zero
	//	  Succ Nat
	DataTypeElem struct {
		TypeConstructor  TypeConstructorElem
		DataConstructors []DeclarationElem
		End              int
		// Start = TypeConstructor.Name.Start
	}

	HaskellStyleDataConstructors struct {
		Vis        Visibility
		Name       token.Token
		Params     []TypeElem
		Start, End int
	}

	HaskellStyleDataTypeElem struct {
		TypeConstructorElem TypeConstructorElem
		DataConstructors    []HaskellStyleDataConstructors
		End                 int
	}

	// `trait` { constraint } IDENT { IDENT } `where` decl { decl }
	SpecElem struct {
		// { constraint } IDENT { IDENT }
		Constraint         []token.Token
		Head               []token.Token
		MethodDeclarations []DeclarationElem
		Start, End         int
	}

	// IDENT `of` { type } `where` def { def }
	InstanceElem struct {
		TraitName   token.Token
		Instance    []token.Token
		Definitions []BindingElem
		End         int
		// Start = TraitName.Start
	}

	TokensElem []token.Token

	ExtensionElem struct {
		// TODO: implement
		Start, End int
	}

	AnnotationElem struct {
		Name token.Token
		Args []token.Token
		End  int
	}
)

const (
	// flags for DeclarationElem
	implicitDecl declarationFlags = 1 << iota
	isPure
	noInline
	suggestInline
	isBuiltin
	isExternal
	isDataConstructor
	isTest
)

func (TokensElem) _e_()     {}
func (LetBindingElem) _e_() {}

func (WhereClause) _c_()     {}
func (MutualBlockElem) _c_() {}

func (m ModuleElem) Pos() (start, end int) {
	return m.Start, m.End
}

func (m ModuleElem) String() string {
	return fmt.Sprintf("module %s", m.Name.Value)
}

func (s ExtensionElem) Pos() (start, end int) {
	return s.Start, s.End
}

func (s ExtensionElem) String() string {
	return "_extension_" // TODO
}

func (decl DeclarationElem) GetSyntaxClass() SyntaxClass {
	return DeclClass
}

func (DataTypeElem) GetSyntaxClass() SyntaxClass {
	return TypeClass
}

func (SpecElem) GetSyntaxClass() SyntaxClass {
	return TraitClass
}

func (InstanceElem) GetSyntaxClass() SyntaxClass {
	return InstanceClass
}

func (BindingElem) GetSyntaxClass() SyntaxClass {
	return FuncClass
}

func (TypeConstructorElem) GetSyntaxClass() SyntaxClass {
	return TypeConsClass
}

func (ExtensionElem) GetSyntaxClass() SyntaxClass {
	return ExtensionClass
}

func (MutualBlockElem) GetSyntaxClass() SyntaxClass {
	return MutualClass
}

func joinTokensStringed(tokens []token.Token, sep string) string {
	var b strings.Builder
	switch len(tokens) {
	case 0:
		return ""
	}

	b.WriteString(tokens[0].Value)
	for _, tok := range tokens[1:] {
		b.WriteString(sep)
		b.WriteString(tok.Value)
	}
	return b.String()
}

func joinElemsStringed[T SyntacticElem](elems []T, sep string) string {
	var b strings.Builder
	switch len(elems) {
	case 0:
		return ""
	}

	b.WriteString(elems[0].String())
	for _, elem := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(elem.String())
	}
	return b.String()
}

func (dec DeclarationElem) Pos() (start, end int) {
	return dec.Start, dec.End
}

func (dec DeclarationElem) String() string {
	var vis string
	if dec.Vis != Private {
		vis = dec.Vis.String() + " "
	}
	start := fmt.Sprintf("%s%s : %s", vis, dec.Name.Value, dec.Typing.String())
	if dec.localDefinition != nil {
		return start + dec.localDefinition.String()
	}
	return start
}

func (typ HaskellStyleDataConstructors) Pos() (start, end int) {
	return typ.Start, typ.End
}

func (typ HaskellStyleDataConstructors) String() string {
	if len(typ.Params) > 0 {
		return fmt.Sprintf("%s %s %s : Type", typ.Vis.String(), typ.Name.Value, joinElemsStringed(typ.Params, " "))
	}
	return fmt.Sprintf("%s %s : Type", typ.Vis.String(), typ.Name.Value)
}

func (typ HaskellStyleDataTypeElem) Pos() (start, end int) {
	return typ.TypeConstructorElem.Start, typ.End
}

func (typ HaskellStyleDataTypeElem) String() string {
	return fmt.Sprintf("%s where {%s}", typ.TypeConstructorElem.String(), joinElemsStringed(typ.DataConstructors, "; "))
}

func (typ TypeElem) Pos() (start, end int) {
	return typ.Start, typ.End
}

func (typ TypeElem) String() string {
	if len(typ.Constraint) > 0 {
		return fmt.Sprintf("%s %s", joinTokensStringed(typ.Constraint, " "), joinTokensStringed(typ.Type, " "))
	}
	return joinTokensStringed(typ.Type, " ")
}

func (local localDefinition) Pos() (start, end int) {
	return local.Start, local.End
}

func (local localDefinition) String() string {
	if len(local.Expression) == 0 {
		return ""
	}
	start := fmt.Sprintf(" := %s", joinElemsStringed(local.Expression, " "))
	clause := local.Clause.String()
	if len(clause) != 0 {
		return start + " " + clause
	}
	return start
}

func (bind BindingElem) Pos() (start, end int) {
	return bind.Start, bind.End
}

func (bind BindingElem) String() string {
	start := fmt.Sprintf("%s = %s", joinTokensStringed(bind.Head, " "), joinElemsStringed(bind.Expression, " "))
	clause := bind.Clause.String()
	if len(clause) != 0 {
		return start + " " + clause
	}
	return start
}

func (mut MutualBlockElem) Pos() (start, end int) {
	return mut.Start, mut.End
}

func (mut MutualBlockElem) String() string {
	traits := joinElemsStringed(mut.Traits, "; ")
	types := joinElemsStringed(mut.Types, "; ")
	insts := joinElemsStringed(mut.Instances, "; ")
	decls := joinElemsStringed(mut.Declarations, "; ")
	bindings := joinElemsStringed(mut.Bindings, "; ")
	return fmt.Sprintf("mutual {%s} {%s} {%s} {%s} {%s}", traits, types, insts, decls, bindings)
}

func (where WhereClause) Pos() (start, end int) {
	return where.Start, where.End
}

func (where WhereClause) String() string {
	if len(where.Elems) == 0 {
		return ""
	}
	return fmt.Sprintf("where %s", joinElemsStringed(where.Elems, "; "))
}

func (let LetBindingElem) Pos() (start, end int) {
	return let.Start, let.End
}

func (let LetBindingElem) String() string {
	return fmt.Sprintf("let %s in", joinElemsStringed(let.Elems, "; "))
}

func (tok TokensElem) Pos() (start, end int) {
	if len(tok) == 0 {
		return 0, 0
	}
	return tok[0].Start, tok[len(tok)-1].End
}

func (tok TokensElem) String() string {
	return joinTokensStringed(tok, " ")
}

func (cons TypeConstructorElem) Pos() (start, end int) {
	return cons.DeclarationElem.Pos()
}

func (cons TypeConstructorElem) String() string {
	// if len(cons.Args) > 0 {
	// 	return fmt.Sprintf("%s %s : %s", cons.Name.Value, joinTokensStringed(cons.Args, " "), cons.Typing.String())
	// }
	return cons.DeclarationElem.String()
	//return fmt.Sprintf("%s : %s", cons.Name.Value, cons.Typing.String())
}

func (data DataTypeElem) Pos() (start, end int) {
	return data.TypeConstructor.Name.Start, data.End
}

func (data DataTypeElem) String() string {
	return fmt.Sprintf("%s where {%s}", data.TypeConstructor.String(), joinElemsStringed(data.DataConstructors, "; "))
}

func (trait SpecElem) Pos() (start, end int) {
	return trait.Start, trait.End
}

func (trait SpecElem) String() string {
	var start string
	if len(trait.Constraint) > 0 {
		start = fmt.Sprintf("trait %s %s where", joinTokensStringed(trait.Constraint, " "), joinTokensStringed(trait.Head, " "))
	} else {
		start = fmt.Sprintf("trait %s where", joinTokensStringed(trait.Head, " "))
	}
	return fmt.Sprintf("%s {%s}", start, joinElemsStringed(trait.MethodDeclarations, "; "))
}

func (inst InstanceElem) Pos() (start, end int) {
	return inst.TraitName.Start, inst.End
}

func (inst InstanceElem) String() string {
	return fmt.Sprintf("%s of %s where {%s}", inst.TraitName, joinTokensStringed(inst.Instance, " "), joinElemsStringed(inst.Definitions, "; "))
}

func (a AnnotationElem) Pos() (start int, end int) {
	return a.Name.Start, a.End
}

func (a AnnotationElem) String() string {
	if a.Name.Value == "" {
		return ""
	}

	if len(a.Args) == 0 {
		return "%" + a.Name.Value
	}
	args := joinTokensStringed(a.Args, " ")
	// %annotation tok0 tok1 .. tokN
	return fmt.Sprintf("%%%s%s", a.Name.Value, args)
}
