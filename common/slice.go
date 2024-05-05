// =================================================================================================
// Alex Peters - January 27, 2024
// =================================================================================================
package common

// prepends x to xs
func Prepend[T any](x T, xs []T) []T { return append([]T{x}, xs...) }

func Find[T any](predicate func(T) bool, xs []T) (found bool, elem T) {
	for _, elem = range xs {
		if found = predicate(elem); found {
			return
		}
	}
	return
}

func Filter[T any](predicate func(T) bool, xs []T) (elems []T) {
	if len(xs) == 0 {
		return []T{}
	}
	
	elems = make([]T, 0, len(xs))
	n := 0
	for _, elem := range xs {
		if predicate(elem) {
			n++
			elems = append(elems, elem)
		}
	}
	return elems[:n]
}