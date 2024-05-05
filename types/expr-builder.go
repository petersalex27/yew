package types

import "fmt"

type (
	CaseArm struct {
		Pattern    Term
		Expression Term
	}

	Case struct {
		Scrutinee     Term
		ScrutineeType Type
		Arms          []CaseArm
	}
)

func Scrutinize(term Term, arms []CaseArm) Case {
	return Case{Scrutinee: term, Arms: arms}
}

func MakeCaseArm(pattern, expression Term) CaseArm {
	return CaseArm{Pattern: pattern, Expression: expression}
}

func (c *Case) AppendArm(pattern, expression Term) {
	c.Arms = append(c.Arms, MakeCaseArm(pattern, expression))
}

func (arm *CaseArm) substitute(u Variable, s Term) {
	arm.Expression.Substitute(&arm.Expression, u, s)
}

func (arm CaseArm) String() string {
	return fmt.Sprintf("%s => %s", arm.Pattern, arm.Expression)
}

func (c *Case) Substitute(dest *Term, u Variable, s Term) {
	c.Scrutinee.Substitute(&c.Scrutinee, u, s)
	for _, arm := range c.Arms {
		arm.substitute(u, s)
	}
}

func (c *Case) String() string {
	return fmt.Sprintf("case %s of {%s}", c.Scrutinee.String(), joinStringed(c.Arms, "; "))
}
