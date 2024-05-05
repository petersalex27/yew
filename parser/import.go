// =================================================================================================
// Alex Peters - March 04, 2024
// =================================================================================================
package parser

import "github.com/petersalex27/yew/token"

type Import struct {
	Lookup token.Token // name used to lookup import
	Id     token.Token // name used in source code
}

type ImportTable map[string]Import

func (parser *Parser) parseImportLineHelper() (lookup, id token.Token, ok bool) {
	lookup, ok = parser.get(token.Id)
	if !ok {
		parser.error(ExpectedIdentifier)
		return
	}

	id = lookup

	if parser.Peek().Type != token.As {
		return // regular package import
	}

	// parse "qualified-as" name
	_ = parser.Advance()
	id, ok = parser.get(token.Id)
	if !ok {
		parser.error(ExpectedIdentifier)
		return
	}
	return
}

// adds an import to an import table `imports` from the "module-lookup-name" `lookup` and the "as-name" `id`
//
// ok is true on successful registration, else false and an error message is returned
func (imports *ImportTable) register(lookupName, asName token.Token) (ok bool, errorMessage string) {
	if _, found := (*imports)[asName.Value]; found {
		return false, DuplicateImportName
	}

	(*imports)[asName.Value] = Import{Lookup: lookupName, Id: asName}
	return true, ""
}

func parseImportLine(parser *Parser) (ok bool) {
	var lookup, asId token.Token
	lookup, asId, ok = parser.parseImportLineHelper()
	if !ok {
		return // error already reported in parseImportLineHelper
	}

	var errorMsg string
	ok, errorMsg = parser.imports.register(lookup, asId)
	if !ok {
		parser.error(errorMsg)
	}
	return
}

func (parser *Parser) parseImports() (ok bool) {
	if parser.Peek().Type != token.Import {
		return true
	}
	again := func(p *Parser, i int) bool { return p.equalIndent(i) }
	return parser.parseSection(parseImportLine, again)
}
