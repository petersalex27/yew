package parser


type termInfo struct {
	// binding power
	bp uint8
	// associates right?
	rAssoc bool
	// arity
	arity uint
}

func (i termInfo) Arity() uint { return i.arity }

func (i termInfo) AssociatesRight() bool { return i.rAssoc }

func (i termInfo) Bp() uint8 { return i.bp }

func (i termInfo) decrementArity() (termInfo, bool) {
	if i.arity == 0 {
		return termInfo{}, false
	} else if i.arity == 1 {
		return termInfo{}, true // make into non-applicable term
	}
	i.arity--
	return i, true
}

func (parser *Parser) top() termElem {
	term, stat := parser.terms.Peek()
	if stat.NotOk() {
		panic("bug: empty stack")
	}
	return term
}
