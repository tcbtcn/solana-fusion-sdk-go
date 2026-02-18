package sdk

import (
	"math/big"

	"github.com/tcbtcn/solana-fusion-sdk-go/api/quoter"
)

// Preset represents a preset configuration
type Preset struct {
	StartAuctionIn     uint32
	AuctionDuration    uint32
	InitialRateBump    uint16
	AuctionStartAmount *big.Int
	AuctionEndAmount   *big.Int
	CostInDstToken     *big.Int
	Points             []PresetPoint
}

// PresetPoint represents a point in a preset
type PresetPoint struct {
	Delay       uint16
	Coefficient uint16
}

// PresetFromJSON creates a Preset from JSON DTO
func PresetFromJSON(json quoter.PresetDTO) *Preset {
	points := make([]PresetPoint, len(json.Points))
	for i, p := range json.Points {
		points[i] = PresetPoint{
			Delay:       p.Delay,
			Coefficient: p.Coefficient,
		}
	}

	auctionStartAmount, _ := new(big.Int).SetString(json.AuctionStartAmount, 10)
	auctionEndAmount, _ := new(big.Int).SetString(json.AuctionEndAmount, 10)
	costInDstToken, _ := new(big.Int).SetString(json.CostInDstToken, 10)

	return &Preset{
		StartAuctionIn:     json.StartAuctionIn,
		AuctionDuration:    json.AuctionDuration,
		InitialRateBump:    json.InitialRateBump,
		AuctionStartAmount: auctionStartAmount,
		AuctionEndAmount:   auctionEndAmount,
		CostInDstToken:     costInDstToken,
		Points:             points,
	}
}
