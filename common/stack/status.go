package stack

type StackStatus uint

const (
	Ok StackStatus = iota
	Empty
	Overflow
	IllegalOperation
	IllegalReturn
)

func (stat StackStatus) NotOk() bool {
	return stat != Ok
}

func (stat StackStatus) IsOk() bool {
	return stat == Ok
}

func (stat StackStatus) IsEmpty() bool {
	return stat == Empty
}

func (stat StackStatus) IsOverflow() bool {
	return stat == Overflow
}

func (stat StackStatus) IsIllegalOperation() bool {
	return stat == IllegalOperation
}

func (stat StackStatus) Is(stat2 StackStatus) bool {
	return stat == stat2
}

func (stat StackStatus) String() string {
	switch stat {
	case Ok:
		return "Ok"
	case Empty:
		return "Empty"
	case Overflow:
		return "Overflow"
	case IllegalOperation:
		return "IllegalOperation"
	default:
		return "StatusUndefined"
	}
}