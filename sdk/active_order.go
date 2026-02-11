package sdk

import (
	"math/big"

	"github.com/dawitel/solana-fusion-sdk-go/api/orders"
	"github.com/dawitel/solana-fusion-sdk-go/domains"
	fusionorder "github.com/dawitel/solana-fusion-sdk-go/fusion-order"
)

// ActiveOrder represents an active order
type ActiveOrder struct {
	CreationTxSignature  string
	Maker                *domains.Address
	Order                *fusionorder.FusionOrder
	RemainingMakerAmount *big.Int
}

// ActiveOrderFromJSON creates an ActiveOrder from JSON DTO
func ActiveOrderFromJSON(json *orders.OrderInfoDTO) (*ActiveOrder, error) {
	maker, err := domains.NewAddress(json.Maker)
	if err != nil {
		return nil, err
	}

	// Convert order DTO to FusionOrder
	orderJSON := &fusionorder.FusionOrderJSON{
		ID:                          json.Order.ID,
		Receiver:                    json.Order.Receiver,
		CancellationAuctionDuration: json.Order.CancellationAuctionDuration,
		SrcMint:                     json.Order.SrcMint,
		DstMint:                     json.Order.DstMint,
		SrcAmount:                   json.Order.SrcAmount,
		MinDstAmount:                json.Order.MinDstAmount,
		EstimatedDstAmount:          json.Order.EstimatedDstAmount,
		ExpirationTime:              json.Order.ExpirationTime,
		SrcAssetIsNative:            json.Order.SrcAssetIsNative,
		DstAssetIsNative:            json.Order.DstAssetIsNative,
		Fee: fusionorder.FeeJSON{
			ProtocolFee:            json.Order.Fee.ProtocolFee,
			IntegratorFee:          json.Order.Fee.IntegratorFee,
			SurplusPercentage:      json.Order.Fee.SurplusPercentage,
			MaxCancellationPremium: json.Order.Fee.MaxCancellationPremium,
		},
		DutchAuctionData: fusionorder.DutchAuctionDataJSON{
			StartTime:           json.Order.DutchAuctionData.StartTime,
			Duration:            json.Order.DutchAuctionData.Duration,
			InitialRateBump:     json.Order.DutchAuctionData.InitialRateBump,
			PointsAndTimeDeltas: make([]fusionorder.PointAndTimeDeltaJSON, len(json.Order.DutchAuctionData.PointsAndTimeDeltas)),
		},
	}

	for i, p := range json.Order.DutchAuctionData.PointsAndTimeDeltas {
		orderJSON.DutchAuctionData.PointsAndTimeDeltas[i] = fusionorder.PointAndTimeDeltaJSON{
			RateBump:  p.RateBump,
			TimeDelta: p.TimeDelta,
		}
	}

	if json.Order.Fee.ProtocolDstAta != "" {
		orderJSON.Fee.ProtocolDstAta = &json.Order.Fee.ProtocolDstAta
	}
	if json.Order.Fee.IntegratorDstAta != "" {
		orderJSON.Fee.IntegratorDstAta = &json.Order.Fee.IntegratorDstAta
	}

	order, err := fusionorder.FromJSON(orderJSON)
	if err != nil {
		return nil, err
	}

	remainingMakerAmount, _ := new(big.Int).SetString(json.RemainingMakerAmount, 10)

	return &ActiveOrder{
		CreationTxSignature:  json.TxSignature,
		Maker:                maker,
		Order:                order,
		RemainingMakerAmount: remainingMakerAmount,
	}, nil
}
