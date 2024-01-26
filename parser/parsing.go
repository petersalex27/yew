// =================================================================================================
// Alex Peters - January 24, 2024
//
// The actual parsing functions/methods
// =================================================================================================

package parser

import "github.com/petersalex27/yew/token"

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
	p.module = &mod // TODO: remove--this just here to avoid compiler errors rn
}

func (parser *Parser) equalToken() (eq token.Token, ok bool) {
	return parser.getToken(token.Equal, ExpectedEqual)
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