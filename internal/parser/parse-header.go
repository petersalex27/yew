package parser

import (
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/common/data"
)

// rule:
//
//	```
//	header = [module], {{"\n"}, [annotations_], import} ;
//	```
func parseHeader(p Parser) data.Either[data.Ers, data.Maybe[header]] {
	origin := getOrigin(p)
	annotEs, mAnnots, isMAnnots := parseAnnotations(p).Break()
	if !isMAnnots {
		return data.PassErs[data.Maybe[header]](annotEs)
	}

	es, mod, ok := parseModule(p).Break()
	if !ok {
		return data.PassErs[data.Maybe[header]](es)
	}

	var imports data.List[importStatement]
	var isImports bool

	moduleExists, somethingIsAnnotated := !mod.IsNothing(), !mAnnots.IsNothing()
	annotateModule := moduleExists && somethingIsAnnotated
	
	if annotateModule {
		m, _ := mod.Break()
		m.annotate(mAnnots) // whether annots target mod. will be determined later
		mod = data.Just(m)
	}
	
	if !annotateModule && somethingIsAnnotated {
		// give annotations to imports if no module is present
		es, imports, isImports = parseImportsHelper(p, mAnnots).Break()
	} else {
		es, imports, isImports = parseImports(p).Break()
	}

	if !isImports {
		return data.PassErs[data.Maybe[header]](es)
	} else if mod.IsNothing() && imports.IsEmpty() && !mAnnots.IsNothing() {
		// reset to the position before the annotations
		//lint:ignore SA4006 ignore unused variable warning
		p = resetOrigin(p, origin)
	}
	
	return data.Ok(makePossibleHeader(mod, imports))
}

func makePossibleHeader(mod data.Maybe[module], imports data.List[importStatement]) data.Maybe[header] {
	if mod.IsNothing() && imports.IsEmpty() {
		return data.Nothing[header]()
	}
	return data.Just(data.EMakePair[header](mod, imports))
}

func parseModule(p Parser) data.Either[data.Ers, data.Maybe[module]] {
	moduleToken, found := getKeywordAtCurrent(p, token.Module)
	if !found {
		return data.Ok(data.Nothing[module](p))
	}

	id, just := parseLowerIdent(p).Break()
	if !just {
		return data.Fail[data.Maybe[module]](ExpectedModuleId, p)
	}

	mod := module{data.Nothing[annotations](), data.One[lowerIdent](id), id.Position}
	mod.Position = mod.Update(moduleToken)
	return data.Ok(data.Just[module](mod))
}

func parseImportsHelper(p Parser, as data.Maybe[annotations]) data.Either[data.Ers, data.List[importStatement]] {
	p.dropNewlines()
	importStatements := data.Nil[importStatement]()
	for {
		// only reset if nothing is returned for `maybeParseImport`
		tokenPosition := getOrigin(p)
		var es data.Ers
		var isMAnnots bool
		if isMAnnots = !as.IsNothing(); !isMAnnots { // this allows the argument to be used if non-Nothing
			es, as, isMAnnots = parseAnnotations(p).Break()
		}

		if !isMAnnots {
			return data.PassErs[data.List[importStatement]](es)
		}

		pEs, mImport := maybeParseImport(p)
		if pEs != nil {
			return data.PassErs[data.List[importStatement]](*pEs)
		}

		im, just := mImport.Break()
		if !just {
			// reset to the position before the annotations

			// ignored linter b/c this serves as a reminder that the position is being reset
			//lint:ignore SA4006 ignore unused variable warning
			p = resetOrigin(p, tokenPosition)

			// no more imports, annot, if exists, is for body element or footer
			return data.Ok(importStatements)
		}

		stmt := data.EMakePair[importStatement](as, im)
		importStatements = importStatements.Snoc(stmt) // keep position, don't reset
		// clear annotations for next import statement
		as = data.Nothing[annotations]()
	}
}

// Note, this is not an actual rule in the grammar, but a helper function to parse the imports. That
// said, if it was a rule, it would be the following:
//
//	```
//	imports = {{"\n"}, [annotations_], import} ;
//	```
func parseImports(p Parser) data.Either[data.Ers, data.List[importStatement]] {
	return parseImportsHelper(p, data.Nothing[annotations]())
}

// rule:
//
//	```
//	import = "import", {"\n"},
//		( package import
//		| "(", {"\n"}, package import, {{"\n"}, package import}, {"\n"}, ")"
//		) ;
func maybeParseImport(p Parser) (*data.Ers, data.Maybe[importing]) {
	importToken, found := getKeywordAtCurrent(p, token.Import)
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
	if token.Underscore.Match(p.current()) {
		p.advance()
		hiddenSelection := data.Nothing[data.NonEmpty[name]]()
		return data.Ok(hiddenSelection) // hides all names from imported namespace
	}

	type group struct{ data.NonEmpty[name] }
	lparen, found := getKeywordAtCurrent(p, token.LeftParen)
	es, symbols, ok := parseSepSequenced[group](p, IllegalEmptyUsingClause, token.Comma, parseMaybeName).Break()
	if !ok { // error occurred
		return data.PassErs[data.Maybe[data.NonEmpty[name]]](es)
	}

	if found {
		rparen, found := getKeywordAtCurrent(p, token.RightParen)
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
	if as, foundAs := getKeywordAtCurrent(p, token.As); foundAs {
		id, ok := parseLowerIdent(p).Break()
		if !ok {
			e := data.Nil[data.Err](1).Snoc(data.MkErr(ExpectedNamespaceAlias, p))
			return &e, data.Nothing[selections](p)
		}

		asClause := data.Inl[data.Maybe[data.NonEmpty[name]]](id).Update(as)
		return nil, data.Just(asClause)
	}

	using, foundUsing := getKeywordAtCurrent(p, token.Using)
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
	pathLiteral, found := getKeywordAtCurrent(p, token.ImportPath)
	if !found {
		return nil, data.Nothing[packageImport](p)
	}

	path := data.EOne[importPathIdent](pathLiteral)

	es, selections := maybeParseImportSpecification(p)
	if es != nil {
		return es, data.Nothing[packageImport](p)
	}

	pi := data.EMakePair[packageImport](path, selections)
	return nil, data.Just(pi)
}
