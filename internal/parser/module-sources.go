package parser

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util/code"
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

	src, err := code.FileSource(path)
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

func getImports(ys YewSource) []string {
	header, isHeader := ys.header.Break()
	if !isHeader {
		return []string{}
	}

	imports := header.Snd()
	if imports.IsEmpty() {
		return []string{}
	}

	out := []string{}

	for _, i := range imports.Elements() {
		importing := i.Snd()
		for _, pi := range importing.Elements() {
			out = append(out, pi.Fst().Break().String())
		}
	}

	return out
}

type moduleGraph struct {
	dependencies map[string]*moduleGraph
	module       YewModule // 0 value in a root, non-zero value in the leaves
}

func buildModuleGraph(mods map[string]YewModule) *moduleGraph {
	graph := new(moduleGraph)
	graph.dependencies = make(map[string]*moduleGraph, len(mods))

	for _, mod := range mods {
		graph.dependencies[mod.path] = func(mod YewModule) *moduleGraph {
			graph := new(moduleGraph)
			graph.module = mod
			graph.dependencies = make(map[string]*moduleGraph, len(mod.required))
			return graph
		}(mod)
	}

	for _, node := range graph.dependencies {
		for dep := range node.module.required {
			node.dependencies[dep] = graph.dependencies[dep]
		}
	}

	return graph
}

func detectCycle(buf map[string][]string, visited map[string]bool, visits map[string]map[string]struct{}, path string) (map[string][]string, bool) {
	// if the path has already been visited, then a cycle has been detected
	if visited[path] {
		return buf, true
	}

	// mark the path as visited
	visited[path] = true
	// un-mark the path as visited when the function exits
	defer func() { visited[path] = false }()

	// get the paths that the current path depends on
	dependsOn, ok := visits[path]
	if !ok {
		return buf, false
	}

	// check if the paths that the current path depends on have cycles
	var cycleFound bool
	for dependency := range dependsOn {
		var found bool
		buf, found = detectCycle(buf, visited, visits, dependency)
		if !found {
			continue
		} else {
			// keep processing even if a cycle is found
			cycleFound = true
		}

		// if there's a cycle, record the cycle
		if len(buf[dependency]) == 0 {
			buf[dependency] = make([]string, 0, 1)
		}
		buf[dependency] = append(buf[dependency], path)
	}

	return buf, cycleFound
}

type cycleError struct {
	path  string
	cycle []string
}

func (ce cycleError) Error() string {
	return "import cycle detected between " + ce.path + " and the following: " + strings.Join(ce.cycle, ", ")
}

func detectCycles(yms map[string]YewModule) []error {
	// initialize data
	visited := make(map[string]bool)
	type visitMap map[string]map[string]struct{}
	visits := make(visitMap)

	for _, ym := range yms {
		visited[ym.path] = false
		visits[ym.path] = ym.required
	}

	// detect cycles
	cycles := make(map[string][]string)
	var found bool
	for path := range visits {
		var cyclic bool
		cycles, cyclic = detectCycle(cycles, visited, visits, path)
		if cyclic {
			found = true
		}
	}

	if !found {
		return []error{} // no cycles found
	}

	// construct errors
	var es []error
	for path, cycle := range cycles {
		es = append(es, &cycleError{path, cycle})
	}

	return es
}
