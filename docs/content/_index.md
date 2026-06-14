---
title: "alphavantage"
description: "A command line for the Alpha Vantage financial data API."
heroTitle: "Alpha Vantage, from the command line"
heroLead: "A command line for the Alpha Vantage financial data API. One pure-Go binary, stock quotes, symbol search, crude oil prices, and US inflation data that pipe into the rest of your tools."
heroPrimaryURL: "/getting-started/quick-start/"
heroPrimaryText: "Get started"
---

`alphavantage` reads financial data from Alpha Vantage over plain HTTPS, shapes
it into clean records, and gets out of your way.

```bash
alphavantage quote IBM               # real-time stock quote
alphavantage search tesco            # search for symbols
alphavantage wti                     # WTI crude oil prices
alphavantage brent --interval daily  # Brent crude, daily
alphavantage inflation --limit 5     # US inflation, last 5 years
```

The `demo` API key works for IBM (quote), tesco (search), WTI, BRENT, and
INFLATION. Pass `--key YOUR_KEY` to use your own key from
[api.nasa.gov](https://www.alphavantage.co/support/#api-key) for any symbol.

Output adapts to where it goes: an aligned table on your terminal, JSONL the
moment you pipe it somewhere.

## Two ways to use it

- **As a command** for reading financial data by hand or in a script. Start with
  the [quick start](/getting-started/quick-start/).
- **As a resource-URI driver** so a host like
  [ant](https://github.com/tamnd/ant) can address Alpha Vantage as
  `alphavantage://` URIs and follow links across sites. See
  [resource URIs](/guides/resource-uris/).

Both are the same code: one operation, declared once, is a CLI command, an HTTP
route, an MCP tool, and a URI dereference.

## Where to go next

- New here? Read the [introduction](/getting-started/introduction/), then the
  [quick start](/getting-started/quick-start/).
- Installing? See [installation](/getting-started/installation/).
- Doing a specific job? The [guides](/guides/) are task-first.
- Need every flag? The [CLI reference](/reference/cli/) is the full surface.
