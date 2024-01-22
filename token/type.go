package token

type Type uint

func (tokenType Type) Make(val string) Token {
	return Token{Value: val, Type: tokenType}
}