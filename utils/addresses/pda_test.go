package addresses

import (
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
)

func TestGetPda_Success(t *testing.T) {
	programID := domains.MustAddressFromString("HNarfxC3kYMMhFkxUFeYb8wHVdPzY5t9pupqW5fL2meM")
	seeds := [][]byte{
		[]byte("escrow"),
		[]byte("test-maker"),
		[]byte("test-hash"),
	}

	pda, err := GetPda(programID, seeds)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if pda == nil {
		t.Fatal("Expected non-nil PDA")
	}
	if len(pda.ToBuffer()) != 32 {
		t.Errorf("Expected 32-byte address, got %d bytes", len(pda.ToBuffer()))
	}
}

func TestGetPda_Deterministic(t *testing.T) {
	programID := domains.MustAddressFromString("HNarfxC3kYMMhFkxUFeYb8wHVdPzY5t9pupqW5fL2meM")
	seeds := [][]byte{
		[]byte("escrow"),
		[]byte("test-maker"),
		[]byte("test-hash"),
	}

	pda1, err1 := GetPda(programID, seeds)
	if err1 != nil {
		t.Fatalf("Expected no error, got %v", err1)
	}

	pda2, err2 := GetPda(programID, seeds)
	if err2 != nil {
		t.Fatalf("Expected no error, got %v", err2)
	}

	if !pda1.Equal(pda2) {
		t.Error("Expected same seeds to produce same PDA")
	}
}

func TestGetPda_EmptySeeds(t *testing.T) {
	programID := domains.MustAddressFromString("HNarfxC3kYMMhFkxUFeYb8wHVdPzY5t9pupqW5fL2meM")
	seeds := [][]byte{}

	pda, err := GetPda(programID, seeds)
	if err != nil {
		t.Fatalf("Expected no error with empty seeds, got %v", err)
	}
	if pda == nil {
		t.Fatal("Expected non-nil PDA even with empty seeds")
	}
}
