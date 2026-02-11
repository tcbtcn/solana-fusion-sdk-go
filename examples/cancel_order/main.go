package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dawitel/solana-fusion-sdk-go/api"
	"github.com/dawitel/solana-fusion-sdk-go/api/http"
	"github.com/dawitel/solana-fusion-sdk-go/contracts"
	"github.com/dawitel/solana-fusion-sdk-go/domains"
	"github.com/dawitel/solana-fusion-sdk-go/sdk"
)

func main() {
	config := api.ApiConfig{
		BaseURL: "https://api.1inch.dev/fusion",
		AuthKey: "your-auth-key-here", // Replace with your actual auth key
		Version: "v1.0",
	}

	httpClient := http.NewHTTPClient(config.BaseURL, config.AuthKey)
	client := sdk.NewSdk(httpClient, config)
	ctx := context.Background()

	// Get order status first
	orderHash := "YourOrderHashHere" // Replace with actual order hash
	fmt.Printf("Getting order status for: %s\n", orderHash)

	status, err := client.GetOrderStatus(ctx, orderHash)
	if err != nil {
		log.Fatalf("Failed to get order status: %v", err)
	}

	if !status.Cancelable {
		log.Fatal("Order is not cancelable")
	}

	fmt.Println("Order is cancelable, creating cancel instruction...")

	// Create cancel instruction
	contract := contracts.DefaultFusionSwapContract()
	maker := status.Maker

	instruction, err := contract.CancelOwnOrder(status.Order, struct {
		Maker           *domains.Address
		SrcTokenProgram *domains.Address
	}{
		Maker:           maker,
		SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
	})
	if err != nil {
		log.Fatalf("Failed to create cancel instruction: %v", err)
	}

	fmt.Printf("Cancel instruction created!\n")
	fmt.Printf("Program ID: %s\n", instruction.ProgramID.ToString())
	fmt.Printf("Number of accounts: %d\n", len(instruction.Accounts))
	fmt.Printf("Instruction data length: %d bytes\n", len(instruction.Data))

	// Note: In a real application, you would:
	// 1. Build a Solana transaction with this instruction
	// 2. Sign the transaction with the maker's private key
	// 3. Send the transaction to the Solana network
	fmt.Println("\nNext steps:")
	fmt.Println("1. Build a Solana transaction with this instruction")
	fmt.Println("2. Sign the transaction with the maker's private key")
	fmt.Println("3. Send the transaction to the Solana network")
}
