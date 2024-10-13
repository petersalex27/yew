// =================================================================================================
// Alex Peters - January 29, 2024
// =================================================================================================
package source

import (
	"strings"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util"
)

type SourceCode struct {
	path string
	// source file as an array of strings for each non-empty line, does not include newline chars
	Source []byte
	// records end (exclusive) position for all lines n at index n-1.
	//
	// for example, given
	//	EndPositions = []int{10, 23, 56}
	// line one ends at position 10, line two ends at position 23, and line three ends at position 56
	endPositions []int
}

func (source SourceCode) String() string {
	return string(source.Source)
}

func (src SourceCode) Lines() int {
	return len(src.endPositions)
}

func (src SourceCode) EndPositions() []int {
	return src.endPositions
}

func (src SourceCode) Copy() SourceCode {
	source := make([]byte, len(src.Source))
	positions := make([]int, len(src.EndPositions()))
	out := SourceCode{
		path:         src.path,
		Source:       source,
		endPositions: positions,
	}
	copy(out.Source, src.Source)
	copy(out.endPositions, src.endPositions)
	return out
}

// IMPORTANT: This mutates the source code by appending the given string to the end of the source code.
func (src *SourceCode) AppendSource(addition string) {
	src.Source = append(src.Source, []byte(addition)...)
	src.endPositions = makeEndPositions(src.endPositions, string(src.Source))
}

func (src SourceCode) Path() string {
	return src.path
}

func (src SourceCode) PrepareForWindowing() api.SourceCode {
	// add final newline--this is necessary for (SourceCode).window to work correctly.
	// We want empty files to contain a single--but, importantly--empty line
	freeSource := util.FreeSource(src.path, string(src.Source)+"\n")
	return (SourceCode{}).Set(freeSource)
}

func makeEndPositions(dest []int, content string) []int {
	n := strings.Count(content, "\n")

	if dest == nil {
		dest = make([]int, 0, n+1)
	}

	for i, b := range content {
		if b == '\n' {
			dest = append(dest, i+1)
		}
	}

	// add final line
	dest = append(dest, len(content))
	if len(content) > 0 && content[len(content)-1] == '\n' {
		dest[len(dest)-1] += 1
	}

	return dest
}

func (src SourceCode) Set(source api.Source) api.SourceCode {
	content := source.String()
	src.path = source.Path()
	src.Source = []byte(content)
	src.endPositions = makeEndPositions(nil, content)
	return src
}
