package errors

import (
	"fmt"

	"github.com/petersalex27/yew/common"
)

// error message to be printed
type ErrorMessage struct {
	isFatal    bool   // true when compilation should be ended
	SourceName string // source name where error occurred
	Message    string // actual message part of error
	Line       int    // line number
	LineEnd    int
	Start      int // start (inclusive) char number
	End        int // end (exclusive) char number
}

// true if and only if e is a fatal error (not a warning or other kind of message)
func (e ErrorMessage) IsFatal() bool { return e.isFatal }

// given:
//
// if path is not given, uses "_"
//
//   - line # = n, start # = s, end # = e
//     locationString() == "\n[/path/to/src:n:s-e]"
//   - line # = n, start # = s
//     locationString() == "\n[/path/to/src:n:s]"
//   - line # = n
//     locationString() == "\n[/path/to/src:n]"
//   - no numbers, just non-empty path
//     locationString() == "\n[/path/to/src]"
//   - no number, no path
//     locationString() == ""
func (e ErrorMessage) locationString() string {
	srcName := "_"
	if e.SourceName != "" {
		srcName = e.SourceName
	}

	if e.Line <= 0 && srcName == "_" {
		return ""
	}

	s := ""
	if e.Line <= 0 {
		return fmt.Sprintf("\n[%s]", srcName)
	}

	if e.LineEnd-1 > e.Line {
		s = fmt.Sprintf("\n[%s:%d-%d", srcName, e.Line, e.LineEnd-1)
	} else {
		s = fmt.Sprintf("\n[%s:%d", srcName, e.Line)
	}

	if e.Start <= 0 {
		return s + "]"
	}

	if e.LineEnd-1 > e.Line {
		// start is on a diff line than end
		if e.End <= 1 {
			return s + fmt.Sprintf(":%d]", e.Start)
		}
		return s + fmt.Sprintf(":%d,%d]", e.Start, e.End-1)
	}

	// start is on same line as end
	if (e.End - 1) <= e.Start {
		return s + fmt.Sprintf(":%d]", e.Start)
	}
	return s + fmt.Sprintf(":%d-%d]", e.Start, e.End-1)
}

func (e ErrorMessage) Error() string {
	return fmt.Sprintf("%s%s", e.Message, e.locationString())
}

// return string representing underlying struct for error message
//
// NOTE: This is not the method you call to get the error message string, this is for debug/testing
func (e ErrorMessage) String() string {
	return fmt.Sprintf("ErrorMessage{isFatal: %t, Message: \"%s\", Line: %d, Start: %d, End: %d,}", e.isFatal, e.Message, e.Line, e.Start, e.End)
}

// return true if and only if an option below is true:
//   - all fields are set, i.e., e.Line, e.Start, and e.End are set
//   - just e.Line is set
func (e ErrorMessage) hasCodeSnippet() bool {
	return (e.Line > 0 && e.Start > 0 && e.End > 0) || e.Line > 0
}

// set location fields of error message
//
//   - 0 arguments: nothing is set
//
//   - 1 argument: line number is set
//
//   - 2 arguments: arg 0 is line number, arg 1 is start character number
//
//   - 3+ arguments: arg 0 is line number, arg 1 and arg 2 are start and end character number
//     respectively; any remaining arguments are ignored
func (e ErrorMessage) setLocation(line_start_end ...int) ErrorMessage {
	// reference location fields of error
	setter := []*int{&e.Line, &e.LineEnd, &e.Start, &e.End}
	// assign to location fields of error in this order for up to as many as are available (max 3
	// fields)
	for i := 0; i < common.Min(4, len(line_start_end)); i++ {
		*setter[i] = line_start_end[i]
	}

	return e
}

// creates an error message
//
// indexes greater than 2 will be ignored for lineStartEnd. The indexes will be used in that order:
// [0]=line number, [1]=start (inclusive) char, [2]=end (exclusive) char. Char numbers should be
// for chars on line number
func message(message string, line_start_end ...int) ErrorMessage {
	return (ErrorMessage{Message: message}).setLocation(line_start_end...)
}

// Wraps non-empty exception type in parens and appends it (with a single space for padding) to the
// messageType string, then appends ": ". String returned is always non-empty.
//
// Examples:
//
//	header("A", "B") == "A (B): "
//	header("", "B") == "Output (B): "
//	header("A", "") == "A: "
//	header("", "") == "Output: "
//
// NOTE: a default value is given to messageType of "Output" when messageType is an empty string
func header(messageType string, exceptionType string) string {
	if exceptionType != "" {
		exceptionType = " (" + exceptionType + ")"
	}
	if messageType == "" {
		messageType = "Output"
	}

	return messageType + exceptionType + ": "
}

// returns a default value if the message string is empty, else the argument passed
func validateMessageString(msg string) string {
	if msg == "" {
		msg = "unexpected exception" // some unaccounted for and exceptional circumstance
	}
	return msg
}

const (
	errorMsgType   string = "Error"
	warningMsgType string = "Warning"
	logMsgType     string = "Log"
)

// creates a fatal error
func MakeError(subType string, msg string, line_start_end ...int) ErrorMessage {
	msg = header(errorMsgType, subType) + validateMessageString(msg)
	e := message(msg, line_start_end...)
	e.isFatal = true
	return e
}

// creates a warning (non-fatal error)
func MakeWarning(subType string, msg string, line_start_end ...int) ErrorMessage {
	msg = header(warningMsgType, subType) + validateMessageString(msg)
	w := message(msg, line_start_end...)
	return w
}

// creates a log
func MakeLog(subType string, msg string, line_start_end ...int) ErrorMessage {
	msg = header(logMsgType, subType) + validateMessageString(msg)
	log := message(msg, line_start_end...)
	return log
}
