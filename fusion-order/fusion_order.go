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
	"github.com/dawitel/solana-fusion-sdk-go/utils/math"
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

	// Validate values
	if err := validation.AssertUInteger(orderExpirationDelay, nil); err != nil {
		return nil, err
	}
	if err := validation.AssertUInteger(deadline, nil); err != nil {
		return nil, err
	}
	if err := validation.AssertUInteger(orderInfo.ID, nil); err != nil {
		return nil, err
	}

	// Validate amounts (u64 max)
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

	// Handle native tokens - use wrapped native address
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

// GetEscrow returns the escrow ATA for src token
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

	// Get PDA for escrow
	orderHash := f.GetOrderHash()
	seeds := [][]byte{
		[]byte("escrow"),
		maker.ToBuffer(),
		orderHash,
	}
	escrow, err := addresses.GetPda(programId, seeds)
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow PDA: %w", err)
	}

	// Get ATA for escrow
	ata, err := addresses.GetAta(escrow, f.orderConfig.SrcMint, srcTokenProgram)
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow ATA: %w", err)
	}

	return ata, nil
}

// CalcTakingAmount calculates required taking amount to fill order for passed makingAmount at block time
// Note: This is a simplified version that doesn't use AmountCalculator to avoid import cycle
// For full calculation with auction and fees, use the amount-calculator package directly
func (f *FusionOrder) CalcTakingAmount(makingAmount *big.Int, time uint32) *big.Int {
	takingAmount := math.CalcTakingAmount(
		makingAmount,
		f.orderConfig.SrcAmount,
		f.orderConfig.MinDstAmount,
	)

	// Note: For full calculation with auction rate bump, use amount-calculator package:
	//   calculator := amountcalculator.NewAmountCalculator(...)
	//   return calculator.GetRequiredTakingAmount(takingAmount, time)
	return takingAmount
}

// GetUserReceiveAmount returns how much user will receive in dst token
// Note: This method requires the amount-calculator package. Use it externally to avoid import cycles.
func (f *FusionOrder) GetUserReceiveAmount(makingAmount *big.Int, time uint32) *big.Int {
	// This method is deprecated - use amount-calculator package directly
	// takingAmount := math.CalcTakingAmount(makingAmount, f.orderConfig.SrcAmount, f.orderConfig.MinDstAmount)
	// Use amount-calculator package to get full calculation with fees
	return big.NewInt(0)
}

// GetIntegratorFee returns fee in dstToken which integrator gets
// Note: This method requires the amount-calculator package. Use it externally to avoid import cycles.
func (f *FusionOrder) GetIntegratorFee(time uint32, makingAmount *big.Int) *big.Int {
	// This method is deprecated - use amount-calculator package directly
	return big.NewInt(0)
}

// GetProtocolFee returns fee in dstToken which protocol gets
// Note: This method requires the amount-calculator package. Use it externally to avoid import cycles.
func (f *FusionOrder) GetProtocolFee(time uint32, makingAmount *big.Int) *big.Int {
	// This method is deprecated - use amount-calculator package directly
	return big.NewInt(0)
}

// IsExpiredAt checks if order expired at a given time
func (f *FusionOrder) IsExpiredAt(time uint32) bool {
	return time > f.orderConfig.ExpirationTime
}

// GetCalculatorInfo returns information needed to create an AmountCalculator for this order
// To create the calculator, use the amount-calculator package:
//
//	import amountcalculator "github.com/dawitel/solana-fusion-sdk-go/amount-calculator"
//	auctionCalc := amountcalculator.FromAuctionData(order.AuctionDetails())
//	feeCalc := amountcalculator.FromFeeConfig(order.Fees())
//	calculator := amountcalculator.NewAmountCalculator(auctionCalc, feeCalc)
//
// Then use the calculator with helper methods like CalcTakingAmountWithCalculator
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

// SerializeBorsh serializes the order config to Borsh format
func (c *ContractOrderConfig) SerializeBorsh() ([]byte, error) {
	var data []byte

	// id: u32
	idBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(idBytes, c.ID)
	data = append(data, idBytes...)

	// srcAmount: u64
	srcAmountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(srcAmountBytes, c.SrcAmount.Uint64())
	data = append(data, srcAmountBytes...)

	// minDstAmount: u64
	minDstAmountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(minDstAmountBytes, c.MinDstAmount.Uint64())
	data = append(data, minDstAmountBytes...)

	// estimatedDstAmount: u64
	estimatedDstAmountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(estimatedDstAmountBytes, c.EstimatedDstAmount.Uint64())
	data = append(data, estimatedDstAmountBytes...)

	// expirationTime: u32
	expirationTimeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(expirationTimeBytes, c.ExpirationTime)
	data = append(data, expirationTimeBytes...)

	// srcAssetIsNative: bool
	if c.SrcAssetIsNative {
		data = append(data, 1)
	} else {
		data = append(data, 0)
	}

	// dstAssetIsNative: bool
	if c.DstAssetIsNative {
		data = append(data, 1)
	} else {
		data = append(data, 0)
	}

	// fee: struct
	// protocolFee: u16
	protocolFeeBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(protocolFeeBytes, c.Fee.ProtocolFee)
	data = append(data, protocolFeeBytes...)

	// integratorFee: u16
	integratorFeeBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(integratorFeeBytes, c.Fee.IntegratorFee)
	data = append(data, integratorFeeBytes...)

	// surplusPercentage: u8
	data = append(data, byte(c.Fee.SurplusPercentage))

	// maxCancellationPremium: u64
	maxCancellationPremiumBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(maxCancellationPremiumBytes, c.Fee.MaxCancellationPremium.Uint64())
	data = append(data, maxCancellationPremiumBytes...)

	// dutchAuctionData: struct
	// startTime: u32
	startTimeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(startTimeBytes, c.DutchAuctionData.StartTime)
	data = append(data, startTimeBytes...)

	// duration: u32
	durationBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(durationBytes, c.DutchAuctionData.Duration)
	data = append(data, durationBytes...)

	// initialRateBump: u16
	initialRateBumpBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(initialRateBumpBytes, c.DutchAuctionData.InitialRateBump)
	data = append(data, initialRateBumpBytes...)

	// pointsAndTimeDeltas: array
	pointsCount := len(c.DutchAuctionData.PointsAndTimeDeltas)
	pointsCountBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(pointsCountBytes, uint32(pointsCount))
	data = append(data, pointsCountBytes...)

	for _, point := range c.DutchAuctionData.PointsAndTimeDeltas {
		// rateBump: u16
		rateBumpBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(rateBumpBytes, point.RateBump)
		data = append(data, rateBumpBytes...)

		// timeDelta: u16
		timeDeltaBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(timeDeltaBytes, point.TimeDelta)
		data = append(data, timeDeltaBytes...)
	}

	// cancellationAuctionDuration: u32
	cancellationAuctionDurationBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(cancellationAuctionDurationBytes, c.CancellationAuctionDuration)
	data = append(data, cancellationAuctionDurationBytes...)

	return data, nil
}

// SerializeOptionalAddress serializes an optional address for Borsh
func SerializeOptionalAddress(addr *domains.Address) []byte {
	if addr == nil {
		return []byte{0} // None
	}
	data := []byte{1} // Some
	data = append(data, addr.ToBuffer()...)
	return data
}

// SerializeAddress serializes an address for Borsh
func SerializeAddress(addr *domains.Address) []byte {
	return addr.ToBuffer()
}

// GetOrderHash returns the order hash (SHA256 of Borsh-serialized order)
// Note: This method panics if serialization fails, which should never happen
// for a properly constructed order. If you need error handling, use GetOrderHashWithError.
func (f *FusionOrder) GetOrderHash() []byte {
	orderConfig := f.Build()
	borshData, err := orderConfig.SerializeBorsh()
	if err != nil {
		panic(fmt.Sprintf("failed to serialize order config for hash: %v", err))
	}

	// Hash the order config
	hash := sha256.Sum256(borshData)

	// Append optional addresses
	if f.orderConfig.Fees != nil {
		hash = sha256.Sum256(append(hash[:], SerializeOptionalAddress(f.orderConfig.Fees.ProtocolDstAta)...))
		hash = sha256.Sum256(append(hash[:], SerializeOptionalAddress(f.orderConfig.Fees.IntegratorDstAta)...))
	}

	// Append addresses
	hash = sha256.Sum256(append(hash[:], SerializeAddress(f.orderConfig.SrcMint)...))
	hash = sha256.Sum256(append(hash[:], SerializeAddress(f.orderConfig.DstMint)...))
	hash = sha256.Sum256(append(hash[:], SerializeAddress(f.orderConfig.Receiver)...))

	return hash[:]
}

// GetOrderHashBase58 returns the base58 encoded order hash
func (f *FusionOrder) GetOrderHashBase58() string {
	return base58.Encode(f.GetOrderHash())
}
