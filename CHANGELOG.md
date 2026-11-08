## [1.0.39] - 2026-11-08
### Added
- **System Tray Controller & Graceful Shutdown**: Added `/api/system/tray` (telemetry, uptime, platform health) and `/api/system/shutdown` (safe exit for daemons and Go server) in `backend-go/systray.go`.
- **Wired Backend Handlers**: Exposed `/api/system/refactor`, `/api/system/pipeline`, `/api/system/tray`, and `/api/system/shutdown` in `backend-go/main.go`.
- **Unified Single-Page Dashboard**: Completely consolidated all workspace features into a single high-priority dashboard in `src/App.jsx`. Features comprehensive interactive form controls, tooltips, real-time telemetry, model routing, mass search/replace, and pipeline dispatching without subpage fragmentation.
- **Unit Tests**: Added test cases for system tray telemetry and shutdown endpoints in `backend-go/main_test.go`.

## [1.0.38] - 2026-06-30
### Fixed
- Merged pending submodules and CI changes.
