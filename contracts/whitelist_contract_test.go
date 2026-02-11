package contracts

import (
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
)

func TestDefaultWhitelistContract(t *testing.T) {
	contract := DefaultWhitelistContract()

	if contract == nil {
		t.Fatal("Expected non-nil contract")
	}
	if contract.Address == nil {
		t.Fatal("Expected non-nil address")
	}
	if !contract.Address.Equal(WhitelistContractAddress) {
		t.Error("Expected address to match WhitelistContractAddress")
	}
}

func TestWhitelistContractAddress(t *testing.T) {
	if WhitelistContractAddress == nil {
		t.Fatal("Expected non-nil WhitelistContractAddress")
	}
	expected := domains.MustAddressFromString("5jzZhrzqkbdwp5d3J1XbmaXMRnqeXimM1mDMoGHyvR7S")
	if !WhitelistContractAddress.Equal(expected) {
		t.Error("Expected WhitelistContractAddress to match expected value")
	}
}
