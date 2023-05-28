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
	var builder strings.Builder

	builder.WriteString(writePrefix(e.path, e.line, e.char))
	if NONE != e.subtype {
		builder.WriteString(subtypeMap[e.subtype] + " ")
	}
	builder.WriteString(ERROR + ": ")
	builder.WriteString(e.message)
	builder.WriteString(".")
	//builder.WriteString(writeCodePointer(in, sourceIndex))
	return builder.String()
}

func (e Error) Print() int {
	out, _ := fmt.Fprintf(os.Stderr, "%s%s\n", errorSprint(e.ToString()), resetSprint(""))
	return out
}