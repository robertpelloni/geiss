# Jules Autopilot Manual

Identity: muse-wonder-1.0 in Muse Code in C:/Users/hyper/workspace/bobsaver on laptop. Heavy processes run at BelowNormal priority.

## Start / stop
- Frontend dev: `bun install; bun run dev` in geiss (Vite on :3006, proxies /api and /system to :8080).
- Backend: `cd geiss/backend-go; $env:CGO_ENABLED='0'; go build -o jules-backend.exe .; .\jules-backend.exe` (serves API + ../dist on :8080). Set priority BelowNormal after start.
- Canonical build: `geiss/scripts/build.ps1` (bun install/build + go vet/build).
- Stop: Quit Server button (dashboard header) or tray Quit; both call the single authenticated shutdown path. One Quit only.

## Auth
All mutating routes require `X-API-KEY: $JULES_API_KEY` and fail closed (401) when the key is unset or wrong. Read-only GETs (status, tray, logs, drift, conflicts, backup list, prune scan) are open.

## Cards -> route -> auth -> effect
| Card | Route | Auth | Effect |
|---|---|---|---|
| Shadow Pilot / Auto-Fix | GET /api/shadow/diff; POST /api/shadow/autofix | POST auth | Auto-fix is simulation-only (501, nothing executed) |
| Task Router | POST /api/tasks/route | auth | Returns assigned_model + reasoning |
| Global Refactor | POST /api/system/refactor | auth | DESTRUCTIVE unless dry_run (default on) |
| State Snapshots | GET/POST /api/system/backup | POST auth | JSON submodule-state snapshot in backend-go/backups/; no restore/delete endpoint |
| Drift / Conflicts | GET /api/system/drift, /api/conflicts/resolve | open | Read-only scans |
| Scan for Broken | GET /api/system/prune | open (GET scan) | Scan only; POST (auth) is the destructive prune |
| Pipeline | POST /api/system/pipeline | auth | Build/test presence checks; deploy simulation-only except local dry-run |
| UI Audit | GET /api/system/audit | auth | Returns submodule, status, backend_files, frontend_files |
| Telemetry | POST /api/system/telemetry | auth | Injects stubs; target_dir must stay inside workspace |
| Commit History | GET /api/system/logs | open | git log output, not application logs |
| Queue | GET /api/queue/telemetry | open | In-memory volatile queue; tasks lost on restart |
| Tray/Server | GET /api/system/tray, /api/system/server/status; POST shutdown routes | POST auth | Single shared start time and shutdown handler |

Ports: 8080 backend, 3006 Vite dev. Version source of truth: root VERSION.md (geiss/VERSION.md and package.json mirror it).
