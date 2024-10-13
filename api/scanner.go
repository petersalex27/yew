package api

type Scanner interface {
	// Scan the next token and return it
	Scan() Token
	// Report whether the scanner has reached the end of its known input
	Eof() bool
	// Signal that the scanner should stop scanning
	Stop()
}

type ScannerPlus interface {
	Scanner
	// Add source code to the scanner to scan
	AppendSource(addition string)
	// Restore the scanner to its state before the last AppendSource call
	//
	// Noop when the scanner is already in its initial state
	Restore()
	SrcCode() SourceCode
}