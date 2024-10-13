package util_test

import (
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util"
	"github.com/petersalex27/yew/internal/source"
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

func TestWindow(t *testing.T) {
	tests := []struct {
		name     string
		source   api.Source
		start    int
		end      int
		expected string
	}{
		{
			name: "Newline sequence window",
			source: mockSource{
				path:    "/path/to/source",
				content: "\n\n\n\n",
			},
			start:    0,
			end:      5,
			expected: "1 | \n2 | \n3 | \n4 | \n5 | ",
		},
		{
			name: "Single line window",
			source: mockSource{
				path:    "/path/to/source",
				content: "single line content",
			},
			start:    0,
			end:      19,
			expected: "1 | single line content",
		},
		{
			name: "Multiple lines window",
			source: mockSource{
				path:    "/path/to/source",
				content: "line 1\nline 2\nline 3\n",
			},
			start:    0,
			end:      21,
			expected: "1 | line 1\n2 | line 2\n3 | line 3",
		},
		{
			name: "Multiple lines window, full",
			source: mockSource{
				path:    "/path/to/source",
				content: "line 1\nline 2\nline 3\n",
			},
			start:    0,
			end:      22,
			expected: "1 | line 1\n2 | line 2\n3 | line 3\n4 | ",
		},
		{
			name: "Multiple lines window, misaligned",
			source: mockSource{
				path:    "/path/to/source",
				content: "line 1\nline 2\nline 3\n",
			},
			start:    3,
			end:      20,
			expected: "1 | line 1\n2 | line 2\n3 | line 3",
		},
		{
			name: "Partial lines window",
			source: mockSource{
				path:    "/path/to/source",
				content: "line 1\nline 2\nline 3\n",
			},
			start:    7,
			end:      14,
			expected: "2 | line 2",
		},
		{
			name: "Partial lines window, misaligned",
			source: mockSource{
				path:    "/path/to/source",
				content: "line 1\nline 2\nline 3\n",
			},
			start:    9,
			end:      10,
			expected: "2 | line 2",
		},
		{
			name: "Empty source window",
			source: mockSource{
				path:    "/path/to/source",
				content: "",
			},
			start:    0,
			end:      0,
			expected: "1 | ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcCode := (source.SourceCode{}).Set(tt.source)
			result := util.Window(srcCode, tt.start, tt.end)
			if result != tt.expected {
				t.Errorf("expected window (raw strings):\n`%s`\n, got \n`%s`\n(interpreted strings):\n%q\n%q\n", tt.expected, result, tt.expected, result)
			}
		})
	}
}
