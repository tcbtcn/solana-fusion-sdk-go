package wsapi

import (
	"fmt"
)

const (
	// Version is the WebSocket API version
	Version = "v1.0"
	// NetworkEnumSolana is the Solana network identifier
	NetworkEnumSolana = 501
)

// WebSocketApi is the main WebSocket API client
type WebSocketApi struct {
	RPC      *RpcWebsocketApi
	Order    *ActiveOrdersWebSocketApi
	provider WsProviderConnector
}

// NewWebSocketApi creates a new WebSocketApi
func NewWebSocketApi(configOrProvider interface{}) (*WebSocketApi, error) {
	var provider WsProviderConnector

	switch v := configOrProvider.(type) {
	case WsProviderConnector:
		provider = v
	case WsApiConfig:
		url := castURL(v.URL)
		// Add version and network to URL
		urlWithNetwork := fmt.Sprintf("%s/%s/%d", url, Version, NetworkEnumSolana)
		config := WsApiConfig{
			URL:      urlWithNetwork,
			AuthKey:  v.AuthKey,
			LazyInit: v.LazyInit,
		}
		provider = NewWebSocketClient(config)
	default:
		return nil, fmt.Errorf("invalid config or provider type")
	}

	return &WebSocketApi{
		RPC:      NewRpcWebsocketApi(provider),
		Order:    NewActiveOrdersWebSocketApi(provider),
		provider: provider,
	}, nil
}

// FromConfig creates a WebSocketApi from config
func FromConfig(config WsApiConfig) (*WebSocketApi, error) {
	return NewWebSocketApi(config)
}

// FromProvider creates a WebSocketApi from provider
func FromProvider(provider WsProviderConnector) (*WebSocketApi, error) {
	return NewWebSocketApi(provider)
}

// Init initializes the WebSocket connection
func (w *WebSocketApi) Init() {
	w.provider.Init()
}

// On subscribes to a WebSocket event
func (w *WebSocketApi) On(event WebSocketEvent, cb interface{}) {
	w.provider.On(event, cb)
}

// Off unsubscribes from a WebSocket event
func (w *WebSocketApi) Off(event WebSocketEvent, cb interface{}) {
	w.provider.Off(event, cb)
}

// OnOpen subscribes to the open event
func (w *WebSocketApi) OnOpen(cb OnOpenCb) {
	w.provider.OnOpen(cb)
}

// OnClose subscribes to the close event
func (w *WebSocketApi) OnClose(cb OnCloseCb) {
	w.provider.OnClose(cb)
}

// OnError subscribes to the error event
func (w *WebSocketApi) OnError(cb OnErrorCb) {
	w.provider.OnError(cb)
}

// OnMessage subscribes to the message event
func (w *WebSocketApi) OnMessage(cb OnMessageCb) {
	w.provider.OnMessage(cb)
}

// Send sends a message through the WebSocket
func (w *WebSocketApi) Send(message interface{}) error {
	return w.provider.Send(message)
}

// Close closes the WebSocket connection
func (w *WebSocketApi) Close() error {
	return w.provider.Close()
}
