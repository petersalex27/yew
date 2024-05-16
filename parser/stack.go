package parser

type termInfo struct {
	// binding power
	bp int8
	// associates right?
	rAssoc bool
	// arity
	arity uint
	// infixed
	infixed bool
}

func (i termInfo) Arity() uint { return i.arity }

func (i termInfo) AssociatesRight() bool { return i.rAssoc }

func (i termInfo) Bp() int8 { return i.bp }

func (i termInfo) toNonApplicable() termInfo {
	return termInfo{0, false, 0, false}
}

func (i termInfo) decrementArity() (termInfo, bool) {
	if i.arity > 0 && i.infixed {
		i.infixed = false
	}
	
	if i.arity == 0 {
		return termInfo{}, false
	} else if i.arity == 1 {
		return i.toNonApplicable(), true // make into non-applicable term
	}
	i.arity--
	return i, true
}