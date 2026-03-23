package bitbank

type Signal struct {
	ID             *int     `json:"id,omitempty"`
	CurrencyPair   string   `json:"currency_pair"`
	Timeframe      string   `json:"timeframe"`
	Timestamp      string   `json:"timestamp,omitempty"`
	BuyPrice       float64  `json:"buy_price"`
	SellPrice      float64  `json:"sell_price"`
	TradeAmountPct *float64 `json:"trade_amount_pct,omitempty"`
	Confidence     *float64 `json:"confidence,omitempty"`
	SignalType     string   `json:"signal_type"`
	ChronosLow     *float64 `json:"chronos_low,omitempty"`
	ChronosMedian  *float64 `json:"chronos_median,omitempty"`
	ChronosHigh    *float64 `json:"chronos_high,omitempty"`
	ModelType      string   `json:"model_type,omitempty"`
}

func (s Signal) SpreadPct() float64 {
	if s.BuyPrice <= 0 {
		return 0
	}
	return (s.SellPrice - s.BuyPrice) / s.BuyPrice
}

type SignalsResponse struct {
	Signals   []Signal `json:"signals"`
	Timestamp string   `json:"timestamp,omitempty"`
}

type PublicSignalLatest struct {
	CurrencyPair string   `json:"currency_pair"`
	Timestamp    string   `json:"timestamp,omitempty"`
	BuyPrice     float64  `json:"buy_price"`
	SellPrice    float64  `json:"sell_price"`
	SignalType   string   `json:"signal_type"`
	Confidence   *float64 `json:"confidence,omitempty"`
	PnL7dPct     float64  `json:"pnl_7d_pct"`
	Trades7d     int      `json:"trades_7d"`
	WinRate7d    *float64 `json:"win_rate_7d,omitempty"`
}

type PublicSignalsLatestResponse struct {
	Signals   []PublicSignalLatest `json:"signals"`
	Timestamp string               `json:"timestamp,omitempty"`
}

type FastForecast struct {
	Pair   string   `json:"pair"`
	Ts     string   `json:"ts"`
	Buy    float64  `json:"buy"`
	Sell   float64  `json:"sell"`
	Signal string   `json:"signal"`
	Conf   *float64 `json:"conf,omitempty"`
}

type FastForecastsResponse struct {
	Forecasts []FastForecast `json:"forecasts"`
	Count     int            `json:"count"`
	CacheTTL  float64        `json:"cache_ttl"`
	Ts        string         `json:"ts,omitempty"`
}

type SignalExportResponse struct {
	Days      int                    `json:"days"`
	Pairs     []string               `json:"pairs"`
	PairCount int                    `json:"pair_count"`
	Signals   map[string]interface{} `json:"signals"`
	PnL7d     map[string]interface{} `json:"pnl_7d"`
	Timestamp string                 `json:"timestamp,omitempty"`
}

type SignalHistoryResponse struct {
	Signals []map[string]interface{} `json:"signals"`
	Pair    string                   `json:"pair"`
	Count   int                      `json:"count"`
}

type WebhookInfo map[string]interface{}

type PublicSummary struct {
	TotalReturnPct *float64 `json:"total_return_pct,omitempty"`
	MaxDrawdownPct *float64 `json:"max_drawdown_pct,omitempty"`
	WinRatePct     *float64 `json:"win_rate_pct,omitempty"`
	TradeCount     int      `json:"trade_count"`
	PeriodDays     int      `json:"period_days"`
}

type PairPerformance struct {
	Pair           string   `json:"pair"`
	TotalReturnPct float64  `json:"total_return_pct"`
	SharpeRatio    *float64 `json:"sharpe_ratio,omitempty"`
	SortinoRatio   *float64 `json:"sortino_ratio,omitempty"`
	TradeCount     int      `json:"trade_count"`
	WinRatePct     *float64 `json:"win_rate_pct,omitempty"`
}

type PublicPerformance struct {
	TotalReturnPct *float64          `json:"total_return_pct,omitempty"`
	SharpeRatio    *float64          `json:"sharpe_ratio,omitempty"`
	SortinoRatio   *float64          `json:"sortino_ratio,omitempty"`
	MaxDrawdownPct *float64          `json:"max_drawdown_pct,omitempty"`
	WinRatePct     *float64          `json:"win_rate_pct,omitempty"`
	TradeCount     int               `json:"trade_count"`
	PeriodDays     int               `json:"period_days"`
	Pairs          []PairPerformance `json:"pairs"`
}

type Trade struct {
	ID              *int     `json:"id,omitempty"`
	CurrencyPair    string   `json:"currency_pair"`
	Timeframe       string   `json:"timeframe,omitempty"`
	Side            string   `json:"side"`
	EntryPrice      float64  `json:"entry_price"`
	ExitPrice       float64  `json:"exit_price"`
	EntryTimestamp  string   `json:"entry_timestamp,omitempty"`
	ExitTimestamp   string   `json:"exit_timestamp,omitempty"`
	PositionSizePct *float64 `json:"position_size_pct,omitempty"`
	PnLPct          float64  `json:"pnl_pct"`
	PnLUSD          float64  `json:"pnl_usd"`
	FeesPct         float64  `json:"fees_pct"`
}

type TradesResponse struct {
	Results []Trade `json:"results"`
	Total   int     `json:"total"`
	Limit   int     `json:"limit"`
	Offset  int     `json:"offset"`
}

func (t TradesResponse) HasNext() bool {
	return t.Offset+t.Limit < t.Total
}

type PnlSummary struct {
	TotalTrades   int     `json:"total_trades"`
	WinningTrades int     `json:"winning_trades"`
	LosingTrades  int     `json:"losing_trades"`
	WinRate       float64 `json:"win_rate"`
	TotalPnLPct   float64 `json:"total_pnl_pct"`
	TotalPnLUSD   float64 `json:"total_pnl_usd"`
	AvgWinPct     float64 `json:"avg_win_pct"`
	AvgLossPct    float64 `json:"avg_loss_pct"`
	BestTradePct  float64 `json:"best_trade_pct"`
	WorstTradePct float64 `json:"worst_trade_pct"`
}

type EquityPoint struct {
	Timestamp   string  `json:"timestamp"`
	Value       float64 `json:"value"`
	DrawdownPct float64 `json:"drawdown_pct"`
	TradeCount  int     `json:"trade_count"`
	WinRate     float64 `json:"win_rate"`
}

type EquitySummary struct {
	TotalReturnPct float64 `json:"total_return_pct"`
	MaxDrawdownPct float64 `json:"max_drawdown_pct"`
	WinRate        float64 `json:"win_rate"`
	TradeCount     int     `json:"trade_count"`
}

type EquityCurveResponse struct {
	Data          []EquityPoint  `json:"data"`
	Summary       *EquitySummary `json:"summary,omitempty"`
	InitialEquity float64        `json:"initial_equity"`
	Pair          string         `json:"pair,omitempty"`
	Timeframe     string         `json:"timeframe"`
	Trades        []Trade        `json:"trades,omitempty"`
}

type ErrorSummary struct {
	MAEPct   *float64 `json:"mae_pct,omitempty"`
	StdevPct *float64 `json:"stdev_pct,omitempty"`
	Samples  int      `json:"samples"`
}

type PairAccuracy struct {
	Pair         string        `json:"pair"`
	Horizon1h    *ErrorSummary `json:"horizon_1h,omitempty"`
	Horizon1d    *ErrorSummary `json:"horizon_1d,omitempty"`
	TotalSamples int           `json:"total_samples"`
}

type DailyAccuracy struct {
	Date         string        `json:"date"`
	Horizon1h    *ErrorSummary `json:"horizon_1h,omitempty"`
	TotalSamples int           `json:"total_samples"`
}

type ForecastAccuracy struct {
	LookbackDays int                    `json:"lookback_days"`
	AsOf         string                 `json:"as_of,omitempty"`
	Source       string                 `json:"source,omitempty"`
	Horizon1h    *ErrorSummary          `json:"horizon_1h,omitempty"`
	Horizon1d    *ErrorSummary          `json:"horizon_1d,omitempty"`
	ByPair       []PairAccuracy         `json:"by_pair"`
	TimeSeries   []DailyAccuracy        `json:"time_series"`
	Volatility   map[string]interface{} `json:"volatility"`
}

type Forecast struct {
	CurrencyPair  string   `json:"currency_pair"`
	Timestamp     string   `json:"timestamp,omitempty"`
	BuyPrice      *float64 `json:"buy_price,omitempty"`
	SellPrice     *float64 `json:"sell_price,omitempty"`
	SignalType    string   `json:"signal_type,omitempty"`
	Confidence    *float64 `json:"confidence,omitempty"`
	ChronosLow    *float64 `json:"chronos_low,omitempty"`
	ChronosMedian *float64 `json:"chronos_median,omitempty"`
	ChronosHigh   *float64 `json:"chronos_high,omitempty"`
}

type ForecastBar struct {
	Timestamp string   `json:"timestamp"`
	Open      *float64 `json:"open,omitempty"`
	High      *float64 `json:"high,omitempty"`
	Low       *float64 `json:"low,omitempty"`
	Close     *float64 `json:"close,omitempty"`
}

type ForecastBarsResponse struct {
	Bars []ForecastBar `json:"bars"`
	Pair string        `json:"pair"`
}

type HourlyForecast map[string]interface{}
type DailyForecastResponse map[string]interface{}
type HomepageForecastResponse map[string]interface{}

type LineForecastRequest struct {
	Series           []float64 `json:"series"`
	PredictionLength int       `json:"prediction_length"`
	IntervalMinutes  int       `json:"interval_minutes"`
}

type LineForecastValues struct {
	Low    []float64 `json:"low"`
	Median []float64 `json:"median"`
	High   []float64 `json:"high"`
}

type LineForecastResult struct {
	SeriesLength       int                `json:"series_length"`
	PredictionLength   int                `json:"prediction_length"`
	IntervalMinutes    int                `json:"interval_minutes"`
	Forecast           LineForecastValues `json:"forecast"`
	Engine             string             `json:"engine,omitempty"`
	GeneratedAt        string             `json:"generated_at,omitempty"`
	CreditsUsed        float64            `json:"credits_used"`
	CreditsRemaining   float64            `json:"credits_remaining"`
	ForecastTimestamps []string           `json:"forecast_timestamps,omitempty"`
}

type Coin struct {
	CurrencyPair string   `json:"currency_pair"`
	Name         string   `json:"name"`
	Symbol       string   `json:"symbol"`
	Price        *float64 `json:"price,omitempty"`
	Change24h    *float64 `json:"change_24h,omitempty"`
	Volume24h    *float64 `json:"volume_24h,omitempty"`
}

type CoinsResponse struct {
	Results []Coin `json:"results"`
}

type UserInfo struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Success bool      `json:"success"`
	User    *UserInfo `json:"user,omitempty"`
	Secret  string    `json:"secret"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

type SignupResponse struct {
	Success bool      `json:"success"`
	User    *UserInfo `json:"user,omitempty"`
	Message string    `json:"message"`
}

type SessionResponse struct {
	User *UserInfo `json:"user,omitempty"`
}

type ApiUsage struct {
	Requests []map[string]interface{} `json:"requests"`
	Total    int                      `json:"total"`
}

type UsageEntry struct {
	Endpoint  string  `json:"endpoint"`
	Count     int     `json:"count"`
	TotalCost float64 `json:"total_cost"`
}

type AutotopupSettings struct {
	AutotopupEnabled   bool    `json:"autotopup_enabled"`
	AutotopupThreshold float64 `json:"autotopup_threshold"`
	AutotopupAmount    float64 `json:"autotopup_amount"`
}

type CreditsPurchase struct {
	Success     bool    `json:"success"`
	Credits     float64 `json:"credits"`
	CheckoutURL string  `json:"checkout_url"`
}

type PairsResponse struct {
	Pairs []string `json:"pairs"`
}
