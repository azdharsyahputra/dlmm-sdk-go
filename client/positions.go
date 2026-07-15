package client

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// GetPositionHistoricalEventsParams represents parameters for GetPositionHistoricalEvents.
type GetPositionHistoricalEventsParams struct {
	EventType      *string
	OrderDirection *string
}

// GetPositionHistoricalEvents retrieves the historical events for a position.
func (c *Client) GetPositionHistoricalEvents(ctx context.Context, address string, params *GetPositionHistoricalEventsParams) (*GetPositionHistoricalEventsResponse, error) {
	if address == "" {
		return nil, errors.New("position address cannot be empty")
	}

	q := url.Values{}
	if params != nil {
		if params.EventType != nil {
			q.Set("event_type", *params.EventType)
		}
		if params.OrderDirection != nil {
			q.Set("order_direction", *params.OrderDirection)
		}
	}

	path := fmt.Sprintf("/positions/%s/historical", address)
	var resp GetPositionHistoricalEventsResponse
	if err := c.doRequest(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetPoolPositionPnLParams represents parameters for GetPoolPositionPnL.
type GetPoolPositionPnLParams struct {
	User     string // required
	Status   *string
	Page     *int
	PageSize *int
}

// GetPoolPositionPnL retrieves position PnL for a wallet in a specific pool.
func (c *Client) GetPoolPositionPnL(ctx context.Context, poolAddress string, params *GetPoolPositionPnLParams) (*GetPoolPositionPnLResponse, error) {
	if poolAddress == "" {
		return nil, errors.New("pool address cannot be empty")
	}

	q := url.Values{}
	q.Set("user", params.User)
	if params.Status != nil {
		q.Set("status", *params.Status)
	}
	if params.Page != nil {
		q.Set("page", strconv.Itoa(*params.Page))
	}
	if params.PageSize != nil {
		q.Set("page_size", strconv.Itoa(*params.PageSize))
	}

	path := fmt.Sprintf("/positions/%s/pnl", poolAddress)
	var resp GetPoolPositionPnLResponse
	if err := c.doRequest(ctx, path, q, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
