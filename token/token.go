package token

import "fmt"

type Token struct {
	Value string
	Type
	Start, End, Line int
}

func (a Token) Equals(b Token) bool {
	return a.Value == b.Value && 
		a.Type == b.Type &&
		a.Start == b.Start &&
		a.End == b.End &&
		a.Line == b.Line
}

func (token Token) String() string {
	return fmt.Sprintf("Token{Value: \"%s\", Type: %v, Start: %d, End: %d, Line: %d}", token.Value, token.Type, token.Start, token.End, token.Line)
}