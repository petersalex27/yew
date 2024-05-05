package types

type Sort interface {
	Type
	Known() bool
}