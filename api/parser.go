package api

type Parser interface {
	// Parse the input and return true if successful. Parsed result can be accessed via Ast method.
	Parse() bool
	// returns the current AST root
	Ast() Node
	// record an error that occurred during parsing
	AddError(error)
	// return all errors that occurred during parsing
	Errors() []error
	// This is used to initialize the parser only. Simply return a pointer to the field corresponding to the scanner.
	//
	// If the returned value is nil, a dummy scanner will be set that panics when used.
	ReferenceScanner() *Scanner
	// this should return a new instance of the parser with all fields reset to their initial state
	Clear() Parser
}
