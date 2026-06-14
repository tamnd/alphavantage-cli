# alphavantage

A command line for the [Alpha Vantage](https://www.alphavantage.co) financial data API.

Stock quotes, symbol search, crude oil prices, and US inflation data. One pure-Go binary, output that pipes into the rest of your tools.

```bash
alphavantage quote IBM           # real-time stock quote
alphavantage search tesco        # search for symbols
alphavantage wti                 # WTI crude oil prices
alphavantage brent --limit 5    # Brent crude, last 5 monthly points
alphavantage inflation           # US consumer price inflation
```

The `demo` API key works out of the box for IBM, tesco search, WTI, BRENT, and
INFLATION. Register a free key at
[alphavantage.co](https://www.alphavantage.co/support/#api-key) and pass it as
`--key YOUR_KEY` to use any symbol.

## Installation

```bash
# Homebrew
brew install tamnd/tap/alphavantage

# Go
go install github.com/tamnd/alphavantage-cli/cmd/alphavantage@latest

# Container
docker run --rm ghcr.io/tamnd/alphavantage quote IBM
```

Or download a release binary from the
[releases page](https://github.com/tamnd/alphavantage-cli/releases).

## Usage

```
alphavantage quote <symbol>           real-time stock quote
alphavantage search <keywords>        symbol search
alphavantage wti [--interval daily|weekly|monthly] [--limit N]
alphavantage brent [--interval daily|weekly|monthly] [--limit N]
alphavantage inflation [--limit N]

Global flags:
  --key string     Alpha Vantage API key (default: demo)
  -o string        Output format: table, json, jsonl, csv, tsv, url, raw
  --fields string  Comma-separated list of fields to include
  --limit int      Maximum number of records
```

## Output

Default output adapts to context: an aligned table on the terminal, JSONL into
a pipe.

```bash
alphavantage quote IBM -o json
alphavantage wti -o jsonl | jq .value
alphavantage inflation --fields date,value -o csv
```

## License

Apache 2.0. See [LICENSE](LICENSE).
