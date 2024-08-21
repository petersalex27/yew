// =================================================================================================
// Alex Peters - March 04, 2024
// =================================================================================================
package parser

import (
	"github.com/petersalex27/yew/token"
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

func (parser *Parser) importModules() {
	// TODO: implement
}

func (parser *Parser) Begin() (module Module, ok bool) {
	//parser.tryIndentPush()
	// NOTE: leading indents probably cause an issue--figure out what to do with them
	//	- do they decide the base indent? Yes? Yes, but only sometimes (when?)? No?
	parser.drop()
	module.name, ok = parser.parseModuleLine()
	if !ok {
		return
	}

	parser.drop()
	for ok && parser.Peek().Type == token.Import {
		ok = parser.parseImports()
		parser.drop()
	}
	if !ok {
		return
	}

	parser.firstPass()
	// TODO: at this point imported modules should be parsed (can be done in parallel)
	parser.importModules()

	// TODO: now, finish parsing this file with the information from now parsed imports
	return
}
