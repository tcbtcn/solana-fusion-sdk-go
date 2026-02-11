package fusionorder

import (
	"math/big"
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/domains"
)

func TestNewFeeConfig_Success(t *testing.T) {
	protocolDstAta := domains.MustAddressFromString("11111111111111111111111111111111")
	protocolFee := domains.BpsFromPercent(1.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)

	feeConfig, err := NewFeeConfig(protocolDstAta, nil, protocolFee, domains.ZeroBps, surplusShare)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if feeConfig == nil {
		t.Fatal("Expected non-nil fee config")
	}
	if !feeConfig.ProtocolDstAta.Equal(protocolDstAta) {
		t.Error("Expected ProtocolDstAta to match")
	}
}

func TestNewFeeConfig_WithIntegrator(t *testing.T) {
	integratorDstAta := domains.MustAddressFromString("11111111111111111111111111111112")
	integratorFee := domains.BpsFromPercent(2.0, nil)

	feeConfig, err := NewFeeConfig(nil, integratorDstAta, domains.ZeroBps, integratorFee, domains.ZeroBps)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if feeConfig == nil {
		t.Fatal("Expected non-nil fee config")
	}
	if !feeConfig.IntegratorDstAta.Equal(integratorDstAta) {
		t.Error("Expected IntegratorDstAta to match")
	}
}

func TestNewFeeConfig_ProtocolFeeMismatch(t *testing.T) {
	protocolDstAta := domains.MustAddressFromString("11111111111111111111111111111111")
	// Invalid: protocolDstAta is provided but both protocolFee and surplusShare are zero
	_, err := NewFeeConfig(protocolDstAta, nil, domains.ZeroBps, domains.ZeroBps, domains.ZeroBps)
	if err == nil {
		t.Fatal("Expected error for protocol fee config mismatch")
	}
}

func TestNewFeeConfig_IntegratorFeeMismatch(t *testing.T) {
	integratorFee := domains.BpsFromPercent(2.0, nil)

	_, err := NewFeeConfig(nil, nil, domains.ZeroBps, integratorFee, domains.ZeroBps)
	if err == nil {
		t.Fatal("Expected error for integrator fee config mismatch")
	}
}

func TestNewFeeConfig_MaxFeeExceeded(t *testing.T) {
	protocolFee := domains.BpsFromPercent(66.0, nil)

	_, err := NewFeeConfig(nil, nil, protocolFee, domains.ZeroBps, domains.ZeroBps)
	if err == nil {
		t.Fatal("Expected error for max fee exceeded")
	}
}

func TestNewFeeConfig_SurplusSharePrecision(t *testing.T) {
	protocolFee := domains.BpsFromPercent(1.0, nil)
	surplusShare, _ := domains.NewBps(big.NewInt(150))

	_, err := NewFeeConfig(nil, nil, protocolFee, domains.ZeroBps, surplusShare)
	if err == nil {
		t.Fatal("Expected error for surplus share precision")
	}
}

func TestFeeConfig_IsZero(t *testing.T) {
	zeroConfig := &FeeConfig{
		ProtocolDstAta:   nil,
		IntegratorDstAta: nil,
	}
	if !zeroConfig.IsZero() {
		t.Error("Expected zero config to be zero")
	}

	nonZeroConfig := &FeeConfig{
		ProtocolDstAta:   domains.MustAddressFromString("11111111111111111111111111111111"),
		IntegratorDstAta: nil,
	}
	if nonZeroConfig.IsZero() {
		t.Error("Expected non-zero config to not be zero")
	}
}

func TestOnlyProtocol(t *testing.T) {
	protocolDstAta := domains.MustAddressFromString("11111111111111111111111111111111")
	protocolFee := domains.BpsFromPercent(1.0, nil)
	surplusShare := domains.BpsFromPercent(50.0, nil)

	feeConfig, err := OnlyProtocol(protocolDstAta, protocolFee, surplusShare)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if feeConfig.IntegratorDstAta != nil {
		t.Error("Expected IntegratorDstAta to be nil")
	}
	if !feeConfig.IntegratorFee.IsZero() {
		t.Error("Expected IntegratorFee to be zero")
	}
}

func TestOnlyIntegrator(t *testing.T) {
	integratorDstAta := domains.MustAddressFromString("11111111111111111111111111111112")
	integratorFee := domains.BpsFromPercent(2.0, nil)

	feeConfig, err := OnlyIntegrator(integratorDstAta, integratorFee)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if feeConfig.ProtocolDstAta != nil {
		t.Error("Expected ProtocolDstAta to be nil")
	}
	if !feeConfig.ProtocolFee.IsZero() {
		t.Error("Expected ProtocolFee to be zero")
	}
}
