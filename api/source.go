package api

type (
	Source interface {
		// returns string representation of the source
		String() string
		// returns the path of the source
		Path() string
	}

	// if the instance of SourceCode needs to be modified before being used to create a window,
	// implement this method:
	//
	// 	PrepareForWindowing() SourceCode
	SourceCode interface {
		Source
		// returns line count of the source
		Lines() int
		EndPositions() []int
		Set(Source) SourceCode
	}

	UnpreparedWindowingSourceCode interface {
		Source
		PrepareForWindowing() SourceCode
	}
)
