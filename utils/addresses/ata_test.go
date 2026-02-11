package addresses

import (
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
)

func TestGetAta_Success(t *testing.T) {
	walletAddress := domains.MustAddressFromString("11111111111111111111111111111111")
	tokenMintAddress := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	tokenProgramID := domains.TOKEN_PROGRAM_ID

	ata, err := GetAta(walletAddress, tokenMintAddress, tokenProgramID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if ata == nil {
		t.Fatal("Expected non-nil ATA")
	}
	if len(ata.ToBuffer()) != 32 {
		t.Errorf("Expected 32-byte address, got %d bytes", len(ata.ToBuffer()))
	}
}

func TestGetAta_Deterministic(t *testing.T) {
	walletAddress := domains.MustAddressFromString("11111111111111111111111111111111")
	tokenMintAddress := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	tokenProgramID := domains.TOKEN_PROGRAM_ID

	ata1, err1 := GetAta(walletAddress, tokenMintAddress, tokenProgramID)
	if err1 != nil {
		t.Fatalf("Expected no error, got %v", err1)
	}

	ata2, err2 := GetAta(walletAddress, tokenMintAddress, tokenProgramID)
	if err2 != nil {
		t.Fatalf("Expected no error, got %v", err2)
	}

	if !ata1.Equal(ata2) {
		t.Error("Expected same inputs to produce same ATA")
	}
}

func TestGetAta_DifferentTokens(t *testing.T) {
	walletAddress := domains.MustAddressFromString("11111111111111111111111111111111")
	tokenMint1 := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	tokenMint2 := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	tokenProgramID := domains.TOKEN_PROGRAM_ID

	ata1, err1 := GetAta(walletAddress, tokenMint1, tokenProgramID)
	if err1 != nil {
		t.Fatalf("Expected no error, got %v", err1)
	}

	ata2, err2 := GetAta(walletAddress, tokenMint2, tokenProgramID)
	if err2 != nil {
		t.Fatalf("Expected no error, got %v", err2)
	}

	if ata1.Equal(ata2) {
		t.Error("Expected different tokens to produce different ATAs")
	}
}
