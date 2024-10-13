package token

import (
	"fmt"
	"testing"

	"github.com/petersalex27/yew/api"
)

func Test_assert_DescribableNode(t *testing.T) {
	var _ api.DescribableNode = Alias.Make()
	// yippee!
}

func Test_TokenType_String(t *testing.T) {
	for i := 0; i <= int(_token_types_proper_end_); i++ {
		name := fmt.Sprintf("Keyword #%d", i)
		t.Run(name, func(t *testing.T) {
			t.Log(name, ":", Type(i).String())
			if Type(i).String() == UnknownTokenTypeString {
				t.Errorf("Keyword #%d has no string representation", i)
			}
		})
	}
}

func Test_tokenStringMap(t *testing.T) {
	for i := 0; i <= int(_token_types_proper_end_); i++ {
		name := fmt.Sprintf("Keyword #%d", i)
		t.Run(name, func(t *testing.T) {
			t.Log(name, ":", Type(i).String())
			if _, found := tokenStringMap[Type(i)]; !found {
				t.Errorf("Keyword #%d is not found in keywordMap", i)
			}
		})
	}
}