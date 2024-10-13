package util_test

import (
	"strings"
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util"
)

type node struct {
	name     string
	children []api.Node
}

func (n node) Pos() (int, int) { return 0, 0 }

func (n node) GetPos() api.Position { return api.ZeroPosition() }

func (n node) Type() api.NodeType { return token.USER_DEFINED_START }

func (n node) Describe() (string, []api.Node) { return n.name, n.children }

func mk(name string, children ...api.Node) node {
	return node{name, children}
}

func TestPrintTree(t *testing.T) {
	tests := []struct {
		input node
		want  string
	}{
		{
			input: mk(""),
			want:  ``,
		},
		{
			input: mk("root"),
			want:  `root`,
		},
		{
			input: mk("root",
				mk("child1",
					mk("child1.1"),
					mk("child1.2")),
				mk("child2",
					mk("child2.1"),
					mk("child2.2"))),
			want: `root
├── child1
│   ├── child1.1
│   └── child1.2
└── child2
    ├── child2.1
    └── child2.2`,
		},
		{
			input: mk("root",
				mk("child1",
					mk("child1.1",
						mk("child1.1.1"),
						mk("child1.1.2")),
					mk("child1.2")),
				mk("child2",
					mk("child2.1"),
					mk("child2.2",
						mk("child2.2.1")))),
			want: `root
├── child1
│   ├── child1.1
│   │   ├── child1.1.1
│   │   └── child1.1.2
│   └── child1.2
└── child2
    ├── child2.1
    └── child2.2
        └── child2.2.1`,
		},
		{
			input: mk("root",
				mk("child1",
					mk("child1.1",
						mk("child1.1.1"),
						mk("child1.1.2")),
					mk("child1.2"))),
			want: `root
└── child1
    ├── child1.1
    │   ├── child1.1.1
    │   └── child1.1.2
    └── child1.2`,
		},
		{
			input: mk("root",
				mk("child1",
					mk("child1.1",
						mk("child1.1.1",
							mk("child1.1.1.1")),
						mk("child1.1.2"),
						mk("child1.1.3")))),
			want: `root
└── child1
    └── child1.1
        ├── child1.1.1
        │   └── child1.1.1.1
        ├── child1.1.2
        └── child1.1.3`,
		},
		{
			input: mk("root",
				mk("child1",
					mk("child1.1",
						mk("child1.1.x2",
							mk("child1.1.x3",
								mk("child1.1.x4"))),
						mk("child1.1.2",
							mk("child1.1.2.1"))))),
			want: `root
└── child1
    └── child1.1
        ├── child1.1.x2
        │   └── child1.1.x3
        │       └── child1.1.x4
        └── child1.1.2
            └── child1.1.2.1`,    
		},
	}

	for _, tt := range tests {
		t.Run(tt.input.name, func(t *testing.T) {
			b := &strings.Builder{}
			util.PrintTree(b, tt.input)
			got := b.String()
			if got != tt.want {
				t.Errorf("printTree() = \n%s\n, want \n%s", got, tt.want)
			}

			//fmt.Printf("tree result:\n%s\n", got)
		})
	}
}
