package module

const (
	Application TermType = iota
	Abstraction
	Variable
	Constant
	Pi
	Universe
)

type (
	ModuleInterface interface {
		AsModule() Module
	}

	TermType byte

	SerializedTerm struct {
		// header

		// term type of the term
		TermTermType TermType
		// term type of the type
		TypeTermType TermType
		// implicit: termLen      [8]byte
		// implicit: typeLen      [8]byte

		// payload

		// serialized term
		Term         []byte
		// serialized type
		Type         []byte
	}

	SerializableTerm interface {
		SerializeTerm() SerializedTerm
	}

	/*
		(name) [len; 8] [name[0], name[1], .., name[len-1]; len]
		(dependencies) [nDeps; 8]
		(dependency 0) [dlen0; 8] [dependencies[0][0], dependencies[0][1], .., dependencies[0][dlen0-1]; dlen0]
		(dependency 1) [dlen1; 8] [dependencies[1][0], dependencies[1][1], .., dependencies[1][dlen1-1]; dlen1]
		...
		(dependency nDeps-1) [dlen(nDeps-1); 8] [dependencies[nDeps-1][0], dependencies[nDeps-1][1], .., dependencies[nDeps-1][dlen(nDeps-1)-1]; dlen(nDeps-1)]
	*/
	Module struct {
		name         string
		dependencies []string
		// terms[n][0] is the Nth term, terms[n][1] is Nth term's type
		terms [][2][]byte
	}
)
