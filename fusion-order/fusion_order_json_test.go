package fusionorder

import (
	"math/big"
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
	"github.com/dawitel/solana-fusion-sdk-go/utils/time"
)

func TestFusionOrder_ToJSON_Success(t *testing.T) {
	orderInfo := OrderInfoData{
		ID:                 1,
		SrcAmount:          big.NewInt(1000000000000000000),
		MinDstAmount:       big.NewInt(1420000000),
		EstimatedDstAmount: big.NewInt(1420000000),
		Receiver:           domains.MustAddressFromString("11111111111111111111111111111111"),
		SrcMint:            domains.WRAPPED_NATIVE,
		DstMint:            domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
	}

	auctionDetails, err := NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []AuctionPoint
	}{
		StartTime:       uint32(time.Now()),
		Duration:        180,
		InitialRateBump: 0,
		Points:          []AuctionPoint{},
	})
	if err != nil {
		t.Fatalf("Failed to create auction details: %v", err)
	}

	order, err := NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *FeeConfig
		ResolverCancellationConfig *ResolverCancellationConfig
	}{
		SrcAssetIsNative:           false,
		DstAssetIsNative:           false,
		OrderExpirationDelay:       12,
		Fees:                       nil,
		ResolverCancellationConfig: AlmostZeroResolverCancellationConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	json := order.ToJSON()
	if json == nil {
		t.Fatal("Expected non-nil JSON")
	}
	if json.ID != 1 {
		t.Errorf("Expected ID 1, got %d", json.ID)
	}
	if json.SrcAmount != "1000000000000000000" {
		t.Errorf("Expected SrcAmount '1000000000000000000', got %s", json.SrcAmount)
	}
	if json.MinDstAmount != "1420000000" {
		t.Errorf("Expected MinDstAmount '1420000000', got %s", json.MinDstAmount)
	}
}

func TestFusionOrder_ToJSON_WithFees(t *testing.T) {
	orderInfo := OrderInfoData{
		ID:                 1,
		SrcAmount:          big.NewInt(1000000000000000000),
		MinDstAmount:       big.NewInt(1420000000),
		EstimatedDstAmount: big.NewInt(1420000000),
		Receiver:           domains.MustAddressFromString("11111111111111111111111111111111"),
		SrcMint:            domains.WRAPPED_NATIVE,
		DstMint:            domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
	}

	auctionDetails, err := NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []AuctionPoint
	}{
		StartTime:       uint32(time.Now()),
		Duration:        180,
		InitialRateBump: 0,
		Points:          []AuctionPoint{},
	})
	if err != nil {
		t.Fatalf("Failed to create auction details: %v", err)
	}

	protocolDstAta := domains.MustAddressFromString("44444444444444444444444444444444")
	protocolFee := domains.BpsFromPercent(1.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)
	fees, err := NewFeeConfig(protocolDstAta, nil, protocolFee, domains.ZeroBps, surplusShare)
	if err != nil {
		t.Fatalf("Failed to create fee config: %v", err)
	}

	order, err := NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *FeeConfig
		ResolverCancellationConfig *ResolverCancellationConfig
	}{
		SrcAssetIsNative:           false,
		DstAssetIsNative:           false,
		OrderExpirationDelay:       12,
		Fees:                       fees,
		ResolverCancellationConfig: AlmostZeroResolverCancellationConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	json := order.ToJSON()
	if json.Fee.ProtocolDstAta == nil {
		t.Error("Expected non-nil ProtocolDstAta")
	}
	if *json.Fee.ProtocolDstAta != protocolDstAta.ToString() {
		t.Errorf("Expected ProtocolDstAta %s, got %s", protocolDstAta.ToString(), *json.Fee.ProtocolDstAta)
	}
}

func TestFromJSON_Success(t *testing.T) {
	json := &FusionOrderJSON{
		ID:                         1,
		Receiver:                   "11111111111111111111111111111111",
		CancellationAuctionDuration: 100,
		SrcMint:                     "So11111111111111111111111111111111111111112",
		DstMint:                     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		SrcAmount:                   "1000000000000000000",
		MinDstAmount:                "1420000000",
		EstimatedDstAmount:         "1420000000",
		ExpirationTime:              1000000000,
		SrcAssetIsNative:            false,
		DstAssetIsNative:            false,
		Fee: FeeJSON{
			ProtocolFee:       0,
			IntegratorFee:     0,
			SurplusPercentage: 0,
			MaxCancellationPremium: "1000000",
		},
		DutchAuctionData: DutchAuctionDataJSON{
			StartTime:           1000000000,
			Duration:            180,
			InitialRateBump:     0,
			PointsAndTimeDeltas: []PointAndTimeDeltaJSON{},
		},
	}

	order, err := FromJSON(json)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if order == nil {
		t.Fatal("Expected non-nil order")
	}
	if order.ID() != 1 {
		t.Errorf("Expected ID 1, got %d", order.ID())
	}
	if order.SrcAmount().Cmp(big.NewInt(1000000000000000000)) != 0 {
		t.Errorf("Expected SrcAmount %s, got %s", "1000000000000000000", order.SrcAmount().String())
	}
}

func TestFromJSON_InvalidAddress(t *testing.T) {
	json := &FusionOrderJSON{
		ID:                         1,
		Receiver:                   "invalid-address",
		CancellationAuctionDuration: 100,
		SrcMint:                     "So11111111111111111111111111111111111111112",
		DstMint:                     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		SrcAmount:                   "1000000000000000000",
		MinDstAmount:                "1420000000",
		EstimatedDstAmount:         "1420000000",
		ExpirationTime:              1000000000,
		SrcAssetIsNative:            false,
		DstAssetIsNative:            false,
		Fee: FeeJSON{
			ProtocolFee:       0,
			IntegratorFee:     0,
			SurplusPercentage: 0,
			MaxCancellationPremium: "1000000",
		},
		DutchAuctionData: DutchAuctionDataJSON{
			StartTime:           1000000000,
			Duration:            180,
			InitialRateBump:     0,
			PointsAndTimeDeltas: []PointAndTimeDeltaJSON{},
		},
	}

	_, err := FromJSON(json)
	if err == nil {
		t.Fatal("Expected error for invalid address")
	}
}

func TestFromJSON_InvalidAmount(t *testing.T) {
	json := &FusionOrderJSON{
		ID:                         1,
		Receiver:                   "11111111111111111111111111111111",
		CancellationAuctionDuration: 100,
		SrcMint:                     "So11111111111111111111111111111111111111112",
		DstMint:                     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		SrcAmount:                   "invalid",
		MinDstAmount:                "1420000000",
		EstimatedDstAmount:         "1420000000",
		ExpirationTime:              1000000000,
		SrcAssetIsNative:            false,
		DstAssetIsNative:            false,
		Fee: FeeJSON{
			ProtocolFee:       0,
			IntegratorFee:     0,
			SurplusPercentage: 0,
			MaxCancellationPremium: "1000000",
		},
		DutchAuctionData: DutchAuctionDataJSON{
			StartTime:           1000000000,
			Duration:            180,
			InitialRateBump:     0,
			PointsAndTimeDeltas: []PointAndTimeDeltaJSON{},
		},
	}

	_, err := FromJSON(json)
	if err == nil {
		t.Fatal("Expected error for invalid amount")
	}
}

func TestFusionOrder_ToJSON_FromJSON_RoundTrip(t *testing.T) {
	orderInfo := OrderInfoData{
		ID:                 1,
		SrcAmount:          big.NewInt(1000000000000000000),
		MinDstAmount:       big.NewInt(1420000000),
		EstimatedDstAmount: big.NewInt(1420000000),
		Receiver:           domains.MustAddressFromString("11111111111111111111111111111111"),
		SrcMint:            domains.WRAPPED_NATIVE,
		DstMint:            domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
	}

	auctionDetails, err := NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []AuctionPoint
	}{
		StartTime:       uint32(time.Now()),
		Duration:        180,
		InitialRateBump: 0,
		Points:          []AuctionPoint{},
	})
	if err != nil {
		t.Fatalf("Failed to create auction details: %v", err)
	}

	originalOrder, err := NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *FeeConfig
		ResolverCancellationConfig *ResolverCancellationConfig
	}{
		SrcAssetIsNative:           false,
		DstAssetIsNative:           false,
		OrderExpirationDelay:       12,
		Fees:                       nil,
		ResolverCancellationConfig: AlmostZeroResolverCancellationConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	json := originalOrder.ToJSON()
	decodedOrder, err := FromJSON(json)
	if err != nil {
		t.Fatalf("Failed to decode from JSON: %v", err)
	}

	if decodedOrder.ID() != originalOrder.ID() {
		t.Errorf("Expected ID %d, got %d", originalOrder.ID(), decodedOrder.ID())
	}
	if decodedOrder.SrcAmount().Cmp(originalOrder.SrcAmount()) != 0 {
		t.Errorf("Expected SrcAmount %s, got %s", originalOrder.SrcAmount().String(), decodedOrder.SrcAmount().String())
	}
	if decodedOrder.MinDstAmount().Cmp(originalOrder.MinDstAmount()) != 0 {
		t.Errorf("Expected MinDstAmount %s, got %s", originalOrder.MinDstAmount().String(), decodedOrder.MinDstAmount().String())
	}
}
