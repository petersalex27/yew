// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package types

type Constant struct {
	C          string
	Start, End int
}

func (c Constant) CollectVariables(m map[string]Variable) map[string]Variable { return m }

func (Constant) TypeClassification() typeClassification { return constantClass }

func (c Constant) GetKind() (Term, Type) {
	panic("cannot get the kind of a constant term")
}

func (c *Constant) Parse(s strPos) bool {
	start, end := s.Pos()
	*c = Constant{s.String(), start, end} // TODO: set start and end
	return true
}

func (c Constant) Pos() (start, end int) {
	return c.Start, c.End
}

// value of sort is known, it's a constant
func (Constant) Known() bool { return true }

func (Constant) Locate(Variable) bool { return false }

func (c Constant) Substitute(dest *Term, _ Variable, _ Term) {
	*dest = c
}

func (c Constant) String() string {
	return c.C
}

func (Constant) _identifier_() {}
