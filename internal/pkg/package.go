package pkg

import "github.com/petersalex27/yew/internal/pkg/module"

func (pkg Data) Name() string {
	return pkg.name
}	

func (pkg Data) Lookup(name string) (*module.Data, bool) {
	m, ok := pkg.data.modules[name]
	return m, ok
}