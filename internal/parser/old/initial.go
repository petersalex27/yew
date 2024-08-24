// =================================================================================================
// Alex Peters - March 04, 2024
// =================================================================================================
package parser

import (
	"github.com/petersalex27/yew/internal/token"
)

var defaultModuleIdentifier Ident = Ident{Name: "_", Start: 0, End: 0}

func (parser *Parser) parseModuleLine() (moduleIdent Ident, ok bool) {
	if _, ok = parser.get(token.Module); !ok {
		// no module, return "_" as name
		return defaultModuleIdentifier, true
	}

	var notPascalCase bool
	var ident Ident
	ident, ok, notPascalCase = parser.getPascalCaseIdent()
	if !ok {
		if notPascalCase {
			parser.errorOn(RequirePascalCaseModule, ident)
			return
		}
		parser.errorOn(ExpectedIdentifier, ident)
		return
	}

	moduleIdent = ident
	return
}

// drops all contiguous comments starting at current token
func (parser *tokenInfo) dropComments() {
	for parser.Peek().Type == token.Comment {
		_ = parser.Advance() // advance past comment
	}
}

func (parser *Parser) Parse() (module Module, ok bool) {
	module = Module{name: defaultModuleIdentifier}
	ok = parser.readEnvironment() &&
		parser.declareModule(&module) &&
		parser.readImports() &&
		firstPass(parser)
	return module, ok
}
