package module

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/petersalex27/yew/internal/kind"
	"github.com/petersalex27/yew/internal/symbol"
)

var pointer = types.NewPointer(types.Void)

func mockModule() *ir.Module {
	mod := ir.NewModule()

	_ = mod.NewTypeDef("Type", pointer)
	Int := mod.NewTypeDef("Int", types.I64)

	addInt := mod.NewFunc("addInt", Int, ir.NewParam("x", Int), ir.NewParam("y", Int))
	addInt.Linkage = enum.LinkageExternal
	{
		b := addInt.NewBlock("")
		b.NewRet(b.NewAdd(addInt.Params[0], addInt.Params[1]))
	}

	return mod
}

func mockSymbols() map[string]symbol.Symbolic {
	Type := kind.Sym2("Type", kind.Sym2( /* never add Type1 to the symbol table!!! */ "Type1", nil))
	Int := kind.Sym2("Int", Type)
	// addInt : Int -> Int -> Int
	addIntKind := symbol.Pi(
		kind.Sym2("", Int),
		symbol.Pi(kind.Sym2("", Int), Int))
	return map[string]symbol.Symbolic{
		"Type":   Type,
		"Int":    Int,
		"addInt": kind.Sym2("addInt", addIntKind),
	}
}

func Mock() *Data {
	m := new(Data)
	m.name = "mock"
	m.internal = mockModule()
	m.symbols = mockSymbols()
	return m
}
