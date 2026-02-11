package domains

import (
	"testing"
)

func TestNewAddress(t *testing.T) {
	// Test valid address
	addr, err := NewAddress("11111111111111111111111111111111")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if addr == nil {
		t.Fatal("Expected non-nil address")
	}

	// Test invalid address (wrong length)
	_, err = NewAddress("invalid")
	if err == nil {
		t.Fatal("Expected error for invalid address")
	}

	// Test invalid base58
	_, err = NewAddress("0")
	if err == nil {
		t.Fatal("Expected error for invalid base58")
	}
}

func TestAddressToString(t *testing.T) {
	addr := MustAddressFromString("11111111111111111111111111111111")
	str := addr.ToString()
	if str != "11111111111111111111111111111111" {
		t.Errorf("Expected '11111111111111111111111111111111', got '%s'", str)
	}
}

func TestAddressEqual(t *testing.T) {
	addr1 := MustAddressFromString("11111111111111111111111111111111")
	addr2 := MustAddressFromString("11111111111111111111111111111111")
	addr3 := MustAddressFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")

	if !addr1.Equal(addr2) {
		t.Error("Expected addresses to be equal")
	}
	if addr1.Equal(addr3) {
		t.Error("Expected addresses to be different")
	}
}

func TestAddressIsNative(t *testing.T) {
	addr := MustAddressFromString("SoNative11111111111111111111111111111111111")
	if !addr.IsNative() {
		t.Error("Expected address to be native")
	}

	addr2 := MustAddressFromString("11111111111111111111111111111111")
	if addr2.IsNative() {
		t.Error("Expected address to not be native")
	}
}

func TestAddressFromBuffer(t *testing.T) {
	buf := make([]byte, 32)
	addr := AddressFromBuffer(buf)
	if len(addr.ToBuffer()) != 32 {
		t.Errorf("Expected 32 bytes, got %d", len(addr.ToBuffer()))
	}
}

func TestAddress_NilHandling(t *testing.T) {
	var nilAddr *Address

	// Test ToString with nil
	if nilAddr.ToString() != "" {
		t.Error("Expected empty string for nil address ToString")
	}

	// Test ToBuffer with nil
	if nilAddr.ToBuffer() != nil {
		t.Error("Expected nil for nil address ToBuffer")
	}

	// Test Bytes with nil
	if nilAddr.Bytes() != nil {
		t.Error("Expected nil for nil address Bytes")
	}

	// Test IsNative with nil
	if nilAddr.IsNative() {
		t.Error("Expected false for nil address IsNative")
	}

	// Test Equal with nil
	validAddr := MustAddressFromString("11111111111111111111111111111111")
	if nilAddr.Equal(validAddr) {
		t.Error("Expected false when comparing nil address with valid address")
	}
	if !nilAddr.Equal(nil) {
		t.Error("Expected true when comparing nil address with nil")
	}
}
