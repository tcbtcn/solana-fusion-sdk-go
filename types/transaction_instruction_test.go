package types

import (
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
)

func TestNewTransactionInstruction(t *testing.T) {
	programID := domains.MustAddressFromString("HNarfxC3kYMMhFkxUFeYb8wHVdPzY5t9pupqW5fL2meM")
	accounts := []AccountMeta{
		{
			Pubkey:     domains.MustAddressFromString("11111111111111111111111111111111"),
			IsSigner:   true,
			IsWritable: true,
		},
	}
	data := []byte{1, 2, 3, 4}

	ix := NewTransactionInstruction(programID, accounts, data)

	if ix == nil {
		t.Fatal("Expected non-nil instruction")
	}
	if !ix.ProgramID.Equal(programID) {
		t.Error("Expected ProgramID to match")
	}
	if len(ix.Accounts) != 1 {
		t.Errorf("Expected 1 account, got %d", len(ix.Accounts))
	}
	if len(ix.Data) != 4 {
		t.Errorf("Expected 4 bytes of data, got %d", len(ix.Data))
	}
}

func TestNewTransactionInstruction_EmptyAccounts(t *testing.T) {
	programID := domains.MustAddressFromString("HNarfxC3kYMMhFkxUFeYb8wHVdPzY5t9pupqW5fL2meM")
	accounts := []AccountMeta{}
	data := []byte{}

	ix := NewTransactionInstruction(programID, accounts, data)

	if ix == nil {
		t.Fatal("Expected non-nil instruction")
	}
	if len(ix.Accounts) != 0 {
		t.Errorf("Expected 0 accounts, got %d", len(ix.Accounts))
	}
	if len(ix.Data) != 0 {
		t.Errorf("Expected 0 bytes of data, got %d", len(ix.Data))
	}
}

func TestAccountMeta(t *testing.T) {
	pubkey := domains.MustAddressFromString("11111111111111111111111111111111")
	meta := AccountMeta{
		Pubkey:     pubkey,
		IsSigner:   true,
		IsWritable: false,
	}

	if !meta.Pubkey.Equal(pubkey) {
		t.Error("Expected Pubkey to match")
	}
	if !meta.IsSigner {
		t.Error("Expected IsSigner to be true")
	}
	if meta.IsWritable {
		t.Error("Expected IsWritable to be false")
	}
}
