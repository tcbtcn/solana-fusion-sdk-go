package fusionorder

import (
	"github.com/tcbtcn/solana-fusion-sdk-go/utils/validation"
)

// AuctionPoint represents a point in the auction curve
type AuctionPoint struct {
	Coefficient uint16 // Rate bump coefficient
	Delay       uint16 // Time delta from previous point
}

// AuctionDetails represents auction configuration
type AuctionDetails struct {
	StartTime       uint32
	Duration        uint32
	InitialRateBump uint16
	Points          []AuctionPoint
}

// NewAuctionDetails creates a new AuctionDetails
func NewAuctionDetails(auction struct {
	StartTime       uint32
	InitialRateBump uint16
	Duration        uint32
	Points          []AuctionPoint
}) (*AuctionDetails, error) {
	// Validate points
	for _, point := range auction.Points {
		if err := validation.AssertUInteger(point.Delay, nil); err != nil {
			return nil, err
		}
		if err := validation.AssertUInteger(point.Coefficient, nil); err != nil {
			return nil, err
		}
	}

	// Validate main fields
	if err := validation.AssertUInteger(auction.StartTime, nil); err != nil {
		return nil, err
	}
	if err := validation.AssertUInteger(auction.Duration, nil); err != nil {
		return nil, err
	}
	if err := validation.AssertUInteger(auction.InitialRateBump, nil); err != nil {
		return nil, err
	}

	return &AuctionDetails{
		StartTime:       auction.StartTime,
		Duration:        auction.Duration,
		InitialRateBump: auction.InitialRateBump,
		Points:          auction.Points,
	}, nil
}

// NoAuction creates an AuctionDetails with no auction (zero rate bump)
func NoAuction(startTime, duration uint32) *AuctionDetails {
	return &AuctionDetails{
		StartTime:       startTime,
		Duration:        duration,
		InitialRateBump: 0,
		Points:          []AuctionPoint{},
	}
}

// FinishTime returns the finish time of the auction
func (a *AuctionDetails) FinishTime() uint32 {
	return a.StartTime + a.Duration
}
