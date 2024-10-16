//go:build test
// +build test

package parser

import (
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
)

// TODO: add tests for ensuring visibility modifiers are correctly parsed and applied
func TestParseBody(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  body
		end   int
	}{
		{
			"typing",
			// x : x
			[]api.Token{id_x_tok, colon, id_x_tok},
			body_typing, -1,
		},
		{
			"type def",
			// MyId : x where MyId : x
			[]api.Token{id_MyId_tok, colon, id_x_tok, where, id_MyId_tok, colon, id_x_tok},
			body_typeDef, -1,
		},
		{
			"spec def",
			// spec MyId x where x : x
			[]api.Token{spec, id_MyId_tok, id_x_tok, where, id_x_tok, colon, id_x_tok},
			body_specDef, -1,
		},
		{
			"spec inst",
			// inst MyId x where x = x
			[]api.Token{inst, id_MyId_tok, id_x_tok, where, id_x_tok, equal, id_x_tok},
			body_specInst, -1,
		},
		{
			"type alias",
			// alias MyId = MyId
			[]api.Token{alias, id_MyId_tok, equal, id_MyId_tok},
			body_alias, -1,
		},
		{
			"syntax",
			// syntax `my` x = x
			[]api.Token{syntaxTok, raw_my_tok, id_x_tok, equal, id_x_tok},
			body_syntax, -1,
		},
		{
			"def",
			// x = x
			[]api.Token{id_x_tok, equal, id_x_tok},
			data.EMakes[body](defNode.asBodyElement()), -1,
		},
		{
			"annotated",
			// --@test\nx : x
			[]api.Token{annot, newline, id_x_tok, colon, id_x_tok},
			data.EMakes[body](annotTypingNode.asBodyElement()), -1,
		},
		{
			"multiple",
			// x : x\nx = x
			[]api.Token{id_x_tok, colon, id_x_tok, newline, id_x_tok, equal, id_x_tok},
			data.EMakes[body](typingNode.asBodyElement(), defNode.asBodyElement()), -1,
		},
		{
			// this was an issue in the old grammar; basically, one could end a definition w/ impossible
			// and directly begin another definition--although it's not a parsing issue, it's not
			// desirable from a code-health perspective
			`avoids "impossible chaining"`,
			[]api.Token{id_x_tok, impossibleTok, id_x_tok, colon, id_x_tok},
			//                                 ^-- should stop here
			body_defImpossible, 2,
		},
		{
			"unconsumed tail annotation - 0",
			// x : x [@test]
			[]api.Token{id_x_tok, colon, id_x_tok, annotOpenTok, id_test_tok, rbracket},
			//                                   ^-- should stop here
			data.EMakes[body](typingNode.asBodyElement()), 3,
		},
		{
			"unconsumed tail annotation - 1",
			// x : x\n[@test]
			[]api.Token{id_x_tok, colon, id_x_tok, newline, annotOpenTok, id_test_tok, rbracket},
			//                                   ^-- should stop here
			data.EMakes[body](typingNode.asBodyElement()), 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, data.Just(test.want), parseBody, test.end))
	}
}

func TestParseBodyElement_TestMaybeParseMainElement(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  mainElement
		end   int
	}{
		{
			"typing",
			// x : x
			[]api.Token{id_x_tok, colon, id_x_tok},
			typingNode, -1,
		},
		{
			"type def",
			// MyId : x where MyId : x
			[]api.Token{id_MyId_tok, colon, id_x_tok, where, id_MyId_tok, colon, id_x_tok},
			typeDefNode, -1,
		},
		{
			"spec def",
			// spec MyId x where x : x
			[]api.Token{spec, id_MyId_tok, id_x_tok, where, id_x_tok, colon, id_x_tok},
			specDefNode, -1,
		},
		{
			"spec inst",
			// inst MyId x where x = x
			[]api.Token{inst, id_MyId_tok, id_x_tok, where, id_x_tok, equal, id_x_tok},
			specInstNode, -1,
		},
		{
			"type alias",
			// alias MyId = MyId
			[]api.Token{alias, id_MyId_tok, equal, id_MyId_tok},
			aliasNode, -1,
		},
		{
			"syntax",
			// syntax `my` x = x
			[]api.Token{syntaxTok, raw_my_tok, id_x_tok, equal, id_x_tok},
			syntaxNode, -1,
		},
		{
			"def",
			// x = x
			[]api.Token{id_x_tok, equal, id_x_tok},
			defNode, -1,
		},
		{
			// this was an issue in the old grammar; basically, one could end a definition w/ impossible
			// and directly begin another definition--although it's not a parsing issue, it's not
			// desirable from a code-health perspective
			`avoids "impossible chaining"`,
			[]api.Token{id_x_tok, impossibleTok, id_x_tok, colon, id_x_tok},
			//                                 ^-- should stop here
			defImpossibleNode, 2,
		},
		{
			"unconsumed tail annotation - 0",
			// x : x [@test]
			[]api.Token{id_x_tok, colon, id_x_tok, annotOpenTok, id_test_tok, rbracket},
			//                                   ^-- should stop here
			typingNode, 3,
		},
		{
			"unconsumed tail annotation - 1",
			// x : x\n[@test]
			[]api.Token{id_x_tok, colon, id_x_tok, newline, annotOpenTok, id_test_tok, rbracket},
			//                                   ^-- should stop here
			typingNode, 3,
		},
	}

	t.Run("TestParseBodyElement", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want.asBodyElement(), parseBodyElement, test.end))
		}
	})

	t.Run("TestMaybeParseMainElem", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.name, maybeOutputFUT_endCheck(test.input, data.Just(test.want), maybeParseMainElement, test.end))
		}

		// run one final test that tests for annotations being applied
		test := struct {
			name  string
			input []api.Token
			want  mainElement
		}{
			"annotated",
			// --@test\nx : x
			[]api.Token{annot, newline, id_x_tok, colon, id_x_tok},
			annotTypingNode,
		}
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, data.Just(test.want), maybeParseMainElement, -1))
	})
}

func TestParseConstructorNameErrors(t *testing.T) {
	tests := []struct {
		name  string
		input api.Token
		want  string
	}{
		{
			"error - lower ident",
			id_x_tok,
			IllegalLowercaseConstructorName,
		},
		{
			"error - lower infix ident",
			infix_x_tok,
			IllegalLowercaseConstructorName,
		},
		{
			"error - method name",
			method_run_tok,
			IllegalMethodTypeConstructor,
		},
		{
			"error - non-name",
			alias,
			ExpectedTypeConstructorName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := typeConstructorNameError(test.input)
			if res != test.want {
				t.Errorf("failed: got %q, want %q", res, test.want)
			}
		})
	}
}

// rule:
//
//	```
//	constructor name = infix upper ident | upper ident | symbol | infix symbol ;
//	```
func TestParseConstructorName(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[name]
	}{
		{
			"infix upper ident",
			[]api.Token{infix_MyId_tok},
			data.Just(name_infix_MyId),
		},
		{
			"upper ident",
			[]api.Token{id_MyId_tok},
			data.Just(name_MyId),
		},
		{
			"symbol",
			[]api.Token{id_dollar_tok},
			data.Just(name_dollar),
		},
		{
			"infix symbol",
			[]api.Token{infix_dollar_tok},
			data.Just(name_infix_dollar),
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, maybeParseConstructorName, -1))
	}
}

// rule:
//
//	```
//	type constructor = constructor name, {{"\n"}, ",", {"\n"}, constructor name}, {"\n"}, ":", {"\n"}, type ;
//	```
func TestParseTypeConstructor(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.NonEmpty[typeConstructor]
	}{
		{
			"single - 00",
			[]api.Token{id_MyId_tok, colon, id_x_tok},
			singleConsNode,
		},
		{
			"single - 01",
			[]api.Token{id_MyId_tok, colon, newline, id_x_tok},
			singleConsNode,
		},
		{
			"single - 10",
			[]api.Token{id_MyId_tok, newline, colon, id_x_tok},
			singleConsNode,
		},
		{
			"single - 11",
			[]api.Token{id_MyId_tok, newline, colon, newline, id_x_tok},
			singleConsNode,
		},
		{
			"multiple - 00",
			[]api.Token{id_MyId_tok, comma, id_MyId_tok, colon, id_x_tok},
			multiConsNode,
		},
		{
			"multiple - 01",
			[]api.Token{id_MyId_tok, comma, newline, id_MyId_tok, colon, id_x_tok},
			multiConsNode,
		},
		{
			"multiple - 10",
			[]api.Token{id_MyId_tok, newline, comma, id_MyId_tok, colon, id_x_tok},
			multiConsNode,
		},
		{
			"multiple - 11",
			[]api.Token{id_MyId_tok, newline, comma, newline, id_MyId_tok, colon, id_x_tok},
			multiConsNode,
		},
		{
			"single - trailing comma",
			[]api.Token{id_MyId_tok, comma, colon, id_x_tok},
			data.Construct(makeCons(name_MyId, typ_x)),
		},
		{
			"multiple - trailing comma",
			[]api.Token{id_MyId_tok, comma, id_MyId_tok, comma, colon, id_x_tok},
			multiConsNode,
		},
	}

	for _, test := range tests {
		fut := fun.Bind1stOf2(parseTypeConstructor, data.Nothing[annotations]())
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, fut, -1))
	}
}

func TestParseDef(t *testing.T) {

}

// rule:
//
//	```
//	def body = (with clause | "=", {"\n"}, expr), [where clause] | "impossible" ;
//	```
func TestParseDefBody(t *testing.T) {
	tests := []struct {
		name string
		input []api.Token
		want defBody
		end int
	}{
		{
			"expr - 0",
			[]api.Token{equal, id_x_tok},
			defBodyNode, -1,
		},
		{
			"expr - 1",
			[]api.Token{equal, newline, id_x_tok},
			defBodyNode, -1,
		},
		{
			"expr - end correctness",
			[]api.Token{equal, id_x_tok, newline},
			//                         ^-- should stop here
			defBodyNode, 2,
		},
		{
			"with clause",
			[]api.Token{with, id_x_tok, of, id_x_tok, thickArrow, id_x_tok},
			data.EInr[defBody](data.EMakePair[defBodyPossible](data.Inl[expr](withClauseNode), data.Nothing[whereClause]())),
			-1,
		},
		{
			"double where",
			[]api.Token{
				// with x of (
				with, id_x_tok, of, lparen,
				//   x => x where x = x
					id_x_tok, thickArrow, id_x_tok, where, id_x_tok, equal, id_x_tok,
				// ) where x = x
				rparen, where, id_x_tok, equal, id_x_tok,
			},
			// with x of (x => x where x = x) where x = x
			data.EInr[defBody](data.EMakePair[defBodyPossible](data.Inl[expr](withClauseNodeWhere), data.Just(whereClauseNode))),
			-1,
		},
		{
			"impossible",
			[]api.Token{impossibleTok},
			data.EInl[defBody](impossibleNode), -1,
		},
		{
			"impossible - end correctness",
			[]api.Token{impossibleTok, id_x_tok},
			//                       ^-- should stop here
			data.EInl[defBody](impossibleNode), 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseDefBody, test.end))
	}
}

// rule:
//
//	```
//	deriving clause = "deriving", {"\n"}, deriving body ;
//	```
func TestParseDerivingClause(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  deriving
	}{
		{
			"single",
			[]api.Token{derivingTok, id_MyId_tok, id_x_tok},
			derivingNode,
		},
		{
			"newline",
			[]api.Token{derivingTok, newline, id_MyId_tok, id_x_tok},
			derivingNode,
		},
		{
			"enclosed single",
			[]api.Token{derivingTok, lparen, id_MyId_tok, id_x_tok, rparen},
			derivingNode,
		},
		{
			"enclosed single trailing comma",
			[]api.Token{derivingTok, lparen, id_MyId_tok, id_x_tok, comma, rparen},
			derivingNode,
		},
		{
			"multiple",
			[]api.Token{derivingTok, lparen, id_MyId_tok, id_x_tok, comma, id_MyId_tok, id_x_tok, rparen},
			derivingNode2,
		},
		{
			"multiple trailing comma",
			[]api.Token{derivingTok, lparen, id_MyId_tok, id_x_tok, comma, id_MyId_tok, id_x_tok, comma, rparen},
			derivingNode2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, data.Just(test.want), parseOptionalDerivingClause, -1))
	}
}

// rule:
//
//	```
//	syntax = "syntax", {"\n"}, syntax rule, {"\n"}, "=", {"\n"}, expr ;
//	```
//
// Ensures the following:
//  1. keywords are read
//  2. newlines are accounted for b/w the integrated production rules
func TestParseSyntax(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  syntax
	}{
		{
			"syntax - 000",
			[]api.Token{syntaxTok, raw_my_tok, id_x_tok, equal, id_x_tok},
			syntaxNode,
		},
		{
			"syntax - 001",
			[]api.Token{syntaxTok, raw_my_tok, id_x_tok, equal, newline, id_x_tok},
			syntaxNode,
		},
		{
			"syntax - 010",
			[]api.Token{syntaxTok, raw_my_tok, id_x_tok, newline, equal, id_x_tok},
			syntaxNode,
		},
		{
			"syntax - 011",
			[]api.Token{syntaxTok, raw_my_tok, id_x_tok, newline, equal, newline, id_x_tok},
			syntaxNode,
		},
		{
			"syntax - 100",
			[]api.Token{syntaxTok, newline, raw_my_tok, id_x_tok, equal, id_x_tok},
			syntaxNode,
		},
		{
			"syntax - 101",
			[]api.Token{syntaxTok, newline, raw_my_tok, id_x_tok, equal, newline, id_x_tok},
			syntaxNode,
		},
		{
			"syntax - 110",
			[]api.Token{syntaxTok, newline, raw_my_tok, id_x_tok, newline, equal, id_x_tok},
			syntaxNode,
		},
		{
			"syntax - 111",
			[]api.Token{syntaxTok, newline, raw_my_tok, id_x_tok, newline, equal, newline, id_x_tok},
			syntaxNode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseSyntax, -1))
	}
}

// rule:
//
//	```
//	binding syntax ident = "{", {"\n"}, ident, {"\n"}, "}" ;
func TestParseBindingSyntaxIdent(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  syntaxSymbol
	}{
		{
			"binding syntax lower ident",
			[]api.Token{lbrace, id_x_tok, rbrace},
			bindingIdSymNode,
		},
		{
			"binding syntax upper ident",
			[]api.Token{lbrace, id_MyId_tok, rbrace},
			data.Inl[syntaxRawKeyword](makeBindingSyntaxRuleIdent(upperId)),
		},
		{
			"binding syntax ident - 00",
			[]api.Token{lbrace, id_x_tok, rbrace},
			bindingIdSymNode,
		},
		{
			"binding syntax ident - 01",
			[]api.Token{lbrace, id_x_tok, newline, rbrace},
			bindingIdSymNode,
		},
		{
			"binding syntax ident - 10",
			[]api.Token{lbrace, newline, id_x_tok, rbrace},
			bindingIdSymNode,
		},
		{
			"binding syntax ident - 11",
			[]api.Token{lbrace, newline, id_x_tok, newline, rbrace},
			bindingIdSymNode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseBindingSyntaxIdent, -1))
	}
}

// rule:
//
//	```
//	syntax symbol = ident | "{", {"\n"}, ident, {"\n"}, "}" | raw keyword ;
//	raw keyword = ? RAW STRING OF JUST A VALID NON INFIX ident OR symbol ? ;
//	```
func TestMaybeParseSyntaxSymbol(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[syntaxSymbol]
	}{
		{
			"ident",
			[]api.Token{id_x_tok},
			data.Just(idSymNode),
		},
		{
			"binging syntax ident",
			[]api.Token{lbrace, id_x_tok, rbrace},
			data.Just(bindingIdSymNode),
		},
		{
			"raw keyword",
			[]api.Token{raw_my_tok},
			data.Just(rawSym),
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, maybeParseSyntaxSymbol, -1))
	}
}

// rule:
//
//	```
//	syntax rule = {syntax symbol, {"\n"}}, raw keyword, {{"\n"}, syntax symbol} ;
//	```
func TestParseSyntaxRule(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  syntaxRule
	}{
		{
			"key",
			[]api.Token{raw_my_tok},
			data.EConstruct[syntaxRule](rawSym),
		},
		{
			"id,key - 0",
			[]api.Token{id_x_tok, raw_my_tok},
			data.EConstruct[syntaxRule](idSymNode, rawSym),
		},
		{
			"id,key - 1",
			[]api.Token{id_x_tok, newline, raw_my_tok},
			data.EConstruct[syntaxRule](idSymNode, rawSym),
		},
		{
			"key,id - 0",
			[]api.Token{raw_my_tok, id_x_tok},
			data.EConstruct[syntaxRule](rawSym, idSymNode),
		},
		{
			"key,id - 1",
			[]api.Token{raw_my_tok, newline, id_x_tok},
			data.EConstruct[syntaxRule](rawSym, idSymNode),
		},
		{
			"id,key,id - 00",
			[]api.Token{id_x_tok, raw_my_tok, id_x_tok},
			data.EConstruct[syntaxRule](idSymNode, rawSym, idSymNode),
		},
		{
			"id,key,id - 01",
			[]api.Token{id_x_tok, raw_my_tok, newline, id_x_tok},
			data.EConstruct[syntaxRule](idSymNode, rawSym, idSymNode),
		},
		{
			"id,key,id - 10",
			[]api.Token{id_x_tok, newline, raw_my_tok, id_x_tok},
			data.EConstruct[syntaxRule](idSymNode, rawSym, idSymNode),
		},
		{
			"id,key,id - 11",
			[]api.Token{id_x_tok, newline, raw_my_tok, newline, id_x_tok},
			data.EConstruct[syntaxRule](idSymNode, rawSym, idSymNode),
		},
	}
	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseSyntaxRule, -1))
	}
}

// rule:
//
//	```
//	type alias = "alias", {"\n"}, name, {"\n"}, "=", {"\n"}, type ;
//	```
func TestParseTypeAlias(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  typeAlias
	}{
		{
			"000",
			[]api.Token{alias, id_MyId_tok, equal, id_MyId_tok},
			aliasNode,
		},
		{
			"001",
			[]api.Token{alias, id_MyId_tok, equal, newline, id_MyId_tok},
			aliasNode,
		},
		{
			"010",
			[]api.Token{alias, id_MyId_tok, newline, equal, id_MyId_tok},
			aliasNode,
		},
		{
			"011",
			[]api.Token{alias, id_MyId_tok, newline, equal, newline, id_MyId_tok},
			aliasNode,
		},
		{
			"100",
			[]api.Token{alias, newline, id_MyId_tok, equal, id_MyId_tok},
			aliasNode,
		},
		{
			"101",
			[]api.Token{alias, newline, id_MyId_tok, equal, newline, id_MyId_tok},
			aliasNode,
		},
		{
			"110",
			[]api.Token{alias, newline, id_MyId_tok, newline, equal, id_MyId_tok},
			aliasNode,
		},
		{
			"111",
			[]api.Token{alias, newline, id_MyId_tok, newline, equal, newline, id_MyId_tok},
			aliasNode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseTypeAlias, -1))
	}
}

// rule:
//
//	```
//	type def = typing, {"\n"}, "where", {"\n"}, type def body, [{"\n"}, deriving clause] ;
//	```
func TestParseTypeDef(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  mainElement
	}{
		{
			"00",
			[]api.Token{id_MyId_tok, colon, id_x_tok, where, id_MyId_tok, colon, id_x_tok},
			typeDefNode,
		},
		{
			"01",
			[]api.Token{id_MyId_tok, colon, id_x_tok, where, newline, id_MyId_tok, colon, id_x_tok},
			typeDefNode,
		},
		{
			"10",
			[]api.Token{id_MyId_tok, colon, id_x_tok, newline, where, id_MyId_tok, colon, id_x_tok},
			typeDefNode,
		},
		{
			"11",
			[]api.Token{id_MyId_tok, colon, id_x_tok, newline, where, newline, id_MyId_tok, colon, id_x_tok},
			typeDefNode,
		},
		{
			"with deriving - 00",
			[]api.Token{id_MyId_tok, colon, id_x_tok, where, id_MyId_tok, colon, id_x_tok, derivingTok, id_MyId_tok, id_x_tok},
			typeDefNodeWithDeriving,
		},
		{
			"with deriving - 01",
			[]api.Token{id_MyId_tok, colon, id_x_tok, where, id_MyId_tok, colon, id_x_tok, derivingTok, newline, id_MyId_tok, id_x_tok},
			typeDefNodeWithDeriving,
		},
		{
			"with deriving - 10",
			[]api.Token{id_MyId_tok, colon, id_x_tok, where, id_MyId_tok, colon, id_x_tok, newline, derivingTok, id_MyId_tok, id_x_tok},
			typeDefNodeWithDeriving,
		},
		{
			"with deriving - 11",
			[]api.Token{id_MyId_tok, colon, id_x_tok, where, id_MyId_tok, colon, id_x_tok, newline, derivingTok, newline, id_MyId_tok, id_x_tok},
			typeDefNodeWithDeriving,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseTypeDefOrTyping, -1))
	}
}

func TestParseTypeDefBody(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  typeDefBody
		end   int
	}{
		{
			"impossible",
			[]api.Token{impossibleTok},
			data.Inr[data.NonEmpty[typeConstructor]](impossibleNode), -1,
		},
		{
			"non-grouped",
			[]api.Token{id_MyId_tok, colon, id_x_tok},
			data.Inl[impossible](singleConsNode), -1,
		},
		{
			"enclosed - 00",
			[]api.Token{lparen, id_MyId_tok, colon, id_x_tok, rparen},
			data.Inl[impossible](singleConsNode), -1,
		},
		{
			"enclosed - 01",
			[]api.Token{lparen, id_MyId_tok, colon, id_x_tok, newline, rparen},
			data.Inl[impossible](singleConsNode), -1,
		},
		{
			"enclosed - 10",
			[]api.Token{lparen, newline, id_MyId_tok, colon, id_x_tok, rparen},
			data.Inl[impossible](singleConsNode), -1,
		},
		{
			"enclosed - 11",
			[]api.Token{lparen, newline, id_MyId_tok, colon, id_x_tok, newline, rparen},
			data.Inl[impossible](singleConsNode), -1,
		},
		{
			"enclosed - multiple",
			[]api.Token{lparen, id_MyId_tok, colon, id_x_tok, newline, id_MyId_tok, colon, id_x_tok, rparen},
			data.Inl[impossible](multiConsNode), -1,
		},
		{
			"input end correctness - 1",
			[]api.Token{impossibleTok, id_MyId_tok, colon, id_x_tok},
			//                       ^-- should stop here
			data.Inr[data.NonEmpty[typeConstructor]](impossibleNode), 1,
		},
		{
			"input end correctness - 1",
			[]api.Token{id_MyId_tok, colon, id_x_tok, newline, id_MyId_tok, colon, id_x_tok},
			//                                      ^-- should stop here
			data.Inl[impossible](singleConsNode), 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseTypeDefBody, test.end))
	}
}

// not much to test here, just make sure the name parse, colon parse, and type parse are
// correctly sequenced to allow for newlines in appropriate places
func TestParseTyping(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  typing
	}{
		{
			"00",
			[]api.Token{id_x_tok, colon, id_x_tok},
			typingNode,
		},
		{
			"01",
			[]api.Token{id_x_tok, colon, newline, id_x_tok},
			typingNode,
		},
		{
			"10",
			[]api.Token{id_x_tok, newline, colon, id_x_tok},
			typingNode,
		},
		{
			"11",
			[]api.Token{id_x_tok, newline, colon, newline, id_x_tok},
			typingNode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseTypeSig, -1))
	}
}

func TestParseOptionalVisibility(t *testing.T) {
	
}

// rule:
//
//	```
//	where body = main elem | "(", {"\n"}, main elem, {then, main elem}, {"\n"}, ")" ;
//	```
func TestParseWhereBody(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  whereClause
	}{
		{
			"single",
			[]api.Token{id_x_tok, colon, id_x_tok},
			data.EConstruct[whereClause](mainElement(typingNode)),
		},
		{
			"enclosed - 00",
			[]api.Token{lparen, id_x_tok, colon, id_x_tok, rparen},
			data.EConstruct[whereClause](mainElement(typingNode)),
		},
		{
			"enclosed - 01",
			[]api.Token{lparen, id_x_tok, colon, id_x_tok, newline, rparen},
			data.EConstruct[whereClause](mainElement(typingNode)),
		},
		{
			"enclosed - 10",
			[]api.Token{lparen, newline, id_x_tok, colon, id_x_tok, rparen},
			data.EConstruct[whereClause](mainElement(typingNode)),
		},
		{
			"enclosed - 11",
			[]api.Token{lparen, newline, id_x_tok, colon, id_x_tok, newline, rparen},
			data.EConstruct[whereClause](mainElement(typingNode)),
		},
		{
			"multiple",
			[]api.Token{lparen, newline, id_x_tok, colon, id_x_tok, newline, id_x_tok, colon, id_x_tok, newline, rparen},
			data.EConstruct[whereClause](mainElement(typingNode), mainElement(typingNode)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseWhereBody, -1))
	}
}

// other cases are covered by `TestParseWhereBody`
//
// parseOptionalWhereClause just grabs where token and calls parseWhereBody
func TestParseOptionalWhereClause(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[whereClause]
		end   int
	}{
		{
			"empty",
			[]api.Token{},
			data.Nothing[whereClause](),
			-1,
		},
		{
			"non-empty",
			[]api.Token{where, id_x_tok, colon, id_x_tok},
			data.Just(data.EConstruct[whereClause](mainElement(typingNode))),
			-1,
		},
		{
			"non-where clause",
			[]api.Token{id_x_tok, colon, id_x_tok},
			data.Nothing[whereClause](),
			0,
		},
		{
			"where clause, followed by more",
			[]api.Token{where, id_x_tok, colon, id_x_tok, newline, id_x_tok, colon, id_x_tok},
			data.Just(data.EConstruct[whereClause](mainElement(typingNode))),
			4, // should NOT read newline
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT(test.input, test.want, parseOptionalWhereClause))
	}
}

// rule:
//
//	```
//	with clause = "with", {"\n"}, pattern, {"\n"}, "of", {"\n"}, with clause arms ;
//	```
func TestParseWithClause(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  withClause
	}{
		{
			"000",
			[]api.Token{with, id_x_tok, of, id_x_tok, thickArrow, id_x_tok},
			withClauseNode,
		},
		{
			"001",
			[]api.Token{with, id_x_tok, of, newline, id_x_tok, thickArrow, id_x_tok},
			withClauseNode,
		},
		{
			"010",
			[]api.Token{with, id_x_tok, newline, of, id_x_tok, thickArrow, id_x_tok},
			withClauseNode,
		},
		{
			"011",
			[]api.Token{with, id_x_tok, newline, of, newline, id_x_tok, thickArrow, id_x_tok},
			withClauseNode,
		},
		{
			"100",
			[]api.Token{with, newline, id_x_tok, of, id_x_tok, thickArrow, id_x_tok},
			withClauseNode,
		},
		{
			"101",
			[]api.Token{with, newline, id_x_tok, of, newline, id_x_tok, thickArrow, id_x_tok},
			withClauseNode,
		},
		{
			"110",
			[]api.Token{with, newline, id_x_tok, newline, of, id_x_tok, thickArrow, id_x_tok},
			withClauseNode,
		},
		{
			"111",
			[]api.Token{with, newline, id_x_tok, newline, of, newline, id_x_tok, thickArrow, id_x_tok},
			withClauseNode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseWithClause, -1))
	}
}

// nothing to test, just wraps a call to parseGroup--if other functions for with clause work, this one is correct
// under the assumption that parseGroup is correct
// - func TestParseWithClauseArms(t *testing.T)

// rule:
//
//	```
//	with clause arm = [view refined pattern, {"\n"}], pattern, {"\n"}, def body thick arrow ;
//	view refined pattern = pattern, {"\n"}, "|" ;
//	```
func TestMaybeParseWithClauseArm(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[withClauseArm]
	}{
		{
			"no view refined - 0",
			[]api.Token{id_x_tok, thickArrow, id_x_tok},
			data.Just(withClauseArmNode),
		},
		{
			"no view refined - 1",
			[]api.Token{id_x_tok, newline, thickArrow, id_x_tok},
			data.Just(withClauseArmNode),
		},
		{
			"view refined - 000",
			[]api.Token{id_x_tok, bar, id_x_tok, thickArrow, id_x_tok},
			data.Just(withClauseVRNode),
		},
		{
			"view refined - 001",
			[]api.Token{id_x_tok, bar, id_x_tok, newline, thickArrow, id_x_tok},
			data.Just(withClauseVRNode),
		},
		{
			"view refined - 010",
			[]api.Token{id_x_tok, bar, newline, id_x_tok, thickArrow, id_x_tok},
			data.Just(withClauseVRNode),
		},
		{
			"view refined - 011",
			[]api.Token{id_x_tok, bar, newline, id_x_tok, newline, thickArrow, id_x_tok},
			data.Just(withClauseVRNode),
		},
		{
			"view refined - 100",
			[]api.Token{id_x_tok, newline, bar, id_x_tok, thickArrow, id_x_tok},
			data.Just(withClauseVRNode),
		},
		{
			"view refined - 101",
			[]api.Token{id_x_tok, newline, bar, id_x_tok, newline, thickArrow, id_x_tok},
			data.Just(withClauseVRNode),
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, maybeParseWithClauseArm, -1))
	}
}
