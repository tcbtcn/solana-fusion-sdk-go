package sdk

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/api"
	"github.com/dawitel/solana-fusion-sdk-go/api/http"
	"github.com/dawitel/solana-fusion-sdk-go/api/orders"
	"github.com/dawitel/solana-fusion-sdk-go/api/quoter"
	"github.com/dawitel/solana-fusion-sdk-go/domains"
)

type mockHTTPClient struct {
	getFunc  func(ctx context.Context, url string, result interface{}) error
	postFunc func(ctx context.Context, url string, data interface{}, result interface{}) error
}

var _ http.Client = (*mockHTTPClient)(nil)

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

func TestNewSdk(t *testing.T) {
	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
		AuthKey: "test-key",
	}
	mockClient := &mockHTTPClient{}

	sdk := NewSdk(mockClient, config)

	if sdk == nil {
		t.Fatal("Expected non-nil SDK")
	}
	if sdk.ordersApi == nil {
		t.Error("Expected non-nil ordersApi")
	}
	if sdk.quoterApi == nil {
		t.Error("Expected non-nil quoterApi")
	}
}

func TestSdk_GetQuote_Success(t *testing.T) {
	quoteID := "test-quote-id"
	mockClient := &mockHTTPClient{
		getFunc: func(ctx context.Context, url string, result interface{}) error {
			if quoteDTO, ok := result.(*quoter.QuoteDTO); ok {
				quoteDTO.QuoteID = &quoteID
				quoteDTO.SrcAmount = "1000000000000000000"
				quoteDTO.DstAmount = "1420000000"
				quoteDTO.Presets = quoter.PresetsDTO{
					Fast: quoter.PresetDTO{
						StartAuctionIn:     10,
						AuctionDuration:    180,
						InitialRateBump:   0,
						AuctionStartAmount: "1500000000",
						AuctionEndAmount:   "1420000000",
						CostInDstToken:     "0",
						Points:             []quoter.PresetPointDTO{},
					},
				}
				quoteDTO.RecommendedPreset = quoter.PresetTypeFast
				quoteDTO.PriceImpactPercent = 0.5
			}
			return nil
		},
	}

	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
		AuthKey: "test-key",
	}
	sdk := NewSdk(mockClient, config)

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	amount := big.NewInt(1000000000000000000)
	signer := domains.MustAddressFromString("11111111111111111111111111111111")
	slippage := domains.BpsFromPercent(1.0, nil)

	quote, err := sdk.GetQuote(context.Background(), srcToken, dstToken, amount, signer, slippage)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if quote == nil {
		t.Fatal("Expected non-nil quote")
	}
	if quote.QuoteID != quoteID {
		t.Errorf("Expected QuoteID %s, got %s", quoteID, quote.QuoteID)
	}
}

func TestSdk_GetQuote_Error(t *testing.T) {
	mockClient := &mockHTTPClient{
		getFunc: func(ctx context.Context, url string, result interface{}) error {
			return errors.New("network error")
		},
	}

	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
		AuthKey: "test-key",
	}
	sdk := NewSdk(mockClient, config)

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	amount := big.NewInt(1000000000000000000)
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	_, err := sdk.GetQuote(context.Background(), srcToken, dstToken, amount, signer, nil)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestSdk_CreateOrder_Success(t *testing.T) {
	quoteID := "test-quote-id"
	mockClient := &mockHTTPClient{
		getFunc: func(ctx context.Context, url string, result interface{}) error {
			if quoteDTO, ok := result.(*quoter.QuoteDTO); ok {
				quoteDTO.QuoteID = &quoteID
				quoteDTO.SrcAmount = "1000000000000000000"
				quoteDTO.DstAmount = "1420000000"
				quoteDTO.Presets = quoter.PresetsDTO{
					Fast: quoter.PresetDTO{
						StartAuctionIn:     10,
						AuctionDuration:    180,
						InitialRateBump:   0,
						AuctionStartAmount: "1500000000",
						AuctionEndAmount:   "1420000000",
						CostInDstToken:     "0",
						Points:             []quoter.PresetPointDTO{},
					},
				}
				quoteDTO.RecommendedPreset = quoter.PresetTypeFast
			}
			return nil
		},
	}

	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
		AuthKey: "test-key",
	}
	sdk := NewSdk(mockClient, config)

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	amount := big.NewInt(1000000000000000000)
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	order, err := sdk.CreateOrder(context.Background(), srcToken, dstToken, amount, signer, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if order == nil {
		t.Fatal("Expected non-nil order")
	}
	if order.SrcAmount().Cmp(amount) != 0 {
		t.Errorf("Expected SrcAmount %s, got %s", amount.String(), order.SrcAmount().String())
	}
}

func TestSdk_GetOrderStatus_Success(t *testing.T) {
	orderHash := "test-hash"
	mockClient := &mockHTTPClient{
		getFunc: func(ctx context.Context, url string, result interface{}) error {
			if statusDTO, ok := result.(*orders.OrderStatusDTO); ok {
				statusDTO.Maker = "11111111111111111111111111111111"
				statusDTO.OrderHash = orderHash
				statusDTO.Status = orders.OrderStatusInProgress
				statusDTO.Order = orders.OrderDTO{
					ID:                         1,
					Receiver:                   "11111111111111111111111111111111",
					CancellationAuctionDuration: 0,
					SrcMint:                     "So11111111111111111111111111111111111111112",
					DstMint:                     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
					SrcAmount:                   "1000000000000000000",
					MinDstAmount:                "1420000000",
					EstimatedDstAmount:         "1420000000",
					ExpirationTime:              1000000000,
					SrcAssetIsNative:            false,
					DstAssetIsNative:            false,
					Fee: orders.FeeDTO{
						ProtocolFee:       0,
						IntegratorFee:     0,
						SurplusPercentage: 0,
					},
					DutchAuctionData: orders.DutchAuctionDataDTO{
						StartTime:           1000000000,
						Duration:            180,
						InitialRateBump:     0,
						PointsAndTimeDeltas: []orders.PointAndTimeDeltaDTO{},
					},
				}
				statusDTO.ApproximateTakingAmount = "1420000000"
				statusDTO.ExpirationTime = 1000000000
				statusDTO.Fills = []orders.FillDTO{}
				statusDTO.CreatedAt = 1000000000
				statusDTO.SrcTokenPriceUsd = 1.0
				statusDTO.DstTokenPriceUsd = 1.0
				statusDTO.Cancelable = true
			}
			return nil
		},
	}

	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
		AuthKey: "test-key",
	}
	sdk := NewSdk(mockClient, config)

	status, err := sdk.GetOrderStatus(context.Background(), orderHash)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if status == nil {
		t.Fatal("Expected non-nil status")
	}
	if status.OrderHash != orderHash {
		t.Errorf("Expected OrderHash %s, got %s", orderHash, status.OrderHash)
	}
}

func TestSdk_GetActiveOrders_Success(t *testing.T) {
	mockClient := &mockHTTPClient{
		getFunc: func(ctx context.Context, url string, result interface{}) error {
			if pagination, ok := result.(*api.Pagination[orders.OrderInfoDTO]); ok {
				pagination.Meta = api.PaginationMeta{
					CurrentPage:  1,
					ItemsPerPage: 10,
					TotalItems:   1,
					TotalPages:   1,
				}
				pagination.Items = []orders.OrderInfoDTO{
					{
						OrderHash:   "test-hash",
						TxSignature: "tx-sig",
						Maker:       "11111111111111111111111111111111",
						Order: orders.OrderDTO{
							ID:                         1,
							Receiver:                   "11111111111111111111111111111111",
							CancellationAuctionDuration: 0,
							SrcMint:                     "So11111111111111111111111111111111111111112",
							DstMint:                     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
							SrcAmount:                   "1000000000000000000",
							MinDstAmount:                "1420000000",
							EstimatedDstAmount:         "1420000000",
							ExpirationTime:              1000000000,
							SrcAssetIsNative:            false,
							DstAssetIsNative:            false,
							Fee: orders.FeeDTO{
								ProtocolFee:       0,
								IntegratorFee:     0,
								SurplusPercentage: 0,
							},
							DutchAuctionData: orders.DutchAuctionDataDTO{
								StartTime:           1000000000,
								Duration:            180,
								InitialRateBump:     0,
								PointsAndTimeDeltas: []orders.PointAndTimeDeltaDTO{},
							},
						},
						RemainingMakerAmount: "500000000000000000",
					},
				}
			}
			return nil
		},
	}

	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
		AuthKey: "test-key",
	}
	sdk := NewSdk(mockClient, config)

	activeOrders, err := sdk.GetActiveOrders(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(activeOrders.Items) != 1 {
		t.Errorf("Expected 1 active order, got %d", len(activeOrders.Items))
	}
}

func TestSdk_GetCancellableOrders_Success(t *testing.T) {
	mockClient := &mockHTTPClient{
		getFunc: func(ctx context.Context, url string, result interface{}) error {
			if pagination, ok := result.(*api.Pagination[orders.OrderCancellableByResolverInfoDTO]); ok {
				pagination.Meta = api.PaginationMeta{
					CurrentPage:  1,
					ItemsPerPage: 10,
					TotalItems:   1,
					TotalPages:   1,
				}
				pagination.Items = []orders.OrderCancellableByResolverInfoDTO{
					{
						OrderHash:   "test-hash",
						TxSignature: "tx-sig",
						Maker:       "11111111111111111111111111111111",
						Order: orders.OrderDTO{
							ID:                         1,
							Receiver:                   "11111111111111111111111111111111",
							CancellationAuctionDuration: 100,
							SrcMint:                     "So11111111111111111111111111111111111111112",
							DstMint:                     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
							SrcAmount:                   "1000000000000000000",
							MinDstAmount:                "1420000000",
							EstimatedDstAmount:         "1420000000",
							ExpirationTime:              1000000000,
							SrcAssetIsNative:            false,
							DstAssetIsNative:            false,
							Fee: orders.FeeDTO{
								ProtocolFee:           0,
								IntegratorFee:         0,
								SurplusPercentage:     0,
								MaxCancellationPremium: "1000000",
							},
							DutchAuctionData: orders.DutchAuctionDataDTO{
								StartTime:           1000000000,
								Duration:            180,
								InitialRateBump:     0,
								PointsAndTimeDeltas: []orders.PointAndTimeDeltaDTO{},
							},
						},
					},
				}
			}
			return nil
		},
	}

	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
		AuthKey: "test-key",
	}
	sdk := NewSdk(mockClient, config)

	cancellableOrders, err := sdk.GetOrdersCancellableByResolver(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(cancellableOrders.Items) != 1 {
		t.Errorf("Expected 1 cancellable order, got %d", len(cancellableOrders.Items))
	}
}
