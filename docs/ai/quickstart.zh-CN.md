# qqnotify-go AI 快速启动

完整的安装与 AI 引导流程已经迁移到下面这个单一文档：

- [`.codex/INSTALL.md`](../../.codex/INSTALL.md)

如果你希望让 Codex、Claude 或其他 AI 直接带你完成配置，请直接使用：

```text
Fetch and follow instructions from https://raw.githubusercontent.com/wangyaxings/qqnotify-go/refs/heads/main/.codex/INSTALL.md
```

这份安装文档会覆盖完整流程，包括：

- 打开 QQ 开放平台并扫码登录
- 获取 `QQ_APP_ID`
- 获取 `QQ_APP_SECRET`
- 抓取 `QQ_USER_OPENID`
- 把 3 个值写入 `./examples/smoke/.env.local`
- 运行 `go run ./examples/smoke`
- 直到用户确认 QQ 收到真实测试消息才算完成
