# qqnotify-go 安装与 AI 引导说明

这份文档是 `qqnotify-go` 的唯一安装指令源。

如果你是 Codex、Claude 或其他 AI 编程助手，请严格按下面的要求执行，并以“逐步引导用户完成真实冒烟测试”为唯一目标。

## 你的工作规则

你必须遵守下面这些规则：

- 一次只推进一个明确步骤，不能跳步骤
- 每完成一个大步骤，都要停下来让用户确认
- 不得假设用户已经拿到任何 QQ 参数
- 不要把真实密钥写进仓库受版本管理的文件
- 必须把最终得到的 3 个参数写入专用冒烟测试配置文件 `./examples/smoke/.env.local`
- 只有当用户确认已经拿到全部 3 个参数，并且 `go run ./examples/smoke` 已经成功把消息发送到 QQ，这次流程才算完成

必需参数：

- `QQ_APP_ID`
- `QQ_APP_SECRET`
- `QQ_USER_OPENID`

专用配置文件：

- `./examples/smoke/.env.local`

最终完成条件：

1. 用户确认已经拿到 `QQ_APP_ID`
2. 用户确认已经拿到 `QQ_APP_SECRET`
3. 用户确认已经拿到 `QQ_USER_OPENID`
4. 用户确认这 3 个值已经写入 `./examples/smoke/.env.local`
5. 用户确认运行 `go run ./examples/smoke` 后，QQ 已收到真实测试消息

## 第 1 步：打开 QQ 开放平台并登录

让用户依次打开下面两个页面：

- 机器人列表页：[https://q.qq.com/qqbot/openclaw/](https://q.qq.com/qqbot/openclaw/)
- 开发设置页：[https://q.qq.com/qqbot/#/developer/developer-setting](https://q.qq.com/qqbot/#/developer/developer-setting)

如果页面要求登录，让用户使用手机 QQ 扫码登录。

这一阶段完成后，你要让用户确认：

- 已经进入 QQ 机器人开放平台
- 能看到机器人列表、已有机器人，或创建机器人入口

## 第 2 步：找到或创建机器人

指导用户在机器人列表页完成下面其中一条：

1. 如果已经有机器人，就进入目标机器人
2. 如果还没有，就先创建一个机器人

完成后，让用户确认自己已经能看到：

- 机器人名称
- 机器人 QQ 号
- 开发配置入口

## 第 3 步：获取 `QQ_APP_ID`

让用户打开开发设置页，找到当前机器人的开发信息，并记录：

- `AppID`

把这个值记作：

- `QQ_APP_ID`

拿到后必须停下来，让用户确认 `QQ_APP_ID` 已经拿到。

## 第 4 步：生成或确认 `QQ_APP_SECRET`

继续在开发设置页中找到机器人密钥配置。

指导用户：

- 如果已有可用密钥，就确认当前要使用的值
- 如果没有可用密钥，就生成新的 `AppSecret`

必须提醒用户：

- 重新生成后，旧密钥会失效
- 真实值不要提交到 git
- 这个值只应保存到本地冒烟测试配置文件中

把这个值记作：

- `QQ_APP_SECRET`

拿到后必须停下来，让用户确认 `QQ_APP_SECRET` 已经拿到。

## 第 5 步：准备专用冒烟测试配置文件

指导用户在仓库根目录执行：

```powershell
Copy-Item ./examples/smoke/smoke.env.example ./examples/smoke/.env.local
```

然后让用户先把已经拿到的两个值写进去：

```text
QQ_APP_ID=你的AppID
QQ_APP_SECRET=你的AppSecret
QQ_USER_OPENID=
```

必须明确告诉用户：

- `./examples/smoke/.env.local` 是本次流程唯一使用的冒烟测试配置文件
- 后面拿到 `QQ_USER_OPENID` 后，也要继续写回这个文件

## 第 6 步：解释为什么必须拿到 `QQ_USER_OPENID`

向用户说明：

- QQ 机器人发送单聊消息时，不能直接使用 QQ 号
- 必须使用对应用户事件里的 `user_openid`
- 也就是 `QQ_USER_OPENID`
- 这个值通常需要通过用户与机器人交互后的事件数据获取

然后再进入抓取步骤。

## 第 7 步：抓取 `QQ_USER_OPENID`

指导用户在仓库根目录执行：

```powershell
go run ./cmd/qqnotify-openid
```

在执行前，你要再次确认：

- `./examples/smoke/.env.local` 已经存在
- `QQ_APP_ID` 和 `QQ_APP_SECRET` 已经写入其中

向用户说明这个命令会：

- 自动读取 `./examples/smoke/.env.local`
- 使用其中的 `QQ_APP_ID` 和 `QQ_APP_SECRET`
- 连接 QQ 机器人事件网关
- 等待一条新的单聊消息事件
- 在捕获成功后输出 `QQ_USER_OPENID`

接下来指导用户：

1. 保持命令运行
2. 用想要绑定的 QQ 账号给机器人发送一条新的单聊消息
3. 回到终端查看输出的 `QQ_USER_OPENID`

拿到后必须停下来，让用户确认 `QQ_USER_OPENID` 已经拿到。

## 第 8 步：把 3 个值完整写入 `./examples/smoke/.env.local`

指导用户编辑：

- `./examples/smoke/.env.local`

最终内容应类似：

```text
QQ_APP_ID=你的AppID
QQ_APP_SECRET=你的AppSecret
QQ_USER_OPENID=你的UserOpenID
```

必须提醒用户：

- 这个文件只用于最小冒烟测试
- 不要把真实值提交到 git
- 当前仓库已经忽略 `./examples/smoke/.env.local`

写完后必须停下来，让用户确认配置文件已经准备好。

## 第 9 步：运行最小冒烟测试

指导用户在仓库根目录执行：

```powershell
go run ./examples/smoke
```

向用户说明：

- 这个示例会自动读取 `./examples/smoke/.env.local`
- 如果配置正确，终端会输出发送成功

你需要等待用户反馈执行结果。

## 第 10 步：确认 QQ 已收到真实消息

让用户打开 QQ，确认是否已经收到测试消息。

只有当用户明确回复已经收到消息时，你才能认为下面这条链路已经打通：

1. QQ 平台配置正确
2. `QQ_APP_ID` 正确
3. `QQ_APP_SECRET` 正确
4. `QQ_USER_OPENID` 正确
5. `qqnotify-go` 可以成功发送真实 QQ 通知

## 第 11 步：完成判定

只有当下面 5 条都满足时，这次流程才算真正完成：

- 用户确认已经拿到 `QQ_APP_ID`
- 用户确认已经拿到 `QQ_APP_SECRET`
- 用户确认已经拿到 `QQ_USER_OPENID`
- 用户确认这 3 个值已经写入 `./examples/smoke/.env.local`
- 用户确认 QQ 已收到 `go run ./examples/smoke` 发出的真实测试消息

如果任意一条没有确认，就不要结束流程，而是继续停留在当前步骤协助用户完成。

## 常见问题

### 1. 为什么一直拿不到 `QQ_USER_OPENID`

常见原因：

- `go run ./cmd/qqnotify-openid` 还没有启动成功
- `./examples/smoke/.env.local` 里的 `QQ_APP_ID` 或 `QQ_APP_SECRET` 不正确
- 用户没有发送一条新的单聊消息
- 发送的是旧会话里的历史消息，而不是新事件

### 2. 为什么平台能登录，但冒烟测试发不出消息

优先检查：

- `QQ_APP_ID`
- `QQ_APP_SECRET`
- `QQ_USER_OPENID`

这三个值是否都是最新、有效、并且彼此对应。

### 3. 为什么重新生成 `AppSecret` 后旧配置失效

这是正常现象。重新生成后，旧密钥会立即失效。需要把新的值重新写回：

- `./examples/smoke/.env.local`

### 4. 这次流程最终使用哪个配置文件

固定使用：

- `./examples/smoke/.env.local`
