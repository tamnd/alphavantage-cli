---
title: "Output formats"
description: "The output formats alphavantage supports and when to use each."
weight: 20
---

Every command accepts `-o <format>`. The default is `auto`.

| Format | When to use |
|--------|-------------|
| `auto` | Table to a terminal, JSONL into a pipe |
| `table` | Force the aligned table |
| `json` | Pretty-printed JSON array |
| `jsonl` | One JSON object per line, ready for `jq` |
| `csv` | Comma-separated values with a header row |
| `tsv` | Tab-separated values with a header row |
| `url` | One URL per line (for resolver types only) |
| `raw` | The first field of each record, one per line |

## Filtering fields

`--fields symbol,price,change_percent` keeps only those columns in any format.

## Go templates

`--template '{{.Price}}'` applies a Go template to each record. Field names
match the JSON keys with an initial capital letter:

```bash
alphavantage quote IBM --template '{{.Symbol}}: {{.Price}} ({{.ChangePercent}})'
```
