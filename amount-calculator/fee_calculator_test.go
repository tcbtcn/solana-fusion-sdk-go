package amountcalculator

import (
	"math/big"
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
	fusionorder "github.com/dawitel/solana-fusion-sdk-go/fusion-order"
)

func TestNewFeeCalculator(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	integratorFee := domains.BpsFromPercent(2.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)

	calculator := NewFeeCalculator(protocolFee, integratorFee, surplusShare)

	if calculator == nil {
		t.Fatal("Expected non-nil calculator")
	}
}

func TestFromFeeConfig_Valid(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	integratorFee := domains.BpsFromPercent(2.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)

	feeConfig, err := fusionorder.NewFeeConfig(nil, nil, protocolFee, integratorFee, surplusShare)
	if err != nil {
		t.Fatalf("Failed to create fee config: %v", err)
	}

	calculator := FromFeeConfig(feeConfig)
	if calculator == nil {
		t.Fatal("Expected non-nil calculator")
	}
}

func TestFromFeeConfig_Nil(t *testing.T) {
	calculator := FromFeeConfig(nil)
	if calculator == nil {
		t.Fatal("Expected non-nil calculator")
	}
}

func TestFeeCalculator_GetIntegratorFee(t *testing.T) {
	integratorFee := domains.BpsFromPercent(2.0, nil)
	calculator := NewFeeCalculator(domains.ZeroBps, integratorFee, domains.ZeroBps)

	auctionTakingAmount := big.NewInt(1000)
	fee := calculator.GetIntegratorFee(auctionTakingAmount)

	expected := big.NewInt(20)
	if fee.Cmp(expected) != 0 {
		t.Errorf("Expected integrator fee %s, got %s", expected.String(), fee.String())
	}
}

func TestFeeCalculator_GetIntegratorFee_Zero(t *testing.T) {
	calculator := NewFeeCalculator(domains.ZeroBps, domains.ZeroBps, domains.ZeroBps)

	auctionTakingAmount := big.NewInt(1000)
	fee := calculator.GetIntegratorFee(auctionTakingAmount)

	expected := big.NewInt(0)
	if fee.Cmp(expected) != 0 {
		t.Errorf("Expected integrator fee %s, got %s", expected.String(), fee.String())
	}
}

func TestFeeCalculator_GetProtocolFee_NoSurplus(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	calculator := NewFeeCalculator(protocolFee, domains.ZeroBps, domains.ZeroBps)

	auctionTakingAmount := big.NewInt(1000)
	estimatedTakingAmount := big.NewInt(1000)

	fee := calculator.GetProtocolFee(auctionTakingAmount, estimatedTakingAmount)
	expected := big.NewInt(10)

	if fee.Cmp(expected) != 0 {
		t.Errorf("Expected protocol fee %s, got %s", expected.String(), fee.String())
	}
}

func TestFeeCalculator_GetProtocolFee_WithSurplus(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)
	calculator := NewFeeCalculator(protocolFee, domains.ZeroBps, surplusShare)

	auctionTakingAmount := big.NewInt(1000)
	estimatedTakingAmount := big.NewInt(500)

	fee := calculator.GetProtocolFee(auctionTakingAmount, estimatedTakingAmount)
	expected := big.NewInt(255)

	if fee.Cmp(expected) != 0 {
		t.Errorf("Expected protocol fee with surplus %s, got %s", expected.String(), fee.String())
	}
}

func TestFeeCalculator_GetUserReceiveAmount_NoSurplus(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	integratorFee := domains.BpsFromPercent(2.0, nil)
	calculator := NewFeeCalculator(protocolFee, integratorFee, domains.ZeroBps)

	auctionTakingAmount := big.NewInt(1000)
	estimatedTakingAmount := big.NewInt(1000)

	amount := calculator.GetUserReceiveAmount(auctionTakingAmount, estimatedTakingAmount)
	expected := big.NewInt(990)

	if amount.Cmp(expected) != 0 {
		t.Errorf("Expected user receive amount %s, got %s", expected.String(), amount.String())
	}
}

func TestFeeCalculator_GetUserReceiveAmount_WithSurplus(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	integratorFee := domains.BpsFromPercent(2.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)
	calculator := NewFeeCalculator(protocolFee, integratorFee, surplusShare)

	auctionTakingAmount := big.NewInt(1000)
	estimatedTakingAmount := big.NewInt(500)

	amount := calculator.GetUserReceiveAmount(auctionTakingAmount, estimatedTakingAmount)
	expected := big.NewInt(745)

	if amount.Cmp(expected) != 0 {
		t.Errorf("Expected user receive amount with surplus %s, got %s", expected.String(), amount.String())
	}
}

func TestFeeCalculator_ZeroFees(t *testing.T) {
	calculator := NewFeeCalculator(domains.ZeroBps, domains.ZeroBps, domains.ZeroBps)

	auctionTakingAmount := big.NewInt(1000)
	estimatedTakingAmount := big.NewInt(1000)

	protocolFee := calculator.GetProtocolFee(auctionTakingAmount, estimatedTakingAmount)
	integratorFee := calculator.GetIntegratorFee(auctionTakingAmount)
	userAmount := calculator.GetUserReceiveAmount(auctionTakingAmount, estimatedTakingAmount)

	if protocolFee.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected zero protocol fee, got %s", protocolFee.String())
	}
	if integratorFee.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Expected zero integrator fee, got %s", integratorFee.String())
	}
	if userAmount.Cmp(auctionTakingAmount) != 0 {
		t.Errorf("Expected user amount to equal taking amount, got %s", userAmount.String())
	}
}
