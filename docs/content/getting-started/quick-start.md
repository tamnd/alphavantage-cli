---
title: "Quick start"
description: "Fetch your first record with alphavantage."
weight: 30
---

Once `alphavantage` is on your `PATH`, fetch a real-time stock quote:

```bash
alphavantage quote IBM
```

The `demo` key works for IBM. You get an aligned table by default. Ask for
JSON when you want to pipe it:

```bash
alphavantage quote IBM -o json
```

```json
[
  {
    "symbol": "IBM",
    "price": "272.2400",
    "open": "274.8500",
    "high": "275.5700",
    "low": "271.3700",
    "volume": "4014849",
    "change": "-2.6100",
    "change_percent": "-0.9496%",
    "previous_close": "274.8500",
    "latest_day": "2026-06-13"
  }
]
```

## Search for symbols

```bash
alphavantage search tesco
```

## Crude oil prices

```bash
alphavantage wti                        # WTI, last 10 monthly points
alphavantage wti --interval daily       # daily prices
alphavantage brent --limit 5            # Brent, last 5 points
```

## Inflation data

```bash
alphavantage inflation                  # last 10 annual data points
alphavantage inflation --limit 3        # just the 3 most recent years
```

## Use your own API key

```bash
alphavantage quote AAPL --key YOUR_KEY
alphavantage search microsoft --key YOUR_KEY
```

Register a free key at
[alphavantage.co](https://www.alphavantage.co/support/#api-key).

## Shape the output

The same flags work on every command:

```bash
alphavantage quote IBM --fields symbol,price,change_percent
alphavantage wti -o jsonl | jq .value
alphavantage inflation -o csv
```

`-o` takes `table`, `json`, `jsonl`, `csv`, `tsv`, `url`, or `raw`. Left to
`auto`, it prints a table to a terminal and JSONL into a pipe.
