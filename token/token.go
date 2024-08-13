package token

import "fmt"

type Token struct {
	Value      string
	Type       Type
	Start, End int
}

func (token Token) Pos() (int, int) {
	return token.Start, token.End
}

func (a Token) Equals(b Token) bool {
	return a.Value == b.Value &&
		a.Type == b.Type &&
		a.Start == b.Start &&
		a.End == b.End
}

func (token Token) String() string {
	return token.Value
}

func (token Token) Debug() string {
	return fmt.Sprintf("Token{Value: \"%s\", Type: %v, Start: %d, End: %d}", token.Value, token.Type, token.Start, token.End)
}
