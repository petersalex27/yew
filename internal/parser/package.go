package parser

import (
	// "path/filepath"
	// "sync"

	// "github.com/petersalex27/yew/api/util"
	// "github.com/petersalex27/yew/internal/lexer"
	"github.com/petersalex27/yew/internal/pkg"
	//mod "github.com/petersalex27/yew/internal/pkg/module"
)

// package file nodes
type (
	YewPackageSource struct {
		name      string
		requiring []requirement
	}

	requirement struct {
		path    string
		version version
	}

	version struct{ major, minor, patch uint }

	pkgTokenType byte

	pkgToken struct {
		ty         pkgTokenType
		val        string
		start, end int
	}
)

const (
	InvalidRequirePath = "invalid require path"
	InvalidVersion     = "invalid required version"
	CouldNotLocate     = "could not locate required package"
)

// builds a package given its directory
func Build(dir string) (pkg.Data, []error) {
	panic("not implemented")
	// path := filepath.Join(dir, "yew.package")
	// src, err := util.FileSource(path)
	// if err != nil {
	// 	return pkg.Data{}, []error{err}
	// }

	// // parse the package file
	// p := Init(lexer.InitPackageLexer(src))
	// yps, es := parseYewPackageSource(p)
	// if len(es) > 0 {
	// 	return pkg.Data{}, es
	// }

	// reqs := make([]mod.Data, len(yps.requiring))
	// var wg sync.WaitGroup

	// // locate required resources
	// //
	// // each resource has its own index, so we can safely assign the result there from the 
	// // corresponding goroutine
	// for i, req := range yps.requiring {
	// 	// validate the path
	// 	if e := req.valid(); e != nil {
	// 		es = append(es, e)
	// 		continue // keep going to provide more errors if needed
	// 	}

	// 	wg.Add(1)
	// 	// run goroutine to locate the resource
	// 	go func() {
	// 		defer wg.Done()
			
	// 	}()
	// }

	// wg.Wait()

	// // TODO
	// panic("not implemented")
}

// simple top-down parser for the package file
// func parseYewPackageSource(p parser) (YewPackageSource, []error) {

// }