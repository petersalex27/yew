// =================================================================================================
// Alex Peters - January 25, 2024
//
// Handles parsing of imports
// =================================================================================================
package parser

import "github.com/petersalex27/yew/token"

// parses import
func (parser *Parser) parseImport() (im Import, ok bool) {
	// var importKeyword token.Token
	// if importKeyword, ok = parser.importToken(); !ok {
	// 	return
	// }

	// TODO: finish

	return
}

// parses an import block
func (parser *Parser) parserImportBlock()

func (parser *Parser) parseImportData() (im Import, ok bool) {
	var name token.Token

	if name, ok = parser.idToken(); !ok {
		return
	}

	im.ImportName = name

	im.Start = name.Start

	var done bool
	done, ok = parser.parseImportOptionalWhere(&im, name)
	if done {
		return
	}

	ok = parser.parseImportContext(&im)
	return
}

// requires parsing of context ('where' and an assignment)
func (parser *Parser) parseImportContext(im *Import) (ok bool) {
	var destId, srcId token.Token
	endStopped := parser.StopOptional()
	defer endStopped()

	if destId, ok = parser.idToken(); !ok {
		return
	}

	if _, ok = parser.equalToken(); !ok {
		return
	}

	if srcId, ok = parser.idToken(); !ok {
		return
	}

	noContextAdded := im.ImportName.Value != destId.Value
	if noContextAdded {
		// no context is added b/c the assignment is pointless as it assigns some unused var to another
		// thing
		parser.warning2(UnusedContext, destId.Start, srcId.End)
		im.LookupName = im.ImportName
	} else {
		// context added is destId = srcId, and destId has the same value as ImportName
		im.LookupName = srcId
	}

	im.End = srcId.End
	return
}

func (parser *Parser) parseImports() []Import {
	return nil // TODO
}

func (parser *Parser) parseImportOptionalWhere(im *Import, name token.Token) (done, ok bool) {
	end := parser.StartOptional()
	defer end()

	ok = true // unconditionally

	if _, found := parser.whereToken(); !found {
		done = true
		im.End = name.End
		im.LookupName = name
	}
	return
}
