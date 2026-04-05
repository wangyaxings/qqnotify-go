# GitHub Actions Example

Use `qqnotify-go` to forward CI results to QQ after a workflow completes.

Files:

- `main.go`: minimal Go notifier example
- `workflow.yml`: copyable GitHub Actions workflow example

## Suggested Flow

1. Store `QQ_APP_ID`, `QQ_APP_SECRET`, and `QQ_USER_OPENID` as repository secrets.
2. Run a small Go helper or call `qqnotifyd` after your build and test jobs.
3. Send a success or failure notification with the workflow name and run URL.

## Included Workflow Example

See [workflow.yml](./workflow.yml) for a complete example that:

- checks out the repository
- runs Go tests
- sends a success notification on green runs
- sends a failure notification with workflow context on failed runs

## Example Payload

```json
{
  "title": "GitHub Actions failed",
  "body": "The release workflow failed on job build-linux.",
  "status": "error",
  "source": "github-actions",
  "trace_id": "run-123456"
}
```
