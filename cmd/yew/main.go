package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/petersalex27/yew/cmd/yew/repl"
)

var (
	interactive = flag.Bool("i", false, "-i")
	file        = flag.String("file", "", "-file <path>")
)

func init() {
	flag.Usage = func() {
		// TODO
		fmt.Fprintf(os.Stderr, "usage: yew [-i] [-file <path>]\n")
		flag.PrintDefaults()
	}
}

func compileFile(path string) {
	// lex := lexer.Init(source.FilePath(path))
	// lex.Write()
	// /*tokens*/ _, ok := lex.Tokenize()
	// if !ok {
	// 	msgs := lex.FlushMessages()
	// 	errors.PrintErrors(msgs)

	// 	os.Exit(1)
	// }

	// p := util.InitParser[*parser.Parser](lex.SourceCode)
	// _, _ = p.Begin() // TODO
}

func main() {
	// if true { // TODO: remove
	// 	repl()
	// }

	flag.Parse()
	fmt.Printf("%t, %s\n", *interactive, *file)
	if *interactive {
		repl.Run()
	} else if *file != "" {
		compileFile(*file)
		os.Exit(0)
	} else {
		flag.Usage()
	}
}
