# Global TODO

This document tracks granular short-term features, immediate bug fixes, and explicit tasks across the Omni-Workspace.

## Immediate Tasks (Current Sprint)
- [x] Migrate Python scripts into native Go counterparts (e.g. `scripts/task_router.py` -> `backend-go/task_router.go`).
- [x] Ensure all API endpoints in `backend-go` are fully covered by unit tests.
- [x] Add a visual log tailer component to the React app to display `backend-go` logs.
- [ ] Connect the `Shadow Pilot` anomaly report to an explicit interactive form for immediate action.
- [x] Create detailed tooltips across the React UI explaining each submodule status.

## Backlog
- [ ] Automate database backups for `jules.db`.
- [ ] Explore WebSockets for real-time `queue telemetry` instead of HTTP polling.
- [ ] Set up end-to-end testing for the React application using Playwright or Cypress.
