package parser

import (
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/lexer"
)

// package file nodes
type (
	YewPackageSource struct {
		name lowerIdent
		lang language
		body data.List[packageBodyElement]
		mods data.List[YewModule]
		api.Position
	}

	packageBodyElement interface {
		api.Node
		packageBodyElement() packageBodyElement
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

	yps, es := parseYewPackageSource(p, modSources)
	if len(es) > 0 {
		return yps, es
	}

	return yps, nil
}

func BuildPackage(packageRootDir string) (YewPackageSource, []error) {
	var err error
	packageRootDir, err = filepath.Abs(packageRootDir)
	if err != nil {
		return YewPackageSource{}, []error{err}
	}

	yewDotPackage := filepath.Join(packageRootDir, "yew.package")
	src, err := util.FileSource(yewDotPackage)
	if err != nil {
		return YewPackageSource{}, []error{err}
	}

	yps, errs := parsePackage(lexer.InitPackageLexer(src), packageRootDir)
	if len(errs) > 0 {
		return yps, errs
	}

	return yps, nil
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

func validateModuleName(p parser, mHeader data.Maybe[header], expected string) error {
	actual := getModuleName(mHeader)
	if actual != expected {
		e := data.MkErr(ModuleNameErrorMsg(expected), mHeader)
		return parseError(p, e)
	}

	return nil
}

// parses sources for a module
func parseModuleSources(sources []api.Source) (YewModule, []error) {
	var wg sync.WaitGroup
	var modMu sync.Mutex
	var esMu sync.Mutex

	es := []error{}
	mod := YewModule{
		name: "", // TODO
		List: data.Nil[YewSource](),
	}

	first := true
	for _, src := range sources {
		wg.Add(1)

		go func(src api.Source) {
			defer wg.Done()

			p := parseYewSource(Init(lexer.Init(src)))

			pes := GetErrors(p)
			if len(pes) > 0 {
				esMu.Lock()
				es = append(es, pes...)
				esMu.Unlock()
			} else {
				modMu.Lock()

				ast := p.(*ParserState).ast
				ok := true
				if first {
					first = false
					mod.name = getModuleName(ast.header)
				} else if err := validateModuleName(p, ast.header, mod.name); err != nil {
					ok = false
					esMu.Lock()
					es = append(es, err)
					esMu.Unlock()
				}

				if ok {
					mod.List = mod.List.Snoc(ast)
				}

				modMu.Unlock()
			}
		}(src)
	}

	wg.Wait()

	if len(es) > 0 {
		return mod, es
	}
	return mod, nil
}

// simple top-down parser for the package file
func parseYewPackageSource(p parser, modules map[string][]api.Source) (YewPackageSource, []error) {
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

	for _, sources := range modules {
		wg.Add(1)
		go func(sources []api.Source) {
			defer wg.Done()

			m, es := parseModuleSources(sources)
			if len(es) > 0 {
				ersMu.Lock()
				ers = append(ers, es...)
				ersMu.Unlock()
				return
			}

			modsMu.Lock()
			mods = mods.Snoc(m)
			modsMu.Unlock()
		}(sources)
	}

	wg.Wait()

	if len(ers) > 0 {
		return YewPackageSource{}, ers
	}

	yps := makeYewPackageSource(pkgName, lang, pkgBody, mods)
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
