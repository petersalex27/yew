// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package parser

import "github.com/petersalex27/yew/internal/types"


type definitionParent struct {
	// used to report errors like: number of params doesn't match, definitions split, etc.
	Name string
	Start, End int
	// this may differ from the number of params, this happens when the function returns a function
	//
	// initially this should be set to -1 to indicate that the arity is not yet known
	//
	// the first time a definition is added, this should be set to the arity of the definition and all
	// subsequent definitions should be checked against this value
	arity int64
	params   []types.Variable
	arms     []definition
}

func (def definitionParent) String() string {
	return def.Name
}

func (def definitionParent) Pos() (int, int) {
	return def.Start, def.End
}

type definition struct {
	// pattern to match representing arguments to function
	scrutinee []types.Term
	// type resulting from application of scrutinee
	typ types.Type
	// map from local functions from `where` clause to their demangled names in global scope
	local types.Locals
	// expression part of function--part bound to the name and scrutinee
	expression termElem
}

func (parser *Parser) verifyThenAppend(parent definitionParent, def definition) (definitionParent, bool) {
	parent.End = def.expression.End // update end point of parent definition

	// verify that the arity of the definition is so far consistent
	if parent.arity < 0 {
		// First definition, set arity. It's now known
		parent.arity = int64(len(def.scrutinee))
	} else if int64(len(def.scrutinee)) != parent.arity {
		// inconsistent arity
		parser.errorOn(ArityMismatch, parent)
		return parent, false
	}

	// append the definition
	parent.arms = append(parent.arms, def)
	return parent, true
}

func (parser *Parser) commitDefinition(def definitionParent) bool {
	name := parser.inParent + def.String()
	if defn, found := parser.definitions[name]; found {
		parser.errorOn(FunctionDefinitionSplit, defn)
		return false
	}

	// re-slice the params to the correct length, it should now match the arity
	//
	// arity should never be -1 here because we only get here if we have a definition
	def.params = def.params[:def.arity]
	parser.definitions[name] = def
	return true
}

func makeDefParent(name stringPos) definitionParent {
	p := definitionParent{Name: name.String()}
	p.Start, p.End = name.Pos()
	p.arity = -1 // arity is not yet known
	return p
}

func (parser *Parser) exchangeDefining(name stringPos) (definitionParent, bool) {
	if parser.defining.Empty() {
		return makeDefParent(name), true
	}

	def, _ := parser.defining.Pop()

	// is the definition we are defining the same as the last definition?
	if def.String() != name.String() { // no
		// okay, actually commit definition `def` now (the last thing defined)
		if ok := parser.commitDefinition(def); !ok {
			return definitionParent{}, false
		}

		return makeDefParent(name), true // create a new definition parent
	}

	// yes, we are defining the same definition as the last one, so just return it
	return def, true
}

func (parser *Parser) define(head types.Term, pattern []types.Term, appliedType types.Type, body termElem, locals types.Locals) bool {
	parent, ok := parser.exchangeDefining(head)
	if !ok {
		return false
	}
	// okay, now attach the definition to the parent
	// TODO: does the typ and body.GetKind() need to be unified?
	// definition arm
	def := definition{pattern, appliedType, locals, body}
	parent, ok = parser.verifyThenAppend(parent, def)
	if !ok {
		return false
	}

	// push the parent (possibly back) onto the stack
	parser.defining.Push(parent)
	return true
}