package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/petersalex27/yew/lexer"
)

func promptRepl(lex *lexer.Lexer) int {
	fmt.Print(">> ")
	return lex.Write()
}

func respondRepl(lex *lexer.Lexer, i, result int) (i_end, result_end int) {
	for result > 0 {
		print("<< ", lex.Source[i])
		i++
		result--
	}
	return i, result
}

func repl() {
	// print initial message
	fmt.Printf("Yew (interactive) v0.0.1\nUse ctrl+C to exit\n\n")

	// initialize quit signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println("\nexiting...")
		os.Exit(0)
	}()

	lex := lexer.Init(lexer.StdinSpec)
	i := 0

	for result := 0; result >= 0; {
		result = promptRepl(lex)
		i, result = respondRepl(lex, i, result)
	}
}

func main() {
	repl()
}
