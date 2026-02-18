package amountcalculator

import (
	"math/big"

	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
	fusionorder "github.com/tcbtcn/solana-fusion-sdk-go/fusion-order"
	"github.com/tcbtcn/solana-fusion-sdk-go/utils/math"
)

// FeeCalculator calculates fees for orders
type FeeCalculator struct {
	protocolFee   *domains.Bps
	integratorFee *domains.Bps
	surplusShare  *domains.Bps
}

// NewFeeCalculator creates a new FeeCalculator
func NewFeeCalculator(
	protocolFee *domains.Bps,
	integratorFee *domains.Bps,
	surplusShare *domains.Bps,
) *FeeCalculator {
	return &FeeCalculator{
		protocolFee:   protocolFee,
		integratorFee: integratorFee,
		surplusShare:  surplusShare,
	}
}

// FromFeeConfig creates a FeeCalculator from FeeConfig
func FromFeeConfig(feeConfig *fusionorder.FeeConfig) *FeeCalculator {
	if feeConfig == nil {
		return NewFeeCalculator(
			domains.ZeroBps,
			domains.ZeroBps,
			domains.ZeroBps,
		)
	}
	return NewFeeCalculator(
		feeConfig.ProtocolFee,
		feeConfig.IntegratorFee,
		feeConfig.SurplusShare,
	)
}

// GetIntegratorFee calculates the integrator fee
func (f *FeeCalculator) GetIntegratorFee(auctionTakingAmount *big.Int) *big.Int {
	if f.integratorFee.IsZero() {
		return big.NewInt(0)
	}

	// BPS is 0-10000 where 10000 = 100%
	// Fee = (amount * bps) / 10000
	return math.MulDiv(auctionTakingAmount, f.integratorFee.Value(), big.NewInt(10000), math.RoundingFloor)
}

// GetUserReceiveAmount calculates the amount the user will receive
func (f *FeeCalculator) GetUserReceiveAmount(
	auctionTakingAmount *big.Int,
	estimatedTakingAmount *big.Int,
) *big.Int {
	amounts := f.getAmounts(auctionTakingAmount, estimatedTakingAmount)
	return amounts.UserReceiveAmount
}

// GetProtocolFee calculates the protocol fee
func (f *FeeCalculator) GetProtocolFee(
	auctionTakingAmount *big.Int,
	estimatedTakingAmount *big.Int,
) *big.Int {
	amounts := f.getAmounts(auctionTakingAmount, estimatedTakingAmount)
	return amounts.ProtocolFeeAmount
}

type amountsResult struct {
	ProtocolFeeAmount *big.Int
	UserReceiveAmount *big.Int
}

func (f *FeeCalculator) getAmounts(
	auctionTakingAmount *big.Int,
	estimatedTakingAmount *big.Int,
) amountsResult {
	var protocolFee *big.Int
	if f.protocolFee.IsZero() {
		protocolFee = big.NewInt(0)
	} else {
		// BPS is 0-10000 where 10000 = 100%
		// Fee = (amount * bps) / 10000
		protocolFee = math.MulDiv(auctionTakingAmount, f.protocolFee.Value(), big.NewInt(10000), math.RoundingFloor)
	}

	integratorFee := f.GetIntegratorFee(auctionTakingAmount)

	userAmountWithoutFee := new(big.Int).Sub(auctionTakingAmount, protocolFee)
	userAmountWithoutFee = userAmountWithoutFee.Sub(userAmountWithoutFee, integratorFee)

	if userAmountWithoutFee.Cmp(estimatedTakingAmount) > 0 {
		surplus := new(big.Int).Sub(userAmountWithoutFee, estimatedTakingAmount)
		if !f.surplusShare.IsZero() {
			// Surplus share is stored as BPS (0-10000 where 10000 = 100%)
			// Surplus fee = (surplus * surplusShare) / 10000
			surplusFee := math.MulDiv(surplus, f.surplusShare.Value(), big.NewInt(10000), math.RoundingFloor)
			protocolFee = protocolFee.Add(protocolFee, surplusFee)
		}
	}

	userReceiveAmount := new(big.Int).Sub(auctionTakingAmount, protocolFee)
	userReceiveAmount = userReceiveAmount.Sub(userReceiveAmount, integratorFee)

	return amountsResult{
		ProtocolFeeAmount: protocolFee,
		UserReceiveAmount: userReceiveAmount,
	}
}
