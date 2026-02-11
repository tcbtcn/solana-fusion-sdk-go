package contracts

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
	fusionorder "github.com/dawitel/solana-fusion-sdk-go/fusion-order"
	"github.com/dawitel/solana-fusion-sdk-go/idl"
	"github.com/dawitel/solana-fusion-sdk-go/types"
	"github.com/dawitel/solana-fusion-sdk-go/utils/addresses"
)

// FusionSwapContractAddress is the default FusionSwap contract address
var FusionSwapContractAddress = domains.MustAddressFromString(idl.FusionSwapProgramAddress)

// FusionSwapContract represents the FusionSwap contract
type FusionSwapContract struct {
	ProgramID *domains.Address
}

// DefaultFusionSwapContract returns the default FusionSwap contract
func DefaultFusionSwapContract() *FusionSwapContract {
	return &FusionSwapContract{
		ProgramID: FusionSwapContractAddress,
	}
}

// NewFusionSwapContract creates a new FusionSwapContract with a custom program ID
func NewFusionSwapContract(programID *domains.Address) *FusionSwapContract {
	return &FusionSwapContract{
		ProgramID: programID,
	}
}

// optionalAccountMeta returns an optional account meta (uses program ID as placeholder)
func (f *FusionSwapContract) optionalAccountMeta() types.AccountMeta {
	return types.AccountMeta{
		Pubkey:     f.ProgramID,
		IsSigner:   false,
		IsWritable: false,
	}
}

// optional returns the address or program ID if nil
func (f *FusionSwapContract) optional(addr *domains.Address) *domains.Address {
	if addr == nil {
		return f.ProgramID
	}
	return addr
}

// Create creates a create order instruction
func (f *FusionSwapContract) Create(
	order *fusionorder.FusionOrder,
	accounts struct {
		Maker           *domains.Address
		SrcTokenProgram *domains.Address
	},
) (*types.TransactionInstruction, error) {
	escrow, err := addresses.GetPda(f.ProgramID, [][]byte{
		[]byte("escrow"),
		accounts.Maker.ToBuffer(),
		order.GetOrderHash(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow PDA: %w", err)
	}

	escrowSrcAta, err := addresses.GetAta(escrow, order.SrcMint(), accounts.SrcTokenProgram)
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow ATA: %w", err)
	}

	accountMetas := []types.AccountMeta{
		// 0. system_program
		types.AccountMeta{Pubkey: domains.SYSTEM_PROGRAM_ID, IsSigner: false, IsWritable: false},
		// 1. escrow
		types.AccountMeta{Pubkey: escrow, IsSigner: false, IsWritable: false},
		// 2. src_mint
		types.AccountMeta{Pubkey: order.SrcMint(), IsSigner: false, IsWritable: false},
		// 3. src_token_program
		types.AccountMeta{Pubkey: accounts.SrcTokenProgram, IsSigner: false, IsWritable: false},
		// 4. escrow_src_ata
		types.AccountMeta{Pubkey: escrowSrcAta, IsSigner: false, IsWritable: true},
		// 5. maker
		types.AccountMeta{Pubkey: accounts.Maker, IsSigner: true, IsWritable: true},
	}

	// 6. maker_src_ata (optional if native)
	if order.SrcAssetIsNative() {
		accountMetas = append(accountMetas, f.optionalAccountMeta())
	} else {
		makerSrcAta, err := addresses.GetAta(accounts.Maker, order.SrcMint(), accounts.SrcTokenProgram)
		if err != nil {
			return nil, fmt.Errorf("failed to get maker ATA: %w", err)
		}
		accountMetas = append(accountMetas, types.AccountMeta{
			Pubkey:     makerSrcAta,
			IsSigner:   false,
			IsWritable: true,
		})
	}

	accountMetas = append(accountMetas,
		// 7. dst_mint
		types.AccountMeta{Pubkey: order.DstMint(), IsSigner: false, IsWritable: false},
		// 8. maker_receiver
		types.AccountMeta{Pubkey: order.Receiver(), IsSigner: false, IsWritable: true},
		// 9. associated_token_program
		types.AccountMeta{Pubkey: domains.ASSOCIATED_TOKEN_PROGRAM_ID, IsSigner: false, IsWritable: false},
	)

	// 10. protocol_dst_ata (optional)
	fees := order.Fees()
	var protocolDstAta, integratorDstAta *domains.Address
	if fees != nil {
		protocolDstAta = fees.ProtocolDstAta
		integratorDstAta = fees.IntegratorDstAta
	}
	accountMetas = append(accountMetas,
		types.AccountMeta{Pubkey: f.optional(protocolDstAta), IsSigner: false, IsWritable: false},
		// 11. integrator_dst_ata
		types.AccountMeta{Pubkey: f.optional(integratorDstAta), IsSigner: false, IsWritable: false},
	)

	// Encode instruction data (Borsh)
	data, err := f.encodeCreateInstruction(order)
	if err != nil {
		return nil, fmt.Errorf("failed to encode create instruction: %w", err)
	}

	return types.NewTransactionInstruction(f.ProgramID, accountMetas, data), nil
}

// encodeCreateInstruction encodes the create instruction data
func (f *FusionSwapContract) encodeCreateInstruction(order *fusionorder.FusionOrder) ([]byte, error) {
	orderConfig := order.Build()
	borshData, err := orderConfig.SerializeBorsh()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize order config: %w", err)
	}

	// Append order config data
	return append(idl.CreateOrderDiscriminator, borshData...), nil
}

// Fill creates a fill order instruction
func (f *FusionSwapContract) Fill(
	order *fusionorder.FusionOrder,
	amount *big.Int,
	accounts struct {
		Taker           *domains.Address
		Maker           *domains.Address
		SrcTokenProgram *domains.Address
		DstTokenProgram *domains.Address
		TakerSrcAccount *domains.Address
		Whitelist       *domains.Address
	},
) (*types.TransactionInstruction, error) {
	whitelist := accounts.Whitelist
	if whitelist == nil {
		whitelist = WhitelistContractAddress
	}

	escrow, err := addresses.GetPda(f.ProgramID, [][]byte{
		[]byte("escrow"),
		accounts.Maker.ToBuffer(),
		order.GetOrderHash(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow PDA: %w", err)
	}

	resolverAccess, err := addresses.GetPda(whitelist, [][]byte{
		[]byte("resolver_access"),
		accounts.Taker.ToBuffer(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get resolver access PDA: %w", err)
	}

	escrowSrcAta, err := addresses.GetAta(escrow, order.SrcMint(), accounts.SrcTokenProgram)
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow ATA: %w", err)
	}

	takerSrcAta := accounts.TakerSrcAccount
	if takerSrcAta == nil {
		takerSrcAta, err = addresses.GetAta(accounts.Taker, order.SrcMint(), accounts.SrcTokenProgram)
		if err != nil {
			return nil, fmt.Errorf("failed to get taker ATA: %w", err)
		}
	}

	accountMetas := []types.AccountMeta{
		// 0. taker
		types.AccountMeta{Pubkey: accounts.Taker, IsSigner: true, IsWritable: true},
		// 1. resolver_access
		types.AccountMeta{Pubkey: resolverAccess, IsSigner: false, IsWritable: false},
		// 2. maker
		types.AccountMeta{Pubkey: accounts.Maker, IsSigner: false, IsWritable: true},
		// 3. maker_receiver
		types.AccountMeta{Pubkey: order.Receiver(), IsSigner: false, IsWritable: order.DstAssetIsNative() || !accounts.Maker.Equal(order.Receiver())},
		// 4. src_mint
		types.AccountMeta{Pubkey: order.SrcMint(), IsSigner: false, IsWritable: false},
		// 5. dst_mint
		types.AccountMeta{Pubkey: order.DstMint(), IsSigner: false, IsWritable: false},
		// 6. escrow
		types.AccountMeta{Pubkey: escrow, IsSigner: false, IsWritable: false},
		// 7. escrow_src_ata
		types.AccountMeta{Pubkey: escrowSrcAta, IsSigner: false, IsWritable: true},
		// 8. taker_src_ata
		types.AccountMeta{Pubkey: takerSrcAta, IsSigner: false, IsWritable: true},
		// 9. src_token_program
		types.AccountMeta{Pubkey: accounts.SrcTokenProgram, IsSigner: false, IsWritable: false},
		// 10. dst_token_program
		types.AccountMeta{Pubkey: accounts.DstTokenProgram, IsSigner: false, IsWritable: false},
		// 11. system_program
		types.AccountMeta{Pubkey: domains.SYSTEM_PROGRAM_ID, IsSigner: false, IsWritable: false},
		// 12. associated_token_program
		types.AccountMeta{Pubkey: domains.ASSOCIATED_TOKEN_PROGRAM_ID, IsSigner: false, IsWritable: false},
	}

	// 13. maker_dst_ata (optional if native)
	if order.DstAssetIsNative() {
		accountMetas = append(accountMetas, f.optionalAccountMeta())
	} else {
		makerDstAta, err := addresses.GetAta(order.Receiver(), order.DstMint(), accounts.DstTokenProgram)
		if err != nil {
			return nil, fmt.Errorf("failed to get maker dst ATA: %w", err)
		}
		accountMetas = append(accountMetas, types.AccountMeta{
			Pubkey:     makerDstAta,
			IsSigner:   false,
			IsWritable: true,
		})
	}

	// 14. taker_dst_ata (optional if native)
	if order.DstAssetIsNative() {
		accountMetas = append(accountMetas, f.optionalAccountMeta())
	} else {
		takerDstAta, err := addresses.GetAta(accounts.Taker, order.DstMint(), accounts.DstTokenProgram)
		if err != nil {
			return nil, fmt.Errorf("failed to get taker dst ATA: %w", err)
		}
		accountMetas = append(accountMetas, types.AccountMeta{
			Pubkey:     takerDstAta,
			IsSigner:   false,
			IsWritable: true,
		})
	}

	// 15. protocol_dst_ata (optional)
	fees := order.Fees()
	if order.DstAssetIsNative() || fees == nil || fees.ProtocolDstAta == nil {
		accountMetas = append(accountMetas, f.optionalAccountMeta())
	} else {
		accountMetas = append(accountMetas, types.AccountMeta{
			Pubkey:     fees.ProtocolDstAta,
			IsSigner:   false,
			IsWritable: true,
		})
	}

	// 16. integrator_dst_ata (optional)
	if order.DstAssetIsNative() || fees == nil || fees.IntegratorDstAta == nil {
		accountMetas = append(accountMetas, f.optionalAccountMeta())
	} else {
		accountMetas = append(accountMetas, types.AccountMeta{
			Pubkey:     fees.IntegratorDstAta,
			IsSigner:   false,
			IsWritable: true,
		})
	}

	// Encode instruction data
	data, err := f.encodeFillInstruction(order, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode fill instruction: %w", err)
	}

	return types.NewTransactionInstruction(f.ProgramID, accountMetas, data), nil
}

// encodeFillInstruction encodes the fill instruction data
func (f *FusionSwapContract) encodeFillInstruction(order *fusionorder.FusionOrder, amount *big.Int) ([]byte, error) {
	orderConfig := order.Build()
	borshData, err := orderConfig.SerializeBorsh()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize order config: %w", err)
	}

	// Append order config
	result := append(idl.FillOrderDiscriminator, borshData...)

	// Append amount (u64, little-endian)
	amountBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(amountBytes, amount.Uint64())
	result = append(result, amountBytes...)

	return result, nil
}

// CancelOwnOrder creates a cancel own order instruction
func (f *FusionSwapContract) CancelOwnOrder(
	order *fusionorder.FusionOrder,
	accounts struct {
		Maker           *domains.Address
		SrcTokenProgram *domains.Address
	},
) (*types.TransactionInstruction, error) {
	orderHash := order.GetOrderHash()
	escrow, err := addresses.GetPda(f.ProgramID, [][]byte{
		[]byte("escrow"),
		accounts.Maker.ToBuffer(),
		orderHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow PDA: %w", err)
	}

	escrowSrcAta, err := addresses.GetAta(escrow, order.SrcMint(), accounts.SrcTokenProgram)
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow ATA: %w", err)
	}

	accountMetas := []types.AccountMeta{
		// 1. maker
		types.AccountMeta{Pubkey: accounts.Maker, IsSigner: true, IsWritable: true},
		// 2. src_mint
		types.AccountMeta{Pubkey: order.SrcMint(), IsSigner: false, IsWritable: false},
		// 3. escrow
		types.AccountMeta{Pubkey: escrow, IsSigner: false, IsWritable: false},
		// 4. escrow_src_ata
		types.AccountMeta{Pubkey: escrowSrcAta, IsSigner: false, IsWritable: true},
	}

	// 5. maker_src_ata (optional if native)
	if order.SrcAssetIsNative() {
		accountMetas = append(accountMetas, f.optionalAccountMeta())
	} else {
		makerSrcAta, err := addresses.GetAta(accounts.Maker, order.SrcMint(), accounts.SrcTokenProgram)
		if err != nil {
			return nil, fmt.Errorf("failed to get maker ATA: %w", err)
		}
		accountMetas = append(accountMetas, types.AccountMeta{
			Pubkey:     makerSrcAta,
			IsSigner:   false,
			IsWritable: true,
		})
	}

	accountMetas = append(accountMetas,
		// 6. src_token_program
		types.AccountMeta{Pubkey: accounts.SrcTokenProgram, IsSigner: false, IsWritable: false},
	)

	// Encode instruction data
	data := f.encodeCancelInstruction(order)

	return types.NewTransactionInstruction(f.ProgramID, accountMetas, data), nil
}

// encodeCancelInstruction encodes the cancel instruction data
func (f *FusionSwapContract) encodeCancelInstruction(order *fusionorder.FusionOrder) []byte {
	// Append order hash (32 bytes)
	orderHash := order.GetOrderHash()
	result := append(idl.CancelOrderDiscriminator, orderHash...)

	// Append orderSrcAssetIsNative (bool, 1 byte)
	if order.SrcAssetIsNative() {
		result = append(result, 1)
	} else {
		result = append(result, 0)
	}

	return result
}

// CancelOrderByResolver creates a cancel by resolver instruction
func (f *FusionSwapContract) CancelOrderByResolver(
	order *fusionorder.FusionOrder,
	accounts struct {
		Maker           *domains.Address
		Resolver        *domains.Address
		SrcTokenProgram *domains.Address
		Whitelist       *domains.Address
	},
	rewardLimit *big.Int,
) (*types.TransactionInstruction, error) {
	resolverConfig := order.ResolverCancellationConfig()
	if resolverConfig == nil {
		return nil, fmt.Errorf("order can not be cancelled by resolver")
	}
	if rewardLimit == nil {
		rewardLimit = resolverConfig.MaxCancellationPremium
	}

	whitelist := accounts.Whitelist
	if whitelist == nil {
		whitelist = WhitelistContractAddress
	}

	orderHash := order.GetOrderHash()
	escrow, err := addresses.GetPda(f.ProgramID, [][]byte{
		[]byte("escrow"),
		accounts.Maker.ToBuffer(),
		orderHash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow PDA: %w", err)
	}

	resolverAccess, err := addresses.GetPda(whitelist, [][]byte{
		[]byte("resolver_access"),
		accounts.Resolver.ToBuffer(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get resolver access PDA: %w", err)
	}

	escrowSrcAta, err := addresses.GetAta(escrow, order.SrcMint(), accounts.SrcTokenProgram)
	if err != nil {
		return nil, fmt.Errorf("failed to get escrow ATA: %w", err)
	}

	accountMetas := []types.AccountMeta{
		// 1. resolver
		types.AccountMeta{Pubkey: accounts.Resolver, IsSigner: true, IsWritable: true},
		// 2. resolver_access
		types.AccountMeta{Pubkey: resolverAccess, IsSigner: false, IsWritable: false},
		// 3. maker
		types.AccountMeta{Pubkey: accounts.Maker, IsSigner: false, IsWritable: true},
		// 4. maker_receiver
		types.AccountMeta{Pubkey: order.Receiver(), IsSigner: false, IsWritable: order.SrcAssetIsNative()},
		// 5. src_mint
		types.AccountMeta{Pubkey: order.SrcMint(), IsSigner: false, IsWritable: false},
		// 6. dst_mint
		types.AccountMeta{Pubkey: order.DstMint(), IsSigner: false, IsWritable: false},
		// 7. escrow
		types.AccountMeta{Pubkey: escrow, IsSigner: false, IsWritable: false},
		// 8. escrow_src_ata
		types.AccountMeta{Pubkey: escrowSrcAta, IsSigner: false, IsWritable: true},
	}

	// 9. maker_src_ata (optional if native)
	if order.SrcAssetIsNative() {
		accountMetas = append(accountMetas, f.optionalAccountMeta())
	} else {
		makerSrcAta, err := addresses.GetAta(accounts.Maker, order.SrcMint(), accounts.SrcTokenProgram)
		if err != nil {
			return nil, fmt.Errorf("failed to get maker ATA: %w", err)
		}
		accountMetas = append(accountMetas, types.AccountMeta{
			Pubkey:     makerSrcAta,
			IsSigner:   false,
			IsWritable: true,
		})
	}

	accountMetas = append(accountMetas,
		// 10. src_token_program
		types.AccountMeta{Pubkey: accounts.SrcTokenProgram, IsSigner: false, IsWritable: false},
		// 11. system_program
		types.AccountMeta{Pubkey: domains.SYSTEM_PROGRAM_ID, IsSigner: false, IsWritable: false},
	)

	// 12. protocol_dst_ata (optional)
	fees := order.Fees()
	if order.DstAssetIsNative() || fees == nil || fees.ProtocolDstAta == nil {
		accountMetas = append(accountMetas, f.optionalAccountMeta())
	} else {
		accountMetas = append(accountMetas, types.AccountMeta{
			Pubkey:     fees.ProtocolDstAta,
			IsSigner:   false,
			IsWritable: true,
		})
	}

	// 13. integrator_dst_ata (optional)
	if order.DstAssetIsNative() || fees == nil || fees.IntegratorDstAta == nil {
		accountMetas = append(accountMetas, f.optionalAccountMeta())
	} else {
		accountMetas = append(accountMetas, types.AccountMeta{
			Pubkey:     fees.IntegratorDstAta,
			IsSigner:   false,
			IsWritable: true,
		})
	}

	// Encode instruction data
	data, err := f.encodeCancelByResolverInstruction(order, rewardLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to encode cancel by resolver instruction: %w", err)
	}

	return types.NewTransactionInstruction(f.ProgramID, accountMetas, data), nil
}

// encodeCancelByResolverInstruction encodes the cancel by resolver instruction data
func (f *FusionSwapContract) encodeCancelByResolverInstruction(order *fusionorder.FusionOrder, rewardLimit *big.Int) ([]byte, error) {
	orderConfig := order.Build()
	borshData, err := orderConfig.SerializeBorsh()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize order config: %w", err)
	}

	// Append order config
	result := append(idl.CancelOrderByResolverDiscriminator, borshData...)

	// Append rewardLimit (u64, little-endian)
	rewardLimitBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(rewardLimitBytes, rewardLimit.Uint64())
	result = append(result, rewardLimitBytes...)

	return result, nil
}
