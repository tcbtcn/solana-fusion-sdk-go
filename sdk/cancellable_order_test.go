package sdk

import (
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/api/orders"
	"github.com/dawitel/solana-fusion-sdk-go/domains"
	fusionorder "github.com/dawitel/solana-fusion-sdk-go/fusion-order"
)

func TestCancellableOrderFromJSON(t *testing.T) {
	json := &orders.OrderCancellableByResolverInfoDTO{
		OrderHash:   "test-hash",
		TxSignature: "tx-signature",
		Maker:       "11111111111111111111111111111111",
		Order: orders.OrderDTO{
			ID:                          1,
			Receiver:                    "11111111111111111111111111111111",
			CancellationAuctionDuration: 100,
			SrcMint:                     "So11111111111111111111111111111111111111112",
			DstMint:                     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
			SrcAmount:                   "1000000000000000000",
			MinDstAmount:                "1420000000",
			EstimatedDstAmount:          "1420000000",
			ExpirationTime:              1000000000,
			SrcAssetIsNative:            false,
			DstAssetIsNative:            false,
			Fee: orders.FeeDTO{
				ProtocolFee:       0,
				IntegratorFee:     0,
				SurplusPercentage: 0,
			},
			DutchAuctionData: orders.DutchAuctionDataDTO{
				StartTime:           1000000000,
				Duration:            180,
				InitialRateBump:     0,
				PointsAndTimeDeltas: []orders.PointAndTimeDeltaDTO{},
			},
		},
	}

	maker, err := domains.NewAddress(json.Maker)
	if err != nil {
		t.Fatalf("Failed to create maker address: %v", err)
	}

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
			PointsAndTimeDeltas: []fusionorder.PointAndTimeDeltaJSON{},
		},
	}

	order, err := fusionorder.FromJSON(orderJSON)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	cancellableOrder := &CancellableOrder{
		Maker: maker,
		Order: order,
	}

	if !cancellableOrder.Maker.Equal(maker) {
		t.Error("Expected maker to match")
	}
	if cancellableOrder.Order == nil {
		t.Fatal("Expected non-nil order")
	}
}
