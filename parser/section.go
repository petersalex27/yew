// =================================================================================================
// Alex Peters -2024
//
// general parsing for syntactic sections of code
// =================================================================================================
package parser

import "github.com/petersalex27/yew/token"

// transfers and receives information related to parsing sections
type sectionHelper func()

func (c sectionHelper) SectionOpened() bool { return c != nil }

func (c sectionHelper) Clean() {
	if c.SectionOpened() {
		c()
	}
}

// true iff the next token is one which can be used as the name of something declarable
func declarationValidator(parser *Parser) bool {
	ty := parser.Peek().Type
	return ty == token.Id || ty == token.Affixed
}

// returns true iff the block marker `marker` requires an indent to follow it
func indentRequired(marker token.Type, newlinesDropped int) bool {
	// `where` and mutual always require a succeeding indent; other tokens only require it if a
	// newline directly follows `marker`
	return marker == token.Where || marker == token.Mutual || newlinesDropped > 0
}

func (parser *Parser) maybeReportErrorWhenNoIndentFound(marker token.Token, newlinesDropped int) {
	if indentRequired(marker.Type, newlinesDropped) {
		// error: a newline followed one of the section-flagging keywords but no indent was found
		parser.error(ExpectedIndentation)
	}
	// no indent required (token is not 'where' and no newlines were dropped)
}

// tries to open a new syntactic section
//
// check if an error occurred with parser.Panicking()
func (parser *Parser) openSection(validator func(*Parser) bool) (openedFrom token.Token, c sectionHelper, validated bool) {
	if validated = validator(parser); !validated {
		return
	}

	// token that opens the section
	openedFrom = parser.Advance()

	// drop anything droppable following keyword
	newlinesDropped := parser.drop()

	// NOTE: no need to check for a newline because an indent can only occur after a newline

	// if next token is not an indent, `success` will be set to false
	indentTok, success := parser.findMeaningfulIndent()
	if !success {
		parser.maybeReportErrorWhenNoIndentFound(openedFrom, newlinesDropped)
		return openedFrom, nil, validated
	}

	indent := len(indentTok.Value)
	if success = parser.testIndentGT(indent); !success {
		// error: indentation is not larger than enclosing context's indentation
		parser.error(ExpectedGreaterIndent)
		return
	}

	parser.indentation.Push(indent)
	c = parser.indentation.silentPop
	return
}

// tries to open a new syntactic section
//
// errors are always reported in this function (this is different than `openSection`)
func (parser *Parser) openMutualWhereSection() (openedFrom token.Token, c sectionHelper, validated bool) {
	if validated = parser.Peek().Type == token.Mutual; !validated {
		parser.error(ExpectedMutual)
		return
	}

	// token that opens the section
	openedFrom = parser.Advance()

	// drop anything droppable following keyword
	parser.drop()

	indentTok, nonInlineWhere := parser.findMeaningfulIndent()
	if !nonInlineWhere {
		// inline where, should be:
		//	mutual where
		//		...
		_, c, validated = parser.openSection(func(p *Parser) bool { return p.Peek().Type == token.Where })
		if !validated {
			parser.error(ExpectedWhere)
			return
		}
		return
	}

	// indent, should be:
	//	mutual
	//		where
	//			...

	indent := len(indentTok.Value)
	if success := parser.testIndentGT(indent); !success {
		// error: indentation is not larger than enclosing context's indentation
		parser.error(ExpectedGreaterIndent)
		return
	}

	parser.indentation.Push(indent)

	// now do where block validation
	var cWhere sectionHelper
	_, cWhere, validated = parser.openSection(func(p *Parser) bool { return p.Peek().Type == token.Where })
	if !validated {
		parser.error(ExpectedWhere)
		return
	}
	validated = !parser.panicking
	if !validated {
		return
	}

	// clean from inner most to outermost
	c = func() {
		cWhere.Clean() // first clean where block
		parser.indentation.silentPop() // now clean mutual block
	}
	return
}

// when action can happen multiple times, this function should be called. It drives the parsing of
// the section and handling of the indentation signifying whether the section is continued or ended
func (parser *Parser) sectionActionLoop(action func(*Parser) bool, againCondition func(*Parser, int) bool) (ok bool) {
	for again := true; again; {
		ok = action(parser)
		parser.drop()
		pos, success := parser.locateMeaningfulIndent()
		again = success && ok
		if !again {
			return ok
		}

		level := len(parser.tokens[pos].Value)
		again = againCondition(parser, level)
		if !again {
			parser.tokenPos = pos // rewind
		}
	}
	return ok
}

// atMostOnce should be true when `cleaner.Ok() == false`; it signifies that, at most, the action should be run once
func (parser *Parser) runSection(atMostOnce bool, action func(*Parser) (ok bool), againCondition func(*Parser, int) bool) (ok bool) {
	if !atMostOnce {
		return parser.sectionActionLoop(action, againCondition)
	}

	// if not panicking, do action just once
	ok = !parser.panicking && action(parser)
	return
}

func (parser *Parser) parseSection2(action func(*Parser) (ok bool), againCondition func(*Parser, int) bool) (ok bool, openedFrom token.Token) {
	// create a new context and make sure it's removed when this function returns
	openedFrom, cleaner, _ := parser.openSection(func(*Parser) bool { return true }) // will always validate as true
	defer cleaner.Clean()

	return parser.runSection(!cleaner.SectionOpened(), action, againCondition), openedFrom
}

// assumes section marker has already been validated
func (parser *Parser) parseSection(action func(*Parser) (ok bool), againCondition func(*Parser, int) bool) (ok bool) {
	ok, _ = parser.parseSection2(action, againCondition)
	return ok
}
