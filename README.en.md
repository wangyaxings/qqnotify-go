[简体中文](./README.md) | [English](./README.en.md)

# qqnotify-go

Send QQ notifications from Go, AI tools, and automation jobs.

`qqnotify-go` sends app results and job status updates to QQ.

[![CI](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml/badge.svg)](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/wangyaxings/qqnotify-go/qqnotify.svg)](https://pkg.go.dev/github.com/wangyaxings/qqnotify-go/qqnotify)

## Installation

Ask your AI assistant to fetch and follow the install guide:

```text
Fetch and follow instructions from https://raw.githubusercontent.com/wangyaxings/qqnotify-go/refs/heads/main/.codex/INSTALL.md
```

The same guide is available in this repository:

- [`.codex/INSTALL.md`](./.codex/INSTALL.md)

## Usage

```powershell
$env:QQ_APP_ID="your-app-id"
$env:QQ_APP_SECRET="your-app-secret"
$env:QQ_USER_OPENID="your-user-openid"
go run ./examples/smoke
```

Use the Go client directly:

```go
client := qqnotify.NewClient(cfg, nil)
err := client.Send(ctx, qqnotify.NewCodexNotification(qqnotify.CodexTemplate{
    Task:    "Refactor notification bridge",
    Summary: "All tests passed.",
    Status:  "success",
}))
```

## Docs

- [Full install guide](./.codex/INSTALL.md)
- [Smoke example guide](./examples/smoke/README.md)
- [GitHub Actions example](./examples/github-actions/workflow.yml)
- [HTTP bridge example](./examples/http-bridge/README.md)
- [CHANGELOG](./CHANGELOG.md)
