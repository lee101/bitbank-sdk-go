package bitbank

import "context"

func (c *Client) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	var out LoginResponse
	if err := c.doPost(ctx, "/api/login", LoginRequest{Email: email, Password: password}, &out); err != nil {
		return nil, err
	}
	if out.Secret != "" {
		c.apiKey = out.Secret
	}
	return &out, nil
}

func (c *Client) Signup(ctx context.Context, email, password, name string) (*SignupResponse, error) {
	var out SignupResponse
	if err := c.doPost(ctx, "/api/signup", SignupRequest{Email: email, Password: password, Name: name}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Session(ctx context.Context) (*SessionResponse, error) {
	var out SessionResponse
	if err := c.doGet(ctx, "/api/session", nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Logout(ctx context.Context) error {
	if err := c.doPost(ctx, "/api/logout", nil, nil); err != nil {
		return err
	}
	c.apiKey = ""
	return nil
}
