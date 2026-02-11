package amountcalculator

import (
	"math/big"
)

// AmountCalculator calculates fees and amounts with accounting to auction
type AmountCalculator struct {
	auctionCalculator *AuctionCalculator
	feeCalculator     *FeeCalculator
}

// NewAmountCalculator creates a new AmountCalculator
func NewAmountCalculator(
	auctionCalculator *AuctionCalculator,
	feeCalculator *FeeCalculator,
) *AmountCalculator {
	return &AmountCalculator{
		auctionCalculator: auctionCalculator,
		feeCalculator:     feeCalculator,
	}
}

// CalcTakingAmount is deprecated - use utils/math.CalcTakingAmount instead
// This function is kept for backward compatibility but delegates to utils/math
func CalcTakingAmount(
	swapMakerAmount *big.Int,
	orderMakerAmount *big.Int,
	orderTakerAmount *big.Int,
) *big.Int {
	// Import utils/math to avoid duplication, but this creates a dependency
	// For now, keep the implementation here to avoid circular dependencies
	if orderMakerAmount.Sign() == 0 {
		return big.NewInt(0)
	}

	// (swapMakerAmount * orderTakerAmount + orderMakerAmount - 1) / orderMakerAmount
	numerator := new(big.Int).Mul(swapMakerAmount, orderTakerAmount)
	numerator = numerator.Add(numerator, orderMakerAmount)
	numerator = numerator.Sub(numerator, big.NewInt(1))

	return new(big.Int).Div(numerator, orderMakerAmount)
}

// GetRequiredTakingAmount returns how much resolver must pay to fill order
func (a *AmountCalculator) GetRequiredTakingAmount(takingAmount *big.Int, time uint32) *big.Int {
	return a.GetAuctionBumpedAmount(takingAmount, time)
}

// GetTotalFee returns total fee = integrator + protocol
func (a *AmountCalculator) GetTotalFee(
	takingAmount *big.Int,
	estimatedTakingAmount *big.Int,
	time uint32,
) *big.Int {
	if a.feeCalculator == nil {
		return big.NewInt(0)
	}

	auctionAmount := a.GetAuctionBumpedAmount(takingAmount, time)
	integratorFee := a.feeCalculator.GetIntegratorFee(auctionAmount)
	protocolFee := a.feeCalculator.GetProtocolFee(auctionAmount, estimatedTakingAmount)

	return new(big.Int).Add(integratorFee, protocolFee)
}

// GetUserReceiveAmount returns amount which will receive user
func (a *AmountCalculator) GetUserReceiveAmount(
	takingAmount *big.Int,
	estimatedTakingAmount *big.Int,
	time uint32,
) *big.Int {
	auctionAmount := a.GetRequiredTakingAmount(takingAmount, time)

	if a.feeCalculator == nil {
		return auctionAmount
	}

	return a.feeCalculator.GetUserReceiveAmount(auctionAmount, estimatedTakingAmount)
}

// GetIntegratorFee returns fee in dstToken which integrator gets
func (a *AmountCalculator) GetIntegratorFee(takingAmount *big.Int, time uint32) *big.Int {
	if a.feeCalculator == nil {
		return big.NewInt(0)
	}

	return a.feeCalculator.GetIntegratorFee(a.GetAuctionBumpedAmount(takingAmount, time))
}

// GetProtocolFee returns fee in dstToken which protocol gets
func (a *AmountCalculator) GetProtocolFee(
	takingAmount *big.Int,
	estimatedTakingAmount *big.Int,
	time uint32,
) *big.Int {
	if a.feeCalculator == nil {
		return big.NewInt(0)
	}

	return a.feeCalculator.GetProtocolFee(
		a.GetAuctionBumpedAmount(takingAmount, time),
		estimatedTakingAmount,
	)
}

// GetAuctionBumpedAmount returns the taking amount with auction bump applied
func (a *AmountCalculator) GetAuctionBumpedAmount(takingAmount *big.Int, time uint32) *big.Int {
	rateBump := a.auctionCalculator.CalcRateBump(time)
	return CalcAuctionTakingAmount(takingAmount, rateBump)
}
