package error

import (
	"fmt"
	"os"
	"strings"
	"yew/info"
)

// Warning reports non-fatal errors
type Warning message

func (w Warning) GetLine() int { return w.line }
func (w Warning) GetChar() int { return w.char }
func (w Warning) GetLocation() info.Location {
	return info.Path(w.path).MakeLocation(w.line, w.char)
}

func (Warning) _primaryMessageStruct_() {}

func (w Warning) shouldAbort() bool { return false }

func (w Warning) ToString() string {
	var builder strings.Builder

	builder.WriteString(writePrefix(w.path, w.line, w.char))
	if NONE != w.subtype {
		builder.WriteString(subtypeMap[w.subtype] + " ")
	}
	builder.WriteString(WARNING + ": ")
	builder.WriteString(w.message)
	builder.WriteString(".")
	//builder.WriteString(writeCodePointer(in, sourceIndex))
	return builder.String()
}

func (w Warning) Print() int {
	out, _ := fmt.Fprintf(os.Stderr, "%s%s\n", warningSprint(w.ToString()), resetSprint(""))
	return out
}