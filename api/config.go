package api

type Config interface {
	Set(raw string) error
	Get(key string) any
	All() map[string]any
}
