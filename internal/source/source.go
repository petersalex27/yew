// =================================================================================================
// Alex Peters - January 29, 2024
// =================================================================================================
package source

import (
	"fmt"
	//"os"
	"strings"

	"github.com/petersalex27/yew/common"
)

type SourceCode struct {
	Path PathSpec
	// source file as an array of strings for each non-empty line, does not include newline chars
	Source []byte
	// records end (exclusive) position for all lines n at index n-1.
	//
	// for example, given
	//	PositionRanges = []int{10, 23, 56}
	// line one ends at position 10, line two ends at position 23, and line three ends at position 56
	PositionRanges []int
}

// calculates line and char number given a source code position
//
// first return value is line number, second is char
func (source SourceCode) CalcLocation(pos int, isEndPos bool) (line, char int) {
	if len(source.PositionRanges) == 0 {
		return 0, 0
	}

	line = 1 + common.SearchRange(source.PositionRanges, pos, isEndPos) // 1 + result = 0 or greater
	if line > 0 {
		char = (source.PositionRanges[line-1] + 1) - pos
	}
	return
}

// calculates a location range given a `start` and `end` position
//
// first two return values are start line number and end line number respectively. Last two return
// values are start and end char number respectively
func (source SourceCode) CalcLocationRange(start, end int) (line1, line2, char1, char2 int) {
	line1, char1 = source.CalcLocation(start, false)
	line2, char2 = source.CalcLocation(end, true)
	return
}

func (source SourceCode) PointedWindow(start, end int) string {
	//fmt.Fprintf(os.Stderr, "%d %d\n", start, end) // TODO: remove
	window := source.Window(start, end)
	if strings.ContainsRune(window, '\n') {
		return window // cannot point to any area b/c window spans more than one line
	}

	res := source.windowCalculations(start, end)
	if res.empty {
		return ""
	}

	// calculate length of line number header
	initialSkip := len(fmt.Sprintf(res.format, 1))
	// calculate offset from start for pointer
	pointerOffset := start - res.sourceStart
	// calculate pointer length
	pointerLength := end - start
	pointerLine := "\n" +
		strings.Repeat(" ", initialSkip+pointerOffset) +
		strings.Repeat("^", pointerLength)

	return source.windowWrite(res) + pointerLine
}

type windowCalcResult struct {
	empty                                      bool
	lineStart, lineEnd, sourceStart, sourceEnd int
	format                                     string
}

func (source SourceCode) windowCalculations(start, end int) (res windowCalcResult) {
	if len(source.Source) == 0 || len(source.PositionRanges) == 0 {
		res.empty = true
		return res
	}

	if start > end {
		panic("illegal arguments: start > end")
	}

	res.lineStart, _ = source.CalcLocation(start, false)
	res.lineEnd, _ = source.CalcLocation(end, true)
	if res.lineStart < 1 || res.lineEnd < 1 {
		res.empty = true
		return res
	}

	totalLines := len(source.PositionRanges)
	width := common.NumDigits(uint64(totalLines), 10)
	res.format = fmt.Sprintf("%%%dd | ", width)

	// start and end points of actual line
	if res.lineStart != 1 {
		res.sourceStart = source.PositionRanges[res.lineStart-2]
	}
	res.sourceEnd = source.PositionRanges[res.lineEnd-1]
	return res
}

func (source SourceCode) windowWrite(res windowCalcResult) string {
	if res.empty {
		return ""
	}

	// write window w/ line numbers
	var builder strings.Builder
	line := res.lineStart
	builder.WriteString(fmt.Sprintf(res.format, line))
	line++
	src := source.Source[res.sourceStart:res.sourceEnd]
	srcLen := len(src)
	for i, b := range src {
		if b == '\n' {
			if i != srcLen-1 { // don't write final newline
				builder.WriteByte(b)
				builder.WriteString(fmt.Sprintf(res.format, line))
				line++
			}
		} else {
			builder.WriteByte(b)
		}
	}

	return builder.String()
}

// returns a window for the source code according to the start to the end position.
//
// example:
//
//	1 | -- original the source code
//	2 | module example where
//	3 |   trait Functor f where
//	4 |     map : (a -> b) -> f a -> f b
//	5 |     _<$_ : a -> f b -> f a
//	6 |   end
//	7 | end
//	mySrc.Window(x, y) // where x is somewhere on line 3, and y is somewhere of line 4
//	3 |   trait Functor f where
//	4 |     map : (a -> b) -> f a -> f b
func (source SourceCode) Window(start, end int) string {
	res := source.windowCalculations(start, end)
	return source.windowWrite(res)
}
