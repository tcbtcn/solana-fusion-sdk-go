package domains

import (
	"errors"
	"fmt"
	"math/big"
)

// Bps represents basis points in range [0, 10000] (0-100%)
// 1 bps = 0.01%
type Bps struct {
	value *big.Int
}

// ZeroBps is a Bps with value 0
var ZeroBps = MustBps(big.NewInt(0))

// NewBps creates a new Bps value
// value must be in range [0, 10000]
func NewBps(value *big.Int) (*Bps, error) {
	if value == nil {
		return nil, errors.New("bps value cannot be nil")
	}
	if value.Sign() < 0 {
		return nil, fmt.Errorf("invalid bps %s: must be >= 0", value.String())
	}
	if value.Cmp(big.NewInt(10000)) > 0 {
		return nil, fmt.Errorf("invalid bps %s: must be <= 10000", value.String())
	}
	return &Bps{value: new(big.Int).Set(value)}, nil
}

// MustBps creates a new Bps value, panicking on error
func MustBps(value *big.Int) *Bps {
	bps, err := NewBps(value)
	if err != nil {
		panic(err)
	}
	return bps
}

// BpsFromPercent creates BPS from percent value
// If value has precision more than 1bps (with accounting to base), it will be rounded down
func BpsFromPercent(val float64, base *big.Int) *Bps {
	if base == nil {
		base = big.NewInt(1)
	}
	// Calculate: (val * 100) / base
	numerator := big.NewInt(int64(val * 100))
	result := new(big.Int).Div(numerator, base)
	return MustBps(result)
}

// BpsFromFraction creates BPS from fraction value
// If value has precision more than 1bps (with accounting to base), it will be rounded down
func BpsFromFraction(val float64, base *big.Int) *Bps {
	if base == nil {
		base = big.NewInt(1)
	}
	// Calculate: (val * 10000) / base
	numerator := big.NewInt(int64(val * 10000))
	result := new(big.Int).Div(numerator, base)
	return MustBps(result)
}

// Value returns the basis points value
func (b *Bps) Value() *big.Int {
	return new(big.Int).Set(b.value)
}

// Equal checks if two Bps values are equal
func (b *Bps) Equal(other *Bps) bool {
	if b == nil || other == nil {
		return b == other
	}
	return b.value.Cmp(other.value) == 0
}

// IsZero checks if the Bps value is zero
func (b *Bps) IsZero() bool {
	return b.value.Sign() == 0
}

// ToPercent converts BPS to percent value
// base represents what 100% is
func (b *Bps) ToPercent(base *big.Int) float64 {
	if base == nil {
		base = big.NewInt(1)
	}
	// Calculate: (value * base) / 100
	// Use floating point division to preserve precision
	valueFloat := float64(b.value.Int64())
	baseFloat := float64(base.Int64())
	return (valueFloat * baseFloat) / 100.0
}

// ToFraction converts BPS to fraction value
// base represents what 100% is
func (b *Bps) ToFraction(base *big.Int) float64 {
	if base == nil {
		base = big.NewInt(1)
	}
	// Calculate: (value * base) / 10000
	// Use floating point division to preserve precision
	valueFloat := float64(b.value.Int64())
	baseFloat := float64(base.Int64())
	return (valueFloat * baseFloat) / 10000.0
}

// String returns the string representation
func (b *Bps) String() string {
	return b.value.String()
}
