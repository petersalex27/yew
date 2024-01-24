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

func respondRepl(lex *lexer.Lexer, i, t, result int) (i_end, t_end, result_end int) {
	if result == 0 {
		return i, t, result
	}
	tokens, ok := lex.Tokenize()
	if !ok {
		messages := lex.FlushMessages()
		for _, m := range messages {
			fmt.Fprintf(os.Stderr, "%s\n", m.Error())
		}
		// remove erroneous tokens and source
		lex.Tokens = lex.Tokens[:t]
		lex.Source = lex.Source[:i]
		return i, t, result
	}

	print("<<")
	for _, tok := range tokens[t:] {
		t++
		print(" ", tok.String())
	}
	print("\n")
	i += result
	result = 0
	return i, t, result
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
	t := 0

	for result := 0; result >= 0; {
		result = promptRepl(lex)
		i, t, result = respondRepl(lex, i, t, result)
	}
}

func main() {
	repl()
}
