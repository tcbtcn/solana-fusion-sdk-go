package fusionorder

import (
	"math/big"
	"testing"
)

func TestNewResolverCancellationConfig_Success(t *testing.T) {
	maxCancellationPremium := big.NewInt(1000000)
	cancellationAuctionDuration := uint32(100)

	config, err := NewResolverCancellationConfig(maxCancellationPremium, cancellationAuctionDuration)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config")
	}
	if config.MaxCancellationPremium.Cmp(maxCancellationPremium) != 0 {
		t.Errorf("Expected MaxCancellationPremium %s, got %s", maxCancellationPremium.String(), config.MaxCancellationPremium.String())
	}
	if config.CancellationAuctionDuration != cancellationAuctionDuration {
		t.Errorf("Expected CancellationAuctionDuration %d, got %d", cancellationAuctionDuration, config.CancellationAuctionDuration)
	}
}

func TestNewResolverCancellationConfig_Zero(t *testing.T) {
	maxCancellationPremium := big.NewInt(0)
	cancellationAuctionDuration := uint32(0)

	config, err := NewResolverCancellationConfig(maxCancellationPremium, cancellationAuctionDuration)
	if err != nil {
		t.Fatalf("Expected no error for zero config, got %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config")
	}
}

func TestNewResolverCancellationConfig_Inconsistent(t *testing.T) {
	maxCancellationPremium := big.NewInt(0)
	cancellationAuctionDuration := uint32(100)

	_, err := NewResolverCancellationConfig(maxCancellationPremium, cancellationAuctionDuration)
	if err == nil {
		t.Fatal("Expected error for inconsistent config")
	}
}

func TestNewResolverCancellationConfig_Inconsistent2(t *testing.T) {
	maxCancellationPremium := big.NewInt(1000000)
	cancellationAuctionDuration := uint32(0)

	_, err := NewResolverCancellationConfig(maxCancellationPremium, cancellationAuctionDuration)
	if err == nil {
		t.Fatal("Expected error for inconsistent config")
	}
}

func TestResolverCancellationConfig_IsZero(t *testing.T) {
	zeroConfig := &ResolverCancellationConfig{
		MaxCancellationPremium:      big.NewInt(0),
		CancellationAuctionDuration: 0,
	}
	if !zeroConfig.IsZero() {
		t.Error("Expected zero config to be zero")
	}

	nonZeroConfig := &ResolverCancellationConfig{
		MaxCancellationPremium:      big.NewInt(1000000),
		CancellationAuctionDuration: 100,
	}
	if nonZeroConfig.IsZero() {
		t.Error("Expected non-zero config to not be zero")
	}
}

func TestDisableResolverCancellation(t *testing.T) {
	config := DisableResolverCancellation()
	if config == nil {
		t.Fatal("Expected non-nil config")
	}
	if !config.IsZero() {
		t.Error("Expected disabled config to be zero")
	}
}
