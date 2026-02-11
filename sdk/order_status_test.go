package sdk

import (
	"math/big"
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/api/orders"
)

func TestOrderStatusFromJSON_Success(t *testing.T) {
	json := &orders.OrderStatusDTO{
		Maker:     "11111111111111111111111111111111",
		OrderHash: "test-hash",
		Status:    orders.OrderStatusInProgress,
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
		ApproximateTakingAmount: "1420000000",
		ExpirationTime:          1000000000,
		Fills:                   []orders.FillDTO{},
		CreatedAt:               1000000000,
		SrcTokenPriceUsd:        1.0,
		DstTokenPriceUsd:        1.0,
		Cancelable:              true,
	}

	status, err := OrderStatusFromJSON(json)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if status == nil {
		t.Fatal("Expected non-nil status")
	}
	if status.OrderHash != "test-hash" {
		t.Errorf("Expected OrderHash 'test-hash', got %s", status.OrderHash)
	}
	if status.Status != orders.OrderStatusInProgress {
		t.Errorf("Expected Status %d, got %d", orders.OrderStatusInProgress, status.Status)
	}
	if !status.IsActive() {
		t.Error("Expected order to be active")
	}
}

func TestOrderStatusFromJSON_WithFills(t *testing.T) {
	json := &orders.OrderStatusDTO{
		Maker:     "11111111111111111111111111111111",
		OrderHash: "test-hash",
		Status:    orders.OrderStatusInProgress,
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
		ApproximateTakingAmount: "1420000000",
		ExpirationTime:          1000000000,
		Fills: []orders.FillDTO{
			{
				TxSignature:              "fill-tx",
				FilledMakerAmount:        "500000000000000000",
				FilledAuctionTakerAmount: "710000000",
			},
		},
		CreatedAt:        1000000000,
		SrcTokenPriceUsd: 1.0,
		DstTokenPriceUsd: 1.0,
		Cancelable:       true,
	}

	status, err := OrderStatusFromJSON(json)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(status.Fills) != 1 {
		t.Errorf("Expected 1 fill, got %d", len(status.Fills))
	}
	if status.Fills[0].TxSignature != "fill-tx" {
		t.Errorf("Expected TxSignature 'fill-tx', got %s", status.Fills[0].TxSignature)
	}
	if status.Fills[0].FilledMakerAmount.Cmp(big.NewInt(500000000000000000)) != 0 {
		t.Errorf("Expected FilledMakerAmount 500000000000000000, got %s", status.Fills[0].FilledMakerAmount.String())
	}
}

func TestOrderStatusFromJSON_InvalidMaker(t *testing.T) {
	json := &orders.OrderStatusDTO{
		Maker:     "invalid-address",
		OrderHash: "test-hash",
		Status:    orders.OrderStatusInProgress,
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
		ApproximateTakingAmount: "1420000000",
		ExpirationTime:          1000000000,
		Fills:                   []orders.FillDTO{},
		CreatedAt:               1000000000,
		SrcTokenPriceUsd:        1.0,
		DstTokenPriceUsd:        1.0,
		Cancelable:              true,
	}

	_, err := OrderStatusFromJSON(json)
	if err == nil {
		t.Fatal("Expected error for invalid maker address")
	}
}

func TestOrderStatus_IsActive(t *testing.T) {
	status := &OrderStatus{
		Status: orders.OrderStatusInProgress,
	}
	if !status.IsActive() {
		t.Error("Expected order to be active")
	}

	status.Status = orders.OrderStatusFilled
	if status.IsActive() {
		t.Error("Expected order to not be active when filled")
	}

	status.Status = orders.OrderStatusCancelled
	if status.IsActive() {
		t.Error("Expected order to not be active when cancelled")
	}
}
