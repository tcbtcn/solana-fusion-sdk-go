# Solana Fusion SDK for Go

Production-grade Go SDK for creating and managing fusion orders on Solana through 1inch Fusion API.

## Installation

```bash
go get github.com/tcbtcn/solana-fusion-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
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
        AuthKey: "your-auth-key",
        Version: "v1.0",
    }

    httpClient := http.NewHTTPClient(config.BaseURL, config.AuthKey)
    client := sdk.NewSdk(httpClient, config)
    ctx := context.Background()

    // Create order: 1 SOL -> USDC
    srcToken := domains.NATIVE // SOL
    dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v") // USDC
    amount := big.NewInt(1_000_000_000) // 1 SOL
    signer := domains.MustAddressFromString("YourWalletAddress")

    order, err := client.CreateOrder(ctx, srcToken, dstToken, amount, signer, nil)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Order created: %s", order.GetOrderHashBase58())
}
```

## Features

- Quote generation for Solana token swaps
- Order creation and management
- Order status tracking
- Active orders listing
- Resolver cancellation support
- **WebSocket API** for real-time order updates and RPC functionality

## API Reference

### SDK Client

#### NewSdk

Creates a new SDK client instance.

```go
client := sdk.NewSdk(httpClient, api.ApiConfig{
    BaseURL: "https://api.1inch.dev/fusion",
    AuthKey: "your-auth-key",
    Version: "v1.0",
})
```

#### GetQuote

Gets a quote for a Solana token swap.

```go
quote, err := client.GetQuote(ctx, srcToken, dstToken, amount, signer, slippage)
```

#### CreateOrder

Creates a fusion order from a quote.

```go
order, err := client.CreateOrder(ctx, srcToken, dstToken, amount, signer, slippage)
```

#### GetOrderStatus

Gets the status of an order.

```go
status, err := client.GetOrderStatus(ctx, orderHash)
```

#### GetActiveOrders

Gets active orders with pagination.

```go
orders, err := client.GetActiveOrders(ctx, page, limit)
```

#### GetOrdersCancellableByResolver

Gets orders that can be cancelled by a resolver.

```go
orders, err := client.GetOrdersCancellableByResolver(ctx, page, limit)
```

### WebSocket API

The WebSocket API provides real-time order updates and RPC functionality.

#### NewWebSocketApi

Creates a new WebSocket API client.

```go
import "github.com/tcbtcn/solana-fusion-sdk-go/ws-api"

ws, err := wsapi.NewWebSocketApi(wsapi.WsApiConfig{
    URL:     "wss://api.1inch.dev/fusion/ws",
    AuthKey: "your-auth-key",
    LazyInit: false, // Set to true for lazy initialization
})
if err != nil {
    log.Fatal(err)
}

// Initialize connection (if lazyInit was true)
ws.Init()
```

#### Order Events

Subscribe to real-time order updates.

```go
// Subscribe to all order events
ws.Order.OnOrder(func(data wsapi.OrderEventType) {
    switch data.GetEvent() {
    case wsapi.EventTypeCreate:
        if created, ok := data.(*wsapi.OrderCreatedEvent); ok {
            fmt.Printf("Order created: %s\n", created.Result.OrderHash)
        }
    case wsapi.EventTypeFill:
        if filled, ok := data.(*wsapi.OrderFilledEvent); ok {
            fmt.Printf("Order filled: %s\n", filled.Result.OrderHash)
        }
    case wsapi.EventTypeCancel:
        if cancelled, ok := data.(*wsapi.OrderCancelledEvent); ok {
            fmt.Printf("Order cancelled: %s\n", cancelled.Result.OrderHash)
        }
    }
})

// Subscribe to specific order events
ws.Order.OnOrderCreated(func(data *wsapi.OrderCreatedEvent) {
    fmt.Printf("Order created: %s\n", data.Result.OrderHash)
})

ws.Order.OnOrderFilled(func(data *wsapi.OrderFilledEvent) {
    fmt.Printf("Order filled: %s\n", data.Result.OrderHash)
})

ws.Order.OnOrderCancelled(func(data *wsapi.OrderCancelledEvent) {
    fmt.Printf("Order cancelled: %s\n", data.Result.OrderHash)
})
```

#### RPC Methods

Call RPC methods for querying data.

```go
// Ping/Pong
ws.RPC.Ping()
ws.RPC.OnPong(func() {
    fmt.Println("Pong received")
})

// Get active orders
page := 1
limit := 10
ws.RPC.GetActiveOrders(&page, &limit)
ws.RPC.OnGetActiveOrders(func(result interface{}) {
    if pagination, ok := result.(api.Pagination[wsapi.ActiveOrderWS]); ok {
        fmt.Printf("Active orders: %d\n", len(pagination.Items))
    }
})

// Get allowed methods
ws.RPC.GetAllowedMethods()
ws.RPC.OnGetAllowedMethods(func(result []string) {
    fmt.Printf("Allowed methods: %v\n", result)
})
```

#### Connection Management

Handle WebSocket connection lifecycle.

```go
// Connection events
ws.OnOpen(func() {
    fmt.Println("WebSocket connected")
})

ws.OnClose(func() {
    fmt.Println("WebSocket disconnected")
})

ws.OnError(func(err error) {
    fmt.Printf("WebSocket error: %v\n", err)
})

// Close connection
defer ws.Close()
```

#### WebSocket Features

- **Order Events**: Subscribe to real-time order updates (created, filled, cancelled)
- **RPC Methods**: Call RPC methods like `getActiveOrders`, `getAllowedMethods`, `ping`
- **Connection Management**: Handle connection lifecycle (open, close, error)
- **Type Safety**: Strongly typed event structures matching the TypeScript SDK

## Contracts

### FusionSwapContract

The SDK provides contract instruction builders for interacting with the FusionSwap program on Solana.

#### Create

Creates a create order instruction.

```go
contract := contracts.DefaultFusionSwapContract()
instruction := contract.Create(order, struct {
    Maker         *domains.Address
    SrcTokenProgram *domains.Address
}{
    Maker: makerAddress,
    SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
})
```

#### Fill

Creates a fill order instruction.

```go
instruction := contract.Fill(order, amount, struct {
    Taker          *domains.Address
    Maker          *domains.Address
    SrcTokenProgram *domains.Address
    DstTokenProgram *domains.Address
}{
    Taker: takerAddress,
    Maker: makerAddress,
    SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
    DstTokenProgram: domains.TOKEN_PROGRAM_ID,
})
```

#### CancelOwnOrder

Creates a cancel own order instruction.

```go
instruction := contract.CancelOwnOrder(order, struct {
    Maker         *domains.Address
    SrcTokenProgram *domains.Address
}{
    Maker: makerAddress,
    SrcTokenProgram: domains.TOKEN_PROGRAM_ID,
})
```

## Examples

See the `examples/` directory for complete working examples:

- `examples/create_order/` - Create and submit order example
- `examples/cancel_order/` - Cancel order example
- `examples/fill_order/` - Fill order example (resolver)
- `examples/websocket_order_events/` - WebSocket API example for subscribing to order events
- `examples/websocket_rpc/` - WebSocket API example for RPC methods

## Instruction Decoding

The SDK supports decoding FusionOrders from Solana transaction instructions:

```go
import (
    "github.com/tcbtcn/solana-fusion-sdk-go/fusionorder"
    "github.com/tcbtcn/solana-fusion-sdk-go/types"
)

// Decode from create instruction
order, err := fusionorder.FromCreateInstruction(instruction)

// Decode from fill instruction
order, err := fusionorder.FromFillInstruction(instruction)

// Decode from cancelByResolver instruction
order, err := fusionorder.FromResolverCancelInstruction(instruction)

// Decode from contract order config
order, err := fusionorder.FromContractOrder(config, accounts)
```

## Amount Calculator

Calculate fees and amounts with auction accounting:

```go
import (
    amountcalculator "github.com/tcbtcn/solana-fusion-sdk-go/amount-calculator"
    fusionorder "github.com/tcbtcn/solana-fusion-sdk-go/fusion-order"
)

// Create calculator from order
info := order.GetCalculatorInfo()
auctionCalc := amountcalculator.FromAuctionData(info.AuctionDetails)
feeCalc := amountcalculator.FromFeeConfig(info.Fees)
calculator := amountcalculator.NewAmountCalculator(auctionCalc, feeCalc)

// Use calculator with order helper methods
takingAmount := order.CalcTakingAmountWithCalculator(calculator, makingAmount, time)
userReceiveAmount := order.GetUserReceiveAmountWithCalculator(calculator, makingAmount, time)
integratorFee := order.GetIntegratorFeeWithCalculator(calculator, time, makingAmount)
protocolFee := order.GetProtocolFeeWithCalculator(calculator, time, makingAmount)
```

## Requirements

- Go 1.21 or higher
- Valid 1inch API auth key from [1inch Portal](https://portal.1inch.dev)

## License

See LICENSE file for details.
