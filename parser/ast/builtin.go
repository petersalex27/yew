package ast

import (
	"yew/symbol"
	//"yew/type"
)

func IncludeBuiltin(table *symbol.SymbolTable, key string) bool {
	panic("TODO") // TODO
}
/*
// Int -> Int
var intToInt = types.Function{Domain: types.Int{}, Codomain: types.Int{}}
// Int -> (Int -> Int)
var intToIntToInt = types.Function{Domain: types.Int{}, Codomain: intToInt}
// Int -> (Int -> Bool)
var intToIntToBool = types.Function{
	Domain: types.Int{}, 
	Codomain: types.Function{
		Domain: types.Int{},
		Codomain: types.Bool{},
	},
}
// Float -> Float
var floatToFloat = types.Function{Domain: types.Float{}, Codomain: types.Float{}}
// Float -> (Float -> Float)
var floatToFloatToFloat = types.Function{Domain: types.Float{}, Codomain: floatToFloat}
// Float -> (Float -> Bool)
var floatToFloatToBool = types.Function{
	Domain: types.Float{},
	Codomain: types.Function{
		Domain: types.Float{},
		Codomain: types.Bool{},
	},
}
// Bool -> Bool -> Bool
var boolToBoolToBool = types.Function{
	Domain: types.Bool{},
	Codomain: types.Function{
		Domain: types.Bool{},
		Codomain: types.Bool{},
	},
}
// Char -> Char -> Bool
var charToCharToBool = types.Function{
	Domain: types.Char{},
	Codomain: types.Function{
		Domain: types.Char{},
		Codomain: types.Bool{},
	},
}

// creates a declaration for builtins
func createDeclaration(name string, t types.Types) Declaration {
	var builtinLocation = symbol.MakeLocation("builtin", 0, 0)
	var emptySymbolMap = map[string]symbol.SymbolUse{}
	return Declaration(Id{id: symbol.MakeSymbol_testable(name, t, builtinLocation, emptySymbolMap)})
}

var builtins = []Declaration{
	// Int builtins
	createDeclaration("addInt", intToIntToInt),
	createDeclaration("subInt", intToIntToInt),
	createDeclaration("mulInt", intToIntToInt),
	createDeclaration("divInt", intToIntToInt),
	createDeclaration("negInt", intToInt),
	createDeclaration("equalsInt", intToIntToBool),
	createDeclaration("notEqualsInt", intToIntToBool),
	createDeclaration("greaterInt", intToIntToBool),
	createDeclaration("lesserInt", intToIntToBool),
	createDeclaration("greaterOrEqualInt", intToIntToBool),
	createDeclaration("lesserOrEqualInt", intToIntToBool),
	createDeclaration("powInt", intToIntToInt),
	createDeclaration("remainderInt", intToIntToInt),
	// Float builtins
	createDeclaration("addFloat", floatToFloatToFloat),
	createDeclaration("subFloat", floatToFloatToFloat),
	createDeclaration("mulFloat", floatToFloatToFloat),
	createDeclaration("divFloat", floatToFloatToFloat),
	createDeclaration("negFloat", floatToFloat),
	createDeclaration("equalsFloat", floatToFloatToBool),
	createDeclaration("notEqualsFloat", floatToFloatToBool),
	createDeclaration("greaterFloat", floatToFloatToBool),
	createDeclaration("lesserFloat", floatToFloatToBool),
	createDeclaration("greaterOrEqualFloat", floatToFloatToBool),
	createDeclaration("lesserOrEqualFloat", floatToFloatToBool),
	// createDeclaration("powFloat", floatToFloatToFloat) // TODO
	// Bool builtins
	createDeclaration("equalsBool", boolToBoolToBool),
	createDeclaration("notEqualsBool", boolToBoolToBool),
	// Char builtins
	createDeclaration("equalsChar", charToCharToBool),
	createDeclaration("notEqualsChar", charToCharToBool),
}

func InitializeBuiltins(table *symbol.SymbolTable) bool {
	for _, fn := range builtins {
		addError, added := table.AddSymbol(fn.id)
		if !added {
			addError.ToError().Print()
			return false
		}
	}
	return true
}*/