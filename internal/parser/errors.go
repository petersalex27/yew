package parser

import (
	"fmt"
	"os"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/errors"
)

const (
	BadImport                    string = "expected package name or import group"                                        // bad-import
	ExpectedAliasBinding         string = "expected '=' to follow type alias name"                                       // expected-alias-binding
	ExpectedBindingTerm          string = "expected a binding term"                                                      // expected-binding-term
	ExpectedBodyElement          string = "expected body element"                                                        // expected-body-element
	ExpectedCaseArmThickArrow    string = "expected '=>' to follow case arm pattern"                                     // expected-case-arm-thick-arrow
	ExpectedColonEqual           string = "expected ':='"                                                                // expected-colon-equal
	ExpectedCommand              string = "expected command"                                                             // expected-command
	ExpectedConstrainer          string = "expected constrainer"                                                         // expected-constrainer
	ExpectedDef                  string = "expected definition"                                                          // expected-def
	ExpectedEndOfFile            string = "expected end of file"                                                         // expected-eof
	ExpectedForallBinder         string = "expected forall binder"                                                       // expected-forall-binder
	ExpectedForallIn             string = "expected 'in' to follow 'forall' binders"                                     // expected-forall-in
	ExpectedId                   string = "expected identifier"                                                          // expected-id
	ExpectedIn                   string = "expected 'in'"                                                                // expected-in
	ExpectedInstWhere            string = "expected 'where' clause to follow inst declaration"                           // expected-inst-where
	ExpectedLambdaAbstraction    string = "expected lambda abstraction"                                                  // expected-lambda-abstraction
	ExpectedLambdaThickArrow     string = "expected '=>' to follow lambda binders"                                       // expected-lambda-thick-arrow
	ExpectedMainElement          string = "expected main element"                                                        // expected-main-element
	ExpectedModuleId             string = "expected module name"                                                         // expected-module-id
	ExpectedName                 string = "expected name"                                                                // expected-name
	ExpectedNamespaceAlias       string = "expected namespace alias to follow 'as'"                                      // expected-namespace-alias
	ExpectedOf                   string = "expected 'of' to follow case scrutinee"                                       // expected-of
	ExpectedPattern              string = "expected pattern"                                                             // expected-pattern
	ExpectedPatternUnit          string = "expected pattern unit"                                                        // expected-pattern-unit
	ExpectedRightBrace           string = "expected '}'"                                                                 // expected-right-brace
	ExpectedRightParen           string = "expected ')'"                                                                 // expected-right-paren
	ExpectedSlashModuleId        string = "expected module identifier to follow '/'"                                     // expected-slash-module-id
	ExpectedSpecDef              string = "expected spec definition"                                                     // expected-spec-def
	ExpectedSpecInst             string = "expected spec instance"                                                       // expected-spec-inst
	ExpectedSpecWhere            string = "expected 'where' clause to follow spec declaration"                           // expected-spec-where
	ExpectedStringLit            string = "expected string literal"                                                      // expected-string-lit
	ExpectedSymbol               string = "expected symbol"                                                              // expected-symbol
	ExpectedSyntax               string = "expected syntax definition"                                                   // expected-syntax
	ExpectedSyntaxBinding        string = "expected '=' to follow syntax rule"                                           // expected-syntax-binding
	ExpectedSyntaxBindingId      string = "expected syntax binding identifier"                                           // expected-syntax-binding-id
	ExpectedSyntaxRule           string = "expected syntax rule"                                                         // expected-syntax-rule
	ExpectedTypeAlias            string = "expected type alias"                                                          // expected-type-alias
	ExpectedTypeAliasName        string = "expected a name to follow 'alias'"                                            // expected-type-alias-name
	ExpectedTypeConstructor      string = "expected type constructor"                                                    // expected-type-constructor
	ExpectedTypeConstructorName  string = "expected type constructor name"                                               // expected-type-constructor-name
	ExpectedTypeSig              string = "expected type signature"                                                      // expected-type-sig
	ExpectedTyping               string = "expected typing"                                                              // expected-typing
	ExpectedTypingOrDef          string = "expected a typing or definition"                                              // expected-typing-or-def
	ExpectedUpperId              string = "expected uppercase identifier"                                                // expected-upper-id
	ExpectedWithArmThickArrow    string = "expected '=>' to follow with arm pattern"                                     // expected-with-arm-thick-arrow
	ExpectedWithClause           string = "expected 'with' clause"                                                       // expected-with-clause
	ExpectedWithClauseArm        string = "expected 'with' clause arm"                                                   // expected-with-clause-arm
	IllegalMultipleEnclosure     string = "illegal multiply enclosed term"                                               // illegal-multiple-enclosure
	IllegalNamespaceAlias        string = "illegal namespace alias, expected lowercase identifier"                       // illegal-namespace-alias
	IllegalOpenModifier          string = "modifier 'open' can only target data type definitions"                        // illegal-open-modifier
	IllegalOpenModifierTyping    string = "modifier 'open' targeted a typing, but no constructors were found"            // illegal-open-modifier-typing
	IllegalVisibleDef            string = "visibility modifiers cannot be applied to definitions, only their signatures" // illegal-visible-def
	InvalidAnnotationTarget      string = "cannot find a valid target for annotations"                                   // invalid-annotation-target
	UnexpectedEOF                string = "unexpected end of file"                                                       // unexpected-eof
	UnexpectedStructure          string = "unexpected structure in source body"                                          // unexpected-structure
	UnexpectedToken              string = "unexpected token"                                                             // unexpected-token
	ExpectedBoundExpr            string = "let-binding requires a bound expression"                                      // expected-bound-expr
	ExpectedType                 string = "expected a type"                                                              // expected-type
	ExpectedLetExpr              string = "expected let expression"                                                      // expected-let-expr
	ExpectedExpr                 string = "expected expression"                                                          // expected-expr
	ExpectedConstraint           string = "expected type constraint"                                                     // expected-type-constraint
	ExpectedDerivingBody         string = "expected body for deriving clause"                                            // expected-deriving-body
	ExpectedCaseArm              string = "expected case arm"                                                            // expected-case-arm
	ExpectedModalId              string = "modality must be followed by an identifier"                                   // expected-modal-id
	IllegalVisibilityTarget      string = "illegal target for visibility modifier"                                       // illegal-visibility-target
	ExpectedImportPath           string = "expected import path"                                                         // expected-import-path
	IllegalEmptyUsingClause      string = "illegal empty using clause"                                                   // illegal-empty-using-clause
	IllegalUnenclosedUsingClause string = "illegal unenclosed symbol selection in using clause"                          // illegal-unenclosed-using-clause
	ExpectedAccessDot            string = "expected '.'"                                                                 // expected-access-dot
)

var makePos = (api.Positioned).GetPos

func mkErr(msg string, positioned api.Positioned) data.Err {
	return data.MkErr(msg, makePos(positioned))
}

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
