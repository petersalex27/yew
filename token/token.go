package token

type Token struct {
	Value string
	Type
	Start, End, Line int
}