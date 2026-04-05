# qqnotify-go 安装与 AI 引导说明

这份文档是 `qqnotify-go` 的唯一安装指令源。

如果你是 Codex、Claude 或其他 AI 编程助手，请严格按下面的协议执行，并以“尽量减少交互轮次、引导用户完成真实 QQ 冒烟测试”为唯一目标。

## 总目标

你必须帮助用户最终完成下面 3 件事：

1. 拿到 `QQ_APP_ID`
2. 拿到 `QQ_APP_SECRET`
3. 拿到 `QQ_USER_OPENID`

并且：

- 将这 3 个值写入 `./examples/smoke/.env.local`
- 成功运行 `go run ./examples/smoke`
- 用户确认 QQ 已收到真实测试消息

只有全部满足，流程才算完成。

---

## 关键规则

你必须遵守下面这些规则：

- 以“阶段”为单位推进，不要按细碎步骤频繁停顿
- 每一轮只推动一个明确动作，不要同时布置多个目标
- 能由 AI 自己验证的内容，不要反复要求用户确认
- 不得假设用户已经拿到任何 QQ 参数
- 不要把真实密钥写入受版本管理的文件
- 最终只允许把真实值写入 `./examples/smoke/.env.local`
- 除最终“QQ 是否收到消息”外，其他内容优先由 AI 自己检查
- 如果文档描述与代码真实行为不一致，以代码真实行为为准，并明确说明阻塞点

---

## 必需参数

- `QQ_APP_ID`
- `QQ_APP_SECRET`
- `QQ_USER_OPENID`

## 唯一允许写入真实值的文件

- `./examples/smoke/.env.local`

---

## AI 输出格式规则

每一轮尽量使用下面的固定结构：

1. 当前阶段结论
2. 现在请用户执行的唯一动作
3. 用户完成后只需要回复的短句

避免输出多个并列任务。避免让用户做你自己可以完成的事。

例如，不要说：

- 请确认已经打开平台、进入机器人、看到开发设置并准备复制参数

而应该说：

- 现在请发送机器人配置截图，或直接发送机器人名称、AppID、AppSecret。完成后无需解释，我会自动处理。

---

# 阶段 A：工作区就绪

## 阶段目标

确认当前目录就是 `qqnotify-go` 仓库根目录，并且具备执行冒烟测试所需的最小文件结构。

## AI 必须先检查

至少检查下面这些路径或等价标志是否存在：

- `go.mod`
- `./cmd/qqnotify-openid`
- `./examples/smoke`

## 如果检查失败

不要继续进入 QQ 平台步骤。你必须先修正工作区，只允许给出下面两种选择之一：

1. 请用户切换到 `qqnotify-go` 仓库根目录
2. 在用户允许的前提下，帮用户克隆仓库到当前目录

## 进入下一阶段的条件

下面条件全部成立后，才能进入阶段 B：

- 当前目录是 `qqnotify-go` 仓库根目录
- `./examples/smoke` 存在
- `./cmd/qqnotify-openid` 存在

---

# 阶段 B：一次性收集机器人信息

## 阶段目标

一次性拿到下面 3 项：

- 机器人名称
- `QQ_APP_ID`
- `QQ_APP_SECRET`

## 本阶段的正式输入方式

用户可以提供下面任一输入：

1. 目标机器人的配置截图
2. 文字形式的：
   - 机器人名称
   - `AppID`
   - `AppSecret`

截图输入是正式支持的主流程，不是可选备注。

## AI 在本阶段开始前必须先提醒

你必须明确告诉用户：

- 你可以直接把 `AppSecret` 发给 AI，我会自动帮你配置本地冒烟测试
- 为了降低风险，建议你在本次测试完成后，立即重新生成新的 `AppSecret`
- 当前机器人和当前这组参数仅建议用于这次真实冒烟测试
- 我只会把这些值写入 `./examples/smoke/.env.local`，不会写入受版本控制文件

## AI 推荐输出模板

你应优先使用类似下面的话术：

```text
现在进入参数收集阶段。

你可以直接把下面任一内容发给我，我会自动帮你完成本地测试配置：
- 目标机器人的配置截图
- 或者直接发送：机器人名称、AppID、AppSecret

说明：
- 你可以直接把 AppSecret 发给 AI，让我自动配置
- 为了降低风险，建议你在本次测试完成后，立即重新生成新的 AppSecret
- 当前机器人和这组参数仅建议用于本次真实冒烟测试
- 我只会把这些值写入 ./examples/smoke/.env.local，不会写入受版本控制文件
```

## AI 处理截图的规则

如果用户发送的是截图，你必须先尝试从截图中提取：

- 机器人名称
- `QQ_APP_ID`
- `QQ_APP_SECRET`

如果缺少任一项，不要让用户全部重发。你只能精确补问缺少的字段。

例如：

- 还缺少机器人名称，请直接发机器人名称文字，或补一张包含机器人名称的截图
- 还缺少 `AppSecret`，请直接把 `AppSecret` 发给我

## 本阶段禁止事项

- 不要把获取 `AppID` 和 `AppSecret` 拆成两个阶段
- 不要在拿到 `QQ_APP_ID` 后先停一次
- 不要在拿到 `QQ_APP_SECRET` 后再停一次
- 不要优先要求用户自己手工编辑配置文件

## 进入下一阶段的条件

下面条件全部成立后，才能进入阶段 C：

- 已拿到机器人名称
- 已拿到 `QQ_APP_ID`
- 已拿到 `QQ_APP_SECRET`

---

# 阶段 C：自动写入本地冒烟配置

## 阶段目标

创建并写入本次冒烟测试专用配置文件：

- `./examples/smoke/.env.local`

## AI 必须做的事

优先自己完成下面动作，而不是让用户手工编辑：

1. 如果 `./examples/smoke/.env.local` 不存在，则创建它
2. 写入：
   - `QQ_APP_ID`
   - `QQ_APP_SECRET`
3. 预留空的 `QQ_USER_OPENID`

最终内容应类似：

```text
QQ_APP_ID=你的AppID
QQ_APP_SECRET=你的AppSecret
QQ_USER_OPENID=
```

## AI 必须提醒

- `./examples/smoke/.env.local` 是本次流程唯一使用的冒烟测试配置文件
- 后面拿到 `QQ_USER_OPENID` 后，也必须继续写回这个文件
- 真实值不要提交到 git
- 如果用户已经把 `AppSecret` 发到了聊天中，在本阶段结束时再次提醒：
  - 本次测试完成后建议立即重新生成新的 `AppSecret`

## 进入下一阶段的条件

下面条件全部成立后，才能进入阶段 D：

- `./examples/smoke/.env.local` 已存在
- `QQ_APP_ID` 已写入
- `QQ_APP_SECRET` 已写入
- `QQ_USER_OPENID` 位置已预留

---

# 阶段 D：抓取 QQ_USER_OPENID

## 阶段目标

启动 openid 捕获程序，监听新的单聊消息事件，并拿到真实的 `QQ_USER_OPENID`。

## AI 在启动前必须说明

你必须向用户明确说明：

- QQ 机器人发送单聊消息时，不能直接使用 QQ 号
- 必须使用用户事件里的 `user_openid`
- 我现在会启动监听程序
- 只有“监听启动之后的一条新单聊消息”才能被捕获

## AI 启动前必须检查

- `./examples/smoke/.env.local` 已存在
- `QQ_APP_ID` 已写入
- `QQ_APP_SECRET` 已写入

## 启动命令

默认使用：

```powershell
go run ./cmd/qqnotify-openid
```

## AI 必须验证的不是“命令启动”，而是“程序正在监听”

只有当输出中出现明确的监听日志时，才能要求用户发消息。类似下面的日志可视为成功进入监听：

- `listening for a fresh QQ direct message`
- `send a new message to the bot`

如果程序没有进入等待状态，就不要让用户发消息。

## 如果程序秒退或报错

优先检查下面问题：

1. `QQ_APP_ID` 或 `QQ_APP_SECRET` 是否无效
2. 当前目录是否不是仓库根目录
3. 代码实现是否把 `QQ_USER_OPENID` 也错误地当成了启动前必填项

如果真实代码当前要求 `QQ_USER_OPENID` 非空，导致监听在启动前退出，则你可以：

- 仅在本地临时写入一个占位值到 `./examples/smoke/.env.local`
- 只用于让监听程序真正启动
- 一旦捕获到真实 `QQ_USER_OPENID`，必须立刻覆盖写回真实值

严禁把占位值当作最终配置保留。

## 监听成功后，AI 对用户的提示必须足够明确

你应使用类似下面的话术：

```text
我已经启动了 openid 抓取监听，并确认它正在等待新的单聊消息。

现在请你做这 2 步：
1. 用你想接收通知的那个 QQ 账号，给机器人发送一条新的单聊消息
2. 发送完成后，只回复我：已发送

注意：
- 必须是监听启动之后的新消息
- 历史消息或更早发送的消息不会被捕获
```

## 用户回复“已发送”后

AI 必须自行读取监听输出，并尝试提取：

- `QQ_USER_OPENID`
- 最近一条消息内容（如果程序有输出）

一旦拿到真实 `QQ_USER_OPENID`，必须立刻写回：

- `./examples/smoke/.env.local`

最终内容应类似：

```text
QQ_APP_ID=你的AppID
QQ_APP_SECRET=你的AppSecret
QQ_USER_OPENID=你的UserOpenID
```

## 进入下一阶段的条件

下面条件全部成立后，才能进入阶段 E：

- 已拿到真实 `QQ_USER_OPENID`
- `QQ_USER_OPENID` 已写回 `./examples/smoke/.env.local`

---

# 阶段 E：运行真实冒烟测试并完成判定

## 阶段目标

使用 `./examples/smoke/.env.local` 发送真实 QQ 通知，并确认 QQ 已收到消息。

## 运行命令

```powershell
go run ./examples/smoke
```

## AI 必须做的事

- 自行运行命令
- 自行读取终端结果
- 如果命令失败，优先根据报错继续排查
- 如果命令成功，不要再要求用户重复确认配置文件内容

## 成功后的用户确认

只有在命令成功后，你才需要让用户做最后一次人工确认：

- 请打开 QQ，确认是否已经收到刚才的真实测试消息

## 最终完成条件

只有下面 3 条都满足时，流程才算真正完成：

1. `./examples/smoke/.env.local` 中已经包含：
   - `QQ_APP_ID`
   - `QQ_APP_SECRET`
   - `QQ_USER_OPENID`
2. `go run ./examples/smoke` 已成功执行
3. 用户明确回复：QQ 已收到真实测试消息

如果任意一条未满足，就不要结束流程，而是继续停留在当前阶段协助完成。

---

# 常见问题

## 1. 为什么一直拿不到 `QQ_USER_OPENID`

优先检查：

- `go run ./cmd/qqnotify-openid` 是否真的进入监听状态
- `./examples/smoke/.env.local` 中的 `QQ_APP_ID` 和 `QQ_APP_SECRET` 是否正确
- 用户发送的是否是“监听启动之后的新单聊消息”
- 代码是否错误要求 `QQ_USER_OPENID` 也作为监听启动前必填项

## 2. 为什么平台能登录，但冒烟测试发不出消息

优先检查：

- `QQ_APP_ID`
- `QQ_APP_SECRET`
- `QQ_USER_OPENID`

这三个值是否都是最新、有效、并且彼此对应。

## 3. 为什么重新生成 `AppSecret` 后旧配置失效

这是正常现象。重新生成后，旧密钥会立即失效。需要把新的值重新写回：

- `./examples/smoke/.env.local`

## 4. 这次流程最终只使用哪个配置文件

固定使用：

- `./examples/smoke/.env.local`
