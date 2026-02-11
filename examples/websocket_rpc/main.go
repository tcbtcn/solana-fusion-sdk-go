package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dawitel/solana-fusion-sdk-go/api"
	"github.com/dawitel/solana-fusion-sdk-go/ws-api"
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

		// Get allowed methods
		fmt.Println("Requesting allowed methods...")
		if err := ws.RPC.GetAllowedMethods(); err != nil {
			fmt.Printf("Error requesting allowed methods: %v\n", err)
		}

		// Get active orders
		page := 1
		limit := 10
		fmt.Printf("Requesting active orders (page %d, limit %d)...\n", page, limit)
		if err := ws.RPC.GetActiveOrders(&page, &limit); err != nil {
			fmt.Printf("Error requesting active orders: %v\n", err)
		}

		// Send ping
		fmt.Println("Sending ping...")
		ws.RPC.Ping()
	})

	ws.OnClose(func() {
		fmt.Println("WebSocket disconnected")
	})

	ws.OnError(func(err error) {
		fmt.Printf("WebSocket error: %v\n", err)
	})

	// RPC response handlers
	ws.RPC.OnGetAllowedMethods(func(result []string) {
		fmt.Printf("Allowed methods: %v\n", result)
	})

	ws.RPC.OnGetActiveOrders(func(result interface{}) {
		// Type assert to get the actual pagination result
		if pagination, ok := result.(api.Pagination[wsapi.ActiveOrderWS]); ok {
			fmt.Printf("Active orders received:\n")
			fmt.Printf("Total items: %d\n", pagination.Meta.TotalItems)
			fmt.Printf("Items per page: %d\n", pagination.Meta.ItemsPerPage)
			fmt.Printf("Total pages: %d\n", pagination.Meta.TotalPages)
			fmt.Printf("Current page: %d\n", pagination.Meta.CurrentPage)
			fmt.Printf("Number of items: %d\n", len(pagination.Items))

			for i, order := range pagination.Items {
				fmt.Printf("  Order %d: %s (Maker: %s)\n", i+1, order.OrderHash, order.Maker)
			}
		} else {
			fmt.Printf("Received active orders (unexpected type): %+v\n", result)
		}
	})

	ws.RPC.OnPong(func() {
		fmt.Println("Pong received")
	})

	// Keep connection alive and send periodic pings
	fmt.Println("WebSocket RPC client started")
	fmt.Println("Press Ctrl+C to exit")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Sending periodic ping...")
			ws.RPC.Ping()
		}
	}
}
