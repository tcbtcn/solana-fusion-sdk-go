package orders

import (
	"context"
	"fmt"

	"github.com/dawitel/solana-fusion-sdk-go/api"
	"github.com/dawitel/solana-fusion-sdk-go/api/http"
)

// OrdersApi provides order functionality
type OrdersApi struct {
	baseURL string
	client  http.Client
}

// NewOrdersApi creates a new OrdersApi
func NewOrdersApi(config api.ApiConfig, client http.Client) *OrdersApi {
	baseURL := fmt.Sprintf("%s/orders/%s/501", config.BaseURL, config.Version)

	return &OrdersApi{
		baseURL: baseURL,
		client:  client,
	}
}

// GetActiveOrders gets active orders with pagination
func (o *OrdersApi) GetActiveOrders(ctx context.Context, page, limit int) (*api.Pagination[OrderInfoDTO], error) {
	path := fmt.Sprintf("/order/active?limit=%d&page=%d", limit, page)

	var result api.Pagination[OrderInfoDTO]
	err := o.client.Get(ctx, o.baseURL+path, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOrderStatus gets the status of an order
func (o *OrdersApi) GetOrderStatus(ctx context.Context, orderHash string) (*OrderStatusDTO, error) {
	path := fmt.Sprintf("/order/status/%s", orderHash)

	var result OrderStatusDTO
	err := o.client.Get(ctx, o.baseURL+path, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOrdersCancellableByResolver gets orders cancellable by resolver with pagination
func (o *OrdersApi) GetOrdersCancellableByResolver(ctx context.Context, page, limit int) (*api.Pagination[OrderCancellableByResolverInfoDTO], error) {
	path := fmt.Sprintf("/order/cancelable-by-resolvers?limit=%d&page=%d", limit, page)

	var result api.Pagination[OrderCancellableByResolverInfoDTO]
	err := o.client.Get(ctx, o.baseURL+path, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
