package util

import (
	"fmt"

	"github.com/petersalex27/yew/api"
)

// Tokenize scans all tokens from the scanner (from its current input position) and returns them in a slice.
//
// If an error occurs during scanning, the function will return whatever tokens were scanned up to
// that point with a non-nil value for 'errorToken'.
func Tokenize(scanner api.Scanner, buf []api.Token) (_ []api.Token, errorToken *api.Token) {
	if buf == nil {
		const initialCapacity int = 128
		buf = make([]api.Token, 0, initialCapacity)
	}

	for !scanner.Eof() {
		token := scanner.Scan()
		if token.Error() != nil {
			return buf, &token
		}

		buf = append(buf, token)
	}

	return buf, nil
}

func ExposeScanner(scanner api.Scanner) string {
	return "Scanner{Eof: " + fmt.Sprint(scanner.Eof()) + "}"
}