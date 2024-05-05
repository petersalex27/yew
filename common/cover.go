// =================================================================================================
// Alex Peters - February 16, 2024
//
// Checks for coverage of a set
//
// =================================================================================================

package common

import (
	"fmt"
)

type (
	SetCoverage[T any] struct {
		covered    *Set[T]
		notCovered *Set[T]
	}

	ord interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
			~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
			~float32 | ~float64
	}

	floating interface {
		~float32 | ~float64
	}

	range_t[T ord] struct{ low, high T }

	RangeCoverage[T ord] struct {
		minNonZero T
		max, min   T
		coverage   []range_t[T]
	}
)

func (rc RangeCoverage[T]) doMerge(ins range_t[T], mid int) {
	rc.coverage[mid] = rc.coverage[mid].merge(ins)
	// check prev and next for merge
	if mid-1 > 0 {
		if rc.coverage[mid-1].overlap(rc.coverage[mid], rc.minNonZero, rc.min) {
			tmp := rc.coverage[mid]
			rc.coverage[mid-1] = rc.coverage[mid-1].merge(tmp)
			rc.coverage = append(rc.coverage[:mid], rc.coverage[mid+1:]...)
			mid--
		}
	}
	if mid+1 < len(rc.coverage) {
		if rc.coverage[mid+1].overlap(rc.coverage[mid], rc.minNonZero, rc.min) {
			tmp := rc.coverage[mid+1]
			rc.coverage[mid] = rc.coverage[mid].merge(tmp)
			rc.coverage = append(rc.coverage[:mid+1], rc.coverage[mid+2:]...)
			mid--
		}
	}
	return
}

func (rc RangeCoverage[T]) Insert(val T) {
	ins := range_t[T]{low: val, high: val}

	left, right := 0, len(rc.coverage)
	for left < right {
		mid := left + (right-left)/2
		if rc.coverage[mid].overlap(ins, rc.minNonZero, rc.min) {
			rc.doMerge(ins, mid)
			return
		} else if rc.coverage[mid].low > ins.high {
			right = mid
		} else if rc.coverage[mid].high < ins.low {
			left = mid + 1
		} else {
			panic("should've found overlap") // here incase overlap function is wrong--shouldn't be though
		}
	}

	tmp := make([]range_t[T], 0, len(rc.coverage)+1)
	i := 0
	j := 0
	for ; i < left; i++ {
		tmp = append(tmp, rc.coverage[j])
		j++
	}

	mergedBack := false
	if j > 0 && tmp[i-1].overlap(ins, rc.minNonZero, rc.min) {
		ins = tmp[i-1].merge(ins)
		tmp[i-1] = ins
		mergedBack = true
	}

	if j < len(rc.coverage) && rc.coverage[j].overlap(ins, rc.minNonZero, rc.min) {
		ins = rc.coverage[j].merge(ins)
		if mergedBack {
			tmp[i-1] = ins
		}
		j++
	}

	if !mergedBack {
		tmp = append(tmp, ins)
	}

	tmp = append(tmp, rc.coverage[j:]...)
	rc.coverage = tmp
	return
}

// returns true iff range r1 overlaps with r2
func (r1 range_t[T]) overlap(r2 range_t[T], minNonZero, min T) bool {
	if r1.low == min || r2.low == min {
		return r1.low <= r2.high && r2.low <= r1.high
	}
	return r1.low-minNonZero <= r2.high && r2.low-minNonZero <= r1.high
}

func (r1 range_t[T]) merge(r2 range_t[T]) range_t[T] {
	return range_t[T]{low: Min(r1.low, r2.low), high: Max(r1.high, r2.high)}
}

func (rc RangeCoverage[T]) RangeCovered() bool {
	if len(rc.coverage) != 1 {
		return false
	}
	return rc.coverage[0].low == rc.min && rc.coverage[0].high == rc.max
}

func CoverageOf[T fmt.Stringer](elems ...T) SetCoverage[T] {
	ncset := Create[T](elems...)
	cset := Create[T]()
	return SetCoverage[T]{&cset, &ncset}
}

func (sc SetCoverage[T]) Cover(elem T) (valid bool) {
	if sc.covered.Contains(elem) {
		return true
	}

	before := len(sc.notCovered.internal)
	sc.notCovered.MutRemove(elem)
	after := len(sc.notCovered.internal)
	valid = before == after // something actually removed?
	if valid {
		sc.covered.MutAdd(elem)
	}
	return
}

func (sc SetCoverage[T]) SetCovered() bool {
	return len(sc.notCovered.internal) == 0
}
