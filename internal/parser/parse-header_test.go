//go:build test
// +build test

package parser

import (
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
)

func TestParseHeader(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[header]
	}{
		{
			"empty",
			[]api.Token{}, //
			data.Nothing[header](),
		},
		{
			"just module",
			[]api.Token{moduleTok, id_x_tok}, // module x
			data.Just(data.EMakePair[header](data.Just(module_x), data.Nil[importStatement]())),
		},
		{
			"just import",
			[]api.Token{importTok, importPathTok}, // import "a/b/c"
			data.Just(data.EMakePair[header](data.Nothing[module](), data.Nil[importStatement](1).Snoc(importStmtNode))),
		},
		{
			"module and import",
			[]api.Token{moduleTok, id_x_tok, newline, importTok, importPathTok}, // module x\nimport "a/b/c"
			data.Just(data.EMakePair[header](data.Just(module_x), data.Nil[importStatement](1).Snoc(importStmtNode))),
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseHeader, -1))
	}
}

func TestParseImports(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.List[importStatement]
	}{
		{
			"empty",
			[]api.Token{},
			data.Nil[importStatement](), //
		},
		{
			"single",
			[]api.Token{importTok, importPathTok},
			data.Nil[importStatement](1).Snoc(importStmtNode), // import "a/b/c"
		},
		{
			"annotated",
			[]api.Token{annotOpenTok, id_test_tok, rbracket, importTok, importPathTok},
			data.Nil[importStatement](1).Snoc(annotatedImportStmtNode), // [@test] import "a/b/c"
		},
		{
			"multiple",
			[]api.Token{importTok, importPathTok, newline, importTok, importPathTok},
			data.Nil[importStatement](2).Snoc(importStmtNode).Snoc(importStmtNode), // import "a/b/c"\nimport "a/b/c"
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseImports, -1))
	}
}

// rule:
//
//	```
//	import = "import", {"\n"},
//		( package import
//		| "(", {"\n"}, package import, {{"\n"}, package import}, {"\n"}, ")"
//		) ;
func TestMaybeParseImport(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[importing]
	}{
		{
			"empty",
			[]api.Token{},
			data.Nothing[importing](), //
		},
		{
			"single",
			[]api.Token{importTok, importPathTok},
			data.Just(imports_abc), // import "a/b/c"
		},
		{
			"import group single - 00",
			[]api.Token{importTok, lparen, importPathTok, rparen},
			data.Just(imports_abc), // import ("a/b/c")
		},
		{
			"import group single - 01",
			[]api.Token{importTok, lparen, importPathTok, newline, rparen},
			data.Just(imports_abc), // import ("a/b/c"\n)
		},
		{
			"import group single - 10",
			[]api.Token{importTok, lparen, newline, importPathTok, rparen},
			data.Just(imports_abc), // import (\n"a/b/c")
		},
		{
			"import group single - 11",
			[]api.Token{importTok, lparen, newline, importPathTok, newline, rparen},
			data.Just(imports_abc), // import (\n"a/b/c"\n)
		},
		{
			"import group multiple",
			[]api.Token{importTok, lparen, importPathTok, newline, importPathTok, rparen},
			data.Just(data.EConstruct[importing](pkgImport_abc, pkgImport_abc)), // import ("a/b/c"\n"a/b/c")
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, maybeParseImport, -1))
	}
}

// rule:
//
//	```
//	package import = import path, [{"\n"}, import specification] ;
//	```
func TestMaybeParsePackageImport(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[packageImport]
	}{
		{
			"empty",
			[]api.Token{},
			data.Nothing[packageImport](), //
		},
		{
			"single",
			[]api.Token{importPathTok},
			data.Just(pkgImport_abc), // "a/b/c"
		},
		{
			"as clause - 00",
			[]api.Token{importPathTok, as, id_x_tok},
			data.Just(data.EMakePair[packageImport](abc_path, data.Just(as_x))), // "a/b/c" as x
		},
		{
			"as clause - 01",
			[]api.Token{importPathTok, as, newline, id_x_tok},
			data.Just(data.EMakePair[packageImport](abc_path, data.Just(as_x))), // "a/b/c" as x
		},
		{
			"as clause - 10",
			[]api.Token{importPathTok, newline, as, id_x_tok},
			data.Just(data.EMakePair[packageImport](abc_path, data.Just(as_x))), // "a/b/c" as x
		},
		{
			"as clause - 11",
			[]api.Token{importPathTok, newline, as, newline, id_x_tok},
			data.Just(data.EMakePair[packageImport](abc_path, data.Just(as_x))), // "a/b/c" as x
		},
		{
			"using clause - 00",
			[]api.Token{importPathTok, using, underscoreTok},
			data.Just(data.EMakePair[packageImport](abc_path, data.Nothing[selections]())), // "a/b/c" using _
		},
		{
			"using clause - 01",
			[]api.Token{importPathTok, using, newline, underscoreTok},
			data.Just(data.EMakePair[packageImport](abc_path, data.Nothing[selections]())), // "a/b/c" using _
		},
		{
			"using clause - 10",
			[]api.Token{importPathTok, newline, using, underscoreTok},
			data.Just(data.EMakePair[packageImport](abc_path, data.Nothing[selections]())), // "a/b/c" using _
		},
		{
			"using clause - 11",
			[]api.Token{importPathTok, newline, using, newline, underscoreTok},
			data.Just(data.EMakePair[packageImport](abc_path, data.Nothing[selections]())), // "a/b/c" using _
		},
		{
			"symbol selection group",
			[]api.Token{importPathTok, using, lparen, id_x_tok, comma, id_x_tok, rparen},
			data.Just(data.EMakePair[packageImport](abc_path, data.Just(using_x_x))), // "a/b/c" using (x, x)
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, maybeParsePackageImport, -1))
	}
}

func TestMaybeParseModule(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[module]
	}{
		{
			"empty",
			[]api.Token{},
			data.Nothing[module](),
		},
		{
			"module - 0",
			[]api.Token{moduleTok, id_x_tok},
			data.Just[module](module_x),
		},
		{
			"module - 1",
			[]api.Token{moduleTok, newline, id_x_tok},
			data.Just[module](module_x),
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseModule, -1))
	}
}
