package bitbank

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) Signals(ctx context.Context) (*SignalsResponse, error) {
	var out SignalsResponse
	if err := c.doGet(ctx, "/api/trading-bot/signals", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Signal(ctx context.Context, pair string) (*SignalsResponse, error) {
	var out SignalsResponse
	if err := c.doGet(ctx, fmt.Sprintf("/api/trading-bot/signals/%s", pair), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) EquityCurve(ctx context.Context, pair string, includeTrades bool) (*EquityCurveResponse, error) {
	params := url.Values{}
	if pair != "" {
		params.Set("pair", pair)
	}
	if includeTrades {
		params.Set("include_trades", "true")
	}
	var out EquityCurveResponse
	if err := c.doGet(ctx, "/api/trading-bot/equity-curve", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Trades(ctx context.Context, limit, offset int, pair string) (*TradesResponse, error) {
	params := url.Values{}
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("offset", fmt.Sprintf("%d", offset))
	if pair != "" {
		params.Set("pair", pair)
	}
	var out TradesResponse
	if err := c.doGet(ctx, "/api/trading-bot/trades", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PnLSummary(ctx context.Context) (*PnlSummary, error) {
	var out PnlSummary
	if err := c.doGet(ctx, "/api/trading-bot/pnl-summary", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Pairs(ctx context.Context) ([]string, error) {
	var out PairsResponse
	if err := c.doGet(ctx, "/api/trading-bot/pairs", nil, &out); err != nil {
		return nil, err
	}
	return out.Pairs, nil
}

func (c *Client) PublicSummary(ctx context.Context) (*PublicSummary, error) {
	var out PublicSummary
	if err := c.doGet(ctx, "/api/trading-bot/public-summary", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PublicPerformance(ctx context.Context) (*PublicPerformance, error) {
	var out PublicPerformance
	if err := c.doGet(ctx, "/api/trading-bot/public-performance", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ForecastAccuracy(ctx context.Context, lookbackDays *int) (*ForecastAccuracy, error) {
	params := url.Values{}
	if lookbackDays != nil {
		params.Set("lookback_days", fmt.Sprintf("%d", *lookbackDays))
	}
	var out ForecastAccuracy
	if err := c.doGet(ctx, "/api/trading-bot/forecast-accuracy", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PublicSignalsHistory(ctx context.Context, pair string, days *int) (*SignalHistoryResponse, error) {
	params := url.Values{}
	params.Set("pair", pair)
	if days != nil {
		params.Set("days", fmt.Sprintf("%d", *days))
	}
	var out SignalHistoryResponse
	if err := c.doGet(ctx, "/api/trading-bot/public-signals-history", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SignalExport(ctx context.Context, days int) (*SignalExportResponse, error) {
	params := url.Values{}
	params.Set("days", fmt.Sprintf("%d", days))
	var out SignalExportResponse
	if err := c.doGet(ctx, "/api/trading-bot/signal-export", params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) PublicSignalsLatest(ctx context.Context) (*PublicSignalsLatestResponse, error) {
	var out PublicSignalsLatestResponse
	if err := c.doGet(ctx, "/api/trading-bot/public-signals-latest", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) FastForecasts(ctx context.Context) (*FastForecastsResponse, error) {
	var out FastForecastsResponse
	if err := c.doGet(ctx, "/api/trading-bot/fast-forecasts", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) WebhookInfo(ctx context.Context) (WebhookInfo, error) {
	var out WebhookInfo
	if err := c.doGet(ctx, "/api/trading-bot/webhook-info", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
