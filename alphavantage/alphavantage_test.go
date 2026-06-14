package alphavantage_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/tamnd/alphavantage-cli/alphavantage"
)

func newTestClient(baseURL string) *alphavantage.Client {
	cfg := alphavantage.DefaultConfig()
	cfg.BaseURL = baseURL
	cfg.APIKey = "TEST_KEY"
	cfg.Rate = 0
	cfg.Retries = 0
	cfg.Timeout = 5 * time.Second
	return alphavantage.NewClient(cfg)
}

func TestGetQuote(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/query" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("function") != "GLOBAL_QUOTE" {
			t.Errorf("expected function=GLOBAL_QUOTE, got %s", r.URL.Query().Get("function"))
		}
		if r.URL.Query().Get("symbol") != "IBM" {
			t.Errorf("expected symbol=IBM, got %s", r.URL.Query().Get("symbol"))
		}
		if r.URL.Query().Get("apikey") != "TEST_KEY" {
			t.Errorf("expected apikey=TEST_KEY, got %s", r.URL.Query().Get("apikey"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"Global Quote": {
				"01. symbol": "IBM",
				"02. open": "274.8500",
				"03. high": "275.5700",
				"04. low": "271.3700",
				"05. price": "272.2400",
				"06. volume": "4014849",
				"07. latest trading day": "2026-06-13",
				"08. previous close": "274.8500",
				"09. change": "-2.6100",
				"10. change percent": "-0.9496%"
			}
		}`))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	q, err := c.GetQuote(context.Background(), "IBM")
	if err != nil {
		t.Fatalf("GetQuote: %v", err)
	}
	if q.Symbol != "IBM" {
		t.Errorf("Symbol = %q, want IBM", q.Symbol)
	}
	if q.Price != "272.2400" {
		t.Errorf("Price = %q, want 272.2400", q.Price)
	}
	if q.Open != "274.8500" {
		t.Errorf("Open = %q, want 274.8500", q.Open)
	}
	if q.High != "275.5700" {
		t.Errorf("High = %q, want 275.5700", q.High)
	}
	if q.Low != "271.3700" {
		t.Errorf("Low = %q, want 271.3700", q.Low)
	}
	if q.Volume != "4014849" {
		t.Errorf("Volume = %q, want 4014849", q.Volume)
	}
	if q.Change != "-2.6100" {
		t.Errorf("Change = %q, want -2.6100", q.Change)
	}
	if q.ChangePercent != "-0.9496%" {
		t.Errorf("ChangePercent = %q, want -0.9496%%", q.ChangePercent)
	}
	if q.PreviousClose != "274.8500" {
		t.Errorf("PreviousClose = %q, want 274.8500", q.PreviousClose)
	}
	if q.LatestDay != "2026-06-13" {
		t.Errorf("LatestDay = %q, want 2026-06-13", q.LatestDay)
	}
}

func TestGetQuoteEmptySymbol(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Global Quote": {}}`))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	_, err := c.GetQuote(context.Background(), "ZZZZZ")
	if err == nil {
		t.Error("GetQuote with empty response: want error, got nil")
	}
}

func TestSearch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("function") != "SYMBOL_SEARCH" {
			t.Errorf("expected function=SYMBOL_SEARCH, got %s", r.URL.Query().Get("function"))
		}
		if r.URL.Query().Get("keywords") != "tesco" {
			t.Errorf("expected keywords=tesco, got %s", r.URL.Query().Get("keywords"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"bestMatches": [
				{
					"1. symbol": "TSCO.LON",
					"2. name": "Tesco PLC",
					"3. type": "Equity",
					"4. region": "United Kingdom",
					"8. currency": "GBX",
					"9. matchScore": "0.7273"
				},
				{
					"1. symbol": "TSCDF",
					"2. name": "Tesco Plc",
					"3. type": "Equity",
					"4. region": "United States",
					"8. currency": "USD",
					"9. matchScore": "0.7273"
				}
			]
		}`))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	matches, err := c.Search(context.Background(), "tesco")
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("got %d matches, want 2", len(matches))
	}
	if matches[0].Symbol != "TSCO.LON" {
		t.Errorf("matches[0].Symbol = %q, want TSCO.LON", matches[0].Symbol)
	}
	if matches[0].Name != "Tesco PLC" {
		t.Errorf("matches[0].Name = %q, want Tesco PLC", matches[0].Name)
	}
	if matches[0].Region != "United Kingdom" {
		t.Errorf("matches[0].Region = %q, want United Kingdom", matches[0].Region)
	}
	if matches[0].Currency != "GBX" {
		t.Errorf("matches[0].Currency = %q, want GBX", matches[0].Currency)
	}
}

func TestGetWTI(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("function") != "WTI" {
			t.Errorf("expected function=WTI, got %s", r.URL.Query().Get("function"))
		}
		if r.URL.Query().Get("interval") != "monthly" {
			t.Errorf("expected interval=monthly, got %s", r.URL.Query().Get("interval"))
		}
		payload := map[string]interface{}{
			"name":     "Crude Oil Prices WTI",
			"interval": "monthly",
			"unit":     "dollars per barrel",
			"data": []map[string]string{
				{"date": "2026-05-01", "value": "102.13"},
				{"date": "2026-04-01", "value": "98.50"},
				{"date": "2026-03-01", "value": "95.20"},
			},
		}
		b, _ := json.Marshal(payload)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	pts, err := c.GetWTI(context.Background(), "monthly", 2)
	if err != nil {
		t.Fatalf("GetWTI: %v", err)
	}
	if len(pts) != 2 {
		t.Fatalf("got %d points with limit 2, want 2", len(pts))
	}
	if pts[0].Date != "2026-05-01" {
		t.Errorf("pts[0].Date = %q, want 2026-05-01", pts[0].Date)
	}
	if pts[0].Value != "102.13" {
		t.Errorf("pts[0].Value = %q, want 102.13", pts[0].Value)
	}
	if pts[0].Name != "Crude Oil Prices WTI" {
		t.Errorf("pts[0].Name = %q, want Crude Oil Prices WTI", pts[0].Name)
	}
	if pts[0].Unit != "dollars per barrel" {
		t.Errorf("pts[0].Unit = %q, want dollars per barrel", pts[0].Unit)
	}
}

func TestGetBrent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("function") != "BRENT" {
			t.Errorf("expected function=BRENT, got %s", r.URL.Query().Get("function"))
		}
		payload := map[string]interface{}{
			"name":     "Crude Oil Prices Brent",
			"interval": "monthly",
			"unit":     "dollars per barrel",
			"data": []map[string]string{
				{"date": "2026-05-01", "value": "105.20"},
			},
		}
		b, _ := json.Marshal(payload)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	pts, err := c.GetBrent(context.Background(), "monthly", 10)
	if err != nil {
		t.Fatalf("GetBrent: %v", err)
	}
	if len(pts) != 1 {
		t.Fatalf("got %d points, want 1", len(pts))
	}
	if pts[0].Value != "105.20" {
		t.Errorf("Value = %q, want 105.20", pts[0].Value)
	}
}

func TestGetInflation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("function") != "INFLATION" {
			t.Errorf("expected function=INFLATION, got %s", r.URL.Query().Get("function"))
		}
		// Inflation has no interval param
		if r.URL.Query().Get("interval") != "" {
			t.Errorf("expected no interval param, got %s", r.URL.Query().Get("interval"))
		}
		payload := map[string]interface{}{
			"name":     "Inflation - US Consumer Prices",
			"interval": "annual",
			"unit":     "percent",
			"data": []map[string]string{
				{"date": "2024-01-01", "value": "2.94952520485207"},
				{"date": "2023-01-01", "value": "4.12448557483777"},
				{"date": "2022-01-01", "value": "8.00279473422739"},
			},
		}
		b, _ := json.Marshal(payload)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	pts, err := c.GetInflation(context.Background(), 2)
	if err != nil {
		t.Fatalf("GetInflation: %v", err)
	}
	if len(pts) != 2 {
		t.Fatalf("got %d points with limit 2, want 2", len(pts))
	}
	if pts[0].Date != "2024-01-01" {
		t.Errorf("pts[0].Date = %q, want 2024-01-01", pts[0].Date)
	}
	if pts[0].Unit != "percent" {
		t.Errorf("pts[0].Unit = %q, want percent", pts[0].Unit)
	}
}

func TestRetryOn503(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if hits < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"Global Quote": {
				"01. symbol": "IBM",
				"05. price": "272.24",
				"07. latest trading day": "2026-06-13",
				"08. previous close": "274.85",
				"09. change": "-2.61",
				"10. change percent": "-0.9496%"
			}
		}`))
	}))
	defer srv.Close()

	cfg := alphavantage.DefaultConfig()
	cfg.BaseURL = srv.URL
	cfg.APIKey = "TEST_KEY"
	cfg.Rate = 0
	cfg.Retries = 5
	cfg.Timeout = 5 * time.Second
	c := alphavantage.NewClient(cfg)

	_, err := c.GetQuote(context.Background(), "IBM")
	if err != nil {
		t.Fatalf("unexpected error after retry: %v", err)
	}
	if hits < 3 {
		t.Errorf("expected at least 3 hits, got %d", hits)
	}
}

func TestApikeyInURL(t *testing.T) {
	var gotAPIKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.URL.Query().Get("apikey")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bestMatches":[]}`))
	}))
	defer srv.Close()

	c := newTestClient(srv.URL)
	_, _ = c.Search(context.Background(), "test")
	if gotAPIKey != "TEST_KEY" {
		t.Errorf("apikey in URL = %q, want TEST_KEY", gotAPIKey)
	}
}
