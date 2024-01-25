// =================================================================================================
// Alex Peters - January 24, 2024
//
// The actual parsing functions/methods
// =================================================================================================

package parser

import "github.com/petersalex27/yew/token"

// parses import
func (p *Parser) parseImport() (_ Import, ok bool) {
	var importKeyword, name token.Token
	if importKeyword, ok = p.importToken(); !ok {
		return
	}

	if name, ok = p.idToken(); !ok {
		return
	}

	end := p.StartOptional()
	defer end()

	if _, ok = p.whereToken(); !ok {
		start, end := importKeyword.Start, name.End
		return Import{Start: start, End: end, ImportName: name, LookupName: name}, true
	}

	
}

func (p *Parser) parseImports() []Import {
	
}

// builds a module from a newly initialized module
func (p *Parser) BuildModule() {
	// for now, comments get dropped // TODO: include comments in the AST
	p.skip(token.Comment)

	end := p.StartOptional()
	//p.CheckAdvance(token.Import)
	end()

	var module, name token.Token
	var ok bool

	if module, ok = p.moduleToken(); !ok {
		return
	}
	
	if name, ok = p.idToken(); !ok {
		return
	}

	if _, ok = p.whereToken(); !ok {
		return
	}

	mod := Module{
		Start: module.Start,
		ModuleName: name,
	}
}

func (p *Parser) getToken(ty token.Type, errorMessage string) (tok token.Token, ok bool) {
	tok, ok = p.ConditionalAdvance(ty)
	if !ok && !p.optionalFlag {
		p.error(errorMessage)
	}
	return
}

// gets id token, not node
func (p *Parser) idToken() (idTok token.Token, ok bool) {
	return p.getToken(token.Id, ExpectedIdentifier)
}

// gets import token
func (p *Parser) importToken() (token.Token, bool) {
	return p.getToken(token.Import, UnexpectedToken)
}

// parses 'module' keyword token
func (p *Parser) moduleToken() (tok token.Token, ok bool) {
	return p.getToken(token.Module, ExpectedModule)
}

// parses 'where' keyword token
func (p *Parser) whereToken() (token.Token, bool) {
	return p.getToken(token.Where, ExpectedWhere)
}