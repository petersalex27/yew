// =================================================================================================
// Alex Peters - May 17, 2024
// =================================================================================================
package parser

type (
	annotatable interface {
		SymbolicElem
	}

	builtinAnnotatable interface {
		markBuiltin(parser *Parser)
	}

	errorAnnotatable interface {
		throwError(parser *Parser, msg string)
	}

	testAnnotatable interface {
		markTest(parser *Parser)
	}

	todoAnnotatable interface {
		reportTodo(parser *Parser, msg string)
	}

	warnAnnotatable interface {
		throwWarning(parser *Parser, msg string)
	}

	deprecatedAnnotatable interface {
		markDeprecated(msg string)
	}

	externalAnnotatable interface {
		markExternal(parser *Parser)
	}

	inlineAnnotatable interface {
		markInline(parser *Parser)
	}

	noInlineAnnotatable interface {
		markNoInline(parser *Parser)
	}

	specializeAnnotatable interface {
		specializeWith(names []string)
	}

	noAliasAnnotatable interface {
		markNoAlias(parser *Parser)
	}

	noGcAnnotatable interface {
		markNoGC(parser *Parser)
	}

	pureAnnotatable interface {
		markPure(parser *Parser)
	}

	infixAnnotatable interface {
		markInfix(parser *Parser, isRightAssociative bool, bp int8)
	}
)

func genAnnotate(annot any, args ...any) (annotate func(parser *Parser, a any) bool) {
	switch annot.(type) {
	case builtinAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(builtinAnnotatable); ok {
				v.markBuiltin(parser)
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case errorAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(errorAnnotatable); ok {
				v.throwError(parser, args[0].(string))
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case todoAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(todoAnnotatable); ok {
				v.reportTodo(parser, args[0].(string))
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case warnAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(warnAnnotatable); ok {
				v.throwWarning(parser, args[0].(string))
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case deprecatedAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(deprecatedAnnotatable); ok {
				v.markDeprecated(args[0].(string))
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case externalAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(externalAnnotatable); ok {
				v.markExternal(parser)
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case inlineAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(inlineAnnotatable); ok {
				v.markInline(parser)
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case noInlineAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(noInlineAnnotatable); ok {
				v.markNoInline(parser)
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case specializeAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(specializeAnnotatable); ok {
				v.specializeWith(args[0].([]string))
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case noAliasAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(noAliasAnnotatable); ok {
				v.markNoAlias(parser)
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case noGcAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(noGcAnnotatable); ok {
				v.markNoGC(parser)
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case pureAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(pureAnnotatable); ok {
				v.markPure(parser)
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	case infixAnnotatable:
		return func(parser *Parser, a any) bool {
			if v, ok := a.(infixAnnotatable); ok {
				v.markInfix(parser, args[0].(bool), args[1].(int8))
				return true
			}
			parser.cannotAnnotateLanguageConstructWith(annot, a.(SymbolicElem))
			return false
		}
	}

	panic("bug: invalid annotation type")
}

// markBuiltin marks the declaration as builtin, and reports an error if isExternal is set
//
// use parser.Panicking() to check for errors
func (decl *DeclarationElem) markBuiltin(parser *Parser) {
	// builtin and external are mutually exclusive, report an error if external is set
	if decl.flags & isExternal != 0 {
		parser.errorOn(BuiltinExternalConflict, decl)
		return
	}
	decl.flags |= isBuiltin
}

// markExternal marks the declaration as external, and reports an error if isBuiltin is set
//
// use parser.Panicking() to check for errors
func (decl *DeclarationElem) markExternal(parser *Parser) {
	// external and builtin are mutually exclusive, report an error if builtin is set
	if decl.flags & isBuiltin != 0 {
		parser.errorOn(BuiltinExternalConflict, decl)
		return
	}
	decl.flags |= isExternal
}

// markInline marks the declaration as suggestInline, and reports an error if noInline is set
//
// use parser.Panicking() to check for errors
func (decl *DeclarationElem) markInline(parser *Parser) {
	// suggestInline and noInline are mutually exclusive, report an error if noInline is set
	if decl.flags & noInline != 0 {
		parser.errorOn(InlineNotAllowed, decl)
		return
	}
	decl.flags |= suggestInline
}

// markNoInline marks the declaration as noInline, and reports an error if suggestInline is set
//
// use parser.Panicking() to check for errors
func (decl *DeclarationElem) markNoInline(parser *Parser) {
	// noInline and suggestInline are mutually exclusive, report an error if suggestInline is set
	if decl.flags & suggestInline != 0 {
		parser.errorOn(IllegalNoInlineSuggestion, decl)
		return
	}
	decl.flags |= noInline
}

// markPure marks the declaration as pure, and warns the user if the function cannot be verified as pure
func (decl *DeclarationElem) markPure(parser *Parser) {
	// if decl is external, warn the user that the function cannot be verified as pure, but still 
	// mark it as pure
	if decl.flags & isExternal != 0 {
		parser.warningOn(UnverifiablePureAnnotation, decl)
	}
	decl.flags |= isPure
}

func (decl *DeclarationElem) markInfix(parser *Parser, isRightAssociative bool, bp int8) {
	decl.rAssoc = isRightAssociative
	decl.bp = bp
}

func (decl *DeclarationElem) markTest(parser *Parser) {
	decl.flags |= isTest
	tests := parser.classifications["test"]
	parser.classifications["test"] = append(tests, decl)
}