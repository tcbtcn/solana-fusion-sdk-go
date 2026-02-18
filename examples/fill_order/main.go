package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tcbtcn/solana-fusion-sdk-go/api"
	"github.com/tcbtcn/solana-fusion-sdk-go/api/http"
	"github.com/tcbtcn/solana-fusion-sdk-go/contracts"
	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
	"github.com/tcbtcn/solana-fusion-sdk-go/sdk"
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

	// Get active orders
	fmt.Println("Fetching active orders...")
	activeOrders, err := client.GetActiveOrders(ctx, 1, 10)
	if err != nil {
		log.Fatalf("Failed to get active orders: %v", err)
	}

	if len(activeOrders.Items) == 0 {
		log.Fatal("No active orders found")
	}

	// Use the first active order
	activeOrder := activeOrders.Items[0]
	order := activeOrder.Order

	fmt.Printf("Found active order:\n")
	fmt.Printf("Order Hash: %s\n", order.GetOrderHashBase58())
	fmt.Printf("Maker: %s\n", activeOrder.Maker.ToString())
	fmt.Printf("Remaining Maker Amount: %s\n", activeOrder.RemainingMakerAmount.String())

	// Create fill instruction
	contract := contracts.DefaultFusionSwapContract()
	taker := domains.MustAddressFromString("YourTakerAddressHere") // Replace with resolver/taker address

	// Fill for the remaining amount
	fillAmount := activeOrder.RemainingMakerAmount

	instruction, err := contract.Fill(order, fillAmount, struct {
		Taker           *domains.Address
		Maker           *domains.Address
		SrcTokenProgram *domains.Address
		DstTokenProgram *domains.Address
		TakerSrcAccount *domains.Address
		Whitelist       *domains.Address
	}{
		Taker:           taker,
		Maker:           activeOrder.Maker,
		SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
		DstTokenProgram: domains.TOKEN_PROGRAM_ID,
		TakerSrcAccount: nil, // Will be calculated automatically
		Whitelist:       nil, // Will use default whitelist
	})
	if err != nil {
		log.Fatalf("Failed to create fill instruction: %v", err)
	}

	fmt.Printf("\nFill instruction created!\n")
	fmt.Printf("Program ID: %s\n", instruction.ProgramID.ToString())
	fmt.Printf("Number of accounts: %d\n", len(instruction.Accounts))
	fmt.Printf("Instruction data length: %d bytes\n", len(instruction.Data))
	fmt.Printf("Fill amount: %s\n", fillAmount.String())

	// Note: In a real application, you would:
	// 1. Build a Solana transaction with this instruction
	// 2. Sign the transaction with the taker's private key
	// 3. Send the transaction to the Solana network
	fmt.Println("\nNext steps:")
	fmt.Println("1. Build a Solana transaction with this instruction")
	fmt.Println("2. Sign the transaction with the taker's private key")
	fmt.Println("3. Send the transaction to the Solana network")
}
