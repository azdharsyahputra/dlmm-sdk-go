package client

import (
	"context"
	"net/url"
	"strconv"
)

// GetDailyMetricParams represents parameters for daily metrics endpoints.
type GetDailyMetricParams struct {
	StartTime *int64
	EndTime   *int64
	View      *string
}

func (c *Client) getDailyMetric(ctx context.Context, path string, params *GetDailyMetricParams) (*DailyMetricResponse, error) {
	q := url.Values{}
	if params != nil {
		if params.StartTime != nil {
			q.Set("start_time", strconv.FormatInt(*params.StartTime, 10))
		}
		if params.EndTime != nil {
			q.Set("end_time", strconv.FormatInt(*params.EndTime, 10))
		}
		if params.View != nil {
			q.Set("view", *params.View)
		}
	}

	var resp DailyMetricResponse
	if err := c.doRequest(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDailyProtocolFees retrieves daily protocol fees.
func (c *Client) GetDailyProtocolFees(ctx context.Context, params *GetDailyMetricParams) (*DailyMetricResponse, error) {
	return c.getDailyMetric(ctx, "/stats/daily/protocol_fees", params)
}

// GetDailyTradingFees retrieves daily trading fees.
func (c *Client) GetDailyTradingFees(ctx context.Context, params *GetDailyMetricParams) (*DailyMetricResponse, error) {
	return c.getDailyMetric(ctx, "/stats/daily/trading_fees", params)
}

// GetDailyTradingVolume retrieves daily trading volume.
func (c *Client) GetDailyTradingVolume(ctx context.Context, params *GetDailyMetricParams) (*DailyMetricResponse, error) {
	return c.getDailyMetric(ctx, "/stats/daily/volume", params)
}

// GetProtocolOverview retrieves aggregated protocol overview metrics.
func (c *Client) GetProtocolOverview(ctx context.Context) (*ProtocolMetricsResponse, error) {
	var resp ProtocolMetricsResponse
	if err := c.doRequest(ctx, "/stats/protocol_metrics", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
