package symbol

import (
	"math"
	"strings"
	"sync"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/stack"
)

type Table struct {
	// all public methods lock this mutex, then call the private method
	mu sync.Mutex
	// freeCounter is the number of free variables declared in this table
	freeCounter int
	// symbols is a scoped map from the name of a symbol to the symbol itself
	symbols    *stack.MapStack[string, sym]
}

// create a new symbol table
func New() *Table {
	return &Table{
		symbols: stack.NewMap[string, sym](),
	}
}

// only enterable from Declare
func (t *Table) declare(tok api.Token) {
	ty := t.newVar(t.freeCounter)
	x := declare(tok, ty) // create a new symbol
	/*success*/_ = t.symbols.MapNew(tok.String(), x)
	t.freeCounter++
	panic("TODO: implement") // TODO: implement
}

// create a new free variable with a unique name
//
// variable will be of the form
//
//	`~a, ~a1, ~a2, ..., ~a9, ~b, ~b1, ..., ~z9, ~a', ~a1', ...`
//
// adding `'` for every 260 variables
func (t *Table) newVar(num int) variable {
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
	return variable(str)
}

func (t *Table) Declare(tok api.Token) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.declare(tok)
}
