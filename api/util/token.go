package util

import (
	"fmt"

	"github.com/petersalex27/yew/api"
)

// ExposeToken returns a string representation of the token's exposed data
func ExposeToken(token api.Token) string {
	start, end := token.Pos()
	return fmt.Sprintf("Token{start: %d, end: %d, type: %s, value: %q}", start, end, token.Type().String(), token.String())
}