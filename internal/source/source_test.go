package source

import (
	"testing"

	"github.com/petersalex27/yew/api"
)

type mockSource struct {
	path    string
	content string
}

func (m mockSource) Path() string {
	return m.path
}

func (m mockSource) String() string {
	return m.content
}

func TestMake(t *testing.T) {
	tests := []struct {
		name     string
		source   api.Source
		expected SourceCode
	}{
		{
			name: "Single line source",
			source: mockSource{
				path:    "/path/to/source",
				content: "single line content",
			},
			expected: SourceCode{
				path:         "/path/to/source",
				Source:       []byte("single line content"),
				endPositions: []int{19},
			},
		},
		{
			name: "Multiple lines source",
			source: mockSource{
				path:    "/path/to/source",
				content: "line 1\nline 2\nline 3\n",
			},
			expected: SourceCode{
				path:         "/path/to/source",
				Source:       []byte("line 1\nline 2\nline 3\n"),
				endPositions: []int{7, 14, 21, 22},
			},
		},
		{
			name: "Leading newline",
			source: mockSource{
				path:    "/path/to/source",
				content: "\nline 1\nline 2\nline 3\n",
			},
			expected: SourceCode{
				path:         "/path/to/source",
				Source:       []byte("\nline 1\nline 2\nline 3\n"),
				endPositions: []int{1, 8, 15, 22, 23},
			},
		},
		{
			name: "Newline sequence",
			source: mockSource{
				path:    "/path/to/source",
				content: "\n\n\n\n",
			},
			expected: SourceCode{
				path:         "/path/to/source",
				Source:       []byte("\n\n\n\n"),
				endPositions: []int{1, 2, 3, 4, 5},
			},
		},
		{
			name: "Empty source",
			source: mockSource{
				path:    "/path/to/source",
				content: "",
			},
			expected: SourceCode{
				path:         "/path/to/source",
				Source:       []byte(""),
				endPositions: []int{0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := (SourceCode{}).Set(tt.source).(SourceCode)
			if result.path != tt.expected.path {
				t.Errorf("expected path %s, got %s", tt.expected.path, result.path)
			}
			if string(result.Source) != string(tt.expected.Source) {
				t.Errorf("expected source %s, got %s", string(tt.expected.Source), string(result.Source))
			}
			if len(result.endPositions) != len(tt.expected.endPositions) {
				t.Fatalf("expected position ranges length %d, got %d", len(tt.expected.endPositions), len(result.endPositions))
			}
			for i, v := range result.endPositions {
				if v != tt.expected.endPositions[i] {
					t.Errorf("expected position range at index %d to be %d, got %d", i, tt.expected.endPositions[i], v)
				}
			}
		})
	}
}
