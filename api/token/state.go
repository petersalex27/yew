package token

import (
	"sync"
)

var (
	mu sync.Mutex
	inReplMode bool = false
)

func SetReplMode(truthy bool) {
	mu.Lock()
	defer mu.Unlock()
	inReplMode = truthy
}

func InReplMode() bool {
	mu.Lock()
	defer mu.Unlock()
	return inReplMode
}