package sdk

import (
	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
	fusionorder "github.com/tcbtcn/solana-fusion-sdk-go/fusion-order"
)

// CancellableOrder represents an order that can be cancelled by a resolver
type CancellableOrder struct {
	Maker *domains.Address
	Order *fusionorder.FusionOrder
}
