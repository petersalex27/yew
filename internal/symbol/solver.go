package symbol

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/stack"
	"github.com/petersalex27/yew/internal/kind"
)

func (s *solver) GiveType(tok api.Token, ty kind.Type) string {
	s.table.mu.Lock()
	defer s.table.mu.Unlock()

	if v, ok := s.table.symbols.Get(tok.String()); ok {
		if !Unify(s, v.Typing(), ty) {
			return UnificationFailure(v.Typing(), ty)
		}

		v = v.UseRepType(s)                  // use representative type
		s.table.symbols.Map(tok.String(), v) // remap the symbol
	} else {
		return UndeclaredSymbol
	}

	return ""
}

// create a new solver that uses the given table
func Solver(table *Table) *solver {
	return &solver{
		stack.NewMap[string, kind.Type](),
		stack.New[int](),
		table,
		make([]errorRecord, 0),
		stack.New[kind.Type](),
	}
}

func (s *solver) newScope() {
	s.freeRange.Push(s.table.freeCounter)
	s.unifications.Push(make(map[string]kind.Type))
}

func (s *solver) exitScope() {
	_, _ = s.unifications.Pop()
	if saved, ok := s.freeRange.Pop(); ok {
		s.table.freeCounter = saved
	} // else, I guess just leave the counter alone ...
}

func (s *solver) RecordNames() {
	s.table.symbols.Push(map[string]kind.Symbol{})
}

func (s *solver) EndRecording() map[string]kind.Symbol {
	res, ok := s.table.symbols.Pop()
	if !ok {
		panic("unbalanced scopes")
	}
	return res
}

func (s *solver) Declare(tok api.Token) {
	s.table.mu.Lock()
	defer s.table.mu.Unlock()

	// ty := s.table.NewVar()
	if _, found := s.table.symbols.Get(tok.String()); found {
		return // already declared
	}

	s.table.symbols.Map(tok.String(), kind.Sym(tok, s.table.newVar()))
}

func (s *solver) NewVar() kind.Type {
	return s.table.NewVar()
}

func (s *solver) RecordError(err string, pos api.Positioned) {
	s.errors = append(s.errors, errorRecord{err, pos})
}
