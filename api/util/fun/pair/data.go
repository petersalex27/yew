package pair

import "cmp"

type (
	Data[a, b any] struct {
		fst a
		snd b
	}
)

// Pair constructor
func Make[a, b any](fst a, snd b) Data[a, b] { return Data[a, b]{fst, snd} }

// Pair left de-constructor
func Fst[a, b any](p Data[a, b]) a { return p.fst }

// Pair right de-constructor
func Snd[a, b any](p Data[a, b]) b { return p.snd }

// Pair both de-constructor
func (p Data[a, b]) Both() (a, b) { return p.fst, p.snd }

type mapper[a, b, c, d any] struct {
	Data[a, b]
}

// Pair map
func Map[leftTarget, rightTarget, a, b any](pair Data[a, b]) mapper[a, b, leftTarget, rightTarget] {
	return mapper[a, b, leftTarget, rightTarget]{pair}
}

func (m mapper[a, b, c, d]) Both(f func(a) c, g func(b) d) Data[c, d] {
	return Data[c, d]{f(m.fst), g(m.snd)}
}

func (m mapper[a, b, c, d]) Fst(f func(a) c) Data[c, b] {
	return Data[c, b]{f(m.fst), m.snd}
}

func (m mapper[a, b, c, d]) Snd(g func(b) d) Data[a, d] {
	return Data[a, d]{m.fst, g(m.snd)}
}

func Min[a, b interface{ ~int | ~uint }](p1, p2 Data[a, b]) Data[a, b] {
	fst := min(p1.fst, p2.fst)
	snd := min(p1.snd, p2.snd)
	return Data[a, b]{fst, snd}
}

func Max[a, b interface{ ~int | ~uint }](p1, p2 Data[a, b]) Data[a, b] {
	fst := max(p1.fst, p2.fst)
	snd := max(p1.snd, p2.snd)
	return Data[a, b]{fst, snd}
}

func MinPositive[a, b interface{ ~int | ~uint }](p1, p2 Data[a, b]) Data[a, b] {
	var fst a
	var snd b
	if p1.fst <= 0 {
		fst = max(p2.fst, 0)
	} else if p2.fst <= 0 {
		fst = max(p1.fst, 0)
	} else {
		fst = min(p1.fst, p2.fst)
	}

	if p1.snd <= 0 {
		snd = max(p2.snd, 0)
	} else if p2.snd <= 0 {
		snd = max(p1.snd, 0)
	} else {
		snd = min(p1.snd, p2.snd)
	}
	return Data[a, b]{fst, snd}
}

func MaxPositive[a, b interface{ ~int | ~uint }](p1, p2 Data[a, b]) Data[a, b] {
	var fst a
	var snd b
	if p1.fst <= 0 {
		fst = max(p2.fst, 0)
	} else if p2.fst <= 0 {
		fst = max(p1.fst, 0)
	} else {
		fst = max(p1.fst, p2.fst)
	}

	if p1.snd <= 0 {
		snd = max(p2.snd, 0)
	} else if p2.snd <= 0 {
		snd = max(p1.snd, 0)
	} else {
		snd = max(p1.snd, p2.snd)
	}
	return Data[a, b]{fst, snd}
}

type EmbedsPair[a, b any] interface{ ~struct{ Data[a, b] } }

// widens the range of the pair, reordering the values if necessary in ascending order
func WeakenRange[a interface{ ~int | ~uint }](p1, p2 Data[a, a]) Data[a, a] {
	p1, p2 = Ascend(p1), Ascend(p2)
	q1, q2 := MinPositive(p1, p2), MaxPositive(p1, p2)
	return Data[a, a]{q1.fst, q2.snd}
}

func Ascend[a cmp.Ordered](p Data[a, a]) Data[a, a] {
	if p.fst > p.snd {
		return Data[a, a]{p.snd, p.fst}
	}
	return p
}

func Descend[a cmp.Ordered](p Data[a, a]) Data[a, a] {
	if p.fst < p.snd {
		return Data[a, a]{p.snd, p.fst}
	}
	return p
}
