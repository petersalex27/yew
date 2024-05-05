// =================================================================================================
// Alex Peters - February 01, 2024
//
// Yew module files
// =================================================================================================
package module

import (
	"bytes"
	"io"
)

type Yewmod struct {
	// paths to files in module
	paths []string
	// required module paths (http), (name, path)
	remoteReqs map[string]string
	// required module paths (local), (name, path)
	localReqs map[string]string
	// compiled locations of deps
	compiled map[string]string
}

type YewmodReader struct {
	*bytes.Reader
	raw []byte
}

// opens a module, creating a reader for the given the module directory `dir`.
//
// `err` != nil if and only if there was an error
//
// `err` is of type `error.ErrorMessage`
func Open(dir string) (reader io.ReaderAt, err error) {
	//os.DirFS()
	return // TODO
}

func Read() *Yewmod {
	return nil // TODO
}