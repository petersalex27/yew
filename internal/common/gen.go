// =================================================================================================
// Alex Peters - February 13, 2024
//
// # Functions for creating generators and methods for those generators
//
// =================================================================================================
package common

import (
	"fmt"
	//"sync"
)

type (
	// give prefix p and suffix s, creates a unique name 
	UniqueNameGenerator struct {
		prefix, suffix string
		counter        uint64
		//mu sync.Mutex
	}
)

// generates a name
func (ung *UniqueNameGenerator) Generate() (name string) {
	// ung.mu.Lock()
	// defer ung.mu.Unlock()
	name = fmt.Sprintf("%s%d%s", ung.prefix, ung.counter, ung.suffix)
	ung.counter++
	return
}

// creates and initializes unique name generator with provided prefix and suffix 
func InitUniqueNameGenerator(prefix, suffix string) *UniqueNameGenerator {
	return &UniqueNameGenerator{prefix: prefix, suffix: suffix, counter: 0}
}

// resets internal counter
func (ung *UniqueNameGenerator) Reset() {
	// ung.mu.Lock()
	// defer ung.mu.Unlock()
	ung.counter = 0
}
