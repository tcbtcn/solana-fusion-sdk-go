package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/tcbtcn/solana-fusion-sdk-go/api"
	"github.com/tcbtcn/solana-fusion-sdk-go/api/http"
	"github.com/tcbtcn/solana-fusion-sdk-go/domains"
	"github.com/tcbtcn/solana-fusion-sdk-go/sdk"
)

func main() {
	config := api.ApiConfig{
		BaseURL: "https://api.1inch.dev/fusion",
		AuthKey: "your-auth-key-here", // Replace with your actual auth key
		Version: "v1.0",
	}

	// Create HTTP client - it needs the full base URL
	httpClient := http.NewHTTPClient(config.BaseURL, config.AuthKey)
	client := sdk.NewSdk(httpClient, config)
	ctx := context.Background()

	// Create order: 1 SOL -> USDC
	srcToken := domains.NATIVE                                                                // SOL
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v") // USDC
	amount := big.NewInt(1_000_000_000)                                                       // 1 SOL (in lamports)
	signer := domains.MustAddressFromString("YourWalletAddressHere")                          // Replace with your wallet address

	fmt.Println("Creating order...")
	order, err := client.CreateOrder(ctx, srcToken, dstToken, amount, signer, nil)
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}

	fmt.Printf("Order created successfully!\n")
	fmt.Printf("Order Hash (Base58): %s\n", order.GetOrderHashBase58())
	fmt.Printf("Order ID: %d\n", order.ID())
	fmt.Printf("Source Amount: %s\n", order.SrcAmount().String())
	fmt.Printf("Min Destination Amount: %s\n", order.MinDstAmount().String())
	fmt.Printf("Estimated Destination Amount: %s\n", order.EstimatedDstAmount().String())
}
