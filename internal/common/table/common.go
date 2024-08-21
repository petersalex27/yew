package table

import "fmt"

type internalPair [T fmt.Stringer, U any]struct{Key T; Value U}

type internalTable [T fmt.Stringer, U any]map[string]internalPair[T, U]