package inf

import (
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
	"github.com/petersalex27/yew-packages/util/stack"
)

// shadow-able symbol
type Symbol[T nameable.Nameable] struct {
	data *stack.Stack[bridge.JudgmentAsExpression[T, expr.Const[T]]]
}

func MakeSymbol[T nameable.Nameable]() Symbol[T] {
	// this assumes that most symbols will not be shadowed
	const initialCapacity uint = 1
	stk := stack.NewStack[bridge.JudgmentAsExpression[T, expr.Const[T]]](initialCapacity)
	return Symbol[T]{data: stk}
}

// get symbol (and judgment)
func (sym *Symbol[T]) Get() bridge.JudgmentAsExpression[T, expr.Const[T]] {
	judgment, _ := sym.data.Peek()
	return judgment
}

// create/shadow symbol
func (sym *Symbol[T]) Shadow(name expr.Const[T], ty types.Type[T]) {
	judgment := bridge.Judgment(name, ty)
	sym.data.Push(judgment)
}

func (sym *Symbol[T]) IncludeInExport(name expr.Const[T], ty types.Type[T]) (added bool) {
	added = sym.data.GetCount() != 0
	if !added {
		return
	}

	sym.Shadow(name, ty)
	return
}

// remove/unshadow symbol
func (sym *Symbol[T]) Unshadow() (remove bool) {
	_, _ = sym.data.Pop()
	remove = sym.data.GetCount() == 0
	return remove
}

// export symbol
func (sym *Symbol[T]) Export() (export *Symbol[T], exported bool) {
	// check if symbol should be exported
	exported = sym.data.GetCount() == 1
	if exported {
		return sym, true
	}

	return nil, false
}
