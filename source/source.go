package source

import err "yew/error"

type Source []string

func (s Source) GetLine(line int) string {
	if line > len(s) || line < 1 {
		err.PrintBug()
		panic("")
	}
	return s[line-1]
}

func (s Source) GetLineSlice(line int, char int) string {
	sourceLine := s.GetLine(line)
	if char > len(sourceLine) || char < 1 {
		err.PrintBug()
		panic("")
	}
	return sourceLine[char-1:]
}

func (s Source) GetLineSliceN(line int, char int, n int) string {
	slice := s.GetLineSlice(line, char)
	if n > len(slice) || n < 0 {
		return slice
	}

	return slice[:n]
}