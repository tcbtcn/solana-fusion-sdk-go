package types

import "github.com/dawitel/solana-fusion-sdk-go/domains"

// AccountMeta represents account metadata for a transaction instruction
type AccountMeta struct {
	Pubkey     *domains.Address
	IsSigner   bool
	IsWritable bool
}

// TransactionInstruction represents a Solana transaction instruction
type TransactionInstruction struct {
	ProgramID *domains.Address
	Accounts  []AccountMeta
	Data      []byte
}

// NewTransactionInstruction creates a new TransactionInstruction
func NewTransactionInstruction(
	programID *domains.Address,
	accounts []AccountMeta,
	data []byte,
) *TransactionInstruction {
	return &TransactionInstruction{
		ProgramID: programID,
		Accounts:  accounts,
		Data:      data,
	}
}
