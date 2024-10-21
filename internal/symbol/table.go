package symbol

import (
	"math"
	"strings"
	//"sync"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/stack"
	"github.com/petersalex27/yew/internal/common"
	"github.com/petersalex27/yew/internal/kind"
)

type (
	Table struct {
		// all public methods lock this mutex, then call the private method
		//mu sync.Mutex
		// freeCounter is the number of free variables declared in this table
		freeCounter int
		// symbols is a scoped map from the name of a symbol to the symbol itself
		symbols *stack.MapStack[string, kind.Symbol]
		// cycles will be detected after declaration
		aliases *stack.MapStack[string, kind.Type]
	}

	errorRecord struct {
		msg string
		api.Positioned
	}

	solver struct {
		// map variables to a representative type
		unifications *stack.MapStack[string, kind.Type]
		freeRange    *stack.Stack[int]
		table        *Table
		errors       []errorRecord
		// occurs stack
		terms *stack.Stack[kind.Type]
	}

	Symbol = kind.Symbol
)

func UnificationFailure(a, b kind.Type) string {
	return "failed to unify " + a.String() + " with " + b.String()
}

// create a new symbol table
func New() *Table {
	return &Table{
		symbols: stack.NewMap[string, kind.Symbol](),
		aliases: stack.NewMap[string, kind.Type](),
	}
}

// this performs no cycle detection, just adds the alias to the table
func (t *Table) defineAliasInitial(tok api.Token, ty kind.Type) string {
	name := tok.String()
	// validate the name
	if !common.IsConstructorId(name) {
		return InvalidAliasIdent
	}

	// add the alias to the table (if it doesn't already exist)
	if !t.aliases.MapNew(name, ty) {
		return AliasRedefined
	}
	return ""
}

// only enterable from Declare
func (t *Table) declare(tok api.Token, ty kind.Type) string {
	//ty := t.newVar(t.freeCounter) // create a type for the declaration

	x := kind.Sym2(tok.String(), ty) // create a new symbol

	success := t.symbols.MapNew(tok.String(), x)
	if !success {
		return DuplicateDeclaration
	}

	// t.freeCounter++

	// if !success {
	// 	t.freeCounter--
	// }

	return ""
}

// create a new free variable with a unique name
//
// variable will be of the form
//
//	`~a, ~a1, ~a2, ..., ~a9, ~b, ~b1, ..., ~z9, ~a', ~a1', ...`
//
// adding `'` for every 260 variables
func (t *Table) newVar() variable {
	num := t.freeCounter
	defer func() { t.freeCounter++ }()

	if num > math.MaxInt32 {
		panic("too many free variables in scope")
	}

	// try to avoid a num >= 10
	offset := (num / 10) % 26 // advance the letter every 10 variables
	num = num % 10
	var out = append(make([]byte, 0, 3), '~', 'a'+byte(offset))
	if num > 0 { // don't use '~a0', just use '~a'
		out = append(out, byte('0'+num)) // num guaranteed to be < 10
	}
	str := string(out)
	if num >= 260 { // add ' for every 260 variables
		str = str + strings.Repeat("'", num/260)
	}
	return variable{str, api.ZeroPosition()}
}

func (t *Table) NewVar() kind.Type {
	// t.mu.Lock()
	// defer t.mu.Unlock()

	v := t.newVar()
	return v
}

func (t *Table) Declare(tok api.Token, ty kind.Type) string {
	// t.mu.Lock()
	// defer t.mu.Unlock()

	return t.declare(tok, ty)
}

func (t *Table) use(tok api.Token) string {
	// get the symbol from an unknown context
	v, ok := t.symbols.Get(tok.String())
	if !ok {
		return "" // must be a free variable, do nothing
	}

	// update multiplicity
	v, ok = v.AccountUse()
	if !ok {
		return OutsideErasedContext
	}

	// update the symbol in the top-most context
	t.symbols.Map(v.Name(), v)
	return ""
}

func (t *Table) Use(tok api.Token) {
	// t.mu.Lock()
	// defer t.mu.Unlock()
	t.use(tok)
}

type Symbolic interface {
	Name() string
	Typing() kind.Type
}

func (t *Table) Track(s Symbolic) {
	// t.mu.Lock()
	// defer t.mu.Unlock()
	// TODO: implement this method
}
