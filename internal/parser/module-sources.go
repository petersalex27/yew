package parser

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util"
)


type moduleSourceState struct {
	wg              *sync.WaitGroup
	modulesMu, esMu *sync.Mutex
	modules         *map[string][]api.Source
	es              *[]error
}

// locate a source of the module
func (msw moduleSourceState) getSource(path string) {
	defer msw.wg.Done()

	// check if the module source exists already (this is possible if the module is pointed to by a symbolic link)
	msw.modulesMu.Lock()
	if _, ok := (*msw.modules)[path]; ok {
		msw.modulesMu.Unlock()
		return
	}
	msw.modulesMu.Unlock()

	src, err := util.FileSource(path)
	if err != nil {
		msw.esMu.Lock()
		*msw.es = append(*msw.es, err)
		msw.esMu.Unlock()
		return
	}

	msw.modulesMu.Lock()
	(*msw.modules)[path] = append((*msw.modules)[path], src)
	msw.modulesMu.Unlock()
}

func (msw *moduleSourceState) locateModuleSources(root, path string) {
	defer msw.wg.Done()

	m, ers := getModuleSources(root, path)
	if len(ers) > 0 {
		msw.esMu.Lock()
		*msw.es = append(*msw.es, ers...)
		msw.esMu.Unlock()
		return
	}

	msw.modulesMu.Lock()
	for k, v := range m {
		(*msw.modules)[k] = v
	}
	msw.modulesMu.Unlock()
}

func resolveFile(root, dir string, f fs.DirEntry) (string, bool) {
	// check if the file is a symbolic link
	info, err := f.Info()
	if err != nil {
		return "", false
	}

	path := filepath.Join(dir, f.Name())

	if info.Mode()&os.ModeSymlink != 0 {
		// read the link
		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			return "", false
		}
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return "", false
	}

	// make sure the path is within the root directory
	return path, strings.HasPrefix(path, root)
}

// `root` and `dir` MUST be an absolute path!
func getModuleSources(root, dir string) (map[string][]api.Source, []error) {
	var wg sync.WaitGroup
	var modulesMu, esMu sync.Mutex
	es := []error{}

	fs, err := os.ReadDir(dir)
	if err != nil {
		return nil, []error{err}
	}

	modules := make(map[string][]api.Source)
	msw := moduleSourceState{&wg, &modulesMu, &esMu, &modules, &es}

	for _, f := range fs {
		var ok bool
		path, ok := resolveFile(root, dir, f)
		if !ok {
			continue
		}

		if !f.IsDir() {
			if filepath.Ext(path) != ".yew" {
				continue
			}
			
			wg.Add(1)
			go msw.getSource(path)
		} else {
			wg.Add(1)
			go msw.locateModuleSources(root, path)
		}
	}

	wg.Wait()

	if len(es) > 0 {
		return nil, es
	}
	return modules, nil
}