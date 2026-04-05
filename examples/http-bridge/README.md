# HTTP Bridge Example

Start the bridge:

```powershell
$env:QQ_APP_ID="your-app-id"
$env:QQ_APP_SECRET="your-app-secret"
$env:QQ_USER_OPENID="your-user-openid"
go run ./cmd/qqnotifyd
```

Call the bridge:

```powershell
Invoke-RestMethod -Method Post -Uri http://127.0.0.1:8080/notify `
  -ContentType 'application/json' `
  -Body '{"title":"Deployment finished","body":"Version v1.0.0 is live","status":"success"}'
```

Supported fields:

- `title`
- `body`
- `status`
- `source`
- `trace_id`
- `timestamp`
