// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package types

type Constant string

func (c *Constant) Parse(s string) bool {
	*c = Constant(s)
	return true
}

// value of sort is known, it's a constant
func (Constant) Known() bool { return true }

func (Constant) Locate(Variable) bool { return false }

func (c Constant) Substitute(*Term, Variable, Term) {}

func (c Constant) String() string {
	return string(c)
}

func (Constant) TypeClassification() typeClassification {
	return constantClass
}

func (Constant) _identifier_() {}
