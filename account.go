package bitbank

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) APIUsage(ctx context.Context) (*ApiUsage, error) {
	var out ApiUsage
	if err := c.doPost(ctx, "/api/get-api-usage", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) UsageHistory(ctx context.Context, userID string) ([]UsageEntry, error) {
	raw, err := c.doGetRaw(ctx, fmt.Sprintf("/api/usage-history/%s", userID), nil)
	if err != nil {
		return nil, err
	}
	var list []UsageEntry
	if err := json.Unmarshal(raw, &list); err == nil {
		return list, nil
	}
	var wrap struct {
		Results []UsageEntry `json:"results"`
	}
	if err := json.Unmarshal(raw, &wrap); err != nil {
		return nil, err
	}
	return wrap.Results, nil
}

func (c *Client) RegenerateKey(ctx context.Context) (map[string]interface{}, error) {
	var out map[string]interface{}
	if err := c.doPost(ctx, "/api/regenerate-key", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) PurchaseCredits(ctx context.Context, amountUSD float64) (*CreditsPurchase, error) {
	body := map[string]interface{}{"amount": amountUSD}
	var out CreditsPurchase
	if err := c.doPost(ctx, "/api/purchase-credits", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) SaveAutotopup(ctx context.Context, enabled bool, threshold, amount float64) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"autotopup_enabled":   enabled,
		"autotopup_threshold": threshold,
		"autotopup_amount":    amount,
	}
	var out map[string]interface{}
	if err := c.doPost(ctx, "/api/save-autotopup-settings", body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetAutotopup(ctx context.Context) (*AutotopupSettings, error) {
	var out AutotopupSettings
	if err := c.doPost(ctx, "/api/autotopup-settings", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
