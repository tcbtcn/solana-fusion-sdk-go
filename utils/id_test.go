package utils

import (
	"testing"
)

func TestID_Uniqueness(t *testing.T) {
	ids := make(map[uint32]bool)
	for range 1000 {
		if _, ok := ids[ID()]; ok {
			t.Error("Duplicate ID found")
		}
		ids[ID()] = true
	}
}
