// ensures that each example satisfies their respective interface methods
package examples_test

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/examples"
)

var _ api.Translator[int, int] = (*examples.PositiveToNonPositive[int])(nil)
