package time

import (
	"testing"
)

func TestNow_ReturnsValidUint32(t *testing.T) {
	now := Now()
	if now == 0 {
		t.Error("Expected non-zero time")
	}
}

func TestNow_ReasonableValue(t *testing.T) {
	now := Now()
	if now < 1000000000 {
		t.Errorf("Expected time >= 1000000000 (year 2001), got %d", now)
	}
	if now > 4000000000 {
		t.Errorf("Expected time < 4000000000 (year 2096), got %d", now)
	}
}

func TestNow_Increasing(t *testing.T) {
	now1 := Now()
	now2 := Now()
	if now2 < now1 {
		t.Errorf("Expected time to be non-decreasing, got %d < %d", now2, now1)
	}
}
