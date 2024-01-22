// =============================================================================
// Alex Peters - January 18, 2024
//
// Error message and related tests
// =============================================================================
package errors

import "testing"

func TestHasCodeSnippet(t *testing.T) {
	tests := []struct{
		line, start, end int
		expected bool
	}{
		{line: 0, start: 0, end: 0, expected: false},
		{line: 1, start: 0, end: 0, expected: true},
		{line: 0, start: 1, end: 0, expected: false},
		{line: 1, start: 1, end: 0, expected: true},
		{line: 0, start: 0, end: 1, expected: false},
		{line: 1, start: 0, end: 1, expected: true},
		{line: 0, start: 1, end: 1, expected: false},
		{line: 1, start: 1, end: 1, expected: true},
	}

	for _, test := range tests {
		// mock error w/ just relevant data, i.e., line, start, and end
		e := ErrorMessage{Line: test.line, Start: test.start, End: test.end}
		actual := e.hasCodeSnippet()

		if actual != test.expected {
			t.Fatalf("unexpected result (line=%d, start=%d, end=%d): got %t", test.line, test.start, test.end, actual)
		}
	}
}

func TestSetLocation(t *testing.T) {
	tests := []struct{
		args []int
		expected ErrorMessage
	}{
		{
			args: []int{}, 
			expected: ErrorMessage{},
		},
		{
			args: []int{1}, 
			expected: ErrorMessage{Line: 1,},
		},
		{
			args: []int{1, 2}, 
			expected: ErrorMessage{Line: 1, Start: 2,},
		},
		{
			args: []int{1, 2, 3}, 
			expected: ErrorMessage{Line: 1, Start: 2, End: 3},
		},
		{
			args: []int{1, 2, 3, 4}, 
			expected: ErrorMessage{Line: 1, Start: 2, End: 3},
		},
	}
	
	for _, test := range tests {
		actual := (ErrorMessage{}).setLocation(test.args...)
		if actual.Line != test.expected.Line {
			t.Errorf("unexpected line number (receiver=%v): got %d", test.expected, actual.Line)
		}
		if actual.Start != test.expected.Start {
			t.Errorf("unexpected start number (receiver=%v): got %d", test.expected, actual.Start)
		}
		if actual.End != test.expected.End {
			t.Errorf("unexpected end number (receiver=%v): got %d", test.expected, actual.End)
		}

		if t.Failed() {
			t.FailNow()
		}
	}
}

func TestHeader(t *testing.T) {
	tests := []struct{
		msgType string
		subType string
		expected string
	}{
		{
			msgType: "A",
			subType: "B",
			expected: "A (B): ",
		},
		{
			msgType: "",
			subType: "B",
			expected: "Output (B): ",
		},
		{
			msgType: "A",
			subType: "",
			expected: "A: ",
		},
		{
			msgType: "",
			subType: "",
			expected: "Output: ",
		},
	}

	for _, test := range tests {
		actual := header(test.msgType, test.subType)
		if actual != test.expected {
			t.Fatalf("unexpected result from \"%s\", \"%s\": got \"%s\"", test.msgType, test.subType, actual)
		}
	}
}

func TestLocationString(t *testing.T) {
	tests := []struct{
		line, start, end int
		path string
		expected string
	}{
		{
			line: 0, start: 0, end: 0,
			path: "",
			expected: "",
		},
		{
			line: 1, start: 0, end: 0,
			path: "",
			expected: "\n[_:1]",
		},
		{
			line: 1, start: 0, end: 0,
			path: "a",
			expected: "\n[a:1]",
		},
		{
			line: 0, start: 1, end: 0,
			path: "a",
			expected: "\n[a]",
		},
		{
			line: 1, start: 1, end: 0,
			path: "a",
			expected: "\n[a:1:1]",
		},
		{
			line: 1, start: 1, end: 0,
			path: "",
			expected: "\n[_:1:1]",
		},
		{
			line: 1, start: 1, end: 3,
			path: "a",
			expected: "\n[a:1:1-2]",
		},
		{
			line: 1, start: 1, end: 3,
			path: "",
			expected: "\n[_:1:1-2]",
		},
		{
			line: 1, start: 1, end: 2,
			path: "",
			expected: "\n[_:1:1]",
		},
	}

	for _, test := range tests {
		e := ErrorMessage{
			Line: test.line,
			Start: test.start,
			End: test.end,
			SourceName: test.path,
		}

		actual := e.locationString()
		if actual != test.expected {
			t.Fatalf("unexpected location string from line=%d, start=%d, end=%d, path=\"%s\": got \"%s\"", test.line, test.start, test.end, test.path, actual)
		}
	}
}