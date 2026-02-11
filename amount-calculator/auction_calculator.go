package amountcalculator

import (
	"math/big"

	fusionorder "github.com/dawitel/solana-fusion-sdk-go/fusion-order"
	"github.com/dawitel/solana-fusion-sdk-go/utils/math"
)

// RateBumpDenominator is 100% (100,000)
var RateBumpDenominator = big.NewInt(100_000)

// AuctionCalculator calculates auction rates and amounts
type AuctionCalculator struct {
	startTime       uint32
	duration        uint32
	initialRateBump uint16
	points          []fusionorder.AuctionPoint
}

// NewAuctionCalculator creates a new AuctionCalculator
func NewAuctionCalculator(
	startTime uint32,
	duration uint32,
	initialRateBump uint16,
	points []fusionorder.AuctionPoint,
) *AuctionCalculator {
	return &AuctionCalculator{
		startTime:       startTime,
		duration:        duration,
		initialRateBump: initialRateBump,
		points:          points,
	}
}

// FromAuctionData creates an AuctionCalculator from AuctionDetails
func FromAuctionData(details *fusionorder.AuctionDetails) *AuctionCalculator {
	return NewAuctionCalculator(
		details.StartTime,
		details.Duration,
		details.InitialRateBump,
		details.Points,
	)
}

// FinishTime returns the finish time of the auction
func (a *AuctionCalculator) FinishTime() uint32 {
	return a.startTime + a.duration
}

// CalcInitialRateBump calculates the initial rate bump from start and end amounts
func CalcInitialRateBump(startAmount, endAmount *big.Int) uint16 {
	if endAmount.Sign() == 0 {
		return 0
	}

	// bump = (RATE_BUMP_DENOMINATOR * startAmount) / endAmount - RATE_BUMP_DENOMINATOR
	numerator := new(big.Int).Mul(RateBumpDenominator, startAmount)
	bump := new(big.Int).Div(numerator, endAmount)
	bump = bump.Sub(bump, RateBumpDenominator)

	return uint16(bump.Uint64())
}

// CalcAuctionTakingAmount calculates the taking amount with rate bump applied
func CalcAuctionTakingAmount(takingAmount *big.Int, rate uint16) *big.Int {
	// (takingAmount * (rate + RATE_BUMP_DENOMINATOR)) / RATE_BUMP_DENOMINATOR
	rateBig := big.NewInt(int64(rate))
	denominator := new(big.Int).Add(rateBig, RateBumpDenominator)

	return math.MulDiv(takingAmount, denominator, RateBumpDenominator, math.RoundingCeil)
}

// CalcAuctionTakingAmountAtTime calculates the taking amount at a specific block time
func (a *AuctionCalculator) CalcAuctionTakingAmount(takingAmount *big.Int, blockTime uint32) *big.Int {
	rateBump := a.CalcRateBump(blockTime)
	return CalcAuctionTakingAmount(takingAmount, rateBump)
}

// CalcRateBump calculates the rate bump at a specific block time
func (a *AuctionCalculator) CalcRateBump(blockTime uint32) uint16 {
	auctionFinishTime := a.FinishTime()

	if blockTime <= a.startTime {
		return a.initialRateBump
	} else if blockTime >= auctionFinishTime {
		return 0
	}

	currentRateBump := big.NewInt(int64(a.initialRateBump))
	currentPointTime := big.NewInt(int64(a.startTime))
	blockTimeBN := big.NewInt(int64(blockTime))
	finishTimeBN := big.NewInt(int64(auctionFinishTime))

	for _, point := range a.points {
		nextPointTime := new(big.Int).Add(big.NewInt(int64(point.Delay)), currentPointTime)

		if blockTimeBN.Cmp(nextPointTime) <= 0 {
			// Linear interpolation
			timeDiff := new(big.Int).Sub(blockTimeBN, currentPointTime)
			nextRateBump := big.NewInt(int64(point.Coefficient))

			// ((blockTime - currentPointTime) * nextRateBump + (nextPointTime - blockTime) * currentRateBump) / (nextPointTime - currentPointTime)
			term1 := new(big.Int).Mul(timeDiff, nextRateBump)
			term2 := new(big.Int).Sub(nextPointTime, blockTimeBN)
			term2 = term2.Mul(term2, currentRateBump)
			numerator := new(big.Int).Add(term1, term2)
			denominator := new(big.Int).Sub(nextPointTime, currentPointTime)

			result := new(big.Int).Div(numerator, denominator)
			return uint16(result.Uint64())
		}

		currentPointTime = nextPointTime
		currentRateBump = big.NewInt(int64(point.Coefficient))
	}

	// After all points, interpolate to finish time
	timeDiff := new(big.Int).Sub(finishTimeBN, blockTimeBN)
	numerator := new(big.Int).Mul(timeDiff, currentRateBump)
	denominator := new(big.Int).Sub(finishTimeBN, currentPointTime)

	result := new(big.Int).Div(numerator, denominator)
	return uint16(result.Uint64())
}
