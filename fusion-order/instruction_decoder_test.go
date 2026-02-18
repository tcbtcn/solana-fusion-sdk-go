package fusionorder

import (
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
	"github.com/tcbtcn/solana-fusion-sdk-go/idl"
	"github.com/tcbtcn/solana-fusion-sdk-go/types"
	"github.com/tcbtcn/solana-fusion-sdk-go/utils/time"
)

func TestFromCreateInstruction(t *testing.T) {
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

	// Create instruction manually (can't import contracts due to import cycle)
	config := originalOrder.Build()
	borshData, _ := config.SerializeBorsh()

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	programID := domains.MustAddressFromString(idl.FusionSwapProgramAddress)

	instruction := types.NewTransactionInstruction(
		programID,
		[]types.AccountMeta{
			{Pubkey: domains.SYSTEM_PROGRAM_ID, IsSigner: false, IsWritable: false},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: false}, // escrow
			{Pubkey: originalOrder.SrcMint(), IsSigner: false, IsWritable: false},                                           // src_mint (index 2)
			{Pubkey: domains.TOKEN_PROGRAM_ID, IsSigner: false, IsWritable: false},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},
			{Pubkey: maker, IsSigner: true, IsWritable: true},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},
			{Pubkey: originalOrder.DstMint(), IsSigner: false, IsWritable: false}, // dst_mint (index 7)
			{Pubkey: originalOrder.Receiver(), IsSigner: false, IsWritable: true}, // receiver (index 8)
			{Pubkey: domains.ASSOCIATED_TOKEN_PROGRAM_ID, IsSigner: false, IsWritable: false},
			{Pubkey: programID, IsSigner: false, IsWritable: false}, // protocol_dst_ata (index 10, optional = programID)
			{Pubkey: programID, IsSigner: false, IsWritable: false}, // integrator_dst_ata (index 11, optional = programID)
		},
		append(idl.CreateOrderDiscriminator, borshData...),
	)

	// Decode instruction
	decodedOrder, err := FromCreateInstruction(instruction)
	if err != nil {
		t.Fatalf("Failed to decode instruction: %v", err)
	}

	// Verify decoded order matches original
	if decodedOrder.ID() != originalOrder.ID() {
		t.Errorf("Expected ID %d, got %d", originalOrder.ID(), decodedOrder.ID())
	}
	if decodedOrder.SrcAmount().Cmp(originalOrder.SrcAmount()) != 0 {
		t.Errorf("Expected srcAmount %s, got %s", originalOrder.SrcAmount().String(), decodedOrder.SrcAmount().String())
	}
	if decodedOrder.MinDstAmount().Cmp(originalOrder.MinDstAmount()) != 0 {
		t.Errorf("Expected minDstAmount %s, got %s", originalOrder.MinDstAmount().String(), decodedOrder.MinDstAmount().String())
	}
}

func TestFromFillInstruction(t *testing.T) {
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

	config := originalOrder.Build()
	borshData, _ := config.SerializeBorsh()

	programID := domains.MustAddressFromString(idl.FusionSwapProgramAddress)
	fillAmount := big.NewInt(100)

	// Instruction data structure: [discriminator][order config][amount]
	// Encode amount as u64 little-endian
	fillAmountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(fillAmountBytes, fillAmount.Uint64())

	instruction := types.NewTransactionInstruction(
		programID,
		[]types.AccountMeta{
			{Pubkey: domains.SYSTEM_PROGRAM_ID, IsSigner: false, IsWritable: false},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: false},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: false},
			{Pubkey: originalOrder.Receiver(), IsSigner: false, IsWritable: true},
			{Pubkey: originalOrder.SrcMint(), IsSigner: false, IsWritable: false},
			{Pubkey: originalOrder.DstMint(), IsSigner: false, IsWritable: false},
			{Pubkey: domains.TOKEN_PROGRAM_ID, IsSigner: false, IsWritable: false},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},
			{Pubkey: programID, IsSigner: false, IsWritable: false},
			{Pubkey: programID, IsSigner: false, IsWritable: false},
		},
		append(append(idl.FillOrderDiscriminator, borshData...), fillAmountBytes...),
	)

	decodedOrder, err := FromFillInstruction(instruction)
	if err != nil {
		t.Fatalf("Failed to decode instruction: %v", err)
	}

	if decodedOrder.ID() != originalOrder.ID() {
		t.Errorf("Expected ID %d, got %d", originalOrder.ID(), decodedOrder.ID())
	}
}

func TestFromFillInstruction_InvalidDiscriminator(t *testing.T) {
	programID := domains.MustAddressFromString(idl.FusionSwapProgramAddress)
	instruction := types.NewTransactionInstruction(
		programID,
		[]types.AccountMeta{},
		[]byte{1, 2, 3, 4, 5, 6, 7, 8},
	)

	_, err := FromFillInstruction(instruction)
	if err == nil {
		t.Fatal("Expected error for invalid discriminator")
	}
}

func TestFromFillInstruction_InsufficientAccounts(t *testing.T) {
	programID := domains.MustAddressFromString(idl.FusionSwapProgramAddress)
	instruction := types.NewTransactionInstruction(
		programID,
		[]types.AccountMeta{},
		append(idl.FillOrderDiscriminator, make([]byte, 100)...),
	)

	_, err := FromFillInstruction(instruction)
	if err == nil {
		t.Fatal("Expected error for insufficient accounts")
	}
}

func TestFromResolverCancelInstruction(t *testing.T) {
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

	maxCancellationPremium := big.NewInt(1000000)
	resolverConfig, err := NewResolverCancellationConfig(maxCancellationPremium, 100)
	if err != nil {
		t.Fatalf("Failed to create resolver config: %v", err)
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
		ResolverCancellationConfig: resolverConfig,
	})
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	config := originalOrder.Build()
	borshData, _ := config.SerializeBorsh()

	maker := domains.MustAddressFromString("11111111111111111111111111111111")
	resolver := domains.MustAddressFromString("11111111111111111111111111111113")
	programID := domains.MustAddressFromString(idl.FusionSwapProgramAddress)
	rewardLimit := big.NewInt(500000)

	// Instruction data structure: [discriminator][order config][rewardLimit]
	// Encode rewardLimit as u64 little-endian
	rewardLimitBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(rewardLimitBytes, rewardLimit.Uint64())

	instruction := types.NewTransactionInstruction(
		programID,
		[]types.AccountMeta{
			{Pubkey: resolver, IsSigner: true, IsWritable: true},                                                            // 0: resolver
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: false}, // 1: resolverAccess
			{Pubkey: maker, IsSigner: false, IsWritable: true},                                                              // 2: maker
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: false}, // 3: makerReceiver
			{Pubkey: originalOrder.SrcMint(), IsSigner: false, IsWritable: false},                                           // 4: srcMint
			{Pubkey: originalOrder.DstMint(), IsSigner: false, IsWritable: false},                                           // 5: dstMint
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},  // 6: escrow
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},  // 7: escrowSrcAta
			{Pubkey: domains.MustAddressFromString("11111111111111111111111111111111"), IsSigner: false, IsWritable: true},  // 8: makerSrcAta (optional)
			{Pubkey: domains.TOKEN_PROGRAM_ID, IsSigner: false, IsWritable: false},                                          // 9: srcTokenProgram
			{Pubkey: domains.SYSTEM_PROGRAM_ID, IsSigner: false, IsWritable: false},                                         // 10: systemProgram
			{Pubkey: programID, IsSigner: false, IsWritable: false},                                                         // 11: protocolDstAta (optional, use programID when nil)
			{Pubkey: programID, IsSigner: false, IsWritable: false},                                                         // 12: integratorDstAta (optional, use programID when nil)
		},
		append(append(idl.CancelOrderByResolverDiscriminator, borshData...), rewardLimitBytes...),
	)

	decodedOrder, err := FromResolverCancelInstruction(instruction)
	if err != nil {
		t.Fatalf("Failed to decode instruction: %v", err)
	}

	if decodedOrder.ID() != originalOrder.ID() {
		t.Errorf("Expected ID %d, got %d", originalOrder.ID(), decodedOrder.ID())
	}
}

func TestFromResolverCancelInstruction_InvalidDiscriminator(t *testing.T) {
	programID := domains.MustAddressFromString(idl.FusionSwapProgramAddress)
	instruction := types.NewTransactionInstruction(
		programID,
		[]types.AccountMeta{},
		[]byte{1, 2, 3, 4, 5, 6, 7, 8},
	)

	_, err := FromResolverCancelInstruction(instruction)
	if err == nil {
		t.Fatal("Expected error for invalid discriminator")
	}
}
