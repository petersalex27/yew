// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package errors

type Error struct {
	// what led to the error? what was teh code trying to do when it failed?
	Context string
	// what exactly failed?
	ErrorMessage
	// what needs to be done in order to overcome the error
	Mitigation string
}