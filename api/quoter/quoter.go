package quoter

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	"github.com/dawitel/solana-fusion-sdk-go/api"
	"github.com/dawitel/solana-fusion-sdk-go/api/http"
	"github.com/dawitel/solana-fusion-sdk-go/domains"
)

// QuoterApi provides quote functionality
type QuoterApi struct {
	baseURL string
	client  http.Client
}

// NewQuoterApi creates a new QuoterApi
func NewQuoterApi(config api.ApiConfig, client http.Client) *QuoterApi {
	baseURL := fmt.Sprintf("%s/quoter/%s/501", config.BaseURL, config.Version)

	return &QuoterApi{
		baseURL: baseURL,
		client:  client,
	}
}

// GetQuote gets a quote for a swap
func (q *QuoterApi) GetQuote(
	ctx context.Context,
	srcToken *domains.Address,
	dstToken *domains.Address,
	amount *big.Int,
	signer *domains.Address,
	enableEstimate bool,
	slippage *domains.Bps,
) (*QuoteDTO, error) {
	params := url.Values{}
	params.Set("srcToken", srcToken.ToString())
	params.Set("dstToken", dstToken.ToString())
	params.Set("amount", amount.String())
	params.Set("wallet", signer.ToString())
	params.Set("enableEstimate", fmt.Sprintf("%t", enableEstimate))

	if slippage != nil {
		params.Set("slippage", fmt.Sprintf("%f", slippage.ToPercent(nil)))
	}

	var result QuoteDTO
	path := fmt.Sprintf("/quote?%s", params.Encode())

	err := q.client.Get(ctx, q.baseURL+path, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
