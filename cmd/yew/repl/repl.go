package repl

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/petersalex27/yew/api"
	// "github.com/petersalex27/yew/api/util"
	// "github.com/petersalex27/yew/internal/lexer"
	// "github.com/petersalex27/yew/internal/parser"
)

func prompt() { fmt.Print("yew< ") }

func respond(output *os.File, resp string) { fmt.Fprintf(output, "yew> %s\n", resp) }

func throw(err error) { respond(os.Stderr, err.Error()) }

func read(input *bufio.Reader) string {
	text, _ := input.ReadString('\n')
	return text
}

func promptRepl(scanner api.ScannerPlus, input *bufio.Reader) {
	prompt()
	scanner.AppendSource(read(input))
}

// return value will be recognized by the regex
//
//	`(v[0-9]+(\.[0-9]+)*)|^$`
//
// NOTE: prepends space if not empty
func version() string {
	if version, err := os.ReadFile("./version.txt"); err == nil {
		return " v" + string(version)
	}
	return ""
}

func reportErrors(errors []error) {
	for _, error := range errors {
		throw(error)
	}
}

func Run() {
	// print initial message
	fmt.Printf("Yew (interactive)" + version() + "\nUse ctrl+C to exit\n\n")

	// initialize quit signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		switch <-sigs {
		case syscall.SIGINT:
			fmt.Println("\nctrl+C detected...")
		case syscall.SIGTERM:
			fmt.Println("\nexiting...")
			os.Exit(0)
		}
	}()

	// input := bufio.NewReader(os.Stdin)

	// prompt()
	// lex := lexer.Init(util.FreeSource("<stdin>", read(input)))
	// p := parser.Init(lex)
	
	// for {
	// 	es := []error{}
	// 	switch lex.Command() {
	// 	case "":
	// 		if !p.ReplParse() {
	// 			es = p.Errors()
	// 		}
	// 		// TODO: evaluate ast
	// 	case ":import", ":type", ":instances", ":main", ":run", ":set":
	// 		panic("command not implemented")
	// 	case ":expose":
	// 		es = command.Expose(lex)
	// 	case ":help":
	// 		es = command.Help(lex)
	// 	case ":quit":
	// 		lex.Stop()
	// 		sigs <- syscall.SIGTERM // exit
	// 	}

	// 	if es != nil {
	// 		reportErrors(es)
	// 		lex.Restore()
	// 	}
	// 	promptRepl(lex, input)
	// }
}
