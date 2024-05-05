package types

type TypeConstant string

func (c *TypeConstant) Parse(s string) bool {
	*c = TypeConstant(s)
	return true
}

func (TypeConstant) Locate(Variable) bool { return false }

func (tc TypeConstant) Substitute(*Term, Variable, Term) {}

func (tc TypeConstant) String() string {
	return string(tc)
}

func (TypeConstant) TypeClassification() typeClassification {
	return typeConstantClass
}
