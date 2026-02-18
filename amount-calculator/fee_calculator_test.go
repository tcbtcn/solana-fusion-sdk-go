package amountcalculator

import (
	"math/big"
	"testing"

	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
	fusionorder "github.com/tcbtcn/solana-fusion-sdk-go/fusion-order"
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
	surplusShare := domains.BpsFromPercent(50.0, nil)
	protocolDstAta := domains.MustAddressFromString("11111111111111111111111111111111")

	// Valid: protocolDstAta provided with protocolFee and surplusShare
	feeConfig, err := fusionorder.NewFeeConfig(protocolDstAta, nil, protocolFee, domains.ZeroBps, surplusShare)
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
	// Expected: 1000 - 10 (protocolFee) - 20 (integratorFee) = 970
	expected := big.NewInt(970)

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
	// Expected: protocolFee = 10, integratorFee = 20
	// userAmountWithoutFee = 1000 - 10 - 20 = 970
	// surplus = 970 - 500 = 470
	// surplusFee = 470 * 50% = 235
	// totalProtocolFee = 10 + 235 = 245
	// userReceiveAmount = 1000 - 245 - 20 = 735
	expected := big.NewInt(735)

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
