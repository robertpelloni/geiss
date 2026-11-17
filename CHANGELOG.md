## [1.0.40] - 2026-11-08
### Added
- **Automated Database Backups (`jules.db`)**: Added `/api/system/backup` GET and POST endpoints in `backend-go/backup.go` with timestamped snapshots in `backups/`.
- **Interactive Backup Card**: Added Database Backups card in `src/App.jsx` with instant backup creation, snapshot history listing, and live size telemetry.
- **Unit Tests**: Added `TestDatabaseBackupHandler` in `backend-go/main_test.go`.

## [1.0.39] - 2026-11-08
### Added
- **System Tray Controller & Graceful Shutdown**: Added `/api/system/tray` (telemetry, uptime, platform health) and `/api/system/shutdown` (safe exit for daemons and Go server) in `backend-go/systray.go`.
- **Wired Backend Handlers**: Exposed `/api/system/refactor`, `/api/system/pipeline`, `/api/system/tray`, and `/api/system/shutdown` in `backend-go/main.go`.
- **Unified Single-Page Dashboard**: Completely consolidated all workspace features into a single high-priority dashboard in `src/App.jsx`. Features comprehensive interactive form controls, tooltips, real-time telemetry, model routing, mass search/replace, and pipeline dispatching without subpage fragmentation.
- **Unit Tests**: Added test cases for system tray telemetry and shutdown endpoints in `backend-go/main_test.go`.

## [1.0.38] - 2026-06-30
### Fixed
- Merged pending submodules and CI changes.

## 1.0.43 (2026-11-16) — synthesis pass [muse-wonder-1.0 in Muse Code in C:/Users/hyper/workspace/bobsaver on laptop]
- Fail-closed auth (401 when JULES_API_KEY unset/wrong) on all mutating backend routes; telemetry target_dir traversal rejected.
- Fixed UI contract mismatches (reasoning, drift status/hashes, conflict file/blocks/strategy, audit status/file counts); removed dead tab state; single-page dashboard with backend PID/uptime polling.
- Honest automation: CI auto-fix is simulation-only (501); pipeline deploy is simulation-only except local dry-run and reports missing build/tests as failures.
- Unified shutdown/tray on one handler and one start time; shadow report mutex; drift fetch timeout + parallel scan; logs labelled commit history.
- Backup tests use temp dirs; state-snapshot (no-DB) test added; legacy db.go.bak archived to docs/archive/.
- Version unified at 1.0.43 (root VERSION.md is source of truth); bun is canonical frontend tool; added geiss/MANUAL.md and scripts/build.ps1.
