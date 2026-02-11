package fusionorder

import (
	"errors"
	"math/big"
)

var (
	// Base1E3 is 1000
	Base1E3 = big.NewInt(1000)
)

// ZeroResolverCancellationConfig is a ResolverCancellationConfig with all zeros
var ZeroResolverCancellationConfig = &ResolverCancellationConfig{
	MaxCancellationPremium:      big.NewInt(0),
	CancellationAuctionDuration: 0,
}

// AlmostZeroResolverCancellationConfig is a ResolverCancellationConfig with minimal values
var AlmostZeroResolverCancellationConfig = &ResolverCancellationConfig{
	MaxCancellationPremium:      big.NewInt(1),
	CancellationAuctionDuration: 1,
}

// ResolverCancellationConfig represents resolver cancellation configuration
type ResolverCancellationConfig struct {
	MaxCancellationPremium      *big.Int
	CancellationAuctionDuration uint32
}

// NewResolverCancellationConfig creates a new ResolverCancellationConfig
func NewResolverCancellationConfig(
	maxCancellationPremium *big.Int,
	cancellationAuctionDuration uint32,
) (*ResolverCancellationConfig, error) {
	// Validate consistency
	if (maxCancellationPremium.Sign() == 0 && cancellationAuctionDuration != 0) ||
		(maxCancellationPremium.Sign() != 0 && cancellationAuctionDuration == 0) {
		return nil, errors.New("inconsistent cancellation config")
	}

	return &ResolverCancellationConfig{
		MaxCancellationPremium:      maxCancellationPremium,
		CancellationAuctionDuration: cancellationAuctionDuration,
	}, nil
}

// DisableResolverCancellation returns a zero ResolverCancellationConfig
func DisableResolverCancellation() *ResolverCancellationConfig {
	return ZeroResolverCancellationConfig
}

// IsZero checks if the config is zero
func (r *ResolverCancellationConfig) IsZero() bool {
	return r.MaxCancellationPremium.Sign() == 0
}
