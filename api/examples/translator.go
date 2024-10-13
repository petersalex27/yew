package examples

import (
	"fmt"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util"
)

type PositiveToNonPositive[T interface{ ~int | ~float32 }] struct {
	numbers       []T
	numberPointer int
}

func (p *PositiveToNonPositive[T]) UpperBound() int {
	// because TranslateRemaining is O(n), where n is the number of remaining translations
	return len(p.Untranslated())
}

func (p *PositiveToNonPositive[T]) Translate() (translated T) {
	if p.numbers[p.numberPointer] > 0 {
		translated = -p.numbers[p.numberPointer]
	}
	p.numbers[p.numberPointer] = translated
	p.numberPointer++
	return
}

func (p *PositiveToNonPositive[T]) Untranslated() (untranslated []T) {
	return p.numbers[p.numberPointer:]
}

func (p *PositiveToNonPositive[T]) Translated() (translated []T) {
	return p.numbers[:p.numberPointer]
}

func (p *PositiveToNonPositive[T]) Done() bool {
	return p.numberPointer >= len(p.numbers)
}

func (p *PositiveToNonPositive[T]) Config() api.Config {
	return nil
}

func NewPositiveToNonPositive[T interface{ ~int | ~float32 }](numbers []T) *PositiveToNonPositive[T] {
	return &PositiveToNonPositive[T]{
		numbers:       numbers,
		numberPointer: 0,
	}
}

func ExamplePositiveToNonPositive[T interface{ ~int | ~float32 }](numbers ...T) {
	translator := NewPositiveToNonPositive(numbers)

	fmt.Printf("example/translator (initial): %s\n", util.ExposeTranslator(translator))

	for !translator.Done() {
		_ = translator.Translate()
		fmt.Printf("example/translator: %s\n", util.ExposeTranslator(translator))
	}
}