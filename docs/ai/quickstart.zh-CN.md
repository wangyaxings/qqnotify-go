# qqnotify-go AI 快速启动

这份文档的目标只有一个：

让你在 AI 的带领下，从 QQ 开放平台登录开始，逐步拿到下面 3 个关键参数，并把它们保存到专门用于冒烟测试的配置文件中，最后成功发送一条真实的 QQ 消息。

必需参数：

- `QQ_APP_ID`
- `QQ_APP_SECRET`
- `QQ_USER_OPENID`

专用配置文件：

- `./examples/smoke/.env.local`

只有当你已经：

1. 确认拿到这 3 个参数
2. 已经把它们保存到 `./examples/smoke/.env.local`
3. 已经运行 `go run ./examples/smoke`
4. 已经在 QQ 中确认收到消息

整个流程才算完成。

## 给 AI 的工作规则

如果你正在让 Codex、Cursor 或其他 AI 帮你配置，请要求它遵守下面的规则：

- 必须一步一步进行，不能跳步骤
- 每完成一个大步骤，都要停下来让你确认
- 不得假设任何参数已经存在
- 只有在你确认 3 个参数都拿到、配置文件已写好、QQ 已收到消息后，整个流程才算结束

## 第 1 步：打开平台并扫码登录

依次打开下面两个页面：

- 机器人列表页：[https://q.qq.com/qqbot/openclaw/](https://q.qq.com/qqbot/openclaw/)
- 开发设置页：[https://q.qq.com/qqbot/#/developer/developer-setting](https://q.qq.com/qqbot/#/developer/developer-setting)

如果页面要求登录，请使用手机 QQ 扫码登录。

完成后，确认：

- 你已经进入 QQ 机器人开放平台
- 你能看到机器人列表或创建入口

## 第 2 步：找到或创建机器人

在机器人列表页中：

1. 如果你已经有机器人，进入该机器人
2. 如果没有，就先创建一个机器人

完成后，你需要能够看到：

- 机器人名称
- 机器人 QQ 号
- 开发配置入口

## 第 3 步：获取 `QQ_APP_ID`

打开开发设置页，找到当前机器人的开发信息。

你需要记录：

- `AppID`

把它作为：

- `QQ_APP_ID`

## 第 4 步：生成或确认 `QQ_APP_SECRET`

仍然在开发设置页中，找到机器人密钥配置。

如果此前没有可用密钥，就重新生成一个新的 `AppSecret`。

注意：

- 重新生成后，旧密钥会失效
- 不要把这个值提交到仓库
- 这个值只应该保存在本地测试配置文件中

把它作为：

- `QQ_APP_SECRET`

## 第 5 步：为什么还需要 `QQ_USER_OPENID`

QQ 机器人发送单聊消息时，不能直接使用 QQ 号码，必须使用：

- `user_openid`

也就是：

- `QQ_USER_OPENID`

这个值不能通过 QQ 号直接反查，必须通过用户和机器人交互后的事件数据获取。

## 第 6 步：抓取 `QQ_USER_OPENID`

先复制冒烟测试配置模板：

```powershell
Copy-Item ./examples/smoke/smoke.env.example ./examples/smoke/.env.local
```

暂时先把你已经拿到的两个值填进去：

- `QQ_APP_ID`
- `QQ_APP_SECRET`

然后使用一个 AI 或你已有的本地工具，完成“抓取 openid”这一步。你可以要求 AI 帮你：
然后直接在当前仓库根目录执行：

```powershell
go run ./cmd/qqnotify-openid
```

这个命令会：

- 自动读取 `./examples/smoke/.env.local`
- 使用其中的 `QQ_APP_ID` 和 `QQ_APP_SECRET`
- 连接 QQ 机器人事件网关
- 等待你发送一条新的单聊消息
- 输出对应的 `QQ_USER_OPENID`

你给机器人发送一条新的单聊消息后，拿到对应的：

- `QQ_USER_OPENID`

## 第 7 步：把 3 个值写入专用冒烟配置文件

编辑：

- `./examples/smoke/.env.local`

写成下面这种形式：

```text
QQ_APP_ID=你的AppID
QQ_APP_SECRET=你的AppSecret
QQ_USER_OPENID=你的UserOpenID
```

这个文件只用于最小冒烟测试。

注意：

- 不要把真实值提交到 git
- `./examples/smoke/.env.local` 已被忽略

## 第 8 步：运行最小冒烟测试

在仓库根目录执行：

```powershell
go run ./examples/smoke
```

这个示例会自动读取：

- `./examples/smoke/.env.local`

如果成功，你会看到类似输出：

```text
smoke notification sent successfully using examples/smoke/.env.local
```

## 第 9 步：确认 QQ 已收到消息

检查你的 QQ 是否收到了来自机器人发送的测试消息。

如果收到了，说明下面这条链路已经打通：

1. QQ 平台配置正确
2. `QQ_APP_ID` 正确
3. `QQ_APP_SECRET` 正确
4. `QQ_USER_OPENID` 正确
5. `qqnotify-go` 可以成功发送真实消息

## 第 10 步：完成确认

只有当下面 4 条都满足时，才算这次 AI 启动流程真正完成：

- 你已经确认 `QQ_APP_ID`
- 你已经确认 `QQ_APP_SECRET`
- 你已经确认 `QQ_USER_OPENID`
- 你已经确认 QQ 收到了 `go run ./examples/smoke` 发出的测试消息

## 常见问题

### 1. 为什么还没拿到 `QQ_USER_OPENID`

常见原因：

- `go run ./cmd/qqnotify-openid` 还没有启动
- 你还没有真正给机器人发送一条新的单聊消息
- 发送的是旧消息，不是新事件

### 2. 为什么能登录平台但发不出消息

通常要检查：

- `QQ_APP_ID`
- `QQ_APP_SECRET`
- `QQ_USER_OPENID`

这三个值是否都是最新且互相对应的。

### 3. 为什么重新生成 `AppSecret` 后失效了

这是正常现象。重新生成之后，旧密钥会立即失效，你需要把新的值重新写入：

- `./examples/smoke/.env.local`

### 4. 最终用于冒烟测试的配置文件在哪里

固定使用：

- `./examples/smoke/.env.local`
