// =================================================================================================
// Alex Peters - January 25, 2024
// =================================================================================================

package common

import "fmt"

type internalPair [T fmt.Stringer, U any]struct{key T; value U}

type internalTable [T fmt.Stringer, U any]map[string]internalPair[T, U]

type Table [T fmt.Stringer, U any]struct{
	internal internalTable[T, U]
}

func MakeTable[T fmt.Stringer, U any](cap int) *Table[T, U] {
	return &Table[T, U]{make(internalTable[T, U], cap)}
}

func (table *Table[T, U]) Map(k T, v U) {
	s := k.String()
	table.internal[s] = internalPair[T, U]{k,v}
}

func (table *Table[T, U]) Find(k T) (value U, found bool) {
	s := k.String()
	var pair internalPair[T, U]
	if pair, found = table.internal[s]; found {
		value = pair.value
	}
	return
}

func (table *Table[T, U]) Delete(k T) {
	s := k.String()
	delete(table.internal, s)
}