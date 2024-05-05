package types

import "strings"

func joinStringed[T interface{String()string}](elems []T, sep string) string {
	var b strings.Builder
	switch len(elems) {
	case 0:
		return ""
	}

	b.WriteString(elems[0].String())
	for _, elem := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(elem.String())
	}
	return b.String()
}