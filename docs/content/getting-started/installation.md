---
title: "Installation"
description: "Get the alphavantage binary onto your PATH."
weight: 20
---

## Homebrew (macOS and Linux)

```bash
brew install tamnd/tap/alphavantage
```

## Download a release binary

Grab the archive for your platform from the
[releases page](https://github.com/tamnd/alphavantage-cli/releases), unpack it,
and drop `alphavantage` somewhere on your `PATH`.

```bash
# macOS arm64
curl -L https://github.com/tamnd/alphavantage-cli/releases/latest/download/alphavantage_latest_darwin_arm64.tar.gz | tar xz
sudo mv alphavantage /usr/local/bin/
```

## Linux packages

Each release ships a `.deb`, `.rpm`, and `.apk`. Download from the
[releases page](https://github.com/tamnd/alphavantage-cli/releases) and install
with your package manager.

## Container image

```bash
docker run --rm ghcr.io/tamnd/alphavantage quote IBM
```

## Build from source

```bash
go install github.com/tamnd/alphavantage-cli/cmd/alphavantage@latest
```

Go 1.22 or later is required.
