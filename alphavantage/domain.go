// Package alphavantage exposes the Alpha Vantage financial data API as a kit
// Domain: a driver that a multi-domain host (ant) enables with a single blank
// import,
//
//	import _ "github.com/tamnd/alphavantage-cli/alphavantage"
//
// exactly as a database/sql program enables a driver with `import _
// "github.com/lib/pq"`. The init below registers it; the host then dereferences
// alphavantage:// URIs by routing to the operations Register installs. The same
// Domain also builds the standalone alphavantage binary (see cli.NewApp), so the
// binary and a host share one source of truth.
package alphavantage

import (
	"context"
	"regexp"
	"strings"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

func init() { kit.Register(Domain{}) }

// Domain is the Alpha Vantage driver.
type Domain struct{}

// Info describes the scheme, the hostnames a pasted link is matched against, and
// the identity reused for the binary's help and version.
func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme: "alphavantage",
		Hosts:  []string{Host},
		Identity: kit.Identity{
			Binary: "alphavantage",
			Short:  "CLI for Alpha Vantage financial data API",
			Long: `CLI for Alpha Vantage financial data API

alphavantage reads financial data — stock quotes, symbol search, crude oil
prices, and US inflation data — over plain HTTPS from www.alphavantage.co,
shapes it into clean records, and prints output that pipes into the rest of
your tools.

The "demo" API key works for IBM (quote), tesco (search), and economic
indicators (WTI, BRENT, INFLATION). Register a free key at
https://www.alphavantage.co/support/#api-key for full access.`,
			Site: Host,
			Repo: "https://github.com/tamnd/alphavantage-cli",
		},
	}
}

// Register installs the client factory and every operation onto app.
func (Domain) Register(app *kit.App) {
	app.SetClient(newClient)

	kit.Handle(app, kit.OpMeta{
		Name:     "quote",
		Group:    "read",
		Single:   true,
		Summary:  "Fetch a real-time stock quote",
		URIType:  "symbol",
		Resolver: true,
		Args:     []kit.Arg{{Name: "symbol", Help: "stock symbol (e.g. IBM)"}},
	}, getQuote)

	kit.Handle(app, kit.OpMeta{
		Name:    "search",
		Group:   "read",
		List:    true,
		Summary: "Search for stock symbols",
		Args:    []kit.Arg{{Name: "keywords", Help: "keywords to search (e.g. tesco)"}},
	}, doSearch)

	kit.Handle(app, kit.OpMeta{
		Name:    "wti",
		Group:   "commodities",
		List:    true,
		Summary: "WTI crude oil price data",
	}, getWTI)

	kit.Handle(app, kit.OpMeta{
		Name:    "brent",
		Group:   "commodities",
		List:    true,
		Summary: "Brent crude oil price data",
	}, getBrent)

	kit.Handle(app, kit.OpMeta{
		Name:    "inflation",
		Group:   "economic",
		List:    true,
		Summary: "US consumer price inflation data",
	}, getInflation)
}

// newClient builds the client from the host-resolved config.
func newClient(_ context.Context, cfg kit.Config) (any, error) {
	c := DefaultConfig()
	if cfg.Rate > 0 {
		c.Rate = cfg.Rate
	}
	if cfg.Retries > 0 {
		c.Retries = cfg.Retries
	}
	if cfg.Timeout > 0 {
		c.Timeout = cfg.Timeout
	}
	return NewClient(c), nil
}

// --- input structs ---

type quoteInput struct {
	Symbol string  `kit:"arg" help:"stock symbol (e.g. IBM)"`
	Key    string  `kit:"flag" name:"key" help:"Alpha Vantage API key (default: demo)"`
	Client *Client `kit:"inject"`
}

type searchInput struct {
	Keywords string  `kit:"arg" help:"keywords to search (e.g. tesco)"`
	Key      string  `kit:"flag" name:"key" help:"Alpha Vantage API key (default: demo)"`
	Client   *Client `kit:"inject"`
}

type wtiInput struct {
	Interval string  `kit:"flag" help:"interval: daily, weekly, monthly" default:"monthly"`
	Limit    int     `kit:"flag,inherit" help:"max data points" default:"10"`
	Key      string  `kit:"flag" name:"key" help:"Alpha Vantage API key (default: demo)"`
	Client   *Client `kit:"inject"`
}

type brentInput struct {
	Interval string  `kit:"flag" help:"interval: daily, weekly, monthly" default:"monthly"`
	Limit    int     `kit:"flag,inherit" help:"max data points" default:"10"`
	Key      string  `kit:"flag" name:"key" help:"Alpha Vantage API key (default: demo)"`
	Client   *Client `kit:"inject"`
}

type inflationInput struct {
	Limit  int     `kit:"flag,inherit" help:"max data points" default:"10"`
	Key    string  `kit:"flag" name:"key" help:"Alpha Vantage API key (default: demo)"`
	Client *Client `kit:"inject"`
}

// --- handlers ---

func getQuote(ctx context.Context, in quoteInput, emit func(*Quote) error) error {
	applyKey(in.Client, in.Key)
	q, err := in.Client.GetQuote(ctx, in.Symbol)
	if err != nil {
		return err
	}
	return emit(q)
}

func doSearch(ctx context.Context, in searchInput, emit func(*Match) error) error {
	applyKey(in.Client, in.Key)
	matches, err := in.Client.Search(ctx, in.Keywords)
	if err != nil {
		return err
	}
	for i := range matches {
		if err := emit(&matches[i]); err != nil {
			return err
		}
	}
	return nil
}

func getWTI(ctx context.Context, in wtiInput, emit func(*DataPoint) error) error {
	applyKey(in.Client, in.Key)
	interval := in.Interval
	if interval == "" {
		interval = "monthly"
	}
	pts, err := in.Client.GetWTI(ctx, interval, in.Limit)
	if err != nil {
		return err
	}
	for i := range pts {
		if err := emit(&pts[i]); err != nil {
			return err
		}
	}
	return nil
}

func getBrent(ctx context.Context, in brentInput, emit func(*DataPoint) error) error {
	applyKey(in.Client, in.Key)
	interval := in.Interval
	if interval == "" {
		interval = "monthly"
	}
	pts, err := in.Client.GetBrent(ctx, interval, in.Limit)
	if err != nil {
		return err
	}
	for i := range pts {
		if err := emit(&pts[i]); err != nil {
			return err
		}
	}
	return nil
}

func getInflation(ctx context.Context, in inflationInput, emit func(*DataPoint) error) error {
	applyKey(in.Client, in.Key)
	pts, err := in.Client.GetInflation(ctx, in.Limit)
	if err != nil {
		return err
	}
	for i := range pts {
		if err := emit(&pts[i]); err != nil {
			return err
		}
	}
	return nil
}

// applyKey updates the client's API key from the --key flag if provided.
func applyKey(c *Client, key string) {
	if key != "" {
		c.cfg.APIKey = key
	}
}

// --- Resolver (URI driver) ---

// tickerRe matches an uppercase ticker: 1-5 letters only.
var tickerRe = regexp.MustCompile(`^[A-Z]{1,5}$`)

// Classify turns any accepted input into the canonical (type, id).
// An uppercase 1-5-letter ticker maps to ("symbol", input).
// Anything else maps to ("query", input) for search.
func (Domain) Classify(input string) (uriType, id string, err error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", "", errs.Usage("empty alphavantage reference")
	}
	if tickerRe.MatchString(input) {
		return "symbol", input, nil
	}
	return "query", input, nil
}

// Locate is the inverse: the live https URL for a (type, id).
func (Domain) Locate(uriType, id string) (string, error) {
	switch uriType {
	case "symbol":
		return "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + id + "&apikey=demo", nil
	case "query":
		return "https://www.alphavantage.co/query?function=SYMBOL_SEARCH&keywords=" + id + "&apikey=demo", nil
	default:
		return "", errs.Usage("alphavantage has no resource type %q", uriType)
	}
}
