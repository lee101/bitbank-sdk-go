package bitbank

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := NewClient(WithBaseURL(srv.URL), WithAPIKey("test-key"))
	return c, srv
}

func jsonHandler(t *testing.T, v interface{}) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(v)
	}
}

func TestPublicSummary(t *testing.T) {
	c, _ := testServer(t, jsonHandler(t, PublicSummary{TradeCount: 42, PeriodDays: 30}))
	ctx := context.Background()
	out, err := c.PublicSummary(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if out.TradeCount != 42 {
		t.Errorf("got %d, want 42", out.TradeCount)
	}
}

func TestPublicPerformance(t *testing.T) {
	ret := 12.5
	c, _ := testServer(t, jsonHandler(t, PublicPerformance{TotalReturnPct: &ret, TradeCount: 10}))
	out, err := c.PublicPerformance(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if *out.TotalReturnPct != 12.5 {
		t.Errorf("got %f, want 12.5", *out.TotalReturnPct)
	}
}

func TestPublicSignalsLatest(t *testing.T) {
	resp := PublicSignalsLatestResponse{
		Signals: []PublicSignalLatest{{CurrencyPair: "BTC_USDT", BuyPrice: 50000, SellPrice: 51000}},
	}
	c, _ := testServer(t, jsonHandler(t, resp))
	out, err := c.PublicSignalsLatest(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Signals) != 1 {
		t.Fatalf("got %d signals, want 1", len(out.Signals))
	}
	if out.Signals[0].CurrencyPair != "BTC_USDT" {
		t.Errorf("got %s, want BTC_USDT", out.Signals[0].CurrencyPair)
	}
}

func TestFastForecasts(t *testing.T) {
	resp := FastForecastsResponse{
		Forecasts: []FastForecast{{Pair: "ETH_USDT", Signal: "buy", Buy: 3000, Sell: 3100}},
		Count:     1,
	}
	c, _ := testServer(t, jsonHandler(t, resp))
	out, err := c.FastForecasts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if out.Count != 1 {
		t.Errorf("got %d, want 1", out.Count)
	}
}

func TestTrades(t *testing.T) {
	resp := TradesResponse{
		Results: []Trade{{CurrencyPair: "BTC_USDT", Side: "buy", PnLPct: 2.5}},
		Total:   100,
		Limit:   50,
		Offset:  0,
	}
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", r.URL.Query().Get("limit"))
		}
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Errorf("expected X-API-Key header")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	out, err := c.Trades(context.Background(), 50, 0, "")
	if err != nil {
		t.Fatal(err)
	}
	if !out.HasNext() {
		t.Error("expected HasNext=true")
	}
	if out.Results[0].PnLPct != 2.5 {
		t.Errorf("got %f, want 2.5", out.Results[0].PnLPct)
	}
}

func TestPnLSummary(t *testing.T) {
	resp := PnlSummary{TotalTrades: 200, WinRate: 0.65, TotalPnLPct: 45.2}
	c, _ := testServer(t, jsonHandler(t, resp))
	out, err := c.PnLSummary(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if out.TotalTrades != 200 {
		t.Errorf("got %d, want 200", out.TotalTrades)
	}
}

func TestForecastAccuracy(t *testing.T) {
	mae := 1.5
	resp := ForecastAccuracy{
		LookbackDays: 30,
		Horizon1h:    &ErrorSummary{MAEPct: &mae, Samples: 1000},
	}
	c, _ := testServer(t, jsonHandler(t, resp))
	out, err := c.ForecastAccuracy(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if *out.Horizon1h.MAEPct != 1.5 {
		t.Errorf("got %f, want 1.5", *out.Horizon1h.MAEPct)
	}
}

func TestLineForecast(t *testing.T) {
	resp := LineForecastResult{
		SeriesLength:     100,
		PredictionLength: 24,
		Forecast:         LineForecastValues{Median: []float64{1.0, 2.0, 3.0}},
		CreditsUsed:      0.1,
	}
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var req LineForecastRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.PredictionLength != 24 {
			t.Errorf("got %d, want 24", req.PredictionLength)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	out, err := c.LineForecast(context.Background(), LineForecastRequest{
		Series:           []float64{100, 101, 99.5},
		PredictionLength: 24,
		IntervalMinutes:  60,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Forecast.Median) != 3 {
		t.Errorf("got %d medians, want 3", len(out.Forecast.Median))
	}
}

func TestLogin(t *testing.T) {
	resp := LoginResponse{Success: true, Secret: "new-secret", User: &UserInfo{Email: "a@b.com"}}
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Email != "a@b.com" {
			t.Errorf("got %s, want a@b.com", req.Email)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	out, err := c.Login(context.Background(), "a@b.com", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if !out.Success {
		t.Error("expected success=true")
	}
}

func TestErrorMapping(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized", "message": "bad key"})
	})
	_, err := c.Signals(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*AuthenticationError); !ok {
		t.Errorf("expected AuthenticationError, got %T", err)
	}
}

func TestRateLimitError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		json.NewEncoder(w).Encode(map[string]string{"error": "rate_limit", "message": "slow down"})
	})
	_, err := c.PublicSummary(context.Background())
	if _, ok := err.(*RateLimitError); !ok {
		t.Errorf("expected RateLimitError, got %T", err)
	}
}

func TestCoins(t *testing.T) {
	resp := CoinsResponse{Results: []Coin{{CurrencyPair: "BTC_USDT", Name: "Bitcoin", Symbol: "BTC"}}}
	c, _ := testServer(t, jsonHandler(t, resp))
	out, err := c.Coins(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Symbol != "BTC" {
		t.Errorf("unexpected coins: %+v", out)
	}
}

func TestEquityCurve(t *testing.T) {
	resp := EquityCurveResponse{
		Data:          []EquityPoint{{Timestamp: "2024-01-01", Value: 10000}},
		Summary:       &EquitySummary{TotalReturnPct: 25.0},
		InitialEquity: 10000,
	}
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("pair") != "BTC_USDT" {
			t.Errorf("expected pair=BTC_USDT")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	out, err := c.EquityCurve(context.Background(), "BTC_USDT", false)
	if err != nil {
		t.Fatal(err)
	}
	if out.Summary.TotalReturnPct != 25.0 {
		t.Errorf("got %f, want 25.0", out.Summary.TotalReturnPct)
	}
}

func TestSignalSpreadPct(t *testing.T) {
	s := Signal{BuyPrice: 100, SellPrice: 105}
	if s.SpreadPct() != 0.05 {
		t.Errorf("got %f, want 0.05", s.SpreadPct())
	}
	s2 := Signal{BuyPrice: 0, SellPrice: 105}
	if s2.SpreadPct() != 0 {
		t.Errorf("got %f, want 0", s2.SpreadPct())
	}
}

func TestSession(t *testing.T) {
	resp := SessionResponse{User: &UserInfo{Email: "test@test.com", ID: "123"}}
	c, _ := testServer(t, jsonHandler(t, resp))
	out, err := c.Session(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if out.User.Email != "test@test.com" {
		t.Errorf("got %s, want test@test.com", out.User.Email)
	}
}

func TestAPIUsage(t *testing.T) {
	resp := ApiUsage{Total: 500}
	c, _ := testServer(t, jsonHandler(t, resp))
	out, err := c.APIUsage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if out.Total != 500 {
		t.Errorf("got %d, want 500", out.Total)
	}
}

func TestForecastBars(t *testing.T) {
	o := 100.0
	resp := ForecastBarsResponse{
		Bars: []ForecastBar{{Timestamp: "2024-01-01", Open: &o}},
		Pair: "BTC_USDT",
	}
	c, _ := testServer(t, jsonHandler(t, resp))
	out, err := c.ForecastBars(context.Background(), "BTC_USDT", 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Bars) != 1 {
		t.Errorf("got %d bars, want 1", len(out.Bars))
	}
}

func TestNewClientDefaults(t *testing.T) {
	c := NewClient()
	if c.baseURL != DefaultBaseURL {
		t.Errorf("got %s, want %s", c.baseURL, DefaultBaseURL)
	}
}

func TestNewClientOptions(t *testing.T) {
	c := NewClient(WithAPIKey("mykey"), WithBaseURL("https://custom.url"))
	if c.apiKey != "mykey" {
		t.Errorf("got %s, want mykey", c.apiKey)
	}
	if c.baseURL != "https://custom.url" {
		t.Errorf("got %s, want https://custom.url", c.baseURL)
	}
}

func TestServerError(t *testing.T) {
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal", "message": "boom"})
	})
	_, err := c.PublicSummary(context.Background())
	if _, ok := err.(*ServerError); !ok {
		t.Errorf("expected ServerError, got %T", err)
	}
}

func TestSignalExport(t *testing.T) {
	resp := SignalExportResponse{Days: 30, PairCount: 5, Pairs: []string{"BTC_USDT"}}
	c, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("days") != "60" {
			t.Errorf("expected days=60, got %s", r.URL.Query().Get("days"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	out, err := c.SignalExport(context.Background(), 60)
	if err != nil {
		t.Fatal(err)
	}
	if out.PairCount != 5 {
		t.Errorf("got %d, want 5", out.PairCount)
	}
}
