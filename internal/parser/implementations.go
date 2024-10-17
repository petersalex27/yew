package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
)

var (
	patternNameAsPattern = fun.Bind1stOf3(data.Cases, (hole).asPattern, (name).asPattern)
	patternNameAsTyp     = fun.Bind1stOf3(data.Cases, (hole).asTyp, (name).asTyp)
	patternNameAsExpr    = fun.Bind1stOf3(data.Cases, (hole).asExpr, (name).asExpr)
	patternAtomAsTyp     = fun.Bind1stOf3(data.Cases, (literal).asTyp, patternNameAsTyp)
	patternAtomAsExpr    = fun.Bind1stOf3(data.Cases, (literal).asExpr, patternNameAsExpr)
	exprAtomAsTyp        = fun.Bind1stOf3(data.Cases, patternAtomAsTyp, (lambdaAbstraction).asTyp)
	exprAtomAsExpr       = fun.Bind1stOf3(data.Cases, patternAtomAsExpr, (lambdaAbstraction).asExpr)
	exprAtomAsExprRes    = fun.Compose(data.Inr[data.Ers], exprAtomAsTyp)
)

// = access ========================================================================================

// access implements expr
func (expr access) updatePosExpr(p api.Positioned) expr {
	expr.Position = expr.Update(p)
	return expr
}

// access implements pattern
func (pat access) updatePosPattern(p api.Positioned) pattern {
	pat.Position = pat.Update(p)
	return pat
}

// access implements typ
func (ty access) updatePosTyp(p api.Positioned) typ {
	ty.Position = ty.Update(p)
	return ty
}

// = appType =======================================================================================

// appType implements typ
func (ty appType) updatePosTyp(p api.Positioned) typ {
	ty.Position = ty.Update(p)
	return ty
}

// = caseExpr ======================================================================================

func (e caseExpr) asExpr() expr { return expr(e) }

// caseExpr implements expr
func (expr caseExpr) updatePosExpr(p api.Positioned) expr {
	expr.Position = expr.Update(p)
	return expr
}

// = constraintUnverified ==========================================================================

// constraintUnverified implements constraint
func (c constraintUnverified) asConstraint() constraint { return c }

// = constraintVerified ============================================================================

// constraintVerified implements constraint
func (c constraintVerified) asConstraint() constraint { return c }

// = constrainedType ===============================================================================

// constrainedType implements typ
func (ty constrainedType) updatePosTyp(p api.Positioned) typ {
	ty.Position = ty.Update(p)
	return ty
}

// = def ===========================================================================================

// def implements bodyElement
func (def def) setAnnotation(as data.Maybe[annotations]) mainElement {
	def.annotate(as)
	return def
}

func (def def) asBodyElement() bodyElement { return data.EInl[bodyElement](def) }

// def implements mainElement
func (def def) asMainElement() mainElement { return def }

// = enclosedType ==================================================================================

// enclosedType implements typ
func (ty enclosedType) updatePosTyp(p api.Positioned) typ {
	ty.typ = ty.typ.updatePosTyp(p)
	return ty
}

// = exprApp =======================================================================================

// exprApp implements expr
func (expr exprApp) updatePosExpr(p api.Positioned) expr {
	expr.Position = expr.Update(p)
	return expr
}

// = forallType ====================================================================================

// forallType implements typ
func (ty forallType) updatePosTyp(p api.Positioned) typ {
	ty.Position = ty.Update(p)
	return ty
}

// = functionType ==================================================================================

// functionType implements typ
func (ty functionType) updatePosTyp(p api.Positioned) typ {
	ty.Position = ty.Update(p)
	return ty
}

// = hole ==========================================================================================

func (h hole) asPattern() pattern { return h }

func (h hole) asExpr() expr { return h }

func (h hole) asTyp() typ { return h }

// hole implements pattern
func (h hole) updatePosPattern(p api.Positioned) pattern {
	h.Position = h.Update(p)
	return h
}

// hole implements expr
func (h hole) updatePosExpr(p api.Positioned) expr {
	h.Position = h.Update(p)
	return h
}

// hole implements typ
func (h hole) updatePosTyp(p api.Positioned) typ {
	h.Position = h.Update(p)
	return h
}

// = innerTypeTerms ================================================================================

// innerTypeTerms implements typ
func (ty innerTypeTerms) updatePosTyp(p api.Positioned) typ {
	ty.Position = ty.Update(p)
	return ty
}

// = innerTyping ===================================================================================

func (ty innerTyping) asTyp() typ { return ty }

// innerTyping implements typ
func (ty innerTyping) updatePosTyp(p api.Positioned) typ {
	ty.Position = ty.Update(p)
	return ty
}

// = implicitTyping ================================================================================

// implicitTyping implements typ
func (ty implicitTyping) updatePosTyp(p api.Positioned) typ {
	ty.Position = ty.Update(p)
	return ty
}

// = lambdaAbstraction =============================================================================

func (e lambdaAbstraction) asExpr() expr { return expr(e) }

func (ty lambdaAbstraction) asTyp() typ { return typ(ty) }

// lambdaAbstraction implements expr
func (expr lambdaAbstraction) updatePosExpr(p api.Positioned) expr {
	expr.Position = expr.Update(p)
	return expr
}

// lambdaAbstraction implements typ
func (ty lambdaAbstraction) updatePosTyp(p api.Positioned) typ {
	ty.Position = ty.Update(p)
	return ty
}

// = letExpr =======================================================================================

func (e letExpr) asExpr() expr { return expr(e) }

// letExpr implements expr
func (expr letExpr) updatePosExpr(p api.Positioned) expr {
	expr.Position = expr.Update(p)
	return expr
}

// = literal =======================================================================================

func (lit literal) asPattern() pattern { return lit }

func (lit literal) asExpr() expr { return lit }

func (lit literal) asTyp() typ { return lit }

// literal implements pattern
func (lit literal) updatePosPattern(p api.Positioned) pattern {
	lit.Position = lit.Update(p)
	return lit
}

// literal implements expr
func (lit literal) updatePosExpr(p api.Positioned) expr {
	lit.Position = lit.Update(p)
	return lit
}

// literal implements typ (as a term, type checker will catch illegal uses)
func (lit literal) updatePosTyp(p api.Positioned) typ {
	lit.Position = lit.Update(p)
	return lit
}

// = name ==========================================================================================

func (n name) asPattern() pattern { return n }

func (n name) asExpr() expr { return n }

func (n name) asTyp() typ { return n }

// name implements pattern
func (n name) updatePosPattern(p api.Positioned) pattern {
	n.Position = n.Update(p)
	return n
}

// name implements expr
func (n name) updatePosExpr(p api.Positioned) expr {
	n.Position = n.Update(p)
	return n
}

// name implements typ
func (n name) updatePosTyp(p api.Positioned) typ {
	n.Position = n.Update(p)
	return n
}

// = patternApp ====================================================================================

// patternApp implements pattern
func (app patternApp) updatePosPattern(p api.Positioned) pattern {
	app.Position = app.Update(p)
	return app
}

// = patternEnclosed ===============================================================================

// patternEnclosed implements pattern
func (enclosed patternEnclosed) updatePosPattern(p api.Positioned) pattern {
	enclosed.Position = enclosed.Update(p)
	return enclosed
}

// = specDef =======================================================================================

// specDef implements bodyElement
func (spec specDef) setAnnotation(as data.Maybe[annotations]) mainElement {
	spec.annotate(as)
	return spec
}

func (spec specDef) asBodyElement() bodyElement { return data.EInr[bodyElement](visibleBodyElement(spec)) }

// specDef implements mainElement
func (spec specDef) asMainElement() mainElement { return spec }

// specDef implements visibleBodyElement
func (spec specDef) setVisibility(mv data.Maybe[visibility]) mainElement {
	spec.visibility = mv
	spec.Position = spec.Update(mv)
	return spec
}

// = specInst ======================================================================================

// specInst implements bodyElement
func (inst specInst) setAnnotation(as data.Maybe[annotations]) mainElement {
	inst.annotate(as)
	return inst
}

func (inst specInst) asBodyElement() bodyElement { return data.EInr[bodyElement](visibleBodyElement(inst)) }

// specInst implements mainElement
func (inst specInst) asMainElement() mainElement { return inst }

// specInst implements visibleBodyElement
func (inst specInst) setVisibility(mv data.Maybe[visibility]) mainElement {
	inst.visibility = mv
	inst.Position = inst.Update(mv)
	return inst
}

// = syntax ========================================================================================

// syntax implements bodyElement
func (syntax syntax) setAnnotation(as data.Maybe[annotations]) mainElement {
	syntax.annotate(as)
	return syntax
}

func (syntax syntax) asBodyElement() bodyElement { return data.EInr[bodyElement](visibleBodyElement(syntax)) }

// syntax implements mainElement
func (syntax syntax) asMainElement() mainElement { return syntax }

// syntax implements visibleBodyElement
func (syntax syntax) setVisibility(mv data.Maybe[visibility]) mainElement {
	syntax.visibility = mv
	syntax.Position = syntax.Update(mv)
	return syntax
}

// = typeAlias =====================================================================================

// typeAlias implements bodyElement
func (alias typeAlias) setAnnotation(as data.Maybe[annotations]) mainElement {
	alias.annotate(as)
	return alias
}

func (alias typeAlias) asBodyElement() bodyElement { return data.EInr[bodyElement](visibleBodyElement(alias)) }

// typeAlias implements mainElement
func (typeAlias typeAlias) asMainElement() mainElement { return typeAlias }

// typeAlias implements visibleBodyElement
func (alias typeAlias) setVisibility(mv data.Maybe[visibility]) mainElement {
	alias.visibility = mv
	alias.Position = alias.Update(mv)
	return alias
}

// = typeDef =======================================================================================

// typeDef implements bodyElement
func (typeDef typeDef) setAnnotation(as data.Maybe[annotations]) mainElement {
	typeDef.annotate(as)
	return typeDef
}

func (typeDef typeDef) asBodyElement() bodyElement { return data.EInr[bodyElement](visibleBodyElement(typeDef)) }

// typeDef implements mainElement
func (typeDef typeDef) asMainElement() mainElement { return typeDef }

// typeDef implements visibleBodyElement
func (typeDef typeDef) setVisibility(mv data.Maybe[visibility]) mainElement {
	typeDef.visibility = mv
	typeDef.Position = typeDef.Update(mv)
	return typeDef
}

// = typing ========================================================================================

func (typing typing) markAuto(auto api.Token) typing {
	typing.automatic = true
	typing.Position = typing.Position.Update(auto)
	return typing
}

// typing implements bodyElement
func (typing typing) setAnnotation(as data.Maybe[annotations]) mainElement {
	typing.annotate(as)
	return typing
}

func (typing typing) asBodyElement() bodyElement { return data.EInr[bodyElement](visibleBodyElement(typing)) }

// typing implements mainElement
func (typing typing) asMainElement() mainElement { return typing }

// typing implements visibleBodyElement
func (typing typing) setVisibility(mv data.Maybe[visibility]) mainElement {
	typing.visibility = mv
	typing.Position = typing.Update(mv)
	return typing
}

// = unitType ======================================================================================

// unitType implements typ
func (ty unitType) updatePosTyp(p api.Positioned) typ {
	ty.Position = ty.Update(p)
	return ty
}

// = wildcard ======================================================================================

// wildcard represents a wildcard type
func (ty wildcard) updatePosTyp(p api.Positioned) typ {
	ty.Position = ty.Update(p)
	return ty
}

// wildcard implements pattern
func (w wildcard) updatePosPattern(p api.Positioned) pattern {
	w.Position = w.Update(p)
	return w
}
