package parser

import "github.com/petersalex27/yew/internal/common/math"

func (a app[T]) pos() (int, int) {
	if len(a.elems) == 0 {
		return 0, 0
	}

	start, end := a.elems[0].pos()

	if len(a.elems) > 1 {
		_, end = a.elems[1].pos()
	}
	
	return start, end
}

func (ts list[T]) pos() (int, int) {
	if len(ts) == 0 {
		return 0, 0
	}

	start, end := ts[0].pos()

	if len(ts) > 1 {
		_, end = ts[len(ts)-1].pos()
	}
	
	return start, end
}

func (old position) reposition(new position) (r position) {
	r.start = math.Min(old.start, new.start)
	r.end = math.Max(old.end, new.end)
	return r
}