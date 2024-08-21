package holes

import "strings"

type Hole struct {
	scope []interface{ String() string }
	name  interface{ String() string }
}

// creates a string of the form ('|' not included):
//
//	   var : Type
//	   x : Just (Just Thing)
//	------------------------
//	?someHole : SomeType Ty
func (hole Hole) String() string {
	var b strings.Builder

	name := hole.name.String()
	max := len(name)

	for _, v := range hole.scope {
		res := "   " + v.String()
		if len(res) > max {
			max = len(res)
		}
		b.WriteString(res + "\n") // write newline here so it isn't counted in length
	}

	line := strings.Repeat("-", max) + "\n"
	b.WriteString(line)
	b.WriteString(name + "\n") // write newline here so it isn't counted in initial length
	return b.String()
}
