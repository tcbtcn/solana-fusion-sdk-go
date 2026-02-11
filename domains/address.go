package domains

import (
	"encoding/json"
	"errors"
	"fmt"

	solana "github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

// AddressLike is an interface for types that can be converted to an Address
type AddressLike interface {
	ToBuffer() []byte
}

// Address represents a Solana address (32-byte base58 encoded)
type Address struct {
	buf []byte
}

// Solana address constants
var (
	// ASSOCIATED_TOKEN_PROGRAM_ID is the Associated Token Program ID
	ASSOCIATED_TOKEN_PROGRAM_ID = MustAddressFromString("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")
	// TOKEN_PROGRAM_ID is the Token Program ID
	TOKEN_PROGRAM_ID = MustAddressFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	// TOKEN_2022_PROGRAM_ID is the Token 2022 Program ID
	TOKEN_2022_PROGRAM_ID = MustAddressFromString("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb")
	// SYSTEM_PROGRAM_ID is the System Program ID
	SYSTEM_PROGRAM_ID = MustAddressFromString("11111111111111111111111111111111")
	// WRAPPED_NATIVE is the wrapped SOL address
	WRAPPED_NATIVE = MustAddressFromString("So11111111111111111111111111111111111111112")
	// NATIVE is the native SOL address
	NATIVE = MustAddressFromString("SoNative11111111111111111111111111111111111")
)

// NewAddress creates a new Address from a base58 string
func NewAddress(value string) (*Address, error) {
	buf, err := base58.Decode(value)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid address: %w", value, err)
	}

	if len(buf) != 32 {
		return nil, fmt.Errorf("%s is not a valid address: expected 32 bytes, got %d", value, len(buf))
	}

	return &Address{buf: buf}, nil
}

// MustAddressFromString creates a new Address from a base58 string, panicking on error
func MustAddressFromString(value string) *Address {
	addr, err := NewAddress(value)
	if err != nil {
		panic(err)
	}
	return addr
}

// AddressFromUnknown creates an Address from various types
func AddressFromUnknown(val interface{}) (*Address, error) {
	if val == nil {
		return nil, errors.New("invalid address: value is nil")
	}

	switch v := val.(type) {
	case string:
		return NewAddress(v)
	case *Address:
		return v, nil
	case AddressLike:
		return AddressFromBuffer(v.ToBuffer()), nil
	case *solana.PublicKey:
		return AddressFromBuffer(v.Bytes()), nil
	default:
		return nil, fmt.Errorf("invalid address type: %T", val)
	}
}

// AddressFromPublicKey creates an Address from a PublicKey-like interface
func AddressFromPublicKey(publicKey AddressLike) *Address {
	return AddressFromBuffer(publicKey.ToBuffer())
}

// AddressFromBuffer creates an Address from a byte buffer
func AddressFromBuffer(buf []byte) *Address {
	return &Address{buf: buf}
}

// AddressFromBigInt creates an Address from a big integer (64-byte hex string)
func AddressFromBigInt(val string) (*Address, error) {
	// Convert hex string to bytes (pad to 32 bytes)
	if len(val) < 2 || val[0:2] != "0x" {
		return nil, errors.New("invalid hex string")
	}

	hexStr := val[2:]
	// Pad to 64 hex chars (32 bytes)
	for len(hexStr) < 64 {
		hexStr = "0" + hexStr
	}

	buf := make([]byte, 32)
	for i := 0; i < 32; i++ {
		var b byte
			_, _ = fmt.Sscanf(hexStr[i*2:(i+1)*2], "%02x", &b)
		buf[i] = b
	}

	return AddressFromBuffer(buf), nil
}

// ToString returns the base58 encoded string representation
func (a *Address) ToString() string {
	return base58.Encode(a.buf)
}

// String implements fmt.Stringer
func (a *Address) String() string {
	return a.ToString()
}

// ToBuffer returns the byte buffer representation
func (a *Address) ToBuffer() []byte {
	return a.buf
}

// Bytes returns the byte buffer (alias for ToBuffer)
func (a *Address) Bytes() []byte {
	return a.buf
}

// Equal checks if two addresses are equal
func (a *Address) Equal(other *Address) bool {
	if a == nil || other == nil {
		return a == other
	}
	if len(a.buf) != len(other.buf) {
		return false
	}
	for i := range a.buf {
		if a.buf[i] != other.buf[i] {
			return false
		}
	}
	return true
}

// IsNative checks if the address is the native SOL address
func (a *Address) IsNative() bool {
	return a.Equal(NATIVE)
}

// ToPublicKey converts the Address to a solana.PublicKey
func (a *Address) ToPublicKey() solana.PublicKey {
	return solana.PublicKey(a.buf)
}

// MarshalJSON implements json.Marshaler
func (a *Address) MarshalJSON() ([]byte, error) {
	return []byte(`"` + a.ToString() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (a *Address) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	addr, err := NewAddress(s)
	if err != nil {
		return err
	}
	*a = *addr
	return nil
}
