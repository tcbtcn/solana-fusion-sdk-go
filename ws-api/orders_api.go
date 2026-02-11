package wsapi

import (
	"encoding/json"
)

// ActiveOrdersWebSocketApi provides active orders event functionality
type ActiveOrdersWebSocketApi struct {
	provider WsProviderConnector
}

// NewActiveOrdersWebSocketApi creates a new ActiveOrders WebSocket API
func NewActiveOrdersWebSocketApi(provider WsProviderConnector) *ActiveOrdersWebSocketApi {
	return &ActiveOrdersWebSocketApi{
		provider: provider,
	}
}

// OnOrder subscribes to all order events
func (a *ActiveOrdersWebSocketApi) OnOrder(cb OnOrderCb) {
	a.provider.OnMessage(func(data interface{}) {
		if event := isOrderEvent(data); event != nil {
			cb(event)
		}
	})
}

// OnOrderCreated subscribes to order created events
func (a *ActiveOrdersWebSocketApi) OnOrderCreated(cb OnOrderCreatedCb) {
	a.provider.OnMessage(func(data interface{}) {
		if event := isOrderCreatedEvent(data); event != nil {
			cb(event)
		}
	})
}

// OnOrderFilled subscribes to order filled events
func (a *ActiveOrdersWebSocketApi) OnOrderFilled(cb OnOrderFilledCb) {
	a.provider.OnMessage(func(data interface{}) {
		if event := isOrderFilledEvent(data); event != nil {
			cb(event)
		}
	})
}

// OnOrderCancelled subscribes to order cancelled events
func (a *ActiveOrdersWebSocketApi) OnOrderCancelled(cb OnOrderCancelledCb) {
	a.provider.OnMessage(func(data interface{}) {
		if event := isOrderCancelledEvent(data); event != nil {
			cb(event)
		}
	})
}

// Helper functions to check event types
func isOrderEvent(data interface{}) OrderEventType {
	if msg, ok := data.(map[string]interface{}); ok {
		if eventStr, ok := msg["event"].(string); ok {
			jsonData, err := json.Marshal(msg)
			if err != nil {
				return nil
			}

			switch EventType(eventStr) {
			case EventTypeCreate:
				var event OrderCreatedEvent
				if err := json.Unmarshal(jsonData, &event); err == nil {
					return &event
				}
			case EventTypeFill:
				var event OrderFilledEvent
				if err := json.Unmarshal(jsonData, &event); err == nil {
					return &event
				}
			case EventTypeCancel:
				var event OrderCancelledEvent
				if err := json.Unmarshal(jsonData, &event); err == nil {
					return &event
				}
			}
		}
	}
	return nil
}

func isOrderCreatedEvent(data interface{}) *OrderCreatedEvent {
	if event := isOrderEvent(data); event != nil {
		if created, ok := event.(*OrderCreatedEvent); ok {
			return created
		}
	}
	return nil
}

func isOrderFilledEvent(data interface{}) *OrderFilledEvent {
	if event := isOrderEvent(data); event != nil {
		if filled, ok := event.(*OrderFilledEvent); ok {
			return filled
		}
	}
	return nil
}

func isOrderCancelledEvent(data interface{}) *OrderCancelledEvent {
	if event := isOrderEvent(data); event != nil {
		if cancelled, ok := event.(*OrderCancelledEvent); ok {
			return cancelled
		}
	}
	return nil
}
