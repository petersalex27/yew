package scan

import (
	"fmt"
	"testing"
)

var expectedImagePtrs = []struct{input Input; sourceLine string; index int} {
	{Input{1, 0, 0, 0, 3, "test0", "abc", 0, 0}, "abc", -1}, 		// test before user input
	{Input{1, 0, 1, 1, 3, "test1", "abc", 0, 0}, "abc", 0},  		// test start of user input
	{Input{1, 0, 3, 3, 3, "test2", "abc", 0, 0}, "abc", 2},  		// test end of user input
	{Input{1, 0, 4, 4, 3, "test3", "abc", 0, 0}, "abc ", 3}, 		// test end of file adds space to end
	{Input{2, 3, 0, 4, 7, "test4", "abc\n123", 0, 0}, "123", -1},	// test start of new line
	{Input{2, 3, 3, 7, 7, "test5", "abc\n123", 0, 0}, "123", 2},	// test on line other than first
}
func TestGetInputImageWithPointer(test *testing.T) {
	for _, expected := range expectedImagePtrs {
		sourceLine, index := expected.input.GetInputImageWithPointer()
		if sourceLine != expected.sourceLine {
			// backticks are here to show whitespace
			fmt.Printf("Expected: sourceLine=`%s`\nActual: sourceLine=`%s`\n", expected.sourceLine, sourceLine)
			test.FailNow()
		}
		if index != expected.index {
			fmt.Printf("Expected: index=%d\nActual: index=%d\n", expected.index, index)
			test.FailNow()
		}
	}
}

func TestSourcePointerString(test *testing.T) {
	expected0 := pointerColorSprint(">") + sourceColorSprint("bad before starting")
	actual0 := GetSourcePointerString("bad before starting", -1)
	if actual0 != expected0 {
		fmt.Printf("Expected:\n(%v)\n%s\nActual:\n(%v)\n%s\n", 
				[]byte(expected0), expected0, 
				[]byte(actual0), actual0)
		test.FailNow()
	}
	
	expected1 :=
			"blah blah bloh blah" + "\n" +
			"            " + pointerColorSprint(
						"^") + sourceColorSprint(
						 "      ")
	actual1 := GetSourcePointerString("blah blah bloh blah", 12)
	if actual1 != expected1 {
		fmt.Printf("Expected:\n(%v)\n%s\nActual:\n(%v)\n%s\n", 
				[]byte(expected1), expected1, 
				[]byte(actual1), actual1)
		test.FailNow()
	}
}