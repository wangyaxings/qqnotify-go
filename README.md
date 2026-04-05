[简体中文](./README.md) | [English](./README.en.md)

# qqnotify-go

用 Go、AI 工具和自动化任务发送 QQ 通知。

`qqnotify-go` 可以把程序结果和任务状态直接发到 QQ。

[![CI](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml/badge.svg)](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/wangyaxings/qqnotify-go/qqnotify.svg)](https://pkg.go.dev/github.com/wangyaxings/qqnotify-go/qqnotify)

## 安装

让 Codex、Claude 或其他 AI 直接读取安装文档并按步骤执行：

```text
Fetch and follow instructions from https://raw.githubusercontent.com/wangyaxings/qqnotify-go/refs/heads/main/.codex/INSTALL.md
```

本地安装文档也在：

- [`.codex/INSTALL.md`](./.codex/INSTALL.md)

## 使用

```powershell
$env:QQ_APP_ID="your-app-id"
$env:QQ_APP_SECRET="your-app-secret"
$env:QQ_USER_OPENID="your-user-openid"
go run ./examples/smoke
```

Go 代码中直接发送：

```go
client := qqnotify.NewClient(cfg, nil)
err := client.Send(ctx, qqnotify.NewCodexNotification(qqnotify.CodexTemplate{
    Task:    "Refactor notification bridge",
    Summary: "All tests passed.",
    Status:  "success",
}))
```

HTTP bridge：

```powershell
$env:QQNOTIFY_LISTEN_ADDR=":8080"
$env:QQNOTIFY_AUTH_TOKEN="your-bridge-token"
go run ./cmd/qqnotifyd
```

## 文档

- [完整安装文档](./.codex/INSTALL.md)
- [冒烟测试说明](./examples/smoke/README.md)
- [GitHub Actions 示例](./examples/github-actions/workflow.yml)
- [HTTP bridge 示例](./examples/http-bridge/README.md)
- [CHANGELOG](./CHANGELOG.md)
- [v0.1.0 发布说明](./docs/releases/v0.1.0.md)
