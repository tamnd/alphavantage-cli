// Package alphavantage is the library behind the alphavantage command line:
// the HTTP client, request shaping, and the typed data models for the
// Alpha Vantage financial data API at https://www.alphavantage.co.
//
// The API requires an API key in the apikey= query parameter. The "demo" key
// works for a limited set of symbols and endpoints (IBM quote, tesco search,
// WTI/BRENT oil prices, US inflation data). Register a free key at
// https://www.alphavantage.co/support/#api-key for full access.
package alphavantage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Host is the Alpha Vantage API hostname.
const Host = "www.alphavantage.co"

// DefaultUserAgent identifies the client to the API.
const DefaultUserAgent = "alphavantage-cli/0.1.0 (github.com/tamnd/alphavantage-cli)"

// Config holds constructor parameters for Client.
type Config struct {
	BaseURL string
	APIKey  string
	Rate    time.Duration
	Timeout time.Duration
	Retries int
}

// DefaultConfig returns production defaults.
func DefaultConfig() Config {
	return Config{
		BaseURL: "https://www.alphavantage.co",
		APIKey:  "demo",
		Rate:    500 * time.Millisecond,
		Timeout: 30 * time.Second,
		Retries: 3,
	}
}

// Client is the Alpha Vantage HTTP client.
type Client struct {
	cfg  Config
	http *http.Client
	mu   sync.Mutex
	last time.Time
}

// NewClient constructs a Client from cfg.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}

// --- Public output types ---

// Quote is the real-time stock quote record.
type Quote struct {
	Symbol        string `kit:"id" json:"symbol"`
	Price         string `json:"price"`
	Open          string `json:"open"`
	High          string `json:"high"`
	Low           string `json:"low"`
	Volume        string `json:"volume"`
	Change        string `json:"change"`
	ChangePercent string `json:"change_percent"`
	PreviousClose string `json:"previous_close"`
	LatestDay     string `json:"latest_day"`
}

// Match is one result from a symbol search.
type Match struct {
	Symbol     string `kit:"id" json:"symbol"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Region     string `json:"region"`
	Currency   string `json:"currency"`
	MatchScore string `json:"match_score"`
}

// DataPoint is one entry in a time-series (commodities, economic indicators).
type DataPoint struct {
	Date  string `kit:"id" json:"date"`
	Value string `json:"value"`
	Name  string `json:"name"`
	Unit  string `json:"unit"`
}

// --- Wire types (internal) ---

type wireQuote struct {
	GlobalQuote struct {
		Symbol        string `json:"01. symbol"`
		Open          string `json:"02. open"`
		High          string `json:"03. high"`
		Low           string `json:"04. low"`
		Price         string `json:"05. price"`
		Volume        string `json:"06. volume"`
		LatestDay     string `json:"07. latest trading day"`
		PrevClose     string `json:"08. previous close"`
		Change        string `json:"09. change"`
		ChangePercent string `json:"10. change percent"`
	} `json:"Global Quote"`
}

type wireSearch struct {
	BestMatches []struct {
		Symbol     string `json:"1. symbol"`
		Name       string `json:"2. name"`
		Type       string `json:"3. type"`
		Region     string `json:"4. region"`
		Currency   string `json:"8. currency"`
		MatchScore string `json:"9. matchScore"`
	} `json:"bestMatches"`
}

type wireSeries struct {
	Name     string `json:"name"`
	Interval string `json:"interval"`
	Unit     string `json:"unit"`
	Data     []struct {
		Date  string `json:"date"`
		Value string `json:"value"`
	} `json:"data"`
}

// --- Client methods ---

// GetQuote fetches a real-time stock quote for the given symbol.
func (c *Client) GetQuote(ctx context.Context, symbol string) (*Quote, error) {
	u := c.buildURL("GLOBAL_QUOTE", "symbol", symbol)
	body, err := c.get(ctx, u)
	if err != nil {
		return nil, err
	}
	var w wireQuote
	if err := json.Unmarshal(body, &w); err != nil {
		return nil, fmt.Errorf("quote decode: %w", err)
	}
	q := w.GlobalQuote
	if q.Symbol == "" {
		return nil, fmt.Errorf("no quote returned for symbol %q (demo key only works for IBM)", symbol)
	}
	return &Quote{
		Symbol:        q.Symbol,
		Price:         q.Price,
		Open:          q.Open,
		High:          q.High,
		Low:           q.Low,
		Volume:        q.Volume,
		Change:        q.Change,
		ChangePercent: q.ChangePercent,
		PreviousClose: q.PrevClose,
		LatestDay:     q.LatestDay,
	}, nil
}

// Search searches for symbols matching the given keywords.
func (c *Client) Search(ctx context.Context, keywords string) ([]Match, error) {
	u := c.buildURL("SYMBOL_SEARCH", "keywords", keywords)
	body, err := c.get(ctx, u)
	if err != nil {
		return nil, err
	}
	var w wireSearch
	if err := json.Unmarshal(body, &w); err != nil {
		return nil, fmt.Errorf("search decode: %w", err)
	}
	out := make([]Match, len(w.BestMatches))
	for i, m := range w.BestMatches {
		out[i] = Match{
			Symbol:     m.Symbol,
			Name:       m.Name,
			Type:       m.Type,
			Region:     m.Region,
			Currency:   m.Currency,
			MatchScore: m.MatchScore,
		}
	}
	return out, nil
}

// GetWTI fetches WTI crude oil price data.
// interval is daily, weekly, or monthly. limit caps the number of data points returned.
func (c *Client) GetWTI(ctx context.Context, interval string, limit int) ([]DataPoint, error) {
	return c.getSeries(ctx, "WTI", interval, limit)
}

// GetBrent fetches Brent crude oil price data.
// interval is daily, weekly, or monthly. limit caps the number of data points returned.
func (c *Client) GetBrent(ctx context.Context, interval string, limit int) ([]DataPoint, error) {
	return c.getSeries(ctx, "BRENT", interval, limit)
}

// GetInflation fetches US consumer price inflation data.
// limit caps the number of annual data points returned.
func (c *Client) GetInflation(ctx context.Context, limit int) ([]DataPoint, error) {
	return c.getSeries(ctx, "INFLATION", "", limit)
}

// getSeries is the shared helper for commodity and economic indicator endpoints.
func (c *Client) getSeries(ctx context.Context, function, interval string, limit int) ([]DataPoint, error) {
	params := []string{}
	if interval != "" {
		params = append(params, "interval", interval)
	}
	u := c.buildURL(function, params...)
	body, err := c.get(ctx, u)
	if err != nil {
		return nil, err
	}
	var w wireSeries
	if err := json.Unmarshal(body, &w); err != nil {
		return nil, fmt.Errorf("%s decode: %w", function, err)
	}
	data := w.Data
	if limit > 0 && len(data) > limit {
		data = data[:limit]
	}
	out := make([]DataPoint, len(data))
	for i, d := range data {
		out[i] = DataPoint{
			Date:  d.Date,
			Value: d.Value,
			Name:  w.Name,
			Unit:  w.Unit,
		}
	}
	return out, nil
}

// buildURL assembles a full API URL with the apikey parameter and optional
// additional key=value pairs. Pairs with empty values are skipped.
func (c *Client) buildURL(function string, params ...string) string {
	apiKey := c.cfg.APIKey
	if apiKey == "" {
		apiKey = "demo"
	}
	u := c.cfg.BaseURL + "/query?function=" + function + "&apikey=" + apiKey
	for i := 0; i+1 < len(params); i += 2 {
		if params[i+1] != "" {
			u += "&" + params[i] + "=" + params[i+1]
		}
	}
	return u
}

// get fetches a URL with pacing and retries. The body is fully read and closed.
func (c *Client) get(ctx context.Context, rawURL string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		body, retry, err := c.do(ctx, rawURL)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("get %s: %w", rawURL, lastErr)
}

func (c *Client) do(ctx context.Context, rawURL string) (body []byte, retry bool, err error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", DefaultUserAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, true, err
	}
	return b, false, nil
}

func (c *Client) pace() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	d := time.Duration(attempt) * 500 * time.Millisecond
	if d > 5*time.Second {
		d = 5 * time.Second
	}
	return d
}
