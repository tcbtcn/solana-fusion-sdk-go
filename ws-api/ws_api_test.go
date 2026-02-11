package wsapi

import (
	"testing"
)

func TestNewWebSocketApi_FromConfig(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	api, err := NewWebSocketApi(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if api == nil {
		t.Fatal("Expected non-nil WebSocketApi")
	}
	if api.RPC == nil {
		t.Error("Expected non-nil RPC API")
	}
	if api.Order == nil {
		t.Error("Expected non-nil Order API")
	}
	if api.provider == nil {
		t.Error("Expected non-nil provider")
	}
}

func TestNewWebSocketApi_FromProvider(t *testing.T) {
	provider := newMockProvider()

	api, err := NewWebSocketApi(provider)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if api == nil {
		t.Fatal("Expected non-nil WebSocketApi")
	}
	if api.provider != provider {
		t.Error("Expected provider to match")
	}
}

func TestNewWebSocketApi_InvalidType(t *testing.T) {
	_, err := NewWebSocketApi("invalid")
	if err == nil {
		t.Fatal("Expected error for invalid type")
	}
}

func TestFromConfig(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	api, err := FromConfig(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if api == nil {
		t.Fatal("Expected non-nil WebSocketApi")
	}
}

func TestFromProvider(t *testing.T) {
	provider := newMockProvider()

	api, err := FromProvider(provider)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if api == nil {
		t.Fatal("Expected non-nil WebSocketApi")
	}
}

func TestWebSocketApi_Init(t *testing.T) {
	provider := newMockProvider()
	api := &WebSocketApi{
		provider: provider,
	}

	api.Init()

	if !provider.IsConnected() {
		t.Error("Expected provider to be connected after Init")
	}
}

func TestWebSocketApi_OnOff(t *testing.T) {
	provider := newMockProvider()
	api := &WebSocketApi{
		provider: provider,
	}

	callbackCalled := false
	cb := func() {
		callbackCalled = true
	}

	api.On(WebSocketEvent("test"), cb)
	api.Off(WebSocketEvent("test"), cb)

	if callbackCalled {
		t.Error("Expected callback not to be called after Off")
	}
}

func TestWebSocketApi_OnOpen(t *testing.T) {
	provider := newMockProvider()
	api := &WebSocketApi{
		provider: provider,
	}

	callbackCalled := false
	api.OnOpen(func() {
		callbackCalled = true
	})

	provider.Init()

	if !callbackCalled {
		t.Error("Expected open callback to be called")
	}
}

func TestWebSocketApi_OnClose(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := &WebSocketApi{
		provider: provider,
	}

	callbackCalled := false
	api.OnClose(func() {
		callbackCalled = true
	})

	api.Close()

	if !callbackCalled {
		t.Error("Expected close callback to be called")
	}
}

func TestWebSocketApi_OnError(t *testing.T) {
	provider := newMockProvider()
	api := &WebSocketApi{
		provider: provider,
	}

	callbackCalled := false
	api.OnError(func(err error) {
		callbackCalled = true
	})

	provider.errorCallbacks[0](nil)

	if !callbackCalled {
		t.Error("Expected error callback to be called")
	}
}

func TestWebSocketApi_OnMessage(t *testing.T) {
	provider := newMockProvider()
	api := &WebSocketApi{
		provider: provider,
	}

	callbackCalled := false
	api.OnMessage(func(data interface{}) {
		callbackCalled = true
	})

	for _, cb := range provider.msgCallbacks {
		cb(map[string]interface{}{"test": "data"})
	}

	if !callbackCalled {
		t.Error("Expected message callback to be called")
	}
}

func TestWebSocketApi_Send(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := &WebSocketApi{
		provider: provider,
	}

	err := api.Send(map[string]interface{}{"test": "data"})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(provider.sentMessages) != 1 {
		t.Errorf("Expected 1 sent message, got %d", len(provider.sentMessages))
	}
}

func TestWebSocketApi_Close(t *testing.T) {
	provider := newMockProvider()
	provider.Init()

	api := &WebSocketApi{
		provider: provider,
	}

	err := api.Close()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if provider.IsConnected() {
		t.Error("Expected provider to be disconnected after Close")
	}
}
