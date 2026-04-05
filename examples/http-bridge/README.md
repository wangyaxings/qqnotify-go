# HTTP Bridge Example

Start the bridge:

```powershell
$env:QQ_APP_ID="your-app-id"
$env:QQ_APP_SECRET="your-app-secret"
$env:QQ_USER_OPENID="your-user-openid"
$env:QQNOTIFY_LISTEN_ADDR=":8080"
$env:QQNOTIFY_AUTH_TOKEN="your-bridge-token"
go run ./cmd/qqnotifyd
```

Call the bridge:

```powershell
Invoke-RestMethod -Method Post -Uri http://127.0.0.1:8080/notify `
  -Headers @{ Authorization = "Bearer your-bridge-token" } `
  -ContentType 'application/json' `
  -Body '{"title":"Deployment finished","body":"Version v1.0.0 is live","status":"success"}'
```

With curl:

```bash
curl -X POST http://127.0.0.1:8080/notify \
  -H "Authorization: Bearer your-bridge-token" \
  -H "Content-Type: application/json" \
  -d '{"title":"Deployment finished","body":"Version v1.0.0 is live","status":"success"}'
```

Health check:

```bash
curl http://127.0.0.1:8080/healthz
```

Template-aware payloads:

```bash
curl -X POST http://127.0.0.1:8080/notify \
  -H "Authorization: Bearer your-bridge-token" \
  -H "Content-Type: application/json" \
  -d '{"type":"codex","task":"Refactor bridge auth","summary":"All tests passed.","status":"success","files":["internal/httpbridge/handler.go","README.md"]}'
```

Supported template types:

- `codex`
- `ci`
- `cron`

Required fields by template:

- `codex`: `task`
- `ci`: `workflow`
- `cron`: `name`

Supported fields:

- `title`
- `body`
- `status`
- `source`
- `trace_id`
- `timestamp`
