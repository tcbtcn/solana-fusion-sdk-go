package addresses

import (
	"fmt"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
	solana "github.com/gagliardetto/solana-go"
)

// GetPda generates a Program Derived Address (PDA) from seeds
func GetPda(programId domains.AddressLike, seeds [][]byte) (*domains.Address, error) {
	programKey := solana.PublicKeyFromBytes(programId.ToBuffer())

	// Find PDA using seeds
	pda, _, err := solana.FindProgramAddress(seeds, programKey)
	if err != nil {
		return nil, fmt.Errorf("failed to find program address: %w", err)
	}

	return domains.AddressFromBuffer(pda.Bytes()), nil
}
