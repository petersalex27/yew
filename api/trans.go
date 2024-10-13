package api

type (
	Translator[From, To any] interface {
		// Translate the given input to the output
		Translate() To
		// Return current untranslated set
		Untranslated() []From
		// Return current translated set
		Translated() []To
		// Return true if the translator is done; otherwise, return false
		Done() bool
		Config() Config
	}

	AnalyzingTranslator[From, To any] interface {
		Translator[From, To]
		UpperBound() int
	}
)
