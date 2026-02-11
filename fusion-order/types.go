package fusionorder

import (
	"math/big"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
)

// OrderInfoData represents order information data
type OrderInfoData struct {
	ID                 uint32
	SrcAmount          *big.Int // u64
	MinDstAmount       *big.Int // u64
	EstimatedDstAmount *big.Int // u64
	Receiver           *domains.Address
	SrcMint            *domains.Address
	DstMint            *domains.Address
}
