package parser

import (
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/common/data"
)

func parseModulePartOfHeader(p Parser) data.Either[data.Ers, data.Maybe[module]] {
	origin := getOrigin(p)
	annotEs, mAnnots, isMAnnots := parseAnnotations_(p).Break()
	if !isMAnnots {
		return data.PassErs[data.Maybe[module]](annotEs)
	}

	es, mMod, ok := parseModule(p).Break()
	if !ok {
		return data.PassErs[data.Maybe[module]](es)
	}

	if m, just := mMod.Break(); !just {
		resetOrigin(p, origin)
	} else {
		// apply annotations to module
		m.annotate(mAnnots)
		mMod = data.Just(m)
	}
	return data.Ok(mMod)
}

// rule:
//
//	```
//	header =
//		[[annotations_], module, {then, [annotations_], import statement}]
//		| [[annotations_], import statement, {then, [annotations_], import statement}] ;
//	```
func parseHeader(p Parser) data.Either[data.Ers, data.Maybe[header]] {
	es, mMod, isMMod := parseModulePartOfHeader(p).Break()
	if !isMMod {
		return data.PassErs[data.Maybe[header]](es)
	}

	// reset when no imports are found
	origin := getOrigin(p)
	
	// the origin WILL change if and only if endHeaderEarly is false!!
	endHeaderEarly := !mMod.IsNothing() && !then(p) // order is REALLY important here
	if endHeaderEarly { // origin is 
		// module but no `then` => no imports
		return data.Ok(makePossibleHeader(mMod, data.Nil[importStatement]()))
	} 
	
	// might need to reset origin--don't know yet though

	// no module or module and `then` are found successful; either way, input is ready to optionally 
	// parse the import statements
	es, imports, isImports := parseImports(p).Break()
	if !isImports {
		return data.PassErs[data.Maybe[header]](es)
	}

	// if imports are NOT found, and b/c the `then` might have succeeded above, if the origin has
	// changed, the `then` sequences a body element or the footer instead; so, reset the origin
	if imports.IsEmpty() {
		resetOrigin(p, origin)
	}

	return data.Ok(makePossibleHeader(mMod, imports))
}

func makePossibleHeader(mod data.Maybe[module], imports data.List[importStatement]) data.Maybe[header] {
	if mod.IsNothing() && imports.IsEmpty() {
		return data.Nothing[header]()
	}
	return data.Just(data.EMakePair[header](mod, imports))
}

func parseModule(p Parser) data.Either[data.Ers, data.Maybe[module]] {
	moduleToken, found := getKeywordAtCurrent(p, token.Module, dropBeforeAndAfter)
	if !found {
		return data.Ok(data.Nothing[module](p))
	}

	id, just := parseLowerIdent(p).Break()
	if !just {
		return data.Fail[data.Maybe[module]](ExpectedModuleId, p)
	}

	mod := makeModule(id)
	mod.Position = mod.Position.Update(moduleToken)
	return data.Ok(data.Just(mod))
}

// Note, this is not an actual rule in the grammar, but a helper function to parse the imports. That
// said, if it was a rule, it would be the following:
//
//	```
//	imports helper = import statement, {then, import statement} ;
//	```
func parseImportsHelper(p Parser, as data.Maybe[annotations]) data.Either[data.Ers, data.List[importStatement]] {
	importStatements := data.Nil[importStatement]()
	// only reset if nothing is returned for `maybeParseImport`
	origin := getOrigin(p)
	for {
		var es data.Ers
		var isMAnnots bool
		if isMAnnots = !as.IsNothing(); !isMAnnots { // this allows the argument to be used if non-Nothing
			es, as, isMAnnots = parseAnnotations_(p).Break()
		}

		if !isMAnnots {
			return data.PassErs[data.List[importStatement]](es)
		} else if pEs, mImport := maybeParseImport(p); pEs != nil {
			return data.PassErs[data.List[importStatement]](*pEs)
		} else if im, just := mImport.Break(); !just {
			break
		} else {
			stmt := data.EMakePair[importStatement](as, im)
			importStatements = importStatements.Snoc(stmt) // keep position, don't reset
			as = data.Nothing[annotations]()               // clear annotations for next import statement
			// set origin to current position (before `then` and the possible annotations)
			origin = getOrigin(p) 
		}

		if !then(p) {
			break
		}
	}

	resetOrigin(p, origin)
	return data.Ok(importStatements)
}

// Note, this is not an actual rule in the grammar, but a helper function to parse the imports. That
// said, if it was a rule, it would be the following:
//
//	```
//	imports = import statement, {then, import statement} ;
//	```
func parseImports(p Parser) data.Either[data.Ers, data.List[importStatement]] {
	return parseImportsHelper(p, data.Nothing[annotations]())
}

// rule:
//
//	```
//	import = "import", {"\n"},
//		( package import
//		| "(", {"\n"}, package import, {then, package import}, {"\n"}, ")"
//		) ;
func maybeParseImport(p Parser) (*data.Ers, data.Maybe[importing]) {
	importToken, found := getKeywordAtCurrent(p, token.Import, dropAfter)
	if !found {
		return nil, data.Nothing[importing](p) // no imports, this is okay, return empty list
	}

	es, ims, isImports := parseGroup[importing](p, ExpectedImportPath, maybeParsePackageImport).Break()
	if !isImports {
		return &es, data.Nothing[importing](p)
	}
	ims.Position = ims.Update(importToken)
	return nil, data.Just(ims)
}

// wraps call to `maybeParseName`, always returns `nil` for the return value
func parseMaybeName(p Parser) (*data.Ers, data.Maybe[name]) {
	return nil, maybeParseName(p)
}

// rule:
//
//	```
//	symbol selection group =
//		"_"
//		| name
//		| "(", {"\n"}, name, {{"\n"}, ",", {"\n"}, name}, [{"\n"}, ","], {"\n"}, ")" ;
//	```
func parseSymbolSelections(p Parser) data.Either[data.Ers, data.Maybe[data.NonEmpty[name]]] {
	// check for special "_" case, hides all exported symbols
	if underscore, found := getKeywordAtCurrent(p, token.Underscore, dropNone); found {
		p.advance()
		hiddenSelection := data.Nothing[data.NonEmpty[name]]().Update(underscore)
		return data.Ok(hiddenSelection) // hides all names from imported namespace
	}

	type group struct{ data.NonEmpty[name] }
	lparen, found := getKeywordAtCurrent(p, token.LeftParen, dropAfter)
	es, symbols, ok := parseSepSequenced[group](p, IllegalEmptyUsingClause, token.Comma, parseMaybeName).Break()
	if !ok { // error occurred
		return data.PassErs[data.Maybe[data.NonEmpty[name]]](es)
	}

	if found {
		rparen, found := getKeywordAtCurrent(p, token.RightParen, dropBefore)
		if !found {
			return data.Fail[data.Maybe[data.NonEmpty[name]]](ExpectedRightParen, p)
		}
		symbols.Position = symbols.Update(lparen).Update(rparen)
	} else if symbols.NonEmpty.Len() > 1 {
		return data.Fail[data.Maybe[data.NonEmpty[name]]](IllegalUnenclosedUsingClause, p)
	}

	return data.Ok(data.Just(symbols.NonEmpty))
}

// rule:
//
//	```
//	import specification = as clause | using clause ;
//		as clause = "as", {"\n"}, lower ident ;
//		using clause = "using", {"\n"}, "_" | symbol selection group ;
//	```
func maybeParseImportSpecification(p Parser) (*data.Ers, data.Maybe[selections]) {
	if as, foundAs := getKeywordAtCurrent(p, token.As, dropAfter); foundAs {
		id, ok := parseLowerIdent(p).Break()
		if !ok {
			e := data.Nil[data.Err](1).Snoc(data.MkErr(ExpectedNamespaceAlias, p))
			return &e, data.Nothing[selections](p)
		}

		asClause := data.Inl[data.Maybe[data.NonEmpty[name]]](id).Update(as)
		return nil, data.Just(asClause)
	}

	using, foundUsing := getKeywordAtCurrent(p, token.Using, dropAfter)
	if !foundUsing {
		return nil, data.Nothing[selections](p) // no import specification, this is okay, return empty list
	}

	// parse symbol selections
	es, mSymbols, isMSymbols := parseSymbolSelections(p).Break()
	if !isMSymbols {
		return &es, data.Nothing[selections](p)
	}
	sel := data.Inr[lowerIdent](mSymbols).Update(using)
	return nil, data.Just(sel)
}

// rule:
//
//	```
//	package import = import path, [{"\n"}, import specification] ;
//	```
func maybeParsePackageImport(p Parser) (*data.Ers, data.Maybe[packageImport]) {
	pathLiteral, found := getKeywordAtCurrent(p, token.ImportPath, dropNone)
	if !found {
		return nil, data.Nothing[packageImport](p)
	}

	path := data.EOne[importPathIdent](pathLiteral)

	origin := getOrigin(p)
	p.dropNewlines()

	es, selections := maybeParseImportSpecification(p)
	if es != nil {
		return es, data.Nothing[packageImport](p)
	} else if selections.IsNothing() {
		resetOrigin(p, origin)
	}

	pi := data.EMakePair[packageImport](path, selections)
	return nil, data.Just(pi)
}
