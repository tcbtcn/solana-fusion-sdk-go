package fusionorder

import (
	"errors"
	"math/big"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
)

var (
	// Base1E5 is 100% = 100000
	Base1E5 = big.NewInt(100_000)
	// Base1E2 is 100% = 100
	Base1E2 = big.NewInt(100)
)

// ZeroFeeConfig is a FeeConfig with all zeros
// Note: This is created without validation since all values are zero
var ZeroFeeConfig = &FeeConfig{
	ProtocolDstAta:   nil,
	IntegratorDstAta: nil,
	ProtocolFee:      domains.ZeroBps,
	IntegratorFee:    domains.ZeroBps,
	SurplusShare:     domains.ZeroBps,
}

// FeeConfig represents fee configuration
type FeeConfig struct {
	ProtocolDstAta   *domains.Address
	IntegratorDstAta *domains.Address
	ProtocolFee      *domains.Bps
	IntegratorFee    *domains.Bps
	SurplusShare     *domains.Bps
}

// NewFeeConfig creates a new FeeConfig
func NewFeeConfig(
	protocolDstAta *domains.Address,
	integratorDstAta *domains.Address,
	protocolFee *domains.Bps,
	integratorFee *domains.Bps,
	surplusShare *domains.Bps,
) (*FeeConfig, error) {
	// Validate protocol fee config (matches TypeScript logic exactly)
	// Invalid: (protocolDstAta is nil AND (protocolFee OR surplusShare is non-zero)) OR
	//          (protocolDstAta is not nil AND both protocolFee and surplusShare are zero)
	isProtocolFeeInvalid := (protocolDstAta == nil &&
		!(protocolFee.IsZero() && surplusShare.IsZero())) ||
		(protocolDstAta != nil &&
			protocolFee.IsZero() &&
			surplusShare.IsZero())

	if isProtocolFeeInvalid {
		return nil, errors.New("protocol fee config mismatch")
	}

	// Validate integrator fee config (matches TypeScript logic exactly)
	// Invalid: (integratorDstAta is nil AND integratorFee is non-zero) OR
	//          (integratorDstAta is not nil AND integratorFee is zero)
	isIntegratorFeeInvalid := (integratorDstAta == nil && !integratorFee.IsZero()) ||
		(integratorDstAta != nil && integratorFee.IsZero())

	if isIntegratorFeeInvalid {
		return nil, errors.New("integrator fee config mismatch")
	}

	// Validate max fees
	if protocolFee.ToFraction(nil) >= 0.6553 {
		return nil, errors.New("max fee is 65.53%")
	}
	if integratorFee.ToFraction(nil) >= 0.6553 {
		return nil, errors.New("max fee is 65.53%")
	}
	if surplusShare.ToFraction(nil) > 1.0 {
		return nil, errors.New("max surplus share is 100%")
	}

	// Surplus share must have percent precision
	if new(big.Int).Mod(surplusShare.Value(), big.NewInt(100)).Sign() != 0 {
		return nil, errors.New("surplus share must have percent precision: 1%, 2% and so on")
	}

	return &FeeConfig{
		ProtocolDstAta:   protocolDstAta,
		IntegratorDstAta: integratorDstAta,
		ProtocolFee:      protocolFee,
		IntegratorFee:    integratorFee,
		SurplusShare:     surplusShare,
	}, nil
}

// OnlyProtocol creates a FeeConfig with only protocol fees
func OnlyProtocol(protocolDstAta *domains.Address, protocolFee, surplusShare *domains.Bps) (*FeeConfig, error) {
	return NewFeeConfig(protocolDstAta, nil, protocolFee, domains.ZeroBps, surplusShare)
}

// OnlyIntegrator creates a FeeConfig with only integrator fees
func OnlyIntegrator(integratorDstAta *domains.Address, fee *domains.Bps) (*FeeConfig, error) {
	return NewFeeConfig(nil, integratorDstAta, domains.ZeroBps, fee, domains.ZeroBps)
}

// IsZero checks if the fee config is zero
func (f *FeeConfig) IsZero() bool {
	return f.ProtocolDstAta == nil && f.IntegratorDstAta == nil
}
