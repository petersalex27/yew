package util

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common"
)

func CalcLocation(source api.SourceCode, pos int, isEndPos bool) (line, char int) {
	endPositions := source.EndPositions()
	if len(endPositions) == 0 {
		return 0, 0
	}

	line = 1 + common.SearchRange(endPositions, pos, isEndPos) // 1 + result = 0 or greater
	if line > 0 {
		char = (endPositions[line-1] + 1) - pos
	}
	return line, char
}

func CalcLocationRange(source api.SourceCode, start, end int) (line1, line2, char1, char2 int) {
	line1, char1 = CalcLocation(source, start, false)
	line2, char2 = CalcLocation(source, end, true)
	return line1, line2, char1, char2
}

type windowCalcResult struct {
	empty                                      bool
	lineStart, lineEnd, sourceStart, sourceEnd int
	format                                     string
}

func windowCalculations(source api.SourceCode, start, end int) (res windowCalcResult) {
	endPositions := source.EndPositions()
	if len(source.String()) == 0 || len(endPositions) == 0 {
		res.empty = true
		return res
	}

	// trailing newlines will extend the end position beyond the actual source code
	end = min(end, len(source.String()))

	if start > end {
		panic("illegal arguments: start > end")
	}

	res.lineStart, res.lineEnd, _, _ = CalcLocationRange(source, start, end)
	if res.lineStart < 1 || res.lineEnd < 1 {
		res.empty = true
		return res
	}

	totalLines := source.Lines()
	width := common.NumDigits(uint64(totalLines), 10)
	res.format = fmt.Sprintf("%%%dd | ", width)

	// start and end points of actual line
	if res.lineStart != 1 {
		res.sourceStart = endPositions[res.lineStart-2]
	}
	res.sourceEnd = endPositions[res.lineEnd-1]
	return res
}

func windowWrite(source api.SourceCode, res windowCalcResult) string {
	if res.empty {
		return ""
	}

	// write window w/ line numbers
	var builder strings.Builder
	line := res.lineStart
	builder.WriteString(fmt.Sprintf(res.format, line))
	line++
	src := []byte(source.String())[res.sourceStart:res.sourceEnd]
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

func window(source_ api.SourceCode, start, end int, pointed bool) string {
	var source api.SourceCode
	if unprepared, ok := source_.(api.UnpreparedWindowingSourceCode); ok {
		source = unprepared.PrepareForWindowing()
	} else {
		source = source_
	}
	
	res := windowCalculations(source, start, end)
	window := windowWrite(source, res)
	if !pointed || strings.ContainsRune(window, '\n') {
		return window
	}

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

	return window + pointerLine
}

// returns a window for the source code according to the start to the end position.
//
// example:
//
//	1 | -- original the source code
//	2 | module example
//	3 | spec Functor f where
//	4 |   map : (a -> b) -> f a -> f b
//	5 |   (<$) : a -> f b -> f a
//
//	mySrc.Window(x, y), where x is somewhere on line 3, and y is somewhere of line 4
//	3 | spec Functor f where
//	4 |   map : (a -> b) -> f a -> f b
func Window(source api.SourceCode, start, end int) string {
	return window(source, start, end, false)
}

func PointedWindow(source api.SourceCode, start, end int) string {
	return window(source, start, end, true)
}