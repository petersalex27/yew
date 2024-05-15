// =================================================================================================
// Alex Peters - March 04, 2024
// =================================================================================================
package parser

import (
	"fmt"

	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/source"
	"github.com/petersalex27/yew/token"
)

func expectedNodeMessage(ty NodeType) string {
	switch ty {
	case identType:
		return ExpectedIdentifier
	case wildcardType:
		return ExpectedWildcard
	case affixedType:
		return ExpectedAffixedId
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
	case bindingType:
		return ExpectedBinding
	case typingType:
		return ExpectedTyping
	case listType:
		return ExpectedListType
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
	case constrainedTypeType:
		return ExpectedConstraint
	case syntaxExtensionType:
		return ExpectedSyntaxExtension
	default:
		return UnexpectedSection
	}
}

func expectedMessage(ty token.Type) string {
	switch ty {
	case token.Equal:
		return ExpectedEqual
	case token.Id:
		return ExpectedIdentifier
	case token.In:
		return ExpectedIn
	case token.Indent:
		return ExpectedIndentation
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
	DuplicateImportName string = "duplicate import name"

	// = "expected" ==================================================================================
	ExpectedBinding           string = "expected pattern binding"
	ExpectedEqual             string = "expected assignment"
	ExpectedCommaOrThickArrow string = "expected ',' or '=>'"
	ExpectedIdentifier        string = "expected identifier"
	ExpectedLeftRightNone     string = "expected constant 'Left', 'Right', or 'None'"
	ExpectedPascalCase        string = "expected PascalCase identifier"
	ExpectedIn                string = "expected 'in'"
	ExpectedIndentation       string = "expected indentation"
	ExpectedModule            string = "expected 'module'"
	ExpectedName              string = "expected name"
	ExpectedRBrace            string = "expected '}'"
	ExpectedRBracket          string = "expected ']'"
	ExpectedRParen            string = "expected ')'"
	ExpectedThickArrow        string = "expected '=>'"
	ExpectedType              string = "expected type"
	ExpectedTypeApplication   string = "expected type application"
	ExpectedTyping            string = "expected typing"
	ExpectedVariable          string = "expected variable"
	ExpectedWhere             string = "expected 'where'"
	ExpectedMutual            string = "expected 'mutual'"
	ExpectedDeclaration       string = "expected declaration"
	ExpectedScrutinee         string = "expected scrutinee"
	ExpectedExpression        string = "expected expression"
	ExpectedLRN               string = "expected one of the constants `Left`, `Right`, or `None`"
	ExpectedInteger           string = "expected an integer literal"
	ExpectedInteger0to10      string = "expected an integer literal in the range [0, 10]"
	ExpectedGreaterIndent     string = "expected larger indentation than enclosing context's indentation"
	ExpectedLParen            string = "expected '('"
	ExpectedUint              string = "expected unsigned integer"
	ExpectedUintRange1_9      string = "expected unsigned integer in the range 1 .. 9"
	ExpectedConstraint        string = "expected constraint, e.g.,\n\t'(Trait a1, .., Trait aN) => ..', or\n\t'Trait a => ..'"
	ExpectedWildcard          string = "expected wildcard"
	ExpectedAffixedId         string = "expected affixed identifier"
	ExpectedIntLit            string = "expected int literal"
	ExpectedCharLit           string = "expected char literal"
	ExpectedFloatLit          string = "expected float literal"
	ExpectedStringLit         string = "expected string literal"
	ExpectedFunctionType      string = "expected function type"
	ExpectedApplication       string = "expected function application"
	ExpectedLambdaAbstraction string = "expected lambda abstraction"
	ExpectedListType          string = "expected list type"
	ExpectedListExpr          string = "expected list expression"
	ExpectedTuple             string = "expected tuple"
	ExpectedTupleExpr         string = "expected tuple expression"
	ExpectedPairings          string = "expected pairings"
	ExpectedListing           string = "expected listing"
	ExpectedSyntaxExtension   string = "expected syntax extension"
	ExpectedTermAfter         string = "expected another syntactic sub-section following this one"

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

	MultiplyOccurringAffixedIdent string = "multiple occurrences of affixed form of identifier"

	EmptyMutualBlock string = "empty mutual block"

	MalformedAffixAnnotation string = "malformed affix annotation"
	BadIdent                 string = "malformed identifier"

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

	UnexpectedListing_Maybe_       string = "unexpected listing; did you mean to enclose it in "
	UnexpectedListing_MaybeEnclose string = "unexpected listing; did you mean to enclose it in '(' listing ')'?"

	NonBindingVariable_warn string = "variable has no bound occurrences; did you mean to use '_'?"

	ReductionFailure string = "could not reduce further"

	ExcessiveParens string = "excessive '(...)' grouping"
)

func UnexpectedListingMaybeEnclosed(term Term) string {
	if term == nil {
		return UnexpectedListing_MaybeEnclose
	}
	return fmt.Sprintf("%v'(%v)'?", UnexpectedListing_Maybe_, term)
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

func (parser *Parser) errorOn(msg string, p positioned) {
	start, end := p.Pos()
	parser.error2(msg, start, end)
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

func (parser *Parser) nameError2(msg string, startPos, endPos int) {
	line1, line2, char1, char2 := parser.src.CalcLocationRange(startPos, endPos)
	e := makeNameError(msg, parser.src.Path, line1, line2, char1, char2)
	e.Message = fmt.Sprintf("%s\n%s\n", e.Message, parser.src.Window(startPos, endPos))
	parser.addMessage(e)
}
