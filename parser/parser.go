// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package parser

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/llir/llvm/ir"
	"github.com/petersalex27/yew/common/stack"
	"github.com/petersalex27/yew/common/table"
	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/source"
	"github.com/petersalex27/yew/token"
	"github.com/petersalex27/yew/types"
)

type tokenInfo struct {
	// token stream, from lexer
	tokens []token.Token
	// current position in field `tokens`
	tokenPos     int
	saveComments bool
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
	traits   *stack.SaveStack[TraitElem]
	inst     *stack.SaveStack[InstanceElem]
	elems    *stack.SaveStack[SymbolicElem]
}

//func (parser *Parser) allow(c SyntaxClass) bool { return parser.saver.cls&c != 0 }

func (parser *Parser) writeDecl(decl DeclarationElem) {
	// if !parser.allow(DeclClass) {
	// 	parser.errorOnElem(IllegalDeclaration, decl)
	// }
	parser.saver.decls.Push(decl)
	parser.saver.elems.Push(decl)
}

func (parser *Parser) writeBinding(binding BindingElem) {
	// if !parser.allow(FuncClass) {
	// 	parser.errorOnElem(IllegalBinding, binding)
	// }
	parser.saver.bindings.Push(binding)
	parser.saver.elems.Push(binding)
}

func (parser *Parser) writeDataType(data DataTypeElem) {
	// if !parser.allow(TypeClass) {
	// 	parser.errorOnElem(IllegalDataType, data)
	// }
	parser.saver.types.Push(data)
	parser.saver.elems.Push(data)
}

func (parser *Parser) writeTyping(typ TypeElem) {
	parser.saver.typings.Push(typ)
}

func (parser *Parser) writeTrait(trait TraitElem) {
	// if !parser.allow(TraitClass) {
	// 	parser.errorOnElem(IllegalTrait, trait)
	// }
	parser.saver.traits.Push(trait)
	parser.saver.elems.Push(trait)
}

func (parser *Parser) writeInstance(inst InstanceElem) {
	// if !parser.allow(InstanceClass) {
	// 	parser.errorOnElem(IllegalInstance, inst)
	// }
	parser.saver.inst.Push(inst)
	parser.saver.elems.Push(inst)
}

func (parser *Parser) PrintResult() {
	count := parser.saver.elems.GetCount()
	res, _ := parser.saver.elems.MultiCheck(int(count))
	for _, el := range res {
		fmt.Fprintf(os.Stderr, "\t%v\n", el)
	}
}

type firstPassInfo struct {
	// parse-stack for indentation
	indentation indentStack
	// syntactic sections
	Sections []SyntacticElem
	saver    symbolSaver
}

type generalInfo struct {
	// source code
	src source.SourceCode
	// messages: errors, warnings, logs, etc.
	messages []errors.ErrorMessage

	// flags whether or not parser encountered an error
	panicking bool
}

type termElem struct {
	Term
	termInfo
}

type stackInfo struct {
	// parse-stack for terms
	terms *stack.SaveStack[termElem]
	// parse-stack for actions on the terms
	//action *stack.Stack[Action]
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
	decls *table.MultiTable[fmt.Stringer, *Declaration]
	//declarations map[string]*Declaration
	generalInfo
	tokenInfo
	// imported modules/packages
	imports ImportTable
	// annotation map
	annotations map[string]Term
	visibility  *stack.Stack[Visibility]
	env         *types.Environment
	mutualBlock *token.Token
	firstPassInfo
	stackInfo
	termMemory *termElem
	mod        *ir.Module
	inTop      bool
}

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
	saver.traits = stack.NewSaveStack[TraitElem](8)
	saver.elems = stack.NewSaveStack[SymbolicElem](32)
	saver.allow = stack.NewStack[SyntaxClass](8)
	saver.inst = stack.NewSaveStack[InstanceElem](8)
	return
}

func Initialize(src source.SourceCode, saveComments bool) (parser *Parser) {
	parser = new(Parser)
	parser.inTop = true
	parser.messages = make([]errors.ErrorMessage, 0)
	parser.src = src
	//parser.action = stack.NewStack[Action](8)
	parser.terms = stack.NewSaveStack[termElem](8)
	parser.saver = initSaver()
	parser.visibility = stack.NewStack[Visibility](2)
	parser.saveComments = saveComments
	if saveComments {
		parser.comments = make([]token.Token, 0, 32)
	}
	parser.env = types.NewEnvironment()
	parser.indentation = indentStack{stack.NewStack[int](8)}
	parser.imports = make(map[string]Import)
	parser.mod = ir.NewModule()
	parser.decls = table.NewMultiTable[fmt.Stringer, *Declaration](8)
	return
}

func Init(src source.SourceCode) (parser *Parser) {
	return Initialize(src, false)
}

func (parser *Parser) Panicking() bool {
	return parser.panicking
}

func (parser *tokenInfo) Advance() token.Token {
	if parser.tokenPos >= len(parser.tokens) {
		return endOfTokensToken()
	}

	parser.tokenPos++
	return parser.tokens[parser.tokenPos-1]
}

// returns next token but does not advance past it
func (parser *tokenInfo) Peek() token.Token {
	if parser.tokenPos >= len(parser.tokens) {
		return endOfTokensToken()
	}

	return parser.tokens[parser.tokenPos]
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

// assumes types have already been validated
func pullData(id, leftRightNone, power token.Token) (name string, right bool, bp uint8, errorMessage string) {
	//@affix <name> <right> <bp>
	name = id.Value

	switch leftRightNone.Value {
	case "Right":
		right = true
	case "None":
		fallthrough
	case "Left":
		right = false
	default:
		errorMessage = ExpectedLRN
		return
	}

	res, e := strconv.ParseUint(power.Value, 0, 4)
	if e != nil {
		if e == strconv.ErrRange {
			errorMessage = ExpectedInteger0to10
		} else {
			errorMessage = ExpectedInteger
		}
		return
	}
	if res > 10 {
		errorMessage = ExpectedInteger0to10
		return
	}

	bp = uint8(res)
	return
}

func (parser *tokenInfo) pull(i int, allowBreaks bool) (id, leftRightNone, power token.Token, errorMessage string) {
	id, i = parser.getTok_breakable(i, token.Id, allowBreaks)
	leftRightNone, i = parser.getTok_breakable(i, token.Id, allowBreaks)
	power, i = parser.getTok_breakable(i, token.IntValue, allowBreaks)
	if i < 0 {
		errorMessage = MalformedAffixAnnotation
	}
	return
}

func (parser *Parser) Load(tokens []token.Token) {
	parser.tokens = tokens
}
