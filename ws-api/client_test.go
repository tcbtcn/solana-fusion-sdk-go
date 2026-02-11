package wsapi

import (
	"errors"
	"sync"
	"testing"
)

type mockProvider struct {
	mu            sync.RWMutex
	connected     bool
	sentMessages  []interface{}
	callbacks     map[WebSocketEvent][]interface{}
	msgCallbacks  []OnMessageCb
	openCallbacks []OnOpenCb
	closeCallbacks []OnCloseCb
	errorCallbacks []OnErrorCb
}

func newMockProvider() *mockProvider {
	return &mockProvider{
		callbacks:      make(map[WebSocketEvent][]interface{}),
		msgCallbacks:   make([]OnMessageCb, 0),
		openCallbacks:  make([]OnOpenCb, 0),
		closeCallbacks: make([]OnCloseCb, 0),
		errorCallbacks: make([]OnErrorCb, 0),
		connected:      false,
	}
}

func (m *mockProvider) Init() {
	m.mu.Lock()
	m.connected = true
	callbacks := make([]OnOpenCb, len(m.openCallbacks))
	copy(callbacks, m.openCallbacks)
	m.mu.Unlock()

	for _, cb := range callbacks {
		cb()
	}
}

func (m *mockProvider) Send(message interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.connected {
		return errors.New("websocket not connected")
	}

	m.sentMessages = append(m.sentMessages, message)
	return nil
}

func (m *mockProvider) On(event WebSocketEvent, cb interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.callbacks[event] = append(m.callbacks[event], cb)
}

func (m *mockProvider) Off(event WebSocketEvent, cb interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	callbacks := m.callbacks[event]
	for i, existingCb := range callbacks {
		if existingCb == cb {
			m.callbacks[event] = append(callbacks[:i], callbacks[i+1:]...)
			break
		}
	}
}

func (m *mockProvider) OnOpen(cb OnOpenCb) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.openCallbacks = append(m.openCallbacks, cb)
}

func (m *mockProvider) OnClose(cb OnCloseCb) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closeCallbacks = append(m.closeCallbacks, cb)
}

func (m *mockProvider) OnError(cb OnErrorCb) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorCallbacks = append(m.errorCallbacks, cb)
}

func (m *mockProvider) OnMessage(cb OnMessageCb) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.msgCallbacks = append(m.msgCallbacks, cb)
}

func (m *mockProvider) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.connected = false
	callbacks := make([]OnCloseCb, len(m.closeCallbacks))
	copy(callbacks, m.closeCallbacks)

	for _, cb := range callbacks {
		cb()
	}

	return nil
}

func (m *mockProvider) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connected
}

func TestNewWebSocketClient_LazyInit(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	client := NewWebSocketClient(config)

	if client == nil {
		t.Fatal("Expected non-nil client")
	}
	if client.IsConnected() {
		t.Error("Expected client to not be connected with lazy init")
	}
}

func TestNewWebSocketClient_NonLazyInit(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://invalid-url-that-will-fail",
		AuthKey:  "test-key",
		LazyInit: false,
	}

	client := NewWebSocketClient(config)

	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

func TestWebSocketClient_Send_NotConnected(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	client := NewWebSocketClient(config)

	err := client.Send(map[string]interface{}{"test": "data"})
	if err == nil {
		t.Fatal("Expected error when not connected")
	}
}

func TestWebSocketClient_OnOff(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	client := NewWebSocketClient(config)

	callbackCalled := false
	cb := func() {
		callbackCalled = true
	}

	client.On(WebSocketEvent("test"), cb)
	client.Off(WebSocketEvent("test"), cb)

	if callbackCalled {
		t.Error("Expected callback not to be called after Off")
	}
}

func TestWebSocketClient_OnOpen(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	client := NewWebSocketClient(config)

	callbackCalled := false
	client.OnOpen(func() {
		callbackCalled = true
	})

	if callbackCalled {
		t.Error("Expected callback not to be called before Init")
	}
}

func TestWebSocketClient_OnClose(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	client := NewWebSocketClient(config)

	callbackCalled := false
	client.OnClose(func() {
		callbackCalled = true
	})

	client.Close()

	if !callbackCalled {
		t.Error("Expected close callback to be called")
	}
}

func TestWebSocketClient_OnError(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	client := NewWebSocketClient(config)

	callbackCalled := false
	client.OnError(func(err error) {
		callbackCalled = true
	})

	client.handleError(errors.New("test error"))

	if !callbackCalled {
		t.Error("Expected error callback to be called")
	}
}

func TestWebSocketClient_OnMessage(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	client := NewWebSocketClient(config)

	callbackCalled := false
	client.OnMessage(func(data interface{}) {
		callbackCalled = true
	})

	if callbackCalled {
		t.Error("Expected callback not to be called before message")
	}
}

func TestWebSocketClient_Ping(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	client := NewWebSocketClient(config)

	client.Ping()
}

func TestWebSocketClient_IsConnected(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	client := NewWebSocketClient(config)

	if client.IsConnected() {
		t.Error("Expected client to not be connected initially")
	}
}

func TestCastURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"http://example.com", "ws://example.com"},
		{"https://example.com", "wss://example.com"},
		{"ws://example.com", "ws://example.com"},
		{"wss://example.com", "wss://example.com"},
		{"example.com", "wss://example.com"},
		{"http://example.com/", "ws://example.com"},
	}

	for _, tt := range tests {
		result := castURL(tt.input)
		if result != tt.expected {
			t.Errorf("Expected %s, got %s for input %s", tt.expected, result, tt.input)
		}
	}
}
