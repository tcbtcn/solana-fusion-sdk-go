package math

import "math/big"

// Rounding represents rounding mode
type Rounding int

const (
	// RoundingCeil rounds up
	RoundingCeil Rounding = iota
	// RoundingFloor rounds down
	RoundingFloor
)

// MulDiv performs (a * b) / x with optional rounding
func MulDiv(a, b, x *big.Int, rounding Rounding) *big.Int {
	if x.Sign() == 0 {
		panic("division by zero")
	}

	product := new(big.Int).Mul(a, b)
	result := new(big.Int).Div(product, x)

	if rounding == RoundingCeil {
		// Check if there's a remainder - if so, round up
		remainder := new(big.Int).Mod(product, x)
		if remainder.Sign() > 0 {
			result = result.Add(result, big.NewInt(1))
		}
	}

	return result
}

// CalcTakingAmount calculates taker amount by linear proportion
// Returns ceiled taker amount: (swapMakerAmount * orderTakerAmount + orderMakerAmount - 1) / orderMakerAmount
func CalcTakingAmount(
	swapMakerAmount *big.Int,
	orderMakerAmount *big.Int,
	orderTakerAmount *big.Int,
) *big.Int {
	if orderMakerAmount.Sign() == 0 {
		return big.NewInt(0)
	}

	numerator := new(big.Int).Mul(swapMakerAmount, orderTakerAmount)
	numerator = numerator.Add(numerator, orderMakerAmount)
	numerator = numerator.Sub(numerator, big.NewInt(1))

	return new(big.Int).Div(numerator, orderMakerAmount)
}
