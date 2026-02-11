package wsapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketClient implements WsProviderConnector using gorilla/websocket
type WebSocketClient struct {
	config         WsApiConfig
	conn           *websocket.Conn
	mu             sync.RWMutex
	callbacks      map[WebSocketEvent][]interface{}
	msgCallbacks   []OnMessageCb
	openCallbacks  []OnOpenCb
	closeCallbacks []OnCloseCb
	errorCallbacks []OnErrorCb
	connected      bool
	initOnce       sync.Once
}

// NewWebSocketClient creates a new WebSocket client
func NewWebSocketClient(config WsApiConfig) *WebSocketClient {
	client := &WebSocketClient{
		config:         config,
		callbacks:      make(map[WebSocketEvent][]interface{}),
		msgCallbacks:   make([]OnMessageCb, 0),
		openCallbacks:  make([]OnOpenCb, 0),
		closeCallbacks: make([]OnCloseCb, 0),
		errorCallbacks: make([]OnErrorCb, 0),
		connected:      false,
	}

	if !config.LazyInit {
		client.Init()
	}

	return client
}

// Init initializes the WebSocket connection
func (c *WebSocketClient) Init() {
	c.initOnce.Do(func() {
		url := castURL(c.config.URL)

		// Add auth key to headers if provided
		headers := make(http.Header)
		if c.config.AuthKey != "" {
			headers.Set("Authorization", "Bearer "+c.config.AuthKey)
		}

		dialer := websocket.Dialer{}
		conn, _, err := dialer.Dial(url, headers)
		if err != nil {
			c.handleError(err)
			return
		}

		c.mu.Lock()
		c.conn = conn
		c.connected = true
		c.mu.Unlock()

		// Call open callbacks
		c.mu.RLock()
		openCallbacks := make([]OnOpenCb, len(c.openCallbacks))
		copy(openCallbacks, c.openCallbacks)
		c.mu.RUnlock()

		for _, cb := range openCallbacks {
			cb()
		}

		go c.readMessages()
	})
}

func (c *WebSocketClient) readMessages() {
	for {
		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil {
			return
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			c.handleError(err)
			c.Close()
			return
		}

		var data interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			c.handleError(fmt.Errorf("failed to unmarshal message: %w", err))
			continue
		}

		c.mu.RLock()
		callbacks := make([]OnMessageCb, len(c.msgCallbacks))
		copy(callbacks, c.msgCallbacks)
		c.mu.RUnlock()

		for _, cb := range callbacks {
			cb(data)
		}
	}
}

func (c *WebSocketClient) handleError(err error) {
	c.mu.RLock()
	callbacks := make([]OnErrorCb, len(c.errorCallbacks))
	copy(callbacks, c.errorCallbacks)
	c.mu.RUnlock()

	for _, cb := range callbacks {
		cb(err)
	}
}

// Send sends a message through the WebSocket connection
func (c *WebSocketClient) Send(message interface{}) error {
	c.mu.RLock()
	conn := c.conn
	connected := c.connected
	c.mu.RUnlock()

	if !connected || conn == nil {
		return fmt.Errorf("websocket not connected")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

// On subscribes to a WebSocket event
func (c *WebSocketClient) On(event WebSocketEvent, cb interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.callbacks[event] = append(c.callbacks[event], cb)
}

// Off unsubscribes from a WebSocket event
func (c *WebSocketClient) Off(event WebSocketEvent, cb interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	callbacks := c.callbacks[event]
	for i, existingCb := range callbacks {
		// Use reflect to compare function pointers since direct comparison doesn't work
		if reflect.DeepEqual(existingCb, cb) {
			c.callbacks[event] = append(callbacks[:i], callbacks[i+1:]...)
			break
		}
	}
}

// OnOpen subscribes to the open event
func (c *WebSocketClient) OnOpen(cb OnOpenCb) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.openCallbacks = append(c.openCallbacks, cb)
}

// OnClose subscribes to the close event
func (c *WebSocketClient) OnClose(cb OnCloseCb) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closeCallbacks = append(c.closeCallbacks, cb)
}

// OnError subscribes to the error event
func (c *WebSocketClient) OnError(cb OnErrorCb) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errorCallbacks = append(c.errorCallbacks, cb)
}

// OnMessage subscribes to the message event
func (c *WebSocketClient) OnMessage(cb OnMessageCb) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.msgCallbacks = append(c.msgCallbacks, cb)
}

// Ping sends a ping message
func (c *WebSocketClient) Ping() {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn != nil {
		_ = conn.WriteMessage(websocket.PingMessage, []byte{})
	}
}

// OnPong subscribes to pong messages
func (c *WebSocketClient) OnPong(cb OnPongCb) {
	c.OnMessage(func(data interface{}) {
		// Check if it's a pong response
		if msg, ok := data.(map[string]interface{}); ok {
			if method, ok := msg["method"].(string); ok && method == string(RpcMethodPing) {
				cb()
			}
		}
	})
}

// Close closes the WebSocket connection
func (c *WebSocketClient) Close() error {
	c.mu.Lock()

	var err error
	if c.conn != nil {
		err = c.conn.Close()
		c.conn = nil
	}
	c.connected = false

	callbacks := make([]OnCloseCb, len(c.closeCallbacks))
	copy(callbacks, c.closeCallbacks)
	c.mu.Unlock()

	for _, cb := range callbacks {
		cb()
	}

	return err
}

// IsConnected returns whether the WebSocket is connected
func (c *WebSocketClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// castURL converts HTTP URL to WebSocket URL
func castURL(url string) string {
	url = strings.TrimSuffix(url, "/")
	if strings.HasPrefix(url, "http://") {
		return strings.Replace(url, "http://", "ws://", 1)
	}
	if strings.HasPrefix(url, "https://") {
		return strings.Replace(url, "https://", "wss://", 1)
	}
	if !strings.HasPrefix(url, "ws://") && !strings.HasPrefix(url, "wss://") {
		return "wss://" + url
	}
	return url
}
