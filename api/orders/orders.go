package orders

import (
	"context"
	"fmt"
	"strings"

	"github.com/tcbtcn/solana-fusion-sdk-go/api"
	"github.com/tcbtcn/solana-fusion-sdk-go/api/http"
)

// OrdersApi provides order functionality
type OrdersApi struct {
	baseURL string
	client  http.Client
}

// NewOrdersApi creates a new OrdersApi
func NewOrdersApi(config api.ApiConfig, client http.Client) *OrdersApi {
	// Remove trailing slash from BaseURL to avoid double slashes
	baseURL := strings.TrimSuffix(config.BaseURL, "/")
	baseURL = fmt.Sprintf("%s/orders/%s/501", baseURL, config.Version)

	return &OrdersApi{
		baseURL: baseURL,
		client:  client,
	}
}

// GetActiveOrders gets active orders with pagination
func (o *OrdersApi) GetActiveOrders(ctx context.Context, page, limit int) (*api.Pagination[OrderInfoDTO], error) {
	url := fmt.Sprintf("%s/order/active?limit=%d&page=%d", o.baseURL, limit, page)

	var result api.Pagination[OrderInfoDTO]
	err := o.client.Get(ctx, url, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOrderStatus gets the status of an order
func (o *OrdersApi) GetOrderStatus(ctx context.Context, orderHash string) (*OrderStatusDTO, error) {
	// Sanitize orderHash to prevent injection
	orderHash = strings.TrimSpace(orderHash)
	url := fmt.Sprintf("%s/order/status/%s", o.baseURL, orderHash)

	var result OrderStatusDTO
	err := o.client.Get(ctx, url, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOrdersCancellableByResolver gets orders cancellable by resolver with pagination
func (o *OrdersApi) GetOrdersCancellableByResolver(ctx context.Context, page, limit int) (*api.Pagination[OrderCancellableByResolverInfoDTO], error) {
	url := fmt.Sprintf("%s/order/cancelable-by-resolvers?limit=%d&page=%d", o.baseURL, limit, page)

	var result api.Pagination[OrderCancellableByResolverInfoDTO]
	err := o.client.Get(ctx, url, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
