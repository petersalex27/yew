//go:build !debug
// +build !debug

package parser

import (
	"io"
)

func debug_log_reduce(_ io.Writer, _, _, _ termElem) {}

type debug_info_parser struct {}

// noop
func (*debug_info_parser) debug_incTestCounter() {}

// noop
func (*debug_info_parser) debug_resetCounter() {}