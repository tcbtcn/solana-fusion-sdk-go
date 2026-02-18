package sdk

import (
	"math/big"
	"testing"

	"github.com/tcbtcn/solana-fusion-sdk-go/api/orders"
	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
)

func TestActiveOrderFromJSON_Success(t *testing.T) {
	json := &orders.OrderInfoDTO{
		OrderHash:   "test-hash",
		TxSignature: "tx-signature",
		Maker:       "11111111111111111111111111111111",
		Order: orders.OrderDTO{
			ID:                          1,
			Receiver:                    "11111111111111111111111111111111",
			CancellationAuctionDuration: 0,
			SrcMint:                     "So11111111111111111111111111111111111111112",
			DstMint:                     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
			SrcAmount:                   "1000000000000000000",
			MinDstAmount:                "1420000000",
			EstimatedDstAmount:          "1420000000",
			ExpirationTime:              1000000000,
			SrcAssetIsNative:            false,
			DstAssetIsNative:            false,
			Fee: orders.FeeDTO{
				ProtocolFee:            0,
				IntegratorFee:          0,
				SurplusPercentage:      0,
				MaxCancellationPremium: "0",
			},
			DutchAuctionData: orders.DutchAuctionDataDTO{
				StartTime:           1000000000,
				Duration:            180,
				InitialRateBump:     0,
				PointsAndTimeDeltas: []orders.PointAndTimeDeltaDTO{},
			},
		},
		RemainingMakerAmount: "500000000000000000",
	}

	activeOrder, err := ActiveOrderFromJSON(json)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if activeOrder == nil {
		t.Fatal("Expected non-nil active order")
	}
	if activeOrder.CreationTxSignature != "tx-signature" {
		t.Errorf("Expected CreationTxSignature 'tx-signature', got %s", activeOrder.CreationTxSignature)
	}
	if !activeOrder.Maker.Equal(domains.MustAddressFromString("11111111111111111111111111111111")) {
		t.Error("Expected maker to match")
	}
	if activeOrder.RemainingMakerAmount.Cmp(big.NewInt(500000000000000000)) != 0 {
		t.Errorf("Expected RemainingMakerAmount 500000000000000000, got %s", activeOrder.RemainingMakerAmount.String())
	}
	if activeOrder.Order == nil {
		t.Fatal("Expected non-nil order")
	}
}

func TestActiveOrderFromJSON_InvalidMaker(t *testing.T) {
	json := &orders.OrderInfoDTO{
		OrderHash:   "test-hash",
		TxSignature: "tx-signature",
		Maker:       "invalid-address",
		Order: orders.OrderDTO{
			ID:                          1,
			Receiver:                    "11111111111111111111111111111111",
			CancellationAuctionDuration: 0,
			SrcMint:                     "So11111111111111111111111111111111111111112",
			DstMint:                     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
			SrcAmount:                   "1000000000000000000",
			MinDstAmount:                "1420000000",
			EstimatedDstAmount:          "1420000000",
			ExpirationTime:              1000000000,
			SrcAssetIsNative:            false,
			DstAssetIsNative:            false,
			Fee: orders.FeeDTO{
				ProtocolFee:            0,
				IntegratorFee:          0,
				SurplusPercentage:      0,
				MaxCancellationPremium: "0",
			},
			DutchAuctionData: orders.DutchAuctionDataDTO{
				StartTime:           1000000000,
				Duration:            180,
				InitialRateBump:     0,
				PointsAndTimeDeltas: []orders.PointAndTimeDeltaDTO{},
			},
		},
		RemainingMakerAmount: "500000000000000000",
	}

	_, err := ActiveOrderFromJSON(json)
	if err == nil {
		t.Fatal("Expected error for invalid maker address")
	}
}
