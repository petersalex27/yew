package common

import "io"

type StringReader struct {
	i int // current index in byte buffer
	length int // length of field `str`
	str []byte // string being read
}

// creates new string reader for string `s`
func NewStringReader(s string) *StringReader {
	sr := new(StringReader)
	sr.str = make([]byte, len(s))
	sr.length = copy(sr.str, s) // set length to actual number of bytes copied
	return sr
}

// true when string reader is at end of input
func (sr *StringReader) Eof() bool {
	return sr.length <= sr.i
}

func (sr *StringReader) Read(p []byte) (n int, err error) {
	if sr.Eof() { // end of input
		return 0, io.EOF
	}

	// try to read len(b) bytes from sr.str
	n = copy(p, sr.str[sr.i:])
	// update string reader index
	sr.i = Min(sr.i + n, sr.length)
	return
}