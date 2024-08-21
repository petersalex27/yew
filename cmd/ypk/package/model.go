package pkg

import (
	"encoding"
	"os"

	"github.com/petersalex27/ypk/util"
)

type Module struct {
	Name string `pkg:"name"`
	// symbol table:
	//		"public_symbols": [[IDENT, TYPE], ...]
	PublicSymbols [][2]string `pkg:"public_symbols"`
}

type Package struct {
	Name string `pkg:"name"`
	Version []int `pkg:"version"`
	// source of the package
	Source string `pkg:"source"`
	// url of packages this package depends on
	Dependencies []string `pkg:"dependencies"`
	// symbol table:
	//		"public_symbols": [[IDENT, TYPE], ...]
	PublicSymbols [][2]string `pkg:"public_symbols"`
	// modules in the package
	Modules []Module `pkg:"modules"`
}

func Pack(path string) (p Package, err error) {
	f, err := os.Open(path)
	if err != nil {
		return p, err
	}

	defer f.Close()

	lexer := util.NewLexer(f)
	lexer.Scan()
	// get the package name
	p.Name = path
	// get the package version
	p.Version = util.GetLatest(path)
	// get the package source
	p.Source = path
	// get the package dependencies
	p.Dependencies = []string{}
	// get the package public symbols
	p.PublicSymbols = [][2]string{}
	// get the package modules
	p.Modules = []Module{}

	return p, nil
}

func (p *Package) MarshalBinary() ([]byte, error) {

}

func (p *Package) UnmarshalBinary(data []byte) error {

	encoding.BinaryUnmarshaler
}