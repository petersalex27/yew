/*
suggest.go handles generating suggestions on how to fix errors
*/

package error

// Assumption: escape sequence mistakes are (usually) the result of hitting a key near the one you intended
// on accident; and, the correct key is the one closest to the wrong key.
// If two keys are equally close (e.g., 'b' and 'n' are equally close to 'h' on a QWERTY keywoard), then both 
// 'b' and 'n' are offered as suggestions
func IllegalEscapeSuggestion(source string, errorAtIndex int) string {
	// find keyboard layout
	// find closest valid escape key to the one that was hit
	// for example, on a QWERTY keyboard the user types
	// 		\m
	// both the 'n' and 'b' are valid chars for the escape sequence and their respective keys
	// are both close to the 'm' key, but 'n' is closer.
	panic("TODO: Implement!\n")
}
