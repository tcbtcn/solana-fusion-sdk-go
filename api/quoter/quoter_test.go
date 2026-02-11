package quoter

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dawitel/solana-fusion-sdk-go/api"
	clienthttp "github.com/dawitel/solana-fusion-sdk-go/api/http"
	"github.com/dawitel/solana-fusion-sdk-go/domains"
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

func TestNewQuoterApi(t *testing.T) {
	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
	}
	mockClient := &mockHTTPClient{}

	api := NewQuoterApi(config, mockClient)

	if api == nil {
		t.Fatal("Expected non-nil QuoterApi")
	}
	expectedBaseURL := "https://api.example.com/quoter/v1.0/501"
	if api.baseURL != expectedBaseURL {
		t.Errorf("Expected baseURL %s, got %s", expectedBaseURL, api.baseURL)
	}
}

func TestQuoterApi_GetQuote_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("srcToken") == "" {
			t.Error("Expected srcToken parameter")
		}
		if r.URL.Query().Get("dstToken") == "" {
			t.Error("Expected dstToken parameter")
		}
		if r.URL.Query().Get("amount") == "" {
			t.Error("Expected amount parameter")
		}
		if r.URL.Query().Get("wallet") == "" {
			t.Error("Expected wallet parameter")
		}

		quote := QuoteDTO{
			QuoteID:   stringPtr("test-quote-id"),
			SrcAmount: "1000000000000000000",
			DstAmount: "1420000000",
			Presets: PresetsDTO{
				Fast: PresetDTO{
					StartAuctionIn:  10,
					AuctionDuration: 180,
					InitialRateBump: 0,
				},
			},
			RecommendedPreset: PresetTypeFast,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(quote)
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	config := api.ApiConfig{
		BaseURL: server.URL,
		Version: "v1.0",
	}
	quoterApi := NewQuoterApi(config, httpClient)

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	amount := big.NewInt(1000000000000000000)
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	quote, err := quoterApi.GetQuote(context.Background(), srcToken, dstToken, amount, signer, true, nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if quote == nil {
		t.Fatal("Expected non-nil quote")
	}
	if quote.QuoteID == nil || *quote.QuoteID != "test-quote-id" {
		t.Errorf("Expected quote ID 'test-quote-id', got %v", quote.QuoteID)
	}
	if quote.SrcAmount != "1000000000000000000" {
		t.Errorf("Expected srcAmount '1000000000000000000', got %s", quote.SrcAmount)
	}
}

func TestQuoterApi_GetQuote_WithSlippage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slippage := r.URL.Query().Get("slippage")
		if slippage == "" {
			t.Error("Expected slippage parameter")
		}
		quote := QuoteDTO{
			SrcAmount: "1000000000000000000",
			DstAmount: "1420000000",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(quote)
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	config := api.ApiConfig{
		BaseURL: server.URL,
		Version: "v1.0",
	}
	quoterApi := NewQuoterApi(config, httpClient)

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	amount := big.NewInt(1000000000000000000)
	signer := domains.MustAddressFromString("11111111111111111111111111111111")
	slippage := domains.BpsFromPercent(1.0, nil)

	_, err := quoterApi.GetQuote(context.Background(), srcToken, dstToken, amount, signer, true, slippage)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestQuoterApi_GetQuote_Error(t *testing.T) {
	mockClient := &mockHTTPClient{
		getFunc: func(ctx context.Context, url string, result interface{}) error {
			return errors.New("network error")
		},
	}

	config := api.ApiConfig{
		BaseURL: "https://api.example.com",
		Version: "v1.0",
	}
	quoterApi := NewQuoterApi(config, mockClient)

	srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
	dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
	amount := big.NewInt(1000000000000000000)
	signer := domains.MustAddressFromString("11111111111111111111111111111111")

	_, err := quoterApi.GetQuote(context.Background(), srcToken, dstToken, amount, signer, true, nil)
	if err == nil {
		t.Fatal("Expected error")
	}
}

func stringPtr(s string) *string {
	return &s
}

// TestQuoterApi_URLConstruction_NoDoubleSlash tests URL construction with various baseURL formats
func TestQuoterApi_URLConstruction_NoDoubleSlash(t *testing.T) {
	testCases := []struct {
		name    string
		baseURL string
		wantURL string
	}{
		{
			name:    "baseURL without trailing slash",
			baseURL: "https://api.example.com",
			wantURL: "https://api.example.com/quoter/v1.0/501/quote?",
		},
		{
			name:    "baseURL with trailing slash",
			baseURL: "https://api.example.com/",
			wantURL: "https://api.example.com/quoter/v1.0/501/quote?",
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
			quoterApi := NewQuoterApi(config, mockClient)

			srcToken := domains.MustAddressFromString("So11111111111111111111111111111111111111112")
			dstToken := domains.MustAddressFromString("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v")
			amount := big.NewInt(1000000000000000000)
			signer := domains.MustAddressFromString("11111111111111111111111111111111")

			_, _ = quoterApi.GetQuote(context.Background(), srcToken, dstToken, amount, signer, true, nil)

			// Check URL starts with expected base (query params may vary)
			if !strings.HasPrefix(actualURL, tc.wantURL) {
				t.Errorf("Expected URL to start with %s, got %s", tc.wantURL, actualURL)
			}
			// Verify no double slashes (excluding protocol)
			urlWithoutProtocol := strings.TrimPrefix(strings.TrimPrefix(actualURL, "http://"), "https://")
			if strings.Contains(urlWithoutProtocol, "//") {
				t.Errorf("URL contains double slash (excluding protocol): %s", actualURL)
			}
		})
	}
}
