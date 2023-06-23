package yew

import (
	"fmt"
	"os"
	"strings"
	"yew/utils"
)

type ExpectedVsActual[I any, T util.Stringable] struct {
	Input I
	Expected util.Equatable[T]
	Actual util.Equatable[T]
}

func (ea ExpectedVsActual[E, A]) ToString() string {
	var builder strings.Builder
	builder.WriteString("= Expected === \n")
	builder.WriteString(ea.Expected.(util.Stringable).ToString())
	builder.WriteString("\n= Actual ===== \n")
	builder.WriteString(ea.Actual.(util.Stringable).ToString())
	return builder.String()
}

func PrintExpectedVsActual(expected string, actual string) {
	fmt.Fprintf(
		os.Stderr, 
		"= expected (len=%d) ===\n%s\n= actual (len=%d) =====\n%s\n", 
		len(expected), expected,
		len(actual), actual)
}