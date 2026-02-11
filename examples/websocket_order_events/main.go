package main

import (
	"fmt"
	"log"

	wsapi "github.com/dawitel/solana-fusion-sdk-go/ws-api"
)

func main() {
	config := wsapi.WsApiConfig{
		URL:      "wss://api.1inch.dev/fusion/ws",
		AuthKey:  "your-auth-key-here", // Replace with your actual auth key
		LazyInit: false,
	}

	ws, err := wsapi.NewWebSocketApi(config)
	if err != nil {
		log.Fatalf("Failed to create WebSocket API: %v", err)
	}

	// Connection event handlers
	ws.OnOpen(func() {
		fmt.Println("WebSocket connected")
	})

	ws.OnClose(func() {
		fmt.Println("WebSocket disconnected")
	})

	ws.OnError(func(err error) {
		fmt.Printf("WebSocket error: %v\n", err)
	})

	// Order event handlers
	ws.Order.OnOrder(func(data wsapi.OrderEventType) {
		switch data.GetEvent() {
		case wsapi.EventTypeCreate:
			if created, ok := data.(*wsapi.OrderCreatedEvent); ok {
				fmt.Printf("Order created: %s\n", created.Result.OrderHash)
				fmt.Printf("Maker: %s\n", created.Result.Maker)
				fmt.Printf("Transaction: %s\n", created.Result.TransactionSignature)
			}
		case wsapi.EventTypeFill:
			if filled, ok := data.(*wsapi.OrderFilledEvent); ok {
				fmt.Printf("Order filled: %s\n", filled.Result.OrderHash)
				fmt.Printf("Resolver: %s\n", filled.Result.Resolver)
				fmt.Printf("Transaction: %s\n", filled.Result.TransactionSignature)
			}
		case wsapi.EventTypeCancel:
			if cancelled, ok := data.(*wsapi.OrderCancelledEvent); ok {
				fmt.Printf("Order cancelled: %s\n", cancelled.Result.OrderHash)
				fmt.Printf("Maker: %s\n", cancelled.Result.Maker)
			}
		}
	})

	// Specific event handlers
	ws.Order.OnOrderCreated(func(data *wsapi.OrderCreatedEvent) {
		fmt.Printf("Order created event: %s\n", data.Result.OrderHash)
	})

	ws.Order.OnOrderFilled(func(data *wsapi.OrderFilledEvent) {
		fmt.Printf("Order filled event: %s\n", data.Result.OrderHash)
	})

	ws.Order.OnOrderCancelled(func(data *wsapi.OrderCancelledEvent) {
		fmt.Printf("Order cancelled event: %s\n", data.Result.OrderHash)
	})

	// Keep connection alive
	fmt.Println("Listening for order events...")
	fmt.Println("Press Ctrl+C to exit")

	// Keep the program running
	select {}
}
