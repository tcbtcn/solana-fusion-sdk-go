package sdk

import (
	"math/big"
	"testing"

	"github.com/tcbtcn/solana-fusion-sdk-go/api/quoter"
)

func TestPresetFromJSON_Success(t *testing.T) {
	json := quoter.PresetDTO{
		StartAuctionIn:     10,
		AuctionDuration:    180,
		InitialRateBump:    50000,
		AuctionStartAmount: "1500000000",
		AuctionEndAmount:   "1420000000",
		CostInDstToken:     "100",
		Points: []quoter.PresetPointDTO{
			{Delay: 10, Coefficient: 40000},
			{Delay: 20, Coefficient: 30000},
		},
	}

	preset := PresetFromJSON(json)

	if preset == nil {
		t.Fatal("Expected non-nil preset")
	}
	if preset.StartAuctionIn != 10 {
		t.Errorf("Expected StartAuctionIn 10, got %d", preset.StartAuctionIn)
	}
	if preset.AuctionDuration != 180 {
		t.Errorf("Expected AuctionDuration 180, got %d", preset.AuctionDuration)
	}
	if preset.InitialRateBump != 50000 {
		t.Errorf("Expected InitialRateBump 50000, got %d", preset.InitialRateBump)
	}
	if preset.AuctionStartAmount.Cmp(big.NewInt(1500000000)) != 0 {
		t.Errorf("Expected AuctionStartAmount 1500000000, got %s", preset.AuctionStartAmount.String())
	}
	if preset.AuctionEndAmount.Cmp(big.NewInt(1420000000)) != 0 {
		t.Errorf("Expected AuctionEndAmount 1420000000, got %s", preset.AuctionEndAmount.String())
	}
	if preset.CostInDstToken.Cmp(big.NewInt(100)) != 0 {
		t.Errorf("Expected CostInDstToken 100, got %s", preset.CostInDstToken.String())
	}
	if len(preset.Points) != 2 {
		t.Errorf("Expected 2 points, got %d", len(preset.Points))
	}
	if preset.Points[0].Delay != 10 {
		t.Errorf("Expected first point delay 10, got %d", preset.Points[0].Delay)
	}
	if preset.Points[0].Coefficient != 40000 {
		t.Errorf("Expected first point coefficient 40000, got %d", preset.Points[0].Coefficient)
	}
}

func TestPresetFromJSON_EmptyPoints(t *testing.T) {
	json := quoter.PresetDTO{
		StartAuctionIn:     10,
		AuctionDuration:    180,
		InitialRateBump:    0,
		AuctionStartAmount: "1500000000",
		AuctionEndAmount:   "1420000000",
		CostInDstToken:     "0",
		Points:             []quoter.PresetPointDTO{},
	}

	preset := PresetFromJSON(json)

	if preset == nil {
		t.Fatal("Expected non-nil preset")
	}
	if len(preset.Points) != 0 {
		t.Errorf("Expected 0 points, got %d", len(preset.Points))
	}
}

func TestPresetFromJSON_InvalidAmounts(t *testing.T) {
	json := quoter.PresetDTO{
		StartAuctionIn:     10,
		AuctionDuration:    180,
		InitialRateBump:    0,
		AuctionStartAmount: "invalid",
		AuctionEndAmount:   "invalid",
		CostInDstToken:     "invalid",
		Points:             []quoter.PresetPointDTO{},
	}

	preset := PresetFromJSON(json)

	if preset == nil {
		t.Fatal("Expected non-nil preset")
	}
	if preset.AuctionStartAmount != nil {
		t.Error("Expected nil AuctionStartAmount for invalid input")
	}
	if preset.AuctionEndAmount != nil {
		t.Error("Expected nil AuctionEndAmount for invalid input")
	}
	if preset.CostInDstToken != nil {
		t.Error("Expected nil CostInDstToken for invalid input")
	}
}
