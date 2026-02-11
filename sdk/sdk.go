package sdk

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/dawitel/solana-fusion-sdk-go/api"
	"github.com/dawitel/solana-fusion-sdk-go/api/http"
	"github.com/dawitel/solana-fusion-sdk-go/api/orders"
	"github.com/dawitel/solana-fusion-sdk-go/api/quoter"
	"github.com/dawitel/solana-fusion-sdk-go/domains"
	fusionorder "github.com/dawitel/solana-fusion-sdk-go/fusion-order"
)

// Sdk is the main SDK client for interacting with the Solana Fusion API.
// It provides methods for getting quotes, creating orders, and querying order status.
type Sdk struct {
	ordersApi *orders.OrdersApi
	quoterApi *quoter.QuoterApi
}

// NewSdk creates a new SDK instance.
// The provider should be an HTTP client configured with the base URL and auth key.
// The apiConfig specifies the API base URL, version, and authentication key.
func NewSdk(provider http.Client, apiConfig api.ApiConfig) *Sdk {
	return &Sdk{
		ordersApi: orders.NewOrdersApi(apiConfig, provider),
		quoterApi: quoter.NewQuoterApi(apiConfig, provider),
	}
}

// GetQuote gets a quote for a token swap.
// It validates inputs (non-nil addresses, positive amount) and returns a Quote
// that can be used to create a FusionOrder.
// Parameters:
//   - ctx: Context for the request
//   - srcToken: Source token address (must not be nil)
//   - dstToken: Destination token address (must not be nil)
//   - amount: Amount to swap (must be positive, must not be nil)
//   - signer: Signer address (must not be nil)
//   - slippage: Optional slippage tolerance in basis points
func (s *Sdk) GetQuote(
	ctx context.Context,
	srcToken *domains.Address,
	dstToken *domains.Address,
	amount *big.Int,
	signer *domains.Address,
	slippage *domains.Bps,
) (*Quote, error) {
	if srcToken == nil {
		return nil, errors.New("srcToken cannot be nil")
	}
	if dstToken == nil {
		return nil, errors.New("dstToken cannot be nil")
	}
	if amount == nil {
		return nil, errors.New("amount cannot be nil")
	}
	if amount.Sign() <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if signer == nil {
		return nil, errors.New("signer cannot be nil")
	}

	quoteRaw, err := s.quoterApi.GetQuote(ctx, srcToken, dstToken, amount, signer, true, slippage)
	if err != nil {
		return nil, err
	}

	return QuoteFromJSON(srcToken, dstToken, signer, quoteRaw)
}

// CreateOrder creates a fusion order from a quote.
// It first gets a quote using GetQuote, then converts it to a FusionOrder.
// The order uses the recommended preset and the signer as the receiver.
// Parameters are the same as GetQuote.
func (s *Sdk) CreateOrder(
	ctx context.Context,
	srcToken *domains.Address,
	dstToken *domains.Address,
	amount *big.Int,
	signer *domains.Address,
	slippage *domains.Bps,
) (*fusionorder.FusionOrder, error) {
	// Validation is done in GetQuote
	quote, err := s.GetQuote(ctx, srcToken, dstToken, amount, signer, slippage)
	if err != nil {
		return nil, err
	}

	return quote.ToOrder(quote.RecommendedPreset, nil)
}

// GetOrderStatus gets the status of an order by its hash.
// Returns an error if orderHash is empty or whitespace-only.
// Parameters:
//   - ctx: Context for the request
//   - orderHash: Base58-encoded order hash (must not be empty)
func (s *Sdk) GetOrderStatus(ctx context.Context, orderHash string) (*OrderStatus, error) {
	orderHash = strings.TrimSpace(orderHash)
	if orderHash == "" {
		return nil, errors.New("orderHash cannot be empty")
	}

	orderStatus, err := s.ordersApi.GetOrderStatus(ctx, orderHash)
	if err != nil {
		return nil, err
	}

	return OrderStatusFromJSON(orderStatus)
}

// GetOrdersCancellableByResolver gets orders that can be cancelled by a resolver.
// Returns a paginated list of cancellable orders.
// Parameters:
//   - ctx: Context for the request
//   - page: Page number (must be >= 1)
//   - limit: Number of items per page (must be >= 1)
func (s *Sdk) GetOrdersCancellableByResolver(ctx context.Context, page, limit int) (*api.Pagination[*CancellableOrder], error) {
	if page < 1 {
		return nil, errors.New("page must be >= 1")
	}
	if limit < 1 {
		return nil, errors.New("limit must be >= 1")
	}

	res, err := s.ordersApi.GetOrdersCancellableByResolver(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	items := make([]*CancellableOrder, len(res.Items))
	for i, o := range res.Items {
		maker, err := domains.NewAddress(o.Maker)
		if err != nil {
			return nil, err
		}

		// Convert order DTO to FusionOrder
		orderJSON := &fusionorder.FusionOrderJSON{
			ID:                          o.Order.ID,
			Receiver:                    o.Order.Receiver,
			CancellationAuctionDuration: o.Order.CancellationAuctionDuration,
			SrcMint:                     o.Order.SrcMint,
			DstMint:                     o.Order.DstMint,
			SrcAmount:                   o.Order.SrcAmount,
			MinDstAmount:                o.Order.MinDstAmount,
			EstimatedDstAmount:          o.Order.EstimatedDstAmount,
			ExpirationTime:              o.Order.ExpirationTime,
			SrcAssetIsNative:            o.Order.SrcAssetIsNative,
			DstAssetIsNative:            o.Order.DstAssetIsNative,
			Fee: fusionorder.FeeJSON{
				ProtocolFee:            o.Order.Fee.ProtocolFee,
				IntegratorFee:          o.Order.Fee.IntegratorFee,
				SurplusPercentage:      o.Order.Fee.SurplusPercentage,
				MaxCancellationPremium: o.Order.Fee.MaxCancellationPremium,
			},
			DutchAuctionData: fusionorder.DutchAuctionDataJSON{
				StartTime:           o.Order.DutchAuctionData.StartTime,
				Duration:            o.Order.DutchAuctionData.Duration,
				InitialRateBump:     o.Order.DutchAuctionData.InitialRateBump,
				PointsAndTimeDeltas: make([]fusionorder.PointAndTimeDeltaJSON, len(o.Order.DutchAuctionData.PointsAndTimeDeltas)),
			},
		}

		for j, p := range o.Order.DutchAuctionData.PointsAndTimeDeltas {
			orderJSON.DutchAuctionData.PointsAndTimeDeltas[j] = fusionorder.PointAndTimeDeltaJSON{
				RateBump:  p.RateBump,
				TimeDelta: p.TimeDelta,
			}
		}

		if o.Order.Fee.ProtocolDstAta != "" {
			orderJSON.Fee.ProtocolDstAta = &o.Order.Fee.ProtocolDstAta
		}
		if o.Order.Fee.IntegratorDstAta != "" {
			orderJSON.Fee.IntegratorDstAta = &o.Order.Fee.IntegratorDstAta
		}

		order, err := fusionorder.FromJSON(orderJSON)
		if err != nil {
			return nil, err
		}

		items[i] = &CancellableOrder{
			Maker: maker,
			Order: order,
		}
	}

	return &api.Pagination[*CancellableOrder]{
		Meta:  res.Meta,
		Items: items,
	}, nil
}

// GetActiveOrders gets active orders with pagination.
// Returns a paginated list of active orders.
// Parameters:
//   - ctx: Context for the request
//   - page: Page number (must be >= 1)
//   - limit: Number of items per page (must be >= 1)
func (s *Sdk) GetActiveOrders(ctx context.Context, page, limit int) (*api.Pagination[*ActiveOrder], error) {
	if page < 1 {
		return nil, errors.New("page must be >= 1")
	}
	if limit < 1 {
		return nil, errors.New("limit must be >= 1")
	}

	res, err := s.ordersApi.GetActiveOrders(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	items := make([]*ActiveOrder, len(res.Items))
	for i, o := range res.Items {
		activeOrder, err := ActiveOrderFromJSON(&o)
		if err != nil {
			return nil, err
		}
		items[i] = activeOrder
	}

	return &api.Pagination[*ActiveOrder]{
		Meta:  res.Meta,
		Items: items,
	}, nil
}
