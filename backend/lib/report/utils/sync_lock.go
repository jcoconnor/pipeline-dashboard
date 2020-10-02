package utils

import (
	"sync"
)

// Locking method to ensure updates to the trains files are
// synchronised with the reader.

var (
    mutex sync.Mutex
)

func Lock() {
    mutex.Lock()
}

func Unlock() {
    mutex.Unlock()
}
