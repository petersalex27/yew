package inf

type QualificationType byte

const (
	NameQualified QualificationType = iota
	NotQualified
	FullyQualified
)