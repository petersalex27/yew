// =================================================================================================
// Alex Peters - January 26, 2024
//
// for printing source code
// =================================================================================================
package common

import "strings"

type SourcePrinter struct {
	strings.Builder
	indentString []byte
	indent []byte
	start string
}

type PrettyPrintable interface {
	PrettyPrint(*SourcePrinter)
}

func NewSourcePrinter(b strings.Builder, indentSize uint) *SourcePrinter {
	indentString := make([]byte, indentSize)
	for i := range indentString {
		indentString[i] = ' '
	}
	return &SourcePrinter{Builder: b, indentString: indentString, indent: []byte{}}
}

// no-op if indentSize is zero, otherwise increases indent by indentSize spaces
func (printer *SourcePrinter) IncreaseIndent() {
	printer.indent = append(printer.indent, printer.indentString...)
	if printer.start != "" {
		printer.start = string(printer.indent)
	}
}

// no-op if indent cannot be decreased, otherwise decreases indent by indentSize spaces
func (printer *SourcePrinter) DecreaseIndent() {
	
	totLen := len(printer.indent)
	indentLen := len(printer.indentString)
	if totLen < indentLen {
		return
	}

	printer.indent = printer.indent[:totLen-indentLen]
	if printer.start != "" {
		printer.start = string(printer.indent)
	}
}

func (printer *SourcePrinter) Indent() {
	printer.Builder.WriteString(string(printer.indent))
}

func (printer *SourcePrinter) WriteString(s string) {
	printer.Builder.WriteString(printer.start + s)
	printer.start = ""
}

func (printer *SourcePrinter) WriteByte(b byte) error {
	_, e := printer.Builder.WriteString(printer.start + string(b))
	printer.start = ""
	return e
}

func (printer *SourcePrinter) Line() {
	printer.Builder.WriteByte('\n')
	printer.start = string(printer.indent)
}

func PrettyPrintPrintables[P PrettyPrintable](printer *SourcePrinter, printables []P, sep string) {
	if len(printables) > 0 {
		printables[0].PrettyPrint(printer)
	}

	for i := 1; i < len(printables); i++ {
		printer.WriteString(sep)
		printables[i].PrettyPrint(printer)
	}
}