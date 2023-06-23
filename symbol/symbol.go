package symbol

import (
	"sort"
	"strconv"
	"strings"
	err "yew/error"
	info "yew/info"
	scan "yew/lex"
	types "yew/type"
)

type Location info.Location
type definedLocation struct {
	path string
	line int
	char int
}

// SymbolTable keeps track of running list of pointers to scopes' tables
type SymbolTable struct {
	sourcePath string
	table []*ScopeTable
}

func NewScopeTable() *ScopeTable {
	out := new(ScopeTable)
	*out = make(ScopeTable)
	return out
}

func InitSymbolTable(sourcePath string) *SymbolTable {
	out := new(SymbolTable)
	*out = SymbolTable{
		sourcePath: sourcePath,
		table: make([]*ScopeTable, 0, 1),
	}
	out.table = append(out.table, NewScopeTable())
	return out
}

type ScopeTable map[string]struct{symbol Symbolic; isDefined bool}

func NewAddableSymbol(sym Symbolic) struct{symbol Symbolic; isDefined bool} {
	return struct{symbol Symbolic; isDefined bool}{sym, false}
}

func (s *SymbolTable) NewScope() {
	s.AddScope(NewScopeTable())
}
func (s *SymbolTable) AddScope(sc *ScopeTable) {
	// pushes scope to front of array
	s.table = append([]*ScopeTable{sc}, s.table...)
}
func (s *SymbolTable) RemoveScope() {
	// pops scope from front of array
	s.table = s.table[1:]
}

func locToDefLoc(loc Location) definedLocation {
	return definedLocation{
		path: loc.GetPath(),
		line: loc.GetLine(),
		char: loc.GetChar(),
	}
}

func MakeLocation(path string, line, char int) Location {
	return definedLocation{path: path, line: line, char: char}
}

func MakeSymbol_testable(id string, t types.Types, loc Location, uses map[string]SymbolUse) Symbol {
	symbol := Symbol{}
	symbol.idToken = scan.MakeIdToken(id, loc.GetLine(), loc.GetChar())
	symbol.ofType = t
	//symbol.definition = definedLocation{path: loc.GetPath(), line: loc.GetLine(), char: loc.GetChar()}
	symbol.uses = uses
	return symbol
}

func MakeSymbol(id scan.IdToken) (symbol Symbol) {
	symbol = Symbol{}
	symbol.idToken = id
	symbol.ofType = types.GetNewTau()
	//symbol.definition = locToDefLoc(id.GetLocation())
	symbol.uses = make(map[string]SymbolUse)
	return
}

// HasErrorAttached returns true iff `s` has an error attached
func HasErrorAttached(s Symbolic) bool {
	return types.ERROR == s.GetType().GetTypeType()
}

// GetAttachedError returns attached error when one is attached, else nil
func  GetAttachedError(s Symbolic) *err.Error {
	if HasErrorAttached(s) {
		e := s.GetType().(types.Error).ToError()
		return &e
	}
	return nil
}

// SearchLocal searches the most local scope of the symbol table for the symbol mapped to by
// key. When found, `sym` is the found symbol and `found` is true; otherwise, `sym` is nil
// and `found` is false.
func (s *SymbolTable) SearchLocal(key string) (Symbolic, bool) {
	sym, found := (*s.table[0])[key]
	if found {
		return sym.symbol, found
	}
	return Symbol{}, false
} 

func (s *SymbolTable) addToLocal(symbol Symbolic) {
	(*s.table[0])[symbol.GetIdToken().ToString()] = NewAddableSymbol(symbol)
}

type nameErrorType int
const (
	E_ILLEGAL_REDEF nameErrorType = iota
	E_ILLEGAL_REDEC
	E_TYPE_MISMATCH 
)
var NameErrorMessages = map[nameErrorType]string {
	E_ILLEGAL_REDEF: "illegal redefinition",
	E_ILLEGAL_REDEC: "illegal redeclaration",
	E_TYPE_MISMATCH: "type mismatch with",
}

func MakeRedefinedError(original Symbolic, newSymbol Symbolic) err.Error {
	/* = outline ================================= */
	// [<file>:<line>:<char>] Name Error: illegal redefinition.
	// <source>
	// ^ 
	// previous definition of <name> at <file>:<line>:<char>
	// <source>
	// ^

	loc := newSymbol.GetIdToken().GetLocation()
	messageAppend := 
			" " + original.GetIdToken().ToString() + " previously defined at " + 
			original.GetIdToken().GetLocation().ToString() + "."
	return err.CompileMessage(
			NameErrorMessages[E_ILLEGAL_REDEF] + messageAppend, err.ERROR, err.NAME, 
			loc.GetPath(), 
			loc.GetLine(),
			loc.GetChar(), 0, "").(err.Error)
}

func MakeRedeclareError(original Symbolic, newSymbol Symbolic) err.Error {
	loc := newSymbol.GetIdToken().GetLocation()
	messageAppend := 
			" " + original.GetIdToken().ToString() + " previously declared at " + 
			original.GetIdToken().GetLocation().ToString() + "."
	return err.CompileMessage(
			NameErrorMessages[E_ILLEGAL_REDEC] + messageAppend, err.ERROR, err.NAME, 
			loc.GetPath(), 
			loc.GetLine(),
			loc.GetChar(), 0, "").(err.Error)
}

func MakeTypeMismatchError(original Symbolic, newSymbol Symbolic) err.Error {
	loc := newSymbol.GetIdToken().GetLocation()
	messageAppend := 
			" " + original.GetIdToken().ToString() + " at " + 
			original.GetIdToken().GetLocation().ToString() + "."
	return err.TYPE.CompileMessage(NameErrorMessages[E_TYPE_MISMATCH] + messageAppend, err.ERROR,
			loc.GetPath(), loc.GetLine(), loc.GetChar(), 0, "").(err.Error)
}

func (s *SymbolTable) GetElseAdd(id Symbolic) (addedSymbol Symbolic, added bool) {
	key := id.GetIdToken().ToString()

	addedSymbol = s.Get(key)
	if nil != addedSymbol {
		return addedSymbol, false
	}

	e, added := s.AddSymbol(id)
	if !added {
		addedSymbol = id.SetType(types.Error(e))
		return addedSymbol, added // returns error
	}
	return id, added
}

func addHelperGen(errFn func(Symbolic, Symbolic) err.Error) (func(*SymbolTable, Symbolic)(types.Error, bool)) {
	return func(s *SymbolTable, symbol Symbolic) (addError types.Error, added bool) {
		sym, found := s.SearchLocal(symbol.GetIdToken().ToString())
		added = !found
		if found {
			addError = types.Error(errFn(sym, symbol))
			return // error
		}
		
		// else, add symbol to table
		s.addToLocal(symbol)
		return
	}
}

var addHelper_ = addHelperGen(MakeRedefinedError)
var decHelper_ = addHelperGen(MakeRedeclareError)

func (s *SymbolTable) addHelper(symbol Symbolic) (addError types.Error, added bool) {
	return addHelper_(s, symbol)
}

// Add returns added symbol on success, else return symbol with an error attached. 
// Also returns true when added to symbol table, else false.
func (s *SymbolTable) Add_(id scan.IdToken) (addedSymbol Symbolic, added bool) {
	// create new symbol
	addedSymbol = MakeSymbol(id)

	var addError types.Error
	addError, added = s.addHelper(addedSymbol)
	if !added {
		addedSymbol = addedSymbol.SetType(types.Error(addError))
	}
	return
}

func (s *SymbolTable) DefineSymbol(symbol Symbolic) (types.Error, bool) {
	sym, stat := s.setDefined(symbol.GetIdToken().ToString())
	switch stat {
	case setDefinedStatus_alreadyDefined:
		return types.Error(MakeRedefinedError(sym, symbol)), false 
	case setDefinedStatus_ok:
		return types.Error{}, true
	case setDefinedStatus_notFound: 
		// add symbol and try to define it again
		e, added := s.AddSymbol(symbol)
		if !added {
			return e, added
		}
		_, stat = s.setDefined(symbol.GetIdToken().ToString())
		if stat == setDefinedStatus_ok {
			return e, true
		}
		fallthrough // something went very wrong ...
	default:
		err.PrintBug()
		panic("")
	}
}

func (sc *ScopeTable) updateSymbol(symbol Symbolic) (types.Error, bool) {
	sy, found := (*sc)[symbol.GetIdToken().ToString()]
	if !found {
		(*sc)[symbol.GetIdToken().ToString()] = 
				struct{symbol Symbolic; isDefined bool}{symbol, false}
	} else {
		// attach any type information
		ty := types.DoTypeInference(sy.symbol.GetType(), symbol.GetType())
		if ty.GetTypeType() == types.ERROR {
			return types.Error(MakeTypeMismatchError(sy.symbol, symbol)), false
		}
		tmp := struct{symbol Symbolic; isDefined bool}{symbol.SetType(ty), sy.isDefined}
		(*sc)[symbol.GetIdToken().ToString()] = tmp
	}
	return types.Error{}, true
} 

// see GetScopeAtDistance--does the same thing but performs no bounds checking
func (s *SymbolTable) getScopeAtDistance_noCheck(n uint) *ScopeTable {
	return (*s).table[n]
}

// attempts to get scope at a distance of `n` from the local scope
func (s *SymbolTable) GetScopeAtDistance(n uint) (scope *ScopeTable, exists bool) {
	if n >= uint(len(s.table)) {
		return nil, false
	}
	return s.getScopeAtDistance_noCheck(n), true
}

func (s *SymbolTable) GetGlobalScope() *ScopeTable {
	return s.getScopeAtDistance_noCheck(uint(len(s.table) - 1))
}

func (s *SymbolTable) GetLocalScope() *ScopeTable {
	return s.getScopeAtDistance_noCheck(0)
}

func (s *SymbolTable) AddSymbolToGlobal(symbol Symbolic) (types.Error, bool) {
	return s.GetGlobalScope().updateSymbol(symbol)
}

// declares a new symbol in the current local scope. This should only be called for 
// explicit declarations.
func (s *SymbolTable) DeclareLocal(symbol Symbolic, ty types.Types) (types.Error, bool) {
	symbol = symbol.SetType(ty)
	return s.GetLocalScope().updateSymbol(symbol)
}

func (s *SymbolTable) IsDefined(symbol Symbolic) bool {
	_, _, defd := s.get(symbol.GetIdToken().ToString())
	return defd
}

func (s *SymbolTable) AddSymbol(symbol Symbolic) (types.Error, bool) {
	return s.addHelper(symbol)
}

type setDefinedStatus int
const (
	setDefinedStatus_notFound setDefinedStatus = iota
	setDefinedStatus_alreadyDefined
	setDefinedStatus_ok
)

// for n >= 0: (n -> n) and (-n -> len - n)
func setHighRange(len int, high int) (out int) {
	if high < 0 {
		out = len + 1 + high // possibly less than zero
	} else {
		out = high
	}

	return
}

func (s *SymbolTable) setDefinedInRange(key string, low int, high int) (Symbolic, setDefinedStatus) {
	high = setHighRange(len(s.table), high)

	for i, table := range s.table {
		if i < low {
			continue
		} else if i >= high {
			break
		}

		sym, found := (*table)[key]
		if found {
			if sym.isDefined {
				return sym.symbol, setDefinedStatus_alreadyDefined
			}
			sym.isDefined = true
			(*table)[key] = sym
			return sym.symbol, setDefinedStatus_ok
		}
	}
	return nil, setDefinedStatus_notFound
}

func (s *SymbolTable) setDefined(key string) (Symbolic, setDefinedStatus) {
	return s.setDefinedInRange(key, 0, -1)
}

func (s *SymbolTable) get(key string) (symbol Symbolic, tableIndex int, isDefined bool) {
	// loop from most current scope (i.e., table at index zero) through global scope
	for i, table := range s.table {
		sym, found := (*table)[key]
		if found {
			return sym.symbol, i, sym.isDefined
		}
	}

	return nil, -1, false  // symbol not found
}

// assumes sym exists!
func (s *SymbolTable) Update(sym Symbolic) Symbolic {
	s.addToLocal(sym)
	return sym
}

// Get returns symbol on success, else nil
func (s *SymbolTable) Get(key string) Symbolic {
	sy, _, _ := s.get(key)
	return sy
}

type locations []Location
func (a locations) Len() int { return len(a) }
func (a locations) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a locations) Less(i, j int) bool { 
	if a[i].GetLine() == a[j].GetLine() {
		return a[i].GetChar() < a[j].GetChar()
	}
	return a[i].GetLine() < a[j].GetLine()
}

func (s *Symbol) UpdateUse(loc Location) {
	_, found := s.uses[loc.GetPath()]
	if !found {
		tmp := SymbolUse{
			path: loc.GetPath(),
		}
		tmp.locations = make([]Location, 0, 1)
		tmp.locations = append(tmp.locations, loc)
		s.uses[loc.GetPath()] = tmp
		return
	}

	locs := s.uses[loc.GetPath()].locations
	_, found = sort.Find(len(locs), func(i int) int {
		return locs[i].GetLine() - loc.GetLine()
	})
	if !found {
		// update
		locs = append(locs, loc)
		sort.Sort(locations(locs))
		s.uses[loc.GetPath()] = SymbolUse{
			path: loc.GetPath(), 
			locations: locs,
		}
		return
	}
}

type SymbolUse struct {
	path string
	locations []Location
}

type Symbolic interface {
	GetIdToken() scan.IdToken
	GetType() types.Types
	SetType(types.Types) Symbolic
	IsDefined() bool
}

type Symbol struct {
	idToken scan.IdToken
	//definition definedLocation
	//id string // not demangled
	demangler string // append to demangle
	isDefined bool
	ofType types.Types
	uses map[string]SymbolUse
}

func (s *Symbol) GetFullName() string {
	return s.demangler + s.idToken.ToString()
}

func (s Symbol) GetIdToken() scan.IdToken {
	return s.idToken
}

func (s Symbol) GetType() types.Types {
	return s.ofType
} 

func (s Symbol) SetType(ty types.Types) Symbolic { 
	s.ofType = ty 
	return s
}
func (s *Symbol) SetDemangler(demangle string) { s.demangler = demangle}
func (s Symbol) SetLocation(line int, char int, path string) Symbol {
	if s.IsDefined() {
		return s
	}
	s.idToken = scan.MakeIdToken(
			s.idToken.ToString(), 
			s.idToken.GetLocation().GetLine(), 
			s.idToken.GetLocation().GetChar())
	return s
	//s.definition.line = line
	//s.definition.char = char
	//s.definition.path = path
}

func (d definedLocation) GetChar() int { return d.char }
func (d definedLocation) GetLine() int { return d.line }
func (d definedLocation) GetPath() string { return d.path }

func (s *Symbol) GetDefinedLocation() (location Location, isDef bool) {
	return s.idToken.GetLocation(), s.IsDefined()
}

func (d definedLocation) isDefined() bool { return d.path != "" }

func (s Symbol) IsDefined() bool {
	return s.isDefined
} 

var maxDisplayedUsesPerLocation = 8
var maxDisplayedUsesPerSymbol = 4

func (d definedLocation) ToString() string {
	return strconv.Itoa(d.line) + ":" + strconv.Itoa(d.char)
}

func (s *Symbol) stringifyDefinedLocation() string {
	return s.idToken.GetLocation().ToString()
} 

/*
creates the following:
"{definition: <def>; id: <id>; demangler: <id>; ofType: <type>; uses: map[string]SymbolUse}"
*/
func (s *Symbol) ToString() string {
	return "{" +
		"definition: " + s.idToken.GetLocation().ToString() + "; " +
		"idToken: " + s.idToken.ToString() + "; " +
		"demangler: " + s.demangler + "; " +
		"ofType: " + s.ofType.ToString() + "; " + 
		"uses: map[string]SymbolUse}" 
}

func (use SymbolUse) ToString() string {
	var builder strings.Builder 
	builder.WriteString("[" + use.path + "] ")
	leftOver := 0
	loopN := len(use.locations)
	if loopN > maxDisplayedUsesPerLocation {
		leftOver = loopN - maxDisplayedUsesPerLocation
		loopN = maxDisplayedUsesPerLocation
	}

	stringifyLocation := func(loc Location) string {
		return strconv.Itoa(loc.GetLine()) + ":" + strconv.Itoa(loc.GetChar())
	}
	
	if loopN > 0 {
		builder.WriteString(stringifyLocation(use.locations[0]))
	}
	for i := 1; i < loopN; i++ {
		builder.WriteString(", ")
		builder.WriteString(stringifyLocation(use.locations[i]))
	}

	if leftOver > 0 {
		builder.WriteString(" (.. +")
		builder.WriteString(strconv.Itoa(leftOver))
		builder.WriteString(" others)")
	}
	return builder.String()
}

func (s *Symbol) stringifyUses() string {
	var builder strings.Builder 
	leftOver := 0
	loopN := len(s.uses)
	if loopN > maxDisplayedUsesPerSymbol {
		leftOver = loopN - maxDisplayedUsesPerSymbol
		loopN = maxDisplayedUsesPerSymbol
	}

	keys := make([]string, 0, loopN)
	for k := range s.uses {
		keys = append(keys, k)
	}

	if loopN > 0 {
		builder.WriteString(s.uses[keys[0]].ToString())
	} else {
		builder.WriteString("")
	}
	for i := 1; i < loopN; i++ {
		builder.WriteString("\n")
		builder.WriteString(s.uses[keys[i]].ToString())
	}

	if leftOver > 0 {
		builder.WriteString("\n[.. +")
		builder.WriteString(strconv.Itoa(leftOver))
		builder.WriteString(" others)")
	}
	return builder.String()
}