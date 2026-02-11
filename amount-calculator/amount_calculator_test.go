package amountcalculator

import (
	"math/big"
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
	fusionorder "github.com/dawitel/solana-fusion-sdk-go/fusion-order"
	"github.com/dawitel/solana-fusion-sdk-go/utils/time"
)

func TestAmountCalculator_GetTotalFee(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	integratorFee := domains.BpsFromPercent(2.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)
	
	now := uint32(time.Now())
	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       now,
		Duration:        120,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	
	feeConfig, _ := fusionorder.NewFeeConfig(nil, nil, protocolFee, integratorFee, surplusShare)
	calculator := NewAmountCalculator(
		FromAuctionData(auctionDetails),
		FromFeeConfig(feeConfig),
	)

	takingAmount := big.NewInt(1000)
	estimatedTakingAmount := big.NewInt(500)

	totalFee := calculator.GetTotalFee(takingAmount, estimatedTakingAmount, now)
	// Expected: ((1000 - 10 - 20) - 500) * 0.5 + 10 + 20 = 265
	expected := big.NewInt(265)
	if totalFee.Cmp(expected) != 0 {
		t.Errorf("Expected total fee %s, got %s", expected.String(), totalFee.String())
	}
}

func TestAmountCalculator_GetIntegratorFee(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	integratorFee := domains.BpsFromPercent(2.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)
	
	now := uint32(time.Now())
	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       now,
		Duration:        120,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	
	feeConfig, _ := fusionorder.NewFeeConfig(nil, nil, protocolFee, integratorFee, surplusShare)
	calculator := NewAmountCalculator(
		FromAuctionData(auctionDetails),
		FromFeeConfig(feeConfig),
	)

	takingAmount := big.NewInt(1000)

	integratorFeeAmount := calculator.GetIntegratorFee(takingAmount, now)
	// Expected: 1000 * 2% = 20
	expected := big.NewInt(20)
	if integratorFeeAmount.Cmp(expected) != 0 {
		t.Errorf("Expected integrator fee %s, got %s", expected.String(), integratorFeeAmount.String())
	}
}

func TestAmountCalculator_GetProtocolFee(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	integratorFee := domains.BpsFromPercent(2.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)
	
	now := uint32(time.Now())
	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       now,
		Duration:        120,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	
	feeConfig, _ := fusionorder.NewFeeConfig(nil, nil, protocolFee, integratorFee, surplusShare)
	calculator := NewAmountCalculator(
		FromAuctionData(auctionDetails),
		FromFeeConfig(feeConfig),
	)

	// Test with no surplus
	takingAmount := big.NewInt(1000)
	estimatedTakingAmount := big.NewInt(1000)

	protocolFeeAmount := calculator.GetProtocolFee(takingAmount, estimatedTakingAmount, now)
	// Expected: 1000 * 1% = 10 (no surplus)
	expected := big.NewInt(10)
	if protocolFeeAmount.Cmp(expected) != 0 {
		t.Errorf("Expected protocol fee %s, got %s", expected.String(), protocolFeeAmount.String())
	}

	// Test with surplus
	estimatedTakingAmount = big.NewInt(500)
	protocolFeeAmount = calculator.GetProtocolFee(takingAmount, estimatedTakingAmount, now)
	// Expected: ((1000 - 10 - 20) - 500) * 0.5 + 10 = 245
	expected = big.NewInt(245)
	if protocolFeeAmount.Cmp(expected) != 0 {
		t.Errorf("Expected protocol fee with surplus %s, got %s", expected.String(), protocolFeeAmount.String())
	}
}

func TestAmountCalculator_GetUserReceiveAmount(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	integratorFee := domains.BpsFromPercent(2.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)
	
	now := uint32(time.Now())
	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       now,
		Duration:        120,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	
	feeConfig, _ := fusionorder.NewFeeConfig(nil, nil, protocolFee, integratorFee, surplusShare)
	calculator := NewAmountCalculator(
		FromAuctionData(auctionDetails),
		FromFeeConfig(feeConfig),
	)

	takingAmount := big.NewInt(1000)
	estimatedTakingAmount := big.NewInt(500)

	userAmount := calculator.GetUserReceiveAmount(takingAmount, estimatedTakingAmount, now)
	expected := big.NewInt(745)
	if userAmount.Cmp(expected) != 0 {
		t.Errorf("Expected user receive amount %s, got %s", expected.String(), userAmount.String())
	}
}

func TestAmountCalculator_GetRequiredTakingAmount(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	integratorFee := domains.BpsFromPercent(2.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)
	
	now := uint32(time.Now())
	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       now,
		Duration:        120,
		InitialRateBump: 50000,
		Points:          []fusionorder.AuctionPoint{},
	})
	
	feeConfig, _ := fusionorder.NewFeeConfig(nil, nil, protocolFee, integratorFee, surplusShare)
	calculator := NewAmountCalculator(
		FromAuctionData(auctionDetails),
		FromFeeConfig(feeConfig),
	)

	takingAmount := big.NewInt(1000)
	requiredAmount := calculator.GetRequiredTakingAmount(takingAmount, now)

	if requiredAmount.Cmp(takingAmount) <= 0 {
		t.Error("Expected required amount to be greater than taking amount due to auction")
	}
}

func TestAmountCalculator_ZeroAmounts(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	integratorFee := domains.BpsFromPercent(2.0, nil)
	
	now := uint32(time.Now())
	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       now,
		Duration:        120,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	
	feeConfig, _ := fusionorder.NewFeeConfig(nil, nil, protocolFee, integratorFee, domains.ZeroBps)
	calculator := NewAmountCalculator(
		FromAuctionData(auctionDetails),
		FromFeeConfig(feeConfig),
	)

	takingAmount := big.NewInt(0)
	estimatedTakingAmount := big.NewInt(0)

	totalFee := calculator.GetTotalFee(takingAmount, estimatedTakingAmount, now)
	if totalFee.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected zero fee for zero amount, got %s", totalFee.String())
	}
}

func TestAmountCalculator_LargeAmounts(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	integratorFee := domains.BpsFromPercent(2.0, nil)
	
	now := uint32(time.Now())
	auctionDetails, _ := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       now,
		Duration:        120,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	
	feeConfig, _ := fusionorder.NewFeeConfig(nil, nil, protocolFee, integratorFee, domains.ZeroBps)
	calculator := NewAmountCalculator(
		FromAuctionData(auctionDetails),
		FromFeeConfig(feeConfig),
	)

	takingAmount, _ := new(big.Int).SetString("1000000000000000000000000", 10)
	estimatedTakingAmount := takingAmount

	totalFee := calculator.GetTotalFee(takingAmount, estimatedTakingAmount, now)
	if totalFee.Cmp(big.NewInt(0)) <= 0 {
		t.Error("Expected positive fee for large amount")
	}
}
