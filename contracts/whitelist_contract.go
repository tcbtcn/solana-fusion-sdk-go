package contracts

import "github.com/dawitel/solana-fusion-sdk-go/domains"

// WhitelistContractAddress is the whitelist contract address
var WhitelistContractAddress = domains.MustAddressFromString("5jzZhrzqkbdwp5d3J1XbmaXMRnqeXimM1mDMoGHyvR7S")

// WhitelistContract represents the whitelist contract
type WhitelistContract struct {
	Address *domains.Address
}

// DefaultWhitelistContract returns the default whitelist contract
func DefaultWhitelistContract() *WhitelistContract {
	return &WhitelistContract{
		Address: WhitelistContractAddress,
	}
}
