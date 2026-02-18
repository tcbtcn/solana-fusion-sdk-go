package sdk

import (
	"math/big"

	"github.com/tcbtcn/solana-fusion-sdk-go/api/orders"
	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
	fusionorder "github.com/tcbtcn/solana-fusion-sdk-go/fusion-order"
)

// OrderStatus represents the status of an order
type OrderStatus struct {
	Maker                   *domains.Address
	OrderHash               string
	Status                  orders.OrderStatus
	Order                   *fusionorder.FusionOrder
	ApproximateTakingAmount *big.Int
	CancelTx                *string
	ExpirationTime          uint32
	Fills                   []Fill
	CreatedAt               uint32
	SrcTokenPriceUsd        float64
	DstTokenPriceUsd        float64
	Cancelable              bool
}

// Fill represents a fill
type Fill struct {
	TxSignature              string
	FilledMakerAmount        *big.Int
	FilledAuctionTakerAmount *big.Int
}

// OrderStatusFromJSON creates an OrderStatus from JSON DTO
func OrderStatusFromJSON(json *orders.OrderStatusDTO) (*OrderStatus, error) {
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

	approximateTakingAmount, _ := new(big.Int).SetString(json.ApproximateTakingAmount, 10)

	fills := make([]Fill, len(json.Fills))
	for i, f := range json.Fills {
		filledMakerAmount, _ := new(big.Int).SetString(f.FilledMakerAmount, 10)
		filledAuctionTakerAmount, _ := new(big.Int).SetString(f.FilledAuctionTakerAmount, 10)
		fills[i] = Fill{
			TxSignature:              f.TxSignature,
			FilledMakerAmount:        filledMakerAmount,
			FilledAuctionTakerAmount: filledAuctionTakerAmount,
		}
	}

	return &OrderStatus{
		Maker:                   maker,
		OrderHash:               json.OrderHash,
		Status:                  json.Status,
		Order:                   order,
		ApproximateTakingAmount: approximateTakingAmount,
		CancelTx:                json.CancelTx,
		ExpirationTime:          json.ExpirationTime,
		Fills:                   fills,
		CreatedAt:               json.CreatedAt,
		SrcTokenPriceUsd:        json.SrcTokenPriceUsd,
		DstTokenPriceUsd:        json.DstTokenPriceUsd,
		Cancelable:              json.Cancelable,
	}, nil
}

// IsActive checks if the order is active
func (o *OrderStatus) IsActive() bool {
	return o.Status == orders.OrderStatusInProgress
}
