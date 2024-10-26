package parser

import (
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util/code"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/lexer"
	"github.com/petersalex27/yew/internal/symbol"
)

// package file nodes
type (
	YewPackageSource struct {
		name lowerIdent
		lang language
		body data.List[packageBodyElement]
		api.Position

		// non-syntax related data

		mods map[string]YewModule
		// all the symbols required by the package
		dependencies map[string]map[string]symbol.Symbol
		// each module mapped to its symbols
		exports map[string]map[string]symbol.Symbol
		// directory that yew.package file resides in
		path string
	}

	packageBodyElement interface {
		api.Node
		packageBodyElement() packageBodyElement
		process(*YewPackageSource) []error
	}

	requirement = data.Pair[api.Token, version]

	requirements struct{ data.List[requirement] }

	language struct{ data.Pair[lowerIdent, version] }

	versionCanonical struct {
		major, minor, patch uint
		milestone           string
		milestoneNum        uint
		api.Position
	}

	versionRelative struct{ data.Solo[api.Token] }

	version = data.Either[versionCanonical, versionRelative]
)

func parsePackage(scanner api.ScannerPlus, packageRootDir string) (YewPackageSource, []error) {
	modSources, modSrcErs := getModuleSources(packageRootDir, packageRootDir)
	if len(modSrcErs) > 0 {
		return YewPackageSource{}, modSrcErs
	}

	p := Init(scanner)

	yps, es := parseYewPackageSource(packageRootDir, p, modSources)
	if len(es) > 0 {
		return yps, es
	}

	return yps, nil
}

func BuildPackage(packageRootDir string) (yps YewPackageSource, warnings []error, errors []error) {
	// get the absolute path of the package root directory
	var err error
	packageRootDir, err = filepath.Abs(packageRootDir)
	if err != nil {
		return YewPackageSource{}, []error{}, []error{err}
	}

	// get the source for the package file
	yewDotPackage := filepath.Join(packageRootDir, "yew.package")
	src, err := code.FileSource(yewDotPackage)
	if err != nil {
		return YewPackageSource{}, []error{}, []error{err}
	}

	warnings = []error{}

	// parse the package file and its modules
	yps, errors = parsePackage(lexer.InitPackageLexer(src), packageRootDir)
	if len(errors) > 0 {
		return yps, warnings, errors
	}

	// process the package source
	if errors, warnings = yps.process(); len(errors) > 0 {
		return yps, warnings, errors
	}

	return yps, warnings, nil
}

func getModuleName(mHeader data.Maybe[header]) string {
	h, isH := mHeader.Break()
	if !isH {
		return ""
	}

	modDirective, isMD := h.Fst().Break()
	if !isMD {
		return ""
	}

	return modDirective.name.Break().Break().String()
}

func runParseYewSource(src api.Source, nm nameMap) (ys YewSource, reqs []string, es []error) {
	p := parseYewSource(Init(lexer.Init(src)))
	es = GetErrors(p)

	if len(es) != 0 {
		return YewSource{}, reqs, es
	}

	ast := p.(*ParserState).ast
	name := getModuleName(ast.header)
	nm.addSource(name, ast.header, p.SrcCode())
	reqs = getImports(ast)

	return ast, reqs, nil
}

type sourceState struct {
	header data.Maybe[header]
	api.SourceCode
}

type nameMap struct {
	mu   *sync.Mutex
	data *map[string][]sourceState
}

func newNameMap() nameMap {
	// for correct modules, you'd need at MOST 1
	// 	- so, optimize for the correct case?
	m := make(map[string][]sourceState, 1)
	return nameMap{&sync.Mutex{}, &m}
}

func (ss sourceState) SrcCode() api.SourceCode {
	return ss.SourceCode
}

func (nm nameMap) addSource(modName string, mh data.Maybe[header], src api.SourceCode) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if _, ok := (*nm.data)[modName]; !ok {
		(*nm.data)[modName] = make([]sourceState, 0)
	}

	(*nm.data)[modName] = append((*nm.data)[modName], sourceState{mh, src})
}

// returns an empty list if no errors--so safe to append to
func checkModuleNames(nm nameMap) []error {
	dataLength := len(*nm.data)
	if dataLength <= 1 {
		return []error{}
	}

	es := []error{}

	for name, sources := range *nm.data {
		msg := UnexpectedModuleName
		if name == "" {
			msg = ExpectedModuleName
		}

		for _, src := range sources {
			e := data.MkErr(msg, src.header)
			es = append(es, parseError(src, e))
		}
	}

	return es
}

// parses sources for a module
func parseModuleSources(dir string, sources []api.Source) (YewModule, []error) {
	var wg sync.WaitGroup
	var modMu sync.Mutex
	var esMu sync.Mutex
	nameMap := newNameMap()

	es := []error{}
	mod := YewModule{
		name:     "", // TODO
		path:     dir,
		List:     data.Nil[YewSource](),
		required: make(map[string]struct{}),
	}

	for _, src := range sources {
		wg.Add(1)

		go func(src api.Source) {
			defer wg.Done()

			ys, reqs, ysEs := runParseYewSource(src, nameMap)
			if len(ysEs) != 0 {
				esMu.Lock()
				es = append(es, ysEs...)
				esMu.Unlock()
				return
			}

			modMu.Lock()
			mod.List = mod.List.Snoc(ys)
			for _, r := range reqs {
				mod.required[r] = struct{}{}
			}
			modMu.Unlock()
		}(src)
	}

	wg.Wait()

	// check module names, appending any errors to the list
	es = append(checkModuleNames(nameMap), es...)
	return mod, es
}

// simple top-down parser for the package file
func parseYewPackageSource(path string, p parser, modules map[string][]api.Source) (YewPackageSource, []error) {
	// parse the package name
	esName, pkgName, isPkgName := parsePackageName(p).Break()
	if !isPkgName {
		return YewPackageSource{}, parseErrors(p, esName)
	} else if !then(p) { // sequence language
		es := []error{parseError(p, data.MkErr(ExpectedPackageLanguage, p))}
		return YewPackageSource{}, es
	}

	// parse the language
	esLang, lang, isLang := parsePackageLanguage(p).Break()
	if !isLang {
		return YewPackageSource{}, parseErrors(p, esLang)
	}

	// parse package body
	esBody, pkgBody, isPkgBody := parsePackageBody(p).Break()
	if !isPkgBody {
		return YewPackageSource{}, parseErrors(p, esBody)
	}

	// parse the modules
	var wg sync.WaitGroup
	var ersMu, modsMu sync.Mutex

	var ers []error
	mods := data.Nil[YewModule](len(modules))

	for dir, sources := range modules {
		wg.Add(1)
		go func(dir string, sources []api.Source) {
			defer wg.Done()

			m, es := parseModuleSources(dir, sources)
			if len(es) > 0 {
				ersMu.Lock()
				ers = append(ers, es...)
				ersMu.Unlock()
				return
			}

			modsMu.Lock()
			mods = mods.Snoc(m)
			modsMu.Unlock()
		}(dir, sources)
	}

	wg.Wait()

	if len(ers) > 0 {
		return YewPackageSource{}, ers
	}

	yps := makeYewPackageSource("", pkgName, lang, pkgBody, mods)
	return yps, nil
}

func parsePackageName(p parser) data.Result[lowerIdent] {
	pkgTok, found := getKeywordAtCurrent(p, token.Package, dropBeforeAndAfter)
	if !found {
		return data.Fail[lowerIdent](ExpectedPackage, p)
	}

	name, isName := parseLowerIdent(p).Break()
	if !isName {
		return data.Fail[lowerIdent](ExpectedPackageName, p)
	}

	name.Position = name.Position.Update(pkgTok)
	return data.Ok(name)
}

func parsePackageLanguage(p parser) data.Result[language] {
	langName, isLang := parseLowerIdent(p).Break()
	if !isLang {
		return data.Fail[language](ExpectedPackageLanguage, p)
	}

	es, ver, isVer := parseVersion(p).Break()
	if !isVer {
		return data.PassErs[language](es)
	}

	return data.Ok(data.EMakePair[language](langName, ver))
}

func parseRequirementBody(p parser) (*data.Ers, requirements) {
	type requirements2 struct{ data.NonEmpty[requirement] }
	reqs := parseGroup[requirements2](p, ExpectedPackageRequirement, maybeParseRequirement)
	if es, reqs, isReqs := reqs.Break(); !isReqs {
		return &es, requirements{}
	} else {
		return nil, requirements{reqs.Weaken()}
	}
}

func parseMaybePackageBodyElement(p parser) (*data.Ers, data.Maybe[packageBodyElement]) {
	if reqTok, found := getKeywordAtCurrent(p, token.Require, dropAfter); found {
		es, req := parseRequirementBody(p)
		if es != nil {
			return es, data.Nothing[packageBodyElement]()
		}
		req.Position = req.Position.Update(reqTok)
		return nil, data.Just[packageBodyElement](req)
	}

	return nil, data.Nothing[packageBodyElement]()
}

func parsePackageBody(p parser) data.Result[data.List[packageBodyElement]] {
	body := data.Nil[packageBodyElement]()
	origin := getOrigin(p)
	for then(p) {
		es, elem := parseMaybePackageBodyElement(p)
		if es != nil {
			return data.PassErs[data.List[packageBodyElement]](*es)
		} else if el, isEl := elem.Break(); !isEl {
			setOrigin(p, origin)
			break
		} else {
			origin = getOrigin(p)
			body = body.Snoc(el)
		}
	}

	return data.Ok(body)
}

func parseVersion(p parser) data.Result[version] {
	versionTok, found := getKeywordAtCurrent(p, token.Version, dropNone)
	if !found {
		return data.Fail[version](ExpectedVersion, p)
	}

	major_minor_patch := strings.Split(strings.TrimPrefix(versionTok.String(), "v"), ".")
	if len(major_minor_patch) == 4 {
		return getTaggedCanonicalVersion(versionTok, major_minor_patch)
	} else if len(major_minor_patch) == 1 {
		return getRelativeVersion(versionTok, major_minor_patch)
	} else if len(major_minor_patch) != 3 {
		return data.Fail[version](InvalidVersion, versionTok)
	} // else, untagged canonical version

	return getCanonicalVersion(versionTok, major_minor_patch)
}

func getRelativeVersion(versionTok api.Token, ver []string) data.Result[version] {
	// two options: @stable or @latest
	if len(ver) != 1 {
		panic("illegal argument")
	}

	if ver[0] == "@stable" {
		return data.Ok(data.Inr[versionCanonical](data.EOne[versionRelative](versionTok)))
	} else if ver[0] == "@latest" {
		return data.Ok(data.Inr[versionCanonical](data.EOne[versionRelative](versionTok)))
	} else {
		return data.Fail[version](UnknownRelativeVersionTag, versionTok)
	}
}

func getTaggedCanonicalVersion(versionTok api.Token, major_minor_patch_milestone []string) data.Result[version] {
	if len(major_minor_patch_milestone) != 4 {
		panic("illegal argument")
	}

	// get the milestone tag from the patch
	var milestoneTag string
	var found bool
	major_minor_patch_milestone[2], milestoneTag, found = strings.Cut(major_minor_patch_milestone[2], "@")
	if !found {
		return data.Fail[version](InvalidVersion, versionTok)
	}

	// parse as canonical version
	res := getCanonicalVersion(versionTok, major_minor_patch_milestone[:3])
	es, v_, isV := res.Break()
	if !isV {
		return data.PassErs[version](es)
	}

	v, _, _ := v_.Break()
	v.milestone = milestoneTag

	var err error
	bitWidth := 64
	if math.MaxUint == math.MaxUint32 {
		bitWidth = 32
	}
	val, err := strconv.ParseUint(major_minor_patch_milestone[3], 10, bitWidth)
	if err != nil {
		return data.Fail[version](InvalidVersion, versionTok)
	}

	v.milestoneNum = uint(val)
	return data.Ok(data.Inl[versionRelative](v))
}

func getCanonicalVersion(versionTok api.Token, major_minor_patch []string) data.Result[version] {
	if len(major_minor_patch) != 3 {
		panic("illegal argument")
	}

	vs := make([]uint, 3)
	for i, v := range major_minor_patch {
		if n, err := strconv.Atoi(v); err != nil {
			return data.Fail[version](InvalidVersion, versionTok)
		} else {
			vs[i] = uint(n)
		}
	}

	// create the version
	v := versionCanonical{vs[0], vs[1], vs[2], "", 0, versionTok.GetPos()}
	return data.Ok(data.Inl[versionRelative](v))
}

func maybeParseRequirement(p parser) (*data.Ers, data.Maybe[requirement]) {
	path, isPath := getKeywordAtCurrent(p, token.StringValue, dropNone)
	if !isPath {
		return nil, data.Nothing[requirement](p)
	}

	if es, version, isVersion := parseVersion(p).Break(); !isVersion {
		return &es, data.Nothing[requirement](p)
	} else {
		req := data.MakePair(path, version)
		return nil, data.Just(req)
	}
}
