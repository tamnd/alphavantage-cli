---
title: "Introduction"
description: "How alphavantage is put together and why."
weight: 10
---

`alphavantage` is a single Go binary. It reads financial data from the Alpha
Vantage API, shapes it into typed records, and prints output that pipes into
the rest of your tools. There is nothing to install alongside it and no daemon
to run.

## The data it covers

Alpha Vantage provides a wide range of financial data. The `demo` key covers a
focused set that works without registration:

- **Quotes** - real-time stock quote for IBM
- **Search** - symbol search (demo works for "tesco")
- **WTI** - West Texas Intermediate crude oil prices
- **BRENT** - Brent crude oil prices
- **Inflation** - US consumer price inflation

Register a free API key at [alphavantage.co](https://www.alphavantage.co/support/#api-key)
to unlock the full symbol universe.

## One shape, many surfaces

Every command is declared once as a kit operation. The same declaration drives:

- a CLI subcommand
- an HTTP route (`alphavantage serve`)
- an MCP tool (`alphavantage mcp`)
- a URI dereference when used as a driver in a host like `ant`

There is no second implementation to keep in step.

## Output

The default output format depends on where it goes:

- **terminal** - an aligned table, coloured if the terminal supports it
- **pipe** - one JSON object per line (JSONL), ready for `jq` or any processor

Override with `-o json`, `-o csv`, `-o tsv`, `-o url`, or `-o raw`.

## API key

The `demo` key is the default. It works for IBM, tesco search, and the
economic indicators. Pass `--key YOUR_KEY` on any command to use your own key
for any symbol.
