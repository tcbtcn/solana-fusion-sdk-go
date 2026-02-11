package fusionorder

import (
	"math/big"
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
	"github.com/dawitel/solana-fusion-sdk-go/utils/time"
)

func TestFusionOrder_GetOrderHash(t *testing.T) {
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

	hash := order.GetOrderHash()
	if len(hash) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(hash))
	}

	hashBase58 := order.GetOrderHashBase58()
	if hashBase58 == "" {
		t.Error("Expected non-empty base58 hash")
	}
}

func TestFusionOrder_Build(t *testing.T) {
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

	config := order.Build()
	if config.ID != 1 {
		t.Errorf("Expected ID 1, got %d", config.ID)
	}
	if config.SrcAmount.Cmp(big.NewInt(1000000000000000000)) != 0 {
		t.Errorf("Expected srcAmount %s, got %s", "1000000000000000000", config.SrcAmount.String())
	}
}

func TestFusionOrder_SerializeBorsh(t *testing.T) {
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

	config := order.Build()
	borshData, err := config.SerializeBorsh()
	if err != nil {
		t.Fatalf("Failed to serialize Borsh: %v", err)
	}

	if len(borshData) == 0 {
		t.Error("Expected non-empty Borsh data")
	}

	// Verify minimum expected size (at least id + amounts + times + flags)
	if len(borshData) < 30 {
		t.Errorf("Expected Borsh data to be at least 30 bytes, got %d", len(borshData))
	}
}

func TestFusionOrder_GetEscrow(t *testing.T) {
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

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	escrow, err := order.GetEscrow(maker, nil, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if escrow == nil {
		t.Fatal("Expected non-nil escrow")
	}
	if len(escrow.ToBuffer()) != 32 {
		t.Errorf("Expected 32-byte address, got %d bytes", len(escrow.ToBuffer()))
	}
}

func TestFusionOrder_GetEscrow_WithProgramId(t *testing.T) {
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

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	customProgramId := domains.MustAddressFromString("HNarfxC3kYMMhFkxUFeYb8wHVdPzY5t9pupqW5fL2meM")
	escrow, err := order.GetEscrow(maker, nil, customProgramId)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if escrow == nil {
		t.Fatal("Expected non-nil escrow")
	}
}

func TestFusionOrder_AccessorMethods(t *testing.T) {
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

	if order.ID() != 1 {
		t.Errorf("Expected ID 1, got %d", order.ID())
	}
	if !order.SrcMint().Equal(domains.WRAPPED_NATIVE) {
		t.Error("Expected SrcMint to match")
	}
	if !order.DstMint().Equal(domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")) {
		t.Error("Expected DstMint to match")
	}
	if !order.Receiver().Equal(domains.MustAddressFromString("11111111111111111111111111111111")) {
		t.Error("Expected Receiver to match")
	}
	if order.SrcAmount().Cmp(big.NewInt(1000000000000000000)) != 0 {
		t.Errorf("Expected SrcAmount %s, got %s", "1000000000000000000", order.SrcAmount().String())
	}
	if order.MinDstAmount().Cmp(big.NewInt(1420000000)) != 0 {
		t.Errorf("Expected MinDstAmount %s, got %s", "1420000000", order.MinDstAmount().String())
	}
	if order.EstimatedDstAmount().Cmp(big.NewInt(1420000000)) != 0 {
		t.Errorf("Expected EstimatedDstAmount %s, got %s", "1420000000", order.EstimatedDstAmount().String())
	}
	if order.SrcAssetIsNative() {
		t.Error("Expected SrcAssetIsNative to be false")
	}
	if order.DstAssetIsNative() {
		t.Error("Expected DstAssetIsNative to be false")
	}
	if order.Fees() != nil {
		t.Error("Expected Fees to be nil")
	}
	if order.ResolverCancellationConfig() == nil {
		t.Error("Expected non-nil ResolverCancellationConfig")
	}
}

func TestFusionOrder_AccessorMethods_WithFees(t *testing.T) {
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

	protocolDstAta := domains.MustAddressFromString("11111111111111111111111111111114")
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

	if order.Fees() == nil {
		t.Error("Expected non-nil Fees")
	}
	if order.Fees().ProtocolDstAta == nil {
		t.Error("Expected non-nil ProtocolDstAta")
	}
}
