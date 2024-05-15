// =================================================================================================
// Alex Peters - May 05, 2024
// =================================================================================================
package parser

import (
	"fmt"

	"github.com/petersalex27/yew/common/table"
)

type positioned interface {
	Pos() (start, end int)
}

type declMultiTable = table.MultiTable[fmt.Stringer, *Declaration]

type declTable = table.Table[fmt.Stringer, *Declaration]