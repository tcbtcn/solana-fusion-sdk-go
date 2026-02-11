package orders

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/api"
	clienthttp "github.com/dawitel/solana-fusion-sdk-go/api/http"
)

type mockHTTPClient struct {
	getFunc  func(ctx context.Context, url string, result interface{}) error
	postFunc func(ctx context.Context, url string, data interface{}, result interface{}) error
}

func (m *mockHTTPClient) Get(ctx context.Context, url string, result interface{}) error {
	if m.getFunc != nil {
		return m.getFunc(ctx, url, result)
	}
	return errors.New("not implemented")
}

func (m *mockHTTPClient) Post(ctx context.Context, url string, data interface{}, result interface{}) error {
	if m.postFunc != nil {
		return m.postFunc(ctx, url, data, result)
	}
	return errors.New("not implemented")
}

func TestNewOrdersApi(t *testing.T) {
	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
	}
	mockClient := &mockHTTPClient{}

	api := NewOrdersApi(config, mockClient)

	if api == nil {
		t.Fatal("Expected non-nil OrdersApi")
	}
	expectedBaseURL := "https://api.example.com/orders/v1.0/501"
	if api.baseURL != expectedBaseURL {
		t.Errorf("Expected baseURL %s, got %s", expectedBaseURL, api.baseURL)
	}
}

func TestOrdersApi_GetOrderStatus_Success(t *testing.T) {
	orderHash := "test-order-hash"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// URL is constructed as: baseURL/order/status/{hash}
		// where baseURL = server.URL/orders/v1.0/501
		// So full path = /orders/v1.0/501/order/status/{hash}
		expectedPath := "/orders/v1.0/501/order/status/" + orderHash
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		status := OrderStatusDTO{
			Maker:     "11111111111111111111111111111111",
			OrderHash: orderHash,
			Status:    OrderStatusInProgress,
			Order: OrderDTO{
				ID:           1,
				SrcAmount:    "1000000000000000000",
				MinDstAmount: "1420000000",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(status)
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	config := api.ApiConfig{
		BaseURL: server.URL,
		Version: "v1.0",
	}
	ordersApi := NewOrdersApi(config, httpClient)

	status, err := ordersApi.GetOrderStatus(context.Background(), orderHash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if status == nil {
		t.Fatal("Expected non-nil status")
	}
	if status.OrderHash != orderHash {
		t.Errorf("Expected orderHash %s, got %s", orderHash, status.OrderHash)
	}
	if status.Status != OrderStatusInProgress {
		t.Errorf("Expected status %d, got %d", OrderStatusInProgress, status.Status)
	}
}

func TestOrdersApi_GetOrderStatus_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Order not found"))
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	config := api.ApiConfig{
		BaseURL: server.URL,
		Version: "v1.0",
	}
	ordersApi := NewOrdersApi(config, httpClient)

	_, err := ordersApi.GetOrderStatus(context.Background(), "non-existent")
	if err == nil {
		t.Fatal("Expected error for not found order")
	}
}

func TestOrdersApi_GetOrderStatus_Error(t *testing.T) {
	mockClient := &mockHTTPClient{
		getFunc: func(ctx context.Context, url string, result interface{}) error {
			return errors.New("network error")
		},
	}

	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
	}
	ordersApi := NewOrdersApi(config, mockClient)

	_, err := ordersApi.GetOrderStatus(context.Background(), "test-hash")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestOrdersApi_GetActiveOrders_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit := r.URL.Query().Get("limit")
		page := r.URL.Query().Get("page")
		if limit == "" || page == "" {
			t.Error("Expected limit and page parameters")
		}

		result := api.Pagination[OrderInfoDTO]{
			Meta: api.PaginationMeta{
				CurrentPage:  1,
				ItemsPerPage: 10,
				TotalItems:   1,
				TotalPages:   1,
			},
			Items: []OrderInfoDTO{
				{
					OrderHash: "test-hash",
					Maker:     "11111111111111111111111111111111",
					Order: OrderDTO{
						ID:           1,
						SrcAmount:    "1000000000000000000",
						MinDstAmount: "1420000000",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	config := api.ApiConfig{
		BaseURL: server.URL,
		Version: "v1.0",
	}
	ordersApi := NewOrdersApi(config, httpClient)

	result, err := ordersApi.GetActiveOrders(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if len(result.Items) != 1 {
		t.Errorf("Expected 1 order, got %d", len(result.Items))
	}
	if result.Items[0].OrderHash != "test-hash" {
		t.Errorf("Expected orderHash 'test-hash', got %s", result.Items[0].OrderHash)
	}
}

func TestOrdersApi_GetActiveOrders_Empty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := api.Pagination[OrderInfoDTO]{
			Meta: api.PaginationMeta{
				CurrentPage:  1,
				ItemsPerPage: 10,
				TotalItems:   0,
				TotalPages:   0,
			},
			Items: []OrderInfoDTO{},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	config := api.ApiConfig{
		BaseURL: server.URL,
		Version: "v1.0",
	}
	ordersApi := NewOrdersApi(config, httpClient)

	result, err := ordersApi.GetActiveOrders(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(result.Items) != 0 {
		t.Errorf("Expected 0 orders, got %d", len(result.Items))
	}
}

func TestOrdersApi_GetActiveOrders_Error(t *testing.T) {
	mockClient := &mockHTTPClient{
		getFunc: func(ctx context.Context, url string, result interface{}) error {
			return errors.New("network error")
		},
	}

	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
	}
	ordersApi := NewOrdersApi(config, mockClient)

	_, err := ordersApi.GetActiveOrders(context.Background(), 1, 10)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestOrdersApi_GetOrdersCancellableByResolver_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit := r.URL.Query().Get("limit")
		page := r.URL.Query().Get("page")
		if limit == "" || page == "" {
			t.Error("Expected limit and page parameters")
		}

		result := api.Pagination[OrderCancellableByResolverInfoDTO]{
			Meta: api.PaginationMeta{
				CurrentPage:  1,
				ItemsPerPage: 10,
				TotalItems:   1,
				TotalPages:   1,
			},
			Items: []OrderCancellableByResolverInfoDTO{
				{
					OrderHash: "test-hash",
					Maker:     "11111111111111111111111111111111",
					Order: OrderDTO{
						ID:           1,
						SrcAmount:    "1000000000000000000",
						MinDstAmount: "1420000000",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(result)
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	config := api.ApiConfig{
		BaseURL: server.URL,
		Version: "v1.0",
	}
	ordersApi := NewOrdersApi(config, httpClient)

	result, err := ordersApi.GetOrdersCancellableByResolver(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil result")
	}
	if len(result.Items) != 1 {
		t.Errorf("Expected 1 order, got %d", len(result.Items))
	}
}

func TestOrdersApi_GetOrdersCancellableByResolver_Error(t *testing.T) {
	mockClient := &mockHTTPClient{
		getFunc: func(ctx context.Context, url string, result interface{}) error {
			return errors.New("network error")
		},
	}

	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
	}
	ordersApi := NewOrdersApi(config, mockClient)

	_, err := ordersApi.GetOrdersCancellableByResolver(context.Background(), 1, 10)
	if err == nil {
		t.Fatal("Expected error")
	}
}

// TestOrdersApi_URLConstruction_NoDoubleSlash tests URL construction with various baseURL formats
func TestOrdersApi_URLConstruction_NoDoubleSlash(t *testing.T) {
	testCases := []struct {
		name    string
		baseURL string
		wantURL string
	}{
		{
			name:    "baseURL without trailing slash",
			baseURL: "https://api.example.com",
			wantURL: "https://api.example.com/orders/v1.0/501/order/active?limit=10&page=1",
		},
		{
			name:    "baseURL with trailing slash",
			baseURL: "https://api.example.com/",
			wantURL: "https://api.example.com/orders/v1.0/501/order/active?limit=10&page=1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actualURL string
			mockClient := &mockHTTPClient{
				getFunc: func(ctx context.Context, url string, result interface{}) error {
					actualURL = url
					return nil
				},
			}

			config := api.ApiConfig{
				BaseURL: tc.baseURL,
				Version: "v1.0",
			}
			ordersApi := NewOrdersApi(config, mockClient)

			_, _ = ordersApi.GetActiveOrders(context.Background(), 1, 10)

			if actualURL != tc.wantURL {
				t.Errorf("Expected URL %s, got %s", tc.wantURL, actualURL)
			}
			// Verify no double slashes (excluding protocol)
			urlWithoutProtocol := strings.TrimPrefix(strings.TrimPrefix(actualURL, "http://"), "https://")
			if strings.Contains(urlWithoutProtocol, "//") {
				t.Errorf("URL contains double slash (excluding protocol): %s", actualURL)
			}
		})
	}
}
