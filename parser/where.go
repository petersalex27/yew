// =================================================================================================
// Alex Peters - January 25, 2024
//
// Parses a 'where' contextualization. Use "ParseWhere" to parse, other functions are helpers
// =================================================================================================

package parser

import (
	"github.com/petersalex27/yew/common"
	"github.com/petersalex27/yew/token"
)

type WhereMap = common.Table[Ident, ExprNode]

// parses 'where' and its name-expression binding(s)
func (parser *Parser) ParseWhere() (wmap *WhereMap, ok bool) {
	if _, ok = parser.whereToken(); !ok {
		return
	}

	// TODO: where block
	endOptional := parser.StartOptional()
	_, isBlock := parser.getToken(token.LeftBrace, "")
	endOptional()

	wmap = common.MakeTable[Ident, ExprNode](4)

	// parse all name bindings
	again := true
	for ; again; again = isBlock && again {
		again = parser.whereBinding(wmap)
	}

	// allow empty where blocks, i.e.,
	//	stuff where {}
	// but do not allow 'where' not followed by anything
	//	stuff where
	if ok = !(wmap.Len() < 1 && !isBlock); !ok {
		parser.error(IllegalWhere)
		return
	}
	return
}

// parse where binding
func (parser *Parser) whereBinding(wmap *WhereMap) (ok bool) {
	var id Ident
	if id, ok = parser.parseFunctionName(); !ok {
		return
	}

	// TODO: functions w/ args

	if _, ok = parser.equalToken(); !ok {
		return
	}

	var expression ExprNode
	if expression, ok = parser.parseExpression(); ok {
		wmap.Map(id, expression)
	}
	return
}
