---
title: "Resource URIs"
description: "Use alphavantage as a database/sql-style driver so a host program can address it as alphavantage:// URIs."
weight: 10
---

`alphavantage` is a command line, but the `alphavantage` Go package is also a
small driver that makes Alpha Vantage data addressable as a resource URI. A host
program registers it the way a program registers a database driver with
`database/sql`, then dereferences `alphavantage://` URIs without knowing
anything about how the data is fetched.

The host that does this today is [ant](https://github.com/tamnd/ant).

## Mounting the driver

A host enables the driver with one blank import:

```go
import _ "github.com/tamnd/alphavantage-cli/alphavantage"
```

The package's `init` registers a domain with the scheme `alphavantage` for the
host `www.alphavantage.co`. The standalone `alphavantage` binary does not change.

## Addressing records

| URI                                  | What it is                       |
| ------------------------------------ | -------------------------------- |
| `alphavantage://symbol/IBM`          | real-time stock quote for IBM    |
| `alphavantage://query/tesco`         | symbol search for "tesco"        |

```bash
ant get alphavantage://symbol/IBM
ant url alphavantage://symbol/IBM
ant resolve https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=IBM&apikey=demo
```

## Classify and locate

`Classify` turns any input into a (type, id) pair:

- Uppercase 1-5 letter ticker (e.g. `IBM`, `AAPL`) becomes `("symbol", "IBM")`
- Anything else becomes `("query", input)` for search

`Locate` is the inverse, returning the live API URL for a (type, id).
