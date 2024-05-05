package inf

import (
	"github.com/petersalex27/yew-packages/expr"
	"github.com/petersalex27/yew-packages/nameable"
	"github.com/petersalex27/yew-packages/types"
)

type errorReport[N nameable.Nameable] struct {
	During        string
	Status        Status
	TermsInvolved []TypeJudgment[N]
	Names         []expr.Const[N]
	TypesInvolved []types.Type[N]
}

// creates an errorReport for a failed rule
func makeReport[N nameable.Nameable](duringRule string, status Status, withTerms ...TypeJudgment[N]) errorReport[N] {
	return errorReport[N]{duringRule, status, withTerms, nil, nil}
}

// creates an errorReport for a failed context lookup
func makeNameReport[N nameable.Nameable](duringRule string, status Status, withNames ...expr.Const[N]) errorReport[N] {
	return errorReport[N]{duringRule, status, nil, withNames, nil}
}

func makeTypeReport[N nameable.Nameable](during string, status Status, withTypes ...types.Type[N]) errorReport[N] {
	return errorReport[N]{during, status, nil, nil, withTypes}
}
