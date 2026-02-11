package wsapi

import (
	"testing"
)

func TestWsProviderConnector_Interface(t *testing.T) {
	var _ WsProviderConnector = newMockProvider()
}

func TestWsApiConfig(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.example.com",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	if config.URL != "wss://api.example.com" {
		t.Errorf("Expected URL %s, got %s", "wss://api.example.com", config.URL)
	}
	if config.AuthKey != "test-key" {
		t.Errorf("Expected AuthKey %s, got %s", "test-key", config.AuthKey)
	}
	if !config.LazyInit {
		t.Error("Expected LazyInit to be true")
	}
}
