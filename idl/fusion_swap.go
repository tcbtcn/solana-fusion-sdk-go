package idl

// FusionSwap program address
const FusionSwapProgramAddress = "HNarfxC3kYMMhFkxUFeYb8wHVdPzY5t9pupqW5fL2meM"

// Instruction discriminators for FusionSwap program
// These are the first 8 bytes of sha256("global:<instruction_name>")
var (
	// CreateOrderDiscriminator is the discriminator for the "create" instruction
	CreateOrderDiscriminator = []byte{24, 30, 200, 40, 5, 28, 7, 119}

	// FillOrderDiscriminator is the discriminator for the "fill" instruction
	FillOrderDiscriminator = []byte{168, 96, 183, 163, 92, 10, 40, 160}

	// CancelOrderDiscriminator is the discriminator for the "cancel" instruction
	CancelOrderDiscriminator = []byte{232, 219, 223, 41, 219, 236, 220, 190}

	// CancelOrderByResolverDiscriminator is the discriminator for the "cancelByResolver" instruction
	CancelOrderByResolverDiscriminator = []byte{229, 180, 171, 131, 171, 6, 60, 191}
)
