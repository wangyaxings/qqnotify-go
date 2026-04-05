[简体中文](./README.md) | [English](./README.en.md)

# qqnotify-go

A QQ notification middleware for AI and automation workflows.

`qqnotify-go` helps Go apps, Codex, Cursor, CI/CD jobs, cron tasks, and internal tools send results to QQ with minimal setup.

[![CI](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml/badge.svg)](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/wangyaxings/qqnotify-go/qqnotify.svg)](https://pkg.go.dev/github.com/wangyaxings/qqnotify-go/qqnotify)

## Installation

```bash
go get github.com/wangyaxings/qqnotify-go/qqnotify@latest
```

## Tell Codex

```text
Read ./docs/ai/quickstart.zh-CN.md and guide me step by step until qqnotify-go sends a real smoke-test message to QQ.

Help me obtain and confirm:
- QQ_APP_ID
- QQ_APP_SECRET
- QQ_USER_OPENID

Capture QQ_USER_OPENID by running:
go run ./cmd/qqnotify-openid

Then save all three values into ./examples/smoke/.env.local and run:
go run ./examples/smoke

Do not skip steps. Pause for confirmation after each milestone.
```

## Quick Start

1. Copy the smoke config template:

```bash
cp ./examples/smoke/smoke.env.example ./examples/smoke/.env.local
```

2. Follow the Chinese setup guide to obtain and fill:

- `QQ_APP_ID`
- `QQ_APP_SECRET`
- `QQ_USER_OPENID`

Guide:

- [docs/ai/quickstart.zh-CN.md](./docs/ai/quickstart.zh-CN.md)

3. Run the smoke example:

```bash
go run ./examples/smoke
```

## Common Usage

- Go client: `qqnotify.NewClient(...)`
- Reusable templates: Codex / CI / cron
- HTTP bridge: `go run ./cmd/qqnotifyd`

## Docs

- [Chinese quickstart guide](./docs/ai/quickstart.zh-CN.md)
- [Smoke example guide](./examples/smoke/README.md)
- [GitHub Actions example](./examples/github-actions/workflow.yml)
- [HTTP bridge example](./examples/http-bridge/README.md)
- [CHANGELOG](./CHANGELOG.md)
