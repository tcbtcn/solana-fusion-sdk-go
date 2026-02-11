package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewHTTPClient(t *testing.T) {
	baseURL := "https://api.example.com"
	authKey := "test-key"

	client := NewHTTPClient(baseURL, authKey)

	if client == nil {
		t.Fatal("Expected non-nil client")
	}
	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, client.baseURL)
	}
	if client.authKey != authKey {
		t.Errorf("Expected authKey %s, got %s", authKey, client.authKey)
	}
	if client.client == nil {
		t.Error("Expected non-nil http.Client")
	}
	if client.client.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.client.Timeout)
	}
}

func TestHTTPClient_Get_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "")
	var result map[string]string

	err := client.Get(context.Background(), "", &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %s", result["status"])
	}
}

func TestHTTPClient_Get_WithAuth(t *testing.T) {
	authKey := "test-auth-key"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "Bearer " + authKey
		if authHeader != expectedAuth {
			t.Errorf("Expected Authorization header %s, got %s", expectedAuth, authHeader)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, authKey)
	var result map[string]string

	err := client.Get(context.Background(), "", &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestHTTPClient_Get_WithNilResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "")

	err := client.Get(context.Background(), "", nil)
	if err != nil {
		t.Fatalf("Expected no error with nil result, got %v", err)
	}
}

func TestHTTPClient_Get_Non2xxStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "")

	err := client.Get(context.Background(), "", &map[string]interface{}{})
	if err == nil {
		t.Fatal("Expected error for non-2xx status")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("Expected error to contain 404, got %v", err)
	}
}

func TestHTTPClient_Get_NetworkError(t *testing.T) {
	client := NewHTTPClient("http://localhost:99999", "")

	err := client.Get(context.Background(), "/test", &map[string]interface{}{})
	if err == nil {
		t.Fatal("Expected error for network failure")
	}
}

func TestHTTPClient_Get_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "")
	var result map[string]interface{}

	err := client.Get(context.Background(), "", &result)
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "failed to decode") {
		t.Errorf("Expected decode error, got %v", err)
	}
}

func TestHTTPClient_Get_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.Get(ctx, "", &map[string]interface{}{})
	if err == nil {
		t.Fatal("Expected error for cancelled context")
	}
}

func TestHTTPClient_Post_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Expected Content-Type application/json")
		}
		var body map[string]interface{}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["test"] != "data" {
			t.Errorf("Expected test=data in body, got %v", body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "")
	var result map[string]string

	err := client.Post(context.Background(), "", map[string]string{"test": "data"}, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %s", result["status"])
	}
}

func TestHTTPClient_Post_WithAuth(t *testing.T) {
	authKey := "test-auth-key"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "Bearer " + authKey
		if authHeader != expectedAuth {
			t.Errorf("Expected Authorization header %s, got %s", expectedAuth, authHeader)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, authKey)

	err := client.Post(context.Background(), "", map[string]string{"test": "data"}, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestHTTPClient_Post_WithNilData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "")

	err := client.Post(context.Background(), "", nil, nil)
	if err != nil {
		t.Fatalf("Expected no error with nil data, got %v", err)
	}
}

func TestHTTPClient_Post_Non2xxStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Bad Request"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "")

	err := client.Post(context.Background(), "", map[string]string{"test": "data"}, &map[string]interface{}{})
	if err == nil {
		t.Fatal("Expected error for non-2xx status")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("Expected error to contain 400, got %v", err)
	}
}

func TestHTTPClient_Post_EncodeError(t *testing.T) {
	client := NewHTTPClient("http://localhost", "")

	invalidData := make(chan int)
	err := client.Post(context.Background(), "/test", invalidData, nil)
	if err == nil {
		t.Fatal("Expected error for unencodable data")
	}
	if !strings.Contains(err.Error(), "failed to encode") {
		t.Errorf("Expected encode error, got %v", err)
	}
}

func TestHTTPClient_Post_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := client.Post(ctx, "", map[string]string{"test": "data"}, nil)
	if err == nil {
		t.Fatal("Expected error for cancelled context")
	}
}
