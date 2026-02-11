package amountcalculator

import (
	"math/big"
	"testing"

	fusionorder "github.com/dawitel/solana-fusion-sdk-go/fusion-order"
)

func TestAuctionCalculator_CalcRateBump(t *testing.T) {
	auctionStartTime := uint32(1708448252)

	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       auctionStartTime,
		InitialRateBump: 50000,
		Duration:        120,
		Points:          []fusionorder.AuctionPoint{},
	})

	calculator := FromAuctionData(auctionDetails)

	blockTime := auctionStartTime + 60
	rate := calculator.CalcRateBump(blockTime)
	// At halfway point, rate should be half of initial
	expected := uint16(25000)
	if rate != expected {
		t.Errorf("Expected rate bump %d, got %d", expected, rate)
	}

	// Test at start time
	rate = calculator.CalcRateBump(auctionStartTime)
	expected = uint16(50000)
	if rate != expected {
		t.Errorf("Expected rate bump at start %d, got %d", expected, rate)
	}

	// Test after end time
	rate = calculator.CalcRateBump(auctionStartTime + 121)
	expected = uint16(0)
	if rate != expected {
		t.Errorf("Expected rate bump after end %d, got %d", expected, rate)
	}
}

func TestAuctionCalculator_CalcAuctionTakingAmount(t *testing.T) {
	auctionStartTime := uint32(1708448252)

	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       auctionStartTime,
		InitialRateBump: 50000,
		Duration:        120,
		Points:          []fusionorder.AuctionPoint{},
	})

	calculator := FromAuctionData(auctionDetails)

	blockTime := auctionStartTime + 60
	takingAmount := big.NewInt(1420000000)
	auctionTakingAmount := calculator.CalcAuctionTakingAmount(takingAmount, blockTime)
	// Expected: 1775000000 from rate
	expected := big.NewInt(1775000000)
	if auctionTakingAmount.Cmp(expected) != 0 {
		t.Errorf("Expected auction taking amount %s, got %s", expected.String(), auctionTakingAmount.String())
	}
}

func TestAuctionCalculator_CalcRateBump_WithPoints(t *testing.T) {
	auctionStartTime := uint32(1708448252)

	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       auctionStartTime,
		InitialRateBump: 50000,
		Duration:        120,
		Points: []fusionorder.AuctionPoint{
			{Delay: 30, Coefficient: 40000},
			{Delay: 60, Coefficient: 30000},
		},
	})

	calculator := FromAuctionData(auctionDetails)

	rate := calculator.CalcRateBump(auctionStartTime + 15)
	if rate == 0 {
		t.Error("Expected non-zero rate bump")
	}

	rate = calculator.CalcRateBump(auctionStartTime + 45)
	if rate == 0 {
		t.Error("Expected non-zero rate bump at point")
	}

	rate = calculator.CalcRateBump(auctionStartTime + 90)
	if rate == 0 {
		t.Error("Expected non-zero rate bump between points")
	}
}

func TestAuctionCalculator_CalcRateBump_BeforeStart(t *testing.T) {
	auctionStartTime := uint32(1708448252)

	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       auctionStartTime,
		InitialRateBump: 50000,
		Duration:        120,
		Points:          []fusionorder.AuctionPoint{},
	})

	calculator := FromAuctionData(auctionDetails)

	rate := calculator.CalcRateBump(auctionStartTime - 1)
	expected := uint16(50000)
	if rate != expected {
		t.Errorf("Expected rate bump %d before start, got %d", expected, rate)
	}
}

func TestAuctionCalculator_CalcRateBump_AtBoundaries(t *testing.T) {
	auctionStartTime := uint32(1708448252)

	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       auctionStartTime,
		InitialRateBump: 50000,
		Duration:        120,
		Points:          []fusionorder.AuctionPoint{},
	})

	calculator := FromAuctionData(auctionDetails)

	rate := calculator.CalcRateBump(auctionStartTime)
	expected := uint16(50000)
	if rate != expected {
		t.Errorf("Expected rate bump %d at start, got %d", expected, rate)
	}

	finishTime := calculator.FinishTime()
	rate = calculator.CalcRateBump(finishTime)
	expected = uint16(0)
	if rate != expected {
		t.Errorf("Expected rate bump %d at finish, got %d", expected, rate)
	}
}
