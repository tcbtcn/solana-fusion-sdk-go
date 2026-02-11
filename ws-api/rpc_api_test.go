package wsapi

import (
	"testing"
)

func TestNewRpcWebsocketApi(t *testing.T) {
	provider := newMockProvider()

	api := NewRpcWebsocketApi(provider)

	if api == nil {
		t.Fatal("Expected non-nil RPC API")
	}
	if api.provider != provider {
		t.Error("Expected provider to match")
	}
}

func TestRpcWebsocketApi_Ping(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := NewRpcWebsocketApi(provider)

	api.Ping()

	if len(provider.sentMessages) != 1 {
		t.Errorf("Expected 1 sent message, got %d", len(provider.sentMessages))
	}
}

func TestRpcWebsocketApi_OnPong(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := NewRpcWebsocketApi(provider)

	callbackCalled := false
	api.OnPong(func() {
		callbackCalled = true
	})

	pongMessage := map[string]interface{}{
		"method": string(RpcMethodPing),
	}

	for _, cb := range provider.msgCallbacks {
		cb(pongMessage)
	}

	if !callbackCalled {
		t.Error("Expected pong callback to be called")
	}
}

func TestRpcWebsocketApi_GetActiveOrders(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := NewRpcWebsocketApi(provider)

	page := 1
	limit := 10

	err := api.GetActiveOrders(&page, &limit)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(provider.sentMessages) != 1 {
		t.Errorf("Expected 1 sent message, got %d", len(provider.sentMessages))
	}
}

func TestRpcWebsocketApi_GetActiveOrders_Error(t *testing.T) {
	provider := newMockProvider()

	api := NewRpcWebsocketApi(provider)

	page := 1
	limit := 10

	err := api.GetActiveOrders(&page, &limit)
	if err == nil {
		t.Fatal("Expected error when not connected")
	}
}

func TestRpcWebsocketApi_OnGetActiveOrders(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := NewRpcWebsocketApi(provider)

	callbackCalled := false
	api.OnGetActiveOrders(func(result interface{}) {
		callbackCalled = true
	})

	activeOrdersMessage := map[string]interface{}{
		"method": string(RpcMethodGetActiveOrders),
		"result": map[string]interface{}{},
	}

	for _, cb := range provider.msgCallbacks {
		cb(activeOrdersMessage)
	}

	if !callbackCalled {
		t.Error("Expected getActiveOrders callback to be called")
	}
}

func TestRpcWebsocketApi_GetAllowedMethods(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := NewRpcWebsocketApi(provider)

	err := api.GetAllowedMethods()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(provider.sentMessages) != 1 {
		t.Errorf("Expected 1 sent message, got %d", len(provider.sentMessages))
	}
}

func TestRpcWebsocketApi_GetAllowedMethods_Error(t *testing.T) {
	provider := newMockProvider()

	api := NewRpcWebsocketApi(provider)

	err := api.GetAllowedMethods()
	if err == nil {
		t.Fatal("Expected error when not connected")
	}
}

func TestRpcWebsocketApi_OnGetAllowedMethods(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := NewRpcWebsocketApi(provider)

	callbackCalled := false
	api.OnGetAllowedMethods(func(result []string) {
		callbackCalled = true
	})

	allowedMethodsMessage := map[string]interface{}{
		"method": string(RpcMethodGetAllowedMethods),
		"result": []string{},
	}

	for _, cb := range provider.msgCallbacks {
		cb(allowedMethodsMessage)
	}

	if !callbackCalled {
		t.Error("Expected getAllowedMethods callback to be called")
	}
}

func TestIsPingRpcEvent(t *testing.T) {
	pingMessage := map[string]interface{}{
		"method": string(RpcMethodPing),
	}

	if !isPingRpcEvent(pingMessage) {
		t.Error("Expected ping message to be recognized")
	}

	nonPingMessage := map[string]interface{}{
		"method": "other",
	}

	if isPingRpcEvent(nonPingMessage) {
		t.Error("Expected non-ping message to not be recognized")
	}
}

func TestNewPaginationRequest(t *testing.T) {
	page := 1
	limit := 10

	req, err := NewPaginationRequest(&page, &limit)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if req == nil {
		t.Fatal("Expected non-nil request")
	}
	if req.Page == nil || *req.Page != 1 {
		t.Error("Expected Page to be 1")
	}
	if req.Limit == nil || *req.Limit != 10 {
		t.Error("Expected Limit to be 10")
	}
}

func TestNewPaginationRequest_InvalidLimit(t *testing.T) {
	page := 1
	limit := 0

	_, err := NewPaginationRequest(&page, &limit)
	if err == nil {
		t.Fatal("Expected error for invalid limit")
	}

	limit = 501
	_, err = NewPaginationRequest(&page, &limit)
	if err == nil {
		t.Fatal("Expected error for limit > 500")
	}
}

func TestNewPaginationRequest_InvalidPage(t *testing.T) {
	page := 0
	limit := 10

	_, err := NewPaginationRequest(&page, &limit)
	if err == nil {
		t.Fatal("Expected error for invalid page")
	}
}

func TestNewPaginationRequest_NilValues(t *testing.T) {
	req, err := NewPaginationRequest(nil, nil)
	if err != nil {
		t.Fatalf("Expected no error for nil values, got %v", err)
	}
	if req == nil {
		t.Fatal("Expected non-nil request")
	}
}
