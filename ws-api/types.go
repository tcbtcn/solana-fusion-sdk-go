package wsapi

import (
	"errors"

	"github.com/tcbtcn/solana-fusion-sdk-go/api"
)

// WebSocketEvent represents WebSocket connection events
type WebSocketEvent string

const (
	EventClose   WebSocketEvent = "close"
	EventError   WebSocketEvent = "error"
	EventMessage WebSocketEvent = "message"
	EventOpen    WebSocketEvent = "open"
)

// EventType represents order event types
type EventType string

const (
	EventTypeCreate EventType = "create"
	EventTypeFill   EventType = "fill"
	EventTypeCancel EventType = "cancel"
)

// RpcMethod represents RPC method names
type RpcMethod string

const (
	RpcMethodGetAllowedMethods RpcMethod = "getAllowedMethods"
	RpcMethodGetActiveOrders   RpcMethod = "getActiveOrders"
	RpcMethodPing              RpcMethod = "ping"
)

// OrderEventType represents any order event
type OrderEventType interface {
	GetEvent() EventType
}

// CreateOrderEventPayload represents the payload for order created event
type CreateOrderEventPayload struct {
	TransactionSignature     string      `json:"transactionSignature"`
	SlotNumber               uint64      `json:"slotNumber"`
	BlockTime                uint64      `json:"blockTime"`
	Action                   string      `json:"action"`
	Commitment               string      `json:"commitment"`
	OrderHash                string      `json:"orderHash"`
	Maker                    string      `json:"maker"`
	Order                    interface{} `json:"order"` // FusionOrderJSON
	FilledAuctionTakerAmount string      `json:"filledAuctionTakerAmount"`
	FilledMakerAmount        string      `json:"filledMakerAmount"`
}

// FillOrderEventPayload represents the payload for order filled event
type FillOrderEventPayload struct {
	TransactionSignature     string `json:"transactionSignature"`
	SlotNumber               uint64 `json:"slotNumber"`
	BlockTime                uint64 `json:"blockTime"`
	Action                   string `json:"action"`
	Commitment               string `json:"commitment"`
	OrderHash                string `json:"orderHash"`
	Maker                    string `json:"maker"`
	Resolver                 string `json:"resolver"`
	FilledAuctionTakerAmount string `json:"filledAuctionTakerAmount"`
	FilledMakerAmount        string `json:"filledMakerAmount"`
}

// CancelOrderEventPayload represents the payload for order cancelled event
type CancelOrderEventPayload struct {
	TransactionSignature     string `json:"transactionSignature"`
	SlotNumber               uint64 `json:"slotNumber"`
	BlockTime                uint64 `json:"blockTime"`
	Action                   string `json:"action"`
	Commitment               string `json:"commitment"`
	OrderHash                string `json:"orderHash"`
	Maker                    string `json:"maker"`
	FilledAuctionTakerAmount string `json:"filledAuctionTakerAmount"`
	FilledMakerAmount        string `json:"filledMakerAmount"`
}

// OrderCreatedEvent represents an order created event
type OrderCreatedEvent struct {
	Event  EventType               `json:"event"`
	Result CreateOrderEventPayload `json:"result"`
}

func (e *OrderCreatedEvent) GetEvent() EventType {
	return EventTypeCreate
}

// OrderFilledEvent represents an order filled event
type OrderFilledEvent struct {
	Event  EventType             `json:"event"`
	Result FillOrderEventPayload `json:"result"`
}

func (e *OrderFilledEvent) GetEvent() EventType {
	return EventTypeFill
}

// OrderCancelledEvent represents an order cancelled event
type OrderCancelledEvent struct {
	Event  EventType               `json:"event"`
	Result CancelOrderEventPayload `json:"result"`
}

func (e *OrderCancelledEvent) GetEvent() EventType {
	return EventTypeCancel
}

// RpcEventType represents any RPC event
type RpcEventType interface {
	GetMethod() RpcMethod
}

// GetAllowMethodsRpcEvent represents a getAllowedMethods RPC event
type GetAllowMethodsRpcEvent struct {
	Method RpcMethod `json:"method"`
	Result []string  `json:"result"`
}

func (e *GetAllowMethodsRpcEvent) GetMethod() RpcMethod {
	return RpcMethodGetAllowedMethods
}

// GetActiveOrdersRpcEvent represents a getActiveOrders RPC event
type GetActiveOrdersRpcEvent struct {
	Method RpcMethod                     `json:"method"`
	Result api.Pagination[ActiveOrderWS] `json:"result"`
}

func (e *GetActiveOrdersRpcEvent) GetMethod() RpcMethod {
	return RpcMethodGetActiveOrders
}

// ActiveOrderWS represents an active order in WebSocket API
type ActiveOrderWS struct {
	OrderHash            string      `json:"orderHash"`
	Order                interface{} `json:"order"` // FusionOrderJSON
	TxSignature          string      `json:"txSignature"`
	Maker                string      `json:"maker"`
	RemainingMakerAmount string      `json:"remainingMakerAmount"`
}

// PingRpcEvent represents a ping RPC event
type PingRpcEvent struct {
	Method RpcMethod `json:"method"`
	Result string    `json:"result,omitempty"`
}

func (e *PingRpcEvent) GetMethod() RpcMethod {
	return RpcMethodPing
}

// Callback types
type OnOrderCb func(data OrderEventType)
type OnOrderCreatedCb func(data *OrderCreatedEvent)
type OnOrderFilledCb func(data *OrderFilledEvent)
type OnOrderCancelledCb func(data *OrderCancelledEvent)
type OnPongCb func()

// OnGetActiveOrdersCb is a callback for getActiveOrders responses
// Note: Using interface{} for flexibility in type assertion
type OnGetActiveOrdersCb func(result interface{})
type OnGetAllowedMethodsCb func(result []string)
type OnMessageCb func(data interface{})
type OnErrorCb func(err error)
type OnCloseCb func()
type OnOpenCb func()

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// PaginationRequest represents a pagination request
type PaginationRequest struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// NewPaginationRequest creates a new pagination request with validation
func NewPaginationRequest(page, limit *int) (*PaginationRequest, error) {
	if limit != nil && (*limit < 1 || *limit > 500) {
		return nil, errors.New("limit should be in range between 1 and 500")
	}
	if page != nil && *page < 1 {
		return nil, errors.New("page should be >= 1")
	}
	return &PaginationRequest{Page: page, Limit: limit}, nil
}
