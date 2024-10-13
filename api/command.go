package api

type (
	Command interface {
		// Execute the command with the given arguments, returning an error if the command fails
		Do(args ...any) error
	}
)