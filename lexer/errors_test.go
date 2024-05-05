// =================================================================================================
// Alex Peters - January 21, 2024
// =================================================================================================
package lexer

import (
	"testing"

	"github.com/petersalex27/yew/source"
)

func TestMakeOSError(t *testing.T) {
	const expected string = "Error (OS): msg"
	actual := makeOSError("msg").Error()
	if actual != expected {
		t.Fatalf("unexpected error message (expected=\"%s\"): got \"%s\"", expected, actual)
	}
}

func TestMakeLexicalError(t *testing.T) {
	const expected string = "Error (Lexical): msg\n[src:1:2-3]"
	actual := makeLexicalError("msg", source.FilePath("src"), 1, 2, 2, 4).Error()
	if actual != expected {
		t.Fatalf("unexpected error message (expected=\"%s\"): got \"%s\"", expected, actual)
	}
}
