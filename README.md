# qqnotify-go

A production-ready Go middleware for sending AI and automation notifications to QQ.

`qqnotify-go` is built for developers who want the fastest way to deliver results from Codex, AI agents, cron jobs, CI/CD pipelines, and internal tools to QQ.

[![CI](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml/badge.svg)](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/wangyaxings/qqnotify-go/qqnotify.svg)](https://pkg.go.dev/github.com/wangyaxings/qqnotify-go/qqnotify)

## Why qqnotify-go

- Minimal setup to send the first message
- Clean Go API for apps and services
- Lightweight HTTP bridge for non-Go callers
- Production-friendly defaults for timeout and transport handling
- Retries transient upstream failures when sending messages
- Designed for AI and automation scenarios instead of generic bot framework complexity

## Quick Example

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/wangyaxings/qqnotify-go/qqnotify"
)

func main() {
    cfg, err := qqnotify.LoadConfigFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    client := qqnotify.NewClientWithOptions(cfg, nil, qqnotify.ClientOptions{
        RetryAttempts: 3,
        Timeout:       20 * time.Second,
    })
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    err = client.Send(ctx, qqnotify.NewCodexNotification(qqnotify.CodexTemplate{
        Task:    "Refactor notification bridge",
        Summary: "All tests passed and the patch is ready.",
        Status:  "success",
        TraceID: "job-123",
        Files:   []string{"internal/httpbridge/handler.go", "README.md"},
    }))
    if err != nil {
        log.Fatal(err)
    }
}
```

## Installation

```bash
go get github.com/wangyaxings/qqnotify-go/qqnotify@latest
```

## Supported Use Cases

- Codex task completion notifications
- AI agent execution updates
- Cron job reports
- GitHub Actions or CI failure alerts
- Internal tools and operational notifications

## Environment Variables

Required:

- `QQ_APP_ID`
- `QQ_APP_SECRET`
- `QQ_USER_OPENID`

Optional:

- `QQ_BOT_TOKEN_BASE_URL`
- `QQ_BOT_API_BASE_URL`
- `QQNOTIFY_LISTEN_ADDR` for `qqnotifyd` listen address
- `QQNOTIFY_AUTH_TOKEN` for optional Bearer token auth on `/notify`
- `PORT` as a compatibility fallback for `qqnotifyd`

## Client Options

Use `qqnotify.NewClientWithOptions` when you want to override default behavior:

- `RetryAttempts`: number of send retries for transient 5xx upstream failures
- `Timeout`: default HTTP client timeout used when `nil` client is provided

Defaults:

- `RetryAttempts = 2`
- `Timeout = 10s`

## Quick Start

Send a first message from Go:

```powershell
$env:QQ_APP_ID="your-app-id"
$env:QQ_APP_SECRET="your-app-secret"
$env:QQ_USER_OPENID="your-user-openid"
go run ./cmd/example-send
```

Start the HTTP bridge:

```powershell
$env:QQ_APP_ID="your-app-id"
$env:QQ_APP_SECRET="your-app-secret"
$env:QQ_USER_OPENID="your-user-openid"
$env:QQNOTIFY_LISTEN_ADDR=":8080"
$env:QQNOTIFY_AUTH_TOKEN="your-bridge-token"
go run ./cmd/qqnotifyd
```

Then call it:

```powershell
Invoke-RestMethod -Method Post -Uri http://127.0.0.1:8080/notify `
  -Headers @{ Authorization = "Bearer your-bridge-token" } `
  -ContentType 'application/json' `
  -Body '{"title":"Build finished","body":"CI completed successfully","status":"success"}'
```

Or with `curl`:

```bash
curl -X POST http://127.0.0.1:8080/notify \
  -H "Authorization: Bearer your-bridge-token" \
  -H "Content-Type: application/json" \
  -d '{"title":"Build finished","body":"CI completed successfully","status":"success"}'
```

Template-aware bridge payloads are also supported:

```bash
curl -X POST http://127.0.0.1:8080/notify \
  -H "Authorization: Bearer your-bridge-token" \
  -H "Content-Type: application/json" \
  -d '{"type":"codex","task":"Refactor bridge auth","summary":"All tests passed.","status":"success","files":["internal/httpbridge/handler.go","README.md"]}'
```

Health check:

```powershell
Invoke-RestMethod -Method Get -Uri http://127.0.0.1:8080/healthz
```

## Reusable Templates

`qqnotify-go` includes reusable templates for common automation workflows:

- `qqnotify.NewCodexNotification`
- `qqnotify.NewCINotification`
- `qqnotify.NewCronNotification`

This lets callers avoid handcrafting titles and multiline bodies for the most common notification types.

The HTTP bridge supports the same common templates through JSON payloads:

- `type: "codex"`
- `type: "ci"`
- `type: "cron"`

## Project Layout

```text
cmd/qqnotifyd          HTTP bridge entrypoint
cmd/example-send       Minimal send demo
qqnotify/              Public library API
internal/httpbridge/   HTTP bridge internals
examples/              Scenario examples
docs/                  Specs and plans
```

## Examples

| Scenario | Files |
| --- | --- |
| Codex task result | [examples/codex/main.go](./examples/codex/main.go) |
| Cron / scheduled job | [examples/cron/main.go](./examples/cron/main.go) |
| GitHub Actions / CI | [examples/github-actions/main.go](./examples/github-actions/main.go), [examples/github-actions/README.md](./examples/github-actions/README.md) |
| HTTP bridge | [examples/http-bridge/README.md](./examples/http-bridge/README.md) |

## Versioning

The repository follows semantic versioning.

- `v0.x`: fast-moving pre-1.0 releases while the API is being shaped
- `v1.x`: stable public API with backward-compatible minor releases

The first public milestone for the repository is `v0.1.0`.

## HTTP Bridge Deployment Notes

- `GET /healthz` is always open for liveness probes
- `POST /notify` can be protected with `QQNOTIFY_AUTH_TOKEN`
- `QQNOTIFY_LISTEN_ADDR` is the preferred listen setting
- `PORT` remains supported as a fallback for simple hosting environments

## Scope

Current focus:

- QQ notification sending
- AI and automation delivery workflows
- Reusable middleware API

Not in v1:

- Full bot framework features
- Two-way chat flows
- Command routing
- Session management
- Plugin system
