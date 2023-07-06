package info

import (
	"strconv"
	util "yew/utils"
)

// for things which can be located
type Locatable interface {
	GetLocation() Location
}

type Location interface {
	util.Stringable
	GetLine() int
	GetChar() int
	GetPath() string
}

type Loc struct {
	line int
	char int
}

func (loc Loc) GetLine() int { return loc.line }

func (loc Loc) GetChar() int { return loc.char }

func (loc Loc) GetPath() string { return "" }

func (loc Loc) ToString() string {
	return strconv.Itoa(loc.line) + ":" + strconv.Itoa(loc.char)
}

func DefaultLoc() Loc { return Loc{0, 0} }

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

func MakeLocation(line int, char int) Loc {
	return Loc{line: line, char: char}
}

func (p pathLocation) toNoPath() noPathLocation {
	return noPathLocation{line: p.line, char: p.char}
}

func (p pathLocation) ToString() string {
	return p.path + ":" + p.toNoPath().ToString()
}

func (p noPathLocation) ToString() string {
	return strconv.Itoa(p.line) + ":" + strconv.Itoa(p.char)
}

func (p Path) MakeLocation(line int, char int) Location {
	return pathLocation{string(p), line, char}
}

func (loc pathLocation) GetLine() int    { return loc.line }
func (loc pathLocation) GetChar() int    { return loc.char }
func (loc pathLocation) GetPath() string { return loc.path }

func (loc noPathLocation) GetLine() int    { return loc.line }
func (loc noPathLocation) GetChar() int    { return loc.char }
func (loc noPathLocation) GetPath() string { return "" }
