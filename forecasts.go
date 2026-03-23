package bitbank

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) Coins(ctx context.Context) ([]Coin, error) {
	var out CoinsResponse
	if err := c.doGet(ctx, "/api/coins", nil, &out); err != nil {
		return nil, err
	}
	return out.Results, nil
}

func (c *Client) HourlyForecasts(ctx context.Context) ([]Forecast, error) {
	raw, err := c.doGetRaw(ctx, "/api/coins/forecasts/hourly", nil)
	if err != nil {
		return nil, err
	}
	var list []Forecast
	if err := json.Unmarshal(raw, &list); err == nil {
		return list, nil
	}
	var wrap struct {
		Forecasts []Forecast `json:"forecasts"`
	}
	if err := json.Unmarshal(raw, &wrap); err != nil {
		return nil, err
	}
	return wrap.Forecasts, nil
}

func (c *Client) HourlyForecastPair(ctx context.Context, pair string) (*Forecast, error) {
	var out Forecast
	if err := c.doGet(ctx, fmt.Sprintf("/api/coins/forecasts/hourly/%s", pair), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Forecasts(ctx context.Context) (json.RawMessage, error) {
	return c.doGetRaw(ctx, "/api/forecasts", nil)
}

func (c *Client) ForecastPair(ctx context.Context, pair string) (*Forecast, error) {
	var out Forecast
	if err := c.doGet(ctx, fmt.Sprintf("/api/forecasts/%s", pair), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ForecastBars(ctx context.Context, pair string, steps int) (*ForecastBarsResponse, error) {
	params := url.Values{}
	params.Set("steps", fmt.Sprintf("%d", steps))
	var out ForecastBarsResponse
	if err := c.doGet(ctx, fmt.Sprintf("/api/forecasts/%s/bars", pair), params, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) HomepageForecast(ctx context.Context, pair string) (HomepageForecastResponse, error) {
	var out HomepageForecastResponse
	if err := c.doGet(ctx, fmt.Sprintf("/api/homepage/forecast/%s", pair), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) DailyForecast(ctx context.Context, pair string) (DailyForecastResponse, error) {
	var out DailyForecastResponse
	if err := c.doGet(ctx, fmt.Sprintf("/api/forecast/daily/%s", pair), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) LineForecast(ctx context.Context, req LineForecastRequest) (*LineForecastResult, error) {
	var out LineForecastResult
	if err := c.doPost(ctx, "/api/forecast/line", req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
