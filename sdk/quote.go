package sdk

import (
	"errors"
	"math/big"

	"github.com/dawitel/solana-fusion-sdk-go/api/quoter"
	"github.com/dawitel/solana-fusion-sdk-go/domains"
	fusionorder "github.com/dawitel/solana-fusion-sdk-go/fusion-order"
	"github.com/dawitel/solana-fusion-sdk-go/utils"
	"github.com/dawitel/solana-fusion-sdk-go/utils/time"
)

// Quote represents a quote for a swap
type Quote struct {
	SrcToken           *domains.Address
	DstToken           *domains.Address
	Signer             *domains.Address
	QuoteID            string
	SrcAmount          *big.Int
	DstAmount          *big.Int
	Presets            Presets
	RecommendedPreset  quoter.PresetType
	PriceImpactPercent float64
}

// Presets represents all presets
type Presets struct {
	Fast   *Preset
	Medium *Preset
	Slow   *Preset
}

// QuoteFromJSON creates a Quote from JSON DTO
func QuoteFromJSON(
	srcToken *domains.Address,
	dstToken *domains.Address,
	signer *domains.Address,
	json *quoter.QuoteDTO,
) (*Quote, error) {
	if json.QuoteID == nil {
		return nil, errors.New("quoteId is required. Use enableEstimate=true to generate it")
	}

	srcAmount, ok := new(big.Int).SetString(json.SrcAmount, 10)
	if !ok {
		return nil, errors.New("invalid srcAmount")
	}

	dstAmount, ok := new(big.Int).SetString(json.DstAmount, 10)
	if !ok {
		return nil, errors.New("invalid dstAmount")
	}

	return &Quote{
		SrcToken:  srcToken,
		DstToken:  dstToken,
		Signer:    signer,
		QuoteID:   *json.QuoteID,
		SrcAmount: srcAmount,
		DstAmount: dstAmount,
		Presets: Presets{
			Fast:   PresetFromJSON(json.Presets.Fast),
			Medium: PresetFromJSON(json.Presets.Medium),
			Slow:   PresetFromJSON(json.Presets.Slow),
		},
		RecommendedPreset:  json.RecommendedPreset,
		PriceImpactPercent: json.PriceImpactPercent,
	}, nil
}

// ToOrder converts a Quote to a FusionOrder
// If presetType is empty, uses RecommendedPreset (matching TypeScript default behavior)
// If receiver is nil, uses Signer (matching TypeScript default behavior)
func (q *Quote) ToOrder(presetType quoter.PresetType, receiver *domains.Address) (*fusionorder.FusionOrder, error) {
	if receiver == nil {
		receiver = q.Signer
	}

	if presetType == "" {
		presetType = q.RecommendedPreset
	}

	var preset *Preset
	switch presetType {
	case quoter.PresetTypeFast:
		preset = q.Presets.Fast
	case quoter.PresetTypeMedium:
		preset = q.Presets.Medium
	case quoter.PresetTypeSlow:
		preset = q.Presets.Slow
	default:
		// Fallback to RecommendedPreset if unknown preset type
		switch q.RecommendedPreset {
		case quoter.PresetTypeFast:
			preset = q.Presets.Fast
		case quoter.PresetTypeMedium:
			preset = q.Presets.Medium
		case quoter.PresetTypeSlow:
			preset = q.Presets.Slow
		default:
			preset = q.Presets.Fast
		}
	}

	if preset == nil {
		return nil, errors.New("preset is nil: invalid preset type or missing preset data")
	}

	// Create auction details
	points := make([]fusionorder.AuctionPoint, len(preset.Points))
	for i, p := range preset.Points {
		points[i] = fusionorder.AuctionPoint{
			Coefficient: p.Coefficient,
			Delay:       p.Delay,
		}
	}

	auctionDetails, err := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       uint32(time.Now()) + preset.StartAuctionIn,
		Duration:        preset.AuctionDuration,
		InitialRateBump: preset.InitialRateBump,
		Points:          points,
	})
	if err != nil {
		return nil, err
	}

	// Create order info
	orderInfo := fusionorder.OrderInfoData{
		ID:                 utils.ID(),
		SrcAmount:          q.SrcAmount,
		MinDstAmount:       preset.AuctionEndAmount,
		EstimatedDstAmount: q.DstAmount,
		Receiver:           receiver,
		SrcMint:            q.SrcToken,
		DstMint:            q.DstToken,
	}

	// Create order
	return fusionorder.NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *fusionorder.FeeConfig
		ResolverCancellationConfig *fusionorder.ResolverCancellationConfig
	}{
		SrcAssetIsNative:           q.SrcToken.IsNative(),
		DstAssetIsNative:           q.DstToken.IsNative(),
		OrderExpirationDelay:       0,   // Use default
		Fees:                       nil, // Use default
		ResolverCancellationConfig: nil, // Use default
	})
}
