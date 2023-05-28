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

func InitSymbolTable(sourcePath string) SymbolTable {
	out := SymbolTable{
		sourcePath: sourcePath,
		table: make([]*ScopeTable, 0, 1),
	}
	out.table = append(out.table, NewScopeTable())
	return out
}

type ScopeTable map[string]*Symbol

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

func MakeSymbol_testable(id string, t types.Types, loc Location, uses map[string]SymbolUse) *Symbol {
	symbol := new(Symbol)
	symbol.id = id
	symbol.ofType = t
	symbol.definition = definedLocation{path: loc.GetPath(), line: loc.GetLine(), char: loc.GetChar()}
	symbol.uses = uses
	return symbol
}

func MakeSymbol(id scan.IdToken) (symbol *Symbol) {
	symbol = new(Symbol)
	symbol.id = id.ToString()
	//symbol.ofType = types.GetNewTau()
	symbol.definition = locToDefLoc(id.GetLocation())
	symbol.uses = make(map[string]SymbolUse)
	return
}

// HasErrorAttatched returns true iff `s` has an error attatched
func (s *Symbol) HasErrorAttatched() bool {
	return types.ERROR == s.ofType.GetTypeType()
}

// GetAttatchedError returns attatched error when one is attatched, else nil
func (s *Symbol) GetAttatchedError() *err.Error {
	if s.HasErrorAttatched() {
		e := s.ofType.(types.Error).ToError()
		return &e
	}
	return nil
}

// SearchLocal searches the most local scope of the symbol table for the symbol mapped to by
// key. When found, `sym` is the found symbol and `found` is true; otherwise, `sym` is nil
// and `found` is false.
func (s *SymbolTable) SearchLocal(key string) (sym *Symbol, found bool) {
	sym, found = (*s.table[0])[key]
	return
} 

func (s *SymbolTable) addToLocal(symbol *Symbol) {
	(*s.table[0])[symbol.id] = symbol
}

type nameErrorType int
const (
	E_ILLEGAL_REDEF nameErrorType = iota

)
var NameErrorMessages = map[nameErrorType]string {
	E_ILLEGAL_REDEF: "illegal redefinition",
}

func MakeRedefinedError(original *Symbol, newSymbol *Symbol) err.Error {
	/* = outline ================================= */
	// [<file>:<line>:<char>] Name Error: illegal redefinition.
	// <source>
	// ^ 
	// previous definition of <name> at <file>:<line>:<char>
	// <source>
	// ^

	messageAppend := 
			" " + original.id + "previously defined at " + 
			original.stringifyDefinedLocation() + "."
	return err.CompileMessage(
			NameErrorMessages[E_ILLEGAL_REDEF] + messageAppend, err.ERROR, err.NAME, 
			newSymbol.definition.GetPath(), newSymbol.definition.GetLine(),
			newSymbol.definition.GetChar(), 0, "").(err.Error)
}

func (s *SymbolTable) GetElseAdd(id scan.IdToken) (addedSymbol *Symbol, added bool) {
	key := id.ToString()

	addedSymbol = s.Get(key)
	if nil != addedSymbol {
		return addedSymbol, false
	}

	return s.Add(id)
}

// Add returns added symbol on success, else return symbol with an error attatched. 
// Also returns true when added to symbol table, else false.
func (s *SymbolTable) Add(id scan.IdToken) (addedSymbol *Symbol, added bool) {
	key := id.ToString()

	// create new symbol
	addedSymbol = MakeSymbol(id)

	sym, found := s.SearchLocal(key)
	if found {
		addedSymbol.ofType = types.Error(MakeRedefinedError(sym, addedSymbol))
		return addedSymbol, false // error, redef
	}
	
	// add symbol to table
	s.addToLocal(addedSymbol)
	return addedSymbol, true
}

// Get returns symbol on success, else nil
func (s *SymbolTable) Get(key string) *Symbol {
	// loop from most current scope (i.e., table at index zero) through global scope
	for _, table := range s.table {
		sym, found := (*table)[key]
		if found {
			return sym
		}
	}

	return nil  // symbol not found
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

type Symbol struct {
	definition definedLocation
	id string // not demangled
	demangler string // append to demangle
	ofType types.Types
	uses map[string]SymbolUse
}

func (s *Symbol) GetFullName() string {
	return s.demangler + s.id
}

func (s *Symbol) GetType() types.Types {
	return s.ofType
} 

func (s *Symbol) SetType(ty types.Types) { s.ofType = ty }
func (s *Symbol) SetDemangler(demangle string) { s.demangler = demangle}
func (s *Symbol) SetLocation(line int, char int, path string) {
	if s.definition.isDefined() {
		return
	}
	s.definition.line = line
	s.definition.char = char
	s.definition.path = path
}

func (d definedLocation) GetChar() int { return d.char }
func (d definedLocation) GetLine() int { return d.line }
func (d definedLocation) GetPath() string { return d.path }

func (s *Symbol) GetDefinedLocation() (location Location, isDef bool) {
	return s.definition, s.definition.isDefined()
}

func (d definedLocation) isDefined() bool { return d.path != "" }

func (s Symbol) IsDefined() bool {
	return s.definition.isDefined()
} 

var maxDisplayedUsesPerLocation = 8
var maxDisplayedUsesPerSymbol = 4

func (d definedLocation) ToString() string {
	return strconv.Itoa(d.line) + ":" + strconv.Itoa(d.char)
}

func (s *Symbol) stringifyDefinedLocation() string {
	return s.definition.path + ":" + s.definition.ToString()
} 

/*
creates the following:
"{definition: <def>; id: <id>; demangler: <id>; ofType: <type>; uses: map[string]SymbolUse}"
*/
func (s *Symbol) ToString() string {
	return "{" +
		"definition: " + s.definition.ToString() + "; " +
		"id: " + s.id + "; " +
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