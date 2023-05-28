package builtin

import (
	types "yew/type"
)

func nestFunction(t types.Types, n int) types.Function {
	if 2 <= n {
		return types.Function{Domain: t, Codomain: t}
	}
	return types.Function{Domain: t, Codomain: nestFunction(t, n - 1)}
} 

/*
class Number n where
	(+) :: n -> n -> n;
	(-) :: n -> n -> n;
	(/) :: n -> n -> n;
	(*) :: n -> n -> n;
	(^) :: n -> n -> n;
	negative :: n -> n;
	positive :: n -> n
*/
var Number types.Class = types.Class{
	Name: "Number",
	TypeVariable: "n",
	Functions: map[string]types.Function{
		"+": nestFunction(types.Tau("n"), 3),
		"-": nestFunction(types.Tau("n"), 3),
		"/": nestFunction(types.Tau("n"), 3),
		"*": nestFunction(types.Tau("n"), 3),
		"^": nestFunction(types.Tau("n"), 3),
		"negative": nestFunction(types.Tau("n"), 2),
		"positive": nestFunction(types.Tau("n"), 2),
	},
}

/*
class Equalable e where
	(==) :: e -> e -> Bool;
	(!=) :: e -> e -> Bool
*/
var Equalable types.Class = types.Class{
	Name: "Equalable",
	TypeVariable: "e",
	Functions: map[string]types.Function{
		"==": {types.Tau("e"), types.Function{types.Tau("e"), types.Bool{}}},
		"!=": {types.Tau("e"), types.Function{types.Tau("e"), types.Bool{}}},
	},
}

/*
class Orderable o where 
	(>) :: o -> o -> Bool;
	(<) :: o -> o -> Bool;
	(>=) :: o -> o -> Bool;
	(<=) :: o -> o -> Bool
*/
var Orderable types.Class = types.Class{
	Name: "Orderable",
	TypeVariable: "o",
	Functions: map[string]types.Function{
		">": {types.Tau("o"), types.Function{types.Tau("o"), types.Bool{}}},
		"<": {types.Tau("o"), types.Function{types.Tau("o"), types.Bool{}}},
		">=": {types.Tau("o"), types.Function{types.Tau("o"), types.Bool{}}},
		"<=": {types.Tau("o"), types.Function{types.Tau("o"), types.Bool{}}},
	},
}

/*
class Listable x where 
	(:) :: x -> [x] -> [x];
	(++) :: [x] -> [x] -> [x];
	head :: [x] -> x;
	tail :: [x] -> [x]
*/
var Listable types.Class = types.Class{
	Name: "Listable",
	TypeVariable: "x",
	Functions: map[string]types.Function{
		":": {types.Tau("x"), nestFunction(types.Array{types.Tau("x")}, 2)},
		"++": nestFunction(types.Array{types.Tau("x")}, 3),
		"head": {types.Array{types.Tau("x")}, types.Tau("x")},
		"tail": nestFunction(types.Array{types.Tau("x")}, 2),
	},
}