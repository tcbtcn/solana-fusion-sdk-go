package contracts

import (
	"math/big"
	"testing"

	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
	fusionorder "github.com/tcbtcn/solana-fusion-sdk-go/fusion-order"
	"github.com/tcbtcn/solana-fusion-sdk-go/utils/time"
)

func TestFusionSwapContract_Create(t *testing.T) {
	contract := DefaultFusionSwapContract()

	orderInfo := fusionorder.OrderInfoData{
		ID:                 1,
		SrcAmount:          big.NewInt(1000000000000000000),
		MinDstAmount:       big.NewInt(1420000000),
		EstimatedDstAmount: big.NewInt(1420000000),
		Receiver:           domains.MustAddressFromString("11111111111111111111111111111111"),
		SrcMint:            domains.WRAPPED_NATIVE,
		DstMint:            domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
	}

	auctionDetails, err := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       uint32(time.Now()),
		Duration:        180,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	if err != nil {
		t.Fatalf("Failed to create auction details: %v", err)
	}

	order, err := fusionorder.NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *fusionorder.FeeConfig
		ResolverCancellationConfig *fusionorder.ResolverCancellationConfig
	}{
		SrcAssetIsNative:           false,
		DstAssetIsNative:           false,
		OrderExpirationDelay:       12,
		Fees:                       nil,
		ResolverCancellationConfig: fusionorder.AlmostZeroResolverCancellationConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	instruction, err := contract.Create(order, struct {
		Maker           *domains.Address
		SrcTokenProgram *domains.Address
	}{
		Maker:           maker,
		SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
	})
	if err != nil {
		t.Fatalf("Failed to create instruction: %v", err)
	}

	if instruction.ProgramID == nil {
		t.Error("Expected non-nil program ID")
	}
	if len(instruction.Accounts) == 0 {
		t.Error("Expected non-empty accounts")
	}
	if len(instruction.Data) == 0 {
		t.Error("Expected non-empty instruction data")
	}

	// Verify discriminator is correct
	if len(instruction.Data) < 8 {
		t.Error("Expected instruction data to have at least 8 bytes for discriminator")
	}
}

func TestFusionSwapContract_Fill(t *testing.T) {
	contract := DefaultFusionSwapContract()

	orderInfo := fusionorder.OrderInfoData{
		ID:                 1,
		SrcAmount:          big.NewInt(1000000000000000000),
		MinDstAmount:       big.NewInt(1420000000),
		EstimatedDstAmount: big.NewInt(1420000000),
		Receiver:           domains.MustAddressFromString("11111111111111111111111111111111"),
		SrcMint:            domains.WRAPPED_NATIVE,
		DstMint:            domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
	}

	auctionDetails, err := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       uint32(time.Now()),
		Duration:        180,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	if err != nil {
		t.Fatalf("Failed to create auction details: %v", err)
	}

	order, err := fusionorder.NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *fusionorder.FeeConfig
		ResolverCancellationConfig *fusionorder.ResolverCancellationConfig
	}{
		SrcAssetIsNative:           false,
		DstAssetIsNative:           false,
		OrderExpirationDelay:       12,
		Fees:                       nil,
		ResolverCancellationConfig: fusionorder.AlmostZeroResolverCancellationConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	taker := domains.MustAddressFromString("11111111111111111111111111111112")
	fillAmount := big.NewInt(100)

	instruction, err := contract.Fill(order, fillAmount, struct {
		Taker           *domains.Address
		Maker           *domains.Address
		SrcTokenProgram *domains.Address
		DstTokenProgram *domains.Address
		TakerSrcAccount *domains.Address
		Whitelist       *domains.Address
	}{
		Taker:           taker,
		Maker:           maker,
		SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
		DstTokenProgram: domains.TOKEN_PROGRAM_ID,
		TakerSrcAccount: nil,
		Whitelist:       nil,
	})
	if err != nil {
		t.Fatalf("Failed to create fill instruction: %v", err)
	}

	if instruction.ProgramID == nil {
		t.Error("Expected non-nil program ID")
	}
	if len(instruction.Accounts) == 0 {
		t.Error("Expected non-empty accounts")
	}
	if len(instruction.Data) == 0 {
		t.Error("Expected non-empty instruction data")
	}
}

func TestFusionSwapContract_CancelOwnOrder(t *testing.T) {
	contract := DefaultFusionSwapContract()

	orderInfo := fusionorder.OrderInfoData{
		ID:                 1,
		SrcAmount:          big.NewInt(1000000000000000000),
		MinDstAmount:       big.NewInt(1420000000),
		EstimatedDstAmount: big.NewInt(1420000000),
		Receiver:           domains.MustAddressFromString("11111111111111111111111111111111"),
		SrcMint:            domains.WRAPPED_NATIVE,
		DstMint:            domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
	}

	auctionDetails, err := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       uint32(time.Now()),
		Duration:        180,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	if err != nil {
		t.Fatalf("Failed to create auction details: %v", err)
	}

	order, err := fusionorder.NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *fusionorder.FeeConfig
		ResolverCancellationConfig *fusionorder.ResolverCancellationConfig
	}{
		SrcAssetIsNative:           false,
		DstAssetIsNative:           false,
		OrderExpirationDelay:       12,
		Fees:                       nil,
		ResolverCancellationConfig: fusionorder.AlmostZeroResolverCancellationConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	instruction, err := contract.CancelOwnOrder(order, struct {
		Maker           *domains.Address
		SrcTokenProgram *domains.Address
	}{
		Maker:           maker,
		SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
	})
	if err != nil {
		t.Fatalf("Failed to create cancel instruction: %v", err)
	}

	if instruction.ProgramID == nil {
		t.Error("Expected non-nil program ID")
	}
	if len(instruction.Accounts) == 0 {
		t.Error("Expected non-empty accounts")
	}
	if len(instruction.Data) == 0 {
		t.Error("Expected non-empty instruction data")
	}
}

func TestFusionSwapContract_CancelOrderByResolver(t *testing.T) {
	contract := DefaultFusionSwapContract()

	orderInfo := fusionorder.OrderInfoData{
		ID:                 1,
		SrcAmount:          big.NewInt(1000000000000000000),
		MinDstAmount:       big.NewInt(1420000000),
		EstimatedDstAmount: big.NewInt(1420000000),
		Receiver:           domains.MustAddressFromString("11111111111111111111111111111111"),
		SrcMint:            domains.WRAPPED_NATIVE,
		DstMint:            domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
	}

	auctionDetails, err := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       uint32(time.Now()),
		Duration:        180,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	if err != nil {
		t.Fatalf("Failed to create auction details: %v", err)
	}

	maxCancellationPremium := big.NewInt(1000000)
	resolverConfig, err := fusionorder.NewResolverCancellationConfig(maxCancellationPremium, 100)
	if err != nil {
		t.Fatalf("Failed to create resolver config: %v", err)
	}

	order, err := fusionorder.NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *fusionorder.FeeConfig
		ResolverCancellationConfig *fusionorder.ResolverCancellationConfig
	}{
		SrcAssetIsNative:           false,
		DstAssetIsNative:           false,
		OrderExpirationDelay:       12,
		Fees:                       nil,
		ResolverCancellationConfig: resolverConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	resolver := domains.MustAddressFromString("11111111111111111111111111111113")
	rewardLimit := big.NewInt(500000)

	instruction, err := contract.CancelOrderByResolver(order, struct {
		Maker           *domains.Address
		Resolver        *domains.Address
		SrcTokenProgram *domains.Address
		Whitelist       *domains.Address
	}{
		Maker:           maker,
		Resolver:        resolver,
		SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
		Whitelist:       nil,
	}, rewardLimit)
	if err != nil {
		t.Fatalf("Failed to create cancel by resolver instruction: %v", err)
	}

	if instruction.ProgramID == nil {
		t.Error("Expected non-nil program ID")
	}
	if len(instruction.Accounts) == 0 {
		t.Error("Expected non-empty accounts")
	}
	if len(instruction.Data) == 0 {
		t.Error("Expected non-empty instruction data")
	}
}

func TestFusionSwapContract_CancelOrderByResolver_NoResolverConfig(t *testing.T) {
	contract := DefaultFusionSwapContract()

	orderInfo := fusionorder.OrderInfoData{
		ID:                 1,
		SrcAmount:          big.NewInt(1000000000000000000),
		MinDstAmount:       big.NewInt(1420000000),
		EstimatedDstAmount: big.NewInt(1420000000),
		Receiver:           domains.MustAddressFromString("11111111111111111111111111111111"),
		SrcMint:            domains.WRAPPED_NATIVE,
		DstMint:            domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
	}

	auctionDetails, err := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       uint32(time.Now()),
		Duration:        180,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	if err != nil {
		t.Fatalf("Failed to create auction details: %v", err)
	}

	order, err := fusionorder.NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *fusionorder.FeeConfig
		ResolverCancellationConfig *fusionorder.ResolverCancellationConfig
	}{
		SrcAssetIsNative:           false,
		DstAssetIsNative:           false,
		OrderExpirationDelay:       12,
		Fees:                       nil,
		ResolverCancellationConfig: fusionorder.ZeroResolverCancellationConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	resolver := domains.MustAddressFromString("11111111111111111111111111111113")

	_, err = contract.CancelOrderByResolver(order, struct {
		Maker           *domains.Address
		Resolver        *domains.Address
		SrcTokenProgram *domains.Address
		Whitelist       *domains.Address
	}{
		Maker:           maker,
		Resolver:        resolver,
		SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
		Whitelist:       nil,
	}, nil)
	if err == nil {
		t.Fatal("Expected error for order without resolver config")
	}
}

func TestFusionSwapContract_Create_WithFees(t *testing.T) {
	contract := DefaultFusionSwapContract()

	orderInfo := fusionorder.OrderInfoData{
		ID:                 1,
		SrcAmount:          big.NewInt(1000000000000000000),
		MinDstAmount:       big.NewInt(1420000000),
		EstimatedDstAmount: big.NewInt(1420000000),
		Receiver:           domains.MustAddressFromString("11111111111111111111111111111111"),
		SrcMint:            domains.WRAPPED_NATIVE,
		DstMint:            domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
	}

	auctionDetails, err := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       uint32(time.Now()),
		Duration:        180,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	if err != nil {
		t.Fatalf("Failed to create auction details: %v", err)
	}

	protocolDstAta := domains.MustAddressFromString("11111111111111111111111111111114")
	protocolFee := domains.BpsFromPercent(1.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)
	fees, err := fusionorder.NewFeeConfig(protocolDstAta, nil, protocolFee, domains.ZeroBps, surplusShare)
	if err != nil {
		t.Fatalf("Failed to create fee config: %v", err)
	}

	order, err := fusionorder.NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *fusionorder.FeeConfig
		ResolverCancellationConfig *fusionorder.ResolverCancellationConfig
	}{
		SrcAssetIsNative:           false,
		DstAssetIsNative:           false,
		OrderExpirationDelay:       12,
		Fees:                       fees,
		ResolverCancellationConfig: fusionorder.AlmostZeroResolverCancellationConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	instruction, err := contract.Create(order, struct {
		Maker           *domains.Address
		SrcTokenProgram *domains.Address
	}{
		Maker:           maker,
		SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
	})
	if err != nil {
		t.Fatalf("Failed to create instruction: %v", err)
	}

	if len(instruction.Accounts) < 12 {
		t.Errorf("Expected at least 12 accounts with fees, got %d", len(instruction.Accounts))
	}
}

func TestFusionSwapContract_Create_WithNativeSrc(t *testing.T) {
	contract := DefaultFusionSwapContract()

	orderInfo := fusionorder.OrderInfoData{
		ID:                 1,
		SrcAmount:          big.NewInt(1000000000000000000),
		MinDstAmount:       big.NewInt(1420000000),
		EstimatedDstAmount: big.NewInt(1420000000),
		Receiver:           domains.MustAddressFromString("11111111111111111111111111111111"),
		SrcMint:            domains.NATIVE,
		DstMint:            domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
	}

	auctionDetails, err := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       uint32(time.Now()),
		Duration:        180,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	if err != nil {
		t.Fatalf("Failed to create auction details: %v", err)
	}

	order, err := fusionorder.NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *fusionorder.FeeConfig
		ResolverCancellationConfig *fusionorder.ResolverCancellationConfig
	}{
		SrcAssetIsNative:           true,
		DstAssetIsNative:           false,
		OrderExpirationDelay:       12,
		Fees:                       nil,
		ResolverCancellationConfig: fusionorder.AlmostZeroResolverCancellationConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	instruction, err := contract.Create(order, struct {
		Maker           *domains.Address
		SrcTokenProgram *domains.Address
	}{
		Maker:           maker,
		SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
	})
	if err != nil {
		t.Fatalf("Failed to create instruction: %v", err)
	}

	if len(instruction.Accounts) < 11 {
		t.Errorf("Expected at least 11 accounts for native src, got %d", len(instruction.Accounts))
	}
}

func TestFusionSwapContract_Fill_WithNativeDst(t *testing.T) {
	contract := DefaultFusionSwapContract()

	orderInfo := fusionorder.OrderInfoData{
		ID:                 1,
		SrcAmount:          big.NewInt(1000000000000000000),
		MinDstAmount:       big.NewInt(1420000000),
		EstimatedDstAmount: big.NewInt(1420000000),
		Receiver:           domains.MustAddressFromString("11111111111111111111111111111111"),
		SrcMint:            domains.WRAPPED_NATIVE,
		DstMint:            domains.NATIVE,
	}

	auctionDetails, err := fusionorder.NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []fusionorder.AuctionPoint
	}{
		StartTime:       uint32(time.Now()),
		Duration:        180,
		InitialRateBump: 0,
		Points:          []fusionorder.AuctionPoint{},
	})
	if err != nil {
		t.Fatalf("Failed to create auction details: %v", err)
	}

	order, err := fusionorder.NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *fusionorder.FeeConfig
		ResolverCancellationConfig *fusionorder.ResolverCancellationConfig
	}{
		SrcAssetIsNative:           false,
		DstAssetIsNative:           true,
		OrderExpirationDelay:       12,
		Fees:                       nil,
		ResolverCancellationConfig: fusionorder.AlmostZeroResolverCancellationConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	taker := domains.MustAddressFromString("11111111111111111111111111111112")
	fillAmount := big.NewInt(100)

	instruction, err := contract.Fill(order, fillAmount, struct {
		Taker           *domains.Address
		Maker           *domains.Address
		SrcTokenProgram *domains.Address
		DstTokenProgram *domains.Address
		TakerSrcAccount *domains.Address
		Whitelist       *domains.Address
	}{
		Taker:           taker,
		Maker:           maker,
		SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
		DstTokenProgram: domains.TOKEN_PROGRAM_ID,
		TakerSrcAccount: nil,
		Whitelist:       nil,
	})
	if err != nil {
		t.Fatalf("Failed to create fill instruction: %v", err)
	}

	if len(instruction.Accounts) < 15 {
		t.Errorf("Expected at least 15 accounts for native dst, got %d", len(instruction.Accounts))
	}
}
