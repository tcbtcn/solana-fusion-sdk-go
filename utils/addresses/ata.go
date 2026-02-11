package addresses

import (
	"fmt"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
	solana "github.com/gagliardetto/solana-go"
)

var (
	// SPLAssociatedTokenProgramID is the Associated Token Program ID
	SPLAssociatedTokenProgramID = solana.MustPublicKeyFromBase58("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")
)

// GetAta returns the associated token account for given params
func GetAta(
	walletAddress domains.AddressLike,
	tokenMintAddress domains.AddressLike,
	tokenProgramId domains.AddressLike,
) (*domains.Address, error) {
	seeds := [][]byte{
		walletAddress.ToBuffer(),
		tokenProgramId.ToBuffer(),
		tokenMintAddress.ToBuffer(),
	}

	pda, _, err := solana.FindProgramAddress(seeds, SPLAssociatedTokenProgramID)
	if err != nil {
		return nil, fmt.Errorf("failed to find associated token account: %w", err)
	}

	return domains.AddressFromBuffer(pda.Bytes()), nil
}
