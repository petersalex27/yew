// =================================================================================================
// Alex Peters - January 18, 2024
//
// # Path specifier interface for either a file system path or stdin
//
// This is useful since the relevant information when the input is stdin is likely just that the
// source of input is the standard input rather than its file system path.
// =================================================================================================
package lexer

type pathSpec interface {
	String() string
}

// file system path
type FilePath string

// returns path to file exactly as it is written
func (fp FilePath) String() string { return string(fp) }

// standard input
type standardInput struct{}

// returns "stdin"
func (si standardInput) String() string { return "stdin" }

var StdinSpec = standardInput{}
