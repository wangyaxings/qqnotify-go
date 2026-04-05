# Smoke Example

This is the fastest path to verify that `qqnotify-go` can send a real QQ message.

## Config File

Create the dedicated smoke-test config file:

```powershell
Copy-Item ./examples/smoke/smoke.env.example ./examples/smoke/.env.local
```

Fill in:

- `QQ_APP_ID`
- `QQ_APP_SECRET`
- `QQ_USER_OPENID`

## Get QQ_USER_OPENID

If you already have `QQ_APP_ID` and `QQ_APP_SECRET`, you can capture `QQ_USER_OPENID` with:

```powershell
go run ./cmd/qqnotify-openid
```

Then send a fresh direct message to your bot from the QQ account you want to bind.

## Run the Smoke Example

```powershell
go run ./examples/smoke
```

If successful:

- the terminal prints a success message
- QQ receives a real test notification
