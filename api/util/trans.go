package util

import (
	"fmt"
	"math"

	"github.com/petersalex27/yew/api"
)

func makeClosures[From, To any](translator api.Translator[From, To]) (done func() bool, update func()) {
	done = translator.Done
	update = func() {} // default is just a noop

	if t, ok := translator.(api.AnalyzingTranslator[From, To]); ok {
		translated := 0
		done = func() bool {
			return translator.Done() || t.UpperBound() == translated
		}
		update = func() {
			// avoid overflow (overflowing would make it possible to put a bound on what should be unbounded)
			if translated == math.MaxInt {
				translated = 0 // reset
			} else {
				translated++
			}
		}
	}

	return done, update
}

// dest may be nil, in which case a new slice will be allocated
//
// NOTE: if translator returns true on `translator.Done()`, this function will loop infinitely
//   - this can be prevented by implementing the following method: `UpperBound() int`
//   - an example of how this could be implemented in a reasonable way is shown in `api/examples/translator.go`
//   - if `UpperBound()` returns a negative value, the number of translations will be treated as unbounded
func TranslateRemaining[From, To any](dest []To, translator api.Translator[From, To]) []To {
	if dest == nil {
		cap := len(translator.Untranslated())
		dest = make([]To, 0, cap)
	}

	done, update := makeClosures(translator)

	for !done() {
		dest = append(dest, translator.Translate())
		update()
	}
	return dest
}

func ExposeTranslator[From, To any](translator api.Translator[From, To]) string {
	untranslated := ExposeList(func(item From) string {
		return fmt.Sprintf("%v", item)
	}, translator.Untranslated(), ", ")
	translated := ExposeList(func(item To) string {
		return fmt.Sprintf("%v", item)
	}, translator.Translated(), ", ")
	config := ExposeConfig(translator.Config())
	return fmt.Sprintf(
		"Translator{Untranslated: %s, Translated: %s, Done: %t, Config: %s}",
		untranslated, translated, translator.Done(), config)
}
