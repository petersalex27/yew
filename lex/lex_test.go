package scan

import (
	"fmt"
	"os"

	//"os"
	//"unicode"

	//"os"
	"testing"
	"yew/value"
	//err "yew/error"
)

type TestToken struct {
	expected Token
}

type TestInput struct {
	path           string
	expected       []string
	expectedTokens []Token
}

const folderLoc = "./test/"

var (
	INT_TOKEN           = OtherToken{tokenType: INT, line: 1, char: 1}
	CHAR_TOKEN          = OtherToken{tokenType: CHAR, line: 1, char: 5}
	BOOL_TOKEN          = OtherToken{tokenType: BOOL, line: 1, char: 10}
	FLOAT_TOKEN         = OtherToken{tokenType: FLOAT, line: 1, char: 15}
	STRING_TOKEN        = OtherToken{tokenType: STRING, line: 1, char: 21}
	LET_TOKEN           = OtherToken{tokenType: LET, line: 1, char: 28}
	MUT_TOKEN           = OtherToken{tokenType: MUT, line: 1, char: 32}
	CONST_TOKEN         = OtherToken{tokenType: CONST, line: 1, char: 36}
	WHERE_TOKEN         = OtherToken{WHERE, 0, 1}
	CLASS_TOKEN         = OtherToken{CLASS, 0, 1}
	TRUE_TOKEN          = ValueToken{value.Bool(true), 0, 1}
	FALSE_TOKEN         = ValueToken{value.Bool(false), 0, 1}
	NL_TOKEN            = OtherToken{NEW_LINE, 0, 2}
	LPAREN_TOKEN        = OtherToken{LPAREN, 0, 2, }
	RPAREN_TOKEN        = OtherToken{RPAREN, 0, 2, }
	LBRACK_TOKEN        = OtherToken{LBRACK, 0, 2, }
	RBRACK_TOKEN        = OtherToken{RBRACK, 0, 2, }
	LCURL_TOKEN         = OtherToken{LCURL, 0, 2, }
	RCURL_TOKEN         = OtherToken{RCURL, 0, 2, }
	PLUS_PLUS_TOKEN     = OtherToken{PLUS_PLUS, 0, 2, }
	PLUS_TOKEN          = OtherToken{PLUS, 0, 2, }
	MINUS_TOKEN         = OtherToken{MINUS, 0, 2, }
	STAR_TOKEN          = OtherToken{STAR, 0, 2, }
	SLASH_TOKEN         = OtherToken{SLASH, 0, 2, }
	HAT_TOKEN           = OtherToken{HAT, 0, 2, }
	EQUALS_TOKEN        = OtherToken{EQUALS, 0, 2, }
	PLUS_EQUALS_TOKEN   = OtherToken{PLUS_EQUALS, 0, 2, }
	MINUS_EQUALS_TOKEN  = OtherToken{MINUS_EQUALS, 0, 2, }
	STAR_EQUALS_TOKEN   = OtherToken{STAR_EQUALS, 0, 2, }
	SLASH_EQUALS_TOKEN  = OtherToken{SLASH_EQUALS, 0, 2, }
	COMMA_TOKEN         = OtherToken{COMMA, 0, 2, }
	DOT_TOKEN           = OtherToken{DOT, 0, 2, }
	DOT_DOT_TOKEN       = OtherToken{DOT_DOT, 0, 2, }
	BANG_TOKEN          = OtherToken{BANG, 0, 2, }
	BANG_EQUALS_TOKEN   = OtherToken{BANG_EQUALS, 0, 2, }
	EQUALS_EQUALS_TOKEN = OtherToken{EQUALS_EQUALS, 0, 2, }
	AMPER_AMPER_TOKEN   = OtherToken{AMPER_AMPER, 0, 2, }
	BAR_BAR_TOKEN       = OtherToken{BAR_BAR, 0, 2, }
	BAR_TOKEN           = OtherToken{BAR, 0, 2, }
	GREAT_TOKEN         = OtherToken{GREAT, 0, 2, }
	GREAT_EQUALS_TOKEN  = OtherToken{GREAT_EQUALS, 0, 2, }
	LESS_TOKEN          = OtherToken{LESS, 0, 2, }
	LESS_EQUALS_TOKEN   = OtherToken{LESS_EQUALS, 0, 2, }
	ARROW_TOKEN         = OtherToken{ARROW, 0, 2, }
	FAT_ARROW_TOKEN     = OtherToken{FAT_ARROW, 0, 2, }
	SEMI_COLON_TOKEN    = OtherToken{SEMI_COLON, 0, 2, }
	COLON_TOKEN         = OtherToken{COLON, 0, 2, }
	COLON_COLON_TOKEN   = OtherToken{COLON_COLON, 0, 2, }
	EOF_TOKEN           = OtherToken{EOF, 0, 2, }
)

var inputCases = []TestInput{
	{
		path: folderLoc + "in.yw",
		expected: []string{
			"Int Char Bool Float String let mut const where class True False\n",
			"()[]{} ++ + - * / ^ = += -= *= /= , . .. ! != == && || | > >= < <= -> => ; : ::\n",
			"?",
		},
		expectedTokens: []Token{
			INT_TOKEN, CHAR_TOKEN, BOOL_TOKEN, FLOAT_TOKEN, STRING_TOKEN, LET_TOKEN, MUT_TOKEN, CONST_TOKEN,
			WHERE_TOKEN, CLASS_TOKEN, TRUE_TOKEN, FALSE_TOKEN, NL_TOKEN, LPAREN_TOKEN, RPAREN_TOKEN, LBRACK_TOKEN,
			RBRACK_TOKEN, LCURL_TOKEN, RCURL_TOKEN, PLUS_PLUS_TOKEN, PLUS_TOKEN, MINUS_TOKEN, STAR_TOKEN, SLASH_TOKEN,
			HAT_TOKEN, EQUALS_TOKEN, PLUS_EQUALS_TOKEN, MINUS_EQUALS_TOKEN, STAR_EQUALS_TOKEN, SLASH_EQUALS_TOKEN,
			COMMA_TOKEN, DOT_TOKEN, DOT_DOT_TOKEN, BANG_TOKEN, BANG_EQUALS_TOKEN, EQUALS_EQUALS_TOKEN, AMPER_AMPER_TOKEN,
			BAR_BAR_TOKEN, BAR_TOKEN, GREAT_TOKEN, GREAT_EQUALS_TOKEN, LESS_TOKEN, LESS_EQUALS_TOKEN, ARROW_TOKEN,
			FAT_ARROW_TOKEN, SEMI_COLON_TOKEN, COLON_TOKEN, COLON_COLON_TOKEN,
			OtherToken{NEW_LINE, -1, 3, }, OtherToken{QUESTION, -1, 3, }, OtherToken{EOF, -1, 3, },
		},
	},
	{
		path:     folderLoc + "int.yw",
		expected: []string{`1 12 1_000_000 0xab_cdef 0xAB_CDEF 0Xf 0o7 0O7 0b1 0B1`},
		expectedTokens: []Token{
			ValueToken{value.Int(1), 0, 0, },
			ValueToken{value.Int(12), 0, 0, },
			ValueToken{value.Int(1000000), 0, 0, },
			ValueToken{value.Int(0xabcdef), 0, 0, },
			ValueToken{value.Int(0xABCDEF), 0, 0, },
			ValueToken{value.Int(0xf), 0, 0, },
			ValueToken{value.Int(0o7), 0, 0, },
			ValueToken{value.Int(0o7), 0, 0, },
			ValueToken{value.Int(0b1), 0, 0, },
			ValueToken{value.Int(0b1), 0, 0, },
			OtherToken{EOF, 0, 0, },
		},
	},
	{
		path:     folderLoc + "char.yw",
		expected: []string{`'a' '1' '\\' '\'' ' '`},
		expectedTokens: []Token{
			ValueToken{value.Char('a'), 0, 0, },
			ValueToken{value.Char('1'), 0, 0, },
			ValueToken{value.Char('\\'), 0, 0, },
			ValueToken{value.Char('\''), 0, 0, },
			ValueToken{value.Char(' '), 0, 0, },
			OtherToken{EOF, 0, 0, },
		},
	},
	{
		path:     folderLoc + "float.yw",
		expected: []string{`1.1 1e1 1.1e1 1.1e+1 1.1e-1 1.1E1`},
		expectedTokens: []Token{
			ValueToken{value.Float(float64(1.1)), 0, 0, },
			ValueToken{value.Float(float64(1e1)), 0, 0, },
			ValueToken{value.Float(float64(1.1e1)), 0, 0, },
			ValueToken{value.Float(float64(1.1e+1)), 0, 0, },
			ValueToken{value.Float(float64(1.1e-1)), 0, 0, },
			ValueToken{value.Float(float64(1.1e1)), 0, 0, },
			OtherToken{EOF, 0, 0, },
		},
	},
	{
		path:     folderLoc + "string.yw",
		expected: []string{`"abc123" "" "\n" "\\\n\r\t\b\'\""`},
		expectedTokens: []Token{
			ValueToken{stringValue("abc123"), 0, 0, },
			ValueToken{stringValue(""), 0, 0, },
			ValueToken{stringValue("\n"), 0, 0, },
			ValueToken{stringValue("\\\n\r\t\b'\""), 0, 0, },
			OtherToken{EOF, 0, 0, },
		},
	},
	/*{
		path: folderLoc + "infix-id.yw",
		expected: "(>>)",
		expectedTokens: []Token{
			OtherToken{LPAREN, 0, 0, },
			IdToken{">>", 0, 0, },
			OtherToken{RPAREN, 0, 0, },
			OtherToken{EOF, 0, 0, },
		},
	},
	{
		path: folderLoc + "prefix-non-ascii-id.yw",
		expected: "`$`",
		expectedTokens: []Token{
			IdToken{"$", 0, 0, },
			OtherToken{EOF, 0, 0, },
		},
	},*/
}

func TestInit(t *testing.T) {
	for _, cs := range inputCases {
		in, e := Init(cs.path)
		if nil != e {
			fmt.Printf("Test Failed (Unexpected Error): %v\n", e)
			t.FailNow()
		} else if len(in.source) != len(cs.expected) {
			fmt.Printf("Test Failed len(expected) != len(actual)\n")
			t.FailNow()
		}
		for i := range in.source {
			if in.source[i] != cs.expected[i] {
				fmt.Printf(
					"Test Failed (expected.source[%d]=%s) != (actual.source[%d]=%s)\n",
					i, cs.expected[i],
					i, in.source[i])
				t.FailNow()
			}
		}
	}
}

// unknown token
var _in0 = Input{1, 0, "test0", []string{`£`}}
var _in0_img = Input{1, 1, "test0", []string{`£`}}

// illegal control
var _in1 = Input{1, 0, "test1", []string{"\"\n\""}}
var _in1_img = Input{1, 2, "test1", []string{"\"\n\""}}

// trailing underscore
var _in2 = Input{1, 0, "test2", []string{`1_`}}
var _in2_img = Input{1, 2, "test2", []string{`1_`}}

// expected char end
var _in3 = Input{1, 0, "test3", []string{`'ab'`}}
var _in3_img = Input{1, 2, "test3", []string{`'ab'`}}

// illegal escape
var _in4 = Input{1, 0, "test4", []string{`'\w'`}}
var _in4_img = Input{1, 3, "test4", []string{`'\w'`}}

// string only escape
var _in5 = Input{1, 0, "test5", []string{`'\u'`}}
var _in5_img = Input{1, 3, "test5", []string{`'\u'`}}

// malformed unicode escape (1)
var _in6 = Input{1, 0, "test6", []string{`"\u"`}}
var _in6_img = Input{1, 3, "test6", []string{`"\u"`}}

// malformed unicode escape (2)
var _in7 = Input{1, 0, "test7", []string{`"\uX"`}}
var _in7_img = Input{1, 3, "test7", []string{`"\uX"`}}

// malformed unicode escape (3)
var _in8 = Input{1, 0, "test8", []string{`"\uXf"`}}
var _in8_img = Input{1, 3, "test8", []string{`"\uXf"`}}

// illegal char lit
var _in9 = Input{1, 0, "test9", []string{`''`}}
var _in9_img = Input{1, 2, "test9", []string{`''`}}

var expectedNextErrors = []struct {
	input         Input
	expectedInput Input
	expected      ErrorToken
}{
	{
		_in0,
		_in0_img,
		inputErrors[E_UNEXPECTED_TOKEN](&_in0_img),
	},
	{
		_in1,
		_in1_img,
		inputErrors[E_UNEXPECTED_CONTROL](&_in1_img),
	},
	{
		_in2,
		_in2_img,
		inputErrors[E_TRAILING_UNDERSCORE](&_in2_img),
	},
	{
		_in3,
		_in3_img,
		inputErrors[E_EXPECTED_CHAR_CLOSE](&_in3_img),
	},
	{
		_in4,
		_in4_img,
		inputErrors[E_ILLEGAL_ESCAPE](&_in4_img),
	},
	{
		_in5,
		_in5_img,
		inputErrors[E_STRING_ONLY_ESCAPE](&_in5_img),
	},
	{
		_in6,
		_in6_img,
		inputErrors[E_BAD_UNICODE](&_in6_img),
	},
	{
		_in7,
		_in7_img,
		inputErrors[E_BAD_UNICODE](&_in7_img),
	},
	{
		_in8,
		_in8_img,
		inputErrors[E_BAD_UNICODE](&_in8_img),
	},
	{
		_in9,
		_in9_img,
		inputErrors[E_ILLEGAL_CHAR](&_in9_img),
	},
}

func TestMatch(test *testing.T) {
	in := InputStream{
		path:            "test",
		streamIndex:     0,
		streamLength:    6,
		streamCapacity:  6,
		source:          []string{"let x Int = 1"},
		asStringPattern: "",
		tokens:          []Token{LET_TOKEN, MakeIdToken("x", 0, 0), INT_TOKEN, EQUALS_TOKEN, ValueToken{value.Int(1), 0, 0}, EOF_TOKEN},
	}
	pats := []TokenPattern{
		CompileTokenPattern([]TokenType{LET, ID, INT, EQUALS, VALUE, EOF}),
		CompileTokenPattern([]TokenType{LET}),
		CompileTokenPattern([]TokenType{_START_GROUP__, LET, ID, INT, _END_GROUP__}),
		CompileTokenPattern([]TokenType{_START_GROUP__, ID, _END_GROUP__, _REPEAT__, 0, LET, ID}),
		CompileTokenPattern([]TokenType{ID}),
	}
	res := make([]int, 5)
	expect := []int{6, 1, 3, 2, 0}
	for i := range res {
		res[i] = in.Match(pats[i])
		if res[i] != expect[i] {
			fmt.Fprintf(os.Stderr, "Expected: %d\nActual: %d\n", expect[i], res[i])
			test.FailNow()
		}
	}
}

func TestNext(test *testing.T) {
	for _, cs := range inputCases {
		in, e := Init(cs.path)
		if nil != e {
			fmt.Printf("Test Failed (Unexpected Error): %v\n", e)
			test.FailNow()
		}

		i := 0
		for t := in.Next(); ; t = in.Next() {
			if t.GetType() != cs.expectedTokens[i].GetType() {
				fmt.Fprintf(os.Stderr, "Failed %s\n", cs.path)
				fmt.Printf("Expected: \"%s\"\nFound: \"%s\"\n",
					cs.expectedTokens[i].ToString(), t.ToString())
				test.FailNow()
			}
			if VALUE == t.GetType() {
				if t.ToString() != cs.expectedTokens[i].ToString() {
					print("Expected Value: \"", cs.expectedTokens[i].ToString(),
						"\"\nFound Value: \"", t.ToString(), "\"\n")
					test.FailNow()
				}
			} else if ID == t.GetType() {
				if t.(IdToken).id != cs.expectedTokens[i].(IdToken).id {
					fmt.Printf("Expected: \"%s\"\nFound: \"%s\"\n",
						cs.expectedTokens[i].(IdToken).id, t.(IdToken).id)
					test.FailNow()
				}
			}
			i++
			if EOF == t.GetType() {
				break
			}
		}
	}

	checkErrorToken := func(actual ErrorToken, exp ErrorToken) bool {
		return actual.err.ToString() == exp.err.ToString()
	}

	// test bad input
	for _, e := range expectedNextErrors {
		bad := e.input
		tok := bad.Next()
		//fmt.Fprintf(os.Stderr, "~~~~~~~~~~~~%t\n", unicode.IsLetter(rune('£')))
		if ERROR != tok.GetType() {
			fmt.Printf("Expected: %s\nActual: %s\n", ERROR.ToString(), tok.GetType().ToString())
			test.FailNow()
		}
		if !bad.equal_test(e.expectedInput) {
			fmt.Printf("Expected: \n%s\n", e.expectedInput.toString_test())
			fmt.Printf("Actual: \n%s\n", bad.toString_test())
			test.FailNow()
		}
		if !checkErrorToken(tok.(ErrorToken), e.expected) {
			fmt.Printf("Expected: %s\nActual: %s\n", e.expected.err.ToString(), tok.(ErrorToken).err.ToString())
			test.FailNow()
		}
	}
}

var testGetNextIndexesSource []string = []string{
	"0 + 0\n",
	"    anotherLine\n",
	"(final . line)",
}

const testGetNextIndexesPath string = "./test/get_next_indexes.yw"

var testGetNextIndexesInputs []Input = []Input{
	// 1 | 0 + 0
	//   ^
	Input{
		lineNumber: 1,
		charNumber: 0,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 1 | 0 + 0
	//     ^
	Input{
		lineNumber: 1,
		charNumber: 1,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 1 | 0 + 0
	//       ^
	Input{
		lineNumber: 1,
		charNumber: 3,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 1 | 0 + 0
	//         ^
	Input{
		lineNumber: 1,
		charNumber: 5,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 1 | 0 + 0
	//          ^(end of line)
	Input{
		lineNumber: 1,
		charNumber: 6,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 2 |     anotherLine
	//         ^
	Input{
		lineNumber: 2,
		charNumber: 5,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 2 |     anotherLine
	//                    ^(end of line)
	Input{
		lineNumber: 2,
		charNumber: 18,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 3 | (final . line)
	//     ^
	Input{
		lineNumber: 3,
		charNumber: 1,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 3 | (final . line)
	//      ^
	Input{
		lineNumber: 3,
		charNumber: 2,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 3 | (final . line)
	//            ^
	Input{
		lineNumber: 3,
		charNumber: 8,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 3 | (final . line)
	//              ^
	Input{
		lineNumber: 3,
		charNumber: 10,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 3 | (final . line)
	//                  ^
	Input{
		lineNumber: 3,
		charNumber: 14,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// 3 | (final . line)
	//                  ^(end of file)
	Input{
		lineNumber: 3,
		charNumber: 14,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
	// calling the function at eof should not change anything no matter how many
	// times it's called
	// 3 | (final . line)
	//                  ^(end of file)
	Input{
		lineNumber: 3,
		charNumber: 14,
		path:       testGetNextIndexesPath,
		source:     testGetNextIndexesSource,
	},
}
/*
func TestGetNextIndexes(t *testing.T) {
	actual, e := Init("./test/get_next_indexes.yw")
	if e != nil {
		fmt.Printf("e.Error(): %v\n", e.Error())
		t.FailNow()
	}

	for _, expected := range testGetNextIndexesInputs {
		line, char, _ := actual.getNextIndexes(); 
		if line != expected.lineNumber || char != expected.charNumber {
			fmt.Fprintf(os.Stderr, 
				"expected: line=%d, char=%d\nactual: line=%d, char=%d\n",
				expected.lineNumber, expected.charNumber,
				line, char)
			t.FailNow()
		}
	}
}*/
