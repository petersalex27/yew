package inf

import (
	"github.com/petersalex27/yew/types"
)

type Conclusion[T types.Matchable] struct {
	Expression T
	Type       types.Monotype
	Status
}

func Conclude[T types.Matchable](e T, t types.Monotype) Conclusion[T] {
	return Conclusion[T]{e, t, Ok}
}

func CannotConclude[T types.Matchable](stat Status) Conclusion[T] {
	return Conclusion[T]{Status: stat}
}
