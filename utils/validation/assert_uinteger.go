package validation

import (
	"fmt"
	"math"
	"math/big"
)

const (
	// Uint32Max is the maximum value for uint32
	Uint32Max = math.MaxUint32
)

// AssertUInteger asserts that a value is a non-negative integer within bounds
func AssertUInteger(val interface{}, max *big.Int) error {
	if max == nil {
		max = big.NewInt(Uint32Max)
	}

	var value *big.Int

	switch v := val.(type) {
	case int:
		value = big.NewInt(int64(v))
	case int64:
		value = big.NewInt(v)
	case uint:
		value = new(big.Int).SetUint64(uint64(v))
	case uint8:
		value = big.NewInt(int64(v))
	case uint16:
		value = big.NewInt(int64(v))
	case uint32:
		value = new(big.Int).SetUint64(uint64(v))
	case uint64:
		value = new(big.Int).SetUint64(v)
	case *big.Int:
		value = v
	default:
		return fmt.Errorf("expected integer type, got %T", val)
	}

	if value.Sign() < 0 {
		return fmt.Errorf("expected %s to be >= 0", value.String())
	}

	if value.Cmp(max) > 0 {
		return fmt.Errorf("expected %s to be <= %s", value.String(), max.String())
	}

	return nil
}
