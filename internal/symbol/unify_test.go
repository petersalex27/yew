package symbol

import (
	"testing"

	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/common/stack"
	"github.com/petersalex27/yew/internal/kind"
)

var (
	vx = Var(token.Id.MakeValued("x"))
	vy = Var(token.Id.MakeValued("y"))
	vz = Var(token.Id.MakeValued("z"))
	vw = Var(token.Id.MakeValued("w"))
	cA = Const(token.Id.MakeValued("A"))
	cB = Const(token.Id.MakeValued("B"))
	cC = Const(token.Id.MakeValued("C"))
)

func TestOccursHelper(t *testing.T) {
	tests := []struct {
		name     string
		v        string
		A        kind.Type
		expected bool
	}{
		{
			name:     "Variable does not occur in constant",
			v:        "x",
			A:        cA,
			expected: false,
		},
		{
			name:     "Variable does not occur in itself",
			v:        "x",
			A:        vx,
			expected: false,
		},
		{
			name:     "Variable occurs in nested type",
			v:        "x",
			A:        Pi(kind.Sym2("", vx), vy),
			expected: true,
		},
		{
			name:     "Variable does not occur in nested type",
			v:        "x",
			A:        Pi(kind.Sym2("", vy), vy),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := occursHelper(tt.v, tt.A, stack.New[kind.Type]())
			if result != tt.expected {
				t.Errorf("occursHelper(%v, %v, _) = %v; want %v", tt.v, tt.A, result, tt.expected)
			}
		})
	}
}
func TestUnify(t *testing.T) {

	tests := []struct {
		name     string
		a        kind.Type
		b        kind.Type
		expected bool
		// state the solver's top unification map should start with
		initialState map[string]kind.Type
		// state the solver's top unification map should end with
		endState map[string]kind.Type
	}{
		{
			name:         "Unify identical constants",
			a:            cA,
			b:            cA,
			expected:     true,
			initialState: map[string]kind.Type{},
			endState:     map[string]kind.Type{},
		},
		{
			name:         "Unify different constants",
			a:            cA,
			b:            cB,
			expected:     false,
			initialState: map[string]kind.Type{},
			endState:     map[string]kind.Type{},
		},
		{
			name:         "Unify variable with constant",
			a:            vx,
			b:            cA,
			expected:     true,
			initialState: map[string]kind.Type{},
			endState:     map[string]kind.Type{"x": cA},
		},
		{
			name:         "Unify variable with constant, state causes failure",
			a:            vx,
			b:            cA,
			expected:     false,
			initialState: map[string]kind.Type{"x": cB},
			endState:     map[string]kind.Type{"x": cB},
		},
		{
			name:         "Unify variable with constant, state causes new unification",
			a:            vx,
			b:            cA,
			expected:     true,
			initialState: map[string]kind.Type{"x": vy},
			endState:     map[string]kind.Type{"x": vy, "y": cA},
		},
		{
			name:         "Unify constant with variable",
			a:            cA,
			b:            vx,
			expected:     true,
			initialState: map[string]kind.Type{},
			endState:     map[string]kind.Type{"x": cA},
		},
		{
			name:         "Unify constant with variable, state causes failure",
			a:            cA,
			b:            vx,
			expected:     false,
			initialState: map[string]kind.Type{"x": cB},
			endState:     map[string]kind.Type{"x": cB},
		},
		{
			name:         "Unify constant with variable, state causes new unification",
			a:            cA,
			b:            vx,
			expected:     true,
			initialState: map[string]kind.Type{"x": vy},
			endState:     map[string]kind.Type{"x": vy, "y": cA},
		},
		{
			name:         "Unify identical variables",
			a:            vx,
			b:            vx,
			expected:     true,
			initialState: map[string]kind.Type{},
			endState:     map[string]kind.Type{"x": vx},
		},
		{
			name:         "Unify identical variables, state already has unification",
			a:            vx,
			b:            vx,
			expected:     true,
			initialState: map[string]kind.Type{"x": vy},
			endState:     map[string]kind.Type{"x": vy, "y": vy},
		},
		{
			name:         "Unify identical variables, state already has unification, init = end",
			a:            vx,
			b:            vx,
			expected:     true,
			initialState: map[string]kind.Type{"x": vx},
			endState:     map[string]kind.Type{"x": vx},
		},
		{
			name:         "Unify different variables",
			a:            vx,
			b:            vy,
			expected:     true,
			initialState: map[string]kind.Type{},
			endState:     map[string]kind.Type{"x": vy},
		},
		{
			name:         "Unify different variables, state already has unification - 1",
			a:            vx,
			b:            vy,
			expected:     true,
			initialState: map[string]kind.Type{"x": vz},
			endState:     map[string]kind.Type{"x": vz, "z": vy},
		},
		{
			name:         "Unify different variables, state already has unification - 2",
			a:            vx,
			b:            vy,
			expected:     true,
			initialState: map[string]kind.Type{"x": vz, "y": vw},
			endState:     map[string]kind.Type{"x": vz, "y": vw, "z": vw},
		},
		{
			name:         "Unify nested types with same structure",
			a:            Pi(kind.Sym2("", cA), cB),
			b:            Pi(kind.Sym2("", cA), cB),
			expected:     true,
			initialState: map[string]kind.Type{},
			endState:     map[string]kind.Type{},
		},
		{
			name:         "Unify nested types with different structure",
			a:            Pi(kind.Sym2("", cA), cB),
			b:            Pi(kind.Sym2("", cC), cB),
			expected:     false,
			initialState: map[string]kind.Type{},
			endState:     map[string]kind.Type{},
		},
		{
			name: "Unification of nested types causes new unification, causing failure",
			// this example should attempt to unify `A` with `B` through the initial unification of `y = A`
			//
			// an example of how this could happen:
			//
			//		fun : (A -> y) -> (y -> B)
			//		fun f = f
			a:            Pi(kind.Sym2("", cA), vy),
			b:            Pi(kind.Sym2("", vy), cB),
			expected:     false,
			initialState: map[string]kind.Type{},
			endState:     map[string]kind.Type{"y": cA},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &solver{
				unifications: stack.NewMap[string, kind.Type](),
				occursStack:  stack.New[kind.Type](),
			}
			s.unifications.Push(tt.initialState)

			result := Unify(s, tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("unify(%v, %v) = %v; want %v", tt.a, tt.b, result, tt.expected)
			}

			// check the final state of the unification map
			endState, ok := s.unifications.Pop()
			if !ok {
				t.Fatalf("unification map was empty")
			}

			if len(endState) != len(tt.endState) {
				t.Fatalf("unification map length mismatch: got %v, want %v", endState, tt.endState)
			}

			for k, v := range tt.endState {
				if endState[k] != v {
					t.Errorf("unification map[%v] = %v; want %v", k, endState[k], v)
				}
			}
		})
	}
}
