package scan

import (
	color "github.com/fatih/color"
)

var pointerColorSprint = color.New(color.FgMagenta).SprintFunc()
var sourceColorSprint = color.New(color.Reset).SprintFunc()

func GetSourcePointerString(sourceLine string, pointerIndex int) string {
	if pointerIndex < 0 {
		// output for GetSourcePointerString("bad before starting", -1)
		// >bad before starting
		return pointerColorSprint(">") + sourceColorSprint(sourceLine)
	}

	pointer := make([]byte, len(sourceLine))
	for i := range pointer {
		pointer[i] = ' '
	}
	str := string(pointer)

	// output for GetSourcePointerString("blah blah bloh blah", 12)
	// blah blah bloh blah
	//             ^
	return sourceLine + "\n" + str[:pointerIndex] + 
			pointerColorSprint("^") + sourceColorSprint(str[pointerIndex + 1:]) 
}

// pointerIndex might be -1!!!
func (in *Input) GetInputImageWithPointer() (sourceLine string, pointerIndex int) {
	copyInput := Input{
		lineNumber: 1,
		prevLineLength: 0,
		charNumber: 0,
		sourceIndex: 0,
		sourceLength: in.sourceLength,
		path: in.path,
		source: in.source,
	}

	// keep reading lines until reaching the current line
	for ; copyInput.lineNumber < in.lineNumber ; {
		if copyInput.readUntil('\n') == 0 {
			return "", 0 // something went wrong
		}
	}

	// line start index
	lineStartIndex := copyInput.sourceIndex 
	// find location of error in current line, i.e., the char - 1
	pointerIndex = in.charNumber - 1
	// get length of line
	copyInput.readUntil('\n') // reads until '\n' or EOF
	lineLength := copyInput.sourceIndex - lineStartIndex // length of line
	sourceLine = in.source[lineStartIndex : lineStartIndex + lineLength]
	if pointerIndex >= lineLength {
		pointerIndex = lineLength
		sourceLine = sourceLine + " "
	}
	return
}