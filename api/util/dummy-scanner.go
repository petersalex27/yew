package util

import "github.com/petersalex27/yew/api"


type dummyScanner struct{}

func (dummyScanner) Scan() api.Token { panic("scanner not initialized") }

func (dummyScanner) Finished() { panic("scanner not initialized") }

func (dummyScanner) AppendSource(string) { panic("scanner not initialized") }

func (dummyScanner) Restore() { panic("scanner not initialized") }

func (dummyScanner) Eof() bool { panic("scanner not initialized") }

func (dummyScanner) Stop() { panic("scanner not initialized") }

func (dummyScanner) SrcCode() api.SourceCode { panic("scanner not initialized") }