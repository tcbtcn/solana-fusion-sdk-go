package domains

import (
	"math/big"
	"testing"
)

func TestNewBps(t *testing.T) {
	// Test valid BPS
	bps, err := NewBps(big.NewInt(100))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if bps == nil {
		t.Fatal("Expected non-nil Bps")
	}

	// Test invalid BPS (negative)
	_, err = NewBps(big.NewInt(-1))
	if err == nil {
		t.Fatal("Expected error for negative BPS")
	}

	// Test invalid BPS (too large)
	_, err = NewBps(big.NewInt(10001))
	if err == nil {
		t.Fatal("Expected error for BPS > 10000")
	}
}

func TestBpsFromPercent(t *testing.T) {
	bps := BpsFromPercent(1.0, nil) // 1%
	if bps.Value().Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Expected 100 bps, got %s", bps.Value().String())
	}

	bps2 := BpsFromPercent(0.5, nil) // 0.5%
	if bps2.Value().Cmp(big.NewInt(50)) != 0 {
		t.Errorf("Expected 50 bps, got %s", bps2.Value().String())
	}
}

func TestBpsFromFraction(t *testing.T) {
	bps := BpsFromFraction(0.01, nil) // 1%
	if bps.Value().Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Expected 100 bps, got %s", bps.Value().String())
	}
}

func TestBpsEqual(t *testing.T) {
	bps1 := MustBps(big.NewInt(100))
	bps2 := MustBps(big.NewInt(100))
	bps3 := MustBps(big.NewInt(200))

	if !bps1.Equal(bps2) {
		t.Error("Expected BPS values to be equal")
	}
	if bps1.Equal(bps3) {
		t.Error("Expected BPS values to be different")
	}
}

func TestBpsIsZero(t *testing.T) {
	bps1 := MustBps(big.NewInt(0))
	bps2 := MustBps(big.NewInt(100))

	if !bps1.IsZero() {
		t.Error("Expected BPS to be zero")
	}
	if bps2.IsZero() {
		t.Error("Expected BPS to not be zero")
	}
}

func TestBpsToPercent(t *testing.T) {
	bps := MustBps(big.NewInt(100)) // 1%
	percent := bps.ToPercent(nil)
	if percent != 1.0 {
		t.Errorf("Expected 1.0%%, got %f", percent)
	}
}

func TestBpsToFraction(t *testing.T) {
	bps := MustBps(big.NewInt(100)) // 1%
	fraction := bps.ToFraction(nil)
	if fraction != 0.01 {
		t.Errorf("Expected 0.01, got %f", fraction)
	}
}
