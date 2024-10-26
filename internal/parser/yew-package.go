package parser

import (
	"path/filepath"

	"github.com/petersalex27/yew/common/data"
)

func (reqs requirements) process(yps *YewPackageSource) []error {
	// TODO: implement
	return []error{}
}

func getNextModule(graph *moduleGraph) (path string, mod YewModule) {
	// search for a terminal node, there must always be at least one since there are no cycles
	// and this is a directed graph
	for p, node := range graph.dependencies {
		if len(node.dependencies) == 0 {
			path = p
			mod = node.module
			break
		}
	}
	
	for _, node := range graph.dependencies {
		delete(node.dependencies, path)
	}
	delete(graph.dependencies, path) // remove the node from the graph

	return path, mod
}

func processModule(yps *YewPackageSource, mod YewModule) (warnings, errors []error) {
	warnings, errors = []error{}, []error{}
	mod.List = data.ListMap(func(ys YewSource) YewSource {
		var ws, es []error
		ys, ws, es = ys.process(yps, mod.path)
		warnings = append(warnings, ws...)
		errors = append(errors, es...)
		return ys
	})(mod.List)

	yps.mods[mod.path] = mod
	return warnings, errors
}

func processModules(yps *YewPackageSource) (warnings, errors []error) {
	// build local dependency graph. It only has a depth of 1
	graph := buildModuleGraph(yps.mods)
	warnings, errors = []error{}, []error{}

	// traverse the graph and evaluate a terminal module node
	path, mod := getNextModule(graph)
	nextMod := mod
	for path != "" {
		mod = nextMod
		// get next module and modify the graph accordingly
		go func() { path, nextMod = getNextModule(graph) }()

		// process the terminal module node
		var ws, es []error
		ws, es = processModule(yps, mod)
		warnings = append(warnings, ws...)
		errors = append(errors, es...)
	}

	return warnings, errors
}

func (yps *YewPackageSource) process() (warnings, errors []error) {
	warnings, errors = []error{}, []error{}

	for _, elem := range yps.body.Elements() {
		es := elem.process(yps)
		errors = append(errors, es...)
	}

	if len(errors) > 0 {
		return warnings, errors
	}

	// check for import cycles among modules of the package
	if errors = detectCycles(yps.mods); len(errors) > 0 {
		return warnings, errors
	}

	warnings, errors = processModules(yps)
	return warnings, errors
}

func (yps YewPackageSource) fromRoot(path string) string {
	return filepath.Join(yps.path, path)
}