package error

import (
	"runtime"
	"strconv"
	"strings"
	"yew/info"
	"yew/utils"
	"github.com/fatih/color"
)

type message struct {
	message string
	subtype ErrorSubType
	path string
	line int
	char int
	index int
} 

type primaryMessageStruct interface {
	UserMessage
	_primaryMessageStruct_()
}

type UserMessage interface {
	util.Stringable
	info.Locatable
	shouldAbort() bool
	Print() int
}

var errorSprint = color.New(color.FgRed).Sprint
var resetSprint = color.New(color.Reset).Sprint
var warningSprint = color.New(color.FgYellow).Sprint

func writePrefix(path string, line int, char int) string {
	var builder strings.Builder
	doPrefix := path != "" || (line > 0 && char > 0)
	if !doPrefix {
		return ""
	}
	builder.WriteString("[")
	if path != "" {
		builder.WriteString(path + ":")
	}
	if line > 0 && char > 0 {
		builder.WriteString(strconv.Itoa(line) + ":" + strconv.Itoa(char))
	}
	
	builder.WriteString("] ")
	return builder.String()
}

func PrintBug() {
	var pc []uintptr
	runtime.Callers(2, pc)
	f, _ := runtime.CallersFrames(pc).Next()
	panic("Bug in " + f.Function + ".\n")
}

type ErrorType string

const (
	ERROR = "Error"
	WARNING = "Warning"
)

type ErrorSubType int
const (
	// TYPE errors occur in many places
	//  - when values (which are also types themselves) are used in illegal places
	//  - when a type cannot be deduced
	//  - when a type doesn't exist
	TYPE ErrorSubType = iota
	// INPUT errors occur when user input (source files) are malformed
	INPUT
	// SYSTEM errors occur when Go-runtime returns an error (e.g., I.O. errors)
	SYSTEM
	// SYNTAX errors occur from bad grammar
	SYNTAX
	// VALUE errors occur from bad values
	VALUE
	// NAME errors occur from bad definitions and uses of IDs 
	NAME
	// LOG errors occur during logging
	LOG
	
	// do not add `ErrorSubType`s below! 
	// only NONE should be here!
	
	// NONE error subtype represents no error
	NONE
)

var subtypeMap = map[ErrorSubType]string {
	TYPE: "Type",
	INPUT: "Input",
	SYSTEM: "System",
	SYNTAX: "Syntax",
	VALUE: "Value",
	NAME: "Name",
	LOG: "Logging",
	NONE: "",
}

func fixSubtype(s ErrorSubType) ErrorSubType {
	if s > NONE || s < TYPE {
		return NONE
	}
	return s
}

type _pathToMessage func(path string) _lineToMessage
type _typeToMessage func(type_ ErrorType) _subtypeToMessage
type _subtypeToMessage func(subtype ErrorSubType) _messageToMessage
type _messageToMessage func(msg string) UserMessage
type _lineToMessage func(line int) _charToMessage
type _charToMessage func(char int) _indexToMessage
type _indexToMessage func(index int) _sourceToMessage
type _sourceToMessage func(source string) _typeToMessage
func curryMessage(path string) _lineToMessage {
	return func(line int) _charToMessage {
		return func(char int) _indexToMessage {
			return func(index int) _sourceToMessage {
				return func(source string) _typeToMessage {
					return func (type_ ErrorType) _subtypeToMessage  {
						return func(subtype ErrorSubType) _messageToMessage { 
							return func(msg string) UserMessage {
								return CompileMessage(msg, type_, subtype, path, line, char, index, source)
							}
						}
					}
				}
			}
		}
	}
}

type ErrorLocation struct {
	line int
	char int
	index int
	path string
	source string
}

var BuiltinErrorLocation = ErrorLocation{
	line: 0,
	char: 0,
	index: 0,
	path: "builtin",
	source: "",
}

func MakeErrorLocation(line int, char int, index int, path string, source string) ErrorLocation {
	return ErrorLocation{
		line: line, 
		char: char, 
		index: index, 
		path: path, 
		source: source, 
	}
}

func (e ErrorLocation) GetLineCharString() string {
	return strconv.Itoa(e.line) + ":" + strconv.Itoa(e.char)
}
func (e ErrorLocation) ToString() string {
	if e.path == "" {
		return e.GetLineCharString()
	}
	return e.path + ":" + e.GetLineCharString()
}
func (e ErrorLocation) GetLine() int { return e.line }
func (e ErrorLocation) GetChar() int { return e.char }
func (e ErrorLocation) GetPath() string { return e.path }
func (e ErrorLocation) GetSource() string { return e.source }
func (e ErrorLocation) GetSourceIndex() int { return e.index }

var NoHeaderMessage = curryMessage("")(0)(0)(0)("")
var SystemError = NoHeaderMessage(ERROR)(SYSTEM)
var SyntaxError = SYNTAX.CompileError

// compileMessage creates an error or warning message from the params
func CompileMessage(
		m string, type_ ErrorType, subtype ErrorSubType, 
		path string, line int, char int, index int, source string) UserMessage {
	subtype = fixSubtype(subtype) // makes sure subtype is always valid
	if ERROR == type_ {
		return Error{message: m, subtype: subtype, path: path, line: line, char: char, index: index}
	} else if WARNING == type_ {
		return Warning{message: m, subtype: subtype, path: path, line: line, char: char, index: index}
	}
	return Error{message: m, subtype: subtype}
}

func (est ErrorSubType) CompileError(m string, e ErrorLocation) Error {
	return CompileMessage(m, ERROR, est, e.path, e.line, e.char, e.index, e.source).(Error)
}
func (est ErrorSubType) CompileWarning(m string, e ErrorLocation) Warning {
	return CompileMessage(m, WARNING, est, e.path, e.line, e.char, e.index, e.source).(Warning)
}

func (e ErrorSubType) CompileMessage(
		msg string, type_ ErrorType, 
		path string, line int, char int, index int, source string) UserMessage {
	return CompileMessage(msg, type_, e, path, line, char, index, source)
}