// =================================================================================================
// Alex Peters - March 06, 2024
//
// various value-test/check/validation functions
// =================================================================================================

package parser

import "unicode"

// true iff identifier starts with an upper-case letter
//
// an upper case letter is one that matches the following regex:
// 	`[A-Z]`
func startsWithUppercase(ident string) bool {
	if len(ident) < 1 {
		return false
	}
	c0 := ident[0]
	return c0 >= 'A' && c0 <= 'Z'
}

func isSymbolicIdent(ident string) bool {
	if len(ident) < 1 {
		return false
	}
	
	for _, r := range ident {
		if unicode.IsLetter(r) {
			return false
		}
		if r == '_' {
			continue
		} else if unicode.IsSymbol(r) {
			return true
		}
	}
	return false // all underscores, whitespace, etc?
}

func validTypeIdent(ident string) bool {
	return startsWithUppercase(ident) || isSymbolicIdent(ident)
}