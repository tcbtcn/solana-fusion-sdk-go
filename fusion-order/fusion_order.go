package fusionorder

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
	"github.com/dawitel/solana-fusion-sdk-go/idl"
	"github.com/dawitel/solana-fusion-sdk-go/utils/addresses"
	"github.com/dawitel/solana-fusion-sdk-go/utils/validation"
	"github.com/mr-tron/base58"
)

const (
	// DefaultOrderExpirationDelay is the default order expiration delay in seconds
	DefaultOrderExpirationDelay = 12
)

// FusionOrder represents a fusion order
type FusionOrder struct {
	orderConfig struct {
		ID                         uint32
		SrcAmount                  *big.Int
		MinDstAmount               *big.Int
		EstimatedDstAmount         *big.Int
		ExpirationTime             uint32
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		Receiver                   *domains.Address
		Fees                       *FeeConfig
		ResolverCancellationConfig *ResolverCancellationConfig
		DutchAuctionData           *AuctionDetails
		SrcMint                    *domains.Address
		DstMint                    *domains.Address
	}
}

// NewFusionOrder creates a new FusionOrder
func NewFusionOrder(
	orderInfo OrderInfoData,
	auctionDetails *AuctionDetails,
	extra struct {
		SrcAssetIsNative           bool
		DstAssetIsNative           bool
		OrderExpirationDelay       uint32
		Fees                       *FeeConfig
		ResolverCancellationConfig *ResolverCancellationConfig
	},
) (*FusionOrder, error) {
	// Validate tokens are different
	if orderInfo.SrcMint.Equal(orderInfo.DstMint) {
		return nil, errors.New("tokens must be different")
	}

	orderExpirationDelay := extra.OrderExpirationDelay
	if orderExpirationDelay == 0 {
		orderExpirationDelay = DefaultOrderExpirationDelay
	}

	deadline := auctionDetails.StartTime + auctionDetails.Duration + orderExpirationDelay

	if err := validation.AssertUInteger(orderExpirationDelay, nil); err != nil {
		return nil, err
	}
	if err := validation.AssertUInteger(deadline, nil); err != nil {
		return nil, err
	}
	if err := validation.AssertUInteger(orderInfo.ID, nil); err != nil {
		return nil, err
	}

	u64Max := big.NewInt(0).SetUint64(^uint64(0))
	if orderInfo.SrcAmount.Cmp(u64Max) > 0 {
		return nil, errors.New("srcAmount exceeds u64 max")
	}
	if orderInfo.EstimatedDstAmount.Cmp(u64Max) > 0 {
		return nil, errors.New("estimatedDstAmount exceeds u64 max")
	}
	if orderInfo.MinDstAmount.Cmp(u64Max) > 0 {
		return nil, errors.New("minDstAmount exceeds u64 max")
	}

	fees := extra.Fees
	if fees == nil {
		fees = ZeroFeeConfig
	}

	resolverCancellationConfig := extra.ResolverCancellationConfig
	if resolverCancellationConfig == nil {
		resolverCancellationConfig = AlmostZeroResolverCancellationConfig
	}

	srcMint := orderInfo.SrcMint
	if srcMint.IsNative() {
		srcMint = domains.WRAPPED_NATIVE
	}
	dstMint := orderInfo.DstMint
	if dstMint.IsNative() {
		dstMint = domains.WRAPPED_NATIVE
	}

	order := &FusionOrder{}
	order.orderConfig.ID = orderInfo.ID
	order.orderConfig.SrcAmount = orderInfo.SrcAmount
	order.orderConfig.MinDstAmount = orderInfo.MinDstAmount
	order.orderConfig.EstimatedDstAmount = orderInfo.EstimatedDstAmount
	order.orderConfig.ExpirationTime = deadline
	order.orderConfig.SrcAssetIsNative = extra.SrcAssetIsNative
	order.orderConfig.DstAssetIsNative = extra.DstAssetIsNative
	order.orderConfig.Receiver = orderInfo.Receiver
	order.orderConfig.Fees = fees
	order.orderConfig.ResolverCancellationConfig = resolverCancellationConfig
	order.orderConfig.DutchAuctionData = auctionDetails
	order.orderConfig.SrcMint = srcMint
	order.orderConfig.DstMint = dstMint

	return order, nil
}

// Fees returns the fee config, or nil if zero
func (f *FusionOrder) Fees() *FeeConfig {
	if f.orderConfig.Fees == nil || f.orderConfig.Fees.IsZero() {
		return nil
	}
	return f.orderConfig.Fees
}

// ResolverCancellationConfig returns the resolver cancellation config, or nil if zero
func (f *FusionOrder) ResolverCancellationConfig() *ResolverCancellationConfig {
	if f.orderConfig.ResolverCancellationConfig == nil || f.orderConfig.ResolverCancellationConfig.IsZero() {
		return nil
	}
	return f.orderConfig.ResolverCancellationConfig
}

// SrcMint returns the source mint address
func (f *FusionOrder) SrcMint() *domains.Address {
	return f.orderConfig.SrcMint
}

// DstMint returns the destination mint address
func (f *FusionOrder) DstMint() *domains.Address {
	return f.orderConfig.DstMint
}

// SrcAmount returns the source amount
func (f *FusionOrder) SrcAmount() *big.Int {
	return f.orderConfig.SrcAmount
}

// MinDstAmount returns the minimum destination amount
func (f *FusionOrder) MinDstAmount() *big.Int {
	return f.orderConfig.MinDstAmount
}

// EstimatedDstAmount returns the estimated destination amount
func (f *FusionOrder) EstimatedDstAmount() *big.Int {
	return f.orderConfig.EstimatedDstAmount
}

// Receiver returns the receiver address
func (f *FusionOrder) Receiver() *domains.Address {
	return f.orderConfig.Receiver
}

// Deadline returns the expiration time in seconds
func (f *FusionOrder) Deadline() uint32 {
	return f.orderConfig.ExpirationTime
}

// AuctionStartTime returns the auction start time in seconds
func (f *FusionOrder) AuctionStartTime() uint32 {
	return f.orderConfig.DutchAuctionData.StartTime
}

// AuctionEndTime returns the auction end time in seconds
func (f *FusionOrder) AuctionEndTime() uint32 {
	return f.orderConfig.DutchAuctionData.StartTime + f.orderConfig.DutchAuctionData.Duration
}

// AuctionDetails returns the auction details
func (f *FusionOrder) AuctionDetails() *AuctionDetails {
	return f.orderConfig.DutchAuctionData
}

// ID returns the order ID
func (f *FusionOrder) ID() uint32 {
	return f.orderConfig.ID
}

// SrcAssetIsNative returns whether the source asset is native
func (f *FusionOrder) SrcAssetIsNative() bool {
	return f.orderConfig.SrcAssetIsNative
}

// DstAssetIsNative returns whether the destination asset is native
func (f *FusionOrder) DstAssetIsNative() bool {
	return f.orderConfig.DstAssetIsNative
}

// GetEscrow returns the escrow ATA for src token.
func (f *FusionOrder) GetEscrow(
	maker domains.AddressLike,
	srcTokenProgram domains.AddressLike,
	programId domains.AddressLike,
) (*domains.Address, error) {
	if srcTokenProgram == nil {
		srcTokenProgram = domains.TOKEN_PROGRAM_ID
	}
	if programId == nil {
		programId = domains.MustAddressFromString(idl.FusionSwapProgramAddress)
	}

	orderHash, err := f.GetOrderHashWithError()
	if err != nil {
		return nil, fmt.Errorf("failed to get order hash: %w", err)
	}

	seeds := [][]byte{
		[]byte("escrow"),
		maker.ToBuffer(),
		orderHash,
	}
	escrow, err := addresses.GetPda(programId, seeds)
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow PDA: %w", err)
	}

	ata, err := addresses.GetAta(escrow, f.orderConfig.SrcMint, srcTokenProgram)
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow ATA: %w", err)
	}

	return ata, nil
}

// IsExpiredAt checks if order expired at a given time.
func (f *FusionOrder) IsExpiredAt(time uint32) bool {
	return time > f.orderConfig.ExpirationTime
}

// GetCalculatorInfo returns information needed to create an AmountCalculator for this order.
func (f *FusionOrder) GetCalculatorInfo() struct {
	AuctionDetails *AuctionDetails
	Fees           *FeeConfig
} {
	return struct {
		AuctionDetails *AuctionDetails
		Fees           *FeeConfig
	}{
		AuctionDetails: f.orderConfig.DutchAuctionData,
		Fees:           f.orderConfig.Fees,
	}
}

// Build returns the ContractOrderConfig for Borsh serialization
func (f *FusionOrder) Build() *ContractOrderConfig {
	auction := f.orderConfig.DutchAuctionData
	fees := f.orderConfig.Fees

	pointsAndTimeDeltas := make([]PointAndTimeDelta, len(auction.Points))
	for i, point := range auction.Points {
		pointsAndTimeDeltas[i] = PointAndTimeDelta{
			RateBump:  point.Coefficient,
			TimeDelta: point.Delay,
		}
	}

	var protocolFee, integratorFee uint16
	var surplusPercentage uint8
	var maxCancellationPremium *big.Int

	if fees != nil {
		protocolFee = uint16(fees.ProtocolFee.ToFraction(Base1E5) * 100000)
		integratorFee = uint16(fees.IntegratorFee.ToFraction(Base1E5) * 100000)
		surplusPercentage = uint8(fees.SurplusShare.ToFraction(Base1E2) * 100)
	}

	if f.orderConfig.ResolverCancellationConfig != nil {
		maxCancellationPremium = f.orderConfig.ResolverCancellationConfig.MaxCancellationPremium
	} else {
		maxCancellationPremium = big.NewInt(0)
	}

	var cancellationAuctionDuration uint32
	if f.orderConfig.ResolverCancellationConfig != nil {
		cancellationAuctionDuration = f.orderConfig.ResolverCancellationConfig.CancellationAuctionDuration
	}

	return &ContractOrderConfig{
		ID:                          f.orderConfig.ID,
		SrcAmount:                   f.orderConfig.SrcAmount,
		MinDstAmount:                f.orderConfig.MinDstAmount,
		EstimatedDstAmount:          f.orderConfig.EstimatedDstAmount,
		ExpirationTime:              f.orderConfig.ExpirationTime,
		SrcAssetIsNative:            f.orderConfig.SrcAssetIsNative,
		DstAssetIsNative:            f.orderConfig.DstAssetIsNative,
		CancellationAuctionDuration: cancellationAuctionDuration,
		Fee: FeeConfigStruct{
			ProtocolFee:            protocolFee,
			IntegratorFee:          integratorFee,
			SurplusPercentage:      surplusPercentage,
			MaxCancellationPremium: maxCancellationPremium,
		},
		DutchAuctionData: DutchAuctionDataStruct{
			StartTime:           auction.StartTime,
			Duration:            auction.Duration,
			InitialRateBump:     auction.InitialRateBump,
			PointsAndTimeDeltas: pointsAndTimeDeltas,
		},
	}
}

// ContractOrderConfig represents the order config for Borsh serialization
type ContractOrderConfig struct {
	ID                          uint32
	SrcAmount                   *big.Int
	MinDstAmount                *big.Int
	EstimatedDstAmount          *big.Int
	ExpirationTime              uint32
	SrcAssetIsNative            bool
	DstAssetIsNative            bool
	CancellationAuctionDuration uint32
	Fee                         FeeConfigStruct
	DutchAuctionData            DutchAuctionDataStruct
}

type FeeConfigStruct struct {
	ProtocolFee            uint16
	IntegratorFee          uint16
	SurplusPercentage      uint8
	MaxCancellationPremium *big.Int
}

type DutchAuctionDataStruct struct {
	StartTime           uint32
	Duration            uint32
	InitialRateBump     uint16
	PointsAndTimeDeltas []PointAndTimeDelta
}

type PointAndTimeDelta struct {
	RateBump  uint16
	TimeDelta uint16
}

// SerializeBorsh serializes the order config to Borsh format.
func (c *ContractOrderConfig) SerializeBorsh() ([]byte, error) {
	var data []byte

	idBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(idBytes, c.ID)
	data = append(data, idBytes...)

	srcAmountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(srcAmountBytes, c.SrcAmount.Uint64())
	data = append(data, srcAmountBytes...)

	minDstAmountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(minDstAmountBytes, c.MinDstAmount.Uint64())
	data = append(data, minDstAmountBytes...)

	estimatedDstAmountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(estimatedDstAmountBytes, c.EstimatedDstAmount.Uint64())
	data = append(data, estimatedDstAmountBytes...)

	expirationTimeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(expirationTimeBytes, c.ExpirationTime)
	data = append(data, expirationTimeBytes...)

	if c.SrcAssetIsNative {
		data = append(data, 1)
	} else {
		data = append(data, 0)
	}

	if c.DstAssetIsNative {
		data = append(data, 1)
	} else {
		data = append(data, 0)
	}

	protocolFeeBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(protocolFeeBytes, c.Fee.ProtocolFee)
	data = append(data, protocolFeeBytes...)

	integratorFeeBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(integratorFeeBytes, c.Fee.IntegratorFee)
	data = append(data, integratorFeeBytes...)

	data = append(data, byte(c.Fee.SurplusPercentage))

	maxCancellationPremiumBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(maxCancellationPremiumBytes, c.Fee.MaxCancellationPremium.Uint64())
	data = append(data, maxCancellationPremiumBytes...)

	startTimeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(startTimeBytes, c.DutchAuctionData.StartTime)
	data = append(data, startTimeBytes...)

	durationBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(durationBytes, c.DutchAuctionData.Duration)
	data = append(data, durationBytes...)

	initialRateBumpBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(initialRateBumpBytes, c.DutchAuctionData.InitialRateBump)
	data = append(data, initialRateBumpBytes...)

	pointsCount := len(c.DutchAuctionData.PointsAndTimeDeltas)
	pointsCountBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(pointsCountBytes, uint32(pointsCount))
	data = append(data, pointsCountBytes...)

	for _, point := range c.DutchAuctionData.PointsAndTimeDeltas {
		rateBumpBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(rateBumpBytes, point.RateBump)
		data = append(data, rateBumpBytes...)

		timeDeltaBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(timeDeltaBytes, point.TimeDelta)
		data = append(data, timeDeltaBytes...)
	}

	cancellationAuctionDurationBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(cancellationAuctionDurationBytes, c.CancellationAuctionDuration)
	data = append(data, cancellationAuctionDurationBytes...)

	return data, nil
}

// SerializeOptionalAddress serializes an optional address for Borsh.
func SerializeOptionalAddress(addr *domains.Address) []byte {
	if addr == nil {
		return []byte{0}
	}
	data := []byte{1}
	data = append(data, addr.ToBuffer()...)
	return data
}

// SerializeAddress serializes an address for Borsh.
func SerializeAddress(addr *domains.Address) []byte {
	return addr.ToBuffer()
}

// GetOrderHash returns the order hash (SHA256 of Borsh-serialized order).
func (f *FusionOrder) GetOrderHash() []byte {
	hash, err := f.GetOrderHashWithError()
	if err != nil {
		panic(fmt.Sprintf("failed to get order hash: %v", err))
	}
	return hash
}

// GetOrderHashWithError returns the order hash (SHA256 of Borsh-serialized order) with error handling.
func (f *FusionOrder) GetOrderHashWithError() ([]byte, error) {
	orderConfig := f.Build()
	borshData, err := orderConfig.SerializeBorsh()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize order config for hash: %w", err)
	}

	hash := sha256.Sum256(borshData)

	var protocolDstAta, integratorDstAta *domains.Address
	if f.orderConfig.Fees != nil {
		protocolDstAta = f.orderConfig.Fees.ProtocolDstAta
		integratorDstAta = f.orderConfig.Fees.IntegratorDstAta
	}
	hash = sha256.Sum256(append(hash[:], SerializeOptionalAddress(protocolDstAta)...))
	hash = sha256.Sum256(append(hash[:], SerializeOptionalAddress(integratorDstAta)...))

	hash = sha256.Sum256(append(hash[:], SerializeAddress(f.orderConfig.SrcMint)...))
	hash = sha256.Sum256(append(hash[:], SerializeAddress(f.orderConfig.DstMint)...))
	hash = sha256.Sum256(append(hash[:], SerializeAddress(f.orderConfig.Receiver)...))

	return hash[:], nil
}

// GetOrderHashBase58 returns the base58 encoded order hash.
func (f *FusionOrder) GetOrderHashBase58() string {
	return base58.Encode(f.GetOrderHash())
}
