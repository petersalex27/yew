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
	TypeVariable: types.Var("n"),
	Functions: map[string]types.Function{
		"+": nestFunction(types.Var("n"), 3),
		"-": nestFunction(types.Var("n"), 3),
		"/": nestFunction(types.Var("n"), 3),
		"*": nestFunction(types.Var("n"), 3),
		"^": nestFunction(types.Var("n"), 3),
		"negative": nestFunction(types.Var("n"), 2),
		"positive": nestFunction(types.Var("n"), 2),
	},
}

/*
class Equalable e where
	(==) :: e -> e -> Bool;
	(!=) :: e -> e -> Bool
*/
/*
var Equalable types.Class = types.Class{
	Name: "Equalable",
	TypeVariable: types.Var("e"),
	Functions: map[string]types.Function{
		"==": {types.Var("e"), types.Function{types.Var("e"), types.Bool{}}},
		"!=": {types.Var("e"), types.Function{types.Var("e"), types.Bool{}}},
	},
}*/

/*
class Orderable o where 
	(>) :: o -> o -> Bool;
	(<) :: o -> o -> Bool;
	(>=) :: o -> o -> Bool;
	(<=) :: o -> o -> Bool
*/
/*
var Orderable types.Class = types.Class{
	Name: "Orderable",
	TypeVariable: "o",
	Functions: map[string]types.Function{
		">": {types.Var("o"), types.Function{types.Var("o"), types.Bool{}}},
		"<": {types.Var("o"), types.Function{types.Var("o"), types.Bool{}}},
		">=": {types.Var("o"), types.Function{types.Var("o"), types.Bool{}}},
		"<=": {types.Var("o"), types.Function{types.Var("o"), types.Bool{}}},
	},
}*/

/*
class Listable x where 
	(:) :: x -> [x] -> [x];
	(++) :: [x] -> [x] -> [x];
	head :: [x] -> x;
	tail :: [x] -> [x]
*/
/*
var Listable types.Class = types.Class{
	Name: "Listable",
	TypeVariable: "x",
	Functions: map[string]types.Function{
		":": {types.Var("x"), nestFunction(types.Array{types.Var("x")}, 2)},
		"++": nestFunction(types.Array{types.Var("x")}, 3),
		"head": {types.Array{types.Var("x")}, types.Var("x")},
		"tail": nestFunction(types.Array{types.Var("x")}, 2),
	},
}*/