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

// finds the value v mapped to by `k` and calls f(v) and then maps the return value to `k`
//
// no-op if map doesn't exist
func (table *Table[T, U]) Apply(k T, f func(U) U) {
	if v, found := table.Find(k); found {
		table.Map(k, f(v))
	}
}

// remove element key-value pair found at `k` from the table
//
// no-op if pair doesn't exist
func (table *Table[T, U]) Delete(k T) {
	s := k.String()
	delete(table.internal, s)
}

// return element v in pair (k, v)
//
// if and only if pair does not exist, found=false
func (table *Table[T, U]) Find(k T) (value U, found bool) {
	s := k.String()
	var pair internalPair[T, U]
	if pair, found = table.internal[s]; found {
		value = pair.value
	}
	return
}

// return number of pairs in table
func (table *Table[T, U]) Len() int {
	return len(table.internal)
}

// creates a new table with enough space to hold at least cap elements
func MakeTable[T fmt.Stringer, U any](cap int) *Table[T, U] {
	return &Table[T, U]{make(internalTable[T, U], cap)}
}

// write a pair (k, v); or, if a pair (k, v') already exists, overwrite it
func (table *Table[T, U]) Map(k T, v U) {
	s := k.String()
	table.internal[s] = internalPair[T, U]{k,v}
}
