package utils

import (
	"math/rand"
)

const (
	// Uint32Max is the maximum value for uint32
	Uint32Max = ^uint32(0)
)

// ID generates a random uint32 ID
func ID() uint32 {
	return rand.Uint32()
}
