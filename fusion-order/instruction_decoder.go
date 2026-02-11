package fusionorder

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
	"github.com/dawitel/solana-fusion-sdk-go/idl"
	"github.com/dawitel/solana-fusion-sdk-go/types"
)

// FromCreateInstruction decodes a FusionOrder from a create instruction
func FromCreateInstruction(ix *types.TransactionInstruction) (*FusionOrder, error) {
	// Check discriminator
	if len(ix.Data) < 8 {
		return nil, errors.New("instruction data too short")
	}
	discriminator := ix.Data[:8]
	if !equalBytes(discriminator, idl.CreateOrderDiscriminator) {
		return nil, fmt.Errorf("invalid instruction discriminator: expected create instruction")
	}

	// Decode order config from Borsh data (skip 8-byte discriminator)
	orderConfig, err := deserializeOrderConfig(ix.Data[8:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode order config: %w", err)
	}

	// Extract accounts (based on create instruction account layout)
	if len(ix.Accounts) < 12 {
		return nil, errors.New("insufficient accounts in instruction")
	}

	srcMint := ix.Accounts[2].Pubkey
	dstMint := ix.Accounts[7].Pubkey
	receiver := ix.Accounts[8].Pubkey
	protocolDstAta := ix.Accounts[10].Pubkey
	integratorDstAta := ix.Accounts[11].Pubkey

	return fromContractOrder(orderConfig, struct {
		SrcMint          *domains.Address
		DstMint          *domains.Address
		Receiver         *domains.Address
		ProtocolDstAta   *domains.Address
		IntegratorDstAta *domains.Address
		ProgramID        *domains.Address
	}{
		SrcMint:          srcMint,
		DstMint:          dstMint,
		Receiver:         receiver,
		ProtocolDstAta:   protocolDstAta,
		IntegratorDstAta: integratorDstAta,
		ProgramID:        ix.ProgramID,
	})
}

// FromFillInstruction decodes a FusionOrder from a fill instruction
func FromFillInstruction(ix *types.TransactionInstruction) (*FusionOrder, error) {
	// Check discriminator
	if len(ix.Data) < 8 {
		return nil, errors.New("instruction data too short")
	}
	discriminator := ix.Data[:8]
	if !equalBytes(discriminator, idl.FillOrderDiscriminator) {
		return nil, fmt.Errorf("invalid instruction discriminator: expected fill instruction")
	}

	// Decode order config from Borsh data (skip 8-byte discriminator, then skip 8-byte amount)
	if len(ix.Data) < 16 {
		return nil, errors.New("instruction data too short")
	}
	orderConfig, err := deserializeOrderConfig(ix.Data[8:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode order config: %w", err)
	}

	// Extract accounts (based on fill instruction account layout)
	if len(ix.Accounts) < 17 {
		return nil, errors.New("insufficient accounts in instruction")
	}

	srcMint := ix.Accounts[4].Pubkey   // index 4
	dstMint := ix.Accounts[5].Pubkey   // index 5
	receiver := ix.Accounts[3].Pubkey  // index 3
	protocolDstAta := ix.Accounts[15].Pubkey // index 15
	integratorDstAta := ix.Accounts[16].Pubkey // index 16

	return fromContractOrder(orderConfig, struct {
		SrcMint          *domains.Address
		DstMint          *domains.Address
		Receiver         *domains.Address
		ProtocolDstAta   *domains.Address
		IntegratorDstAta *domains.Address
		ProgramID        *domains.Address
	}{
		SrcMint:          srcMint,
		DstMint:          dstMint,
		Receiver:         receiver,
		ProtocolDstAta:   protocolDstAta,
		IntegratorDstAta: integratorDstAta,
		ProgramID:        ix.ProgramID,
	})
}

// FromResolverCancelInstruction decodes a FusionOrder from a cancelByResolver instruction
func FromResolverCancelInstruction(ix *types.TransactionInstruction) (*FusionOrder, error) {
	// Check discriminator
	if len(ix.Data) < 8 {
		return nil, errors.New("instruction data too short")
	}
	discriminator := ix.Data[:8]
	if !equalBytes(discriminator, idl.CancelOrderByResolverDiscriminator) {
		return nil, fmt.Errorf("invalid instruction discriminator: expected cancelByResolver instruction")
	}

	// Decode order config from Borsh data (skip 8-byte discriminator, then skip 8-byte rewardLimit)
	if len(ix.Data) < 16 {
		return nil, errors.New("instruction data too short")
	}
	orderConfig, err := deserializeOrderConfig(ix.Data[8:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode order config: %w", err)
	}

	// Extract accounts (based on cancelByResolver instruction account layout)
	if len(ix.Accounts) < 13 {
		return nil, errors.New("insufficient accounts in instruction")
	}

	srcMint := ix.Accounts[4].Pubkey   // index 4
	dstMint := ix.Accounts[5].Pubkey   // index 5
	receiver := ix.Accounts[3].Pubkey  // index 3
	protocolDstAta := ix.Accounts[11].Pubkey // index 11
	integratorDstAta := ix.Accounts[12].Pubkey // index 12

	return fromContractOrder(orderConfig, struct {
		SrcMint          *domains.Address
		DstMint          *domains.Address
		Receiver         *domains.Address
		ProtocolDstAta   *domains.Address
		IntegratorDstAta *domains.Address
		ProgramID        *domains.Address
	}{
		SrcMint:          srcMint,
		DstMint:          dstMint,
		Receiver:         receiver,
		ProtocolDstAta:   protocolDstAta,
		IntegratorDstAta: integratorDstAta,
		ProgramID:        ix.ProgramID,
	})
}

// FromContractOrder creates a FusionOrder from a ContractOrderConfig and account addresses
func FromContractOrder(
	reducedConfig *ContractOrderConfig,
	accounts struct {
		SrcMint          *domains.Address
		DstMint          *domains.Address
		Receiver         *domains.Address
		ProtocolDstAta   *domains.Address
		IntegratorDstAta *domains.Address
		ProgramID        *domains.Address
	},
) (*FusionOrder, error) {
	return fromContractOrder(reducedConfig, accounts)
}

// fromContractOrder is the internal implementation
func fromContractOrder(
	reducedConfig *ContractOrderConfig,
	accounts struct {
		SrcMint          *domains.Address
		DstMint          *domains.Address
		Receiver         *domains.Address
		ProtocolDstAta   *domains.Address
		IntegratorDstAta *domains.Address
		ProgramID        *domains.Address
	},
) (*FusionOrder, error) {
	auction := reducedConfig.DutchAuctionData
	fee := reducedConfig.Fee

	// Calculate order expiration delay
	orderExpirationDelay := reducedConfig.ExpirationTime - auction.Duration - auction.StartTime

	// Parse fees
	var protocolDstAta, integratorDstAta *domains.Address
	if accounts.ProtocolDstAta != nil && !accounts.ProtocolDstAta.Equal(accounts.ProgramID) {
		protocolDstAta = accounts.ProtocolDstAta
	}
	if accounts.IntegratorDstAta != nil && !accounts.IntegratorDstAta.Equal(accounts.ProgramID) {
		integratorDstAta = accounts.IntegratorDstAta
	}

	protocolFee := domains.BpsFromFraction(float64(fee.ProtocolFee)/100000.0, Base1E5)
	integratorFee := domains.BpsFromFraction(float64(fee.IntegratorFee)/100000.0, Base1E5)
	surplusShare := domains.BpsFromFraction(float64(fee.SurplusPercentage)/100.0, Base1E2)

	fees, err := NewFeeConfig(protocolDstAta, integratorDstAta, protocolFee, integratorFee, surplusShare)
	if err != nil {
		return nil, fmt.Errorf("failed to create fee config: %w", err)
	}

	// Parse resolver cancellation config
	resolverCancellationConfig, err := NewResolverCancellationConfig(
		fee.MaxCancellationPremium,
		reducedConfig.CancellationAuctionDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resolver cancellation config: %w", err)
	}

	// Parse auction details
	points := make([]AuctionPoint, len(auction.PointsAndTimeDeltas))
	for i, point := range auction.PointsAndTimeDeltas {
		points[i] = AuctionPoint{
			Coefficient: point.RateBump,
			Delay:       point.TimeDelta,
		}
	}

	auctionDetails, err := NewAuctionDetails(struct {
		StartTime       uint32
		InitialRateBump uint16
		Duration        uint32
		Points          []AuctionPoint
	}{
		StartTime:       auction.StartTime,
		Duration:        auction.Duration,
		InitialRateBump: auction.InitialRateBump,
		Points:          points,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create auction details: %w", err)
	}

	// Create order
	orderInfo := OrderInfoData{
		ID:                 reducedConfig.ID,
		SrcAmount:          reducedConfig.SrcAmount,
		MinDstAmount:       reducedConfig.MinDstAmount,
		EstimatedDstAmount: reducedConfig.EstimatedDstAmount,
		Receiver:           accounts.Receiver,
		SrcMint:            accounts.SrcMint,
		DstMint:            accounts.DstMint,
	}

	return NewFusionOrder(orderInfo, auctionDetails, struct {
		SrcAssetIsNative          bool
		DstAssetIsNative          bool
		OrderExpirationDelay      uint32
		Fees                      *FeeConfig
		ResolverCancellationConfig *ResolverCancellationConfig
	}{
		SrcAssetIsNative:          reducedConfig.SrcAssetIsNative,
		DstAssetIsNative:          reducedConfig.DstAssetIsNative,
		OrderExpirationDelay:      orderExpirationDelay,
		Fees:                      fees,
		ResolverCancellationConfig: resolverCancellationConfig,
	})
}

// deserializeOrderConfig deserializes a ContractOrderConfig from Borsh-encoded data
func deserializeOrderConfig(data []byte) (*ContractOrderConfig, error) {
	if len(data) < 4 {
		return nil, errors.New("data too short for order config")
	}

	offset := 0

	// id: u32
	if offset+4 > len(data) {
		return nil, errors.New("insufficient data for id")
	}
	id := binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	// srcAmount: u64
	if offset+8 > len(data) {
		return nil, errors.New("insufficient data for srcAmount")
	}
	srcAmount := new(big.Int).SetUint64(binary.LittleEndian.Uint64(data[offset:]))
	offset += 8

	// minDstAmount: u64
	if offset+8 > len(data) {
		return nil, errors.New("insufficient data for minDstAmount")
	}
	minDstAmount := new(big.Int).SetUint64(binary.LittleEndian.Uint64(data[offset:]))
	offset += 8

	// estimatedDstAmount: u64
	if offset+8 > len(data) {
		return nil, errors.New("insufficient data for estimatedDstAmount")
	}
	estimatedDstAmount := new(big.Int).SetUint64(binary.LittleEndian.Uint64(data[offset:]))
	offset += 8

	// expirationTime: u32
	if offset+4 > len(data) {
		return nil, errors.New("insufficient data for expirationTime")
	}
	expirationTime := binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	// srcAssetIsNative: bool
	if offset+1 > len(data) {
		return nil, errors.New("insufficient data for srcAssetIsNative")
	}
	srcAssetIsNative := data[offset] != 0
	offset += 1

	// dstAssetIsNative: bool
	if offset+1 > len(data) {
		return nil, errors.New("insufficient data for dstAssetIsNative")
	}
	dstAssetIsNative := data[offset] != 0
	offset += 1

	// fee: struct
	// protocolFee: u16
	if offset+2 > len(data) {
		return nil, errors.New("insufficient data for protocolFee")
	}
	protocolFee := binary.LittleEndian.Uint16(data[offset:])
	offset += 2

	// integratorFee: u16
	if offset+2 > len(data) {
		return nil, errors.New("insufficient data for integratorFee")
	}
	integratorFee := binary.LittleEndian.Uint16(data[offset:])
	offset += 2

	// surplusPercentage: u8
	if offset+1 > len(data) {
		return nil, errors.New("insufficient data for surplusPercentage")
	}
	surplusPercentage := data[offset]
	offset += 1

	// maxCancellationPremium: u64
	if offset+8 > len(data) {
		return nil, errors.New("insufficient data for maxCancellationPremium")
	}
	maxCancellationPremium := new(big.Int).SetUint64(binary.LittleEndian.Uint64(data[offset:]))
	offset += 8

	// dutchAuctionData: struct
	// startTime: u32
	if offset+4 > len(data) {
		return nil, errors.New("insufficient data for startTime")
	}
	startTime := binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	// duration: u32
	if offset+4 > len(data) {
		return nil, errors.New("insufficient data for duration")
	}
	duration := binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	// initialRateBump: u16
	if offset+2 > len(data) {
		return nil, errors.New("insufficient data for initialRateBump")
	}
	initialRateBump := binary.LittleEndian.Uint16(data[offset:])
	offset += 2

	// pointsAndTimeDeltas: array (u32 length prefix)
	if offset+4 > len(data) {
		return nil, errors.New("insufficient data for points count")
	}
	pointsCount := binary.LittleEndian.Uint32(data[offset:])
	offset += 4

	pointsAndTimeDeltas := make([]PointAndTimeDelta, pointsCount)
	for i := uint32(0); i < pointsCount; i++ {
		// rateBump: u16
		if offset+2 > len(data) {
			return nil, fmt.Errorf("insufficient data for point %d rateBump", i)
		}
		rateBump := binary.LittleEndian.Uint16(data[offset:])
		offset += 2

		// timeDelta: u16
		if offset+2 > len(data) {
			return nil, fmt.Errorf("insufficient data for point %d timeDelta", i)
		}
		timeDelta := binary.LittleEndian.Uint16(data[offset:])
		offset += 2

		pointsAndTimeDeltas[i] = PointAndTimeDelta{
			RateBump:  rateBump,
			TimeDelta: timeDelta,
		}
	}

	// cancellationAuctionDuration: u32
	if offset+4 > len(data) {
		return nil, errors.New("insufficient data for cancellationAuctionDuration")
	}
	cancellationAuctionDuration := binary.LittleEndian.Uint32(data[offset:])

	return &ContractOrderConfig{
		ID:                          id,
		SrcAmount:                   srcAmount,
		MinDstAmount:                minDstAmount,
		EstimatedDstAmount:          estimatedDstAmount,
		ExpirationTime:              expirationTime,
		SrcAssetIsNative:            srcAssetIsNative,
		DstAssetIsNative:            dstAssetIsNative,
		CancellationAuctionDuration: cancellationAuctionDuration,
		Fee: FeeConfigStruct{
			ProtocolFee:            protocolFee,
			IntegratorFee:          integratorFee,
			SurplusPercentage:      surplusPercentage,
			MaxCancellationPremium: maxCancellationPremium,
		},
		DutchAuctionData: DutchAuctionDataStruct{
			StartTime:           startTime,
			Duration:            duration,
			InitialRateBump:     initialRateBump,
			PointsAndTimeDeltas: pointsAndTimeDeltas,
		},
	}, nil
}

// equalBytes checks if two byte slices are equal
func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
