package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GetPoolsParams represents query parameters for GetPools.
type GetPoolsParams struct {
	Page     *int
	PageSize *int
	Query    *string
	SortBy   *string
	FilterBy *string
}

// GetPools retrieves a paginated list of pools.
func (c *Client) GetPools(ctx context.Context, params *GetPoolsParams) (*PoolsResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.Query != nil {
			q.Set("query", *params.Query)
		}
		if params.SortBy != nil {
			q.Set("sort_by", *params.SortBy)
		}
		if params.FilterBy != nil {
			q.Set("filter_by", *params.FilterBy)
		}
	}

	var resp PoolsResponse
	if err := c.doRequest(ctx, "/pools", q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPool retrieves metadata and current state for a single pool.
func (c *Client) GetPool(ctx context.Context, address string) (*PoolResponse, error) {
	if address == "" {
		return nil, fmt.Errorf("pool address cannot be empty")
	}

	path := fmt.Sprintf("/pools/%s", address)
	var resp PoolResponse
	if err := c.doRequest(ctx, path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetOHLCVParams represents parameters for GetOHLCV.
type GetOHLCVParams struct {
	Timeframe *string
	StartTime *int64
	EndTime   *int64
}

// GetOHLCV retrieves OHLCV candles for a single pool.
func (c *Client) GetOHLCV(ctx context.Context, address string, params *GetOHLCVParams) (*OHLCVResponse, error) {
	if address == "" {
		return nil, fmt.Errorf("pool address cannot be empty")
	}

	q := url.Values{}
	if params != nil {
		if params.Timeframe != nil {
			q.Set("timeframe", *params.Timeframe)
		}
		if params.StartTime != nil {
			q.Set("start_time", strconv.FormatInt(*params.StartTime, 10))
		}
		if params.EndTime != nil {
			q.Set("end_time", strconv.FormatInt(*params.EndTime, 10))
		}
	}

	path := fmt.Sprintf("/pools/%s/ohlcv", address)
	var resp OHLCVResponse
	if err := c.doRequest(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVolumeHistoryParams represents parameters for GetVolumeHistory.
type GetVolumeHistoryParams struct {
	Timeframe *string
	StartTime *int64
	EndTime   *int64
}

// GetVolumeHistory retrieves historical volume for a pool.
func (c *Client) GetVolumeHistory(ctx context.Context, address string, params *GetVolumeHistoryParams) (*VolumeHistoryResponse, error) {
	if address == "" {
		return nil, fmt.Errorf("pool address cannot be empty")
	}

	q := url.Values{}
	if params != nil {
		if params.Timeframe != nil {
			q.Set("timeframe", *params.Timeframe)
		}
		if params.StartTime != nil {
			q.Set("start_time", strconv.FormatInt(*params.StartTime, 10))
		}
		if params.EndTime != nil {
			q.Set("end_time", strconv.FormatInt(*params.EndTime, 10))
		}
	}

	path := fmt.Sprintf("/pools/%s/volume/history", address)
	var resp VolumeHistoryResponse
	if err := c.doRequest(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
