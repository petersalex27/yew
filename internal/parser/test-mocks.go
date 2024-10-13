//go:build test
// +build test

package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/internal/source"
)



type testScanner struct {
	tokens []api.Token
	counter int
	*source.SourceCode
}

func (s *testScanner) Scan() api.Token {
	if s.Eof() {
		return token.EndOfTokens.Make()
	}
	t := s.tokens[s.counter]
	s.counter++
	return t
}

func (s *testScanner) Eof() bool {
	return s.counter >= len(s.tokens)
}

func (s *testScanner) Stop() {}

func (s *testScanner) Restore() {}

func (s *testScanner) SrcCode() api.SourceCode { return s.SourceCode }

var _ api.ScannerPlus = &testScanner{}

