package inf

import "fmt"

type Status byte

const (
	// default status
	Ok Status = iota
	// constants at same positions in unification did not match
	ConstantMismatch
	// kind-constants at same positions in unification did not match
	KindConstantMismatch
	// application monotype kind did not have same number of type params as other
	// monotype being unified
	ParamLengthMismatch
	// dependent instance did not have same number of type indexes as other
	// monotype being unified
	IndexLengthMismatch
	// data did not have same number of members as other expression being unified
	MemsLengthMismatch
	// unification of variable and monotype failed because variable occurred w/in
	// monotype
	OccursCheckFailed
	// failed to find given name in available context
	NameNotInContext
	// length of judgments slice passed to second part of `Rec` did not match
	// length of names slice passed to the first part of `Rec`
	RecArgsLengthMismatch
	// type already exists in constructor table
	TypeRedef
	// constructor already exists for given type in constructor table
	ConstructorRedef
	// type is not defined but tried to the type anyways
	TypeNotDefined
	// tried to export undefined type
	UndefinedType
	// tried to export undefined constructor
	UndefinedConstructor
	// tried to export undefined function
	UndefinedFunction
	// tried to export an ambiguous function (has multiple definitions)
	AmbiguousFunction
	// illegal name shadowing 
	IllegalShadow
	// unification of variables succeeded, so signals that there is nothing left 
	// to unify
	skipUnify
)

func (stat Status) String() string {
	switch stat {
	case Ok:
		return "Ok"
	case ConstantMismatch:
		return "ConstantMismatch"
	case KindConstantMismatch:
		return "KindConstantMismatch"
	case ParamLengthMismatch:
		return "ParamLengthMismatch"
	case IndexLengthMismatch:
		return "IndexLengthMismatch"
	case MemsLengthMismatch:
		return "MemsLengthMismatch"
	case OccursCheckFailed:
		return "OccursCheckFailed"
	case NameNotInContext:
		return "NameNotInContext"
	case RecArgsLengthMismatch:
		return "RecArgsLengthMismatch"
	case skipUnify:
		return "skipUnify"
	default:
		return fmt.Sprintf("Status(%d)", stat)
	}
}

func (stat Status) IsOk() bool {
	return stat == Ok
}

func (stat Status) NotOk() bool {
	return stat != Ok
}

func (stat Status) Is(stat2 Status) bool {
	return stat == stat2
}