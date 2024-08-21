package util

import (
	"math"
	"strings"
)

type padType byte

const (
	LEFT padType = iota
	RIGHT
	BOTH
)

// StrPad returns the input string padded on the left, right or both sides using padType to the specified padding length padLength.
//
// Example:
// input := "Codes";
// StrPad(input, 10, " ", RIGHT)        // produces "Codes     "
// StrPad(input, 10, "-=", LEFT)        // produces "=-=-=Codes"
// StrPad(input, 10, "_", BOTH)         // produces "__Codes___"
// StrPad(input, 6, "___", RIGHT)       // produces "Codes_"
// StrPad(input, 3, "*", RIGHT)         // produces "Codes"
// taken from // https://gist.github.com/asessa/3aaec43d93044fc42b7c6d5f728cb039
func StrPad(input string, padLength int, padString string, padType padType) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case RIGHT:
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case LEFT:
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case BOTH:
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] + input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
	}

	return output
}