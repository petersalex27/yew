// =================================================================================================
// Alex Peters - January 25, 2024
// =================================================================================================
package common

// given an ordered list (from smallest to largest and assuming the start point is at 0) of end 
// points, finds the index elem belongs to the range of
//
// endPoints is assumed to be continuous, i.e., the value at index[n] is the start point at 
// index[n+1]
//
// endPoints is assumed not to have duplicate values (this is reasonable b/c there shouldn't be a 
// new line w/o at least one newline char)
//
// if not found, returns -1
func SearchRange(endPoints []int, elem int) (index int) {
	if elem < 0 {
		return -1
	}

	left, right := 0, len(endPoints)
	for left < right {
		mid := left + (right - left) / 2
		if endPoints[mid] == elem {
			left = mid + 1
			break
		} else if endPoints[mid] > elem {
			right = mid
		} else {
			left = mid + 1
		}
	}

	if left > len(endPoints) {
		return -1
	}
	return left
} 