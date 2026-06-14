---
title: "CLI reference"
description: "Every alphavantage command and flag."
weight: 10
---

## Global flags

These flags apply to every command:

| Flag | Default | Description |
|------|---------|-------------|
| `-o`, `--output` | `auto` | Output format: `table`, `json`, `jsonl`, `csv`, `tsv`, `url`, `raw` |
| `--fields` | all | Comma-separated list of fields to include |
| `--template` | | Go template for each record |
| `-n`, `--limit` | 0 (all) | Maximum number of records |
| `--no-color` | | Disable color output |
| `--key` | `demo` | Alpha Vantage API key |

## Commands

### quote

Fetch a real-time stock quote.

```bash
alphavantage quote <symbol> [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--key` | `demo` | Alpha Vantage API key |

**Example:**

```bash
alphavantage quote IBM
alphavantage quote AAPL --key YOUR_KEY
```

**Output fields:** `symbol`, `price`, `open`, `high`, `low`, `volume`,
`change`, `change_percent`, `previous_close`, `latest_day`

### search

Search for stock symbols by keywords.

```bash
alphavantage search <keywords> [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--key` | `demo` | Alpha Vantage API key |

**Example:**

```bash
alphavantage search tesco
alphavantage search "microsoft" --key YOUR_KEY
```

**Output fields:** `symbol`, `name`, `type`, `region`, `currency`, `match_score`

### wti

Fetch WTI crude oil price data.

```bash
alphavantage wti [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `monthly` | Interval: `daily`, `weekly`, `monthly` |
| `--limit` | 10 | Maximum number of data points |
| `--key` | `demo` | Alpha Vantage API key |

**Example:**

```bash
alphavantage wti
alphavantage wti --interval daily --limit 30
```

**Output fields:** `date`, `value`, `name`, `unit`

### brent

Fetch Brent crude oil price data.

```bash
alphavantage brent [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `monthly` | Interval: `daily`, `weekly`, `monthly` |
| `--limit` | 10 | Maximum number of data points |
| `--key` | `demo` | Alpha Vantage API key |

**Example:**

```bash
alphavantage brent
alphavantage brent --interval weekly --limit 52
```

**Output fields:** `date`, `value`, `name`, `unit`

### inflation

Fetch US consumer price inflation data.

```bash
alphavantage inflation [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | 10 | Maximum number of data points |
| `--key` | `demo` | Alpha Vantage API key |

**Example:**

```bash
alphavantage inflation
alphavantage inflation --limit 5
```

**Output fields:** `date`, `value`, `name`, `unit`

### serve

Run an HTTP server exposing every command as a route.

```bash
alphavantage serve [--addr :7777]
```

### mcp

Run an MCP server exposing every command as a tool.

```bash
alphavantage mcp
```

### version

Print version information.

```bash
alphavantage version [--short]
```
