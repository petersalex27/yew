package strings

import (
	"fmt"
	"strings"
)

func Join[S fmt.Stringer](s []S, sep string) string {
	ss := make([]string, len(s))
	for i, s := range s {
		ss[i] = s.String()
	}
	return strings.Join(ss, sep)
}

func FJoin[S, T fmt.Stringer](f func(S) string, s []S, sep string) string {
	ss := make([]string, len(s))
	for i, s := range s {
		ss[i] = f(s)
	}
	return strings.Join(ss, sep)
}

func ReplaceAll[S fmt.Stringer](s S, old, new string) string {
	return strings.ReplaceAll(s.String(), old, new)
}

func Contains[S fmt.Stringer](s S, substr string) bool {
	return strings.Contains(s.String(), substr)
}