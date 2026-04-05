# Changelog

All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog, and this project follows Semantic Versioning.

## [0.1.0] - 2026-04-05

### Added

- Initial public `qqnotify` Go client for sending QQ notifications
- Config loading from environment variables for QQ bot credentials
- Configurable client options for retry attempts and HTTP timeout
- Reusable notification templates for Codex, CI, and cron workflows
- `qqnotifyd` HTTP bridge for non-Go integrations
- Optional Bearer token auth for `/notify`
- Health endpoint at `GET /healthz`
- Template-aware bridge payloads with support for `codex`, `ci`, and `cron`
- GitHub Actions CI workflow
- Example programs for Codex, cron, GitHub Actions, and HTTP bridge usage

### Changed

- Refined repository positioning around AI and automation notifications to QQ
- Improved README with installation guidance, versioning notes, and example matrix
- Added bridge payload validation for required template fields

### Security

- Added optional HTTP bridge Bearer token protection for `/notify`
