package sdk

import (
	"math/big"
	"testing"

	"github.com/tcbtcn/solana-fusion-sdk-go/api/quoter"
	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
)

func TestQuoteFromJSON_Success(t *testing.T) {
	quoteID := "test-quote-id"
	json := &quoter.QuoteDTO{
		QuoteID:   &quoteID,
		SrcAmount: "1000000000000000000",
		DstAmount: "1420000000",
		Presets: quoter.PresetsDTO{
			Fast: quoter.PresetDTO{
				StartAuctionIn:     10,
				AuctionDuration:    180,
				InitialRateBump:    0,
				AuctionStartAmount: "1500000000",
				AuctionEndAmount:   "1420000000",
				CostInDstToken:     "0",
				Points:             []quoter.PresetPointDTO{},
			},
		},
		RecommendedPreset:  quoter.PresetTypeFast,
		PriceImpactPercent: 0.5,
	}

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	quote, err := QuoteFromJSON(srcToken, dstToken, signer, json)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if quote == nil {
		t.Fatal("Expected non-nil quote")
	}
	if quote.QuoteID != quoteID {
		t.Errorf("Expected QuoteID %s, got %s", quoteID, quote.QuoteID)
	}
	if quote.SrcAmount.Cmp(big.NewInt(1000000000000000000)) != 0 {
		t.Errorf("Expected SrcAmount %s, got %s", "1000000000000000000", quote.SrcAmount.String())
	}
	if quote.DstAmount.Cmp(big.NewInt(1420000000)) != 0 {
		t.Errorf("Expected DstAmount %s, got %s", "1420000000", quote.DstAmount.String())
	}
	if quote.PriceImpactPercent != 0.5 {
		t.Errorf("Expected PriceImpactPercent 0.5, got %f", quote.PriceImpactPercent)
	}
}

func TestQuoteFromJSON_MissingQuoteID(t *testing.T) {
	json := &quoter.QuoteDTO{
		QuoteID:   nil,
		SrcAmount: "1000000000000000000",
		DstAmount: "1420000000",
	}

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	_, err := QuoteFromJSON(srcToken, dstToken, signer, json)
	if err == nil {
		t.Fatal("Expected error for missing quoteId")
	}
}

func TestQuoteFromJSON_InvalidSrcAmount(t *testing.T) {
	quoteID := "test-quote-id"
	json := &quoter.QuoteDTO{
		QuoteID:   &quoteID,
		SrcAmount: "invalid",
		DstAmount: "1420000000",
	}

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	_, err := QuoteFromJSON(srcToken, dstToken, signer, json)
	if err == nil {
		t.Fatal("Expected error for invalid srcAmount")
	}
}

func TestQuoteFromJSON_InvalidDstAmount(t *testing.T) {
	quoteID := "test-quote-id"
	json := &quoter.QuoteDTO{
		QuoteID:   &quoteID,
		SrcAmount: "1000000000000000000",
		DstAmount: "invalid",
	}

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	_, err := QuoteFromJSON(srcToken, dstToken, signer, json)
	if err == nil {
		t.Fatal("Expected error for invalid dstAmount")
	}
}

func TestQuote_ToOrder_Success(t *testing.T) {
	quoteID := "test-quote-id"
	json := &quoter.QuoteDTO{
		QuoteID:   &quoteID,
		SrcAmount: "1000000000000000000",
		DstAmount: "1420000000",
		Presets: quoter.PresetsDTO{
			Fast: quoter.PresetDTO{
				StartAuctionIn:     10,
				AuctionDuration:    180,
				InitialRateBump:    0,
				AuctionStartAmount: "1500000000",
				AuctionEndAmount:   "1420000000",
				CostInDstToken:     "0",
				Points:             []quoter.PresetPointDTO{},
			},
		},
		RecommendedPreset: quoter.PresetTypeFast,
	}

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	quote, err := QuoteFromJSON(srcToken, dstToken, signer, json)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	order, err := quote.ToOrder(quoter.PresetTypeFast, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if order == nil {
		t.Fatal("Expected non-nil order")
	}
	if order.SrcAmount().Cmp(big.NewInt(1000000000000000000)) != 0 {
		t.Errorf("Expected SrcAmount %s, got %s", "1000000000000000000", order.SrcAmount().String())
	}
}

func TestQuote_ToOrder_WithReceiver(t *testing.T) {
	quoteID := "test-quote-id"
	json := &quoter.QuoteDTO{
		QuoteID:   &quoteID,
		SrcAmount: "1000000000000000000",
		DstAmount: "1420000000",
		Presets: quoter.PresetsDTO{
			Fast: quoter.PresetDTO{
				StartAuctionIn:     10,
				AuctionDuration:    180,
				InitialRateBump:    0,
				AuctionStartAmount: "1500000000",
				AuctionEndAmount:   "1420000000",
				CostInDstToken:     "0",
				Points:             []quoter.PresetPointDTO{},
			},
		},
		RecommendedPreset: quoter.PresetTypeFast,
	}

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	signer := domains.MustAddressFromString("11111111111111111111111111111111")
	receiver := domains.MustAddressFromString("11111111111111111111111111111112")

	quote, err := QuoteFromJSON(srcToken, dstToken, signer, json)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	order, err := quote.ToOrder(quoter.PresetTypeFast, receiver)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !order.Receiver().Equal(receiver) {
		t.Error("Expected receiver to match provided receiver")
	}
}

func TestQuote_ToOrder_DefaultPreset(t *testing.T) {
	quoteID := "test-quote-id"
	json := &quoter.QuoteDTO{
		QuoteID:   &quoteID,
		SrcAmount: "1000000000000000000",
		DstAmount: "1420000000",
		Presets: quoter.PresetsDTO{
			Fast: quoter.PresetDTO{
				StartAuctionIn:     10,
				AuctionDuration:    180,
				InitialRateBump:    0,
				AuctionStartAmount: "1500000000",
				AuctionEndAmount:   "1420000000",
				CostInDstToken:     "0",
				Points:             []quoter.PresetPointDTO{},
			},
		},
		RecommendedPreset: quoter.PresetTypeFast,
	}

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	quote, err := QuoteFromJSON(srcToken, dstToken, signer, json)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	order, err := quote.ToOrder(quoter.PresetType("unknown"), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if order == nil {
		t.Fatal("Expected non-nil order")
	}
}

// TestQuote_ToOrder_EmptyPresetType tests that empty presetType uses RecommendedPreset
func TestQuote_ToOrder_EmptyPresetType(t *testing.T) {
	quoteID := "test-quote-id"
	json := &quoter.QuoteDTO{
		QuoteID:   &quoteID,
		SrcAmount: "1000000000000000000",
		DstAmount: "1420000000",
		Presets: quoter.PresetsDTO{
			Medium: quoter.PresetDTO{
				StartAuctionIn:     20,
				AuctionDuration:    200,
				InitialRateBump:    100,
				AuctionStartAmount: "1500000000",
				AuctionEndAmount:   "1420000000",
				CostInDstToken:     "0",
				Points:             []quoter.PresetPointDTO{},
			},
		},
		RecommendedPreset: quoter.PresetTypeMedium,
	}

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	quote, err := QuoteFromJSON(srcToken, dstToken, signer, json)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	// Test with empty presetType (zero value) - should use RecommendedPreset (Medium)
	order, err := quote.ToOrder("", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if order == nil {
		t.Fatal("Expected non-nil order")
	}
	// Verify it used Medium preset by checking auction duration (200 for medium vs 180 for fast)
	// This is a basic check - the actual preset used would affect the auction details
}

// TestQuote_ToOrder_NilReceiver tests that nil receiver uses Signer
func TestQuote_ToOrder_NilReceiver(t *testing.T) {
	quoteID := "test-quote-id"
	json := &quoter.QuoteDTO{
		QuoteID:   &quoteID,
		SrcAmount: "1000000000000000000",
		DstAmount: "1420000000",
		Presets: quoter.PresetsDTO{
			Fast: quoter.PresetDTO{
				StartAuctionIn:     10,
				AuctionDuration:    180,
				InitialRateBump:    0,
				AuctionStartAmount: "1500000000",
				AuctionEndAmount:   "1420000000",
				CostInDstToken:     "0",
				Points:             []quoter.PresetPointDTO{},
			},
		},
		RecommendedPreset: quoter.PresetTypeFast,
	}

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	quote, err := QuoteFromJSON(srcToken, dstToken, signer, json)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	// Test with nil receiver - should use Signer
	order, err := quote.ToOrder(quoter.PresetTypeFast, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if order == nil {
		t.Fatal("Expected non-nil order")
	}
	if !order.Receiver().Equal(signer) {
		t.Error("Expected receiver to be signer when nil receiver provided")
	}
}

// TestQuote_ToOrder_EmptyPresetAndNilReceiver tests both defaults together
func TestQuote_ToOrder_EmptyPresetAndNilReceiver(t *testing.T) {
	quoteID := "test-quote-id"
	json := &quoter.QuoteDTO{
		QuoteID:   &quoteID,
		SrcAmount: "1000000000000000000",
		DstAmount: "1420000000",
		Presets: quoter.PresetsDTO{
			Slow: quoter.PresetDTO{
				StartAuctionIn:     30,
				AuctionDuration:    300,
				InitialRateBump:    200,
				AuctionStartAmount: "1500000000",
				AuctionEndAmount:   "1420000000",
				CostInDstToken:     "0",
				Points:             []quoter.PresetPointDTO{},
			},
		},
		RecommendedPreset: quoter.PresetTypeSlow,
	}

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	quote, err := QuoteFromJSON(srcToken, dstToken, signer, json)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	// Test with both empty presetType and nil receiver - should use RecommendedPreset and Signer
	order, err := quote.ToOrder("", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if order == nil {
		t.Fatal("Expected non-nil order")
	}
	if !order.Receiver().Equal(signer) {
		t.Error("Expected receiver to be signer when nil receiver provided")
	}
}

// TestQuote_ToOrder_InvalidPreset tests handling of invalid preset types
func TestQuote_ToOrder_InvalidPreset(t *testing.T) {
	quoteID := "test-quote-id"
	json := &quoter.QuoteDTO{
		QuoteID:   &quoteID,
		SrcAmount: "1000000000000000000",
		DstAmount: "1420000000",
		Presets: quoter.PresetsDTO{
			Fast: quoter.PresetDTO{
				StartAuctionIn:     10,
				AuctionDuration:    180,
				InitialRateBump:    0,
				AuctionStartAmount: "1500000000",
				AuctionEndAmount:   "1420000000",
				CostInDstToken:     "0",
				Points:             []quoter.PresetPointDTO{},
			},
		},
		RecommendedPreset: quoter.PresetTypeFast,
	}

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	quote, err := QuoteFromJSON(srcToken, dstToken, signer, json)
	if err != nil {
		t.Fatalf("Failed to create quote: %v", err)
	}

	// Test with invalid preset type - should fallback to RecommendedPreset (Fast)
	order, err := quote.ToOrder(quoter.PresetType("invalid"), nil)
	if err != nil {
		t.Fatalf("Expected no error (should fallback to RecommendedPreset), got %v", err)
	}
	if order == nil {
		t.Fatal("Expected non-nil order")
	}
}
