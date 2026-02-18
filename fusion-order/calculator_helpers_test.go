package fusionorder

import (
	"math/big"
	"testing"

	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
	"github.com/tcbtcn/solana-fusion-sdk-go/utils/time"
)

type mockCalculator struct {
	getRequiredTakingAmount func(takingAmount *big.Int, time uint32) *big.Int
	getUserReceiveAmount    func(takingAmount *big.Int, estimatedTakingAmount *big.Int, time uint32) *big.Int
	getIntegratorFee        func(takingAmount *big.Int, time uint32) *big.Int
	getProtocolFee          func(takingAmount *big.Int, estimatedTakingAmount *big.Int, time uint32) *big.Int
}

func (m *mockCalculator) GetRequiredTakingAmount(takingAmount *big.Int, time uint32) *big.Int {
	if m.getRequiredTakingAmount != nil {
		return m.getRequiredTakingAmount(takingAmount, time)
	}
	return takingAmount
}

func (m *mockCalculator) GetUserReceiveAmount(takingAmount *big.Int, estimatedTakingAmount *big.Int, time uint32) *big.Int {
	if m.getUserReceiveAmount != nil {
		return m.getUserReceiveAmount(takingAmount, estimatedTakingAmount, time)
	}
	return takingAmount
}

func (m *mockCalculator) GetIntegratorFee(takingAmount *big.Int, time uint32) *big.Int {
	if m.getIntegratorFee != nil {
		return m.getIntegratorFee(takingAmount, time)
	}
	return big.NewInt(0)
}

func (m *mockCalculator) GetProtocolFee(takingAmount *big.Int, estimatedTakingAmount *big.Int, time uint32) *big.Int {
	if m.getProtocolFee != nil {
		return m.getProtocolFee(takingAmount, estimatedTakingAmount, time)
	}
	return big.NewInt(0)
}

func TestFusionOrder_CalcTakingAmountWithCalculator(t *testing.T) {
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

	makingAmount := big.NewInt(500000000000000000)
	calculator := &mockCalculator{
		getRequiredTakingAmount: func(takingAmount *big.Int, time uint32) *big.Int {
			return new(big.Int).Mul(takingAmount, big.NewInt(2))
		},
	}

	result := order.CalcTakingAmountWithCalculator(calculator, makingAmount, uint32(time.Now()))
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if result.Cmp(big.NewInt(0)) <= 0 {
		t.Error("Expected positive result")
	}
}

func TestFusionOrder_GetUserReceiveAmountWithCalculator(t *testing.T) {
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

	makingAmount := big.NewInt(500000000000000000)
	calculator := &mockCalculator{
		getUserReceiveAmount: func(takingAmount *big.Int, estimatedTakingAmount *big.Int, time uint32) *big.Int {
			return new(big.Int).Sub(takingAmount, big.NewInt(100))
		},
	}

	result := order.GetUserReceiveAmountWithCalculator(calculator, makingAmount, uint32(time.Now()))
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
}

func TestFusionOrder_GetIntegratorFeeWithCalculator(t *testing.T) {
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

	calculator := &mockCalculator{
		getIntegratorFee: func(takingAmount *big.Int, time uint32) *big.Int {
			return big.NewInt(20)
		},
	}

	result := order.GetIntegratorFeeWithCalculator(calculator, uint32(time.Now()), nil)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if result.Cmp(big.NewInt(20)) != 0 {
		t.Errorf("Expected integrator fee 20, got %s", result.String())
	}
}

func TestFusionOrder_GetProtocolFeeWithCalculator(t *testing.T) {
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

	calculator := &mockCalculator{
		getProtocolFee: func(takingAmount *big.Int, estimatedTakingAmount *big.Int, time uint32) *big.Int {
			return big.NewInt(10)
		},
	}

	result := order.GetProtocolFeeWithCalculator(calculator, uint32(time.Now()), nil)
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if result.Cmp(big.NewInt(10)) != 0 {
		t.Errorf("Expected protocol fee 10, got %s", result.String())
	}
}

func TestFusionOrder_CalcTakingAmount(t *testing.T) {
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

	makingAmount := big.NewInt(500000000000000000)
	result, err := order.CalcTakingAmount(makingAmount, uint32(time.Now()))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if result.Cmp(big.NewInt(0)) <= 0 {
		t.Error("Expected positive result")
	}

	// Verify it matches the WithCalculator version
	calculator := order.GetCalculator()
	expected := order.CalcTakingAmountWithCalculator(calculator, makingAmount, uint32(time.Now()))
	if result.Cmp(expected) != 0 {
		t.Errorf("Expected result to match WithCalculator version: got %s, expected %s", result.String(), expected.String())
	}
}

func TestFusionOrder_GetUserReceiveAmount(t *testing.T) {
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

	makingAmount := big.NewInt(500000000000000000)
	result, err := order.GetUserReceiveAmount(makingAmount, uint32(time.Now()))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Verify it matches the WithCalculator version
	calculator := order.GetCalculator()
	expected := order.GetUserReceiveAmountWithCalculator(calculator, makingAmount, uint32(time.Now()))
	if result.Cmp(expected) != 0 {
		t.Errorf("Expected result to match WithCalculator version: got %s, expected %s", result.String(), expected.String())
	}
}

func TestFusionOrder_GetIntegratorFee(t *testing.T) {
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

	result, err := order.GetIntegratorFee(uint32(time.Now()), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Verify it matches the WithCalculator version
	calculator := order.GetCalculator()
	expected := order.GetIntegratorFeeWithCalculator(calculator, uint32(time.Now()), nil)
	if result.Cmp(expected) != 0 {
		t.Errorf("Expected result to match WithCalculator version: got %s, expected %s", result.String(), expected.String())
	}
}

func TestFusionOrder_GetProtocolFee(t *testing.T) {
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

	result, err := order.GetProtocolFee(uint32(time.Now()), nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Verify it matches the WithCalculator version
	calculator := order.GetCalculator()
	expected := order.GetProtocolFeeWithCalculator(calculator, uint32(time.Now()), nil)
	if result.Cmp(expected) != 0 {
		t.Errorf("Expected result to match WithCalculator version: got %s, expected %s", result.String(), expected.String())
	}
}
