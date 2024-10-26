package parser

import (
	"errors"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
	yewErrors "github.com/petersalex27/yew/internal/errors"
	"github.com/petersalex27/yew/internal/symbol"
)

type modLookup func(api.Token) (error, symbol.Symbol)
type importedSymbols map[string]modLookup //map[string]symbol.Symbol

func (ys YewSource) Header() data.Maybe[header] {
	return ys.header
}

func (ys YewSource) Imports() data.Maybe[data.List[importStatement]] {
	return data.MapMaybe((header).Snd)(ys.header)
}

func (ys YewSource) Module() data.Maybe[module] {
	return data.Bind(ys.header, (header).Fst)
}

func (ys YewSource) Body() data.Maybe[body] {
	return ys.body
}

func (ys YewSource) Footer() data.Maybe[annotations] {
	return ys.footer
}

func (ys YewSource) GetPos() api.Position {
	return ys.pos
}

func (ys YewSource) Pos() (int, int) {
	return ys.pos.Pos()
}

func annotateImports(yps *YewPackageSource, annots annotations, ims importing) (importing, []error) {
	// TODO: implement
	return ims, []error{errors.New("import annotations not yet implemented")}
}

func createLookup(lookupTable map[string]symbol.Symbol, srcCode api.SourceCode, path string) modLookup {
	return func(tok api.Token) (error, symbol.Symbol) {
		sym, found := lookupTable[tok.String()]
		if !found {
			start, end := tok.Pos()
			return yewErrors.Name(srcCode, "undefined "+path+"."+tok.String(), start, end), sym
		}
		return nil, sym
	}
}

func getAccessor(yps *YewPackageSource, srcCode api.SourceCode, im packageImport) (string, modLookup, []error) {
	path := im.Fst().Break().String()
	path = yps.fromRoot(path)
	lookupTable, found := yps.dependencies[path]
	if !found {
		return "", nil, []error{unknownModuleError(srcCode, path, im)}
	}
	return path, createLookup(lookupTable, srcCode, path), []error{}
}

func (ys YewSource) resolveImports(yps *YewPackageSource, ims data.NonEmpty[packageImport]) (YewSource, []error) {
	ys.imported = make(importedSymbols)

	errors := []error{}
	for _, im := range ims.Elements() {
		acc, f, es := getAccessor(yps, ys.srcCode, im)
		if len(es) > 0 {
			errors = append(errors, es...)
		}
		ys.imported[acc] = f
	}

	return ys, errors
}

func foldFunc_(ne data.NonEmpty[packageImport], im importing) data.NonEmpty[packageImport] {
	return ne.Append(im.Elements()...)
}

func (ys YewSource) processHeader(yps *YewPackageSource, modPath string) (_ YewSource, warnings, errors []error) {
	// process the header
	header, isHeader := ys.header.Break()
	var imports data.List[importStatement]
	if isHeader {
		imports = header.Snd()
	} else {
		imports = data.Nil[importStatement](0)
	}

	// process imports
	annotImports := data.ListMap(func(is importStatement) importing {
		annots, ims := is.Split()
		var ws []error
		if annots, isAnnots := annots.Break(); isAnnots {
			ims, ws = annotateImports(yps, annots, ims)
			warnings = append(warnings, ws...)
		}
		return ims
	})

	if len(errors) > 0 {
		return ys, warnings, errors
	}

	ims := annotImports(imports)
	if ims.Len() == 0 {
		tail, _ := ims.Tail().Break()
		ys, errors = ys.resolveImports(yps, data.FoldLeft(foldFunc_, ims.At(0).NonEmpty)(tail))
	}

	return ys, warnings, errors
}

func (ys YewSource) processFooter(yps *YewPackageSource, modPath string) (_ YewSource, warnings, errors []error) {
	// process the footer
	panic("not implemented")
}

func (ys YewSource) processBody(yps *YewPackageSource, modPath string) (_ YewSource, warnings, errors []error) {
	// process the body
	panic("not implemented")
}

func (ys YewSource) process(yps *YewPackageSource, modPath string) (_ YewSource, warnings, errors []error) {
	if ys, warnings, errors = ys.processHeader(yps, modPath); len(errors) > 0 {
		return ys, warnings, errors
	} else if ys, warnings, errors = ys.processBody(yps, modPath); len(errors) > 0 {
		return ys, warnings, errors
	} else if ys, warnings, errors = ys.processFooter(yps, modPath); len(errors) > 0 {
		return ys, warnings, errors
	}

	return ys, warnings, errors
}
