package parser

// import (
// 	"github.com/petersalex27/yew/api"
// )

// type (
// 	cmd struct{ api.Node }

// 	tokenStream struct {
// 		tokens []api.Token
// 	}
// )

// func (cmd *cmd) Describe() (name string, children []api.Node) { return "cmd", []api.Node{cmd.Node} }

// func (cmd *exposeCmd) Describe() (name string, children []api.Node) {
// 	return "exposeCmd", append([]api.Node{cmd.exposing}, cmd.tokens...)
// }

// func (cmd *exposeCmd) Pos() (int, int) {
// 	start, end := cmd.exposing.Pos()
// 	if len(cmd.tokens) == 0 {
// 		return start, end
// 	}

// 	_, end = cmd.tokens[len(cmd.tokens)-1].Pos()
// 	return start, end
// }

// func (cmd *cmd) acceptCommandToken(p *ParserState) {
// 	panic("not implemented")
// }
