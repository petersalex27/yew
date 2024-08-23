package errors

import (
	"fmt"
	"io"
	"os"

	"github.com/petersalex27/yew/internal/common"
)

type MessageType byte
const (
	AnyMessage MessageType = 0
	ErrorType MessageType = 'e'
	WarningType MessageType = 'w'
	LogType MessageType = 'l'
	TodoType MessageType = 't'
)

// error message to be printed
type ErrorMessage struct {
	isFatal    bool   // true when compilation should be ended
	messageType MessageType 
	SourceName string // source name where error occurred
	Message    string // actual message part of error
	Line       int    // line number
	LineEnd    int
	Start      int // start (inclusive) char number
	End        int // end (exclusive) char number
}

// true if and only if e is a fatal error (not a warning or other kind of message)
func (e ErrorMessage) IsFatal() bool { return e.isFatal }

func (e ErrorMessage) GetType() MessageType {
	return e.messageType
}

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
	todoMsgType    string = "Todo"
)

// creates a fatal error
func MakeError(subType string, msg string, line_start_end ...int) ErrorMessage {
	msg = header(errorMsgType, subType) + validateMessageString(msg)
	e := message(msg, line_start_end...)
	e.isFatal = true
	e.messageType = ErrorType
	return e
}

// creates a warning (non-fatal error)
func MakeWarning(subType string, msg string, line_start_end ...int) ErrorMessage {
	msg = header(warningMsgType, subType) + validateMessageString(msg)
	w := message(msg, line_start_end...)
	w.messageType = WarningType
	return w
}

// creates a log
func MakeLog(subType string, msg string, line_start_end ...int) ErrorMessage {
	msg = header(logMsgType, subType) + validateMessageString(msg)
	log := message(msg, line_start_end...)
	log.messageType = LogType
	return log
}

func MakeTodo(subType string, msg string, line_start_end ...int) ErrorMessage {
	msg = header("Todo", subType) + validateMessageString(msg)
	todo := message(msg, line_start_end...)
	todo.messageType = TodoType
	return todo
}

// writes all messages belonging to `mt` in `es` to `w`, returning number of messages belonging to `mt` in `es`
//
// panics if message type `mt` is not one of `AnyMessage`, `ErrorType`, `WarningType`, or `LogType`
func LogMessage(w io.Writer, mt MessageType, messages []ErrorMessage) (n int) {
	switch mt {
	case AnyMessage, ErrorType, WarningType, LogType:
		return unchecked_LogMessage(w, mt, messages)
	default:
		panic("illegal message type: must be one of `AnyMessage`, `ErrorType`, `WarningType`, or `LogType`")
	}
}

func unchecked_LogMessage(w io.Writer, mt MessageType, messages []ErrorMessage) (n int) {
	for _, e := range messages {
		// true if `mt` and `e.messageType` are the same or either is `AnyType`
		if e.messageType ^ mt == 0 {
			n++
			fmt.Fprintf(w, "%s", e.Error())
		}
	}
	return
}

type ordering byte
const (
	First ordering = iota
	Second
	Third
	Fourth
	fifth
	Exclude
)

func PrintOrdered(es []ErrorMessage, errors, warnings, logs, todos ordering) (n int) {
	ors := []ordering{errors, warnings, logs, todos}
	ms := []MessageType{ErrorType, WarningType, LogType, TodoType}
	m := map[ordering]MessageType{}
	m_inverted := map[MessageType]ordering{}
	for i, o := range ors {
		if o == Exclude {
			continue
		}
		_, found := m[o]
		if found {
			panic("PrintOrdered: no two arguments can be the same")
		}
		m[o] = ms[i]
		m_inverted[ms[i]] = o
	}
	m_inverted[AnyMessage] = fifth

	bins := make([][]ErrorMessage, 5)
	for _, e := range es {
		o, found := m_inverted[e.messageType]
		if !found || o == Exclude {
			continue
		}

		bins[o] = append(bins[o], e)
	}

	// now print the bins in order
	for i, bin := range bins {
		if len(bin) == 0 {
			continue
		}

		n += LogMessage(os.Stderr, m[ordering(i)], bin)
	}
	return n
}

// writes all errors in `es` to `os.Stderr`, returning number of errors occurring in `es`
func PrintErrors(es []ErrorMessage) (n int) {
	return unchecked_LogMessage(os.Stderr, ErrorType, es)
}

// writes all warnings in `es` to `os.Stderr`, returning number of warnings occurring in `es`
func PrintWarnings(es []ErrorMessage) (n int) {
	return unchecked_LogMessage(os.Stderr, WarningType, es)
}

// writes all logs in `es` to `os.Stderr`, returning number of logs occurring in `es`
func PrintLogs(es []ErrorMessage) (n int) {
	return unchecked_LogMessage(os.Stderr, LogType, es)
}

// writes all todos in `es` to `os.Stderr`, returning number of todos occurring in `es`
func PrintTodos(es []ErrorMessage) (n int) {
	return unchecked_LogMessage(os.Stderr, TodoType, es)
}