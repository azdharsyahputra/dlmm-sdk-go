package client

import (
	"context"
	"net/url"
	"strconv"
)

// GetPortfolioParams represents parameters for GetPortfolio.
type GetPortfolioParams struct {
	User     string // required
	Page     *int
	PageSize *int
	DaysBack *int
}

// GetPortfolio retrieves user portfolio with all pools containing closed positions.
func (c *Client) GetPortfolio(ctx context.Context, user string, params *GetPortfolioParams) (*PortfolioResponse, error) {
	q := url.Values{}
	q.Set("user", params.User)
	if params.Page != nil {
		q.Set("page", strconv.Itoa(*params.Page))
	}
	if params.PageSize != nil {
		q.Set("page_size", strconv.Itoa(*params.PageSize))
	}
	if params.DaysBack != nil {
		q.Set("days_back", strconv.Itoa(*params.DaysBack))
	}

	var resp PortfolioResponse
	if err := c.doRequest(ctx, "/portfolio", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOpenPortfolioParams represents parameters for GetOpenPortfolio.
type GetOpenPortfolioParams struct {
	User          string // required
	Page          *int
	PageSize      *int
	SortDirection *string
	SortBy        *string
}

// GetOpenPortfolio retrieves user open portfolio.
func (c *Client) GetOpenPortfolio(ctx context.Context, params *GetOpenPortfolioParams) (*OpenPortfolioResponse, error) {
	q := url.Values{}
	q.Set("user", params.User)
	if params.Page != nil {
		q.Set("page", strconv.Itoa(*params.Page))
	}
	if params.PageSize != nil {
		q.Set("page_size", strconv.Itoa(*params.PageSize))
	}
	if params.SortDirection != nil {
		q.Set("sort_direction", *params.SortDirection)
	}
	if params.SortBy != nil {
		q.Set("sort_by", *params.SortBy)
	}

	var resp OpenPortfolioResponse
	if err := c.doRequest(ctx, "/portfolio/open", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPortfolioTotal retrieves total portfolio PnL across all pools.
func (c *Client) GetPortfolioTotal(ctx context.Context, user string) (*PortfolioTotalResponse, error) {
	q := url.Values{}
	q.Set("user", user)

	var resp PortfolioTotalResponse
	if err := c.doRequest(ctx, "/portfolio/total", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
