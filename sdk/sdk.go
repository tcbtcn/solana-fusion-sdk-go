package sdk

import (
	"context"
	"math/big"

	"github.com/dawitel/solana-fusion-sdk-go/api"
	"github.com/dawitel/solana-fusion-sdk-go/api/http"
	"github.com/dawitel/solana-fusion-sdk-go/api/orders"
	"github.com/dawitel/solana-fusion-sdk-go/api/quoter"
	"github.com/dawitel/solana-fusion-sdk-go/domains"
	fusionorder "github.com/dawitel/solana-fusion-sdk-go/fusion-order"
)

// Sdk is the main SDK client
type Sdk struct {
	ordersApi *orders.OrdersApi
	quoterApi *quoter.QuoterApi
}

// NewSdk creates a new SDK instance
// The provider should be an HTTP client configured with the base URL and auth key
func NewSdk(provider http.Client, apiConfig api.ApiConfig) *Sdk {
	return &Sdk{
		ordersApi: orders.NewOrdersApi(apiConfig, provider),
		quoterApi: quoter.NewQuoterApi(apiConfig, provider),
	}
}

// GetQuote gets a quote for a swap
func (s *Sdk) GetQuote(
	ctx context.Context,
	srcToken *domains.Address,
	dstToken *domains.Address,
	amount *big.Int,
	signer *domains.Address,
	slippage *domains.Bps,
) (*Quote, error) {
	quoteRaw, err := s.quoterApi.GetQuote(ctx, srcToken, dstToken, amount, signer, true, slippage)
	if err != nil {
		return nil, err
	}

	return QuoteFromJSON(srcToken, dstToken, signer, quoteRaw)
}

// CreateOrder creates a fusion order from a quote
func (s *Sdk) CreateOrder(
	ctx context.Context,
	srcToken *domains.Address,
	dstToken *domains.Address,
	amount *big.Int,
	signer *domains.Address,
	slippage *domains.Bps,
) (*fusionorder.FusionOrder, error) {
	quote, err := s.GetQuote(ctx, srcToken, dstToken, amount, signer, slippage)
	if err != nil {
		return nil, err
	}

	return quote.ToOrder(quote.RecommendedPreset, nil)
}

// GetOrderStatus gets the status of an order
func (s *Sdk) GetOrderStatus(ctx context.Context, orderHash string) (*OrderStatus, error) {
	orderStatus, err := s.ordersApi.GetOrderStatus(ctx, orderHash)
	if err != nil {
		return nil, err
	}

	return OrderStatusFromJSON(orderStatus)
}

// GetOrdersCancellableByResolver gets orders cancellable by resolver
func (s *Sdk) GetOrdersCancellableByResolver(ctx context.Context, page, limit int) (*api.Pagination[*CancellableOrder], error) {
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

// GetActiveOrders gets active orders
func (s *Sdk) GetActiveOrders(ctx context.Context, page, limit int) (*api.Pagination[*ActiveOrder], error) {
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
