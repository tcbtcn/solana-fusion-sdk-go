package wsapi

import (
	"testing"
)

func TestNewActiveOrdersWebSocketApi(t *testing.T) {
	provider := newMockProvider()

	api := NewActiveOrdersWebSocketApi(provider)

	if api == nil {
		t.Fatal("Expected non-nil ActiveOrders API")
	}
	if api.provider != provider {
		t.Error("Expected provider to match")
	}
}

func TestActiveOrdersWebSocketApi_OnOrder(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := NewActiveOrdersWebSocketApi(provider)

	callbackCalled := false
	api.OnOrder(func(event OrderEventType) {
		callbackCalled = true
	})

	orderCreatedMessage := map[string]interface{}{
		"event": string(EventTypeCreate),
		"data":  map[string]interface{}{},
	}

	for _, cb := range provider.msgCallbacks {
		cb(orderCreatedMessage)
	}

	if !callbackCalled {
		t.Error("Expected order callback to be called")
	}
}

func TestActiveOrdersWebSocketApi_OnOrderCreated(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := NewActiveOrdersWebSocketApi(provider)

	callbackCalled := false
	api.OnOrderCreated(func(event *OrderCreatedEvent) {
		callbackCalled = true
	})

	orderCreatedMessage := map[string]interface{}{
		"event": string(EventTypeCreate),
		"data":  map[string]interface{}{},
	}

	for _, cb := range provider.msgCallbacks {
		cb(orderCreatedMessage)
	}

	if !callbackCalled {
		t.Error("Expected order created callback to be called")
	}
}

func TestActiveOrdersWebSocketApi_OnOrderFilled(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := NewActiveOrdersWebSocketApi(provider)

	callbackCalled := false
	api.OnOrderFilled(func(event *OrderFilledEvent) {
		callbackCalled = true
	})

	orderFilledMessage := map[string]interface{}{
		"event": string(EventTypeFill),
		"data":  map[string]interface{}{},
	}

	for _, cb := range provider.msgCallbacks {
		cb(orderFilledMessage)
	}

	if !callbackCalled {
		t.Error("Expected order filled callback to be called")
	}
}

func TestActiveOrdersWebSocketApi_OnOrderCancelled(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := NewActiveOrdersWebSocketApi(provider)

	callbackCalled := false
	api.OnOrderCancelled(func(event *OrderCancelledEvent) {
		callbackCalled = true
	})

	orderCancelledMessage := map[string]interface{}{
		"event": string(EventTypeCancel),
		"data":  map[string]interface{}{},
	}

	for _, cb := range provider.msgCallbacks {
		cb(orderCancelledMessage)
	}

	if !callbackCalled {
		t.Error("Expected order cancelled callback to be called")
	}
}

func TestIsOrderEvent(t *testing.T) {
	orderCreatedMessage := map[string]interface{}{
		"event": string(EventTypeCreate),
		"data":  map[string]interface{}{},
	}

	event := isOrderEvent(orderCreatedMessage)
	if event == nil {
		t.Error("Expected order event to be recognized")
	}

	nonOrderMessage := map[string]interface{}{
		"method": "other",
	}

	event = isOrderEvent(nonOrderMessage)
	if event != nil {
		t.Error("Expected non-order message to not be recognized")
	}
}

func TestIsOrderCreatedEvent(t *testing.T) {
	orderCreatedMessage := map[string]interface{}{
		"event": string(EventTypeCreate),
		"data":  map[string]interface{}{},
	}

	event := isOrderCreatedEvent(orderCreatedMessage)
	if event == nil {
		t.Error("Expected order created event to be recognized")
	}
}

func TestIsOrderFilledEvent(t *testing.T) {
	orderFilledMessage := map[string]interface{}{
		"event": string(EventTypeFill),
		"data":  map[string]interface{}{},
	}

	event := isOrderFilledEvent(orderFilledMessage)
	if event == nil {
		t.Error("Expected order filled event to be recognized")
	}
}

func TestIsOrderCancelledEvent(t *testing.T) {
	orderCancelledMessage := map[string]interface{}{
		"event": string(EventTypeCancel),
		"data":  map[string]interface{}{},
	}

	event := isOrderCancelledEvent(orderCancelledMessage)
	if event == nil {
		t.Error("Expected order cancelled event to be recognized")
	}
}
