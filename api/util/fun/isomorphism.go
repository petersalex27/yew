package fun

type Iso[a, b any] interface {
	Trans(a) b
}