// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package types

import (
	"fmt"
	"math/big"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

type Literal interface {
	Parse(strPos) (Literal, bool)
	Term
}

type (
	IntConst struct {
		*constant.Int
		Start, End int
	}

	FloatConst struct {
		*constant.Float
		Start, End int
	}

	CharConst struct {
		*constant.Int
		Start, End int
	}

	StringConst struct {
		*constant.CharArray
		Start, End int
	}

	//Array constant.Array
)

func (IntConst) CollectVariables(m map[string]Variable) map[string]Variable { return m }

func (*IntConst) Parse(s strPos) (_ Literal, ok bool) {
	c := &IntConst{}
	c.Int = &constant.Int{}
	c.X, ok = new(big.Int).SetString(s.String(), 10)
	if c.X.BitLen() > 64 {
		return nil, false
	}
	c.Typ = types.I64
	c.Start, c.End = s.Pos()
	return c, ok
}

var (
	int_type Variable = Variable{false, "Int", 0, Unrestricted, Type0, 0, 0}
	float_type Variable = Variable{false, "Float", 0, Unrestricted, Type0, 0, 0}
	char_type Variable = Variable{false, "Char", 0, Unrestricted, Type0, 0, 0}
	string_type Variable = Variable{false, "String", 0, Unrestricted, Type0, 0, 0}
)

func (c *IntConst) GetKind() (Term, Type) {
	return c, int_type
}

func (*IntConst) Locate(v Variable) bool { return false }

func (c *IntConst) Substitute(dest *Term, _ Variable, _ Term) {
	*dest = c
}

func (c *IntConst) Pos() (start, end int) {
	return c.Start, c.End
}

func (c *IntConst) String() string {
	return c.X.String()
}

func (*FloatConst) CollectVariables(m map[string]Variable) map[string]Variable { return m }

func (*FloatConst) Locate(v Variable) bool { return false }

func (c *FloatConst) GetKind() (Term, Type) {
	return c, float_type
}

func (c *FloatConst) Parse(s strPos) (_ Literal, ok bool) {
	c = &FloatConst{}
	c.Float = &constant.Float{}
	c.X, ok = new(big.Float).SetString(s.String())
	c.Typ = types.Double
	c.Start, c.End = s.Pos()
	return c, ok
}

func (c *FloatConst) Substitute(dest *Term, _ Variable, _ Term) {
	*dest = c
}

func (c *FloatConst) Pos() (start, end int) {
	return c.Start, c.End
}

func (c *FloatConst) String() string {
	return c.X.String()
}

func (*CharConst) CollectVariables(m map[string]Variable) map[string]Variable { return m }

func (*CharConst) Locate(v Variable) bool { return false }

func (c *CharConst) GetKind() (Term, Type) {
	return c, char_type
}

func (c *CharConst) Parse(s strPos) (_ Literal, ok bool) {
	c.Int = &constant.Int{}
	st := s.String()
	if ok = len(st) == 1; !ok {
		return nil, false
	}
	c.X = big.NewInt(int64(st[0]))
	c.Typ = types.I8
	c.Start, c.End = s.Pos()
	return c, ok
}

func (c *CharConst) Substitute(dest *Term, _ Variable, _ Term) {
	*dest = c
}

func (c *CharConst) Pos() (start, end int) {
	return c.Start, c.End
}

func (c *CharConst) String() string {
	return c.X.String()
}

func (*StringConst) CollectVariables(m map[string]Variable) map[string]Variable { return m }

func (*StringConst) Locate(v Variable) bool { return false }

func (c *StringConst) GetKind() (Term, Type) {
	return c, string_type
}

func (c *StringConst) Parse(s strPos) (_ Literal, ok bool) {
	c = &StringConst{}
	c.CharArray = &constant.CharArray{}
	c.X = []byte(s.String())
	c.Typ = types.NewArray(uint64(len(c.X)), types.I8)
	c.Start, c.End = s.Pos()
	return c, true
}

func (c *StringConst) Substitute(dest *Term, _ Variable, _ Term) {
	*dest = c
}

func (c *StringConst) Pos() (start, end int) {
	return c.Start, c.End
}

func (c *StringConst) String() string {
	return fmt.Sprint(c.X)
}