package stringer

import "fmt"

func Join[stringer fmt.Stringer](xs []stringer, sep string) string {
	if len(xs) == 0 {
		return ""
	}

	out := xs[0].String()
	for i := 1; i < len(xs); i++ {
		out += sep + xs[i].String()
	}
	return out
}