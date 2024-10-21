package pkg

import (
	"github.com/petersalex27/yew/internal/pkg/module"
)

type (
	Data struct {
		name string
		data packageData
	}

	// references a loadable package (if that packages exists)
	//
	// this is used so a library can have references to its packages w/o
	// actually loading them until they are needed
	Loader struct{ path string }

	// packageData represents a compiled package ready to be imported
	//
	// It contains llvm IR and a symbol table for Yew code to reference
	packageData struct {
		modules map[string]*module.Data
	}
)

func (loader Loader) Load() (packageData, error) {
	return loader.mockLoad()
}

func (loader Loader) mockLoad() (packageData, error) {
	pkgData := packageData{make(map[string]*module.Data)}
	pkgData.modules["mock"] = module.Mock()
	return pkgData, nil
}
