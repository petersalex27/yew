package parser

import (
	"github.com/petersalex27/yew/common/stack"
	"github.com/petersalex27/yew/token"
)

// wrapper for *stack.Stack[int]
type indentStack struct{ *stack.Stack[int] }

// pops an indent off the stack--no errors are reported and no value is returned
func (stack indentStack) silentPop() { stack.Pop() }

// returns current indentation level
func (stack indentStack) CurrentIndent() (indentLevel int) {
	if !stack.Empty() {
		indentLevel, _ = stack.Peek()
	}
	return
}

// a meaningful indent is one that ...
//   - begins a line (indent token guaranteed to start line by lexer)
//   - contains non-whitespace characters after the indent
//   - contains non-comment tokens after indent
func (parser *tokenInfo) findMeaningfulIndent() (meaningfulIndent token.Token, success bool) {
	meaningfulIndent = parser.Peek()
	success = meaningfulIndent.Type == token.Indent
	if !success {
		return meaningfulIndent, false
	}

	again := true

	nextToken := meaningfulIndent

	// keep looping until meaningful indent was the last indent found
	//	at least one indent is guaranteed to have been found at this point b/c function returns
	//	earlier when no indent is found
	for again {
		meaningfulIndent = nextToken // store most recent indent
		_ = parser.Advance()         // move past most recent indent
		nextToken = parser.Peek()    // peek at next token

		if nextToken.Type == token.Comment {
			parser.dropComments()
			nextToken = parser.Peek()
		}

		if nextToken.Type == token.Newline {
			parser.drop()
			nextToken = parser.Peek() // peek at token after ignorable tokens
			if nextToken.Type != token.Indent {
				return meaningfulIndent, false // expected a meaningful indent to follow
			}
		}

		// loop again if an indent is found
		again = nextToken.Type == token.Indent
	}
	return meaningfulIndent, true
}

// use to advance input just like `findMeaningfulIndent` but also be given the token position of
// the next meaningful indent.
//
// this is useful for when you want to reverse input that's been read if the indent isn't what's
// wanted. To do so, do `parser.tokenPos = position`
//
// a meaningful indent is one that ...
//   - begins a line (indent token guaranteed to start line by lexer)
//   - contains non-whitespace characters after the indent
//   - contains non-comment tokens after indent
func (parser *tokenInfo) locateMeaningfulIndent() (position int, success bool) {
	initial := parser.tokenPos
	meaningfulIndent := parser.Peek()
	success = meaningfulIndent.Type == token.Indent
	if !success {
		return initial, false
	}

	again := true

	nextToken := meaningfulIndent

	position = initial

	// keep looping until meaningful indent was the last indent found
	//	at least one indent is guaranteed to have been found at this point b/c function returns
	//	earlier when no indent is found
	for again {
		position = parser.tokenPos
		meaningfulIndent = nextToken // store most recent indent
		_ = parser.Advance()         // move past most recent indent
		nextToken = parser.Peek()    // peek at next token

		if nextToken.Type == token.Comment {
			parser.dropComments()
			nextToken = parser.Peek()
		}

		if nextToken.Type == token.Newline {
			parser.drop()
			nextToken = parser.Peek() // peek at token after ignorable tokens
			if nextToken.Type != token.Indent {
				return initial, false // expected a meaningful indent to follow
			}
		}

		// loop again if an indent is found
		again = nextToken.Type == token.Indent
	}
	return position, true
}

// true iff argument is the same level as the current saved indent
func (parser *Parser) equalIndent(indent int) bool {
	return indent == parser.indentation.CurrentIndent()
}

// true iff argument is of a greater level than the current saved indent
func (parser *Parser) testIndentGT(indent int) bool {
	return indent > parser.indentation.CurrentIndent()
}