// =================================================================================================
// Modified source code from Go standard lib.
//
// Original license located w/in the same directory as this file inside the file named GO-LICENSE
//
// # Modified by Alex Peters, January 29, 2024
//
// See modification w/in function documentation
//
// Lastly, any help on a better way to deal with this legal stuff would be greatly appreciated
// =================================================================================================
package lexer

import "bytes"

// original documentation:
//
// dropCR drops a terminal \r from the data.
//
// modifications: None, copied directly from Go src/bufio/scan.go source code
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// original documentation:
//
// ScanLines is a split function for a Scanner that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is one optional carriage return followed
// by one mandatory newline. In regular expression notation, it is `\r?\n`.
// The last non-empty line of input will be returned even if it has no
// newline.
//
// modifications:
//   - within the second if statement's block, the line that read
//     `return i + 1, dropCR(data[0:i]), nil` is modified to
//     `return i + 1, dropCR(data[0 : i+1]), nil`
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0 : i+1]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
