package table

import (
	"fmt"
)

type MultiTable[T fmt.Stringer, U any] struct {
	tabs []Table[T, U]
	i    int
}

// finds the value v mapped to by `k` and calls f(v) and then maps the return value to `k`
//
// no-op if map doesn't exist
func (table *MultiTable[T, U]) Apply(k T, f func(U) U) {
	if v, found := table.Find(k); found {
		table.Map(k, f(v))
	}
}

func (table *MultiTable[T, U]) Decrease() (_ Table[T, U], ok bool) {
	if ok = table.i != 0; !ok {
		return
	}
	out := table.tabs[table.i]
	table.i--
	return out, true
}

func (table *MultiTable[T, U]) Increase() {
	if len(table.tabs) > 0 {
		table.i++
	}

	if table.tabs == nil {
		table.tabs = make([]Table[T, U], 0, 8)
	}
	table.tabs = append(table.tabs, *MakeTable[T, U](8))
}

// remove element key-value pair found at `k` from the table
//
// no-op if pair doesn't exist
func (table *MultiTable[T, U]) Delete(k T) {
	table.tabs[table.i].Delete(k)
}

// return element v in pair (k, v)
//
// if and only if pair does not exist, found=false
func (table *MultiTable[T, U]) Find(k T) (value U, found bool) {
	for i := table.i; i >= 0; i-- {
		if value, found = table.tabs[i].Find(k); found {
			return
		}
	}
	return
}

// returns value in key value pair from the top most table and whether it was found
func (table *MultiTable[T, U]) PeekFind(k T) (value U, found bool) {
	return table.tabs[table.i].Find(k)
}

// return number of pairs in top table
func (table *MultiTable[T, U]) Len() int {
	return table.tabs[table.i].Len()
}

// return the number of tables
func (table *MultiTable[T, U]) NumTables() int {
	return len(table.tabs)
}

// creates a new table with enough space to hold at least cap elements
func NewMultiTable[T fmt.Stringer, U any](cap int) *MultiTable[T, U] {
	out := new(MultiTable[T, U])
	if cap > 0 {
		out.tabs = make([]Table[T, U], 0, cap)
	} // otherwise, default is 8 inside of Increase
	out.Increase() 
	return out
}

// write a pair (k, v); or, if a pair (k, v') already exists, overwrite it
func (table *MultiTable[T, U]) Map(k T, v U) {
	table.tabs[table.i].Map(k, v)
}
