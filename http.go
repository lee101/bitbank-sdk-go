package bitbank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) doGet(ctx context.Context, path string, params url.Values, out interface{}) error {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return &ConnectionError{BitbankError{Message: err.Error()}}
	}
	return c.do(req, out)
}

func (c *Client) doPost(ctx context.Context, path string, body interface{}, out interface{}) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, reader)
	if err != nil {
		return &ConnectionError{BitbankError{Message: err.Error()}}
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req, out)
}

func (c *Client) doGetRaw(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	var raw json.RawMessage
	if err := c.doGet(ctx, path, params, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *Client) do(req *http.Request, out interface{}) error {
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &ConnectionError{BitbankError{Message: err.Error()}}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ConnectionError{BitbankError{Message: err.Error()}}
	}

	if resp.StatusCode >= 400 {
		var errBody struct {
			Error   string `json:"error"`
			Message string `json:"message"`
			Detail  string `json:"detail"`
		}
		if json.Unmarshal(data, &errBody) != nil {
			errBody.Error = "unknown"
			errBody.Message = string(data)
		}
		msg := errBody.Message
		if msg == "" {
			msg = errBody.Detail
		}
		if msg == "" {
			msg = fmt.Sprintf("%s", string(data))
		}
		return mapError(resp.StatusCode, errBody.Error, msg)
	}

	if resp.StatusCode == 204 || out == nil {
		return nil
	}
	return json.Unmarshal(data, out)
}
