//go:build !test
// +build !test

package data

import "github.com/petersalex27/yew/api"

func Fail[a api.Node](msg string, positioned api.Positioned) Either[Ers, a] {
	return __Fail[a](msg, positioned)
}

func PassErs[b api.Node](e Ers) Either[Ers, b] { return __PassErs[b](e) }