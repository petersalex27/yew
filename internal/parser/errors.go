package parser

import (
	"fmt"
	"os"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/errors"
)

const (
	BadImport                       = "expected package name or import group"                                        // bad-import
	ExpectedAccessDot               = "expected '.'"                                                                 // expected-access-dot
	ExpectedAliasBinding            = "expected '=' to follow type alias name"                                       // expected-alias-binding
	ExpectedBindingTerm             = "expected a binding term"                                                      // expected-binding-term
	ExpectedBodyElement             = "expected body element"                                                        // expected-body-element
	ExpectedBoundExpr               = "let-binding requires a bound expression"                                      // expected-bound-expr
	ExpectedCaseArm                 = "expected case arm"                                                            // expected-case-arm
	ExpectedCaseArmThickArrow       = "expected '=>' to follow case arm pattern"                                     // expected-case-arm-thick-arrow
	ExpectedColonEqual              = "expected ':='"                                                                // expected-colon-equal
	ExpectedCommand                 = "expected command"                                                             // expected-command
	ExpectedConstrainer             = "expected constrainer"                                                         // expected-constrainer
	ExpectedConstraint              = "expected type constraint"                                                     // expected-type-constraint
	ExpectedDef                     = "expected definition"                                                          // expected-def
	ExpectedDerivingBody            = "expected body for deriving clause"                                            // expected-deriving-body
	ExpectedEndOfFile               = "expected end of file"                                                         // expected-eof
	ExpectedExpr                    = "expected expression"                                                          // expected-expr
	ExpectedForallBinder            = "expected forall binder"                                                       // expected-forall-binder
	ExpectedForallIn                = "expected 'in' to follow 'forall' binders"                                     // expected-forall-in
	ExpectedId                      = "expected identifier"                                                          // expected-id
	ExpectedImportPath              = "expected import path"                                                         // expected-import-path
	ExpectedIn                      = "expected 'in'"                                                                // expected-in
	ExpectedInstWhere               = "expected 'where' clause to follow inst declaration"                           // expected-inst-where
	ExpectedLambdaAbstraction       = "expected lambda abstraction"                                                  // expected-lambda-abstraction
	ExpectedLambdaThickArrow        = "expected '=>' to follow lambda binders"                                       // expected-lambda-thick-arrow
	ExpectedLetExpr                 = "expected let expression"                                                      // expected-let-expr
	ExpectedMainElement             = "expected main element"                                                        // expected-main-element
	ExpectedModalId                 = "modality must be followed by an identifier"                                   // expected-modal-id
	ExpectedModuleId                = "expected module name"                                                         // expected-module-id
	ExpectedName                    = "expected name"                                                                // expected-name
	ExpectedNamespaceAlias          = "expected namespace alias to follow 'as'"                                      // expected-namespace-alias
	ExpectedOf                      = "expected 'of' to follow case scrutinee"                                       // expected-of
	ExpectedPattern                 = "expected pattern"                                                             // expected-pattern
	ExpectedPatternUnit             = "expected pattern unit"                                                        // expected-pattern-unit
	ExpectedRightBrace              = "expected '}'"                                                                 // expected-right-brace
	ExpectedRightParen              = "expected ')'"                                                                 // expected-right-paren
	ExpectedSlashModuleId           = "expected module identifier to follow '/'"                                     // expected-slash-module-id
	ExpectedSpecDef                 = "expected spec definition"                                                     // expected-spec-def
	ExpectedSpecInst                = "expected spec instance"                                                       // expected-spec-inst
	ExpectedSpecWhere               = "expected 'where' clause to follow spec declaration"                           // expected-spec-where
	ExpectedStringLit               = "expected  literal"                                                            // expected--lit
	ExpectedSymbol                  = "expected symbol"                                                              // expected-symbol
	ExpectedSyntax                  = "expected syntax definition"                                                   // expected-syntax
	ExpectedSyntaxBinding           = "expected '=' to follow syntax rule"                                           // expected-syntax-binding
	ExpectedSyntaxBindingId         = "expected syntax binding identifier"                                           // expected-syntax-binding-id
	ExpectedSyntaxRule              = "expected syntax rule"                                                         // expected-syntax-rule
	ExpectedType                    = "expected a type"                                                              // expected-type
	ExpectedTypeAlias               = "expected type alias"                                                          // expected-type-alias
	ExpectedTypeAliasName           = "expected a name to follow 'alias'"                                            // expected-type-alias-name
	ExpectedTypeConstructor         = "expected type constructor"                                                    // expected-type-constructor
	ExpectedTypeConstructorName     = "expected type constructor name"                                               // expected-type-constructor-name
	ExpectedTypeJudgment            = "expected type judgement"                                                      // expected-type-judgement
	ExpectedTypeSig                 = "expected type signature"                                                      // expected-type-sig
	ExpectedTyping                  = "expected typing"                                                              // expected-typing
	ExpectedTypingOrDef             = "expected a typing or definition"                                              // expected-typing-or-def
	ExpectedUpperId                 = "expected uppercase identifier"                                                // expected-upper-id
	ExpectedWithArmThickArrow       = "expected '=>' to follow with arm pattern"                                     // expected-with-arm-thick-arrow
	ExpectedWithClause              = "expected 'with' clause"                                                       // expected-with-clause
	ExpectedWithClauseArm           = "expected 'with' clause arm"                                                   // expected-with-clause-arm
	IllegalEmptyUsingClause         = "illegal empty using clause"                                                   // illegal-empty-using-clause
	IllegalLowercaseConstructorName = "constructor names cannot be lowercase identifiers"                            // illegal-lowercase-constructor-name
	IllegalMethodTypeConstructor    = "type constructors cannot be identified by method identifiers"                 // illegal-method-type-constructor
	IllegalMultipleEnclosure        = "illegal multiply enclosed term"                                               // illegal-multiple-enclosure
	IllegalNamespaceAlias           = "illegal namespace alias, expected lowercase identifier"                       // illegal-namespace-alias
	IllegalOpenModifier             = "modifier 'open' can only target data type definitions"                        // illegal-open-modifier
	IllegalOpenModifierTyping       = "modifier 'open' targeted a typing, but no constructors were found"            // illegal-open-modifier-typing
	IllegalUnenclosedUsingClause    = "illegal unenclosed symbol selection in using clause"                          // illegal-unenclosed-using-clause
	IllegalVisibilityTarget         = "illegal target for visibility modifier"                                       // illegal-visibility-target
	IllegalVisibleDef               = "visibility modifiers cannot be applied to definitions, only their signatures" // illegal-visible-def
	InvalidAnnotationTarget         = "cannot find a valid target for annotations"                                   // invalid-annotation-target
	UnexpectedEOF                   = "unexpected end of file"                                                       // unexpected-eof
	UnexpectedStructure             = "unexpected structure in source body"                                          // unexpected-structure
	UnexpectedToken                 = "unexpected token"                                                             // unexpected-token
	ExpectedConstraintElem          = "expected constraint element"                                                  // expected-constraint-elem
)

var makePos = (api.Positioned).GetPos

func parseError(p Parser, e data.Err) error {
	start, end := e.Pos()
	return errors.Syntax(p.srcCode(), e.Msg(), start, end)
}

func parseErrors(p Parser, es data.Ers) []error {
	errs := make([]error, es.Len())
	for i, e := range es.Elements() {
		errs[i] = parseError(p, e)
	}
	return errs
}

func printErrors(es ...error) {
	for _, e := range es {
		fmt.Fprintf(os.Stderr, "%s\n", e.Error())
	}
}

func makeError(p Parser, msg string, ps ...api.Positioned) error {
	var start, end int
	if len(ps) == 0 {
		start, end = 0, 0
	} else {
		start, end = api.WeakenRangeOver(ps[0], ps[1:]...).Pos()
	}
	return errors.Syntax(p.srcCode(), msg, start, end)
}
