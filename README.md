[简体中文](./README.md) | [English](./README.en.md)

# qqnotify-go

一个面向 AI 服务与自动化任务的 QQ 通知中间件。

`qqnotify-go` 用来帮助 Go 程序、Codex、Cursor、CI/CD、定时任务和内部工具，以较低心智负担把结果稳定发送到 QQ。

[![CI](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml/badge.svg)](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/wangyaxings/qqnotify-go/qqnotify.svg)](https://pkg.go.dev/github.com/wangyaxings/qqnotify-go/qqnotify)

## 安装

```bash
go get github.com/wangyaxings/qqnotify-go/qqnotify@latest
```

## Tell Codex

把下面这段话直接发给 Codex：

```text
Read ./docs/ai/quickstart.zh-CN.md and act as a step-by-step setup assistant for qqnotify-go.

Guide me through the full QQ setup flow:
1. Open the QQ bot platform pages
2. Log in by scanning the QR code
3. Find or create the bot
4. Obtain QQ_APP_ID
5. Generate QQ_APP_SECRET
6. Capture QQ_USER_OPENID by running go run ./cmd/qqnotify-openid
7. Save all three values into ./examples/smoke/.env.local
8. Run the smoke example and verify that QQ receives the message

Do not skip steps.
Do not assume values.
After each major step, ask me to confirm before continuing.
Only finish after I have confirmed all three values and the smoke example has sent a real QQ message successfully.
```

## Tell Cursor

把下面这段话直接发给 Cursor：

```text
Open this repository and follow ./docs/ai/quickstart.zh-CN.md as a strict setup checklist.

Your job is to guide me step by step until qqnotify-go successfully sends a real smoke-test message to QQ.
Make me obtain and confirm these three values first:
- QQ_APP_ID
- QQ_APP_SECRET
- QQ_USER_OPENID

Then store them in ./examples/smoke/.env.local and run the smoke example.
Do not skip steps, and pause for my confirmation after each milestone.
```

## Tell AI

如果你使用其他 AI，可以直接发这段话：

```text
Please act as a setup guide for qqnotify-go by following the repository document ./docs/ai/quickstart.zh-CN.md.

Guide me from opening the QQ bot platform and logging in, all the way to collecting QQ_APP_ID, QQ_APP_SECRET, and QQ_USER_OPENID, saving them into ./examples/smoke/.env.local, and running the smoke example successfully.

You must work step by step.
You must wait for my confirmation after each major step.
You must not finish until all three values are confirmed and QQ has received the smoke-test message.
```

## 最快启动

1. 复制示例配置文件

```powershell
Copy-Item ./examples/smoke/smoke.env.example ./examples/smoke/.env.local
```

2. 按照 [AI 快速启动文档](./docs/ai/quickstart.zh-CN.md) 获取并填写：

- `QQ_APP_ID`
- `QQ_APP_SECRET`
- `QQ_USER_OPENID`

3. 运行冒烟测试

```powershell
go run ./examples/smoke
```

如果成功：

- 终端会输出发送成功
- 你的 QQ 会收到一条 `qqnotify-go` 的测试消息

## 常用方式

Go 代码直接发送：

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

- [中文 AI 快速启动](./docs/ai/quickstart.zh-CN.md)
- [冒烟测试说明](./examples/smoke/README.md)
- [GitHub Actions 示例](./examples/github-actions/workflow.yml)
- [HTTP bridge 示例](./examples/http-bridge/README.md)
- [CHANGELOG](./CHANGELOG.md)
- [v0.1.0 发布说明](./docs/releases/v0.1.0.md)
