// =================================================================================================
// Alex Peters - March 04, 2024
// =================================================================================================
package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/source"
	"github.com/petersalex27/yew/token"
	"github.com/petersalex27/yew/types"
)

func stringsAndSort(m map[token.Type]bool) []string {
	s := make([]string, 0, len(m))
	for ty := range m {
		if v := ty.Make().Value; v != "" {
			s = append(s, ty.Make().Value)
		}
	}
	sort.Strings(s)
	return s
}

func expectedMessageMulti(m map[token.Type]bool, defaultMsg ...string) (msg string) {
	def := ""
	if len(defaultMsg) > 0 {
		def = defaultMsg[0]
	}

	elems := stringsAndSort(m)
	if len(elems) == 0 {
		if def != "" {
			return def
		}
		panic("expectedMessageMulti: no default message provided, assuming bug in caller")
	}

	// surround each element with single quotes
	msg = "'" + strings.Join(elems, "', '") + "'"
	// if there is more than one element, add an 'or' before the last element
	if len(elems) > 1 {
		idx := strings.LastIndex(msg, "'")
		msg = msg[:idx] + " or " + msg[idx:]
	}
	// more than one token:
	//		"expected to find token 'tok1', 'tok2', or 'tok3'"
	// one token:
	//		"expected to find token 'tok1'"
	return ExpectedToFindToken_ + msg
}

func expectedNodeMessage(ty NodeType) string {
	switch ty {
	case identType:
		return ExpectedIdentifier
	case intConstType:
		return ExpectedIntLit
	case charConstType:
		return ExpectedCharLit
	case floatConstType:
		return ExpectedFloatLit
	case stringConstType:
		return ExpectedStringLit
	case functionType:
		return ExpectedFunctionType
	case applicationType:
		return ExpectedApplication
	case lambdaType:
		return ExpectedLambdaAbstraction
	case typingType:
		return ExpectedTyping
	case listExprType:
		return ExpectedListExpr
	case tupleType:
		return ExpectedTuple
	case tupleExprType:
		return ExpectedTupleExpr
	case pairsType:
		return ExpectedPairings
	case listingType:
		return ExpectedListing
	default:
		return UnexpectedSection
	}
}

func expectedSyntax(expected string) string {
	return fmt.Sprintf("expected '%s'", expected)
}

func expectedMessage(ty token.Type) string {
	switch ty {
	case token.Equal:
		return ExpectedEqual
	case token.Id:
		return ExpectedIdentifier
	case token.In:
		return ExpectedIn
	case token.Module:
		return ExpectedModule
	case token.RightBrace:
		return ExpectedRBrace
	case token.RightParen:
		return ExpectedRParen
	case token.RightBracket:
		return ExpectedRBracket
	case token.ThickArrow:
		return ExpectedThickArrow
	case token.Where:
		return ExpectedWhere
	case token.LeftParen:
		return ExpectedLParen
	default:
		return UnexpectedToken
	}
}

const (
	ArityMismatch string = "new definition-case's arity doesn't match previous definition-cases' arities"

	DuplicateImportName string = "duplicate import name"

	// = "expected" ==================================================================================
	ExpectedBinding               string = "expected pattern binding"
	ExpectedEqual                 string = "expected assignment"
	ExpectedCommaOrThickArrow     string = "expected ',' or '=>'"
	ExpectedIdentifier            string = "expected identifier"
	ExpectedLeftRightNone         string = "expected constant 'Left', 'Right', or 'None'"
	ExpectedPascalCase            string = "expected PascalCase identifier"
	ExpectedIn                    string = "expected 'in'"
	ExpectedIndentation           string = "expected indentation"
	ExpectedModule                string = "expected 'module'"
	ExpectedName                  string = "expected name"
	ExpectedRBrace                string = "expected '}'"
	ExpectedRBracket              string = "expected ']'"
	ExpectedRParen                string = "expected ')'"
	ExpectedThickArrow            string = "expected '=>'"
	ExpectedType                  string = "expected type"
	ExpectedTypeApplication       string = "expected type application"
	ExpectedTyping                string = "expected typing"
	ExpectedVariable              string = "expected variable"
	ExpectedWhere                 string = "expected 'where'"
	ExpectedMutual                string = "expected 'mutual'"
	ExpectedDeclaration           string = "expected declaration"
	ExpectedScrutinee             string = "expected scrutinee"
	ExpectedExpression            string = "expected expression"
	ExpectedLRN                   string = "expected one of the constants `Left`, `Right`, or `None`"
	ExpectedInteger               string = "expected an integer literal"
	ExpectedInteger0to10          string = "expected an integer literal in the range [0, 10]"
	ExpectedGreaterIndent         string = "expected larger indentation than enclosing context's indentation"
	ExpectedLParen                string = "expected '('"
	ExpectedUint                  string = "expected unsigned integer"
	ExpectedUintRange1_9          string = "expected unsigned integer in the range 1 .. 9"
	ExpectedConstraint            string = "expected constraint, e.g.,\n\t'(Trait a1, .., Trait aN) => ..', or\n\t'Trait a => ..'"
	ExpectedConstrainedType       string = "expected constrained type, e.g., \n\t'Trait a'"
	ExpectedWildcard              string = "expected wildcard"
	ExpectedAffixedId             string = "expected affixed identifier"
	ExpectedIntLit                string = "expected int literal"
	ExpectedCharLit               string = "expected char literal"
	ExpectedFloatLit              string = "expected float literal"
	ExpectedStringLit             string = "expected string literal"
	ExpectedFunctionType          string = "expected function type"
	ExpectedApplication           string = "expected function application"
	ExpectedLambdaAbstraction     string = "expected lambda abstraction"
	ExpectedListType              string = "expected list type"
	ExpectedListExpr              string = "expected list expression"
	ExpectedTuple                 string = "expected tuple"
	ExpectedTupleExpr             string = "expected tuple expression"
	ExpectedPairings              string = "expected pairings"
	ExpectedListing               string = "expected listing"
	ExpectedSyntaxExtension       string = "expected syntax extension"
	ExpectedTermAfter             string = "expected another syntactic sub-section following this one"
	ExpectedEndOfSection          string = "unexpected continued section"
	ExpectedExpressionAfterLet    string = "expected expression after binding in let-binding"
	FunctionPatternExpected       string = "expected function pattern"
	FunctionDefinedButNotDeclared string = "function defined but not declared"
	ExpectedComma                 string = "expected ','"

	ExpectedToFindToken_ string = "expected to find token "

	BuiltinExternalConflict             string = `'%builtin' and '%external' are mutually exclusive`
	InlineNotAllowed                    string = `'%inline' and '%noInline' are mutually exclusive`
	IllegalNoInlineSuggestion           string = `'%noInline' and '%inline' are mutually exclusive`
	UnverifiablePureAnnotation          string = `external functions cannot be verified as pure`
	UnverifiablePureAnnotation_RESOURCE string = `functions that use the 'RESOURCE' type, or have types that use it, cannot be verified as pure`

	InvalidBuiltinAnnotation string = " '%builtin' annotation"

	KindMismatch string = "kind mismatch"

	// = "illegal" ===================================================================================
	IllegalApplication         string = "illegal application"
	IllegalEnclosedTypeFamily  string = "illegal type family, type family enclosed by type family"
	IllegalImplicitDeclaration string = "illegal declaration, implicit declaration of affixed identifier"
	IllegalImport              string = "illegal import, empty"
	IllegalPattern             string = "illegal pattern, expected initial identifier"
	IllegalRedeclaration       string = "illegal redeclaration"
	IllegalReimport            string = "illegal import, multiply declared module"
	IllegalTraitNoParams       string = "illegal trait, no type parameters"
	IllegalTraitRedeclaration  string = "illegal trait redeclaration"
	IllegalTupleType           string = "illegal tuple type, must be an n-tuple where n > 1"
	IllegalTypeRedeclaration   string = "illegal type redeclaration"
	IllegalUse                 string = "illegal use-import, empty"
	IllegalWhere               string = "illegal 'where'"
	IllegalVarExtElem          string = "extension cannot have multiple contiguous variables"
	IllegalNonExprPosHole      string = "hole outside of expression position"
	IllegalDeclaration         string = "illegal location of declaration"
	IllegalBinding             string = "illegal location of binding"
	IllegalTrait               string = "illegal location of trait definition occurrence"
	IllegalInstance            string = "illegal location of instance definition"
	IllegalDataType            string = "illegal location of data type definition"
	IllegalDataTypeName        string = "illegal data type identifier, must be more than a character followed by numbers and/or single quotes"
	IllegalConstraintPosition  string = "illegal constraint position, must appear before all constrained types"
	IllegalExplicatedArgInType string = "illegal explicated argument, cannot appear in type"
	IllegalModalityLocation    string = "illegal modality location, must appear immediately to the left of a typing"
	IllegalTypeConsList        string = "type constructors cannot be declared in list form"
	IllegalTypeConsLocalDef    string = "type constructors cannot be defined locally"
	IllegalInfixTyping         string = "illegal infix typing"

	FunctionDefinitionSplit string = "function definition is split by other definitions; original definition(s) here:"

	IncompleteFunctionType string = "incomplete function type"
	IncompleteLambda       string = "incomplete lambda abstraction"

	UnexpectedModality string = "unexpected modality"

	MultiplyOccurringAffixedIdent string = "multiple occurrences of affixed form of identifier"

	EmptyMutualBlock string = "empty mutual block"

	MalformedAffixAnnotation string = "malformed affix annotation"
	BadIdent                 string = "malformed identifier"

	NoGrammarExtensionFound string = "no applicable grammar extension found"

	NoAppropriateType string = "no appropriate type found"

	CouldNotParseInt    string = "couldn't parse integer"
	CouldNotParseFloat  string = "couldn't parse float"
	CouldNotParseChar   string = "couldn't parse char"
	CouldNotParseString string = "couldn't parse string"

	CannotReturnNamedArg string = "cannot return a named argument"
	CannotReturnImplicitArg string = "cannot return an implicit argument"

	// = "invalid" ===================================================================================
	InvalidListElementType      string = "invalid list type, no element type"
	InvalidTypeIdentifier       string = "invalid type identifier"
	InvalidDuplicateMutualBlock string = "mutual blocks cannot be nested"

	RequireCamelCase        string = "camelCase identifier is required"
	RequirePascalCaseModule string = "module identifiers must be in PascalCase"

	// TODO: better message??
	UnusedContext    string = "unreferenced context"
	UnusedVisibility string = "unused visibility modifier"

	// = "unexpected" ================================================================================
	UnexpectedEOF      string = "unexpected end of file"
	UnexpectedFinalTok string = "unexpected final token in syntactic structure"
	UnexpectedToken    string = "unexpected token"
	UnexpectedIndent   string = "unexpected indentation"
	UnexpectedRParen   string = "unexpected ')'"
	UnexpectedMetaArgs string = "unexpected arguments in annotation"
	UnexpectedSection  string = "unexpected syntactic section"
	UnexpectedTerm     string = "unexpected term"

	UnmatchedLParen string = "unmatched '('"

	UnexpectedListing_Maybe_       string = "unexpected listing; did you mean to enclose it in "
	UnexpectedListing_MaybeEnclose string = "unexpected listing; did you mean to enclose it in '(' listing ')'?"

	NonBindingVariable_warn string = "variable has no bound occurrences; did you mean to use '_'?"

	ReductionFailure string = "could not reduce further"

	ExcessiveParens string = "excessive '(...)' grouping"

	CouldNotType string = "could not type"
)

func annotatesWhat(a any) string {
	switch a.(type) {
	case builtinAnnotatable:
		return "declarations, traits, trait instances, bindings, modules, data constructors, and type constructors"
	case todoAnnotatable, warnAnnotatable, errorAnnotatable:
		return "any language construct"
	case deprecatedAnnotatable:
		return "declarations, traits, trait instances, data constructors, and type constructors"
	case externalAnnotatable, noInlineAnnotatable, inlineAnnotatable:
		return "declarations and data constructors"
	case noAliasAnnotatable:
		return "data and type constructors"
	case noGcAnnotatable, testAnnotatable:
		return "declarations"
	case pureAnnotatable, infixAnnotatable, specializeAnnotatable:
		return "declarations, data constructors, and type constructors"
	default:
		return "nothing"
	}
}

func (parser *Parser) cannotAnnotateLanguageConstructWith(a any, s SyntacticElem) {
	const format string = "cannot annotate %s with '%s'\n'%s' can annotate %s"
	annotationName := getAnnotationName(a)
	validConstructs := annotatesWhat(a)
	msg := fmt.Sprintf(format, getReadableNameOfLanguageConstruct(s), annotationName, annotationName, validConstructs)
	parser.errorOn(msg, s)
}

// getAnnotationName returns the name of the annotation
func getAnnotationName(a any) string {
	switch a.(type) {
	case builtinAnnotatable:
		return `%builtin`
	case errorAnnotatable:
		return `%error`
	case todoAnnotatable:
		return `%todo`
	case warnAnnotatable:
		return `%warn`
	case deprecatedAnnotatable:
		return `%deprecated`
	case externalAnnotatable:
		return `%external`
	case inlineAnnotatable:
		return `%inline`
	case noInlineAnnotatable:
		return `%noInline`
	case specializeAnnotatable:
		return `%specialize`
	case noAliasAnnotatable:
		return `%noAlias`
	case noGcAnnotatable:
		return `%noGC`
	case pureAnnotatable:
		return `%pure`
	case infixAnnotatable:
		return `%infix`
	default:
		return ""
	}
}

// getReadableNameOfLanguageConstruct returns a human-readable name of the language construct
func getReadableNameOfLanguageConstruct(e SyntacticElem) string {
	switch s := e.(type) {
	case BindingElem:
		return "binding"
	case WhereClause:
		return "where clause"
	case MutualBlockElem:
		return "mutual block"
	case LetBindingElem:
		return "let-binding"
	case DeclarationElem:
		if s.flags&isDataConstructor != 0 {
			return "data constructor"
		}
		return "declaration"
	case DataTypeElem, HaskellStyleDataTypeElem:
		return "data type"
	case TypeConstructorElem:
		return "type constructor"
	case HaskellStyleDataConstructors:
		return "data constructor"
	case TypeElem:
		return "type"
	case TraitElem:
		return "trait"
	case InstanceElem:
		return "instance"
	case AnnotationElem:
		return "annotation"
	case ExtensionElem:
		return "syntax extension"
	case ModuleElem:
		return "module"
	default:
		return "language construct"
	}
}

func UnexpectedListingMaybeEnclosed(term termElem) string {
	if term.Term == nil {
		return UnexpectedListing_MaybeEnclose
	}
	return fmt.Sprintf("%v'(%v)'?", UnexpectedListing_Maybe_, term.Term)
}

const (
	UndefinedName string = "undefined name"
	UndefinedType string = "undefined type"

	UnknownIdent string = "unknown identifier"

	RedefAffix    string = "affix identifier re-specified"
	RedefAlias    string = "illegal definition, type alias redefined"
	RedefFunction string = "illegal definition, function redefined"
	RedefType     string = "illegal definition, type redefined"
	RedefTypeCons string = "illegal definition, type constructor redefined"
)

// creates a syntax error from the arguments
func makeNameError(msg string, path source.PathSpec, line, lineEnd, start, end int) errors.ErrorMessage {
	e := errors.MakeError("Name", msg, line, lineEnd, start, end)
	if path == nil {
		e.SourceName = "unknown"
	} else {
		e.SourceName = path.Path()
	}
	return e
}

// creates a syntax error from the arguments
func makeSyntaxError(msg string, path source.PathSpec, line, lineEnd, start, end int) errors.ErrorMessage {
	e := errors.MakeError("Syntax", msg, line, lineEnd, start, end)
	if path == nil {
		e.SourceName = "unknown"
	} else {
		e.SourceName = path.Path()
	}
	return e
}

// creates a syntax error from the arguments
func makeSyntaxWarning(msg string, path source.PathSpec, line, lineEnd, start, end int) errors.ErrorMessage {
	e := errors.MakeWarning("Syntax", msg, line, lineEnd, start, end)
	if path == nil {
		e.SourceName = "unknown"
	} else {
		e.SourceName = path.Path()
	}
	return e
}

// adds a warning constructed using the arguments
func (parser *Parser) warning2(msg string, startPos, endPos int) {
	line1, line2, char1, char2 := parser.src.CalcLocationRange(startPos, endPos)
	w := makeSyntaxWarning(msg, parser.src.Path, line1, line2, char1, char2)
	w.Message += parser.src.Window(startPos, endPos) + "\n"
	parser.addMessage(w)
}

// adds an error constructed using the arguments
func (parser *Parser) error2(msg string, startPos, endPos int) {
	line1, line2, char1, char2 := parser.src.CalcLocationRange(startPos, endPos)
	e := makeSyntaxError(msg, parser.src.Path, line1, line2, char1, char2)
	e.Message = fmt.Sprintf("%s\n%s\n", e.Message, parser.src.PointedWindow(startPos, endPos))
	parser.addMessage(e)
}

// wrapper for `(*Parser).errorOn` using `ty` and `errorNodeMessage` to generate the error
func (parser *Parser) expectedErrorOn(ty NodeType, p positioned) {
	parser.errorOn(expectedNodeMessage(ty), p)
}

func (parser *Parser) errorEOI(pos positioned) {
	// position of final character of tok
	finalStart, end := pos.Pos()
	finalStart = end - 1
	if finalStart <= 0 {
		finalStart = 0
	}

	// write an error that points to final character of token
	parser.error2(UnexpectedFinalTok, finalStart, end)
}

func (parser *Parser) kindMismatch(x Term, A Term, u types.Type, s types.Term) {
	msg := fmt.Sprintf("could not unify '%s' with '%s' from '%s : %s' and '%s : %s'", A, u, x, u, A, s)
	start, end := getTermsPos(x, A)
	parser.error2(msg, start, end)
}

func (parser *Parser) errorOn(msg string, p positioned) {
	start, end := p.Pos()
	parser.error2(msg, start, end)
}

func (parser *Parser) warningOn(msg string, p positioned) {
	start, end := p.Pos()
	parser.warning2(msg, start, end)
}

func (parser *Parser) illegalModalityError(modality token.Token, term stringPos) {
	msg := "illegal modality; did you mean '" + modality.Value + " _ : "

	name := term.String()
	_, end := term.Pos()
	if name == "" {
		parser.errorOnToken(UnexpectedModality, modality)
		return
	}
	msg = msg + "'?"
	parser.error2(msg, modality.Start, end)
}

// adds an error constructed using parser's data and the message string passed as an argument
func (parser *Parser) errorOnToken(msg string, tok token.Token) {
	start, end := tok.Start, tok.End
	parser.error2(msg, start, end)
}

// adds an error constructed using parser's data and the message string passed as an argument
func (parser *Parser) error(msg string) {
	start, end := parser.Peek().Start, parser.Peek().End
	line1, line2, char1, char2 := parser.src.CalcLocationRange(start, end)
	e := makeSyntaxError(msg, parser.src.Path, line1, line2, char1, char2)
	e.Message = fmt.Sprintf("%s\n%s\n", e.Message, parser.src.PointedWindow(start, end))
	parser.addMessage(e)
}

// reports an error when `cond == true`, otherwise resets the token position to `before`
func (parser *Parser) conditionalError(cond bool, before int, msg string) {
	if cond {
		parser.error(msg)
	} else {
		parser.tokenPos = before
	}
}

func (parser *Parser) nameError2(msg string, startPos, endPos int) {
	line1, line2, char1, char2 := parser.src.CalcLocationRange(startPos, endPos)
	e := makeNameError(msg, parser.src.Path, line1, line2, char1, char2)
	e.Message = fmt.Sprintf("%s\n%s\n", e.Message, parser.src.Window(startPos, endPos))
	parser.addMessage(e)
}
