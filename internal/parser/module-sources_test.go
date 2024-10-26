package parser

import (
	"testing"
)

func TestDetectCycle(t *testing.T) {
	visits := map[string]map[string]struct{}{
		// a -> b -> c
		"a": {"b": {}},
		"b": {"c": {}},
		"c": {"a": {}},

		// d -> e -> f -> g -> h -> i
		"d": {"e": {}, "w": {}},
		"e": {"f": {}, "x": {}, "y": {}},
		"f": {"g": {}},
		"g": {"h": {}},
		"h": {"d": {}},

		// w -> d
		"w": {"d": {}},

		"x": {"z": {}},
		"y": {"z": {}},
		"z": {},

		// make sure that visit[path] gets reset
		//
		// if it doesn't, then this acyclic graph will detect a cycle when it finishes p->q->r and starts p->r
		"p": {"q": {}, "r": {}},
		"q": {"r": {}},
		"r": {},
	}

	tests := []struct {
		start       string
		expectBuf   map[string][]string
		expectCycle bool
	}{
		{
			"a",
			map[string][]string{
				"a": {"c"},
				"b": {"a"},
				"c": {"b"},
			},
			true,
		},
		{
			"b", 
			map[string][]string{
				"a": {"c"},
				"b": {"a"},
				"c": {"b"},
			},
			true,
		},
		{
			"c",
			map[string][]string{
				"a": {"c"},
				"b": {"a"},
				"c": {"b"},
			},
			true,
		},
		{
			"d",
			map[string][]string{
				"h": {"g"},
				"g": {"f"},
				"f": {"e"},
				"e": {"d"},
				"d": {"h", "w"},
				"w": {"d"},
			},
			true,
		},
		{
			"w",
			map[string][]string{
				"h": {"g"},
				"g": {"f"},
				"f": {"e"},
				"e": {"d"},
				"d": {"h", "w"},
				"w": {"d"},
			},
			true,
		},
		{
			"x",
			map[string][]string{},
			false,
		},
		{
			"y",
			map[string][]string{},
			false,
		},
		{
			"z",
			map[string][]string{},
			false,
		},
		{
			"p",
			map[string][]string{},
			false,
		},
		{
			"q",
			map[string][]string{},
			false,
		},
		{
			"r",
			map[string][]string{},
			false,
		},
	}

	for _, test := range tests {
		buf, cycle := detectCycle(map[string][]string{}, make(map[string]bool, len(visits)), visits, test.start)
		if cycle != test.expectCycle {
			t.Errorf("detectCycle(%q) = %v, want %v", test.start, cycle, test.expectCycle)
		}

		if len(buf) != len(test.expectBuf) {
			t.Errorf("detectCycle(%q) = %v, want %v", test.start, buf, test.expectBuf)
		}

		for k, v := range buf {
			if len(v) != len(test.expectBuf[k]) {
				t.Errorf("detectCycle(%q) = %v, want %v", test.start, buf, test.expectBuf)
			}
			for i, p := range v {
				if p != test.expectBuf[k][i] {
					t.Errorf("detectCycle(%q) = %v, want %v", test.start, buf, test.expectBuf)
				}
			}
		}
	}
}
