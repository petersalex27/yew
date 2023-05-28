package info

// for things which can be located
type Locatable interface {
	GetLocation() Location
}

type Location interface {
	GetLine() int
	GetChar() int
	GetPath() string
}

type pathLocation struct {
	path string
	line int
	char int
}
type noPathLocation struct {
	line int
	char int
}
type Path string

func MakeLocation(line int, char int) Location {
	return noPathLocation{line: line, char: char}
}

func (p Path) MakeLocation(line int, char int) Location {
	return pathLocation{string(p), line, char}
}

func (loc pathLocation) GetLine() int { return loc.line }
func (loc pathLocation) GetChar() int { return loc.char }
func (loc pathLocation) GetPath() string { return loc.path }

func (loc noPathLocation) GetLine() int { return loc.line }
func (loc noPathLocation) GetChar() int { return loc.char }
func (loc noPathLocation) GetPath() string { return "" }