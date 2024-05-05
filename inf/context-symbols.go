// =============================================================================
// Author-Date: Alex Peters - November 19, 2023
//
// Content: methods associated w/ Context's symbol table
//
// Notes: -
// =============================================================================
package inf

import (
	"github.com/petersalex27/yew-packages/bridge"
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/types"
)

// removes name binding from context
func (cxt *Context[N]) Remove(name expr.Const[N]) {
	key := name.Name
	sym, ok := cxt.syms.Get(key)
	if !ok {
		// TODO: do nothing, ig?
		return
	}

	// unshadow/remove sym
	remove := sym.Unshadow()
	if remove {
		// symbol is not shadowed, remove it
		cxt.syms.Remove(key)
	}
}

func (cxt *Context[N]) Add(name expr.Const[N], ty types.Type[N]) (added bool) {
	key := name.Name
	// attempt to look up existing symbol
	sym, ok := cxt.syms.Get(key)
	if ok {
		cxt.appendReport(makeNameReport("Declare Name", IllegalShadow, name))
		return false
	}
	// create new, empty symbol to be filled in
	sym = MakeSymbol[N]()
	// create (/shadow empty) symbol
	sym.Shadow(name, ty)
	// add symbol to table
	cxt.syms.Add(name.Name, sym)
	return true
}

// adds judgment to context
func (cxt *Context[N]) Shadow(name expr.Const[N], ty types.Type[N]) {
	key := name.Name
	// attempt to look up existing symbol
	sym, ok := cxt.syms.Get(key)
	if !ok {
		// no symbol in table, create new, empty symbol to be filled in
		sym = MakeSymbol[N]()
	}
	// create/shadow symbol
	sym.Shadow(name, ty)
	// add symbol to table
	cxt.syms.Add(name.Name, sym)
}

// tries to find symbol w/ name
func (cxt *Context[N]) Get(name expr.Const[N]) (judgedName bridge.JudgmentAsExpression[N, expr.Const[N]], found bool) {
	key := name.Name
	var sym Symbol[N]
	sym, found = cxt.syms.Get(key)
	if found {
		judgedName = sym.Get()
	}
	return
}
