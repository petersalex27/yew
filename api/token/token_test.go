package token

import (
	"testing"

	"github.com/petersalex27/yew/api"
)

func Test_assert_DescribableNode(t *testing.T) {
	var _ api.DescribableNode = Alias.Make()
	// yippee!
}
