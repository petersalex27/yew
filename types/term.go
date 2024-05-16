// =================================================================================================
// Alex Peters - February 29, 2024
// =================================================================================================
package types

import (
	"fmt"
	"math/big"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type Term interface {
	// M[x:=e]
	Substitute(*Term, Variable, Term)
	String() string
	Locate(Variable) bool
}

type typing[T Term] struct {
	Term T
	Kind Type
}

type Typing = typing[Term]

type VarTyping = typing[Variable]

func Split(t Term) (c string, terms []Term) {
	const (
		appc string = "_$"
		pic  string = "_->"
		lamc string = "_λ"
	)
	
	if c, ok := t.(Constant); ok {
		return string(c), []Term{}
	} else if tc, ok := t.(TypeConstant); ok {
		return string(tc), []Term{}
	} else if a, ok := t.(Application); ok {
		return appc, a.terms
	} else if p, ok := t.(Pi); ok {
		return pic, []Term{p.binderType, p.dependent}
	} else if _, ok := t.(Literal); ok {
		return t.String(), []Term{}
	} else if lam, ok := t.(*Lambda); ok {
		return lamc, []Term{lam.binder.Term, lam.bound} // TODO: need binder type??
	} else if _, ok := t.(Variable); ok {
		panic("bug: variables should not be passed to Split")
	} else {
		panic("bug: tried to split unknown term")
	}
	// else if _, ok := t.(*TermBind); ok {
	// 	panic("let binding occurrence inside type")
	// }

}

func reduce(a Term) (ra Term) {
	// TODO: we assume here that every function is total! this is not true!!!!
	// 	the compiler will break on non-terminating functions :(
	if wait, ok := a.(Waiting); ok {
		ra = wait.reduce()
	} else {
		ra = a
	}
	return
}

type Lambda struct {
	binder VarTyping
	bound  Term
}

// type TermBind struct {
// 	bound Term
// 	boundType Type
// 	binding *EmbeddedEnvironment
// }

// func (term *TermBind) Locate(v Variable) bool {
// 	return term.bound.Locate(v) || term.boundType.Locate(v) // TODO: search through embedded environment?
// }

// func (term *TermBind) Substitute(dest *Term, v Variable, t Term) {
// 	if _, found := term.binding.Find(v); found {
// 		// no sub., variable v is shadowed by binding environment
// 		*dest = term
// 		return
// 	}

// 	term.bound.Substitute(&term.bound, v, t)
// }

// func (term *TermBind) String() string {
// 	return fmt.Sprintf("%v:%v in {..}#%d", term.bound, term.boundType, term.binding.uid)
// }

func (lambda *Lambda) Substitute(dest *Term, x Variable, e Term) {
	lambda.bound.Substitute(&lambda.bound, x, e)
}

func (lambda *Lambda) Locate(v Variable) bool {
	return lambda.binder.Kind.Locate(v) || lambda.bound.Locate(v)
}

func (lambda *Lambda) String() string {
	var end string
	if lam, ok := lambda.bound.(*Lambda); ok {
		end = lam.continuedString()
	} else {
		end = " => " + lambda.bound.String()
	}
	return fmt.Sprintf("\\(%v:%v)%v", lambda.binder.Term, lambda.binder.Kind, end)
}

func (lambda *Lambda) continuedString() string {
	var end string
	if lam, ok := lambda.bound.(*Lambda); ok {
		end = lam.continuedString()
	} else {
		end = " => " + lambda.bound.String()
	}
	return fmt.Sprintf(", (%v:%v)%v", lambda.binder.Term, lambda.binder.Kind, end)
}

func (lambda *Lambda) betaReduce(e Term) Term {
	lambda.bound.Substitute(&lambda.bound, lambda.binder.Term, e)
	return lambda.bound
}

type (
	IntConst constant.Int

	FloatConst constant.Float

	CharConst constant.Int

	StringConst constant.CharArray

	//Array constant.Array
)

func (c *IntConst) Parse(s string) (ok bool) {
	c.X, ok = new(big.Int).SetString(s, 10)
	if c.X.BitLen() > 64 {
		return false
	}
	c.Typ = types.I64
	return
}

func (*IntConst) Locate(v Variable) bool { return false }

func (c *IntConst) Substitute(dest *Term, _ Variable, _ Term) {
	*dest = c
}

func (c *IntConst) String() string {
	return c.X.String()
}

func (*FloatConst) Locate(v Variable) bool { return false }

func (c *FloatConst) Parse(s string) (ok bool) {
	c.X, ok = new(big.Float).SetString(s)
	c.Typ = types.Double
	return
}

func (c *FloatConst) Substitute(dest *Term, _ Variable, _ Term) {
	*dest = c
}

func (c *FloatConst) String() string {
	return c.X.String()
}

func (*CharConst) Locate(v Variable) bool { return false }

func (c *CharConst) Parse(s string) (ok bool) {
	if ok = len(s) == 1; !ok {
		return
	}
	c.X = big.NewInt(int64(s[0]))
	c.Typ = types.I8
	return
}

func (c *CharConst) Substitute(dest *Term, _ Variable, _ Term) {
	*dest = c
}

func (c *CharConst) String() string {
	return c.X.String()
}

func (*StringConst) Locate(v Variable) bool { return false }

func (c *StringConst) Parse(s string) (ok bool) {
	c.X = []byte(s)
	c.Typ = types.NewArray(uint64(len(c.X)), types.I8)
	return true
}

func (c *StringConst) Substitute(dest *Term, _ Variable, _ Term) {
	*dest = c
}

func (c *StringConst) String() string {
	return string(c.X)
}
