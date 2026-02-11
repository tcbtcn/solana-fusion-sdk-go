package fusionorder

import (
	"testing"
)

func TestNewAuctionDetails_Success(t *testing.T) {
	auction, err := NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []AuctionPoint
	}{
		StartTime:       1000000000,
		Duration:        180,
		InitialRateBump: 50000,
		Points:          []AuctionPoint{},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if auction == nil {
		t.Fatal("Expected non-nil auction details")
	}
	if auction.StartTime != 1000000000 {
		t.Errorf("Expected StartTime 1000000000, got %d", auction.StartTime)
	}
	if auction.Duration != 180 {
		t.Errorf("Expected Duration 180, got %d", auction.Duration)
	}
	if auction.InitialRateBump != 50000 {
		t.Errorf("Expected InitialRateBump 50000, got %d", auction.InitialRateBump)
	}
}

func TestNewAuctionDetails_WithPoints(t *testing.T) {
	auction, err := NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []AuctionPoint
	}{
		StartTime:       1000000000,
		Duration:        180,
		InitialRateBump: 50000,
		Points: []AuctionPoint{
			{Delay: 10, Coefficient: 40000},
			{Delay: 20, Coefficient: 30000},
		},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(auction.Points) != 2 {
		t.Errorf("Expected 2 points, got %d", len(auction.Points))
	}
}

func TestNewAuctionDetails_ZeroDuration(t *testing.T) {
	auction, err := NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []AuctionPoint
	}{
		StartTime:       1000000000,
		Duration:        0,
		InitialRateBump: 0,
		Points:          []AuctionPoint{},
	})

	if err != nil {
		t.Fatalf("Expected no error for zero duration, got %v", err)
	}
	if auction.Duration != 0 {
		t.Errorf("Expected Duration 0, got %d", auction.Duration)
	}
}

func TestNoAuction(t *testing.T) {
	auction := NoAuction(1000000000, 180)

	if auction == nil {
		t.Fatal("Expected non-nil auction details")
	}
	if auction.StartTime != 1000000000 {
		t.Errorf("Expected StartTime 1000000000, got %d", auction.StartTime)
	}
	if auction.Duration != 180 {
		t.Errorf("Expected Duration 180, got %d", auction.Duration)
	}
	if auction.InitialRateBump != 0 {
		t.Errorf("Expected InitialRateBump 0, got %d", auction.InitialRateBump)
	}
	if len(auction.Points) != 0 {
		t.Errorf("Expected 0 points, got %d", len(auction.Points))
	}
}

func TestAuctionDetails_FinishTime(t *testing.T) {
	auction := &AuctionDetails{
		StartTime: 1000000000,
		Duration:  180,
	}

	finishTime := auction.FinishTime()
	expected := uint32(1000000180)

	if finishTime != expected {
		t.Errorf("Expected FinishTime %d, got %d", expected, finishTime)
	}
}
