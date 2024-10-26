package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util/code"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/common/stack"
	"github.com/petersalex27/yew/internal/kind"
	"github.com/petersalex27/yew/internal/symbol"
)

// node w/ a symbol
type symbolic[a api.Node] struct {
	node   a
	symbol symbol.Symbol
}

func (s symbolic[a]) GetPos() api.Position {
	return s.node.GetPos()
}

func (s symbolic[a]) Pos() (int, int) {
	return s.node.Pos()
}

func (s symbolic[a]) Type() api.NodeType {
	return s.node.Type()
}

func (s symbolic[a]) Describe() (string, []api.Node) {
	if dn, ok := api.Node(s.node).(api.DescribableNode); ok {
		return dn.Describe()
	}
	return "symbolic node, " + s.symbol.String(), []api.Node{s.node}
}

type symbols struct {
	typings *stack.MapStack[string, symbolic[declaration]]
}

type state struct {
	scanner api.ScannerPlus
	solver  kind.Solver
	symbols symbols

	tokens       []api.Token
	tokenCounter int

	errors   []error
	warnings []error
}

// State for parsing
type ParserState struct {
	state
	ast YewSource
}

func createState(scanner api.ScannerPlus) state {
	return state{
		scanner:  scanner,
		tokens:   nil,
		errors:   make([]error, 0),
		warnings: make([]error, 0),
		solver:   symbol.Solver(symbol.New()),
		symbols: symbols{
			stack.NewMap[string, symbolic[declaration]](),
		},
	}
}

func (p *ParserState) dropNewlines() {
	for token.Newline.Match(p.current()) {
		p.advance()
	}
}

func (p *ParserState) SrcCode() api.SourceCode {
	return p.scanner.SrcCode()
}

func (p *ParserState) Pos() (int, int) {
	return p.current().Pos()
}

func (p *ParserState) GetPos() api.Position {
	return p.current().GetPos()
}

func (s *symbols) declare(dec declaration, ty kind.Type, fix byte, prec byte) *data.Ers {
	name := data.Cases(dec, fun.Id, func(typing typing) api.Token { return typing.typing.Fst().Break() })
	var sym kind.Symbol
	if token.Infix.Match(name) {
		if fix == kind.None {
			fix = kind.Infix
		}
		sym = kind.InfixSym(name, ty, fix, prec)
	} else {
		sym = kind.Sym(name, ty)
	}

	success := s.typings.MapNew(name.String(), symbolic[declaration]{dec, sym})
	if !success {
		// e := data.MkErr(RedeclaredSymbol, name)
		// es := data.Makes(e)
		// return &es // already declared
	}
	return nil
}

func (p *ParserState) declare(dec declaration, fix_prec ...byte) *data.Ers {
	fix, prec := kind.None, byte(0)
	if len(fix_prec) >= 2 {
		fix, prec = fix_prec[0], fix_prec[1]
	}
	return p.symbols.declare(dec, p.newVar(), fix, prec)
}

func (p *ParserState) newScope() {
	p.solver.RecordNames()
	p.symbols.typings.Push(map[string]symbolic[declaration]{})
}

func (p *ParserState) exitScope() map[string]symbolic[declaration] {
	p.solver.EndRecording()
	typings, ok := p.symbols.typings.Pop()
	if !ok {
		panic("unbalanced scopes")
	}
	return typings
}

func (p *ParserState) newVar() kind.Type {
	return p.solver.NewVar()
}

// adds the error, and if fatal, then returns a *ParserStateFail
func (p *ParserState) report(e error, fatal bool) parser {
	p.AddError(e)
	if fatal {
		return &ParserStateFail{bad: p}
	}
	return p
}

func (parser *ParserState) AppendTokens(tokens ...api.Token) {
	if parser.tokens == nil {
		cap := int64(float64(len(tokens)) * 1.25)
		parser.tokens = make([]api.Token, 0, int64(cap))
	}
	parser.tokens = append(parser.tokens, tokens...)
}

func (parser *ParserState) Ast() api.Node {
	return parser.ast
}

func (parser *ParserState) load() bool {
	tokens, errorToken := code.Tokenize(parser.scanner, nil)
	if errorToken != nil {
		parser.AddError((*errorToken).Error())
		return false
	}
	parser.tokens = tokens
	return true
}

func (parser *ParserState) replStatement() bool {
	panic("not implemented")
}

func (parser *ParserState) ReplParse() bool {
	if !parser.load() {
		return false
	}

	return parser.replStatement()
}

func (parser *ParserState) Parse() bool {
	if !parser.load() {
		return false
	}

	parser.Run()
	return len(parser.Errors()) == 0
}

func (parser *ParserState) AddError(err error) {
	parser.errors = append(parser.errors, err)
}

// ensure that the errors slice is never nil when needed
func (parser *ParserState) Errors() []error {
	if parser.errors == nil {
		parser.errors = make([]error, 0)
	}
	return parser.errors
}

func (parser *ParserState) ReferenceScanner() *api.Scanner {
	s := api.Scanner(parser.scanner)
	return &s
}

func (*ParserState) Clear() api.Parser {
	return &ParserState{}
}

func (parser *ParserState) current() api.Token {
	if parser.tokenCounter >= len(parser.tokens) {
		return token.EndOfTokens.Make()
	}
	return parser.tokens[parser.tokenCounter]
}

func (parser *ParserState) advance() {
	if parser.tokenCounter < len(parser.tokens) {
		parser.tokenCounter++
	}
}

func (p *ParserState) Run() {
	// TODO: Implement
	panic("not implemented")
}
