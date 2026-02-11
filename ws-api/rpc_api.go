package wsapi

import (
	"encoding/json"
)

// RpcWebsocketApi provides RPC functionality over WebSocket
type RpcWebsocketApi struct {
	provider WsProviderConnector
}

// NewRpcWebsocketApi creates a new RPC WebSocket API
func NewRpcWebsocketApi(provider WsProviderConnector) *RpcWebsocketApi {
	return &RpcWebsocketApi{
		provider: provider,
	}
}

// Ping sends a ping message
func (r *RpcWebsocketApi) Ping() {
	_ = r.provider.Send(map[string]interface{}{
		"method": string(RpcMethodPing),
	})
}

// OnPong subscribes to pong responses
func (r *RpcWebsocketApi) OnPong(cb OnPongCb) {
	r.provider.OnMessage(func(data interface{}) {
		if isPingRpcEvent(data) {
			cb()
		}
	})
}

// GetActiveOrders requests active orders with pagination
func (r *RpcWebsocketApi) GetActiveOrders(page, limit *int) error {
	params, err := NewPaginationRequest(page, limit)
	if err != nil {
		return err
	}

	return r.provider.Send(map[string]interface{}{
		"method": string(RpcMethodGetActiveOrders),
		"param":  params,
	})
}

// OnGetActiveOrders subscribes to getActiveOrders responses
func (r *RpcWebsocketApi) OnGetActiveOrders(cb OnGetActiveOrdersCb) {
	r.provider.OnMessage(func(data interface{}) {
		if event := isGetActiveOrdersRpcEvent(data); event != nil {
			cb(event.Result)
		}
	})
}

// GetAllowedMethods requests allowed methods
func (r *RpcWebsocketApi) GetAllowedMethods() error {
	return r.provider.Send(map[string]interface{}{
		"method": string(RpcMethodGetAllowedMethods),
	})
}

// OnGetAllowedMethods subscribes to getAllowedMethods responses
func (r *RpcWebsocketApi) OnGetAllowedMethods(cb OnGetAllowedMethodsCb) {
	r.provider.OnMessage(func(data interface{}) {
		if event := isGetAllowedMethodsRpcEvent(data); event != nil {
			cb(event.Result)
		}
	})
}

// Helper functions to check event types
func isPingRpcEvent(data interface{}) bool {
	if msg, ok := data.(map[string]interface{}); ok {
		if method, ok := msg["method"].(string); ok {
			return method == string(RpcMethodPing)
		}
	}
	return false
}

func isGetActiveOrdersRpcEvent(data interface{}) *GetActiveOrdersRpcEvent {
	if msg, ok := data.(map[string]interface{}); ok {
		if method, ok := msg["method"].(string); ok && method == string(RpcMethodGetActiveOrders) {
			// Try to unmarshal into GetActiveOrdersRpcEvent
			jsonData, err := json.Marshal(msg)
			if err != nil {
				return nil
			}
			var event GetActiveOrdersRpcEvent
			if err := json.Unmarshal(jsonData, &event); err == nil {
				return &event
			}
		}
	}
	return nil
}

func isGetAllowedMethodsRpcEvent(data interface{}) *GetAllowMethodsRpcEvent {
	if msg, ok := data.(map[string]interface{}); ok {
		if method, ok := msg["method"].(string); ok && method == string(RpcMethodGetAllowedMethods) {
			// Try to unmarshal into GetAllowMethodsRpcEvent
			jsonData, err := json.Marshal(msg)
			if err != nil {
				return nil
			}
			var event GetAllowMethodsRpcEvent
			if err := json.Unmarshal(jsonData, &event); err == nil {
				return &event
			}
		}
	}
	return nil
}
