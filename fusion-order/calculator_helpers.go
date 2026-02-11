package fusionorder

import (
	"math/big"

	"github.com/dawitel/solana-fusion-sdk-go/utils/math"
)

// Calculator is an interface for amount calculators
// This allows fusion-order to work with amount-calculator without direct import
type Calculator interface {
	GetRequiredTakingAmount(takingAmount *big.Int, time uint32) *big.Int
	GetUserReceiveAmount(takingAmount *big.Int, estimatedTakingAmount *big.Int, time uint32) *big.Int
	GetIntegratorFee(takingAmount *big.Int, time uint32) *big.Int
	GetProtocolFee(takingAmount *big.Int, estimatedTakingAmount *big.Int, time uint32) *big.Int
}

// CalcTakingAmountWithCalculator calculates taking amount using an external calculator
func (f *FusionOrder) CalcTakingAmountWithCalculator(calculator Calculator, makingAmount *big.Int, time uint32) *big.Int {
	takingAmount := math.CalcTakingAmount(
		makingAmount,
		f.orderConfig.SrcAmount,
		f.orderConfig.MinDstAmount,
	)
	return calculator.GetRequiredTakingAmount(takingAmount, time)
}

// GetUserReceiveAmountWithCalculator returns how much user will receive using an external calculator
func (f *FusionOrder) GetUserReceiveAmountWithCalculator(calculator Calculator, makingAmount *big.Int, time uint32) *big.Int {
	takingAmount := math.CalcTakingAmount(
		makingAmount,
		f.orderConfig.SrcAmount,
		f.orderConfig.MinDstAmount,
	)
	return calculator.GetUserReceiveAmount(
		takingAmount,
		f.orderConfig.EstimatedDstAmount,
		time,
	)
}

// GetIntegratorFeeWithCalculator returns integrator fee using an external calculator
func (f *FusionOrder) GetIntegratorFeeWithCalculator(calculator Calculator, time uint32, makingAmount *big.Int) *big.Int {
	if makingAmount == nil {
		makingAmount = f.orderConfig.SrcAmount
	}
	takingAmount := math.CalcTakingAmount(
		makingAmount,
		f.orderConfig.SrcAmount,
		f.orderConfig.MinDstAmount,
	)
	return calculator.GetIntegratorFee(takingAmount, time)
}

// GetProtocolFeeWithCalculator returns protocol fee using an external calculator
func (f *FusionOrder) GetProtocolFeeWithCalculator(calculator Calculator, time uint32, makingAmount *big.Int) *big.Int {
	if makingAmount == nil {
		makingAmount = f.orderConfig.SrcAmount
	}
	takingAmount := math.CalcTakingAmount(
		makingAmount,
		f.orderConfig.SrcAmount,
		f.orderConfig.MinDstAmount,
	)
	return calculator.GetProtocolFee(
		takingAmount,
		f.orderConfig.EstimatedDstAmount,
		time,
	)
}
