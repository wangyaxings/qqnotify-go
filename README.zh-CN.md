# qqnotify-go

一个面向 AI 服务与自动化任务的 QQ 通知中间件。

`qqnotify-go` 的目标是让 Go 程序、Codex、AI Agent、定时任务、CI/CD 和内部工具，以较低心智负担将结果稳定发送到 QQ。

[![CI](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml/badge.svg)](https://github.com/wangyaxings/qqnotify-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/wangyaxings/qqnotify-go/qqnotify.svg)](https://pkg.go.dev/github.com/wangyaxings/qqnotify-go/qqnotify)

## 为什么用它

- 配置简单，能尽快发出第一条 QQ 通知
- 提供清晰的 Go API
- 提供给非 Go 系统使用的 `qqnotifyd` HTTP bridge
- 提供 Codex、CI、cron 三类常见模板
- 支持 bridge 鉴权和健康检查

## 安装

```bash
go get github.com/wangyaxings/qqnotify-go/qqnotify@latest
```

## 快速示例

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

## 环境变量

必填：

- `QQ_APP_ID`
- `QQ_APP_SECRET`
- `QQ_USER_OPENID`

可选：

- `QQ_BOT_TOKEN_BASE_URL`
- `QQ_BOT_API_BASE_URL`
- `QQNOTIFY_LISTEN_ADDR`
- `QQNOTIFY_AUTH_TOKEN`
- `PORT`

## 模板能力

内置模板：

- `qqnotify.NewCodexNotification`
- `qqnotify.NewCINotification`
- `qqnotify.NewCronNotification`

HTTP bridge 也支持模板化 payload：

- `type: "codex"`，必填 `task`
- `type: "ci"`，必填 `workflow`
- `type: "cron"`，必填 `name`

## HTTP Bridge

启动：

```powershell
$env:QQ_APP_ID="your-app-id"
$env:QQ_APP_SECRET="your-app-secret"
$env:QQ_USER_OPENID="your-user-openid"
$env:QQNOTIFY_LISTEN_ADDR=":8080"
$env:QQNOTIFY_AUTH_TOKEN="your-bridge-token"
go run ./cmd/qqnotifyd
```

调用：

```bash
curl -X POST http://127.0.0.1:8080/notify \
  -H "Authorization: Bearer your-bridge-token" \
  -H "Content-Type: application/json" \
  -d '{"type":"codex","task":"Refactor bridge auth","summary":"All tests passed.","status":"success","files":["internal/httpbridge/handler.go","README.md"]}'
```

探活：

```bash
curl http://127.0.0.1:8080/healthz
```

## 示例

- [Codex 示例](./examples/codex/main.go)
- [Cron 示例](./examples/cron/main.go)
- [GitHub Actions 示例](./examples/github-actions/workflow.yml)
- [HTTP bridge 示例](./examples/http-bridge/README.md)

## 版本说明

- 当前仓库采用语义化版本
- `v0.x` 表示 API 仍在持续打磨
- 第一个正式公开版本是 [`v0.1.0`](./docs/releases/v0.1.0.md)

更多版本信息见：

- [CHANGELOG.md](./CHANGELOG.md)
- [v0.1.0 release notes](./docs/releases/v0.1.0.md)
