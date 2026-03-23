# bitbank-sdk-go

Go SDK for the [bitbank.nz](https://bitbank.nz) cryptocurrency trading API. Get forecasts, trading signals, performance data, and execute account operations.

## Install

```bash
go get github.com/lee101/bitbank-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"

    bitbank "github.com/lee101/bitbank-sdk-go"
)

func main() {
    client := bitbank.NewClient()
    ctx := context.Background()

    // Public endpoints (no API key needed)
    signals, _ := client.PublicSignalsLatest(ctx)
    for _, s := range signals.Signals {
        fmt.Printf("%s: buy=%.2f sell=%.2f pnl_7d=%.2f%%\n",
            s.CurrencyPair, s.BuyPrice, s.SellPrice, s.PnL7dPct)
    }

    forecasts, _ := client.FastForecasts(ctx)
    for _, f := range forecasts.Forecasts {
        fmt.Printf("%s: %s conf=%v\n", f.Pair, f.Signal, f.Conf)
    }

    accuracy, _ := client.ForecastAccuracy(ctx, nil)
    if accuracy.Horizon1h != nil {
        fmt.Printf("1h MAE: %v%%\n", *accuracy.Horizon1h.MAEPct)
    }
}
```

## Authenticated Usage

```go
client := bitbank.NewClient(bitbank.WithAPIKey("your-api-key"))
ctx := context.Background()

// Trading signals (subscriber only)
signals, _ := client.Signals(ctx)

// Equity curve
equity, _ := client.EquityCurve(ctx, "", false)
fmt.Printf("Return: %.2f%%\n", equity.Summary.TotalReturnPct)

// Trade history (paginated)
trades, _ := client.Trades(ctx, 50, 0, "")
for trades.HasNext() {
    trades, _ = client.Trades(ctx, 50, trades.Offset+trades.Limit, "")
}

// Custom time series forecast (costs credits)
result, _ := client.LineForecast(ctx, bitbank.LineForecastRequest{
    Series:           []float64{100.0, 101.0, 99.5},
    PredictionLength: 24,
    IntervalMinutes:  60,
})
fmt.Println(result.Forecast.Median)
```

## WebSocket

```go
client := bitbank.NewClient(bitbank.WithAPIKey("your-key"))
ws := client.WebSocket()

ws.On("features_update", func(msg map[string]interface{}) {
    fmt.Println(msg)
})
ws.On("live_prices", func(msg map[string]interface{}) {
    fmt.Println(msg)
})

ws.Connect()
ws.Subscribe("BTC_USDT")
ws.SubscribeLive([]string{"BTC_USDT", "ETH_USDT"})

// ... keep running ...
ws.Disconnect()
```

## Functional Options

```go
client := bitbank.NewClient(
    bitbank.WithAPIKey("your-key"),
    bitbank.WithBaseURL("https://custom-url.com"),
    bitbank.WithTimeout(60 * time.Second),
)
```

## Available Methods

| Method | Auth | Description |
|--------|------|-------------|
| `PublicSummary` | No | Trading bot summary stats |
| `PublicPerformance` | No | Detailed performance with per-pair breakdown |
| `ForecastAccuracy` | No | Forecast accuracy metrics |
| `PublicSignalsHistory` | No | Historical signals for a pair |
| `SignalExport` | No | Bulk signal export |
| `PublicSignalsLatest` | No | Latest signals for all pairs |
| `FastForecasts` | No | Fast forecast snapshot |
| `WebhookInfo` | No | Webhook configuration info |
| `Coins` | No | List available coins |
| `HourlyForecasts` | No | Hourly forecasts for all coins |
| `HourlyForecastPair` | No | Hourly forecast for a specific pair |
| `Forecasts` | No | All forecasts (raw JSON) |
| `ForecastPair` | No | Forecast for a specific pair |
| `ForecastBars` | No | OHLC forecast bars |
| `HomepageForecast` | No | Homepage forecast for a pair |
| `DailyForecast` | No | Daily forecast for a pair |
| `Signals` | Yes | Current trading signals |
| `Signal` | Yes | Signal for a specific pair |
| `EquityCurve` | Yes | Equity curve data |
| `Trades` | Yes | Trade history (paginated) |
| `PnLSummary` | Yes | PnL summary |
| `Pairs` | Yes | Active trading pairs |
| `LineForecast` | Yes | Custom time series forecast |
| `APIUsage` | Yes | API usage stats |
| `UsageHistory` | Yes | Usage history for a user |
| `RegenerateKey` | Yes | Regenerate API key |
| `PurchaseCredits` | Yes | Purchase credits |
| `SaveAutotopup` | Yes | Save auto top-up settings |
| `GetAutotopup` | Yes | Get auto top-up settings |
| `Login` | No | Login with email/password |
| `Signup` | No | Create account |
| `Session` | Yes | Get current session |
| `Logout` | Yes | Logout |

## Error Handling

```go
signals, err := client.Signals(ctx)
if err != nil {
    switch err.(type) {
    case *bitbank.AuthenticationError:
        fmt.Println("Invalid API key")
    case *bitbank.InsufficientCreditsError:
        fmt.Println("Need more credits")
    case *bitbank.RateLimitError:
        fmt.Println("Rate limited")
    default:
        fmt.Println(err)
    }
}
```

## Related

- [bitbank-sdk-python](https://github.com/lee101/bitbank-sdk-python) - Python SDK
- [bitbank-sdk-typescript](https://github.com/lee101/bitbank-sdk-typescript) - TypeScript SDK
- [bitbank-binance-bot](https://github.com/lee101/bitbank-binance-bot) - Trading bot using these SDKs

## License

MIT
