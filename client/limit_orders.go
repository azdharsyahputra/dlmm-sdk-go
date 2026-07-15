package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GetLimitOrdersParams represents parameters for pagination in limit orders.
type GetLimitOrdersParams struct {
	Page     *int
	PageSize *int
}

func addPaginationParams(q url.Values, params *GetLimitOrdersParams) {
	if params != nil {
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
	}
}

// GetClosedLimitOrderPools retrieves closed limit order pools summary for a wallet.
func (c *Client) GetClosedLimitOrderPools(ctx context.Context, wallet string, params *GetLimitOrdersParams) (*PaginationResponse_ClosedLimitOrderPoolSummary, error) {
	if wallet == "" {
		return nil, fmt.Errorf("wallet address cannot be empty")
	}

	q := url.Values{}
	addPaginationParams(q, params)

	path := fmt.Sprintf("/wallets/%s/limit_orders/closed/pools", wallet)
	var resp PaginationResponse_ClosedLimitOrderPoolSummary
	if err := c.doRequest(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetClosedLimitOrdersForPool retrieves closed limit orders details for a wallet in a pool.
func (c *Client) GetClosedLimitOrdersForPool(ctx context.Context, wallet string, poolAddress string, params *GetLimitOrdersParams) (*ClosedLimitOrdersForPoolResponse, error) {
	if wallet == "" || poolAddress == "" {
		return nil, fmt.Errorf("wallet address and pool address cannot be empty")
	}

	q := url.Values{}
	addPaginationParams(q, params)

	path := fmt.Sprintf("/wallets/%s/limit_orders/closed/pools/%s", wallet, poolAddress)
	var resp ClosedLimitOrdersForPoolResponse
	if err := c.doRequest(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOpenLimitOrderPools retrieves open limit order pools summary for a wallet.
func (c *Client) GetOpenLimitOrderPools(ctx context.Context, wallet string, params *GetLimitOrdersParams) (*PaginationResponse_OpenLimitOrderPoolSummary, error) {
	if wallet == "" {
		return nil, fmt.Errorf("wallet address cannot be empty")
	}

	q := url.Values{}
	addPaginationParams(q, params)

	path := fmt.Sprintf("/wallets/%s/limit_orders/open/pools", wallet)
	var resp PaginationResponse_OpenLimitOrderPoolSummary
	if err := c.doRequest(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOpenLimitOrdersForPool retrieves open limit orders details for a wallet in a pool.
func (c *Client) GetOpenLimitOrdersForPool(ctx context.Context, wallet string, poolAddress string, params *GetLimitOrdersParams) (*OpenLimitOrdersForPoolResponse, error) {
	if wallet == "" || poolAddress == "" {
		return nil, fmt.Errorf("wallet address and pool address cannot be empty")
	}

	q := url.Values{}
	addPaginationParams(q, params)

	path := fmt.Sprintf("/wallets/%s/limit_orders/open/pools/%s", wallet, poolAddress)
	var resp OpenLimitOrdersForPoolResponse
	if err := c.doRequest(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBonusClaimedForPool retrieves realized bonus claimed by a wallet in a pool.
func (c *Client) GetBonusClaimedForPool(ctx context.Context, wallet string, poolAddress string) (*BonusClaimedResponse, error) {
	if wallet == "" || poolAddress == "" {
		return nil, fmt.Errorf("wallet address and pool address cannot be empty")
	}

	path := fmt.Sprintf("/wallets/%s/limit_orders/pools/%s/bonus_claimed", wallet, poolAddress)
	var resp BonusClaimedResponse
	if err := c.doRequest(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetLimitOrderSummary retrieves aggregate limit order totals for a wallet.
func (c *Client) GetLimitOrderSummary(ctx context.Context, wallet string) (*LimitOrderSummaryResponse, error) {
	if wallet == "" {
		return nil, fmt.Errorf("wallet address cannot be empty")
	}

	path := fmt.Sprintf("/wallets/%s/limit_orders/summary", wallet)
	var resp LimitOrderSummaryResponse
	if err := c.doRequest(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWalletPoolTotalClaims retrieves combined claimed fees and rewards for a wallet in a pool.
func (c *Client) GetWalletPoolTotalClaims(ctx context.Context, wallet string, poolAddress string) (*GetWalletTotalClaimsResponse, error) {
	if wallet == "" || poolAddress == "" {
		return nil, fmt.Errorf("wallet address and pool address cannot be empty")
	}

	path := fmt.Sprintf("/wallets/%s/pools/%s/total_claims", wallet, poolAddress)
	var resp GetWalletTotalClaimsResponse
	if err := c.doRequest(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
