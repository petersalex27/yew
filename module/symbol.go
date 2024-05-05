// =================================================================================================
// Alex Peters - February 01, 2024
//
// =================================================================================================
package module

import (
	"fmt"

	"github.com/petersalex27/yew/token"
	"github.com/petersalex27/yew/types"
)

type Associativity = uint8

type Ident struct {
	token.Token
}

const (
	Left  Associativity = 0 // left or non-associative
	Right Associativity = 1 // right associative
)

type Binding uint8

// think of lsb (which, is on the right side of the first byte) as the "right assoc" flag
const (
	InfixLeft    Binding = 0
	InfixRight   Binding = InfixLeft + Binding(Right)
	PrefixLeft   Binding = 2
	PrefixRight  Binding = PrefixLeft + Binding(Right)
	PostfixLeft  Binding = 4
	PostfixRight Binding = PostfixLeft + Binding(Right)
)

func (binding Binding) LeftAssoc() bool {
	const balance Binding = 2
	return (binding % balance) == Binding(Left)
}

func (binding Binding) Postfixed() bool {
	return binding >= PostfixLeft
}

func (binding Binding) RightAssoc() bool {
	const balance Binding = 2
	return (binding % balance) == Binding(Right)
}

// symbol w/in symbol table
//
// Binding: _, infix
// Arity: 3
//
//	_?_!:_
//
// Binding: _, prefix
// Arity: 3
//
//	if_then_else_
//
// Binding: _, postfix
// Arity: 2
//
//	_a_b
//
// basically, given arity of n:
//   - infix will have n "_" and n-1 identifiers placed b/w each "_"
//   - prefix will have n "_" and n identifiers, one placed in front of each "_"
//   - postfix will have n "_" and n identifiers, one placed after each "_"
type Symbol struct {
	Id Ident
	Binding
	Precedence uint8
	Arity      uint
	Judgment   types.Type
}

// symbol for imported module
type ModuleSymbol struct {
	// module identifier
	Id Ident
	// path to module's `yew.root` file
	Path string
	// table of modules symbols
	Table *SymbolTable
}

// symbol table, holds multiple scopes of symbols
type SymbolTable struct {
	constructors map[string]ConstructorSymbol
	imported     map[string]ModuleSymbol
	tables       []map[string]Symbol
	traits       map[string]TraitSymbol
	types        map[string]TypeSymbol
	pos          int
}

type TraitSymbol struct {
	Id      TypeConstant
	Symbols []Symbol
}

type TypeConstant struct {
	token.Token
}

type DependentTypeSymbol struct {
	isAlias bool
	Name    TypeConstant
	Params  []TypeVariable
}

type FamilySymbol struct {
	Name    TypeConstant
	Params  []TypeVariable
	Members []TypeSymbol
}

type TypeSymbol struct {
	isAlias bool
	types.Application
}

// symbol for a constructor of type T
//
//	let (C t0 .. tN) =
//	  ((\p0 .. pN -> (@tag C) p0 .. pN) : t0 -> .. -> tN -> T)
//	in C
type ConstructorSymbol struct {
	// type constructor constructs for
	Constructs types.Type
	// name of constructor
	Name string
	// constructor params
	Params []types.Type
}

func MakeConstructorSymbol(constructs types.Type, name string, params []types.Type) ConstructorSymbol {
	return ConstructorSymbol{constructs, name, params}
}

type TypeVariable struct {
	token.Token
}

// grabs module symbol from table
func (table *SymbolTable) Access(key fmt.Stringer) (module ModuleSymbol, found bool) {
	module, found = table.imported[key.String()]
	return
}

// adds an additional scope or initializes an empty table
func (table *SymbolTable) AddScope() {
	if cap(table.tables) == 0 {
		table.tables = make([]map[string]Symbol, 0, 8)
		table.pos = -1 // set to zero below
	}
	table.tables = append(table.tables, make(map[string]Symbol))
	table.pos++
}

// creates an alias for a type that already exists in the table
//
// fails when type is not in the table (ok==false, found==false) or when alias is already in table
// (ok==false, found==true)
//
// panics if type symbol is not a type alias
//
// panics if type symbol has no location specified (location uses the params field)
func (table *SymbolTable) Alias(alias TypeSymbol) (ok, found bool) {
	if !alias.isAlias {
		panic("illegal argument: type symbol is not an alias")
	}

	if len(alias.Application) < 1 {
		panic("illegal argument: type alias symbol has no alias specified")
	}

	// Params slice is used for locating the aliased type. It is defined like follows:
	//	{[0]: accessor, [1]: accessor, .., [len(Params)-1]: aliased type key}
	// So, iterate through just the accessors, returning the type table each time (initial table is
	// `table.types`).
	// Then, the type table at the end accesses the type key at [len(Params)-1]
	typeTable := table.types
	aliased := alias.TypeConst
	for _, accessor := range alias.Args {
		var modSym ModuleSymbol
		modSym, found = table.imported[accessor.String()]
		if ok = found; !ok {
			return
		}
		typeTable = modSym.Table.types
	}

	found = false // reassign value, aliased type not yet found

	// table is empty? => fail
	ok = len(typeTable) > 0
	if !ok {
		return
	}

	// aliased type not found? => fail
	if _, found = typeTable[aliased.String()]; !found {
		ok = false
		return
	}

	key := alias.String()
	// alias redefined? => fail
	if _, redef := table.types[key]; redef {
		ok, found = false, true
		return
	}

	// define alias
	table.types[key] = alias
	return
}

// defines a module symbol
//
// returns ok==true iff symbol was added. If symbol was not added, it was because it was already
// defined
func (table *SymbolTable) DefineModule(module ModuleSymbol) (ok bool) {
	key := module.Id.String()
	_, found := table.imported[key]
	if ok = !found; !ok {
		return
	}
	table.imported[key] = module
	return
}

// defines a symbol w/in the current scope
//
// returns ok==true iff symbol was added. If symbol was not added, it was because it was already
// defined in the current scope
func (table *SymbolTable) DefineSymbol(symbol Symbol) (ok bool) {
	key := symbol.Id.String()
	_, found := table.lookup_unsafe(table.pos, key)
	if ok = !found; !ok {
		return
	}

	table.tables[table.pos][key] = symbol
	return
}

func (table *SymbolTable) DefineTrait(symbol TraitSymbol) (ok bool) {
	key := symbol.Id.Token.String()
	_, found := table.traits[key]
	if ok = !found; !ok {
		return
	}

	table.traits[key] = symbol
	return
}

// defines a type symbol `symbol`
//
// returns false, and does not redefine the symbol, if and only if type denoted by `symbol` has
// already been defined
//
// panics if a type symbol is a type alias
func (table *SymbolTable) DefineType(symbol TypeSymbol) (ok bool) {
	if symbol.isAlias {
		panic("illegal argument: type symbol is an alias")
	}
	key := symbol.String()
	_, found := table.types[key]
	if ok = !found; !ok {
		return // type redefinition
	}

	table.types[key] = symbol
	return
}

// defines a type constructor symbol `symbol`
//
// returns false, and does not redefine the symbol, if and only if constructor denoted by `symbol` has
// already been defined
func (table *SymbolTable) DefineTypeConstructor(symbol ConstructorSymbol) (ok bool) {
	key := symbol.Name
	_, found := table.constructors[key]
	if ok = !found; !ok {
		return // constructor redefinition
	}

	table.constructors[key] = symbol
	return
}

// looks up symbol with key `key`. Returns symbol (if found) and the truthy of whether it was found
func (table *SymbolTable) Lookup(key fmt.Stringer) (sym Symbol, found bool) {
	for scope := table.pos; scope >= 0 && !found; scope-- {
		sym, found = table.lookup_unsafe(scope, key.String())
	}
	return
}

// looks up type constructor symbol with key `key`. Returns symbol (if found) and the truthy of
// whether it was found
func (table *SymbolTable) LookupConstructor(key fmt.Stringer) (sym ConstructorSymbol, found bool) {
	sym, found = table.constructors[key.String()]
	return
}

// looks up type symbol with key `key`. Returns symbol (if found) and the truthy of whether it was
// found
func (table *SymbolTable) LookupType(key fmt.Stringer) (sym TypeSymbol, found bool) {
	sym, found = table.types[key.String()]
	return
}

// makes type alias symbol
func MakeAlias(alias TypeConstant, accessors []Ident, aliased TypeConstant) TypeSymbol {
	// create location params
	loc := make([]types.Monotype, 0, len(accessors))
	for _, accessor := range accessors {
		loc = append(loc, types.Variable(accessor.String()))
	}

	return TypeSymbol{
		isAlias:     true,
		Application: types.Application{TypeConst: types.Constant(aliased.String()), Args: loc},
	}
}

// makes symbol
func MakeSymbol(ident Ident, binding Binding, prec uint8, arity uint) Symbol {
	return Symbol{
		Id:         ident,
		Binding:    binding,
		Precedence: prec,
		Arity:      arity,
	}
}

// makes trait def symbol
//
// NOTE: uses the decs slice, not a copy of it!
func MakeTrait(name TypeConstant, decs []Symbol) TraitSymbol {
	return TraitSymbol{
		Id:      name,
		Symbols: decs,
	}
}

// makes monotype symbol
//
// NOTE: uses the params slice, not a copy of it!
func MakeType(name TypeConstant, params []types.Variable) TypeSymbol {
	args := make([]types.Monotype, len(params))
	for i, p := range params {
		args[i] = types.Variable(p.String())
	}
	return TypeSymbol{
		isAlias:     false,
		Application: types.Application{TypeConst: types.Constant(name.String()), Args: args},
	}
}

// makes a new symbol table
func NewSymbolTable() (table *SymbolTable) {
	table = new(SymbolTable)
	table.imported = make(map[string]ModuleSymbol)
	table.types = make(map[string]TypeSymbol)
	table.AddScope()
	return
}

// removes the top scope from the symbol table
//
// no-op if empty
func (table *SymbolTable) RemoveScope() {
	if len(table.tables) == 0 {
		return
	}
	table.tables = table.tables[:len(table.tables)-1]
	table.pos--
}

// this call does not verify `scope`; thus, may cause panic from out of bounds access
func (table *SymbolTable) lookup_unsafe(scope int, key string) (sym Symbol, found bool) {
	sym, found = table.tables[scope][key]
	return
}
