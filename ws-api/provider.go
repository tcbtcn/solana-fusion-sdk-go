package wsapi

// WsProviderConnector is an interface for WebSocket providers
type WsProviderConnector interface {
	Init()
	Send(message interface{}) error
	On(event WebSocketEvent, cb interface{})
	Off(event WebSocketEvent, cb interface{})
	OnOpen(cb OnOpenCb)
	OnClose(cb OnCloseCb)
	OnError(cb OnErrorCb)
	OnMessage(cb OnMessageCb)
	Close() error
	IsConnected() bool
}

// WsApiConfig represents WebSocket API configuration
type WsApiConfig struct {
	URL      string
	AuthKey  string
	LazyInit bool
}
