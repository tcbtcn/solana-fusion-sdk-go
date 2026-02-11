package validation

import (
	"math/big"
	"testing"
)

func TestAssertUInteger_Int(t *testing.T) {
	err := AssertUInteger(100, nil)
	if err != nil {
		t.Errorf("Expected no error for int, got %v", err)
	}
}

func TestAssertUInteger_Int64(t *testing.T) {
	err := AssertUInteger(int64(100), nil)
	if err != nil {
		t.Errorf("Expected no error for int64, got %v", err)
	}
}

func TestAssertUInteger_Uint(t *testing.T) {
	err := AssertUInteger(uint(100), nil)
	if err != nil {
		t.Errorf("Expected no error for uint, got %v", err)
	}
}

func TestAssertUInteger_Uint32(t *testing.T) {
	err := AssertUInteger(uint32(100), nil)
	if err != nil {
		t.Errorf("Expected no error for uint32, got %v", err)
	}
}

func TestAssertUInteger_Uint64(t *testing.T) {
	err := AssertUInteger(uint64(100), nil)
	if err != nil {
		t.Errorf("Expected no error for uint64, got %v", err)
	}
}

func TestAssertUInteger_BigInt(t *testing.T) {
	err := AssertUInteger(big.NewInt(100), nil)
	if err != nil {
		t.Errorf("Expected no error for *big.Int, got %v", err)
	}
}

func TestAssertUInteger_NegativeInt(t *testing.T) {
	err := AssertUInteger(-100, nil)
	if err == nil {
		t.Fatal("Expected error for negative int")
	}
}

func TestAssertUInteger_NegativeInt64(t *testing.T) {
	err := AssertUInteger(int64(-100), nil)
	if err == nil {
		t.Fatal("Expected error for negative int64")
	}
}

func TestAssertUInteger_NegativeBigInt(t *testing.T) {
	err := AssertUInteger(big.NewInt(-100), nil)
	if err == nil {
		t.Fatal("Expected error for negative *big.Int")
	}
}

func TestAssertUInteger_InvalidType(t *testing.T) {
	err := AssertUInteger("invalid", nil)
	if err == nil {
		t.Fatal("Expected error for invalid type")
	}
}

func TestAssertUInteger_Zero(t *testing.T) {
	err := AssertUInteger(0, nil)
	if err != nil {
		t.Errorf("Expected no error for zero, got %v", err)
	}
}

func TestAssertUInteger_WithMax(t *testing.T) {
	max := big.NewInt(1000)
	err := AssertUInteger(500, max)
	if err != nil {
		t.Errorf("Expected no error for value within max, got %v", err)
	}

	err = AssertUInteger(1500, max)
	if err == nil {
		t.Fatal("Expected error for value exceeding max")
	}
}
