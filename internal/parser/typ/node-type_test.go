package typ

import (
	"testing"
)

func TestNodeType(t *testing.T) {
	// check ranges are correct
	if _last_ <= MinProperNodeTypeValue { // value fits in a byte?
		t.Errorf("last control node type is %d: expected %d > %d (dec. 255)", _last_, _last_, MinProperNodeTypeValue)
	}

	// need this so that we can encode multiple types in a single nodeType
	if _last_node_type_ > MaxProperNodeTypeValue { // value is no larger than 0xff << 0x08?
		t.Errorf("last node type is %d: expected %d <= %d ", _last_node_type_, _last_node_type_, MaxProperNodeTypeValue)
	}
}

func TestNodeTypeString(t *testing.T) {
	for i := MinProperNodeTypeValue; i < _last_node_type_; i++ {
		if nodeType(i).String() == UnknownNodeType.String() {
			t.Errorf("node type #%d has unknown string representation", i)
		}
	}
}

// func TestNodeTypeString(t *testing.T) {
// 	tests := []struct {
// 		input nodeType
// 		want  string
// 	}{
// 		{EmptyType_, "empty"},
// 		{PairType_.BuildType2(Ident, Ident), "identifier and identifier pair"},
// 		{PairType_.BuildType2(Ident, Integer), "identifier and integer literal pair"},
// 		{PairType_.BuildType2(Integer, Ident), "integer literal and identifier pair"},
// 		{PairType_.BuildType(Ident), "identifier and node pair"},
// 		{PairType_, "node and node pair"},
// 		{YewSource+1, "?unknown"},
// 		{Annotation, "annotation"},
// 		{Annotations, "annotations"},
// 		{AppType, "type application"},
// 		{YewSource, "yew source"},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.want, func(t *testing.T) {
// 			got := tt.input.String()
// 			if got != tt.want {
// 				t.Errorf("String() = %q, want %q", got, tt.want)
// 			}
// 		})
// 	}
// }
