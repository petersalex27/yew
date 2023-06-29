package error

import (
	"fmt"
	"os"
	"strings"
	"yew/info"
)

// Error reports fatal errors
type Error message

func (_ Error) _primaryMessageStruct_() {}
func (e Error) shouldAbort() bool { return true }
func (e Error) GetLine() int { return e.line }
func (e Error) GetChar() int { return e.char }
func (e Error) GetLocation() info.Location {
	return info.Path(e.path).MakeLocation(e.line, e.char)
}

func (e Error) ToString() string {
	loc, title, message, snippet := e.ToSplitString()
	var builder strings.Builder

	if loc != "" {
		builder.WriteString(loc)
		builder.WriteByte(' ')
	}
	builder.WriteString(title)
	builder.WriteString(": ")
	builder.WriteString(message)
	builder.WriteByte('.')
	if snippet != "" {
		builder.WriteByte('\n')
		builder.WriteString(snippet)
	}
	return builder.String()
}

func (e Error) ToSplitString() (loc string, title string, message string, snippet string) {
	loc = writePrefix(e.path, e.line, e.char)
	
	title = ""
	if NONE != e.subtype {
		title = subtypeMap[e.subtype] + " "
	}
	title = title + ERROR
	message = e.message
	snippet = e.snippet
	return
}

func (e Error) Print() int {
	out, _ := fmt.Fprintf(
		os.Stderr, "%s%s\n", errorSprint(e.ToString()), resetSprint(""))
	return out
}