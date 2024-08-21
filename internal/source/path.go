// =================================================================================================
// Alex Peters - January 18, 2024
//
// # Path specifier interface for either a file system path or stdin
//
// This is useful since the relevant information when the input is stdin is likely just that the
// source of input is the standard input rather than its file system path.
// =================================================================================================
package source

type PathSpec interface {
	Path() string
}

// file system path
type FilePath string

// returns path to file exactly as it is written
func (fp FilePath) Path() string { return string(fp) }

// standard input
type StandardInput struct{}

// returns "stdin"
func (si StandardInput) Path() string { return "stdin" }

var StdinSpec = StandardInput{}

type stringInput struct {
	s string
	name string
}

func StringInput(s string) stringInput {
	return NamedStringInput(s, "string-input")
}

func NamedStringInput(s string, name string) stringInput {
	return stringInput{s: s, name: name}
}

func (si stringInput) Path() string { return si.name }

func (si stringInput) GetInput() string {
	return si.s
}

func AssertStringInput(p PathSpec) (si stringInput, ok bool) {
	si, ok = p.(stringInput)
	return
}
