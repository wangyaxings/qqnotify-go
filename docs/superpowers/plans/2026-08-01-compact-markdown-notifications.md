# Compact Markdown Notifications Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Render concise Chinese Markdown notifications and send them through QQ without internal scheduler metadata.

**Architecture:** Notification templates build user-facing Chinese title/body values. The common renderer creates a small Markdown document and omits transport metadata. The QQ client adds a dedicated Markdown request while preserving `SendText` compatibility.

**Tech Stack:** Go 1.24, net/http, QQ C2C message API, standard Go tests.

## Global Constraints

- Do not expose credentials, source IDs, trace IDs, Cron expressions, or internal paths.
- Do not silently report a Markdown send as successful when QQ rejects it.
- Keep the existing text sender API available for callers and test doubles.
- PanSou must pin an immutable qqnotify-go commit.

---

### Task 1: Define compact Markdown rendering

**Files:**
- Modify: `qqnotify/message.go`
- Test: `qqnotify/message_test.go`

**Interfaces:**
- Produces: `RenderNotification(Notification) string` with `### title`, a blank line, and body only.

- [ ] **Step 1: Replace loose containment tests with an exact Markdown output assertion and explicit absence checks for metadata.**
- [ ] **Step 2: Run the focused test and observe the old metadata renderer fail.**
- [ ] **Step 3: Implement the minimal title/body Markdown renderer.**
- [ ] **Step 4: Run the focused test and confirm it passes.**

### Task 2: Localize built-in templates

**Files:**
- Modify: `qqnotify/templates.go`
- Test: `qqnotify/templates_test.go`
- Test: `internal/httpbridge/handler_test.go`

**Interfaces:**
- Produces: Chinese Codex, CI, and Cron notifications; Cron body contains summary only.

- [ ] **Step 1: Add exact expectations that Cron output excludes `Cron`, `Job`, `Schedule`, and the Cron expression.**
- [ ] **Step 2: Run template and bridge tests and observe old English expectations fail.**
- [ ] **Step 3: Implement concise Chinese titles and labels, omitting internal fields.**
- [ ] **Step 4: Run template and bridge tests and confirm they pass.**

### Task 3: Send native QQ Markdown

**Files:**
- Modify: `qqnotify/client.go`
- Test: `qqnotify/client_test.go`
- Modify: `internal/httpbridge/handler.go`
- Test: `internal/httpbridge/handler_test.go`

**Interfaces:**
- Produces: `SendMarkdown(context.Context, string) error` using `msg_type=2` and `markdown.content`.
- Preserves: `SendText(context.Context, string) error` using `msg_type=0`.

- [ ] **Step 1: Add a client test that captures and asserts the Markdown request JSON.**
- [ ] **Step 2: Run the client test and observe `SendMarkdown` is missing.**
- [ ] **Step 3: Add the Markdown request model and sender method; make the HTTP bridge use it when supported by the sender.**
- [ ] **Step 4: Run all Go tests and confirm they pass.**

### Task 4: Publish and pin the gateway

**Files:**
- Modify: `baidu-subscriptions/deploy/z480/Dockerfile.qqnotifyd`

**Interfaces:**
- Consumes: immutable qqnotify-go commit from Tasks 1-3.
- Produces: deployed `qqnotifyd` built from the new GitHub commit.

- [ ] **Step 1: Commit and push `agent/compact-markdown-notifications` to `wangyaxings/qqnotify-go`.**
- [ ] **Step 2: Update both clone checkout and verification hashes in `Dockerfile.qqnotifyd`.**
- [ ] **Step 3: Build and restart `qqnotifyd`, then generate and inspect a representative PanSou Markdown payload.**
- [ ] **Step 4: Commit and push the PanSou deployment pin, then close both issues with evidence.**

