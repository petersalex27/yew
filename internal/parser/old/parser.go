// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package parser

import (
	"fmt"
	"io"
	"os"

	"github.com/llir/llvm/ir"
	"github.com/petersalex27/yew/internal/common/stack"
	"github.com/petersalex27/yew/internal/common/table"
	"github.com/petersalex27/yew/internal/errors"
	"github.com/petersalex27/yew/internal/source"
	"github.com/petersalex27/yew/internal/token"
	"github.com/petersalex27/yew/internal/types"
)

type tokenInfo struct {
	// token stream, from lexer
	tokens []token.Token
	// current position in field `tokens`
	tokenPos     int
	saveComments bool
	keepNewlines bool
	comments     []token.Token
}

type symbolSaver struct {
	mutual bool
	allow  *stack.Stack[SyntaxClass]
	cls    SyntaxClass

	decls    *stack.SaveStack[DeclarationElem]
	bindings *stack.SaveStack[BindingElem]
	types    *stack.SaveStack[DataTypeElem]
	typings  *stack.SaveStack[TypeElem]
	traits   *stack.SaveStack[SpecElem]
	inst     *stack.SaveStack[InstanceElem]
	elems    *stack.SaveStack[SymbolicElem]
}

func (parser *Parser) writeDecls(decls []DeclarationElem) bool {
	for _, decl := range decls {
		if !parser.write(decl) {
			return false
		}
	}
	return true
}

func annotate[T SyntacticElem](parser *Parser, elem T) bool {
	for _, annotation := range parser.annotations {
		if !annotation(parser, &elem) {
			return false
		}
	}
	return true
}

// annotations are applied here
func (parser *Parser) write(elem SyntacticElem) (ok bool) {
	switch elem := elem.(type) {
	case DeclarationElem:
		if ok = annotate(parser, elem); ok {
			parser.saver.decls.Push(elem)
		}
	case BindingElem:
		if ok = annotate(parser, elem); ok {
			parser.saver.bindings.Push(elem)
		}
	case DataTypeElem:
		if ok = annotate(parser, elem); ok {
			parser.saver.types.Push(elem)
		}
	case SpecElem:
		if ok = annotate(parser, elem); ok {
			parser.saver.traits.Push(elem)
		}
	case InstanceElem:
		if ok = annotate(parser, elem); ok {
			parser.saver.inst.Push(elem)
		}
	case TypeElem:
		if ok = annotate(parser, elem); ok {
			parser.saver.typings.Push(elem)
		}
		return ok
	case SymbolicElem: // SymbolicElem is a catch-all for all other cases
		if ok = annotate(parser, elem); !ok {
			parser.saver.elems.Push(elem)
		}
		return ok
	default:
		panic(fmt.Sprintf("bug: unhandled case, %v", elem))
	}

	if !ok {
		return
	} else if s, ok := elem.(SymbolicElem); ok {
		parser.saver.elems.Push(s)
	} else {
		panic("bug: unhandled case, non-SymbolicElem escaped switch")
	}
	return ok
}

func (parser *Parser) PrintResult() {
	count := parser.saver.elems.GetCount()
	res, _ := parser.saver.elems.MultiCheck(int(count))
	for _, el := range res {
		fmt.Fprintf(os.Stderr, "\t%v\n", el)
	}
}

// information related to the parser's first pass
type firstPassInfo struct {
	// parse-stack for indentation
	indentation indentStack
	// syntactic sections
	Sections []SyntacticElem
	saver    symbolSaver
}

// general info about the parser and what it's parsing
type generalInfo struct {
	// source code
	src source.SourceCode
	// messages: errors, warnings, logs, etc.
	messages []errors.ErrorMessage

	// flags whether or not parser encountered an error
	panicking bool
}

// a term with additional precedence, associativity, and arity information for parsing
type termElem struct {
	NodeType
	types.Term
	//Kind types.Type
	termInfo
	Start, End int
}

func (t termElem) Pos() (int, int) {
	return t.Start, t.End
}

func (t termElem) String() string {
	isNil := t.Term == nil
	switch t.NodeType {
	case identType, applicationType, intConstType, charConstType, floatConstType, stringConstType, lambdaType, listExprType, tupleType, tupleExprType, pairsType, listingType:
		if isNil {
			return "?"
		}
		return fmt.Sprintf("%v", t.Term)
	case implicitType:
		if isNil {
			return "?"
		}
		if t.arity == 0 {
			return fmt.Sprintf("{%v : %v}", t.Term, types.GetKind(&t.Term))
		} else {
			return fmt.Sprintf("{%v : ?}", t.Term)
		}
	case functionType:
		if isNil {
			return "? -> ?"
		}
		if t.arity == 0 {
			return fmt.Sprintf("%v", t.Term)
		} else {
			return fmt.Sprintf("%v -> ?", t.Term)
		}
	case labeledFunctionType, implicitFunctionType:
		if isNil {
			return "? -> ?"
		}
		if t.arity == 0 {
			return fmt.Sprintf("%v", t.Term)
		} else if t.NodeType == labeledFunctionType {
			return fmt.Sprintf("(%v) -> ?", types.TypingString(t.Term))
		} else {
			return fmt.Sprintf("{%v} -> ?", types.TypingString(t.Term))
		}
	case typingType:
		if isNil {
			return "?"
		}
		if t.arity == 0 {
			return fmt.Sprintf("%v : %v", t.Term, types.GetKind(&t.Term))
		} else {
			return fmt.Sprintf("%v : ?", t.Term)
		}
	default:
		panic("bug: unhandled case")
	}
}

type annotationMacro struct {
	replace []token.Token
}

type metaType uint

const (
	affixInfoMetaType metaType = iota
	macroDefMetaType
	macroUseMetaType
)

// Structure for information used during parsing
type Parser struct {
	// declarations, each table represents the scope of some binding (or collection of non-interfering
	// bindings when possible)
	declarations *declMultiTable
	defining     *stack.Stack[definitionParent]
	// TODO:
	// 	- need to allow for polymorphism--not just from traits
	// 	- need to allow multiple definitions of the same function, matching different patterns
	// 	- the order of the definitions _must_ be preserved
	definitions     map[string]definitionParent
	classifications map[string][]SyntacticElem
	//locals          *declMultiTable
	//declarations map[string]*Declaration
	generalInfo
	// information about tokens, input from lexer
	tokenInfo
	// imported modules/packages
	imports ImportTable
	// tracks scoped visibility modifier
	visibility *stack.Stack[Visibility]
	// environment for type checking
	env         *types.Environment
	mutualBlock *token.Token
	// information for first pass
	firstPassInfo
	// for passing between calls to (Term) Parse(..)
	termPasser *stack.Stack[termElem]
	// previous term, used to restore previous term on certain operations
	termMemory *stack.Stack[termElem]
	// llvm-ir module: parser outputs data here
	mod *ir.Module
	// true iff in top level
	inTop          bool
	inParent       string
	parsingTypeSig bool
	allowModality  bool
	annotations    []func(*Parser, any) bool
	debug_info_parser
}

// signal that parser is not in top level, parsing local symbols
func (parser *Parser) parsingLocal(name string) func(*Parser) {
	old := parser.inParent
	parser.inParent = old + "_" + name
	return func(p *Parser) {
		p.inParent = old
	}
}

func (parser *Parser) findDeclAsTerm(key stringPos) (term termElem, found bool) {
	declp, ok := parser.declarations.Find(key)
	found = ok
	if !found {
		return
	}

	term.termInfo = *declp.termInfo
	term.Term, _, found = parser.env.Get(key)
	// change value to current occurrence, leaving one in map the same
	term.Start, term.End = key.Pos()
	term.NodeType = identType
	return term, true
}

// func (parser *Parser) findInLocals(key fmt.Stringer) (decl Declaration, found bool) {
// 	var declPtr *Declaration
// 	if declPtr, found = parser.locals.Find(key); found {
// 		decl = *declPtr
// 	}
// 	return
// }

func (parser *Parser) ExploringTopLevel() bool {
	return parser.inTop
}

// flags that parser is inside mutual block
func (parser *Parser) enterMutualBlock(mutualToken token.Token) (ok bool) {
	if ok = parser.mutualBlock == nil; !ok {
		parser.error(InvalidDuplicateMutualBlock)
		return
	}
	parser.mutualBlock = &mutualToken
	return
}

// flags that parser is no longer in mutual block
func (parser *Parser) exitMutualBlock() {
	parser.mutualBlock = nil
}

func (parser *Parser) CurrentVisibility() Visibility {
	if parser.visibility.Empty() {
		return Private
	}
	v, _ := parser.visibility.Peek()
	return v
}

func (parser *Parser) setVisibility(v Visibility) {
	parser.visibility.Push(v)
}

func (parser *Parser) restoreVisibility() {
	if !parser.visibility.Empty() {
		_, _ = parser.visibility.Pop()
	}
}

func (parser *tokenInfo) saveComment(comment token.Token) {
	if parser.saveComments {
		parser.comments = append(parser.comments, comment)
	}
}

// function intended for debugging
func (parser *Parser) PrintImports(w io.Writer) {
	if len(parser.imports) == 0 {
		fmt.Fprintf(w, "(no imports)")
		return
	}

	fmt.Fprintf(w, "import\n")
	for _, im := range parser.imports {
		fmt.Fprintf(w, "  %s as %s\n", im.Lookup, im.Id)
	}
}

func initSaver() (saver symbolSaver) {
	saver.decls = stack.NewSaveStack[DeclarationElem](8)
	saver.bindings = stack.NewSaveStack[BindingElem](8)
	saver.typings = stack.NewSaveStack[TypeElem](8)
	saver.types = stack.NewSaveStack[DataTypeElem](8)
	saver.traits = stack.NewSaveStack[SpecElem](8)
	saver.elems = stack.NewSaveStack[SymbolicElem](32)
	saver.allow = stack.NewStack[SyntaxClass](8)
	saver.inst = stack.NewSaveStack[InstanceElem](8)
	return
}

func Initialize(src source.SourceCode, saveComments bool) (parser *Parser) {
	parser = new(Parser)
	parser.defining = stack.NewStack[definitionParent](8)
	parser.inParent = "" // top level, no module known yet even
	parser.inTop = true
	parser.messages = make([]errors.ErrorMessage, 0)
	parser.src = src
	//parser.action = stack.NewStack[Action](8)
	//parser.terms.SaveStack = stack.NewSaveStack[termElem](8)
	parser.saver = initSaver()
	parser.visibility = stack.NewStack[Visibility](2)
	parser.saveComments = saveComments
	if saveComments {
		parser.comments = make([]token.Token, 0, 32)
	}
	parser.annotations = make([]func(*Parser, any) bool, 0, 8)
	parser.definitions = make(map[string]definitionParent)
	// initialize annotation classifiers
	parser.classifications = make(map[string][]SyntacticElem)
	// initialize test classification
	parser.classifications["test"] = make([]SyntacticElem, 0)
	parser.termPasser = stack.NewStack[termElem](4)
	parser.termMemory = stack.NewStack[termElem](4)
	parser.env = types.NewEnvironment()
	parser.indentation = indentStack{stack.NewStack[int](8)}
	parser.imports = make(map[string]Import)
	parser.mod = ir.NewModule()
	parser.declarations = table.NewMultiTable[fmt.Stringer, *declaration](8)
	//parser.locals = table.NewMultiTable[fmt.Stringer, termInfo](8)
	builtinsPath := "prelude/builtin.yew" // TODO
	parser.declareBuiltin(builtinsPath)
	return
}

func Init(src source.SourceCode) (parser *Parser) {
	return Initialize(src, false)
}

func (parser *Parser) Panicking() bool {
	return parser.panicking
}

func (parser *Parser) openSection() {

}

func (parser *tokenInfo) getNextToken() (tok token.Token, offset int) {
	for offset = 0; ; {
		if parser.tokenPos >= len(parser.tokens) {
			return endOfTokensToken(), offset
		}

		tok := parser.tokens[parser.tokenPos+offset]
		offset++
		if !parser.keepNewlines && tok.Type == token.Newline {
			continue
		}
		return tok, offset
	}
}

func (parser *tokenInfo) Advance() token.Token {
	tok, offset := parser.getNextToken()
	parser.tokenPos += offset
	return tok
}

// returns next token but does not advance past it
func (parser *tokenInfo) Peek() token.Token {
	tok, _ := parser.getNextToken()
	return tok
}

// adds a message to parser's internal messages slice
func (parser *generalInfo) addMessage(e errors.ErrorMessage) {
	parser.messages = append(parser.messages, e)
	parser.panicking = parser.panicking || e.IsFatal()
}

func (parser *generalInfo) Messages() []errors.ErrorMessage {
	return parser.messages
}

func (parser *generalInfo) FlushMessages() []errors.ErrorMessage {
	messages := parser.Messages()
	parser.messages = []errors.ErrorMessage{}
	return messages
}

func (parser *tokenInfo) drop() (newlines int) {
	next := parser.Peek()
	ty := next.Type
	for ty == token.Newline || ty == token.Comment {
		if ty == token.Newline {
			newlines++
		} else {
			parser.saveComment(next)
		}

		_ = parser.Advance()
		next = parser.Peek()
		ty = next.Type
	}
	return
}

func (parser *Parser) dropAndAdvanceGreaterIndent() (greaterOrNoneDropped bool) {
	// drop comments and newlines
	newlines := parser.drop()
	// locate position of next meaningful indent
	pos, located := parser.locateMeaningfulIndent()
	if !located {
		// no indent found, just (maybe) dropped newlines and comments
		return newlines == 0
	}

	// grab indent from tokens
	indent := parser.tokens[pos]
	level := len(indent.Value)
	// test if this indent has a greater level than the section indent level
	if greaterOrNoneDropped = parser.testIndentGT(level); !greaterOrNoneDropped {
		// it does not, rewind to meaningful indent and report no greater indent found
		parser.tokenPos = pos
	} // else, it does; indent remains dropped
	return
}

// returns an "End" token
func endOfTokensToken() token.Token {
	return token.Token{Type: token.EndOfTokens}
}

// returns true iff annotation at `index` allows line-breaks between args
func (parser *tokenInfo) allowBreaks(index int) bool {
	for index = index - 1; index >= 0; index-- {
		ty := parser.tokens[index].Type
		if ty == token.LeftParen {
			return true
		}
		if ty != token.Newline {
			return false
		}
	}
	return false
}

func (parser *tokenInfo) getTok_breakable(i int, ty token.Type, eatBreaks bool) (tok token.Token, iNew int) {
	length := len(parser.tokens)
	iNew = i
	if iNew >= length {
		return
	}

	if eatBreaks {
		for iNew < length {
			tok = parser.tokens[i]
			if tok.Type != token.Newline {
				break
			}
			iNew++
		}
	}
	tok = parser.tokens[iNew]
	if tok.Type != ty {
		iNew = -1
		return
	}
	iNew++
	return
}

func (parser *Parser) Load(tokens []token.Token) {
	parser.tokens = tokens
}
