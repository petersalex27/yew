package parser

// import (
// 	"github.com/petersalex27/yew/api/util/fun"
// 	"github.com/petersalex27/yew/common/data"
// )

// func processWhereClause(wc whereClause, ys YewSource, table *Table) (out whereClause, err error) {
// 	table.Push(wc.TableEntry)
// 	defer func() { out.TableEntry, _ = table.Pop() }()

// 	out.NonEmpty = data.MapNonEmpty(func(me mainElement) mainElement {
// 		if err != nil {
// 			return me // short-circuit
// 		}

// 		me, err = me.processMainElement(ys, table)
// 		return me
// 	})(wc.NonEmpty)

// 	return out, err
// }

// func processWithClause(wc withClause, ys YewSource, table *Table) (out withClause, err error) {
// 	table.Push(wc.TableEntry)
// 	defer func() { out.TableEntry, _ = table.Pop() }()

// 	pattern, arms := wc.Split()
// 	pattern, _, err = pattern.processPattern(ys, table, 0)
// 	if err != nil {
// 		return wc, err
// 	}

// 	arms = data.EmbMapNonEmpty[withClauseArms](
// 		func(arm withClauseArm) withClauseArm {
// 			table.Push(arm.TableEntry)
// 			defer func() { arm.TableEntry, _ = table.Pop() }()

// 			lhs, db := arm.Split()
// 			lhs = data.MapEither(
// 				func(p Pattern) Pattern {
// 					p, _, err = p.processPattern(ys, table, 0)
// 					return p
// 				},
// 				func(ps data.Pair[Pattern, Pattern]) data.Pair[Pattern, Pattern] {
// 					p0, p1 := ps.Split()
// 					p0, _, err = p0.processPattern(ys, table, 0)
// 					if err != nil {
// 						return ps
// 					}
// 					p1, _, err = p1.processPattern(ys, table, 0)
// 					return data.MakePair(p0, p1)
// 				},
// 			)(lhs)

// 			if err != nil {
// 				return arm
// 			}

// 			db, err = processDefBody(db, ys, table)
// 			arm.Pair = data.MakePair(lhs, db)
// 			return arm
// 		},
// 	)(arms.NonEmpty)

// 	out.Pair = data.MakePair(pattern, arms)
// 	return out, err
// }

// func processDefBody(db defBody, ys YewSource, table *Table) (_ defBody, err error) {
// 	return data.EmbMapEither[defBody](fun.Id, func(dbp defBodyPossible) defBodyPossible {
// 		main, mWhereClause := dbp.Split()
// 		// bring the where clause into scope
// 		whereClause, isWhereClause := mWhereClause.Break()
// 		if isWhereClause {
// 			whereClause, err = processWhereClause(whereClause, ys, table)
// 			if err != nil {
// 				return dbp
// 			}

// 			table.Push(whereClause.TableEntry)
// 			mWhereClause = data.Just(whereClause)
// 		}

// 		wc, expr, isExpr := main.Break()
// 		if isExpr {
// 			expr, _, err = expr.processExpr(ys, table, 0)
// 			main = data.Inr[withClause](expr)
// 		} else {
// 			wc, err = processWithClause(wc, ys, table)
// 			main = data.Inl[Expr](wc)
// 		}

// 		if isWhereClause {
// 			whereClause.TableEntry, _ = table.Pop() // pop the where clause table entry
// 			mWhereClause = data.Just(whereClause)
// 		}

// 		return data.EMakePair[defBodyPossible](main, mWhereClause)
// 	})(db), err
// }

// func (ce caseExpr) processExpr(ys YewSource, table *Table, lbp byte) (_ Expr, bp byte, err error) {
// 	// deconstruct the case expression
// 	pattern, arms := ce.Split()
// 	pattern, _, err = pattern.processPattern(ys, table, 0)
// 	if err != nil {
// 		return ce, 0, err
// 	}

// 	arms = data.EmbMapNonEmpty[caseArms](func(ca caseArm) (out caseArm) {
// 		if err != nil {
// 			return ca // short-circuit
// 		}

// 		table.Push(ca.TableEntry)
// 		defer func() { out.TableEntry, _ = table.Pop() }()

// 		pattern, body := ca.Split()
// 		pattern, _, err = pattern.processPattern(ys, table, 0)
// 		if err != nil {
// 			return ca
// 		}

// 		body, err = processDefBody(body, ys, table)
// 		out.Pair = data.MakePair(pattern, body)
// 		return out
// 	})(arms.NonEmpty)

// 	return data.EMakePair[caseExpr](pattern, arms), 0, err
// }

// func (ea ExprApp) processExpr(ys YewSource, table *Table, _ byte) (_ Expr, bp byte, err error) {
// 	first, rest := ea.Split()
// 	first, bp, err = first.processExpr(ys, table, 0)
// 	xs := []Expr{first}
	
// 	// application, then Snoc
// 	// rbp is stronger, parse right first, then Snoc
// 	for _, e := range rest.Elements() {
// 		var rbp byte
// 		e, rbp, err = e.processExpr(ys, table, bp)
// 		if err != nil {
// 			return ea, 0, err
// 		}

// 		if bp >= rbp {
// 			xs = append(xs, e)
// 		}
// 	}
// }

// func (lit aLiteral[any]) processExpr(ys YewSource, table *Table, lbp byte) (_ Expr, bp byte, err error) {
// 	return lit, 0, nil
// }

// func (h hole) processExpr(ys YewSource, table *Table, lbp byte) (_ Expr, bp byte, err error) {
// 	return h, 0, nil
// }

// func (n name) processExpr(ys YewSource, table *Table, lbp byte) (_ Expr, bp byte, err error) {
// 	res, found := table.Get(n.Break().String())
// 	if !found {
// 		err = undefinedNameError(ys.srcCode, n.Break().String(), n)
// 	}
// 	bp = res.symbol.Bp()
// 	return n, 0, err
// }

// func (al accessList) processExpr(ys YewSource, table *Table, lbp byte) (_ Expr, bp byte, err error) {

// }

// func (la LambdaAbstraction) processExpr(ys YewSource, table *Table, lbp byte) (_ Expr, bp byte, err error) {

// }

// func (le LetExpr) processExpr(ys YewSource, table *Table, lbp byte) (_ Expr, bp byte, err error) {

// }

// func (t TypeApp) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (ct constrainedType) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (et enclosedType) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (f forallType) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (lit aLiteral[any]) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (h hole) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (n name) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (al accessList) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (la LambdaAbstraction) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (ft functionType) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (it innerTypeSeq) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (it innerTyping) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (it implicitTyping) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (ut unitType) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (w wildcard) processTyp(ys YewSource, table *Table, lbp byte) (_ Typ, bp byte, err error) {

// }

// func (p patternApp) processPattern(ys YewSource, table *Table, lbp byte) (_ Pattern, bp byte, err error) {

// }

// func (pe patternEnclosed) processPattern(ys YewSource, table *Table, lbp byte) (_ Pattern, bp byte, err error) {

// }

// func (lit aLiteral[any]) processPattern(ys YewSource, table *Table, lbp byte) (_ Pattern, bp byte, err error) {

// }

// func (h hole) processPattern(ys YewSource, table *Table, lbp byte) (_ Pattern, bp byte, err error) {

// }

// func (n name) processPattern(ys YewSource, table *Table, lbp byte) (_ Pattern, bp byte, err error) {

// }

// func (al accessList) processPattern(ys YewSource, table *Table, lbp byte) (_ Pattern, bp byte, err error) {

// }

// func (w wildcard) processPattern(ys YewSource, table *Table, lbp byte) (_ Pattern, bp byte, err error) {

// }

// func (si specInst) processVbe(ys YewSource, table *Table) (_ VisibleBodyElement, e error) {

// }

// func (sd specDef) processVbe(ys YewSource, table *Table) (_ VisibleBodyElement, e error) {

// }

// func (t typing) processVbe(ys YewSource, table *Table) (_ VisibleBodyElement, e error) {

// }

// func (td typeDef) processVbe(ys YewSource, table *Table) (_ VisibleBodyElement, e error) {

// }

// func (ta typeAlias) processVbe(ys YewSource, table *Table) (_ VisibleBodyElement, e error) {

// }

// func (s syntax) processVbe(ys YewSource, table *Table) (_ VisibleBodyElement, e error) {

// }

// func (d Def) processDef(ys YewSource, table *Table) (_ Def, err error) {

// }

// func (si specInst) processMainElement(ys YewSource, table *Table) (_ mainElement, e error) {
// 	return si.processVbe(ys, table)
// }

// func (sd specDef) processMainElement(ys YewSource, table *Table) (_ mainElement, e error) {
// 	return sd.processVbe(ys, table)
// }

// func (t typing) processMainElement(ys YewSource, table *Table) (_ mainElement, e error) {
// 	return t.processVbe(ys, table)
// }

// func (td typeDef) processMainElement(ys YewSource, table *Table) (_ mainElement, e error) {
// 	return td.processVbe(ys, table)
// }

// func (ta typeAlias) processMainElement(ys YewSource, table *Table) (_ mainElement, e error) {
// 	return ta.processVbe(ys, table)
// }

// func (s syntax) processMainElement(ys YewSource, table *Table) (_ mainElement, e error) {
// 	return s.processVbe(ys, table)
// }

// func (d Def) processMainElement(ys YewSource, table *Table) (_ mainElement, err error) {
// 	return d.processDef(ys, table)
// }
