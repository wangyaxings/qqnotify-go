# QQNotify-Go Repository Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将当前实验型 QQ 机器人项目重构为一个面向 AI 与自动化通知场景的标准 Go 中间件仓库 `qqnotify-go`。

**Architecture:** 采用标准 Go 项目布局，将执行入口拆分到 `cmd/`，将可复用能力下沉到 `internal/`，将运行产物、示例与文档分层管理。第一阶段只保留“发 QQ 通知”和最小 HTTP bridge 所需能力，移除与定位不一致的双向 worker 叙事，同时保留必要的底层发送逻辑与测试覆盖。

**Tech Stack:** Go 1.24, net/http, gorilla/websocket（暂存，仅在尚未完全移除接收能力时保留）, Git, GitHub

---

## Planned File Structure

- `cmd/qqnotifyd/main.go`
  - HTTP bridge 进程入口
- `cmd/example-send/main.go`
  - 最小发送示例入口
- `internal/qqapi/config.go`
  - QQ 应用配置与环境变量加载
- `internal/qqapi/client.go`
  - token 获取与消息发送
- `internal/notify/message.go`
  - 通知结构体与文本渲染
- `internal/notify/service.go`
  - 通知发送服务与最小业务编排
- `internal/httpbridge/handler.go`
  - HTTP bridge handler
- `internal/httpbridge/handler_test.go`
  - HTTP bridge 行为测试
- `internal/web/`（如仍保留演示主页）
  - 可选；若产品定位不需要则删除
- `examples/`
  - 面向 Codex / cron / GitHub Actions 的示例
- `.gitignore`
  - 忽略 `artifacts/`, `*.exe`, `.env` 等
- `README.md`
  - 改写为产品定位 README
- `docs/superpowers/specs/2026-04-05-qq-notify-middleware-design.md`
  - 已完成的定位 spec

## Task 1: Baseline Repository Hygiene

**Files:**
- Create: `.gitignore`
- Create: `LICENSE`（如用户未另行指定，先使用 MIT）
- Modify: `README.md`

- [ ] **Step 1: Write the failing test**

在当前代码中新增一个仓库卫生测试并不合适，因此本任务以“仓库约束文件建立”为主，不新增代码测试。

- [ ] **Step 2: Run current test suite to capture baseline**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 3: Add ignore and repository baseline files**

创建 `.gitignore`，至少忽略：

```gitignore
artifacts/
*.exe
*.dll
*.so
*.dylib
.env
.env.*
coverage.out
dist/
bin/
```

新增 `LICENSE`（默认 MIT）。

- [ ] **Step 4: Run tests again to verify no regression**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add .gitignore LICENSE README.md
git commit -m "chore: add repository baseline files"
```

## Task 2: Restructure Package Layout Into cmd/ and internal/

**Files:**
- Create: `cmd/qqnotifyd/main.go`
- Create: `cmd/example-send/main.go`
- Create: `internal/qqapi/config.go`
- Create: `internal/qqapi/client.go`
- Create: `internal/notify/message.go`
- Create: `internal/notify/service.go`
- Modify: existing `qqworker/*.go` by migrating or deleting responsibilities
- Test: migrated tests under `internal/.../*_test.go`

- [ ] **Step 1: Write the failing tests for new package boundaries**

新增最小测试：

```go
func TestLoadConfigFromEnvRequiresQQFields(t *testing.T) {}
func TestServiceSendTextDelegatesToClient(t *testing.T) {}
func TestRenderNotificationIncludesTitleAndBody(t *testing.T) {}
```

- [ ] **Step 2: Run targeted tests to verify they fail**

Run: `go test ./internal/...`
Expected: FAIL with missing packages / symbols

- [ ] **Step 3: Implement minimal new package layout**

迁移现有能力：

- `qqworker/notify.go` -> `internal/qqapi/config.go` + `internal/qqapi/client.go`
- 抽出通知结构与服务到 `internal/notify`
- 删除与 v1 定位无关的 queue / worker-subprocess / receive 入口依赖

- [ ] **Step 4: Run package tests**

Run: `go test ./internal/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/qqnotifyd/main.go cmd/example-send/main.go internal/qqapi internal/notify
git commit -m "refactor: adopt standard go package layout"
```

## Task 3: Add Minimal HTTP Bridge

**Files:**
- Create: `internal/httpbridge/handler.go`
- Create: `internal/httpbridge/handler_test.go`
- Modify: `cmd/qqnotifyd/main.go`

- [ ] **Step 1: Write the failing HTTP bridge tests**

```go
func TestHandlerAcceptsJSONAndSendsNotification(t *testing.T) {}
func TestHandlerRejectsInvalidPayload(t *testing.T) {}
```

- [ ] **Step 2: Run targeted tests to verify they fail**

Run: `go test ./internal/httpbridge -v`
Expected: FAIL with missing handler

- [ ] **Step 3: Implement minimal HTTP handler**

行为要求：

- `POST /notify`
- 接收 JSON payload
- 校验必填字段
- 调用 `notify.Service`
- 返回清晰状态码与 JSON

- [ ] **Step 4: Run tests**

Run: `go test ./internal/httpbridge -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/httpbridge/handler.go internal/httpbridge/handler_test.go cmd/qqnotifyd/main.go
git commit -m "feat: add minimal http notification bridge"
```

## Task 4: Rebuild README and Examples Around Product Positioning

**Files:**
- Modify: `README.md`
- Create: `examples/codex/main.go`
- Create: `examples/cron/main.go`
- Create: `examples/http-bridge/README.md`
- Create: `examples/github-actions/README.md`

- [ ] **Step 1: Write the failing validation step**

不新增自动化测试，改以人工检查 README 结构与示例命令为验证目标。

- [ ] **Step 2: Rewrite README**

README 必须包含：

- one-liner positioning
- 15 秒代码示例
- Quick Start
- 场景列表
- HTTP bridge 示例
- examples 导航

- [ ] **Step 3: Add examples**

最小示例代码或说明覆盖：

- Codex 场景
- Cron 场景
- GitHub Actions 场景
- HTTP bridge 场景

- [ ] **Step 4: Verify docs and examples build or read cleanly**

Run: `go test ./...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add README.md examples
git commit -m "docs: reposition project as qq notification middleware"
```

## Task 5: Final Verification and GitHub Publish Preparation

**Files:**
- Modify: `go.mod` / `go.sum` if dependency cleanup is needed
- Modify: repository root files as needed

- [ ] **Step 1: Run full verification**

Run: `go test ./...`
Expected: PASS

Run: `go build ./cmd/...`
Expected: PASS

- [ ] **Step 2: Remove obsolete files and dependencies**

清理：

- 无定位价值的旧入口
- 旧 `app.exe`
- 不再需要的 worker / queue / capture 代码
- 无用依赖

- [ ] **Step 3: Re-run verification**

Run: `go test ./...`
Expected: PASS

Run: `go build ./cmd/...`
Expected: PASS

- [ ] **Step 4: Initialize git and prepare first publish**

Run:

```bash
git init
git branch -M main
git status
```

如果远端仓库已存在：

```bash
git remote add origin https://github.com/wangyaxings/qqnotify-go.git
git push -u origin main
```

- [ ] **Step 5: Commit**

```bash
git add README.md cmd internal examples go.mod go.sum .gitignore LICENSE
git commit -m "feat: bootstrap qq notification middleware"
```
